package device

import (
	"time"

	deviceGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DeviceInMongo struct {
	UUID        primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"title"`
	Description string             `bson:"description"`

	Identity string `bson:"identity"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (t *DeviceInMongo) ToGRPCDevice(namespaceName string) *deviceGRPC.Device {
	return &deviceGRPC.Device{
		Namespace:   namespaceName,
		Uuid:        t.UUID.Hex(),
		Name:        t.Name,
		Identity:    t.Identity,
		Description: t.Description,
		Created:     timestamppb.New(t.Created),
		Updated:     timestamppb.New(t.Updated),
		Version:     t.Version,
	}
}
