package balena

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/integrations/balena/api"
	telemetryServer "github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/telemetry"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrSyncInternal = errors.New("internal synchronization error")
var ErrSyncVaultSealed = errors.New("synchronization impossible because vault is sealed")

type syncManager struct {
	systemStub      *system.SystemStub
	telemetryServer *telemetryServer.TelemetryServer
	deviceServer    *device.DeviceServer
	logger          *logrus.Entry
	apiClient       api.Client

	workerContext context.Context
	workerCancel  context.CancelFunc
}

type SyncManager interface {
	SyncAllServers(ctx context.Context) error
	SyncServer(ctx context.Context, server BalenaServerInMongo) (*SyncLogInMongo, error)
	Close()
}

func NewSyncManager(systemStub *system.SystemStub, telemetryServer *telemetryServer.TelemetryServer, deviceServer *device.DeviceServer, logger *logrus.Entry, apiClient api.Client) SyncManager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &syncManager{
		systemStub:      systemStub,
		telemetryServer: telemetryServer,
		deviceServer:    deviceServer,
		logger:          logger,
		apiClient:       apiClient,
		workerContext:   ctx,
		workerCancel:    cancel,
	}

	go manager.worker()
	return manager
}

func (m *syncManager) worker() {
	for {
		err := m.SyncAllServers(m.workerContext)
		if err != nil {
			m.logger.Error("background synchronization failed", err.Error())
		}
		select {
		case <-m.workerContext.Done():
			return
		case <-time.After(time.Second * 10):
			continue
		}
	}
}

func (m *syncManager) Close() {
	m.workerCancel()
}

func (m *syncManager) SyncAllServers(ctx context.Context) error {
	collection := getBalenaServerCollection(m.systemStub)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		err := errors.Join(ErrSyncInternal, errors.New("failed to open database cursor to get all balena server"), err)
		return err
	}

	var wg sync.WaitGroup
	runSyncForServer := func(server BalenaServerInMongo) {
		defer wg.Done()
		if _, err := m.SyncServer(ctx, server); err != nil {
			err := errors.Join(ErrSyncInternal, errors.New("balena server synchronization failed with internal error"), err)
			logger := m.logger.WithFields(logrus.Fields{
				"balena.serverUUID":      server.UUID.Hex(),
				"balena.serverName":      server.Name,
				"balena.serverNamespace": server.Namespace,
			})
			logger.Error(err.Error())
		}
	}

	for cur.Next(ctx) {
		var server BalenaServerInMongo
		if err = cur.Decode(&server); err != nil {
			err := errors.Join(ErrSyncInternal, errors.New("failed to decode server information from database"), err)
			return err
		}

		if !server.Enabled {
			continue
		}

		wg.Add(1)
		go runSyncForServer(server)
	}

	wg.Wait()

	return nil
}

func (s *syncManager) decryptAuthToken(ctx context.Context, encryptedToken []byte) (string, error) {
	decryptResponse, err := s.systemStub.Vault.Decrypt(ctx, &vault.DecryptRequest{
		EncryptedData: encryptedToken,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			return "", errors.Join(ErrSyncVaultSealed, err)
		}

		return "", errors.Join(ErrSyncInternal, errors.New("error with system_vault service while decrypting auth token"), err)
	}

	return string(decryptResponse.PlainData), nil
}

func (m *syncManager) SyncServer(ctx context.Context, server BalenaServerInMongo) (*SyncLogInMongo, error) {
	syncStartTime := time.Now().UTC()

	foundedDevicesOnServer := 0
	foundedActiveDevices := 0
	metricsUpdates := 0

	saveLogAndFormatError := func(logError error) (*SyncLogInMongo, error) {
		logTime := time.Now().UTC()
		syncLogEntry := SyncLogInMongo{
			ServerUUID: server.UUID,
			Timestamp:  logTime,
			Status:     SyncStatusOK,
			Error:      "",
			Stats: SyncStats{
				ExecutionTime:          uint64(logTime.Sub(syncStartTime).Milliseconds()),
				FoundedDevicesOnServer: foundedDevicesOnServer,
				FoundedActiveDevices:   foundedActiveDevices,
				MetricsUpdates:         metricsUpdates,
			},
		}

		if logError != nil {
			if errors.Is(logError, ErrSyncInternal) {
				syncLogEntry.Status = SyncStatusInternalError
				syncLogEntry.Error = "Internal sync error"
			} else {
				syncLogEntry.Status = SyncStatusError
				syncLogEntry.Error = logError.Error()
			}
		}

		logCollection := getSyncLogCollection(m.systemStub)
		insertResult, err := logCollection.InsertOne(ctx, syncLogEntry)
		if err != nil {
			return nil, errors.Join(ErrSyncInternal, err)
		}
		syncLogEntry.UUID = insertResult.InsertedID.(primitive.ObjectID)
		if errors.Is(logError, ErrSyncInternal) {
			return &syncLogEntry, logError
		} else {
			return &syncLogEntry, nil
		}
	}

	authToken, err := m.decryptAuthToken(ctx, server.AuthToken)
	if err != nil {
		return saveLogAndFormatError(err)
	}

	foundedDevices, err := m.apiClient.GetAllDevices(ctx, api.BalenaServerInfo{
		BaseURL:  server.BaseURL,
		APIToken: authToken,
	})
	if err != nil {
		return saveLogAndFormatError(errors.Join(errors.New("failed to get devices from balena server"), err))
	}
	foundedDevicesOnServer = len(foundedDevices)

	devicesCollection := getBalenaDeviceCollection(m.systemStub)
	deviceUpdatesModels := make([]mongo.WriteModel, 0, 10)
	for _, foundedDevice := range foundedDevices {
		updateModel := mongo.NewUpdateOneModel()
		updateModel.SetFilter(bson.M{"balenaServerUUID": server.UUID, "balenaDeviceUUID": foundedDevice.UUID})
		updateModel.SetUpdate(bson.M{
			"$set": bson.M{"balenaData": foundedDevice},
			"$setOnInsert": bson.M{
				"balenaServerNamespace": server.Namespace,
				"balenaServerUUID":      server.UUID,
				"balenaDeviceUUID":      foundedDevice.UUID,
				"created":               syncStartTime,
			},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		})
		updateModel.SetUpsert(true)

		deviceUpdatesModels = append(deviceUpdatesModels, updateModel)
	}
	if _, err = devicesCollection.BulkWrite(ctx, deviceUpdatesModels); err != nil {
		return saveLogAndFormatError(errors.Join(ErrSyncInternal, errors.New("database error when bulkwrite device updates"), err))
	}

	metricsToSubmit := make([]*telemetry.SubmitBasicMetricsRequest_NewBasicMetricData, 0, 10)
	for _, foundedDevice := range foundedDevices {
		if !foundedDevice.IsOnline {
			continue
		}

		//Balena has bug when IsOnline is true, but last connectivity event was long time ago and device is offline
		if foundedDevice.LastConnectivityEvent == nil || syncStartTime.Sub(*foundedDevice.LastConnectivityEvent).Seconds() > 60 {
			continue
		}

		foundedActiveDevices += 1

		// Find if device is binded
		var device BalenaDeviceInMongo
		if err = devicesCollection.FindOne(ctx, bson.M{"balenaServerUUID": server.UUID, "balenaDeviceUUID": foundedDevice.UUID}).Decode(&device); err != nil {
			return saveLogAndFormatError(errors.Join(ErrSyncInternal, errors.New("error while getting device information from database"), err))
		}

		if device.BindedDeviceNamespace == nil || device.BindedDeviceUUID == nil {
			continue
		}

		metricsUpdates += 1

		metricsToSubmit = append(metricsToSubmit, &telemetry.SubmitBasicMetricsRequest_NewBasicMetricData{
			Timestamp:       timestamppb.New(syncStartTime),
			DeviceNamespace: *device.BindedDeviceNamespace,
			DeviceUUID:      (*device.BindedDeviceUUID).Hex(),
			Cpu: &telemetry.CPUMetric{
				Usage: float32(device.BalenaData.CPUUsage),
			},
			Ram: &telemetry.RAMMetric{
				Usage: float32(device.BalenaData.MemoryUsage) / float32(device.BalenaData.MemoryTotal),
			},
		})
	}

	_, err = m.telemetryServer.SubmitBasicMetrics(ctx, &telemetry.SubmitBasicMetricsRequest{
		Metrics: metricsToSubmit,
	})
	if err != nil {
		return saveLogAndFormatError(errors.Join(ErrSyncInternal, errors.New("error while submiting metrics"), err))
	}

	return saveLogAndFormatError(nil)
}
