package password

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc/status"

	grpccodes "google.golang.org/grpc/codes"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
)

type PasswordIdentificationService struct {
	grpc.UnimplementedIAMAuthenticationPasswordServiceServer

	nativeStub *native.NativeStub
	systemStub *system.SystemStub

	globalMongoCollection *mongo.Collection
}

type passwordInMongo struct {
	Identity string `bson:"identity"`
	Salt     []byte `bson:"salt"`
	Password []byte `bson:"password"`
}

func collectionByNamespace(s *PasswordIdentificationService, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.globalMongoCollection
	} else {
		dbName := fmt.Sprintf("openbp_namespace_%s", namespace)
		return s.systemStub.DB.Database(dbName).Collection("native_iam_authentication_password")
	}
}

func NewPasswordIdentificationService(ctx context.Context, systemStub *system.SystemStub, nativeStub *native.NativeStub) (*PasswordIdentificationService, error) {
	err := EnsureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("failed to ensure indexes in global namespace: " + err.Error())
	}

	return &PasswordIdentificationService{
		systemStub:            systemStub,
		globalMongoCollection: systemStub.DB.Database("openbp_global").Collection("native_iam_authentication_password"),
		nativeStub:            nativeStub,
	}, nil
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

	dataToValidate := make([]byte, 0, 32+len(in.Password))
	dataToValidate = append(dataToValidate, entry.Salt...)
	dataToValidate = append(dataToValidate, in.Password...)
	verifyResponse, err := s.systemStub.Vault.HMACVerify(ctx, &vault.HMACVerifyRequest{
		Data:      dataToValidate,
		Signature: entry.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == grpccodes.FailedPrecondition {
				return nil, status.Error(grpccodes.FailedPrecondition, "The vault is sealed. Message from system_vault: "+err.Error())
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to verify password: "+err.Error())
	}

	return &grpc.AuthenticateResponse{Authenticated: verifyResponse.Valid}, status.Error(grpccodes.OK, "")
}

func (s *PasswordIdentificationService) CreateOrUpdate(ctx context.Context, in *grpc.CreateOrUpdateRequest) (*grpc.CreateOrUpdateResponse, error) {
	// Check if namespace exists
	if in.Namespace != "" {
		namespaceExistResponse, err := s.nativeStub.Services.Namespace.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to check if namespace exist "+err.Error())
		}
		if !namespaceExistResponse.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	salt := make([]byte, 32, 32+len(in.Password))
	_, err := rand.Read(salt)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to generate salt: "+err.Error())
	}

	passwordHashResponse, err := s.systemStub.Vault.HMACSign(ctx, &vault.HMACSignRequest{
		Data: append(salt, in.Password...),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == grpccodes.FailedPrecondition {
				return nil, status.Error(grpccodes.FailedPrecondition, "The vault is sealed. Message from system_vault: "+err.Error())
			}
		}

		return nil, status.Error(grpccodes.Internal, "Failed to hash password: "+err.Error())
	}

	collection := collectionByNamespace(s, in.Namespace)
	updateResponse, err := collection.UpdateOne(ctx, bson.M{"identity": in.Identity}, bson.M{"$set": bson.M{"password": passwordHashResponse.Signature, "salt": salt}, "$setOnInsert": bson.M{"identity": in.Identity}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error on updating password in database: "+err.Error())
	}

	// TODO: index on identity field

	return &grpc.CreateOrUpdateResponse{Created: updateResponse.UpsertedCount != 0}, status.Error(grpccodes.OK, "")
}

func (s *PasswordIdentificationService) Delete(ctx context.Context, in *grpc.DeleteRequest) (*grpc.DeleteResponse, error) {
	collection := collectionByNamespace(s, in.Namespace)
	deleteResult, err := collection.DeleteOne(ctx, bson.M{"identity": in.Identity})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.DeleteResponse{}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Error on deleting password in database: "+err.Error())
	}
	return &grpc.DeleteResponse{Existed: deleteResult.DeletedCount != 0}, status.Error(grpccodes.OK, "")
}

func (s *PasswordIdentificationService) Exists(ctx context.Context, in *grpc.ExistsRequest) (*grpc.ExistsResponse, error) {
	collection := collectionByNamespace(s, in.Namespace)
	count, err := collection.CountDocuments(ctx, bson.M{"identity": in.Identity}, options.Count().SetLimit(1))
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &grpc.ExistsResponse{Exists: false}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	return &grpc.ExistsResponse{Exists: count == 1}, status.Error(grpccodes.OK, "")
}
