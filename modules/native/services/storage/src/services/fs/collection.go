package fs

import (
	"context"
	"errors"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const fileInfoCollectionPrefix = "native_storage_files_info_"
const fileBucketPrefix = "native_storage_files_data_"
const directoryPrefix = "native_storage_directories_"

func GetFileInfoCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_namespace_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(fileInfoCollectionPrefix + namespace)
}

func GetDirectoryCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_namespace_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(directoryPrefix + namespace)
}

func GetFileBucket(systemStub *system.SystemStub, namespace string, bucketUUID primitive.ObjectID) (*gridfs.Bucket, error) {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_namespace_" + namespace
	}

	bucketName := fileBucketPrefix + namespace + "_" + bucketUUID.Hex()

	bucket, err := gridfs.NewBucket(systemStub.DB.Database(dbName), options.GridFSBucket().SetName(bucketName))
	if err != nil {
		return nil, errors.Join(ErrFailedToGetFileBucket, err)
	}

	return bucket, nil
}

func prepareCollections(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	fileInfoCollection := GetFileInfoCollection(systemStub, namespace)
	_, err := fileInfoCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "bucket", Value: 1}, bson.E{Key: "path", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("bucket_path_unique"),
		},
		{
			Keys:    bson.D{bson.E{Key: "bucket", Value: 1}, bson.E{Key: "baseDirectoryPath", Value: 1}},
			Options: options.Index().SetName("base_directory_path_search"),
		},
	})
	if err != nil {
		err = errors.Join(errors.New("failed to create index for file info collection"), err)
		return err
	}

	directoryCollection := GetDirectoryCollection(systemStub, namespace)
	_, err = directoryCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "bucket", Value: 1}, bson.E{Key: "path", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("bucket_path_unique"),
		},
		{
			Keys:    bson.D{bson.E{Key: "bucket", Value: 1}, bson.E{Key: "baseDirectoryPath", Value: 1}},
			Options: options.Index().SetName("base_directory_path_search"),
		},
	})
	if err != nil {
		err = errors.Join(errors.New("failed to create index for directory collection"), err)
		return err
	}

	return nil
}

func DestroyCollectionsForBucket(ctx context.Context, systemStub *system.SystemStub, namespace string, bucketUUID primitive.ObjectID) error {
	filesBucket, err := GetFileBucket(systemStub, namespace, bucketUUID)
	if err != nil {
		err = errors.Join(errors.New("failed to get file bucket"), err)
		return err
	}

	err = filesBucket.Drop()
	if err != nil {
		err = errors.Join(errors.New("failed to drop file bucket"), err)
		return err
	}

	return nil
}
