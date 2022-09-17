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

	"github.com/golang/protobuf/proto"
	"github.com/slamy-solutions/openbp/modules/system/libs/go/cache"
	grpccodes "google.golang.org/grpc/codes"

	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/policy/src/grpc/native_iam_policy"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/policy/src/grpc/native_namespace"
)

type IAMPolicyServer struct {
	nativeIAmPolicyGRPC.UnimplementedIAMPolicyServiceServer

	mongoClient           *mongo.Client
	mongoGlobalCollection *mongo.Collection
	cacheClient           cache.Cache
	nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient
}

type policyInMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Resources []string           `bson:"resources"`
	Actions   []string           `bson:"actions"`
}

const (
	POLICY_CACHE_TIMEOUT = time.Second * 30
)

func collectionByNamespace(s *IAMPolicyServer, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.mongoGlobalCollection
	} else {
		db := s.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
		return db.Collection("native_iam_policy")
	}
}

func makePolicyCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_policy_data_%s_%s", namespace, uuid)
}

func NewIAMPolicyServer(mongoClient *mongo.Client, cacheClient cache.Cache, nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient) *IAMPolicyServer {
	mongoGlobalCollection := mongoClient.Database("openbp_global").Collection("native_iam_policy")
	return &IAMPolicyServer{
		mongoClient:           mongoClient,
		mongoGlobalCollection: mongoGlobalCollection,
		cacheClient:           cacheClient,
		nativeNamespaceClient: nativeNamespaceClient,
	}
}

func (s *IAMPolicyServer) Create(ctx context.Context, in *nativeIAmPolicyGRPC.CreatePolicyRequest) (*nativeIAmPolicyGRPC.CreatePolicyResponse, error) {
	// Check if namespace really exist
	if in.Namespace != "" {
		r, err := s.nativeNamespaceClient.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	policy := policyInMongo{
		Name:      in.Name,
		Actions:   in.Actions,
		Resources: in.Resources,
	}
	collection := collectionByNamespace(s, in.Namespace)
	insertResponse, err := collection.InsertOne(ctx, policy)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	id := insertResponse.InsertedID.(primitive.ObjectID)

	return &nativeIAmPolicyGRPC.CreatePolicyResponse{
		Policy: &nativeIAmPolicyGRPC.Policy{
			Namespace: in.Namespace,
			Uuid:      id.Hex(),
			Name:      in.Name,
			Actions:   in.Actions,
			Resources: in.Resources,
		},
	}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) Get(ctx context.Context, in *nativeIAmPolicyGRPC.GetPolicyRequest) (*nativeIAmPolicyGRPC.GetPolicyResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makePolicyCacheKey(in.Namespace, in.Uuid)
		byteData, _ := s.cacheClient.Get(ctx, cacheKey)
		if byteData != nil {
			var policy nativeIAmPolicyGRPC.Policy
			err := proto.Unmarshal(byteData, &policy)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Error while unmarshaling policy from cache: "+err.Error())
			}
			return &nativeIAmPolicyGRPC.GetPolicyResponse{Policy: &policy}, status.Error(grpccodes.OK, "")
		}
	}

	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Policy UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	var mongoPolicy policyInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mongoPolicy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Policy with provided UUID and namespace not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	policy := &nativeIAmPolicyGRPC.Policy{
		Namespace: in.Namespace,
		Uuid:      in.Uuid,
		Name:      mongoPolicy.Name,
		Actions:   mongoPolicy.Actions,
		Resources: mongoPolicy.Resources,
	}

	if in.UseCache {
		policyBytes, err := proto.Marshal(policy)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling policy to cache: "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, policyBytes, POLICY_CACHE_TIMEOUT)
	}

	return &nativeIAmPolicyGRPC.GetPolicyResponse{Policy: policy}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) Exist(ctx context.Context, in *nativeIAmPolicyGRPC.ExistPolicyRequest) (*nativeIAmPolicyGRPC.ExistPolicyResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Policy UUID has bad format")
	}

	var cacheKey string
	if in.UseCache {
		cacheKey = makePolicyCacheKey(in.Namespace, in.Uuid)
		exist, err := s.cacheClient.Has(ctx, cacheKey)
		if err == nil && exist {
			return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: true}, status.Error(grpccodes.OK, "")
		}
	}

	collection := collectionByNamespace(s, in.Namespace)
	if in.UseCache {
		var mongoPolicy policyInMongo
		err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mongoPolicy)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: false}, status.Error(grpccodes.OK, "")
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		policy := &nativeIAmPolicyGRPC.Policy{
			Namespace: in.Namespace,
			Uuid:      in.Uuid,
			Name:      mongoPolicy.Name,
			Actions:   mongoPolicy.Actions,
			Resources: mongoPolicy.Resources,
		}
		policyBytes, err := proto.Marshal(policy)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling policy to cache: "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, policyBytes, POLICY_CACHE_TIMEOUT)
		return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: true}, status.Error(grpccodes.OK, "")
	} else {
		// Fast check in mongo if exist without getting data
		count, err := collection.CountDocuments(ctx, bson.M{"_id": id}, options.Count().SetLimit(1))
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}

		return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: count == 1}, status.Error(grpccodes.OK, "")
	}
}

func (s *IAMPolicyServer) Update(ctx context.Context, in *nativeIAmPolicyGRPC.UpdatePolicyRequest) (*nativeIAmPolicyGRPC.UpdatePolicyResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Bad UUID format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	updateData := &policyInMongo{
		Name:      in.Name,
		Actions:   in.Actions,
		Resources: in.Resources,
	}
	r, err := collection.UpdateByID(ctx, id, bson.M{"$set": updateData})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if r.MatchedCount == 0 {
		return nil, status.Error(grpccodes.NotFound, "Policy with specified namespace and uuid not found")
	}

	s.cacheClient.Remove(ctx, makePolicyCacheKey(in.Namespace, in.Uuid))

	return &nativeIAmPolicyGRPC.UpdatePolicyResponse{
		Policy: &nativeIAmPolicyGRPC.Policy{
			Namespace: in.Namespace,
			Uuid:      in.Uuid,
			Name:      in.Name,
			Actions:   in.Actions,
			Resources: in.Resources,
		},
	}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) Delete(ctx context.Context, in *nativeIAmPolicyGRPC.DeletePolicyRequest) (*nativeIAmPolicyGRPC.DeletePolicyResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Bad UUID format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	r, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if r.DeletedCount != 0 {
		s.cacheClient.Remove(ctx, makePolicyCacheKey(in.Namespace, in.Uuid))
	}

	return &nativeIAmPolicyGRPC.DeletePolicyResponse{}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) List(in *nativeIAmPolicyGRPC.ListPoliciesRequest, out nativeIAmPolicyGRPC.IAMPolicyService_ListServer) error {
	ctx := out.Context()

	collection := collectionByNamespace(s, in.Namespace)
	findOptions := options.Find()
	if in.Skip != 0 {
		findOptions = findOptions.SetSkip(int64(in.Skip))
	}
	if in.Limit != 0 {
		findOptions = findOptions.SetLimit(int64(in.Limit))
	}
	if in.Skip != 0 || in.Limit != 0 {
		// Make sure that same data will be returned for same skip and limit parameters
		findOptions = findOptions.SetSort(bson.D{primitive.E{Key: "_id", Value: 1}})
	}

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result policyInMongo
		if err := cursor.Decode(&result); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
		sendData := &nativeIAmPolicyGRPC.ListPoliciesResponse{
			Policy: &nativeIAmPolicyGRPC.Policy{
				Namespace: in.Namespace,
				Uuid:      result.ID.Hex(),
				Name:      result.Name,
				Actions:   result.Actions,
				Resources: result.Resources,
			},
		}
		if err := out.Send(sendData); err != nil {
			return status.Error(grpccodes.Internal, err.Error())
		}
	}
	if err := cursor.Err(); err != nil {
		return status.Error(grpccodes.Internal, err.Error())
	}

	return status.Error(grpccodes.OK, "")
}
