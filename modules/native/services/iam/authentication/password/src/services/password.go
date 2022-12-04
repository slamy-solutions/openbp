package services

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc/status"

	grpccodes "google.golang.org/grpc/codes"

	grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type PasswordIdentificationService struct {
	grpc.UnimplementedIAMAuthenticationPasswordServiceServer

	mongoClient           *mongo.Client
	globalMongoCollection *mongo.Collection

	nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient
}

type passwordInMongo struct {
	Identity string `bson:"identity"`
	Password []byte `bson:"password"`
}

func collectionByNamespace(s *PasswordIdentificationService, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.globalMongoCollection
	} else {
		dbName := fmt.Sprintf("openbp_namespace_%s", namespace)
		return s.mongoClient.Database(dbName).Collection("native_iam_authentication_password")
	}
}

func NewPasswordIdentificationService(mongoClient *mongo.Client, nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient) *PasswordIdentificationService {
	return &PasswordIdentificationService{
		mongoClient:           mongoClient,
		globalMongoCollection: mongoClient.Database("openbp_global").Collection("native_iam_authentication_password"),
		nativeNamespaceClient: nativeNamespaceClient,
	}
}

func (s *PasswordIdentificationService) Authenticate(ctx context.Context, in *grpc.AuthenticateRequest) (*grpc.AuthenticateResponse, error) {
	collection := collectionByNamespace(s, in.Namespace)
	var entry passwordInMongo
	err := collection.FindOne(ctx, bson.M{"identity": in.Identity}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &grpc.AuthenticateResponse{Authenticated: false}, status.Error(grpccodes.OK, "")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.AuthenticateResponse{Authenticated: false}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to fetch information about identity: "+err.Error())
	}

	err = bcrypt.CompareHashAndPassword(entry.Password, []byte(in.Password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return &grpc.AuthenticateResponse{Authenticated: false}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to compare passwords: "+err.Error())
	}

	return &grpc.AuthenticateResponse{Authenticated: true}, status.Error(grpccodes.OK, "")
}

func (s *PasswordIdentificationService) CreateOrUpdate(ctx context.Context, in *grpc.CreateOrUpdateRequest) (*grpc.CreateOrUpdateResponse, error) {
	// Check if namespace exists
	if in.Namespace != "" {
		namespaceExistResponse, err := s.nativeNamespaceClient.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to check if namespace exist "+err.Error())
		}
		if !namespaceExistResponse.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to hash password: "+err.Error())
	}

	collection := collectionByNamespace(s, in.Namespace)
	_, err = collection.UpdateOne(ctx, bson.M{"identity": in.Identity}, bson.M{"$set": bson.M{"password": passwordHash}, "$setOnInsert": bson.M{"identity": in.Identity}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error on updating password in database: "+err.Error())
	}

	// TODO: index on identity field

	return &grpc.CreateOrUpdateResponse{}, status.Error(grpccodes.OK, "")
}

func (s *PasswordIdentificationService) Delete(ctx context.Context, in *grpc.DeleteRequest) (*grpc.DeleteResponse, error) {
	collection := collectionByNamespace(s, in.Namespace)
	_, err := collection.DeleteOne(ctx, bson.M{"identity": in.Identity})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.DeleteResponse{}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Error on deleting password in database: "+err.Error())
	}
	return &grpc.DeleteResponse{}, status.Error(grpccodes.OK, "")
}

func (s *PasswordIdentificationService) Exist(ctx context.Context, in *grpc.ExistRequest) (*grpc.ExistResponse, error) {
	collection := collectionByNamespace(s, in.Namespace)
	count, err := collection.CountDocuments(ctx, bson.M{"identity": in.Identity}, options.Count().SetLimit(1))
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.ExistResponse{Exist: false}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &grpc.ExistResponse{Exist: count == 1}, status.Error(grpccodes.OK, "")
}
