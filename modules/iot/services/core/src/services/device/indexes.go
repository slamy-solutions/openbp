package device

import (
	"context"
	"errors"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const deviceSearchByIndentity = "fast_identity_search"

func CreateIndexesForNamespace(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	deviceCollection := DeviceCollectionByNamespace(systemStub, namespace)
	_, err := deviceCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "identity", Value: "hashed"},
			},
			Options: options.Index().SetName(deviceSearchByIndentity),
		},
	})
	if err != nil {
		return errors.New("failed to create indexes for the devices collection: " + err.Error())
	}

	return nil
}
