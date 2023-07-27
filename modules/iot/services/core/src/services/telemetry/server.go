package telemetry

import (
	"context"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	telemetryGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
	deviceServer "github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TelemetryServer struct {
	telemetryGRPC.UnimplementedTelemetryServiceServer

	logger *logrus.Entry

	systemStub *system.SystemStub
	nativeStub *native.NativeStub

	device *deviceServer.DeviceServer
}

func NewTelemetryServer(ctx context.Context, logger *logrus.Entry, systemStub *system.SystemStub, nativeStub *native.NativeStub, device *deviceServer.DeviceServer) (*TelemetryServer, error) {
	err := CreateIndexesForNamespace(ctx, systemStub, "")
	if err != nil {
		return nil, errors.New("failed to create indexes for global namespace: " + err.Error())
	}

	return &TelemetryServer{
		systemStub: systemStub,
		nativeStub: nativeStub,

		logger: logger,

		device: device,
	}, nil
}

func (s *TelemetryServer) SubmitBasicMetrics(ctx context.Context, in *telemetryGRPC.SubmitBasicMetricsRequest) (*telemetryGRPC.SubmitBasicMetricsResponse, error) {
	metrics := make(map[string]bson.A)
	metricsIndexes := make(map[string][]int)

	for metricIndex, incommingMetric := range in.Metrics {
		deviceUUID, err := primitive.ObjectIDFromHex(incommingMetric.DeviceUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Device UUID [%s] has invalid format.", deviceUUID))
		}

		metric := BasicMetricsInMongo{
			Timestamp:  incommingMetric.Timestamp.AsTime(),
			DeviceUUID: deviceUUID,
			CPU: CPUMetric{
				Usage: incommingMetric.Cpu.Usage,
			},
			RAM: RAMMetric{
				Usage: incommingMetric.Ram.Usage,
			},
			GPU: GPUMetric{
				MemoryUsage: incommingMetric.Gpu.MemoryUsage,
			},
			Network: NetworkMetric{
				Download: incommingMetric.Network.Download,
				Upload:   incommingMetric.Network.Upload,
			},
			Storage: StorageMetric{
				Writes: incommingMetric.Storage.Writes,
				Reads:  incommingMetric.Storage.Reads,
				IOpS:   incommingMetric.Storage.Iops,
			},
		}
		metrics[incommingMetric.DeviceNamespace] = append(metrics[incommingMetric.DeviceNamespace], metric)
		metricsIndexes[incommingMetric.DeviceNamespace] = append(metricsIndexes[incommingMetric.DeviceNamespace], metricIndex)
	}

	metricsUUIDs := make(map[int]string)
	renderAssignedUUIDs := func() []string {
		assignedUUIDs := make([]string, 0, len(in.Metrics))
		for i := 0; i < len(in.Metrics); i++ {
			if uuid, ok := metricsUUIDs[i]; ok {
				assignedUUIDs = append(assignedUUIDs, uuid)
			} else {
				assignedUUIDs = append(assignedUUIDs, "")
			}
		}
		return assignedUUIDs
	}

	for namespaceName, metricsToAdd := range metrics {
		collection := TelemetryBasicMetricCollectionByNamespace(s.systemStub, namespaceName)
		insertResult, err := collection.InsertMany(ctx, metricsToAdd)
		if err != nil {
			err = errors.New("error while inserting basic metrics to the database: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": namespaceName,
			}).Error(err.Error())

			//Return UUIDS that was added
			return &telemetryGRPC.SubmitBasicMetricsResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.Internal, err.Error())
		}
		for i, uuid := range insertResult.InsertedIDs {
			metricsUUIDs[metricsIndexes[namespaceName][i]] = uuid.(primitive.ObjectID).Hex()
		}

		for i, uuid := range insertResult.InsertedIDs {
			metric := metricsToAdd[i].(BasicMetricsInMongo)
			metric.UUID = uuid.(primitive.ObjectID)
			err := RaiseNewBasicMetricsEvent(s.systemStub, metric.ToGRPCBasicMetrics(namespaceName))
			if err != nil {
				err = errors.New("error while raising NewBasicMetric event: " + err.Error())
				s.logger.WithFields(logrus.Fields{
					"namespace": namespaceName,
					"device":    metric.DeviceUUID,
					"metric":    metric.UUID,
				}).Error(err.Error())

				//Return UUIDS that was added
				return &telemetryGRPC.SubmitBasicMetricsResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.Internal, err.Error())
			}
		}
	}

	return &telemetryGRPC.SubmitBasicMetricsResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.OK, "")
}
func (s *TelemetryServer) SubmitLog(ctx context.Context, in *telemetryGRPC.SubmitLogRequest) (*telemetryGRPC.SubmitLogResponse, error) {
	logs := make(map[string]bson.A)
	logsIndexes := make(map[string][]int)

	for entryIndex, incommingEntry := range in.Entries {
		deviceUUID, err := primitive.ObjectIDFromHex(incommingEntry.DeviceUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Device UUID [%s] has invalid format.", deviceUUID))
		}

		logEntry := LogEntryInMongo{
			Timestamp:  incommingEntry.Timestamp.AsTime(),
			DeviceUUID: deviceUUID,
			Message:    incommingEntry.Message,
		}
		logs[incommingEntry.DeviceNamespace] = append(logs[incommingEntry.DeviceNamespace], logEntry)
		logsIndexes[incommingEntry.DeviceNamespace] = append(logsIndexes[incommingEntry.DeviceNamespace], entryIndex)
	}

	logsUUIDs := make(map[int]string)
	renderAssignedUUIDs := func() []string {
		assignedUUIDs := make([]string, 0, len(in.Entries))
		for i := 0; i < len(in.Entries); i++ {
			if uuid, ok := logsUUIDs[i]; ok {
				assignedUUIDs = append(assignedUUIDs, uuid)
			} else {
				assignedUUIDs = append(assignedUUIDs, "")
			}
		}
		return assignedUUIDs
	}

	for namespaceName, entriesToAdd := range logs {
		collection := TelemetryLogCollectionByNamespace(s.systemStub, namespaceName)
		insertResult, err := collection.InsertMany(ctx, entriesToAdd)
		if err != nil {
			err = errors.New("error while inserting logs to the database: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": namespaceName,
			}).Error(err.Error())

			//Return UUIDS that was added
			return &telemetryGRPC.SubmitLogResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.Internal, err.Error())
		}
		for i, uuid := range insertResult.InsertedIDs {
			logsUUIDs[logsIndexes[namespaceName][i]] = uuid.(primitive.ObjectID).Hex()
		}
		for i, uuid := range insertResult.InsertedIDs {
			logEntry := entriesToAdd[i].(LogEntryInMongo)
			logEntry.UUID = uuid.(primitive.ObjectID)
			err := RaiseNewLogEntryEvent(s.systemStub, logEntry.ToGRPCLogEntry(namespaceName))
			if err != nil {
				err = errors.New("error while raising NewLogEntry event: " + err.Error())
				s.logger.WithFields(logrus.Fields{
					"namespace": namespaceName,
					"device":    logEntry.DeviceUUID,
					"logEntry":  logEntry.UUID,
				}).Error(err.Error())

				//Return UUIDS that was added
				return &telemetryGRPC.SubmitLogResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.Internal, err.Error())
			}
		}
	}

	return &telemetryGRPC.SubmitLogResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.OK, "")
}
func (s *TelemetryServer) SubmitEvent(ctx context.Context, in *telemetryGRPC.SubmitEventRequest) (*telemetryGRPC.SubmitEventResponse, error) {
	events := make(map[string]bson.A)
	eventsIndexes := make(map[string][]int)

	for entryIndex, incommingEvent := range in.Events {
		deviceUUID, err := primitive.ObjectIDFromHex(incommingEvent.DeviceUUID)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Device UUID [%s] has invalid format.", deviceUUID))
		}

		event := EventInMongo{
			Timestamp:  incommingEvent.Timestamp.AsTime(),
			DeviceUUID: deviceUUID,
			EventID:    incommingEvent.EventID,
			Data:       incommingEvent.Data,
		}
		if !incommingEvent.Persistent {
			err := RaiseNewDeviceEventEvent(s.systemStub, event.ToGRPCEvent(incommingEvent.DeviceNamespace, incommingEvent.Persistent))
			if err != nil {
				err = errors.New("error while raising NewDeviceEvent for non persistent event: " + err.Error())
				s.logger.WithFields(logrus.Fields{
					"namespace": incommingEvent.DeviceNamespace,
					"device":    incommingEvent.DeviceUUID,
					"eventID":   incommingEvent.EventID,
				}).Error(err.Error())

				//Return UUIDS that was added
				return nil, status.Error(codes.Internal, err.Error())
			}
			continue
		}

		events[incommingEvent.DeviceNamespace] = append(events[incommingEvent.DeviceNamespace], event)
		eventsIndexes[incommingEvent.DeviceNamespace] = append(eventsIndexes[incommingEvent.DeviceNamespace], entryIndex)
	}

	eventsUUIDs := make(map[int]string)
	renderAssignedUUIDs := func() []string {
		assignedUUIDs := make([]string, 0, len(in.Events))
		for i := 0; i < len(in.Events); i++ {
			if uuid, ok := eventsUUIDs[i]; ok {
				assignedUUIDs = append(assignedUUIDs, uuid)
			} else {
				assignedUUIDs = append(assignedUUIDs, "")
			}
		}
		return assignedUUIDs
	}

	for namespaceName, eventsToAdd := range events {
		collection := TelemetryEventCollectionByNamespace(s.systemStub, namespaceName)
		insertResult, err := collection.InsertMany(ctx, eventsToAdd)
		if err != nil {
			err = errors.New("error while inserting events to the database: " + err.Error())
			s.logger.WithFields(logrus.Fields{
				"namespace": namespaceName,
			}).Error(err.Error())

			//Return UUIDS that was added
			return &telemetryGRPC.SubmitEventResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.Internal, err.Error())
		}
		for i, uuid := range insertResult.InsertedIDs {
			eventsUUIDs[eventsIndexes[namespaceName][i]] = uuid.(primitive.ObjectID).Hex()
		}
		for i, uuid := range insertResult.InsertedIDs {
			event := eventsToAdd[i].(EventInMongo)
			event.UUID = uuid.(primitive.ObjectID)
			err := RaiseNewDeviceEventEvent(s.systemStub, event.ToGRPCEvent(namespaceName, true))
			if err != nil {
				err = errors.New("error while raising NewDeviceEvent event: " + err.Error())
				s.logger.WithFields(logrus.Fields{
					"namespace": namespaceName,
					"device":    event.DeviceUUID,
					"event":     event.UUID,
				}).Error(err.Error())

				//Return UUIDS that was added
				return &telemetryGRPC.SubmitEventResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.Internal, err.Error())
			}
		}
	}

	return &telemetryGRPC.SubmitEventResponse{AssignedUUIDs: renderAssignedUUIDs()}, status.Error(codes.OK, "")
}
