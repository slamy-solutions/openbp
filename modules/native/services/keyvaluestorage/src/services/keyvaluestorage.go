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

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	keyValueStorageGRPC "github.com/slamy-solutions/open-erp/modules/native/services/keyvaluestorage/src/grpc/native_keyvaluestorage"
	namespaceGRPC "github.com/slamy-solutions/open-erp/modules/native/services/keyvaluestorage/src/grpc/native_namespace"
)

type KeyValueStorageServer struct {
	keyValueStorageGRPC.UnimplementedKeyValueStorageServiceServer

	dbPrefix        string
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
)

func makeCacheKey(namespace string, key string) string {
	return fmt.Sprintf("native_keyvaluestorage_key_%s_%s", namespace, key)
}

func getCollectionByNamespace(server *KeyValueStorageServer, namespace string) *mongo.Collection {
	var db *mongo.Database
	if namespace == "" {
		db = server.mongoClient.Database(fmt.Sprintf("%sglobal", server.dbPrefix))
	} else {
		db = server.mongoClient.Database(fmt.Sprintf("%snamespace_%s", server.dbPrefix, namespace))
	}
	return db.Collection("native_keyvaluestorage")
}

func NewKeyValueStorageServer(dbPrefix string, mongoClient *mongo.Client, cacheClient cache.Cache, namespaceClient namespaceGRPC.NamespaceServiceClient) *KeyValueStorageServer {
	return &KeyValueStorageServer{
		dbPrefix:        dbPrefix,
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
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	// TODO: Index key field
	collection := getCollectionByNamespace(s, in.Namespace)
	updateData := bson.M{"$setOnInsert": bson.M{"key": in.Key}, "$set": bson.M{"value": in.Value}}
	_, err := collection.UpdateOne(ctx, bson.M{"key": in.Key}, updateData, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.cacheClient.Remove(ctx, makeCacheKey(in.Namespace, in.Key))
	return &keyValueStorageGRPC.SetResponse{}, status.Error(grpccodes.OK, "")
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
	err := collection.FindOne(ctx, bson.M{"key": in.Key}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Value for specified namespaces and key wasnt founded")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		s.cacheClient.Set(ctx, cacheKey, entry.Value, CACHE_KEY_EXPIRATION_TIME)
	}

	return &keyValueStorageGRPC.GetResponse{Value: entry.Value}, status.Error(grpccodes.OK, "")
}
