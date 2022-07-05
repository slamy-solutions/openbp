package services

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	lambdaGRPC "github.com/slamy-solutions/open-erp/modules/native/services/lambda-manager/src/grpc/native_lambda"
)

type LambdaManagerServer struct {
	lambdaGRPC.UnimplementedLambdaManagerServiceServer

	mongoClient           *mongo.Client
	mongoDatabase         *mongo.Database
	mongoInfoCollection   *mongo.Collection
	mongoBundleCollection *mongo.Collection
	cacheClient           cache.Cache
	bigCacheClient        cache.Cache
}

type lambdaInMongo struct {
	Uuid                     string `bson:"uuid,omitempty"`
	Environment              string `bson:"environment,omitempty"`
	Bundle                   []byte `bson:"bundle,omitempty"`
	EnsureExactlyOneDelivery bool   `bson:"ensureExactlyOneDelivery,omitempty"`
}

type lambdaBundleInMongo struct {
	Uuid       []byte `bson:"uuid,omitempty"`
	Data       []byte `bson:"bundleData,omitempty"`
	References int    `bson:"references,omitempty"`
}

func New(mongoClient *mongo.Client, dbPrefix string, cacheClient cache.Cache, bigCacheClient cache.Cache) *LambdaManagerServer {
	dbName := fmt.Sprintf("%sglobal", dbPrefix)
	db := mongoClient.Database(dbName)
	return &LambdaManagerServer{
		mongoClient:           mongoClient,
		mongoDatabase:         db,
		mongoInfoCollection:   db.Collection("native_lambda_manager_info"),
		mongoBundleCollection: db.Collection("native_lambda_manager_bundle"),
		cacheClient:           cacheClient,
		bigCacheClient:        bigCacheClient,
	}
}

func (s *LambdaManagerServer) Create(ctx context.Context, in *lambdaGRPC.CreateLambdaRequest) (*lambdaGRPC.CreateLambdaResponse, error) {
	lambdaMongoInfo := lambdaInMongo{
		Uuid:                     in.Uuid,
		Environment:              in.Environment,
		Bundle:                   in.Bundle,
		EnsureExactlyOneDelivery: in.EnsureExactlyOneDelivery,
	}

	lambdaMongoBundle := lambdaBundleInMongo{
		Uuid:       in.Bundle,
		Data:       in.Data,
		References: 1,
	}

	insertTransactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := s.mongoInfoCollection.InsertOne(sessCtx, lambdaMongoInfo)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, status.Error(grpccodes.AlreadyExists, "Lambda with same UUID and environment already exists")
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		filter := bson.M{"uuid": lambdaMongoBundle.Uuid}
		update := bson.M{"$inc": bson.M{"references": 1}, "$setOnInsert": lambdaMongoBundle}
		_, err = s.mongoBundleCollection.UpdateOne(sessCtx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return nil, nil
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, insertTransactionCallback)
	if err != nil {
		return nil, err
	}

	lambda := &lambdaGRPC.Lambda{
		Uuid:                     in.Uuid,
		Environment:              in.Environment,
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
	err := s.mongoInfoCollection.FindOne(ctx, bson.M{"uuid": in.Uuid}, options.FindOne().SetProjection(projection)).Decode(&existingLambda)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &lambdaGRPC.DeleteLambdaResponse{}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	deleteTransactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		infoDeleteResult, err := s.mongoInfoCollection.DeleteOne(sessCtx, bson.M{"uuid": in.Uuid})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		if infoDeleteResult.DeletedCount == 0 {
			return nil, nil
		}

		_, err = s.mongoBundleCollection.UpdateOne(sessCtx, bson.M{"uuid": existingLambda.Bundle}, bson.M{"$dec": bson.M{"references": 1}})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		_, err = s.mongoBundleCollection.DeleteOne(sessCtx, bson.M{"uuid": existingLambda.Bundle, "references": 0})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return nil, nil
	}

	session, err := s.mongoClient.StartSession()
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, deleteTransactionCallback)
	if err != nil {
		return nil, err
	}

	return &lambdaGRPC.DeleteLambdaResponse{}, status.Error(grpccodes.OK, "")
}

func (s *LambdaManagerServer) Exists(ctx context.Context, in *lambdaGRPC.ExistsLambdaRequest) (*lambdaGRPC.ExistsLambdaResponse, error) {
	err := s.mongoInfoCollection.FindOne(ctx, bson.M{"uuid": in.Uuid}).Err()
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
	err := s.mongoInfoCollection.FindOne(ctx, bson.M{"uuid": in.Uuid}).Decode(&lambda)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Bundle with specified Id not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &lambdaGRPC.GetLambdaResponse{
		Lambda: &lambdaGRPC.Lambda{
			Uuid:                     lambda.Uuid,
			Environment:              lambda.Environment,
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
