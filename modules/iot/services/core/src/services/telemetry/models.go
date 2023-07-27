package telemetry

import (
	"time"

	telemetryGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CPUMetric struct {
	Usage float32 `bson:"usage"`
}

func (m *CPUMetric) ToGRPCMetric() *telemetryGRPC.CPUMetric {
	return &telemetryGRPC.CPUMetric{
		Usage: m.Usage,
	}
}

type RAMMetric struct {
	Usage float32 `bson:"usage"`
}

func (m *RAMMetric) ToGRPCMetric() *telemetryGRPC.RAMMetric {
	return &telemetryGRPC.RAMMetric{
		Usage: m.Usage,
	}
}

type GPUMetric struct {
	MemoryUsage float32 `bson:"memoryUsage"`
}

func (m *GPUMetric) ToGRPCMetric() *telemetryGRPC.GPUMetric {
	return &telemetryGRPC.GPUMetric{
		MemoryUsage: m.MemoryUsage,
	}
}

type NetworkMetric struct {
	Download uint64 `bson:"download"`
	Upload   uint64 `bson:"upload"`
}

func (m *NetworkMetric) ToGRPCMetric() *telemetryGRPC.NetworkMetric {
	return &telemetryGRPC.NetworkMetric{
		Download: m.Download,
		Upload:   m.Upload,
	}
}

type StorageMetric struct {
	Writes uint64 `bson:"writes"`
	Reads  uint64 `bson:"reads"`
	IOpS   uint64 `bson:"iops"`
}

func (m *StorageMetric) ToGRPCMetric() *telemetryGRPC.StorageMetric {
	return &telemetryGRPC.StorageMetric{
		Writes: m.Writes,
		Reads:  m.Reads,
		Iops:   m.IOpS,
	}
}

type BasicMetricsInMongo struct {
	UUID       primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp  time.Time          `bson:"timestamp"`
	DeviceUUID primitive.ObjectID `bson:"deviceUUID"`

	CPU     CPUMetric     `bson:"cpu"`
	RAM     RAMMetric     `bson:"ram"`
	GPU     GPUMetric     `bson:"gpu"`
	Network NetworkMetric `bson:"network"`
	Storage StorageMetric `bson:"storage"`
}

func (m *BasicMetricsInMongo) ToGRPCBasicMetrics(namespace string) *telemetryGRPC.BasicMetrics {
	return &telemetryGRPC.BasicMetrics{
		UUID:            m.UUID.Hex(),
		Timestamp:       timestamppb.New(m.Timestamp),
		DeviceNamespace: namespace,
		DeviceUUID:      m.DeviceUUID.Hex(),
		Cpu:             m.CPU.ToGRPCMetric(),
		Ram:             m.RAM.ToGRPCMetric(),
		Gpu:             m.GPU.ToGRPCMetric(),
		Network:         m.Network.ToGRPCMetric(),
		Storage:         m.Storage.ToGRPCMetric(),
	}
}

type LogEntryInMongo struct {
	UUID       primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp  time.Time          `bson:"timestamp"`
	DeviceUUID primitive.ObjectID `bson:"deviceUUID"`
	Message    []byte             `bson:"message"`
}

func (e *LogEntryInMongo) ToGRPCLogEntry(namespace string) *telemetryGRPC.LogEntry {
	return &telemetryGRPC.LogEntry{
		UUID:            e.UUID.Hex(),
		Timestamp:       timestamppb.New(e.Timestamp),
		DeviceNamespace: namespace,
		DeviceUUID:      e.DeviceUUID.Hex(),
		Message:         e.Message,
	}
}

type EventInMongo struct {
	UUID       primitive.ObjectID `bson:"_id,omitempty"`
	Timestamp  time.Time          `bson:"timestamp"`
	DeviceUUID primitive.ObjectID `bson:"deviceUUID"`

	EventID string `bson:"eventID"`
	Data    []byte `bson:"data"`
}

func (e *EventInMongo) ToGRPCEvent(namespace string, persistent bool) *telemetryGRPC.Event {
	return &telemetryGRPC.Event{
		UUID:            e.UUID.Hex(),
		Timestamp:       timestamppb.New(e.Timestamp),
		DeviceNamespace: namespace,
		DeviceUUID:      e.DeviceUUID.Hex(),
		EventID:         e.EventID,
		Data:            e.Data,
		Persistent:      persistent,
	}
}
