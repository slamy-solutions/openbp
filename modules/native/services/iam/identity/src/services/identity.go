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
	grpccodes "google.golang.org/grpc/codes"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmRoleGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type IAmIdentityServer struct {
	nativeIAmIdentityGRPC.UnimplementedIAMIdentityServiceServer

	mongoGlobalCollection *mongo.Collection

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

const (
	IDENTITY_CACHE_TIMEOUT = time.Second * 30
)

func NewIAmIdentityServer(systemStub *system.SystemStub, nativeStub *native.NativeStub) *IAmIdentityServer {
	mongoGlobalCollection := systemStub.DB.Database("openbp_global").Collection("native_iam_identity")
	return &IAmIdentityServer{
		mongoGlobalCollection: mongoGlobalCollection,
		systemStub:            systemStub,
		nativeStub:            nativeStub,
	}
}

func makeIndetityCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_identity_data_%s_%s", namespace, uuid)
}

func collectionByNamespace(s *IAmIdentityServer, namespace string) *mongo.Collection {
	if namespace == "" {
		return s.mongoGlobalCollection
	} else {
		db := s.systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
		return db.Collection("native_iam_identity")
	}
}

func (s *IAmIdentityServer) Create(ctx context.Context, in *nativeIAmIdentityGRPC.CreateIdentityRequest) (*nativeIAmIdentityGRPC.CreateIdentityResponse, error) {
	if in.Namespace != "" {
		r, err := s.nativeStub.Services.Namespace.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while checking namespace: "+err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	var managed managementTypeInMongo
	managed.ManagementType = identity_managed_no
	if t, ok := in.Managed.(*nativeIAmIdentityGRPC.CreateIdentityRequest_Identity); ok {
		managed.ManagementType = identity_managed_identity
		managed.IdentityNamespace = t.Identity.IdentityNamespace
		managed.IdentityUUID = t.Identity.IdentityUUID
	}
	if t, ok := in.Managed.(*nativeIAmIdentityGRPC.CreateIdentityRequest_Service); ok {
		managed.ManagementType = identity_managed_service
		managed.ServiceName = t.Service.Service
		managed.ServiceReason = t.Service.Reason
		managed.ServiceManagementID = t.Service.ManagementId
	}

	creationTime := time.Now().UTC()
	insertData := identityInMongo{
		Name:     in.Name,
		Managed:  &managed,
		Active:   in.InitiallyActive,
		Policies: []identityPolicyInMongo{},
		Roles:    []identityRoleInMongo{},
		Created:  creationTime,
		Updated:  creationTime,
		Version:  0,
	}

	collection := collectionByNamespace(s, in.Namespace)

	if insertData.Managed.ManagementType == identity_managed_service && insertData.Managed.ServiceManagementID != "" {
		r, err := collection.UpdateOne(
			ctx,
			bson.M{"managed": bson.M{
				"serviceName":         insertData.Managed.ServiceName,
				"serviceManagementId": insertData.Managed.ServiceManagementID,
			}},
			bson.M{
				"$setOnInsert": insertData,
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while adding identity to the database. "+err.Error())
		}
		if r.UpsertedCount == 0 {
			return nil, status.Error(grpccodes.AlreadyExists, "Identity managed by this service with same managementId already exists.")
		}
		insertData.ID = r.UpsertedID.(primitive.ObjectID)
	} else {
		insertResponse, err := collection.InsertOne(ctx, insertData)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error on inserting to DB: "+err.Error())
		}
		insertData.ID = insertResponse.InsertedID.(primitive.ObjectID)
	}

	return &nativeIAmIdentityGRPC.CreateIdentityResponse{
		Identity: insertData.ToGRPCIdentity(in.Namespace),
	}, status.Errorf(grpccodes.OK, "")
}

func (s *IAmIdentityServer) Get(ctx context.Context, in *nativeIAmIdentityGRPC.GetIdentityRequest) (*nativeIAmIdentityGRPC.GetIdentityResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeIndetityCacheKey(in.Namespace, in.Uuid)
		byteData, _ := s.systemStub.Cache.Get(ctx, cacheKey)
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

	identity := mongoIdentity.ToGRPCIdentity(in.Namespace)

	if in.UseCache {
		identityBytes, err := proto.Marshal(identity)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling identiry to cache: "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, identityBytes, IDENTITY_CACHE_TIMEOUT)
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
		s.systemStub.Cache.Remove(ctx, makeIndetityCacheKey(in.Namespace, in.Uuid))
	}

	return &nativeIAmIdentityGRPC.DeleteIdentityResponse{Existed: result.DeletedCount != 0}, status.Error(grpccodes.OK, "")
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

	existResponse, err := s.nativeStub.Services.IamPolicy.Exist(ctx, &nativeIAmPolicyGRPC.ExistPolicyRequest{
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

	policy := &identityPolicyInMongo{
		Namespace: in.PolicyNamespace,
		UUID:      in.PolicyUUID,
	}
	collection := collectionByNamespace(s, in.IdentityNamespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": identityId},
		bson.M{
			"$addToSet":    bson.M{"policies": policy},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity := mongoIdentity.ToGRPCIdentity(in.IdentityNamespace)

	s.systemStub.Cache.Remove(ctx, makeIndetityCacheKey(in.IdentityNamespace, in.IdentityUUID))

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

	policy := &identityPolicyInMongo{
		Namespace: in.PolicyNamespace,
		UUID:      in.PolicyUUID,
	}
	collection := collectionByNamespace(s, in.IdentityNamespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": identityId},
		bson.M{
			"$pull":        bson.M{"policies": policy},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity := mongoIdentity.ToGRPCIdentity(in.IdentityNamespace)

	s.systemStub.Cache.Remove(ctx, makeIndetityCacheKey(in.IdentityNamespace, in.IdentityUUID))

	return &nativeIAmIdentityGRPC.RemovePolicyResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}

func (s *IAmIdentityServer) AddRole(ctx context.Context, in *nativeIAmIdentityGRPC.AddRoleRequest) (*nativeIAmIdentityGRPC.AddRoleResponse, error) {
	identityId, err := primitive.ObjectIDFromHex(in.IdentityUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	_, err = primitive.ObjectIDFromHex(in.RoleUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Role UUID has bad format")
	}

	// validate if role exist. This flow is not atomic but it will prevent most of the problems
	existResponse, err := s.nativeStub.Services.IamRole.Exist(ctx, &nativeIAmRoleGRPC.ExistRoleRequest{
		Namespace: in.RoleNamespace,
		Uuid:      in.RoleUUID,
		UseCache:  true,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to check existance of role: "+err.Error())
	}
	if !existResponse.Exist {
		return nil, status.Error(grpccodes.FailedPrecondition, "Policy doesnt exist")
	}

	role := &identityRoleInMongo{
		Namespace: in.RoleNamespace,
		UUID:      in.RoleUUID,
	}
	collection := collectionByNamespace(s, in.IdentityNamespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": identityId},
		bson.M{
			"$addToSet":    bson.M{"roles": role},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity := mongoIdentity.ToGRPCIdentity(in.IdentityNamespace)

	s.systemStub.Cache.Remove(ctx, makeIndetityCacheKey(in.IdentityNamespace, in.IdentityUUID))

	return &nativeIAmIdentityGRPC.AddRoleResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}
func (s *IAmIdentityServer) RemoveRole(ctx context.Context, in *nativeIAmIdentityGRPC.RemoveRoleRequest) (*nativeIAmIdentityGRPC.RemoveRoleResponse, error) {
	identityId, err := primitive.ObjectIDFromHex(in.IdentityUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	_, err = primitive.ObjectIDFromHex(in.RoleUUID)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Role UUID has bad format")
	}

	role := &identityRoleInMongo{
		Namespace: in.RoleNamespace,
		UUID:      in.RoleUUID,
	}
	collection := collectionByNamespace(s, in.IdentityNamespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": identityId},
		bson.M{
			"$pull":        bson.M{"roles": role},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity := mongoIdentity.ToGRPCIdentity(in.IdentityNamespace)

	s.systemStub.Cache.Remove(ctx, makeIndetityCacheKey(in.IdentityNamespace, in.IdentityUUID))

	return &nativeIAmIdentityGRPC.RemoveRoleResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}

func (s *IAmIdentityServer) SetActive(ctx context.Context, in *nativeIAmIdentityGRPC.SetIdentityActiveRequest) (*nativeIAmIdentityGRPC.SetIdentityActiveResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Identity UUID has bad format")
	}

	collection := collectionByNamespace(s, in.Namespace)
	var mongoIdentity identityInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set":         bson.M{"active": in.Active},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&mongoIdentity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Identity not found")
		}
		return nil, status.Error(grpccodes.Internal, "Error on updating identity: "+err.Error())
	}

	identity := mongoIdentity.ToGRPCIdentity(in.Namespace)

	s.systemStub.Cache.Remove(ctx, makeIndetityCacheKey(in.Namespace, in.Uuid))

	return &nativeIAmIdentityGRPC.SetIdentityActiveResponse{Identity: identity}, status.Error(grpccodes.OK, "")
}
