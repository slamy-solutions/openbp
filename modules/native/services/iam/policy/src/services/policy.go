package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"
	grpccodes "google.golang.org/grpc/codes"

	nativeIAmPolicyGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/policy/src/grpc/native_iam_policy"
	nativeNamespaceGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/policy/src/grpc/native_namespace"
)

type IAMPolicyServer struct {
	nativeIAmPolicyGRPC.UnimplementedIAMPolicyServiceServer

	mongoClient           *mongo.Client
	mongoDbPrefix         string
	mongoGlobalCollection *mongo.Collection
	cacheClient           cache.Cache
	nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient
}

type policyInMongo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"name"`
	Actions []string           `bson:"actions"`
}

const (
	POLICY_CACHE_TIMEOUT = time.Second * 30
)

func collectionByNamespace(s *IAMPolicyServer, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.mongoGlobalCollection
	} else {
		db := s.mongoClient.Database(fmt.Sprintf("%snamespace_%s", s.mongoDbPrefix, namespace))
		return db.Collection("native_iam_policy")
	}
}

func makePolicyCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_policy_data_%s_%s", namespace, uuid)
}

func (s *IAMPolicyServer) Create(ctx context.Context, in *nativeIAmPolicyGRPC.CreatePolicyRequest) (*nativeIAmPolicyGRPC.CreatePolicyResponse, error) {
	// Check if namespace really exist
	r, err := s.nativeNamespaceClient.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Name, UseCache: true})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}
	if !r.Exist {
		return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
	}

	policy := policyInMongo{
		Name:    in.Name,
		Actions: in.Actions,
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

func (s *IAMPolicyServer) Update(ctx context.Context, in *nativeIAmPolicyGRPC.UpdatePolicyRequest) (*nativeIAmPolicyGRPC.UpdatePolicyResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method Update not implemented")
}

func (s *IAMPolicyServer) Delete(ctx context.Context, in *nativeIAmPolicyGRPC.DeletePolicyRequest) (*nativeIAmPolicyGRPC.DeletePolicyResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method Delete not implemented")
}

func (s *IAMPolicyServer) List(in *nativeIAmPolicyGRPC.ListPoliciesRequest, out nativeIAmPolicyGRPC.IAMPolicyService_ListServer) error {
	return status.Errorf(grpccodes.Unimplemented, "method List not implemented")
}
