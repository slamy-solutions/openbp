package services

import (
	"context"
	"fmt"
	grpc "slamy/opencrm/native/catalog/grpc/native_catalog_grpc"
	"time"

	"slamy/opencrm/native/lib/cache"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go.mongodb.org/mongo-driver/"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogServer struct {
	grpc.UnimplementedCatalogServiceServer

	MongoClient *mongo.Client
	Cache       cache.Cache
	Tracer      trace.Tracer
}

const (
	CATALOG_COLLECTION_NAME       = "native_catalog"
	CATALOG_DATA_CACHE_TIMEOUT    = time.Second * 60
	CATALOG_LIST_CACHE_TIMEOUT    = time.Second * 60
	CATALOG_INDEXES_CACHE_TIMEOUT = time.Second * 60
)

func MakeCatalogListCacheKey(namespace string) string {
	return fmt.Sprintf("native_catalog_list_%s", namespace)
}

func MakeCatalogDataCacheKey(namespace string, catalog string) string {
	return fmt.Sprintf("native_catalog_data_%s_%s", namespace, catalog)
}

func MakeCatalogIndexesCacheKey(namescape string, catalog string) string {
	return fmt.Sprintf("native_catalog_indexes_%s_%s", namescape, catalog)
}

func (s *CatalogServer) Create(ctx context.Context, in *grpc.CreateCatalogRequest) (*grpc.CreateCatalogResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "createCatalog")
	defer span.End()

	resultProperties := GRPCPropertiesToBSON(in.Properties)

	creationTime := time.Now().UTC()

	insert := bson.D{
		bson.E{Key: "name", Value: in.Name},
		bson.E{Key: "publicName", Value: in.PublicName},

		bson.E{Key: "properties", Value: resultProperties},

		bson.E{Key: "_created", Value: creationTime},
		bson.E{Key: "_updated", Value: creationTime},
		bson.E{Key: "_version", Value: 1},
	}

	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	// TODO: check if namespace exist

	_, err := db.Collection(CATALOG_COLLECTION_NAME).InsertOne(ctx, insert)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	entryCollectionName := MakeCatalogEntryCollectionName(in.Namespace, in.Name)
	db.CreateCollection(ctx, entryCollectionName)

	s.Cache.Remove(ctx, MakeCatalogListCacheKey(in.Namespace))

	resultCatalog := grpc.Catalog{
		Namespace:  in.Namespace,
		Name:       in.Name,
		PublicName: in.PublicName,
		XCreated:   timestamppb.New(creationTime),
		XUpdated:   timestamppb.New(creationTime),
		XVersion:   1,
		Properties: in.Properties,
	}

	span.SetStatus(codes.Ok, "")
	return &grpc.CreateCatalogResponse{Catalog: &resultCatalog}, nil
}

func (s *CatalogServer) Delete(ctx context.Context, in *grpc.DeleteCatalogRequest) (*grpc.DeleteCatalogResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "deleteCatalog")
	defer span.End()

	deleteFilter := bson.D{
		bson.E{Key: "namespace", Value: in.Namespace},
		bson.E{Key: "name", Value: in.Name},
	}

	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	_, err := db.Collection(CATALOG_COLLECTION_NAME).DeleteOne(ctx, deleteFilter)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.Cache.Remove(ctx, MakeCatalogListCacheKey(in.Namespace), MakeCatalogIndexesCacheKey(in.Namespace, in.Name))
	// TODO: remove all cache for deleted catalog entries

	entryCollectionName := MakeCatalogEntryCollectionName(in.Namespace, in.Name)
	db.Collection(entryCollectionName).Drop(ctx)

	span.SetStatus(codes.Ok, "")
	return &grpc.DeleteCatalogResponse{}, nil
}

func (s *CatalogServer) Update(ctx context.Context, in *grpc.UpdateCatalogRequest) (*grpc.UpdateCatalogReponse, error) {
	ctx, span := s.Tracer.Start(ctx, "updateCatalog")
	defer span.End()

	findFilter := bson.D{
		bson.E{Key: "namespace", Value: in.Namespace},
		bson.E{Key: "name", Value: in.Name},
	}

	updateQuery := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "properties", Value: GRPCPropertiesToBSON(in.Properties)},
			bson.E{Key: "_updated", Value: time.Now().UTC()},
		}},
		bson.E{Key: "$inc", Value: bson.D{
			bson.E{Key: "_version", Value: 1},
		}},
	}

	updateOptions := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var data bson.D
	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	err := db.Collection(CATALOG_COLLECTION_NAME).FindOneAndUpdate(ctx, findFilter, updateQuery, updateOptions).Decode(&data)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.SetStatus(codes.Error, "Catalog not found")
			return nil, status.Errorf(grpccodes.NotFound, "Catalog not found")
		}
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	s.Cache.Remove(ctx, MakeCatalogListCacheKey(in.Namespace), MakeCatalogDataCacheKey(in.Namespace, in.Name))

	updatedCatalog := &grpc.UpdateCatalogReponse{
		Catalog: CatalogFromBSON(data, in.Namespace),
	}

	span.SetStatus(codes.Ok, "")
	return updatedCatalog, nil
}

func (s *CatalogServer) Get(ctx context.Context, in *grpc.GetCatalogRequest) (*grpc.GetCatalogResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "getCatalog")
	defer span.End()

	cacheKey := MakeCatalogDataCacheKey(in.Namespace, in.Name)
	if in.UseCache {
		data, err := s.Cache.Get(ctx, cacheKey)
		if err != nil && data != nil {
			var result bson.D
			bson.Unmarshal(data, &result)
			response := &grpc.GetCatalogResponse{
				Catalog: CatalogFromBSON(result, in.Namespace),
			}
			span.SetStatus(codes.Ok, "")
			return response, nil
		}
	}

	findFilter := bson.D{
		bson.E{Key: "namespace", Value: in.Namespace},
		bson.E{Key: "name", Value: in.Name},
	}

	var data bson.D
	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	err := db.Collection(CATALOG_COLLECTION_NAME).FindOne(ctx, findFilter).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			span.SetStatus(codes.Error, "Catalog not found")
			return nil, status.Errorf(grpccodes.NotFound, "Catalog not found")
		}
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(data)
		s.Cache.Set(ctx, cacheKey, dataBytes, CATALOG_DATA_CACHE_TIMEOUT)
	}

	span.SetStatus(codes.Ok, "")
	response := &grpc.GetCatalogResponse{
		Catalog: CatalogFromBSON(data, in.Namespace),
	}
	return response, nil
}

func (s *CatalogServer) GetIfChanged(ctx context.Context, in *grpc.GetCatalogIfChangedRequest) (*grpc.GetCatalogIfChangedResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "getCatalogIfNotChanged")
	defer span.End()
	response, err := s.Get(ctx, &grpc.GetCatalogRequest{Namespace: in.Namespace, Name: in.Name, UseCache: in.UseCache})
	if err != nil {
		span.SetStatus(codes.Error, "Failed to get catalog")
		return nil, err
	}

	span.SetStatus(codes.Ok, "")

	if response.Catalog.XVersion == in.Version {
		return &grpc.GetCatalogIfChangedResponse{Message: &grpc.GetCatalogIfChangedResponse_Null{}}, nil
	}

	return &grpc.GetCatalogIfChangedResponse{Message: &grpc.GetCatalogIfChangedResponse_Catalog{Catalog: response.Catalog}}, nil
}

func (s *CatalogServer) GetAll(in *grpc.GetAllCatalogsRequest, out grpc.CatalogService_GetAllServer) error {
	ctx, span := s.Tracer.Start(out.Context(), "getAllCatalogs")
	defer span.End()

	cacheKey := MakeCatalogListCacheKey(in.Namespace)
	if in.UseCache {
		data, err := s.Cache.Get(ctx, cacheKey)
		if err != nil && data != nil {
			var result bson.A
			bson.Unmarshal(data, &result)
			for _, catalog := range result {
				err := out.Send(&grpc.GetAllCatalogsResponse{
					Catalog: CatalogFromBSON(catalog.(bson.D), in.Namespace),
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

	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	cursor, err := db.Collection(CATALOG_COLLECTION_NAME).Find(ctx, bson.M{"namespace": in.Namespace})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return status.Errorf(grpccodes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	var catalogs bson.A
	for cursor.Next(ctx) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return status.Errorf(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&grpc.GetAllCatalogsResponse{Catalog: CatalogFromBSON(result, in.Namespace)}); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return status.Errorf(grpccodes.Internal, err.Error())
		}
		if in.UseCache {
			catalogs = append(catalogs, result)
		}
	}
	if err := cursor.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return status.Errorf(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(catalogs)
		s.Cache.Set(ctx, cacheKey, dataBytes, CATALOG_LIST_CACHE_TIMEOUT)
	}

	return nil
}

func (s *CatalogServer) ListIndexes(in *grpc.ListCatalogIndexesRequest, out grpc.CatalogService_ListIndexesServer) error {
	ctx, span := s.Tracer.Start(out.Context(), "listIndexes")
	defer span.End()

	cacheKey := MakeCatalogIndexesCacheKey(in.Namespace, in.Catalog)
	if in.UseCache {
		data, err := s.Cache.Get(ctx, cacheKey)
		if err != nil && data != nil {
			var result bson.A
			bson.Unmarshal(data, &result)
			for _, indexData := range result {
				err := out.Send(&grpc.ListCatalogIndexesResponse{
					Index: CatalogIndexFromBSON(indexData.(bson.D), in.Namespace),
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

	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	cursor, err := db.Collection(MakeCatalogEntryCollectionName(in.Namespace, in.Catalog)).Indexes().List(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return status.Errorf(grpccodes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	var indexes bson.A
	for cursor.Next(ctx) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return status.Errorf(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&grpc.ListCatalogIndexesResponse{Index: CatalogIndexFromBSON(result, in.Namespace)}); err != nil {
			span.SetStatus(codes.Error, err.Error())
			return status.Errorf(grpccodes.Internal, err.Error())
		}
		if in.UseCache {
			indexes = append(indexes, result)
		}
	}
	if err := cursor.Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return status.Errorf(grpccodes.Internal, err.Error())
	}

	if in.UseCache {
		dataBytes, _ := bson.Marshal(indexes)
		s.Cache.Set(ctx, cacheKey, dataBytes, CATALOG_INDEXES_CACHE_TIMEOUT)
	}

	return nil
}
func (s *CatalogServer) EnsureIndex(ctx context.Context, in *grpc.EnsureCatalogIndexRequest) (*grpc.EnsureCatalogIndexResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "ensureIndex")
	defer span.End()

	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	// db.Collection("").Indexes().CreateOne(ctx, mongo.IndexModel{Keys: })
	var keys bson.D
	for _, field := range in.Index.Fields {
		var indexType interface{}
		if field.Type == grpc.CatalogIndex_IndexField_ASCENDING {
			indexType = 1
		} else if field.Type == grpc.CatalogIndex_IndexField_DESCENDING {
			indexType = -1
		} else {
			indexType = "hashed"
		}
		keys = append(keys, bson.E{
			Key:   field.Name,
			Value: indexType,
		})
	}

	mongoIndex := mongo.IndexModel{
		Options: options.Index().SetUnique(in.Index.Unique).SetName(in.Index.Name),
		Keys:    keys,
	}
	_, err := db.Collection(MakeCatalogEntryCollectionName(in.Namespace, in.Catalog)).Indexes().CreateOne(ctx, mongoIndex)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	}

	s.Cache.Remove(ctx, MakeCatalogIndexesCacheKey(in.Namespace, in.Catalog))
	span.SetStatus(codes.Ok, "")
	return &grpc.EnsureCatalogIndexResponse{}, nil
}
func (s *CatalogServer) RemoveIndex(ctx context.Context, in *grpc.RemoveCatalogIndexRequest) (*grpc.RemoveCatalogIndexResponse, error) {
	ctx, span := s.Tracer.Start(ctx, "removeIndex")
	defer span.End()

	db := s.MongoClient.Database(MakeDBName(in.Namespace))
	_, err := db.Collection(MakeCatalogEntryCollectionName(in.Namespace, in.Catalog)).Indexes().DropOne(ctx, in.Index)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, status.Errorf(grpccodes.Internal, err.Error())
	} else {
		s.Cache.Remove(ctx, MakeCatalogIndexesCacheKey(in.Namespace, in.Catalog))
		span.SetStatus(codes.Ok, "")
		return &grpc.RemoveCatalogIndexResponse{}, nil
	}
}
