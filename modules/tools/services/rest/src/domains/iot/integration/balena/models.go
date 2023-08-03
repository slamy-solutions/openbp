package balena

import (
	"time"

	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
)

type formatedServer struct {
	UUID        string    `json:"uuid"`
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Version     uint64    `json:"version"`
}

func formatedServerFromGRPC(server *balena.BalenaServer) formatedServer {
	return formatedServer{
		UUID:        server.Uuid,
		Namespace:   server.Namespace,
		Name:        server.Name,
		Description: server.Description,
		Enabled:     server.Enabled,
		Created:     server.Created.AsTime(),
		Updated:     server.Updated.AsTime(),
		Version:     server.Version,
	}
}

type formatedServerWithMetadata struct {
	Server      formatedServer        `json:"server"`
	LastSyncLog *formatedSyncLogEntry `json:"lastSyncLog"`
}

type formatedBalenaData struct {
	UUID                  string    `json:"uuid"`
	Id                    int32     `json:"id"`
	IsOnline              bool      `json:"isOnline"`
	Status                string    ` json:"status"`
	DeviceName            string    ` json:"deviceName"`
	Longitude             string    `json:"longitude"`
	Latitude              string    ` json:"latitude"`
	Location              string    `json:"location"`
	LastConnectivityEvent time.Time ` json:"lastConnectivityEvent"`
	MemoryUsage           uint32    `json:"memoryUsage"`
	MemoryTotal           uint32    `json:"memoryTotal"`
	StorageUsage          uint32    ` json:"storageUsage"`
	CpuUsage              uint32    ` json:"cpuUsage"`
	CpuTemp               uint32    ` json:"cpuTemp"`
	IsUndervolted         bool      ` json:"isUndervolted"`
}

func formatedBalenaDataFromGRPC(data *balena.BalenaData) formatedBalenaData {
	return formatedBalenaData{
		UUID:                  data.Uuid,
		Id:                    data.Id,
		IsOnline:              data.IsOnline,
		Status:                data.Status,
		DeviceName:            data.DeviceName,
		Longitude:             data.Longitude,
		Latitude:              data.Latitude,
		Location:              data.Location,
		LastConnectivityEvent: data.LastConnectivityEvent.AsTime(),
		MemoryUsage:           data.MemoryUsage,
		MemoryTotal:           data.MemoryTotal,
		StorageUsage:          data.StorageUsage,
		CpuUsage:              data.CpuUsage,
		CpuTemp:               data.CpuTemp,
		IsUndervolted:         data.IsUndervolted,
	}
}

type formatedDevice struct {
	UUID string `json:"uuid"`

	BindedDeviceNamespace string `json:"bindedDeviceNamespace"`
	BindedDeviceUUID      string `json:"bindedDeviceUUID"`

	BalenaServerNamespace string             `json:"balenaServerNamespace"`
	BalenaServerUUID      string             `json:"balenaServerUUID"`
	BalenaData            formatedBalenaData `json:"balenaData"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	Version uint64    `json:"version"`
}

func formatedDeviceFromGRPC(device *balena.BalenaDevice) formatedDevice {
	return formatedDevice{
		UUID:                  device.Uuid,
		BindedDeviceNamespace: device.BindedDeviceNamespace,
		BindedDeviceUUID:      device.BindedDeviceUUID,
		BalenaServerNamespace: device.BalenaServerNamespace,
		BalenaServerUUID:      device.BalenaServerUUID,
		BalenaData:            formatedBalenaDataFromGRPC(device.BalenaData),
		Created:               device.Created.AsTime(),
		Updated:               device.Updated.AsTime(),
		Version:               device.Version,
	}
}

type formatedSyncLogEntryStats struct {
	FoundedDevicesOnServer int32  `json:"foundedDevicesOnServer"`
	FoundedActiveDevices   int32  `json:"foundedActiveDevices"`
	MetricsUpdates         int32  `json:"metricsUpdates"`
	ExecutionTime          uint64 `json:"executionTime"`
}

func formatedSyncLogEntryStatsFromGRPC(stats *balena.SyncLogEntry_Stats) formatedSyncLogEntryStats {
	return formatedSyncLogEntryStats{
		FoundedDevicesOnServer: stats.FoundedDevicesOnServer,
		FoundedActiveDevices:   stats.FoundedActiveDevices,
		MetricsUpdates:         stats.MetricsUpdates,
		ExecutionTime:          stats.ExecutionTime,
	}
}

type formatedSyncLogEntry struct {
	UUID       string                    `json:"uuid"`
	ServerUUID string                    `json:"serverUUID"`
	Timestamp  time.Time                 `json:"timestamp"`
	Status     string                    `json:"status"`
	Error      string                    `json:"error"`
	Stats      formatedSyncLogEntryStats `json:"stats"`
}

func formatedSyncLogEntryFromGRPC(entry *balena.SyncLogEntry) formatedSyncLogEntry {
	status := "OK"
	if entry.Status == balena.SyncLogEntry_ERROR {
		status = "ERROR"
	} else if entry.Status == balena.SyncLogEntry_INTERNAL_ERROR {
		status = "INTERNAL_ERROR"
	}

	return formatedSyncLogEntry{
		UUID:       entry.Uuid,
		ServerUUID: entry.ServerUUID,
		Timestamp:  entry.Timestamp.AsTime(),
		Status:     status,
		Error:      entry.Error,
		Stats:      formatedSyncLogEntryStatsFromGRPC(entry.Stats),
	}
}

/*
	// Number of devices founded on the balena server
	FoundedDevicesOnServer int32 `protobuf:"varint,1,opt,name=foundedDevicesOnServer,proto3" json:"foundedDevicesOnServer,omitempty"`
	// Number of founded devices that are active (online)
	FoundedActiveDevices int32 `protobuf:"varint,2,opt,name=foundedActiveDevices,proto3" json:"foundedActiveDevices,omitempty"`
	// Number of the devices for which metrics where updated
	MetricsUpdates int32 `protobuf:"varint,3,opt,name=metricsUpdates,proto3" json:"metricsUpdates,omitempty"`
	// Execution time in milliseconds
	ExecutionTime uint64 `protobuf:"varint,4,opt,name=executionTime,proto3" json:"executionTime,omitempty"`

*/
