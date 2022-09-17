package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	grpccodes "google.golang.org/grpc/codes"

	"github.com/slamy-solutions/openbp/modules/system/libs/go/cache"

	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/identity/src/grpc/native_iam_identity"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/identity/src/grpc/native_iam_policy"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/services/iam/identity/src/grpc/native_namespace"
)

type IAmIdentityServer struct {
	nativeIAmIdentityGRPC.UnimplementedIAMIdentityServiceServer

	mongoClient           *mongo.Client
	mongoGlobalCollection *mongo.Collection
	cacheClient           cache.Cache
	nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient
	nativeIAmPolicyClient nativeIAmPolicyGRPC.IAMPolicyServiceClient
}

type identityInMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Active   bool               `bson:"active"`
	Policies []string           `bson:"policies"`
}

const (
	IDENTITY_CACHE_TIMEOUT = time.Second * 30
)

func NewIAmIdentityServer(mongoClient *mongo.Client, cacheClient cache.Cache, nativeNamespaceClient nativeNamespaceGRPC.NamespaceServiceClient, nativeIAmPolicyClient nativeIAmPolicyGRPC.IAMPolicyServiceClient) *IAmIdentityServer {
	mongoGlobalCollection := mongoClient.Database("openbp_global").Collection("native_iam_identity")
	return &IAmIdentityServer{
		mongoClient:           mongoClient,
		mongoGlobalCollection: mongoGlobalCollection,
		cacheClient:           cacheClient,
		nativeNamespaceClient: nativeNamespaceClient,
		nativeIAmPolicyClient: nativeIAmPolicyClient,
	}
}

func makeIndetityCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_identity_data_%s_%s", namespace, uuid)
}

func collectionByNamespace(s *IAmIdentityServer, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.mongoGlobalCollection
	} else {
		db := s.mongoClient.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
		return db.Collection("native_iam_identity")
	}
}

func identityFromMongo(mongoIdentity identityInMongo, namespace string, uuid string) (*nativeIAmIdentityGRPC.Identity, error) {
	policies := make([]*nativeIAmIdentityGRPC.Identity_PolicyReference, len(mongoIdentity.Policies))
	for index, policy := range mongoIdentity.Policies {
		splitedData := strings.Split(policy, ":")
		if len(splitedData) != 2 {
			return nil, status.Error(grpccodes.Internal, "Identity policy has bad format.")
		}
		policies[index] = &nativeIAmIdentityGRPC.Identity_PolicyReference{
			Namespace: splitedData[0],
			Uuid:      splitedData[1],
		}
	}
	return &nativeIAmIdentityGRPC.Identity{
		Namespace: namespace,
		Uuid:      uuid,
		Name:      mongoIdentity.Name,
		Active:    mongoIdentity.Active,
		Policies:  policies,
	}, nil
}

func (s *IAmIdentityServer) Create(ctx context.Context, in *nativeIAmIdentityGRPC.CreateIdentityRequest) (*nativeIAmIdentityGRPC.CreateIdentityResponse, error) {
	if in.Namespace != "" {
		r, err := s.nativeNamespaceClient.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while checking namespace: "+err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	collection := collectionByNamespace(s, in.Namespace)
	insertData := identityInMongo{
		Name:     in.Name,
		Active:   in.InitiallyActive,
		Policies: []string{},
	}
	insertResponse, err := collection.InsertOne(ctx, insertData)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error on inserting to DB: "+err.Error())
	}
	uuid := insertResponse.InsertedID.(primitive.ObjectID)

	return &nativeIAmIdentityGRPC.CreateIdentityResponse{
		Identity: &nativeIAmIdentityGRPC.Identity{
			Namespace: in.Namespace,
			Uuid:      uuid.Hex(),
			Name:      in.Name,
			Active:    in.InitiallyActive,
			Policies:  []*nativeIAmIdentityGRPC.Identity_PolicyReference{},
		},
	}, status.Errorf(grpccodes.OK, "")
}

func (s *IAmIdentityServer) Get(ctx context.Context, in *nativeIAmIdentityGRPC.GetIdentityRequest) (*nativeIAmIdentityGRPC.GetIdentityResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeIndetityCacheKey(in.Namespace, in.Uuid)
		byteData, _ := s.cacheClient.Get(ctx, cacheKey)
		if byteData != nil {
			var identity nativeIAmIdentityGRPC.Identity
			err := proto.Unmarshal(byteData, &identity)
			if err != nil {
				return nil, status.Error(grpccodes.Internal, "Error while unmarshaling identity from cache: "+err.Error())
			}
			return &nativeIAmIdentityGRPC.GetIdentityResponse{Identity: &identity}, status.Error(grpccodes.OK, "")
		}
	}

	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	var mongoIdentity identityInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity with provided UUID and namespace not found")
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	identity, err := identityFromMongo(mongoIdentity, in.Namespace, in.Uuid)
	if err != nil {
		return nil, err
	}

	if in.UseCache {
		identityBytes, err := proto.Marshal(identity)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling identiry to cache: "+err.Error())
		}
		s.cacheClient.Set(ctx, cacheKey, identityBytes, IDENTITY_CACHE_TIMEOUT)
	}

	return &nativeIAmIdentityGRPC.GetIdentityResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}

func (s *IAmIdentityServer) Delete(ctx context.Context, in *nativeIAmIdentityGRPC.DeleteIdentityRequest) (*nativeIAmIdentityGRPC.DeleteIdentityResponse, error) {
	identityId, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	result, err := collection.DeleteOne(ctx, bson.M{"_id": identityId})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error while deleting data from mongo: "+err.Error())
	}

	if result.DeletedCount != 0 {
		s.cacheClient.Remove(ctx, makeIndetityCacheKey(in.Namespace, in.Uuid))
	}

	return &nativeIAmIdentityGRPC.DeleteIdentityResponse{}, status.Error(grpccodes.OK, "")
}

func (s *IAmIdentityServer) AddPolicy(ctx context.Context, in *nativeIAmIdentityGRPC.AddPolicyRequest) (*nativeIAmIdentityGRPC.AddPolicyResponse, error) {
	identityId, err := primitive.ObjectIDFromHex(in.IdentityUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	_, err = primitive.ObjectIDFromHex(in.PolicyUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Policy UUID has bad format")
	}

	existResponse, err := s.nativeIAmPolicyClient.Exist(ctx, &nativeIAmPolicyGRPC.ExistPolicyRequest{
		Namespace: in.PolicyNamespace,
		Uuid:      in.PolicyUUID,
		UseCache:  true,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to check existance of policy: "+err.Error())
	}
	if !existResponse.Exist {
		return nil, status.Error(grpccodes.FailedPrecondition, "Policy doesnt exist")
	}

	policyData := in.PolicyNamespace + ":" + in.PolicyUUID
	collection := collectionByNamespace(s, in.IdentityNamespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": identityId}, bson.M{"$addToSet": bson.M{"policies": policyData}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity, err := identityFromMongo(mongoIdentity, in.PolicyNamespace, in.PolicyUUID)
	if err != nil {
		return nil, err
	}

	s.cacheClient.Remove(ctx, makeIndetityCacheKey(in.IdentityNamespace, in.IdentityUUID))

	return &nativeIAmIdentityGRPC.AddPolicyResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}

func (s *IAmIdentityServer) RemovePolicy(ctx context.Context, in *nativeIAmIdentityGRPC.RemovePolicyRequest) (*nativeIAmIdentityGRPC.RemovePolicyResponse, error) {
	identityId, err := primitive.ObjectIDFromHex(in.IdentityUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	_, err = primitive.ObjectIDFromHex(in.PolicyUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Policy UUID has bad format")
	}

	policyData := in.PolicyNamespace + ":" + in.PolicyUUID
	collection := collectionByNamespace(s, in.IdentityNamespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": identityId}, bson.M{"$pull": bson.M{"policies": policyData}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity, err := identityFromMongo(mongoIdentity, in.PolicyNamespace, in.PolicyUUID)
	if err != nil {
		return nil, err
	}

	s.cacheClient.Remove(ctx, makeIndetityCacheKey(in.IdentityNamespace, in.IdentityUUID))

	return &nativeIAmIdentityGRPC.RemovePolicyResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}

func (s *IAmIdentityServer) SetActive(ctx context.Context, in *nativeIAmIdentityGRPC.SetIdentityActiveRequest) (*nativeIAmIdentityGRPC.SetIdentityActiveResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"active": in.Active}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity, err := identityFromMongo(mongoIdentity, in.Namespace, in.Uuid)
	if err != nil {
		return nil, err
	}

	s.cacheClient.Remove(ctx, makeIndetityCacheKey(in.Namespace, in.Uuid))

	return &nativeIAmIdentityGRPC.SetIdentityActiveResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}
