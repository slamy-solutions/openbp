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

	"github.com/slamy-solutions/openbp/modules/system/libs/go/cache"

	lambdaGRPC "github.com/slamy-solutions/openbp/modules/native/services/lambda/manager/src/grpc/native_lambda"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/services/lambda/manager/src/grpc/native_namespace"
)

type LambdaManagerServer struct {
	lambdaGRPC.UnimplementedLambdaManagerServiceServer

	mongoClient           *mongo.Client
	mongoBundleCollection *mongo.Collection
	cacheClient           cache.Cache
	bigCacheClient        cache.Cache
	namespaceClient       namespaceGRPC.NamespaceServiceClient
}

type lambdaInMongo struct {
	Uuid                     string `bson:"uuid,omitempty"`
	Runtime                  string `bson:"runtime,omitempty"`
	Bundle                   []byte `bson:"bundle,omitempty"`
	EnsureExactlyOneDelivery bool   `bson:"ensureExactlyOneDelivery,omitempty"`
}

type lambdaBundleInMongo struct {
	Uuid       []byte `bson:"uuid,omitempty"`
	Data       []byte `bson:"data,omitempty"`
	References int    `bson:"references,omitempty"`
}

func New(ctx context.Context, mongoClient *mongo.Client, cacheClient cache.Cache, bigCacheClient cache.Cache, namespaceClient namespaceGRPC.NamespaceServiceClient) (*LambdaManagerServer, error) {
	dbName := fmt.Sprintf("openbp_global")
	db := mongoClient.Database(dbName)
	server := &LambdaManagerServer{
		mongoClient:           mongoClient,
		mongoBundleCollection: db.Collection("native_lambda_manager_bundle"),
		cacheClient:           cacheClient,
		bigCacheClient:        bigCacheClient,
		namespaceClient:       namespaceClient,
	}

	err := createIndexes(ctx, server)
	return server, err
}

func createIndexes(ctx context.Context, s *LambdaManagerServer) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "uuid", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("unique_uuid"),
	}
	_, err := s.mongoBundleCollection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (s *LambdaManagerServer) getInfoDatabase(namespace string) *mongo.Database {
	return s.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
}

func (s *LambdaManagerServer) getInfoCollection(namespace string) *mongo.Collection {
	return s.getInfoDatabase(namespace).Collection("native_lambda_manager_info")
}

func (s *LambdaManagerServer) Create(ctx context.Context, in *lambdaGRPC.CreateLambdaRequest) (*lambdaGRPC.CreateLambdaResponse, error) {

	// Check if namespace exist.
	existResponse, err := s.namespaceClient.Exists(ctx, &namespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	if !existResponse.Exist {
		return nil, status.Error(grpccodes.FailedPrecondition, "Namespace does not exist")
	}

	// Verify if info collection exist. Collection can not be created in multishard transaction. Create index if not exist
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "uuid", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("unique_uuid"),
	}
	_, err = s.getInfoCollection(in.Namespace).Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	lambdaMongoInfo := lambdaInMongo{
		Uuid:                     in.Uuid,
		Runtime:                  in.Runtime,
		Bundle:                   in.Bundle,
		EnsureExactlyOneDelivery: in.EnsureExactlyOneDelivery,
	}

	lambdaMongoBundle := lambdaBundleInMongo{
		Uuid:       in.Bundle,
		Data:       in.Data,
		References: 0,
	}

	insertTransactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := s.getInfoCollection(in.Namespace).InsertOne(sessCtx, lambdaMongoInfo)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, status.Error(grpccodes.AlreadyExists, "Lambda with same UUID already exists")
			}
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				return nil, err
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		filter := bson.M{"uuid": lambdaMongoBundle.Uuid}
		update := bson.M{"$inc": bson.M{"references": 1}, "$setOnInsert": lambdaMongoBundle}
		_, err = s.mongoBundleCollection.UpdateOne(sessCtx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				return nil, err
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return nil, nil
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	defer session.EndSession(ctx)

	retries := 3
	for {
		retries -= 1
		// Retry transaction several times
		_, err = session.WithTransaction(ctx, insertTransactionCallback)
		if err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && retries >= 0 && cmdErr.HasErrorLabel("TransientTransactionError") {
				continue
			}
			return nil, err
		}
		break
	}

	lambda := &lambdaGRPC.Lambda{
		Namespace:                in.Namespace,
		Uuid:                     in.Uuid,
		Runtime:                  in.Runtime,
		Bundle:                   in.Bundle,
		EnsureExactlyOneDelivery: in.EnsureExactlyOneDelivery,
	}

	return &lambdaGRPC.CreateLambdaResponse{Lambda: lambda}, status.Error(grpccodes.OK, "")
}
func (s *LambdaManagerServer) Delete(ctx context.Context, in *lambdaGRPC.DeleteLambdaRequest) (*lambdaGRPC.DeleteLambdaResponse, error) {
	projection := bson.M{
		"bundle": 1,
	}
	var existingLambda lambdaInMongo
	err := s.getInfoCollection(in.Namespace).FindOne(ctx, bson.M{"uuid": in.Uuid}, options.FindOne().SetProjection(projection)).Decode(&existingLambda)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &lambdaGRPC.DeleteLambdaResponse{}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	deleteTransactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		infoDeleteResult, err := s.getInfoCollection(in.Namespace).DeleteOne(sessCtx, bson.M{"uuid": in.Uuid})
		if err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				return nil, err
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		if infoDeleteResult.DeletedCount == 0 {
			return nil, nil
		}

		_, err = s.mongoBundleCollection.UpdateOne(sessCtx, bson.M{"uuid": existingLambda.Bundle}, bson.M{"$inc": bson.M{"references": -1}})
		if err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				return nil, err
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		_, err = s.mongoBundleCollection.DeleteOne(sessCtx, bson.M{"uuid": existingLambda.Bundle, "references": 0})
		if err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.HasErrorLabel("TransientTransactionError") {
				return nil, err
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return nil, nil
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	defer session.EndSession(ctx)

	retries := 3
	for {
		retries -= 1
		// Retry transaction several times
		_, err = session.WithTransaction(ctx, deleteTransactionCallback)
		if err != nil {
			if cmdErr, ok := err.(mongo.CommandError); ok && retries >= 0 && cmdErr.HasErrorLabel("TransientTransactionError") {
				continue
			}
			return nil, err
		}
		break
	}

	return &lambdaGRPC.DeleteLambdaResponse{}, status.Error(grpccodes.OK, "")
}

func (s *LambdaManagerServer) Exists(ctx context.Context, in *lambdaGRPC.ExistsLambdaRequest) (*lambdaGRPC.ExistsLambdaResponse, error) {
	err := s.getInfoCollection(in.Namespace).FindOne(ctx, bson.M{"uuid": in.Uuid}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &lambdaGRPC.ExistsLambdaResponse{Exists: false}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	return &lambdaGRPC.ExistsLambdaResponse{Exists: true}, status.Error(grpccodes.OK, "")
}

func (s *LambdaManagerServer) Get(ctx context.Context, in *lambdaGRPC.GetLambdaRequest) (*lambdaGRPC.GetLambdaResponse, error) {
	var lambda lambdaInMongo
	err := s.getInfoCollection(in.Namespace).FindOne(ctx, bson.M{"uuid": in.Uuid}).Decode(&lambda)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Bundle with specified Id not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &lambdaGRPC.GetLambdaResponse{
		Lambda: &lambdaGRPC.Lambda{
			Namespace:                in.Namespace,
			Uuid:                     lambda.Uuid,
			Runtime:                  lambda.Runtime,
			Bundle:                   lambda.Bundle,
			EnsureExactlyOneDelivery: lambda.EnsureExactlyOneDelivery,
		},
	}, status.Error(grpccodes.OK, "")
}

func (s *LambdaManagerServer) GetBundle(ctx context.Context, in *lambdaGRPC.GetBundleRequest) (*lambdaGRPC.GetBundleResponse, error) {
	var bundle lambdaBundleInMongo
	projection := bson.M{
		"data": 1,
	}
	err := s.mongoBundleCollection.FindOne(ctx, bson.M{"uuid": in.Bundle}, options.FindOne().SetProjection(projection)).Decode(&bundle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Bundle with specified Id not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &lambdaGRPC.GetBundleResponse{Data: bundle.Data}, status.Error(grpccodes.OK, "")
}
