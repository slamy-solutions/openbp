package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"

	keyValueStorageGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type KeyValueStorageServer struct {
	keyValueStorageGRPC.UnimplementedKeyValueStorageServiceServer

	mongoClient     *mongo.Client
	cacheClient     cache.Cache
	namespaceClient namespaceGRPC.NamespaceServiceClient
}

type keyInMongo struct {
	Key   string `bson:"key"`
	Value []byte `bson:"value"`
}

const (
	MAX_ENTRY_SIZE            = 1024 * 1024 * 15
	CACHE_KEY_EXPIRATION_TIME = 60 * time.Second
	KEY_INDEX_NAME            = "key_hashed"
)

func makeCacheKey(namespace string, key string) string {
	return fmt.Sprintf("native_keyvaluestorage_key_%s_%s", namespace, key)
}

func getCollectionByNamespace(server *KeyValueStorageServer, namespace string) *mongo.Collection {
	var db *mongo.Database
	if namespace == "" {
		db = server.mongoClient.Database("openbp_global")
	} else {
		db = server.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
	}
	return db.Collection("native_keyvaluestorage")
}

func NewKeyValueStorageServer(mongoClient *mongo.Client, cacheClient cache.Cache, namespaceClient namespaceGRPC.NamespaceServiceClient) *KeyValueStorageServer {
	return &KeyValueStorageServer{
		mongoClient:     mongoClient,
		cacheClient:     cacheClient,
		namespaceClient: namespaceClient,
	}
}

func (s *KeyValueStorageServer) Set(ctx context.Context, in *keyValueStorageGRPC.SetRequest) (*keyValueStorageGRPC.SetResponse, error) {
	if len(in.Value)+len(in.Key) > MAX_ENTRY_SIZE {
		return nil, status.Error(grpccodes.InvalidArgument, "Size of key+value is too big. It must be less than 15 megabytes.")
	}

	if in.Namespace != "" {
		r, err := s.namespaceClient.Exists(ctx, &namespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while checking if namespace exists: "+err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	collection := getCollectionByNamespace(s, in.Namespace)
	updateData := bson.M{"$setOnInsert": bson.M{"key": in.Key}, "$set": bson.M{"value": in.Value}}
	_, err := collection.UpdateOne(ctx, bson.M{"key": in.Key}, updateData, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.cacheClient.Remove(ctx, makeCacheKey(in.Namespace, in.Key))
	return &keyValueStorageGRPC.SetResponse{}, status.Error(grpccodes.OK, "")
}
func (s *KeyValueStorageServer) SetIfNotExist(ctx context.Context, in *keyValueStorageGRPC.SetIfNotExistRequest) (*keyValueStorageGRPC.SetIfNotExistResponse, error) {
	if len(in.Value)+len(in.Key) > MAX_ENTRY_SIZE {
		return nil, status.Error(grpccodes.InvalidArgument, "Size of key+value is too big. It must be less than 15 megabytes.")
	}

	if in.Namespace != "" {
		r, err := s.namespaceClient.Exists(ctx, &namespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while checking if namespace exists: "+err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	collection := getCollectionByNamespace(s, in.Namespace)
	updateData := bson.M{"$setOnInsert": bson.M{"key": in.Key, "value": in.Value}}
	updateResult, err := collection.UpdateOne(ctx, bson.M{"key": in.Key}, updateData, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &keyValueStorageGRPC.SetIfNotExistResponse{Seted: updateResult.UpsertedCount != 0}, status.Error(grpccodes.OK, "")

}
func (s *KeyValueStorageServer) Get(ctx context.Context, in *keyValueStorageGRPC.GetRequest) (*keyValueStorageGRPC.GetResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeCacheKey(in.Namespace, in.Key)
		value, _ := s.cacheClient.Get(ctx, cacheKey)
		if value != nil {
			return &keyValueStorageGRPC.GetResponse{
				Value: value,
			}, status.Error(grpccodes.OK, "")
		}
	}

	var entry keyInMongo
	collection := getCollectionByNamespace(s, in.Namespace)
	err := collection.FindOne(ctx, bson.M{"key": in.Key}, options.FindOne()).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Value for specified namespaces and key wasnt founded")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "Value for specified namespaces and key wasnt founded")
			}
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		s.cacheClient.Set(ctx, cacheKey, entry.Value, CACHE_KEY_EXPIRATION_TIME)
	}

	return &keyValueStorageGRPC.GetResponse{Value: entry.Value}, status.Error(grpccodes.OK, "")
}
func (s *KeyValueStorageServer) Remove(ctx context.Context, in *keyValueStorageGRPC.RemoveRequest) (*keyValueStorageGRPC.RemoveResponse, error) {
	collection := getCollectionByNamespace(s, in.Namespace)
	response, err := collection.DeleteOne(ctx, bson.M{"key": in.Key})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &keyValueStorageGRPC.RemoveResponse{Removed: false}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if response.DeletedCount != 0 {
		s.cacheClient.Remove(ctx, makeCacheKey(in.Namespace, in.Key))
	}

	return &keyValueStorageGRPC.RemoveResponse{Removed: response.DeletedCount != 0}, status.Error(grpccodes.OK, "")
}
func (s *KeyValueStorageServer) Exist(ctx context.Context, in *keyValueStorageGRPC.ExistRequest) (*keyValueStorageGRPC.ExistResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeCacheKey(in.Namespace, in.Key)
		has, err := s.cacheClient.Has(ctx, cacheKey)
		if err == nil && has {
			return &keyValueStorageGRPC.ExistResponse{Exist: true}, status.Error(grpccodes.OK, "")
		}
	}

	if in.UseCache {
		_, err := s.Get(ctx, &keyValueStorageGRPC.GetRequest{Namespace: in.Namespace, Key: in.Key, UseCache: true})
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == grpccodes.NotFound {
					return &keyValueStorageGRPC.ExistResponse{Exist: false}, status.Error(grpccodes.OK, "")
				}
			}
			return nil, err
		}
		return &keyValueStorageGRPC.ExistResponse{Exist: true}, status.Error(grpccodes.OK, "")
	} else {
		collection := getCollectionByNamespace(s, in.Namespace)
		count, err := collection.CountDocuments(ctx, bson.M{"key": in.Key}, options.Count().SetLimit(1))
		if err != nil {
			if err, ok := err.(mongo.WriteException); ok {
				if err.HasErrorLabel("InvalidNamespace") {
					return &keyValueStorageGRPC.ExistResponse{Exist: false}, status.Error(grpccodes.OK, "")
				}
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return &keyValueStorageGRPC.ExistResponse{Exist: count == 1}, status.Error(grpccodes.OK, "")
	}
}
