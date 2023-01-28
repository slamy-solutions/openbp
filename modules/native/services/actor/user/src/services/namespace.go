package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"

	nativeActorUserGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
)

type ActorUserNamespaceServer struct {
	nativeActorUserGRPC.UnimplementedActorUserServiceServer

	mongoClient                  *mongo.Client
	mongoUserCollection          *mongo.Collection
	mongoUserNamespaceCollection *mongo.Collection
	cacheClient                  cache.Cache
}

type userNamespaceInMongo struct {
	User      primitive.ObjectID `bson:"user"`
	Namespace string             `bson:"namespace"`
}

const (
	NAMESPACE_INDEX_NAME      = "namespace_user_unique"
	NAMESPACE_USER_INDEX_NAME = "user_namespaces"
)

func NewActorUserNamespaceServer(ctx context.Context, mongoClient *mongo.Client, cacheClient cache.Cache) (*ActorUserNamespaceServer, error) {
	mongoUserCollection := mongoClient.Database("openbp_global").Collection("native_actor_user")
	mongoUserNamespaceCollection := mongoClient.Database("openbp_global").Collection("native_actor_user_namespace")

	// Ensure indexes on the collection
	_, err := mongoUserNamespaceCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Options: options.Index().SetName(NAMESPACE_INDEX_NAME).SetUnique(true),
			Keys:    bson.D{bson.E{Key: "namespace", Value: 1}, bson.E{Key: "user", Value: 1}},
		},
		{
			Options: options.Index().SetName(NAMESPACE_USER_INDEX_NAME),
			Keys:    bson.D{bson.E{Key: "user", Value: "hashed"}},
		},
	})
	if err != nil {
		return nil, errors.New("Failed to create indexes. " + err.Error())
	}

	return &ActorUserNamespaceServer{
		mongoClient:                  mongoClient,
		mongoUserCollection:          mongoUserCollection,
		mongoUserNamespaceCollection: mongoUserNamespaceCollection,
		cacheClient:                  cacheClient,
	}, nil
}

/*
func (s *ActorUserNamespaceServer) GetUsersForNamespace(in *nativeActorUserGRPC.GetUsersForNamespaceRequest, out nativeActorUserGRPC.ActorUserNamespaceService_GetUsersForNamespaceServer) error {
	ctx := out.Context()

	opts :=
	if in.Limit != 0 {
		opts = opts.SetLimit(int64(in.Limit))
	}
	opts.SetHint(NAMESPACE_INDEX_NAME)

	type queryResultInMongo struct {
		User      primitive.ObjectID `bson:"user"`
		Namespace string             `bson:"namespace"`
		UserData  *userInMongo       `bson:"user_data"`
	}

	// cursor, err := s.mongoUserNamespaceCollection.Find(ctx, bson.M{"namespace": in.Namespace}, opts)
	s.mongoUserNamespaceCollection.Aggregate(ctx, bson.A{
		// Match the namespace
		bson.D{
			bson.E{"$match", bson.D{
				bson.E{"namespace", in.Namespace},
			}},
		},

		// Join the information about user
		bson.D{
			bson.E{"$lookup", bson.D{
				bson.E{"from", "native_actor_user"},
				bson.E{"localField", "user"},
				bson.E{"foreignField", "_id"},
				bson.E{"as", "user_data"},
			}},
		},

		bson.D{
			bson.E{"$unwind", "$user_data"},
		},
	})



	if err != nil {
		return status.Error(grpccodes.Internal, "Failed to get users for namespace from the database. "+err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var mapping userNamespaceInMongo
		if err := cursor.Decode(&mapping); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
		if err := out.Send(&nativeActorUserGRPC.GetUsersForNamespaceResponse{User: user.ToProtoUser()}); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
	}
	if err := cursor.Err(); err != nil {
		return status.Error(grpccodes.Internal, "Unexpected cursor error on fetching data from database. "+err.Error())
	}

	return nil

	return status.Errorf(grpccodes.Unimplemented, "method GetUsersForNamespace not implemented")
}
func (s *ActorUserNamespaceServer) AddUserToNamespace(ctx context.Context, in *nativeActorUserGRPC.AddUserToNamespaceRequest) (*nativeActorUserGRPC.AddUserToNamespaceResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method AddUserToNamespace not implemented")
}
func (s *ActorUserNamespaceServer) RemoveUserFromNamespace(ctx context.Context, in *nativeActorUserGRPC.RemoveUserFromNamespaceRequest) (*nativeActorUserGRPC.RemoveUserFromNamespaceResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method RemoveUserFromNamespace not implemented")
}
func (s *ActorUserNamespaceServer) GetNamespacesForUser(ctx context.Context, in *nativeActorUserGRPC.GetNamespacesForUserRequest) (*nativeActorUserGRPC.GetNamespacesForUserResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method GetNamespacesForUser not implemented")
}
func (s *ActorUserNamespaceServer) IsUserInTheNamespace(ctx context.Context, in *nativeActorUserGRPC.IsUserInTheNamespaceRequest) (*nativeActorUserGRPC.IsUserInTheNamespaceResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method IsUserInTheNamespace not implemented")
}
*/
