package bucket

import (
	"context"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bucketCollectionName = "native_storage_buckets"

func GetBucketsCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_namespace_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(bucketCollectionName)
}

func prepareBucketsCollection(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	collection := GetBucketsCollection(systemStub, namespace)

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"name": 1},
		Options: options.Index().SetName("unique_name").SetUnique(true),
	})
	if err != nil {
		return err
	}

	return nil
}
