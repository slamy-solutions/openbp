package services

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	grpc "github.com/slamy-solutions/open-erp/modules/native/services/namespace/src/grpc/native_namespace"
)

type NamespaceServer struct {
	grpc.UnimplementedNamespaceServiceServer

	dbPrefix            string
	namespaceCollection *mongo.Collection
	mongoClient         *mongo.Client
	cache               cache.Cache
	tracer              trace.Tracer
}

const (
	DATABASE_NAME                = "openerp_global"
	COLLECTION_NAME              = "native_namespace"
	NAMESPACE_LIST_CACHE_KEY     = "native_namespace_list"
	NAMESPACE_DATA_CACHE_TIMEOUT = time.Second * 60
	NAMESPACE_LIST_CACHE_TIMEOUT = time.Second * 60
)

var /* const */ nameValidator = regexp.MustCompile(`^[A-Za-z0-9]+$`)

func New(mongoClient *mongo.Client, cache cache.Cache, dbPrefix string) *NamespaceServer {
	return &NamespaceServer{
		dbPrefix:            dbPrefix,
		namespaceCollection: mongoClient.Database(dbPrefix + "global").Collection("namespace"),
		mongoClient:         mongoClient,
		cache:               cache,
		tracer:              otel.Tracer("github.com/slamy-solutions/open-erp/modules/native/services/namespace"),
	}
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
}

func (n *NamespaceInMongo) ToGRPCNamespace() *grpc.Namespace {
	return &grpc.Namespace{
		Name: n.Name,
	}
}

/*func bsonToNamespace(data bson.M) *grpc.Namespace {
	return &grpc.Namespace{
		Name: data["name"].(string),
	}
}*/

func (s *NamespaceServer) Ensure(ctx context.Context, in *grpc.EnsureNamespaceRequest) (*grpc.EnsureNamespaceResponse, error) {
	if !nameValidator.MatchString(in.Name) {
		return nil, status.Error(grpccodes.InvalidArgument, "Namespace name doesnt match regex \"[A-Za-z0-9]+$\"")
	}

	filterData := bson.D{
		bson.E{Key: "name", Value: in.Name},
	}

	updateData := bson.M{
		"$set": NamespaceInMongo{Name: in.Name},
	}

	options := options.Update().SetUpsert(true)

	_, err := s.namespaceCollection.UpdateOne(ctx, filterData, updateData, options)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY)

	return &grpc.EnsureNamespaceResponse{Namespace: &grpc.Namespace{Name: in.Name}}, status.Error(grpccodes.OK, "")
}
func (s *NamespaceServer) Delete(ctx context.Context, in *grpc.DeleteNamespaceRequest) (*grpc.DeleteNamespaceResponse, error) {
	filterData := bson.M{
		"name": in.Name,
	}

	_, err := s.namespaceCollection.DeleteOne(ctx, filterData)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY, makeNamespaceDataCacheKey(in.Name))

	err = s.mongoClient.Database(s.dbPrefix + "namespace_" + in.Name).Drop(ctx)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &grpc.DeleteNamespaceResponse{}, status.Error(grpccodes.OK, "")
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
