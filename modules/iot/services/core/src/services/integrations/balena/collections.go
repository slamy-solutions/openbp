package balena

import (
	"context"
	"errors"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const SyncLogExpirationTime = 60 * 60 * 24 * 7 // 7 days

func getBalenaServerCollection(systemStub *system.SystemStub) *mongo.Collection {
	return systemStub.DB.Database("openbp_global").Collection("iot_integration_balena_server")
}

func getBalenaDeviceCollection(systemStub *system.SystemStub) *mongo.Collection {
	return systemStub.DB.Database("openbp_global").Collection("iot_integration_balena_device")
}

func getSyncLogCollection(systemStub *system.SystemStub) *mongo.Collection {
	return systemStub.DB.Database("openbp_global").Collection("iot_integration_balena_synclog")
}

func setupCollections(ctx context.Context, systemStub *system.SystemStub) error {
	serverCollection := getBalenaServerCollection(systemStub)
	_, err := serverCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "namespace", Value: "hashed"},
			},
			Options: options.Index().SetName("namespace_fast_search"),
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to create indexes for the balena server collection"), err)
	}

	deviceCollection := getBalenaDeviceCollection(systemStub)
	_, err = deviceCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "bindedDeviceNamespace", Value: 1},
				bson.E{Key: "bindedDeviceUUID", Value: "hashed"},
			},
			Options: options.Index().SetName("namespace_fast_binded"),
		},
		{
			Keys: bson.D{
				bson.E{Key: "balenaServerUUID", Value: 1},
				bson.E{Key: "balenaData.id", Value: "hashed"},
			},
			Options: options.Index().SetName("namespace_fast_by_balena_id"),
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to create indexes for the balena server collection"), err)
	}

	syncLogCollectionOptions := options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().
		SetGranularity("hours").
		SetMetaField("serverUUID").
		SetTimeField("timestamp")).
		SetExpireAfterSeconds(SyncLogExpirationTime)
	err = systemStub.DB.Database("openbp_global").CreateCollection(ctx, "iot_integration_balena_synclog", syncLogCollectionOptions)
	if err != nil {
		// If error is not "collection already exist"
		if cmd, ok := err.(mongo.CommandError); !ok || (cmd.Code != 17399 && cmd.Name != "NamespaceExists") {
			return errors.New("failed to create balena synclog collection: " + err.Error())
		}
	}

	return nil
}
