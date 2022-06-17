package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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
	return fmt.Sprintf("native_namespace_%s", namespaceName)
}

func bsonToNamespace(data bson.M) *grpc.Namespace {
	return &grpc.Namespace{
		Name: data["name"].(string),
	}
}

func (s *NamespaceServer) Ensure(ctx context.Context, in *grpc.EnsureNamespaceRequest) (*grpc.EnsureNamespaceResponse, error) {
	ctx, span := s.tracer.Start(ctx, "ensure")
	defer span.End()

	filterData := bson.D{
		bson.E{Key: "name", Value: in.Name},
	}

	updateData := bson.D{}

	options := options.FindOneAndUpdate().SetUpsert(true)

	_, err := s.namespaceCollection.FindOneAndUpdate(ctx, filterData, updateData, options).DecodeBytes()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY)

	span.SetStatus(codes.Ok, "")
	return &grpc.EnsureNamespaceResponse{Namespace: &grpc.Namespace{Name: in.Name}}, nil
}
func (s *NamespaceServer) Delete(ctx context.Context, in *grpc.DeleteNamespaceRequest) (*grpc.DeleteNamespaceResponse, error) {
	ctx, span := s.tracer.Start(ctx, "delete")
	defer span.End()

	filterData := bson.M{
		"name": in.Name,
	}

	_, err := s.namespaceCollection.DeleteOne(ctx, filterData)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	s.cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY, makeNamespaceDataCacheKey(in.Name))

	err = s.mongoClient.Database(s.dbPrefix + "namespace_" + in.Name).Drop(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	span.SetStatus(codes.Ok, "")
	return &grpc.DeleteNamespaceResponse{}, nil
}

func (s *NamespaceServer) Get(ctx context.Context, in *grpc.GetNamespaceRequest) (*grpc.GetNamespaceResponse, error) {
	ctx, span := s.tracer.Start(ctx, "get")
	defer span.End()

	cacheKey := makeNamespaceDataCacheKey(in.Name)
	if in.UseCache {
		data, err := s.cache.Get(ctx, cacheKey)
		if err != nil && data != nil {
			var result bson.M
			bson.Unmarshal(data, &result)
			response := &grpc.GetNamespaceResponse{
				Namespace: bsonToNamespace(result),
			}
			span.SetStatus(codes.Ok, "")
			return response, nil
		}
	}

	var data bson.M
	err := s.namespaceCollection.FindOne(ctx, bson.M{"name": in.Name}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.SetStatus(codes.Error, "Namespace not found")
			return nil, status.Errorf(grpccodes.NotFound, "Namespace not found")
		}
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(data)
		s.cache.Set(ctx, cacheKey, dataBytes, NAMESPACE_DATA_CACHE_TIMEOUT)
	}

	span.SetStatus(codes.Ok, "")
	response := &grpc.GetNamespaceResponse{
		Namespace: bsonToNamespace(data),
	}
	return response, nil

}

func (s *NamespaceServer) GetAll(in *grpc.GetAllNamespacesRequest, out grpc.NamespaceService_GetAllServer) error {
	ctx, span := s.tracer.Start(out.Context(), "getAll")
	defer span.End()

	if in.UseCache {
		data, err := s.cache.Get(ctx, NAMESPACE_LIST_CACHE_KEY)
		if err != nil && data != nil {
			var result bson.A
			bson.Unmarshal(data, &result)
			for _, namespace := range result {
				err := out.Send(&grpc.GetAllNamespacesResponse{
					Namespace: bsonToNamespace(namespace.(bson.M)),
				})
				if err != nil {
					span.SetStatus(codes.Error, err.Error())
					return status.Errorf(grpccodes.Internal, err.Error())
				}
			}
			span.SetStatus(codes.Ok, "")
			return nil
		}
	}

	cursor, err := s.namespaceCollection.Find(ctx, bson.M{})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return status.Errorf(grpccodes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	var namespaces bson.A
	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return status.Errorf(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&grpc.GetAllNamespacesResponse{Namespace: bsonToNamespace(result)}); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return status.Errorf(grpccodes.Internal, err.Error())
		}
		if in.UseCache {
			namespaces = append(namespaces, result)
		}
	}
	if err := cursor.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return status.Errorf(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(namespaces)
		s.cache.Set(ctx, NAMESPACE_LIST_CACHE_KEY, dataBytes, NAMESPACE_LIST_CACHE_TIMEOUT)
	}

	return nil
}

func (s *NamespaceServer) Exists(ctx context.Context, in *grpc.IsNamespaceExistRequest) (*grpc.IsNamespaceExistResponse, error) {
	ctx, span := s.tracer.Start(ctx, "exists")
	defer span.End()

	if in.UseCache {
		inCache, _ := s.cache.Has(ctx, makeNamespaceDataCacheKey(in.Name))
		if inCache {
			return &grpc.IsNamespaceExistResponse{Exist: true}, nil
		}

		// Get entire namespace data and save it to the cache
		_, err := s.Get(ctx, &grpc.GetNamespaceRequest{Name: in.Name, UseCache: true})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				if e.Code() == grpccodes.NotFound {
					span.SetStatus(codes.Ok, "")
					return &grpc.IsNamespaceExistResponse{Exist: false}, nil
				}
			}
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}

		span.SetStatus(codes.Ok, "")
		return &grpc.IsNamespaceExistResponse{Exist: true}, nil
	} else {
		// Fast check in mongo if namespace exists without getting namespace data
		count, err := s.namespaceCollection.CountDocuments(ctx, bson.M{"name": in.Name}, options.Count().SetLimit(1))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, status.Errorf(grpccodes.Internal, err.Error())
		}

		span.SetStatus(codes.Ok, "")
		return &grpc.IsNamespaceExistResponse{Exist: count == 1}, nil
	}
}
