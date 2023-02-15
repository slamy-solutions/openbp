package services

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/golang/protobuf/proto"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	nativeActorUserGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
)

type ActorUserServer struct {
	nativeActorUserGRPC.UnimplementedActorUserServiceServer

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

type userInMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Login    string             `bson:"login"`
	Identity string             `bson:"identity"`

	FullName string `bson:"fullName"`
	Avatar   string `bson:"avatar"`
	Email    string `bson:"email"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Version uint64    `bson:"version"`
}

const (
	UUID_CACHE_TIMEOUT     = time.Second * 60
	IDENTITY_CACHE_TIMEOUT = time.Second * 60
	LOGIN_CACHE_TIMEOUT    = time.Second * 60
)

func makeUserCacheKey(namespace string, userUUID string) string {
	return "native_actor_user_data_uuid_" + namespace + "_" + userUUID
}
func makeIdentityCacheKey(namespace string, identityUUID string) string {
	return "native_actor_user_data_identity_" + namespace + "_" + identityUUID
}
func makeLoginCacheKey(namespace string, login string) string {
	return "native_actor_user_data_login_" + namespace + "_" + login
}

func clearUserCache(ctx context.Context, s *ActorUserServer, namespace string, user *userInMongo) error {
	return s.systemStub.Cache.Remove(
		ctx,
		makeUserCacheKey(namespace, user.ID.Hex()),
		makeIdentityCacheKey(namespace, user.Identity),
		makeLoginCacheKey(namespace, user.Login),
	)
}

func (u *userInMongo) ToGRPCUser(namespace string) *nativeActorUserGRPC.User {
	return &nativeActorUserGRPC.User{
		Namespace: namespace,
		Uuid:      u.ID.Hex(),
		Login:     u.Login,
		Identity:  u.Identity,
		FullName:  u.FullName,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Created:   timestamppb.New(u.Created),
		Updated:   timestamppb.New(u.Updated),
		Version:   u.Version,
	}
}

func NewActorUserServer(ctx context.Context, systemStub *system.SystemStub, nativeStub *native.NativeStub) (*ActorUserServer, error) {
	// Ensure indexes on the collection
	err := ensureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("Failed to ensure indexes. " + err.Error())
	}

	return &ActorUserServer{
		systemStub: systemStub,
		nativeStub: nativeStub,
	}, nil
}

func (s *ActorUserServer) collectionByNamespace(namespace string) *mongo.Collection {
	collection := s.systemStub.DB.Database("openbp_global").Collection("native_actor_user")
	if namespace != "" {
		collection = s.systemStub.DB.Database("openbp_namespace_" + namespace).Collection("native_actor_user")
	}
	return collection
}

func (s *ActorUserServer) Create(ctx context.Context, in *nativeActorUserGRPC.CreateRequest) (*nativeActorUserGRPC.CreateResponse, error) {
	// Create identity for user
	identityResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &nativeIAmIdentityGRPC.CreateIdentityRequest{
		Namespace: in.Namespace,
		Name:      "",
		Managed: &nativeIAmIdentityGRPC.CreateIdentityRequest_Service{
			Service: &nativeIAmIdentityGRPC.ServiceManagedData{
				Service:      "native_actor_user",
				Reason:       "Manage identity for user",
				ManagementId: "",
			},
		},
		InitiallyActive: true,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to create identity for user. "+err.Error())
	}

	createdTime := time.Now().UTC()

	collection := s.collectionByNamespace(in.Namespace)

	user := &userInMongo{
		Login:    in.Login,
		Identity: identityResponse.Identity.Uuid,
		FullName: in.FullName,
		Avatar:   in.Avatar,
		Email:    in.Email,
		Created:  createdTime,
		Updated:  createdTime,
		Version:  0,
	}
	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		// Delete identity on error. Dont care about errors, because even it will occure it is harmless. Chance of error is ridicously low so there is no reason to do something with it.
		s.nativeStub.Services.IamIdentity.Delete(ctx, &nativeIAmIdentityGRPC.DeleteIdentityRequest{
			Namespace: "",
			Uuid:      identityResponse.Identity.Uuid,
		})

		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(grpccodes.AlreadyExists, "User with same login already exist")
		}

		return nil, status.Error(grpccodes.Internal, "Failed to create user in database. "+err.Error())
	}
	user.ID = insertResult.InsertedID.(primitive.ObjectID)

	log.Info("Created user with UUID [" + user.ID.Hex() + "] in namespace [" + in.Namespace + "]")

	return &nativeActorUserGRPC.CreateResponse{
		User: user.ToGRPCUser(in.Namespace),
	}, status.Error(grpccodes.OK, "")
}

func (s *ActorUserServer) Get(ctx context.Context, in *nativeActorUserGRPC.GetRequest) (*nativeActorUserGRPC.GetResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeUserCacheKey(in.Namespace, in.Uuid)
		cacheBytes, _ := s.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var user nativeActorUserGRPC.User
			err := proto.Unmarshal(cacheBytes, &user)
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

	collection := s.collectionByNamespace(in.Namespace)

	var user userInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with this UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "User wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get user from database. "+err.Error())
	}
	userGRPC := user.ToGRPCUser(in.Namespace)

	if in.UseCache {
		userBytes, err := proto.Marshal(userGRPC)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall user to cache. "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, userBytes, UUID_CACHE_TIMEOUT)
	}

	return &nativeActorUserGRPC.GetResponse{
		User: userGRPC,
	}, status.Errorf(grpccodes.OK, "")
}

func (s *ActorUserServer) GetByLogin(ctx context.Context, in *nativeActorUserGRPC.GetByLoginRequest) (*nativeActorUserGRPC.GetByLoginResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeLoginCacheKey(in.Namespace, in.Login)
		cacheBytes, _ := s.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var user nativeActorUserGRPC.User
			err := proto.Unmarshal(cacheBytes, &user)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Failed to unmarshall user from cache. "+err.Error())
			}
			return &nativeActorUserGRPC.GetByLoginResponse{User: &user}, status.Error(grpccodes.OK, "")
		}
	}

	collection := s.collectionByNamespace(in.Namespace)

	var user userInMongo
	err := collection.FindOne(ctx, bson.M{"login": in.Login}, options.FindOne().SetHint(unique_login_index)).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with this login not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "User wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get user from database. "+err.Error())
	}
	userGRPC := user.ToGRPCUser(in.Namespace)

	if in.UseCache {
		userBytes, err := proto.Marshal(userGRPC)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall user to cache. "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, userBytes, LOGIN_CACHE_TIMEOUT)
	}

	return &nativeActorUserGRPC.GetByLoginResponse{
		User: userGRPC,
	}, status.Errorf(grpccodes.OK, "")
}

func (s *ActorUserServer) GetByIdentity(ctx context.Context, in *nativeActorUserGRPC.GetByIdentityRequest) (*nativeActorUserGRPC.GetByIdentityResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeIdentityCacheKey(in.Namespace, in.Identity)
		cacheBytes, _ := s.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var user nativeActorUserGRPC.User
			err := proto.Unmarshal(cacheBytes, &user)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Failed to unmarshall user from cache. "+err.Error())
			}
			return &nativeActorUserGRPC.GetByIdentityResponse{User: &user}, status.Error(grpccodes.OK, "")
		}
	}

	collection := s.collectionByNamespace(in.Namespace)

	var user userInMongo
	err := collection.FindOne(ctx, bson.M{"identity": in.Identity}, options.FindOne().SetHint(fast_identity_search_index)).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with this identity not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "User wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get user from database. "+err.Error())
	}
	userGRPC := user.ToGRPCUser(in.Namespace)

	if in.UseCache {
		userBytes, err := proto.Marshal(userGRPC)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to marshall user to cache. "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, userBytes, IDENTITY_CACHE_TIMEOUT)
	}

	return &nativeActorUserGRPC.GetByIdentityResponse{
		User: userGRPC,
	}, status.Errorf(grpccodes.OK, "")
}

func (s *ActorUserServer) Update(ctx context.Context, in *nativeActorUserGRPC.UpdateRequest) (*nativeActorUserGRPC.UpdateResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "User UUID has bad format")
	}

	collection := s.collectionByNamespace(in.Namespace)

	var user userInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"login":    in.Login,
			"avatar":   in.Avatar,
			"fullName": in.FullName,
			"email":    in.Email,
		},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
		"$inc":         bson.M{"version": 1},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, status.Error(grpccodes.AlreadyExists, "Login already exists")
		}
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "User with specified UUID not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "User wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to update user in database. "+err.Error())
	}

	clearUserCache(ctx, s, in.Namespace, &user)

	return &nativeActorUserGRPC.UpdateResponse{
		User: user.ToGRPCUser(in.Namespace),
	}, status.Error(grpccodes.OK, "")
}

func (s *ActorUserServer) Delete(ctx context.Context, in *nativeActorUserGRPC.DeleteRequest) (*nativeActorUserGRPC.DeleteResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "User UUID has bad format")
	}

	collection := s.collectionByNamespace(in.Namespace)

	var user userInMongo
	err = collection.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &nativeActorUserGRPC.DeleteResponse{Existed: false}, status.Error(grpccodes.OK, "")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &nativeActorUserGRPC.DeleteResponse{Existed: false}, status.Error(grpccodes.OK, "Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to delete user from the database. "+err.Error())
	}

	clearUserCache(ctx, s, in.Namespace, &user)
	log.Info("Deleted user with UUID [" + user.ID.Hex() + "] in namespace [" + in.Namespace + "]")

	_, err = s.nativeStub.Services.IamIdentity.Delete(ctx, &nativeIAmIdentityGRPC.DeleteIdentityRequest{
		Namespace: in.Namespace,
		Uuid:      user.Identity,
	})
	if err != nil {
		log.Error("Failed to delete identity [" + user.Identity + "] of deleted user [" + user.ID.Hex() + "] in namespace [" + in.Namespace + "]. " + err.Error())
	}

	return &nativeActorUserGRPC.DeleteResponse{Existed: true}, status.Error(grpccodes.OK, "")
}

func (s *ActorUserServer) List(in *nativeActorUserGRPC.ListRequest, out nativeActorUserGRPC.ActorUserService_ListServer) error {
	ctx := out.Context()

	collection := s.collectionByNamespace(in.Namespace)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				status.Error(grpccodes.OK, "Namespace doesnt exist.")
			}
		}
		return status.Error(grpccodes.Internal, "Failed to get users from the database. "+err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user userInMongo
		if err := cursor.Decode(&user); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&nativeActorUserGRPC.ListResponse{User: user.ToGRPCUser(in.Namespace)}); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
	}
	if err := cursor.Err(); err != nil {
		return status.Error(grpccodes.Internal, "Unexpected cursor error on fetching data from database. "+err.Error())
	}

	return status.Error(grpccodes.OK, "")
}

func (s *ActorUserServer) Search(in *nativeActorUserGRPC.SearchRequest, out nativeActorUserGRPC.ActorUserService_SearchServer) error {
	ctx := out.Context()

	opts := options.Find()
	if in.Limit != 0 {
		opts = opts.SetLimit(int64(in.Limit))
	}

	collection := s.collectionByNamespace(in.Namespace)

	cursor, err := collection.Find(ctx, bson.M{"$text": bson.M{"$search": in.Match}}, opts)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				status.Error(grpccodes.OK, "Namespace doesnt exist.")
			}
		}
		return status.Error(grpccodes.Internal, "Failed to get users from the database. "+err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user userInMongo
		if err := cursor.Decode(&user); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&nativeActorUserGRPC.SearchResponse{User: user.ToGRPCUser(in.Namespace)}); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
	}
	if err := cursor.Err(); err != nil {
		return status.Error(grpccodes.Internal, "Unexpected cursor error on fetching data from database. "+err.Error())
	}

	return status.Error(grpccodes.OK, "")
}
