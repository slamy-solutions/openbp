package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/nats-io/nats.go"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/golang/protobuf/proto"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"

	grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type NamespaceServer struct {
	grpc.UnimplementedNamespaceServiceServer

	namespaceCollection *mongo.Collection
	mongoClient         *mongo.Client
	cache               cache.Cache
	tracer              trace.Tracer
	jetstreamClient     nats.JetStreamContext
}

const (
	DATABASE_NAME                 = "openbp_global"
	COLLECTION_NAME               = "native_namespace"
	NAMESPACE_UNIQUE_INDEX_NAME   = "unique_name"
	NAMESPACE_LIST_CACHE_KEY      = "native_namespace_list"
	NAMESPACE_DATA_CACHE_TIMEOUT  = time.Second * 60
	NAMESPACE_LIST_CACHE_TIMEOUT  = time.Second * 60
	NAMESPACE_STATS_CACHE_KEY     = "native_namespace_stats"
	NAMESPACE_STATS_CACHE_TIMEOUT = time.Second * 60
)

var /* const */ nameValidator = regexp.MustCompile(`^[A-Za-z0-9]+$`)

func New(ctx context.Context, mongoClient *mongo.Client, cache cache.Cache, js nats.JetStreamContext) (*NamespaceServer, error) {

	// Make sure there is stream for events on system_nats
	cfg := nats.StreamConfig{
		Name:      "native_namespace_event",
		Retention: nats.InterestPolicy,
		Subjects:  []string{"native.namespace.event.created", "native.namespace.event.updated", "native.namespace.event.deleted"},
		Storage:   nats.FileStorage,
		Replicas:  1, // TODO: use envirnment variable to enable HA
	}
	_, err := js.AddStream(&cfg)
	if err != nil {
		return nil, err
	}

	// Add namespaces index to ensure unique
	namespaceCollection := mongoClient.Database(DATABASE_NAME).Collection(COLLECTION_NAME)
	_, err = namespaceCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Options: options.Index().SetName(NAMESPACE_UNIQUE_INDEX_NAME).SetUnique(true),
		Keys:    bson.D{bson.E{Key: "name", Value: 1}},
	})
	if err != nil {
		return nil, errors.New("Failed to create index. " + err.Error())
	}

	return &NamespaceServer{
		namespaceCollection: mongoClient.Database(DATABASE_NAME).Collection(COLLECTION_NAME),
		mongoClient:         mongoClient,
		cache:               cache,
		tracer:              otel.Tracer("github.com/slamy-solutions/openbp/modules/native/services/namespace"),
		jetstreamClient:     js,
	}, nil
}

func makeNamespaceDataCacheKey(namespaceName string) string {
	return fmt.Sprintf("native_namespace_data_%s", namespaceName)
}

type NamespacesList struct {
	Namespaces []NamespaceInMongo `bson:"namespaces"`
}

type NamespaceInMongo struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	FullName    string `bson:"fullName"`
	Description string `bson:"description"`

	Updated time.Time `bson:"updated"`
	Created time.Time `bson:"created"`
	Version uint64    `bson:"version"`
}

func (n *NamespaceInMongo) ToGRPCNamespace() *grpc.Namespace {
	return &grpc.Namespace{
		Name:        n.Name,
		FullName:    n.FullName,
		Description: n.Description,
		Created:     timestamppb.New(n.Created),
		Updated:     timestamppb.New(n.Updated),
		Version:     n.Version,
	}
}

/*func bsonToNamespace(data bson.M) *grpc.Namespace {
	return &grpc.Namespace{
		Name: data["name"].(string),
	}
}*/

func validateNamespaceData(name string, fullName string, description string) error {
	if len(name) == 0 {
		return status.Error(grpccodes.InvalidArgument, "Namespace name can not be empty.")
	}
	if !nameValidator.MatchString(name) {
		return status.Error(grpccodes.InvalidArgument, "Namespace name doesnt match regex \"^[A-Za-z0-9]+$\"")
	}
	if len(name) > 32 {
		return status.Error(grpccodes.InvalidArgument, "Namespace name is too long. Max length is 32 symbols.")
	}

	if len(fullName) > 128 {
		return status.Error(grpccodes.InvalidArgument, "Namespace full name is too long. Max length is 128 symbols.")
	}

	if len(description) > 512 {
		return status.Error(grpccodes.InvalidArgument, "Namespace description is too long. Max length is 512 symbols.")
	}

	return nil
}

func (s *NamespaceServer) Ensure(ctx context.Context, in *grpc.EnsureNamespaceRequest) (*grpc.EnsureNamespaceResponse, error) {
	err := validateNamespaceData(in.Name, in.FullName, in.Description)
	if err != nil {
		return nil, err
	}

	filterData := bson.D{
		bson.E{Key: "name", Value: in.Name},
	}

	updateTime := time.Now().UTC().Truncate(time.Millisecond)
	insertData := NamespaceInMongo{
		Name:        in.Name,
		FullName:    in.FullName,
		Description: in.Description,
		Updated:     updateTime,
		Created:     updateTime,
		Version:     0,
	}
	updateData := bson.M{
		"$setOnInsert": insertData,
	}

	options := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updated NamespaceInMongo
	err = s.namespaceCollection.FindOneAndUpdate(ctx, filterData, updateData, options).Decode(&updated)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	namespaceGrpc := updated.ToGRPCNamespace()

	namespaceCreated := updated.Created == updateTime
	if namespaceCreated {
		s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY)
		namespaceBytes, _ := proto.Marshal(namespaceGrpc)
		s.jetstreamClient.Publish("native.namespace.event.created", namespaceBytes)
	}

	return &grpc.EnsureNamespaceResponse{Namespace: namespaceGrpc, Created: namespaceCreated}, status.Error(grpccodes.OK, "")
}

func (s *NamespaceServer) Create(ctx context.Context, in *grpc.CreateNamespaceRequest) (*grpc.CreateNamespaceResponse, error) {
	err := validateNamespaceData(in.Name, in.FullName, in.Description)
	if err != nil {
		return nil, err
	}

	updateTime := time.Now().UTC()
	insertData := NamespaceInMongo{
		Name:        in.Name,
		FullName:    in.FullName,
		Description: in.Description,
		Updated:     updateTime,
		Created:     updateTime,
		Version:     0,
	}

	_, err = s.namespaceCollection.InsertOne(ctx, insertData)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(grpccodes.AlreadyExists, "Namespace with same name already exist.")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to create namespace. "+err.Error())
	}

	// Namespace was created
	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY)
	namespaceGRPC := insertData.ToGRPCNamespace()
	namespaceBytes, _ := proto.Marshal(namespaceGRPC)
	s.jetstreamClient.Publish("native.namespace.event.created", namespaceBytes)

	return &grpc.CreateNamespaceResponse{Namespace: namespaceGRPC}, status.Error(grpccodes.OK, "")
}

func (s *NamespaceServer) Update(ctx context.Context, in *grpc.UpdateNamespaceRequest) (*grpc.UpdateNamespaceResponse, error) {
	err := validateNamespaceData(in.Name, in.FullName, in.Description)
	if err != nil {
		return nil, err
	}

	updateFilter := bson.M{"name": in.Name}
	updateData := bson.M{
		"$set": bson.M{
			"fullName":    in.FullName,
			"description": in.Description,
			"updated":     time.Now().UTC(),
		},
		"$inc": bson.M{
			"version": 1,
		},
	}
	updateOptions := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var newNamespace NamespaceInMongo
	err = s.namespaceCollection.FindOneAndUpdate(ctx, updateFilter, updateData, updateOptions).Decode(&newNamespace)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(grpccodes.NotFound, "Namespace not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error while updating the namespace: "+err.Error())
	}

	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY, makeNamespaceDataCacheKey(in.Name))

	namespaceGRPC := newNamespace.ToGRPCNamespace()
	namespaceBytes, _ := proto.Marshal(namespaceGRPC)
	s.jetstreamClient.Publish("native.namespace.event.updated", namespaceBytes)

	return &grpc.UpdateNamespaceResponse{Namespace: namespaceGRPC}, status.Error(grpccodes.OK, "")
}

func (s *NamespaceServer) Delete(ctx context.Context, in *grpc.DeleteNamespaceRequest) (*grpc.DeleteNamespaceResponse, error) {
	filterData := bson.M{
		"name": in.Name,
	}

	var dataInMongo NamespaceInMongo
	err := s.namespaceCollection.FindOneAndDelete(ctx, filterData).Decode(&dataInMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &grpc.DeleteNamespaceResponse{Existed: false}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY, makeNamespaceDataCacheKey(in.Name))

	err = s.mongoClient.Database("openbp_namespace_" + in.Name).Drop(ctx)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	namespaceBytes, _ := proto.Marshal(dataInMongo.ToGRPCNamespace())
	s.jetstreamClient.Publish("native.namespace.event.deleted", namespaceBytes)

	return &grpc.DeleteNamespaceResponse{Existed: true}, status.Error(grpccodes.OK, "")
}

func (s *NamespaceServer) Get(ctx context.Context, in *grpc.GetNamespaceRequest) (*grpc.GetNamespaceResponse, error) {
	cacheKey := makeNamespaceDataCacheKey(in.Name)
	if in.UseCache {
		data, err := s.cache.Get(ctx, cacheKey)
		if err == nil && data != nil {
			var result NamespaceInMongo
			bson.Unmarshal(data, &result)
			response := &grpc.GetNamespaceResponse{
				Namespace: result.ToGRPCNamespace(),
			}
			return response, status.Error(grpccodes.OK, "")
		}
	}

	var data NamespaceInMongo
	err := s.namespaceCollection.FindOne(ctx, bson.M{"name": in.Name}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(grpccodes.NotFound, "Namespace not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(data)
		s.cache.Set(ctx, cacheKey, dataBytes, NAMESPACE_DATA_CACHE_TIMEOUT)
	}

	response := &grpc.GetNamespaceResponse{
		Namespace: data.ToGRPCNamespace(),
	}
	return response, status.Error(grpccodes.OK, "")

}

func (s *NamespaceServer) GetAll(in *grpc.GetAllNamespacesRequest, out grpc.NamespaceService_GetAllServer) error {
	ctx := out.Context()
	if in.UseCache {
		data, err := s.cache.Get(ctx, NAMESPACE_LIST_CACHE_KEY)
		if err == nil && data != nil {
			var result NamespacesList
			bson.Unmarshal(data, &result)
			for _, namespace := range result.Namespaces {
				err := out.Send(&grpc.GetAllNamespacesResponse{
					Namespace: namespace.ToGRPCNamespace(),
				})
				if err != nil {
					return status.Error(grpccodes.Internal, err.Error())
				}
			}
			return status.Error(grpccodes.OK, "")
		}
	}

	cursor, err := s.namespaceCollection.Find(ctx, bson.M{})
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	var namespaces []NamespaceInMongo
	for cursor.Next(ctx) {
		var result NamespaceInMongo
		if err := cursor.Decode(&result); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&grpc.GetAllNamespacesResponse{Namespace: result.ToGRPCNamespace()}); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
		if in.UseCache {
			namespaces = append(namespaces, result)
		}
	}
	if err := cursor.Err(); err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(NamespacesList{Namespaces: namespaces})
		s.cache.Set(ctx, NAMESPACE_LIST_CACHE_KEY, dataBytes, NAMESPACE_LIST_CACHE_TIMEOUT)
	}

	return nil
}

func (s *NamespaceServer) Exists(ctx context.Context, in *grpc.IsNamespaceExistRequest) (*grpc.IsNamespaceExistResponse, error) {
	if in.UseCache {
		cacheKey := makeNamespaceDataCacheKey(in.Name)
		inCache, _ := s.cache.Has(ctx, cacheKey)
		if inCache {
			return &grpc.IsNamespaceExistResponse{Exist: true}, nil
		}

		// Get entire namespace data and save it to the cache
		var data NamespaceInMongo
		err := s.namespaceCollection.FindOne(ctx, bson.M{"name": in.Name}).Decode(&data)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return &grpc.IsNamespaceExistResponse{Exist: false}, status.Error(grpccodes.OK, "")
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		dataBytes, _ := bson.Marshal(data)
		s.cache.Set(ctx, cacheKey, dataBytes, NAMESPACE_DATA_CACHE_TIMEOUT)

		return &grpc.IsNamespaceExistResponse{Exist: true}, status.Error(grpccodes.OK, "")
	} else {
		// Fast check in mongo if namespace exists without getting namespace data
		count, err := s.namespaceCollection.CountDocuments(ctx, bson.M{"name": in.Name}, options.Count().SetLimit(1))
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return &grpc.IsNamespaceExistResponse{Exist: count == 1}, status.Error(grpccodes.OK, "")
	}
}

func (s *NamespaceServer) Stat(ctx context.Context, in *grpc.GetNamespaceStatisticsRequest) (*grpc.GetNamespaceStatisticsResponse, error) {
	if in.UseCache {
		data, err := s.cache.Get(ctx, NAMESPACE_STATS_CACHE_KEY)
		if err == nil && data != nil {
			var response grpc.GetNamespaceStatisticsResponse
			err = proto.Unmarshal(data, &response)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Failed to unmarshall namespace statistics from cache. "+err.Error())
			}
			return &response, status.Error(grpccodes.OK, "")
		}
	}

	var data bson.M
	err := s.mongoClient.Database(in.Name).RunCommand(ctx, bson.M{
		"dbStats":     1,
		"scale":       1,
		"freeStorage": 0,
	}).Decode(&data)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to get statistics about DB: "+err.Error())
	}

	stats := &grpc.GetNamespaceStatisticsResponse{
		Db: &grpc.GetNamespaceStatisticsResponse_Db{
			DataSize:  data["dataSize"].(uint64),
			Objects:   data["objects"].(uint64),
			TotalSize: data["totalSize"].(uint64),
		},
	}

	if in.UseCache {
		dataBytes, err := bson.Marshal(stats)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall namespace statistics to cache: "+err.Error())
		}
		s.cache.Set(ctx, NAMESPACE_STATS_CACHE_KEY, dataBytes, NAMESPACE_STATS_CACHE_TIMEOUT)
	}

	return stats, status.Error(grpccodes.OK, "")
}
