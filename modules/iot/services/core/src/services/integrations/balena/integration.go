package balena

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/integrations/balena/api"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/telemetry"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type BalenaIntegration struct {
	ServersServer *ServersServer
	DevicesServer *DevicesServer
	ToolsServer   *ToolsServer
	SyncServer    *SyncServer

	apiClient   api.Client
	syncManager SyncManager
	logger      *logrus.Entry
}

func NewBalenaIntegration(ctx context.Context, systemStub *system.SystemStub, telemetryServer *telemetry.TelemetryServer, deviceServer *device.DeviceServer, logger *logrus.Entry) (*BalenaIntegration, error) {
	err := setupCollections(ctx, systemStub)
	if err != nil {
		return nil, errors.Join(errors.New("failed to setup collections"), err)
	}

	apiClient := api.NewClient()
	syncManager := NewSyncManager(systemStub, telemetryServer, deviceServer, logger.WithField("balenaIntegration.component", "SyncManager"), apiClient)

	serversServer := NewServersServer(logger.WithField("balenaIntegration.component", "grpcServersServer"), systemStub)
	devicesServer := NewDevicesServer(logger.WithField("balenaIntegration.component", "grpcDevicesServer"), systemStub)
	toolsServer := NewToolsServer(logger.WithField("balenaIntegration.component", "grpcToolsServer"))
	syncServer := NewSyncServer(logger.WithField("balenaIntegration.component", "grpcSyncServer"))

	return &BalenaIntegration{
		ServersServer: serversServer,
		DevicesServer: devicesServer,
		ToolsServer:   toolsServer,
		SyncServer:    syncServer,

		apiClient:   apiClient,
		syncManager: syncManager,
		logger:      logger.WithField("balenaIntegration.component", "grpcServer"),
	}, nil
}

func (s *BalenaIntegration) Close() {
	s.apiClient.Close()
	s.syncManager.Close()
}
