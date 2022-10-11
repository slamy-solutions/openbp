package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/golang/protobuf/proto"
	"github.com/slamy-solutions/openbp/modules/system/libs/go/cache"

	nativeIAmTokenGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/token/src/grpc/native_iam_token"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/token/src/grpc/native_namespace"
)

type IAmTokenServer struct {
	nativeIAmTokenGRPC.UnimplementedIAMTokenServiceServer

	mongoClient           *mongo.Client
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
		db := s.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
		return db.Collection("native_iam_token")
	}
}

func makeTokenCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_token_data_%s_%s", namespace, uuid)
}

func NewIAmTokenServer(mongoClient *mongo.Client, cacheClient cache.Cache, nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient) *IAmTokenServer {
	return &IAmTokenServer{
		mongoClient:           mongoClient,
		mongoGlobalCollection: mongoClient.Database("openbp_global").Collection("native_iam_token"),
		cacheClient:           cacheClient,
		nativeNamespaceClient: nativeNamespaceClient,
	}
}

func (s *IAmTokenServer) Create(ctx context.Context, in *nativeIAmTokenGRPC.CreateRequest) (*nativeIAmTokenGRPC.CreateResponse, error) {
	// Validating if namespace for new token exist
	if in.Namespace != "" {
		namespaceExistResponse, err := s.nativeNamespaceClient.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to check if namespace exist: "+err.Error())
		}
		if !namespaceExistResponse.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
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

func (s *IAmTokenServer) Get(ctx context.Context, in *nativeIAmTokenGRPC.GetRequest) (*nativeIAmTokenGRPC.GetResponse, error) {
	// Trying to fast get data from cache if cache enabled
	var cacheKey string
	if in.UseCache {
		cacheKey = makeTokenCacheKey(in.Namespace, in.Uuid)
		cacheBytes, _ := s.cacheClient.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var tokenData nativeIAmTokenGRPC.TokenData
			err := proto.Unmarshal(cacheBytes, &tokenData)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Error while unmarshaling token from cache. "+err.Error())
			}
			return &nativeIAmTokenGRPC.GetResponse{TokenData: &tokenData}, status.Error(grpccodes.OK, "")
		}
	}

	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Token UUID has bad format")
	}
	collection := collectionByNamespace(s, in.Namespace)
	var mongoToken tokenInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mongoToken)
	if err != nil {
		if err, ok := err.(mongo.ServerError); ok {
			if err.HasErrorCode(73) { // InvalidNamespace
				return nil, status.Error(grpccodes.NotFound, "Token not found. Probably namespace doesnt exist.")
			}
		}
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Token with specified UUID not found.")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get token: "+err.Error())
	}
	tokenData := mongoToken.ToProtoTokenData(in.Namespace)

	if in.UseCache {
		tokenDataBytes, err := proto.Marshal(tokenData)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling token to cache. "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, tokenDataBytes, TOKEN_CACHE_EXPIRATION_TIME)
	}

	return &nativeIAmTokenGRPC.GetResponse{TokenData: tokenData}, status.Error(grpccodes.OK, "")
}

func (s *IAmTokenServer) RawGet(ctx context.Context, in *nativeIAmTokenGRPC.RawGetRequest) (*nativeIAmTokenGRPC.RawGetResponse, error) {
	// Decoding JWT
	jwtData, err := JWTDataFromString(in.Token)
	if err != nil {
		if err == ErrInvalidToken {
			return nil, status.Error(grpccodes.InvalidArgument, "Token invalid. Bad token format or signature.")
		}
		if err == ErrTokenExpired {
			return nil, status.Error(grpccodes.InvalidArgument, "Token expired and must not be used.")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to verify token. "+err.Error())
	}

	getResponse, err := s.Get(ctx, &nativeIAmTokenGRPC.GetRequest{
		Namespace: jwtData.Namespace,
		Uuid:      jwtData.UUID,
		UseCache:  in.UseCache,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() != grpccodes.OK {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &nativeIAmTokenGRPC.RawGetResponse{TokenData: getResponse.TokenData}, status.Error(grpccodes.OK, "")
}

func (s *IAmTokenServer) Delete(ctx context.Context, in *nativeIAmTokenGRPC.DeleteRequest) (*nativeIAmTokenGRPC.DeleteResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Token UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if err, ok := err.(mongo.ServerError); ok {
			if err.HasErrorCode(73) { // InvalidNamespace
				return &nativeIAmTokenGRPC.DeleteResponse{}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to delete token: "+err.Error())
	}

	if result.DeletedCount != 0 {
		s.cacheClient.Remove(ctx, makeTokenCacheKey(in.Namespace, in.Uuid))
	}

	return &nativeIAmTokenGRPC.DeleteResponse{}, status.Error(grpccodes.OK, "")
}
func (s *IAmTokenServer) Disable(ctx context.Context, in *nativeIAmTokenGRPC.DisableRequest) (*nativeIAmTokenGRPC.DisableResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Token UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	result, err := collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"disabled": true}})
	if err != nil {
		if err, ok := err.(mongo.ServerError); ok {
			if err.HasErrorCode(73) { // InvalidNamespace
				return nil, status.Error(grpccodes.NotFound, "Cant find token with specified UUID. Probably namespace doesnt exist.")
			}
		}
		return nil, status.Error(grpccodes.Internal, "failed to update token in database: "+err.Error())
	}

	if result.MatchedCount == 0 {
		return nil, status.Error(grpccodes.NotFound, "Cant find token with specified UUID.")
	}

	if result.ModifiedCount != 0 {
		s.cacheClient.Remove(ctx, makeTokenCacheKey(in.Namespace, in.Uuid))
	}

	return &nativeIAmTokenGRPC.DisableResponse{}, status.Error(grpccodes.OK, "")
}

func (s *IAmTokenServer) Validate(ctx context.Context, in *nativeIAmTokenGRPC.ValidateRequest) (*nativeIAmTokenGRPC.ValidateResponse, error) {
	// Decoding JWT
	jwtData, err := JWTDataFromString(in.Token)
	if err != nil {
		if err == ErrInvalidToken {
			return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_INVALID, TokenData: nil}, status.Error(grpccodes.OK, "")
		}
		if err == ErrTokenExpired {
			return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_EXPIRED, TokenData: nil}, status.Error(grpccodes.OK, "JWT token expired")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to verify token. "+err.Error())
	}

	// Trying to fast get data from cache if cache enabled
	var cacheKey string
	var tokenData nativeIAmTokenGRPC.TokenData
	loadedFromCache := false
	if in.UseCache {
		cacheKey = makeTokenCacheKey(jwtData.Namespace, jwtData.UUID)
		cacheBytes, err := s.cacheClient.Get(ctx, cacheKey)
		if cacheBytes != nil {
			err = proto.Unmarshal(cacheBytes, &tokenData)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Error while unmarshaling token from cache. "+err.Error())
			}
			loadedFromCache = true
		}
	}

	// If cache missed - check database
	if !loadedFromCache {
		id, err := primitive.ObjectIDFromHex(jwtData.UUID)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Token UUID from JWT has bad format")
		}

		collection := collectionByNamespace(s, jwtData.Namespace)
		var data tokenInMongo
		err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&data)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_NOT_FOUND, TokenData: nil}, status.Error(grpccodes.OK, "")
			}
			// Handle error in case if namespaces is not valid (not exist)
			if err, ok := err.(mongo.ServerError); ok {
				if err.HasErrorCode(73) { // InvalidNamespace
					return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_NOT_FOUND, TokenData: nil}, status.Error(grpccodes.OK, "")
				}
			}
			return nil, status.Error(grpccodes.Internal, "Failed to get token from database. "+err.Error())
		}

		tokenData = *data.ToProtoTokenData(jwtData.Namespace)
	}

	if tokenData.Disabled {
		return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_DISABLED, TokenData: nil}, status.Error(grpccodes.OK, "")
	}
	currentTime := time.Now().UTC()
	if tokenData.ExpiresAt.AsTime().Before(currentTime) {
		return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_EXPIRED, TokenData: nil}, status.Error(grpccodes.OK, "Token expired in DB")
	}

	if in.UseCache && !loadedFromCache {
		tokenDataBytes, err := proto.Marshal(&tokenData)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling token to cache. "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, tokenDataBytes, TOKEN_CACHE_EXPIRATION_TIME)
	}

	return &nativeIAmTokenGRPC.ValidateResponse{Status: nativeIAmTokenGRPC.ValidateResponse_OK, TokenData: &tokenData}, status.Error(grpccodes.OK, "")
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
			if err.HasErrorCode(73) { // InvalidNamespace
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

func (s *IAmTokenServer) GetTokensForIdentity(in *nativeIAmTokenGRPC.GetTokensForIdentityRequest, out nativeIAmTokenGRPC.IAMTokenService_GetTokensForIdentityServer) error {
	// TODO: add indexes to the "identity" and "disabled + expiresat" + createdAt(for sorting) fields

	ctx := out.Context()

	collection := collectionByNamespace(s, in.Namespace)
	filter := bson.M{"identity": in.Identity}
	switch in.ActiveFilter {
	case nativeIAmTokenGRPC.GetTokensForIdentityRequest_ONLY_ACTIVE:
		filter["$and"] = bson.A{
			bson.M{"disabled": false},
			bson.M{"expireAt": bson.M{"$gt": time.Now().UTC()}},
		}
		break
	case nativeIAmTokenGRPC.GetTokensForIdentityRequest_ONLY_NOT_ACTIVE:
		filter["$or"] = bson.A{
			bson.M{"disabled": true},
			bson.M{"expireAt": bson.M{"$lt": time.Now().UTC()}},
		}
		break
	}

	findOptions := options.Find().SetSkip(int64(in.Skip)).SetLimit(int64(in.Limit)).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return status.Error(grpccodes.Internal, "Failed to fetch data from the database: "+err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var tokenData tokenInMongo
		if err := cursor.Decode(&tokenData); err != nil {
			return status.Error(grpccodes.Internal, "Failed to fetch token from the database: "+err.Error())
		}
		if err := out.Send(&nativeIAmTokenGRPC.GetTokensForIdentityResponse{TokenData: tokenData.ToProtoTokenData(in.Namespace)}); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
	}
	if err := cursor.Err(); err != nil {
		return status.Error(grpccodes.Internal, "Failed to fetch token from database: "+err.Error())
	}

	return status.Error(grpccodes.OK, "")
}
