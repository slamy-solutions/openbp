package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cache "slamy/opencrm/native/lib/cache"

	grpc "slamy/openerp/native/namespace/grpc/native_namespace"
)

type NamespaceServer struct {
	grpc.UnimplementedNamespaceServiceServer

	NamespaceCollection *mongo.Collection
	MongoClient         *mongo.Client
	Cache               cache.Cache
	Tracer              trace.Tracer
}

const (
	DATABASE_NAME                = "openerp_global"
	COLLECTION_NAME              = "native_namespace"
	NAMESPACE_LIST_CACHE_KEY     = "native_namespace_list"
	NAMESPACE_DATA_CACHE_TIMEOUT = time.Second * 60
	NAMESPACE_LIST_CACHE_TIMEOUT = time.Second * 60
)

func makeNamespaceDataCacheKey(namespaceName string) string {
	return fmt.Sprintf("native_namespace_%s", namespaceName)
}

func bsonToNamespace(data bson.M) *grpc.Namespace {
	return &grpc.Namespace{
		Name: data["name"].(string),
	}
}

func (s *NamespaceServer) Ensure(ctx context.Context, in *grpc.EnsureNamespaceRequest) (*grpc.EnsureNamespaceResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "ensure")
	defer span.End()

	filterData := bson.D{
		bson.E{Key: "name", Value: in.Name},
	}

	updateData := bson.D{}

	options := options.FindOneAndUpdate().SetUpsert(true)

	_, err := s.NamespaceCollection.FindOneAndUpdate(ctx, filterData, updateData, options).DecodeBytes()
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	s.Cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY)

	span.SetStatus(codes.Ok, "")
	return &grpc.EnsureNamespaceResponse{Namespace: &grpc.Namespace{Name: in.Name}}, nil
}
func (s *NamespaceServer) Delete(ctx context.Context, in *grpc.DeleteNamespaceRequest) (*grpc.DeleteNamespaceResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "delete")
	defer span.End()

	filterData := bson.M{
		"name": in.Name,
	}

	_, err := s.NamespaceCollection.DeleteOne(ctx, filterData)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	s.Cache.Remove(ctx, NAMESPACE_LIST_CACHE_KEY, makeNamespaceDataCacheKey(in.Name))

	err = s.MongoClient.Database("openerp_namespace_" + in.Name).Drop(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	span.SetStatus(codes.Ok, "")
	return &grpc.DeleteNamespaceResponse{}, nil
}

func (s *NamespaceServer) Get(ctx context.Context, in *grpc.GetNamespaceRequest) (*grpc.GetNamespaceResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "get")
	defer span.End()

	cacheKey := makeNamespaceDataCacheKey(in.Name)
	if in.UseCache {
		data, err := s.Cache.Get(ctx, cacheKey)
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
	err := s.NamespaceCollection.FindOne(ctx, bson.M{"name": in.Name}).Decode(&data)
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
		s.Cache.Set(ctx, cacheKey, dataBytes, NAMESPACE_DATA_CACHE_TIMEOUT)
	}

	span.SetStatus(codes.Ok, "")
	response := &grpc.GetNamespaceResponse{
		Namespace: bsonToNamespace(data),
	}
	return response, nil

}

func (s *NamespaceServer) GetAll(in *grpc.GetAllNamespacesRequest, out grpc.NamespaceService_GetAllServer) error {
	ctx, span := s.Tracer.Start(out.Context(), "getAll")
	defer span.End()

	if in.UseCache {
		data, err := s.Cache.Get(ctx, NAMESPACE_LIST_CACHE_KEY)
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

	cursor, err := s.NamespaceCollection.Find(ctx, bson.M{})
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
		s.Cache.Set(ctx, NAMESPACE_LIST_CACHE_KEY, dataBytes, NAMESPACE_LIST_CACHE_TIMEOUT)
	}

	return nil
}

func (s *NamespaceServer) Exists(ctx context.Context, in *grpc.IsNamespaceExistRequest) (*grpc.IsNamespaceExistResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "exists")
	defer span.End()

	if in.UseCache {
		inCache, _ := s.Cache.Has(ctx, makeNamespaceDataCacheKey(in.Name))
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
		count, err := s.NamespaceCollection.CountDocuments(ctx, bson.M{"name": in.Name}, options.Count().SetLimit(1))
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, status.Errorf(grpccodes.Internal, err.Error())
		}

		span.SetStatus(codes.Ok, "")
		return &grpc.IsNamespaceExistResponse{Exist: count == 1}, nil
	}
}
