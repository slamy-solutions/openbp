package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/golang/protobuf/proto"
	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	nativeIAmTokenGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/token/src/grpc/native_iam_token"
	nativeNamespaceGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/token/src/grpc/native_namespace"
)

type IAmTokenServer struct {
	nativeIAmTokenGRPC.UnimplementedIAMTokenServiceServer

	mongoClient           *mongo.Client
	mongoDbPrefix         string
	mongoGlobalCollection *mongo.Collection
	cacheClient           cache.Cache
	nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient
}

type scopeInMongo struct {
	Namespace string   `bson:"namespace"`
	Resources []string `bson:"resources"`
	Actions   []string `bson:"actions"`
}

type tokenInMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Identity  string             `bson:"identity"`
	Disabled  bool               `bson:"disabled"`
	Scopes    []scopeInMongo     `bson:"scopes"`
	CreatedAt time.Time          `bson:"createdAt"`
	ExpireAt  time.Time          `bson:"expireAt"`
	Metadata  string             `bson:"metadata"`
}

const (
	TOKEN_EXPITE_TIME           = time.Minute * 15
	REFERSH_TOKEN_EXPIRE_TIME   = time.Hour * 24 * 7
	TOKEN_CACHE_EXPIRATION_TIME = time.Second * 30
)

func (t *tokenInMongo) ToProtoTokenData(namespace string) *nativeIAmTokenGRPC.TokenData {
	scopes := make([]*nativeIAmTokenGRPC.Scope, len(t.Scopes))
	for i, scope := range t.Scopes {
		scopes[i] = &nativeIAmTokenGRPC.Scope{
			Namespace: scope.Namespace,
			Resources: scope.Resources,
			Actions:   scope.Actions,
		}
	}
	return &nativeIAmTokenGRPC.TokenData{
		Namespace:        namespace,
		Uuid:             t.ID.Hex(),
		Identity:         t.Identity,
		Disabled:         t.Disabled,
		ExpiresAt:        timestamppb.New(t.ExpireAt),
		CreatedAt:        timestamppb.New(t.CreatedAt),
		CreationMetadata: t.Metadata,
		Scopes:           scopes,
	}
}

func (t *tokenInMongo) ToJWTData(namespace string, refresh bool, maxExpiration time.Time) *JWTData {
	scopes := []JWTScope{}

	if !refresh {
		scopes = make([]JWTScope, len(t.Scopes))
		for i, scope := range t.Scopes {
			scopes[i] = JWTScope{
				Namespace: scope.Namespace,
				Resources: scope.Resources,
				Actions:   scope.Actions,
			}
		}
	}

	expirationTime := time.Now().UTC().Add(TOKEN_EXPITE_TIME)
	if refresh || expirationTime.After(maxExpiration) {
		expirationTime = maxExpiration
	}
	return NewJWTData(t.ID.Hex(), namespace, t.Identity, scopes, refresh, expirationTime)
}

func collectionByNamespace(s *IAmTokenServer, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.mongoGlobalCollection
	} else {
		db := s.mongoClient.Database(fmt.Sprintf("%snamespace_%s", s.mongoDbPrefix, namespace))
		return db.Collection("native_iam_token")
	}
}

func makeTokenCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_token_data_%s_%s", namespace, uuid)
}

func (s *IAmTokenServer) Create(ctx context.Context, in *nativeIAmTokenGRPC.CreateRequest) (*nativeIAmTokenGRPC.CreateResponse, error) {
	// Validating if namespace for new token exist
	namespaceExistResponse, err := s.nativeNamespaceClient.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to check if namespace exist: "+err.Error())
	}
	if !namespaceExistResponse.Exist {
		return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
	}

	// Converting token to the DB format (bson)
	scopes := make([]scopeInMongo, len(in.Scopes))
	for i, scope := range in.Scopes {
		scopes[i].Namespace = scope.Namespace
		scopes[i].Resources = scope.Resources
		scopes[i].Actions = scope.Actions
	}
	creationTime := time.Now().UTC()
	token := &tokenInMongo{
		Identity:  in.Identity,
		Disabled:  false,
		Scopes:    scopes,
		CreatedAt: creationTime,
		ExpireAt:  creationTime.Add(REFERSH_TOKEN_EXPIRE_TIME),
		Metadata:  in.Metadata,
	}

	// Inserting token to the DB
	collection := collectionByNamespace(s, in.Namespace)
	insertResult, err := collection.InsertOne(ctx, token)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to insert token to DB: "+err.Error())
	}
	token.ID = insertResult.InsertedID.(primitive.ObjectID)

	// Creating token and refresh token
	stringToken, err := token.ToJWTData(in.Namespace, false, token.ExpireAt).ToSignedString()
	if err != nil {
		collection.DeleteOne(ctx, bson.M{"_id": token.ID})
		return nil, status.Error(grpccodes.Internal, "Failed to create token: "+err.Error())
	}
	stringRefreshToken, err := token.ToJWTData(in.Namespace, true, token.ExpireAt).ToSignedString()
	if err != nil {
		collection.DeleteOne(ctx, bson.M{"_id": token.ID})
		return nil, status.Error(grpccodes.Internal, "Failed to create refresh token: "+err.Error())
	}

	return &nativeIAmTokenGRPC.CreateResponse{
		Token:        stringToken,
		RefreshToken: stringRefreshToken,
		TokenData:    token.ToProtoTokenData(in.Namespace),
	}, status.Error(grpccodes.OK, "")
}

func (s *IAmTokenServer) GetByUUID(ctx context.Context, in *nativeIAmTokenGRPC.GetByUUIDRequest) (*nativeIAmTokenGRPC.GetByUUIDResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method GetByUUID not implemented")
}
func (s *IAmTokenServer) DisableByUUID(ctx context.Context, in *nativeIAmTokenGRPC.DisableByUUIDRequest) (*nativeIAmTokenGRPC.DisableByUUIDResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method DisableByUUID not implemented")
}
func (s *IAmTokenServer) Authorize(ctx context.Context, in *nativeIAmTokenGRPC.AuthorizeRequest) (*nativeIAmTokenGRPC.AuthorizeResponse, error) {
	// Decoding JWT
	jwtData, err := JWTDataFromString(in.Token)
	if err != nil {
		if err == ErrInvalidToken {
			return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_INVALID, TokenData: nil}, status.Error(grpccodes.OK, "")
		}
		if err == ErrTokenExpired {
			return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_EXPIRED, TokenData: nil}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to verify token. "+err.Error())
	}

	// Trying to fast get data from cache if cache enabled
	var cacheKey string
	if in.UseCache {
		cacheKey = makeTokenCacheKey(jwtData.Namespace, jwtData.Namespace)
		cacheBytes, err := s.cacheClient.Get(ctx, cacheKey)
		if err == nil {
			var tokenData nativeIAmTokenGRPC.TokenData
			err = proto.Unmarshal(cacheBytes, &tokenData)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Error while unmarshaling token from cache. "+err.Error())
			}
			return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_OK, TokenData: &tokenData}, status.Error(grpccodes.OK, "")
		}
	}

	id, err := primitive.ObjectIDFromHex(jwtData.UUID)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Token UUID from JWT has bad format")
	}

	collection := collectionByNamespace(s, jwtData.Namespace)
	var data tokenInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_NOT_FOUND, TokenData: nil}, status.Error(grpccodes.OK, "")
		}
		// Handle error in case if namespaces is not valid (not exist)
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_NOT_FOUND, TokenData: nil}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get token from database. "+err.Error())
	}

	if data.Disabled {
		return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_DISABLED, TokenData: nil}, status.Error(grpccodes.OK, "")
	}
	currentTime := time.Now().UTC()
	if data.ExpireAt.Before(currentTime) {
		return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_EXPIRED, TokenData: nil}, status.Error(grpccodes.OK, "")
	}

	tokenData := data.ToProtoTokenData(jwtData.Namespace)
	if in.UseCache {
		tokenDataBytes, err := proto.Marshal(tokenData)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling token to cache. "+err.Error())
		}
		timeToExpire := data.ExpireAt.Sub(currentTime)
		if timeToExpire > TOKEN_CACHE_EXPIRATION_TIME {
			timeToExpire = TOKEN_CACHE_EXPIRATION_TIME
		}
		s.cacheClient.Set(ctx, cacheKey, tokenDataBytes, timeToExpire)
	}

	return &nativeIAmTokenGRPC.AuthorizeResponse{Status: nativeIAmTokenGRPC.AuthorizeResponse_OK, TokenData: tokenData}, status.Error(grpccodes.OK, "")
}
func (s *IAmTokenServer) Refresh(ctx context.Context, in *nativeIAmTokenGRPC.RefreshRequest) (*nativeIAmTokenGRPC.RefreshResponse, error) {
	// Decoding JWT
	jwtData, err := JWTDataFromString(in.RefreshToken)
	if err != nil {
		if err == ErrInvalidToken {
			return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_INVALID}, status.Error(grpccodes.OK, "")
		}
		if err == ErrTokenExpired {
			return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_EXPIRED}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to verify token. "+err.Error())
	}
	if !jwtData.Refresh {
		return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_NOT_REFRESH_TOKEN}, status.Error(grpccodes.OK, "")
	}

	id, err := primitive.ObjectIDFromHex(jwtData.UUID)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Token UUID from JWT has bad format")
	}

	collection := collectionByNamespace(s, jwtData.Namespace)
	var data tokenInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_NOT_FOUND}, status.Error(grpccodes.OK, "")
		}
		// Handle error in case if namespaces is not valid (not exist)
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_NOT_FOUND}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get token from database. "+err.Error())
	}

	if data.Disabled {
		return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_DISABLED}, status.Error(grpccodes.OK, "")
	}
	if data.ExpireAt.Before(time.Now().UTC()) {
		return &nativeIAmTokenGRPC.RefreshResponse{Status: nativeIAmTokenGRPC.RefreshResponse_EXPIRED}, status.Error(grpccodes.OK, "")
	}

	tokenString, err := data.ToJWTData(jwtData.Namespace, false, data.ExpireAt).ToSignedString()
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to sign token. "+err.Error())
	}

	return &nativeIAmTokenGRPC.RefreshResponse{
		Status:    nativeIAmTokenGRPC.RefreshResponse_OK,
		Token:     tokenString,
		TokenData: data.ToProtoTokenData(jwtData.Namespace),
	}, status.Error(grpccodes.OK, "")
}
func (s *IAmTokenServer) TokensForIdentity(in *nativeIAmTokenGRPC.TokensForIdentityRequest, out nativeIAmTokenGRPC.IAMTokenService_TokensForIdentityServer) error {
	return status.Errorf(grpccodes.Unimplemented, "method TokensForIdentity not implemented")
}
