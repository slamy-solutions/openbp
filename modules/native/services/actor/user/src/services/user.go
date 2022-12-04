package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/golang/protobuf/proto"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"

	nativeActorUserGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
)

type ActorUserServer struct {
	nativeActorUserGRPC.UnimplementedActorUserServiceServer

	mongoClient             *mongo.Client
	mongoCollection         *mongo.Collection
	cacheClient             cache.Cache
	nativeIAmIdentityClient nativeIAmIdentityGRPC.IAMIdentityServiceClient
}

type userInMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Login    string             `bson:"login"`
	Identity string             `bson:"identity"`

	FullName string `bson:"fullName"`
	Avatar   string `bson:"avatar"`
	Email    string `bson:"email"`
}

const (
	UUID_CACHE_KEY_PREFIX    = "native_actor_user_data_uuid_"
	UUID_CACHE_TIMEOUT       = time.Second * 60
	IDENITY_CACHE_KEY_PREFIX = "native_actor_user_data_identity_"
	IDENTITY_CACHE_TIMEOUT   = time.Second * 60
	LOGIN_CACHE_KEY_PREFIX   = "native_actor_user_data_login_"
	LOGIN_CACHE_TIMEOUT      = time.Second * 60
	LOGIN_INDEX_NAME         = "unique_login"
	IDENTITY_INDEX_NAME      = "search_identity"
)

func clearUserCache(ctx context.Context, s *ActorUserServer, user *userInMongo) error {
	return s.cacheClient.Remove(
		ctx,
		UUID_CACHE_KEY_PREFIX+user.ID.Hex(),
		IDENITY_CACHE_KEY_PREFIX+user.Identity,
		LOGIN_CACHE_KEY_PREFIX+user.Login,
	)
}

func (u *userInMongo) ToProtoUser() *nativeActorUserGRPC.User {
	return &nativeActorUserGRPC.User{
		Uuid:     u.ID.Hex(),
		Login:    u.Login,
		Identity: u.Identity,
		FullName: u.FullName,
		Avatar:   u.Avatar,
		Email:    u.Email,
	}
}

func NewActorUserServer(ctx context.Context, mongoClient *mongo.Client, cacheClient cache.Cache, nativeIAmIdentityClient nativeIAmIdentityGRPC.IAMIdentityServiceClient) (*ActorUserServer, error) {
	collection := mongoClient.Database("openbp_global").Collection("native_actor_user")

	// Ensure indexes on the collection
	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Options: options.Index().SetName(LOGIN_INDEX_NAME).SetUnique(true),
			Keys:    bson.D{bson.E{Key: "login", Value: 1}},
		},
		{
			Options: options.Index().SetName(IDENTITY_INDEX_NAME),
			Keys:    bson.D{bson.E{Key: "identity", Value: "hashed"}},
		},
	})
	if err != nil {
		return nil, errors.New("Failed to create indexes. " + err.Error())
	}

	return &ActorUserServer{
		mongoClient:             mongoClient,
		mongoCollection:         collection,
		cacheClient:             cacheClient,
		nativeIAmIdentityClient: nativeIAmIdentityClient,
	}, nil
}

func (s *ActorUserServer) Create(ctx context.Context, in *nativeActorUserGRPC.CreateRequest) (*nativeActorUserGRPC.CreateResponse, error) {
	// Create identity for user
	identityResponse, err := s.nativeIAmIdentityClient.Create(ctx, &nativeIAmIdentityGRPC.CreateIdentityRequest{
		Namespace:       "",
		Name:            "",
		InitiallyActive: true,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to create identity for user. "+err.Error())
	}

	user := &userInMongo{
		Login:    in.Login,
		Identity: identityResponse.Identity.Uuid,
		FullName: in.FullName,
		Avatar:   in.Avatar,
		Email:    in.Email,
	}
	insertResult, err := s.mongoCollection.InsertOne(ctx, user)
	if err != nil {
		// Delete identity on error. Dont care about errors, because even it will occure it is harmless. Chance of error is ridicously low so there is no reason to do something with it.
		s.nativeIAmIdentityClient.Delete(ctx, &nativeIAmIdentityGRPC.DeleteIdentityRequest{
			Namespace: "",
			Uuid:      identityResponse.Identity.Uuid,
		})

		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(grpccodes.AlreadyExists, "User with same login already exist")
		}

		return nil, status.Error(grpccodes.Internal, "Failed to create user in database. "+err.Error())
	}
	user.ID = insertResult.InsertedID.(primitive.ObjectID)

	return &nativeActorUserGRPC.CreateResponse{
		User: user.ToProtoUser(),
	}, status.Error(grpccodes.OK, "")
}

func (s *ActorUserServer) Get(ctx context.Context, in *nativeActorUserGRPC.GetRequest) (*nativeActorUserGRPC.GetResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = UUID_CACHE_KEY_PREFIX + in.Uuid
		cacheBytes, err := s.cacheClient.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var user nativeActorUserGRPC.User
			err = proto.Unmarshal(cacheBytes, &user)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Failed to unmarshall user from cache. "+err.Error())
			}
			return &nativeActorUserGRPC.GetResponse{User: &user}, status.Error(grpccodes.OK, "")
		}
	}

	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "User UUID has bad format")
	}

	var user userInMongo
	err = s.mongoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with this UUID not found")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get user from database. "+err.Error())
	}
	protoUser := user.ToProtoUser()

	if in.UseCache {
		userBytes, err := proto.Marshal(protoUser)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall user to cache. "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, userBytes, UUID_CACHE_TIMEOUT)
	}

	return &nativeActorUserGRPC.GetResponse{
		User: protoUser,
	}, status.Errorf(grpccodes.OK, "")
}

func (s *ActorUserServer) GetByLogin(ctx context.Context, in *nativeActorUserGRPC.GetByLoginRequest) (*nativeActorUserGRPC.GetByLoginResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = LOGIN_CACHE_KEY_PREFIX + in.Login
		cacheBytes, err := s.cacheClient.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var user nativeActorUserGRPC.User
			err = proto.Unmarshal(cacheBytes, &user)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Failed to unmarshall user from cache. "+err.Error())
			}
			return &nativeActorUserGRPC.GetByLoginResponse{User: &user}, status.Error(grpccodes.OK, "")
		}
	}

	var user userInMongo
	err := s.mongoCollection.FindOne(ctx, bson.M{"login": in.Login}, options.FindOne().SetHint(LOGIN_INDEX_NAME)).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with this login not found")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get user from database. "+err.Error())
	}
	protoUser := user.ToProtoUser()

	if in.UseCache {
		userBytes, err := proto.Marshal(protoUser)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall user to cache. "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, userBytes, LOGIN_CACHE_TIMEOUT)
	}

	return &nativeActorUserGRPC.GetByLoginResponse{
		User: protoUser,
	}, status.Errorf(grpccodes.OK, "")
}

func (s *ActorUserServer) GetByIdentity(ctx context.Context, in *nativeActorUserGRPC.GetByIdentityRequest) (*nativeActorUserGRPC.GetByIdentityResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = IDENITY_CACHE_KEY_PREFIX + in.Identity
		cacheBytes, err := s.cacheClient.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var user nativeActorUserGRPC.User
			err = proto.Unmarshal(cacheBytes, &user)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Failed to unmarshall user from cache. "+err.Error())
			}
			return &nativeActorUserGRPC.GetByIdentityResponse{User: &user}, status.Error(grpccodes.OK, "")
		}
	}

	var user userInMongo
	err := s.mongoCollection.FindOne(ctx, bson.M{"identity": in.Identity}, options.FindOne().SetHint(IDENTITY_INDEX_NAME)).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with this identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get user from database. "+err.Error())
	}
	protoUser := user.ToProtoUser()

	if in.UseCache {
		userBytes, err := proto.Marshal(protoUser)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall user to cache. "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, userBytes, IDENTITY_CACHE_TIMEOUT)
	}

	return &nativeActorUserGRPC.GetByIdentityResponse{
		User: protoUser,
	}, status.Errorf(grpccodes.OK, "")
}

func (s *ActorUserServer) Update(ctx context.Context, in *nativeActorUserGRPC.UpdateRequest) (*nativeActorUserGRPC.UpdateResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "User UUID has bad format")
	}

	var user userInMongo
	err = s.mongoCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"login":    in.Login,
		"avatar":   in.Avatar,
		"fullName": in.FullName,
		"email":    in.Email,
	}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(grpccodes.AlreadyExists, "Login already exists")
		}
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with specified UUID not found")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to update user in database. "+err.Error())
	}

	clearUserCache(ctx, s, &user)

	return &nativeActorUserGRPC.UpdateResponse{
		User: user.ToProtoUser(),
	}, status.Error(grpccodes.OK, "")
}

func (s *ActorUserServer) Delete(ctx context.Context, in *nativeActorUserGRPC.DeleteRequest) (*nativeActorUserGRPC.DeleteResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "User UUID has bad format")
	}

	var user userInMongo
	err = s.mongoCollection.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &nativeActorUserGRPC.DeleteResponse{}, status.Error(grpccodes.OK, "")
		}
		return nil, status.Error(grpccodes.Internal, "Failed to delete user from the database. "+err.Error())
	}

	clearUserCache(ctx, s, &user)

	s.nativeIAmIdentityClient.Delete(ctx, &nativeIAmIdentityGRPC.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      user.Identity,
	})

	return &nativeActorUserGRPC.DeleteResponse{}, status.Error(grpccodes.OK, "")
}
