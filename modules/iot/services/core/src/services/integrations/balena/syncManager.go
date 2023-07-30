package balena

import (
	"github.com/sirupsen/logrus"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/integrations/balena/api"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/telemetry"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type syncManager struct {
	systemStub      *system.SystemStub
	telemetryServer *telemetry.TelemetryServer
	deviceServer    *device.DeviceServer
	logger          *logrus.Entry
	apiClient       api.Client
}

type SyncManager interface {
	Close()
}

func NewSyncManager(systemStub *system.SystemStub, telemetryServer *telemetry.TelemetryServer, deviceServer *device.DeviceServer, logger *logrus.Entry, apiClient api.Client) SyncManager {

	return &syncManager{
		systemStub:      systemStub,
		telemetryServer: telemetryServer,
		deviceServer:    deviceServer,
		logger:          logger,
		apiClient:       apiClient,
	}
}

func (m *syncManager) Close() {

}
