package balena

import (
	"time"

	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/integrations/balena/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BalenaDeviceInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`

	BindedDeviceNamespace *string             `bson:"bindedDeviceNamespace"`
	BindedDeviceUUID      *primitive.ObjectID `bson:"bindedDeviceUUID"`

	BalenaServerNamespace string             `bson:"balenaServerNamespace"`
	BalenaServerUUID      primitive.ObjectID `bson:"balenaServerUUID"`
	BalenaData            api.Device         `bson:"balenaData"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (d *BalenaDeviceInMongo) ToGRPCDevice() *balena.BalenaDevice {
	device := &balena.BalenaDevice{
		Uuid:                  d.UUID.Hex(),
		BindedDeviceNamespace: "",
		BindedDeviceUUID:      "",
		BalenaServerNamespace: d.BalenaServerNamespace,
		BalenaServerUUID:      d.BalenaServerUUID.Hex(),
		BalenaData:            d.BalenaData.ToGRPCBalenaData(),
		Created:               timestamppb.New(d.Created),
		Updated:               timestamppb.New(d.Updated),
		Version:               d.Version,
	}

	if d.BindedDeviceNamespace != nil {
		device.BindedDeviceNamespace = *d.BindedDeviceNamespace
	}
	if d.BindedDeviceUUID != nil {
		device.BindedDeviceUUID = (*d.BindedDeviceUUID).Hex()
	}

	return device
}

type BalenaServerInMongo struct {
	UUID      primitive.ObjectID `bson:"_id,omitempty"`
	Namespace string             `bson:"namespace"`

	Name        string `bson:"name"`
	Description string `bson:"description"`

	BaseURL string `bson:"baseURL"`
	// Encrypted AuthToken. Has to be decrypted with system_vault before use
	AuthToken []byte `bson:"authToken"`
	Enabled   bool   `bson:"enabled"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (s *BalenaServerInMongo) ToGRPCServer() *balena.BalenaServer {
	return &balena.BalenaServer{
		Namespace:   s.Namespace,
		Uuid:        s.UUID.Hex(),
		Name:        s.Name,
		Description: s.Description,
		BaseURL:     s.BaseURL,
		Enabled:     s.Enabled,
		Created:     timestamppb.New(s.Created),
		Updated:     timestamppb.New(s.Updated),
		Version:     s.Version,
	}
}

type SyncStats struct {
	FoundedDevicesOnServer int `bson:"foundedDevicesOnServer"`
	FoundedActiveDevices   int `bson:"foundedActiveDevices"`
	MetricsUpdates         int `bson:"metricsUpdates"`
	// Execution time in milliseconds
	ExecutionTime uint64 `bson:"executionTime"`
}

type SyncLogInMongo struct {
	UUID       primitive.ObjectID `bson:"_id,omitempty"`
	ServerUUID primitive.ObjectID `bson:"serverUUID"`
	Timestamp  time.Time          `bson:"timestamp"`

	Status string    `bson:"status"`
	Error  string    `bson:"error,omitempty"`
	Stats  SyncStats `bson:"stats"`
}
