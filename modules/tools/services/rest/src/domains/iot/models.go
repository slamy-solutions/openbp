package iot

import (
	"time"

	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
)

type formatedDevice struct {
	Namespace   string    `json:"namespace"`
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Identity    string    `json:"identity"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Version     uint64    `json:"version"`
}

func formatedDeviceFromGRPC(grpcDevice *device.Device) formatedDevice {
	return formatedDevice{
		Namespace:   grpcDevice.Namespace,
		UUID:        grpcDevice.Uuid,
		Name:        grpcDevice.Name,
		Description: grpcDevice.Description,
		Identity:    grpcDevice.Identity,
		Created:     grpcDevice.Created.AsTime(),
		Updated:     grpcDevice.Updated.AsTime(),
		Version:     grpcDevice.Version,
	}
}

type formatedFleet struct {
	Namespace   string `json:"namespace"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Version uint64    `json:"version"`
}

func formatedFleetFromGRPC(grpcFleet *fleet.Fleet) formatedFleet {
	return formatedFleet{
		Namespace:   grpcFleet.Namespace,
		UUID:        grpcFleet.Uuid,
		Name:        grpcFleet.Name,
		Description: grpcFleet.Description,
		Created:     grpcFleet.Created.AsTime(),
		Updated:     grpcFleet.Updated.AsTime(),
		Version:     grpcFleet.Version,
	}
}

type formatedTelemetryBasicMetrics struct {
	DeviceNamespace string    `json:"deviceNamespace"`
	UUID            string    `json:"uuid"`
	Timestamp       time.Time `json:"timestamp"`
	DeviceUUID      string    `json:"deviceUUID"`

	CPU struct {
		Usage float32 `json:"usage"`
	} `json:"cpu"`

	RAM struct {
		Usage float32 `json:"usage"`
	} `json:"ram"`

	GPU struct {
		MemoryUsage float32 `json:"memoryUsage"`
	} `json:"gpu"`

	Network struct {
		Download uint64 `json:"download"`
		Upload   uint64 `json:"upload"`
	} `json:"network"`

	Storage struct {
		Writes uint64 `json:"writes"`
		Reads  uint64 `json:"reads"`
		IOpS   uint64 `json:"iops"`
	} `json:"storage"`
}

func formatedTelemetryBasicMetricsFromGRPC(grpcBasicMetrics *telemetry.BasicMetrics) formatedTelemetryBasicMetrics {
	return formatedTelemetryBasicMetrics{
		DeviceNamespace: grpcBasicMetrics.DeviceNamespace,
		DeviceUUID:      grpcBasicMetrics.DeviceUUID,
		UUID:            grpcBasicMetrics.UUID,
		Timestamp:       grpcBasicMetrics.Timestamp.AsTime(),

		CPU: struct {
			Usage float32 "json:\"usage\""
		}{
			Usage: grpcBasicMetrics.Cpu.Usage,
		},

		RAM: struct {
			Usage float32 "json:\"usage\""
		}{
			Usage: grpcBasicMetrics.Ram.Usage,
		},

		GPU: struct {
			MemoryUsage float32 "json:\"memoryUsage\""
		}{
			MemoryUsage: grpcBasicMetrics.Gpu.MemoryUsage,
		},

		Network: struct {
			Download uint64 "json:\"download\""
			Upload   uint64 "json:\"upload\""
		}{
			Download: grpcBasicMetrics.Network.Download,
			Upload:   grpcBasicMetrics.Network.Upload,
		},

		Storage: struct {
			Writes uint64 "json:\"writes\""
			Reads  uint64 "json:\"reads\""
			IOpS   uint64 "json:\"iops\""
		}{
			Writes: grpcBasicMetrics.Storage.Writes,
			Reads:  grpcBasicMetrics.Storage.Reads,
			IOpS:   grpcBasicMetrics.Storage.Iops,
		},
	}
}

type formatedTelemetryLogEntry struct {
	DeviceNamespace string    `json:"deviceNamespace"`
	UUID            string    `json:"uuid"`
	Timestamp       time.Time `json:"timestamp"`
	DeviceUUID      string    `json:"deviceUUID"`
	Message         []byte    `json:"message"`
}

func formatedTelemetryLogEntryFromGRPC(grpcLogEntry *telemetry.LogEntry) formatedTelemetryLogEntry {
	return formatedTelemetryLogEntry{
		DeviceNamespace: grpcLogEntry.DeviceNamespace,
		UUID:            grpcLogEntry.UUID,
		DeviceUUID:      grpcLogEntry.DeviceUUID,
		Timestamp:       grpcLogEntry.Timestamp.AsTime(),
		Message:         grpcLogEntry.Message,
	}
}

type formatedTelemetryEvent struct {
	UUID       string    `json:"uuid"`
	Timestamp  time.Time `json:"timestamp"`
	DeviceUUID string    `json:"deviceUUID"`

	EventID string `json:"eventID"`
	Data    []byte `json:"data"`
}

func formatedTelemetryEventFromGRPC(grpcEvent *telemetry.Event) formatedTelemetryEvent {
	return formatedTelemetryEvent{
		UUID:       grpcEvent.UUID,
		Timestamp:  grpcEvent.Timestamp.AsTime(),
		DeviceUUID: grpcEvent.DeviceUUID,
		EventID:    grpcEvent.EventID,
		Data:       grpcEvent.Data,
	}
}
