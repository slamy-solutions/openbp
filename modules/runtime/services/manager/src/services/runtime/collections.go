package runtime

import (
	"context"
	"errors"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const runtimeCollectionName = "runtime_manager_runtime"
const runtimeDataBucketName = "runtime_manager_runtime_data"

func GetRuntimeCollection(systemStub *system.SystemStub) *mongo.Collection {
	return systemStub.DB.Database("openbp_global").Collection(runtimeCollectionName)
}

func GetRuntimeDataBucket(namespace string, systemStub *system.SystemStub) (*gridfs.Bucket, error) {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "namespace_" + namespace
	}
	return gridfs.NewBucket(systemStub.DB.Database(dbName), options.GridFSBucket().SetName(runtimeDataBucketName))
}

func initCollections(ctx context.Context, systemStub *system.SystemStub) error {
	collection := GetRuntimeCollection(systemStub)

	_, err := collection.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.D{
					bson.E{Key: "namespace", Value: 1},
					bson.E{Key: "name", Value: 1},
				},
				Options: options.Index().SetUnique(true).SetName("unique_within_namespace"),
			},
		},
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create indexes for the runtime collection"), err)
		return err
	}

	return nil
}
