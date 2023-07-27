package fleet

import (
	"context"
	"errors"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const fleetDevicesUniqueIndexName = "unique_entries_search"
const fleetDevicesAddedIndex = "added_search"

func CreateIndexesForNamespace(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	fleetDevicesCollection := FleetDevicesCollectionByNamespace(systemStub, namespace)
	_, err := fleetDevicesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "fleetUUID", Value: 1},
				bson.E{Key: "deviceUUID", Value: 1},
			},
			Options: options.Index().SetName(fleetDevicesUniqueIndexName).SetUnique(true),
		},
		{
			Keys: bson.D{
				bson.E{Key: "added", Value: 1},
			},
			Options: options.Index().SetName(fleetDevicesAddedIndex),
		},
	})
	if err != nil {
		return errors.New("failed to create indexes for the fleet devices collection: " + err.Error())
	}

	return nil
}
