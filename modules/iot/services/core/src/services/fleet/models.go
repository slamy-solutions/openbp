package fleet

import (
	"time"

	fleetGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FleetInMongo struct {
	UUID        primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"title"`
	Description string             `bson:"description"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

func (t *FleetInMongo) ToGRPCFleet(namespaceName string) *fleetGRPC.Fleet {
	return &fleetGRPC.Fleet{
		Namespace:   namespaceName,
		Uuid:        t.UUID.Hex(),
		Name:        t.Name,
		Description: t.Description,
		Created:     timestamppb.New(t.Created),
		Updated:     timestamppb.New(t.Updated),
		Version:     t.Version,
	}
}

type FleetDeviceInMongo struct {
	DeviceUUID primitive.ObjectID `bson:"deviceUUID"`
	FleetUUID  primitive.ObjectID `bson:"fleetUUID"`

	Added time.Time `bson:"added"`
}
