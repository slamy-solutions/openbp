package role

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"google.golang.org/grpc/codes"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmRoleGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"

	policy_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/policy"
)

type IAMRoleServer struct {
	nativeIAmRoleGRPC.UnimplementedIAMRoleServiceServer

	systemStub *system.SystemStub
	nativeStub *native.NativeStub

	policyServer *policy_server.IAMPolicyServer
}

func (s *IAMRoleServer) getCollectionByNamespace(namespace string) *mongo.Collection {
	var db *mongo.Database
	if namespace == "" {
		db = s.systemStub.DB.Database("openbp_global")
	} else {
		db = s.systemStub.DB.Database("openbp_namespace_" + namespace)
	}

	return db.Collection("native_iam_role")
}

func NewIAMRoleServer(ctx context.Context, systemStub *system.SystemStub, nativeStub *native.NativeStub, policyServer *policy_server.IAMPolicyServer) (*IAMRoleServer, error) {

	// Ensure indexes for global namespace
	err := ensureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("Failed to ensure indexes for global namespace: " + err.Error())
	}

	// Ensure BuiltIn roles for global namespace
	err = ensureBuiltInsForNamespace(ctx, "", systemStub, nativeStub, policyServer)
	if err != nil {
		return nil, errors.New("Failed to ensure builtIns for global namespace: " + err.Error())
	}

	return &IAMRoleServer{
		systemStub:   systemStub,
		nativeStub:   nativeStub,
		policyServer: policyServer,
	}, nil
}

// TODO: add cache
const ROLE_CACHE_TIMEOUT = time.Second * 30

func makeCountRoleCacheKey(namespace string) string {
	return fmt.Sprintf("native_iam_role_count_%s", namespace)
}

func (s *IAMRoleServer) Create(ctx context.Context, in *nativeIAmRoleGRPC.CreateRoleRequest) (*nativeIAmRoleGRPC.CreateRoleResponse, error) {
	creationTime := time.Now().UTC()

	var managed managementTypeInMongo
	managed.ManagementType = role_managed_no
	if t, ok := in.Managed.(*nativeIAmRoleGRPC.CreateRoleRequest_Identity); ok {
		managed.ManagementType = role_managed_identity
		managed.IdentityNamespace = t.Identity.IdentityNamespace
		managed.IdentityUUID = t.Identity.IdentityUUID
	}
	if t, ok := in.Managed.(*nativeIAmRoleGRPC.CreateRoleRequest_Service); ok {
		managed.ManagementType = role_managed_service
		managed.ServiceName = t.Service.Service
		managed.ServiceReason = t.Service.Reason
		managed.ServiceManagementID = t.Service.ManagementId
	}

	role := &roleInMongo{
		Name:        in.Name,
		Description: in.Description,
		Policies:    []*assignedPolicyInMongo{},
		Managed:     &managed,
		Tags:        []string{},
		Created:     creationTime,
		Updated:     creationTime,
		Version:     0,
	}

	collection := s.getCollectionByNamespace(in.Namespace)

	// If this is managed by service, we have to additionally ensure "ServiceName" and "ManagementId" to be unique
	if role.Managed.ManagementType == role_managed_service && role.Managed.ServiceManagementID != "" {
		r, err := collection.UpdateOne(
			ctx,
			bson.M{
				"managed.serviceName":         role.Managed.ServiceName,
				"managed.serviceManagementId": role.Managed.ServiceManagementID,
			},
			bson.M{
				"$setOnInsert": role,
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error while adding role to the database. "+err.Error())
		}
		if r.UpsertedCount == 0 {
			return nil, status.Error(codes.AlreadyExists, "Role managed by this service with same managementId already exists.")
		}
		role.ID = r.UpsertedID.(primitive.ObjectID)
	} else {
		r, err := collection.InsertOne(ctx, role)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error while inserting role to the database. "+err.Error())
		}
		role.ID = r.InsertedID.(primitive.ObjectID)
	}

	log.Infof("Created new role with UUID [%s] and name [%s] in the [%s] namespace", role.ID.Hex(), in.Name, in.Namespace)
	s.systemStub.Cache.Remove(ctx, makeCountRoleCacheKey(in.Namespace))

	return &nativeIAmRoleGRPC.CreateRoleResponse{Role: role.ToGRPCRole(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) Get(ctx context.Context, in *nativeIAmRoleGRPC.GetRoleRequest) (*nativeIAmRoleGRPC.GetRoleResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Role wasnt founded. Role UUID has bad format.")
	}

	collection := s.getCollectionByNamespace(in.Namespace)
	var role roleInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Role wasnt founded.")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while searhing for role in database. "+err.Error())
	}

	return &nativeIAmRoleGRPC.GetRoleResponse{Role: role.ToGRPCRole(in.Namespace)}, status.Error(codes.OK, "")
}

func (s *IAMRoleServer) OpenGetMultipleChannel(ctx context.Context, in *nativeIAmRoleGRPC.GetMultipleRolesRequest) chan struct {
	Err       error
	Namespace string
	Role      roleInMongo
} {
	// Split all the roles into namespaces
	namespaces := make(map[string][]*nativeIAmRoleGRPC.GetMultipleRolesRequest_RequestedRole)
	for _, role := range in.Roles {
		roles, ok := namespaces[role.Namespace]
		if !ok {
			roles = make([]*nativeIAmRoleGRPC.GetMultipleRolesRequest_RequestedRole, 0, 1)
		}

		roles = append(roles, role)
		namespaces[role.Namespace] = roles
	}

	// For each namespace run mongo request
	var requestsWaitGroup sync.WaitGroup
	requestsWaitGroup.Add(len(namespaces))
	results := make(chan struct {
		Err       error
		Namespace string
		Role      roleInMongo
	})
	for namespace, roles := range namespaces {
		go func(namespace string, roles []*nativeIAmRoleGRPC.GetMultipleRolesRequest_RequestedRole) {
			defer requestsWaitGroup.Done()

			uuids := bson.A{}
			for _, role := range roles {
				id, err := primitive.ObjectIDFromHex(role.Uuid)
				if err == nil {
					uuids = append(uuids, id)
				}
			}

			collection := s.getCollectionByNamespace(namespace)
			cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": uuids}})
			if err != nil {
				if err, ok := err.(mongo.WriteException); ok {
					if err.HasErrorLabel("InvalidNamespace") {
						return //Its ok. Role not founded
					}
				}
				results <- struct {
					Err       error
					Namespace string
					Role      roleInMongo
				}{Err: err, Namespace: "", Role: roleInMongo{}}
				return
			}
			defer cursor.Close(ctx)

			for cursor.Next(ctx) {
				var role roleInMongo
				if err := cursor.Decode(&role); err != nil {
					results <- struct {
						Err       error
						Namespace string
						Role      roleInMongo
					}{Err: err, Namespace: "", Role: roleInMongo{}}
					return
				}
				results <- struct {
					Err       error
					Namespace string
					Role      roleInMongo
				}{Err: nil, Namespace: namespace, Role: role}
			}

			if err := cursor.Err(); err != nil {
				results <- struct {
					Err       error
					Namespace string
					Role      roleInMongo
				}{Err: err, Namespace: "", Role: roleInMongo{}}
			}
		}(namespace, roles)
	}

	// Wait for all requests to finish and close results channel
	go func() {
		requestsWaitGroup.Wait()
		close(results)
	}()

	return results
}

func (s *IAMRoleServer) GetMultiple(in *nativeIAmRoleGRPC.GetMultipleRolesRequest, out nativeIAmRoleGRPC.IAMRoleService_GetMultipleServer) error {
	ctx := out.Context()

	results := s.OpenGetMultipleChannel(ctx, in)

	var err error = nil
	for result := range results {
		if err != nil {
			// We have error. We dont need other data
			continue
		}
		if result.Err != nil {
			err = result.Err
			continue
		}

		err = out.Send(&nativeIAmRoleGRPC.GetMultipleRolesResponse{
			Role: result.Role.ToGRPCRole(result.Namespace),
		})
	}
	if err != nil {
		return status.Error(codes.Internal, "Error while fetching roles from database. "+err.Error())
	}
	return status.Error(codes.OK, "")
}
func (s *IAMRoleServer) List(in *nativeIAmRoleGRPC.ListRolesRequest, out nativeIAmRoleGRPC.IAMRoleService_ListServer) error {
	collection := s.getCollectionByNamespace(in.Namespace)
	options := options.Find().SetSort(bson.D{{Key: "_id", Value: 1}})
	if in.Skip != 0 {
		options.SetSkip(int64(in.Skip))
	}
	if in.Limit != 0 {
		options.SetLimit(int64(in.Limit))
	}
	cursor, err := collection.Find(out.Context(), bson.M{}, options)
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return status.Error(codes.OK, "Namespace doesnt exist")
			}
		}
		return status.Error(codes.Internal, err.Error())
	}
	defer cursor.Close(out.Context())

	for cursor.Next(out.Context()) {
		var role roleInMongo
		err = cursor.Decode(&role)
		if err != nil {
			return status.Error(codes.Internal, "error while decoding role from mongo:"+err.Error())
		}

		err = out.Send(&nativeIAmRoleGRPC.ListRolesResponse{
			Role: role.ToGRPCRole(in.Namespace),
		})
		if err != nil {
			return status.Error(codes.Internal, "error while sending role:"+err.Error())
		}
	}

	return status.Error(codes.OK, "")
}
func (s *IAMRoleServer) Count(ctx context.Context, in *nativeIAmRoleGRPC.CountRolesRequest) (*nativeIAmRoleGRPC.CountRolesResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makeCountRoleCacheKey(in.Namespace)
		byteData, _ := s.systemStub.Cache.Get(ctx, cacheKey)
		if byteData != nil {
			var response nativeIAmRoleGRPC.CountRolesResponse
			err := proto.Unmarshal(byteData, &response)
			if err != nil {
				return nil, status.Error(codes.Internal, "Error while unmarshaling roles count response from cache: "+err.Error())
			}
			return &response, status.Error(codes.OK, "")
		}
	}

	collection := s.getCollectionByNamespace(in.Namespace)
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &nativeIAmRoleGRPC.CountRolesResponse{
					Count: 0,
				}, status.Error(codes.OK, "Namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "error while counting documents in mongo:"+err.Error())
	}

	response := nativeIAmRoleGRPC.CountRolesResponse{
		Count: uint64(count),
	}

	if in.UseCache {
		responseBytes, err := proto.Marshal(&response)
		if err != nil {
			return nil, status.Error(codes.Internal, "Error while marshaling policy count to cache: "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, responseBytes, ROLE_CACHE_TIMEOUT)
	}

	return &response, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) Update(ctx context.Context, in *nativeIAmRoleGRPC.UpdateRoleRequest) (*nativeIAmRoleGRPC.UpdateRoleResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Role wasnt founded. Role UUID has bad format.")
	}

	collection := s.getCollectionByNamespace(in.Namespace)
	var updatedRole roleInMongo
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"name": in.NewName, "description": in.NewDescription,
			},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedRole)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while updating role in database. "+err.Error())
	}

	log.Infof("Updated role with UUID [%s] from the [%s] namespace", in.Uuid, in.Namespace)
	// s.systemStub.Cache.Remove(ctx, makeCountRoleCacheKey(in.Namespace))

	return &nativeIAmRoleGRPC.UpdateRoleResponse{
		Role: updatedRole.ToGRPCRole(in.Namespace),
	}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) Delete(ctx context.Context, in *nativeIAmRoleGRPC.DeleteRoleRequest) (*nativeIAmRoleGRPC.DeleteRoleResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return &nativeIAmRoleGRPC.DeleteRoleResponse{Existed: false}, status.Error(codes.OK, "Role wasnt founded. Role UUID has bad format")
	}

	collection := s.getCollectionByNamespace(in.Namespace)
	r, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &nativeIAmRoleGRPC.DeleteRoleResponse{Existed: false}, status.Error(codes.OK, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while deleting role in database. "+err.Error())
	}

	if r.DeletedCount != 0 {
		log.Infof("Deleted role with UUID [%s] from the [%s] namespace", in.Uuid, in.Namespace)
		s.systemStub.Cache.Remove(ctx, makeCountRoleCacheKey(in.Namespace))
	}

	return &nativeIAmRoleGRPC.DeleteRoleResponse{Existed: r.DeletedCount != 0}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) GetServiceManagedRole(ctx context.Context, in *nativeIAmRoleGRPC.GetServiceManagedRoleRequest) (*nativeIAmRoleGRPC.GetServiceManagedRoleResponse, error) {
	collection := s.getCollectionByNamespace(in.Namespace)
	var role roleInMongo
	err := collection.FindOne(
		ctx,
		bson.M{
			"managed._managementType":     role_managed_service,
			"managed.serviceName":         in.Service,
			"managed.serviceManagementId": in.ManagedId,
		},
		options.FindOne().SetHint(fast_search_service_index),
	).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Role wasnt founded.")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while searhing for role in database. "+err.Error())
	}

	return &nativeIAmRoleGRPC.GetServiceManagedRoleResponse{Role: role.ToGRPCRole(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) GetBuiltInRole(ctx context.Context, in *nativeIAmRoleGRPC.GetBuiltInRoleRequest) (*nativeIAmRoleGRPC.GetBuiltInRoleResponse, error) {
	builtInType := role_builtin_empty
	switch in.Type {
	case nativeIAmRoleGRPC.BuiltInRoleType_NAMESPACE_ROOT:
		builtInType = role_builtin_namespace_root
	case nativeIAmRoleGRPC.BuiltInRoleType_GLOBAL_ROOT:
		builtInType = role_builtin_global_root
	case nativeIAmRoleGRPC.BuiltInRoleType_EMPTY:
		builtInType = role_builtin_empty
	}

	collection := s.getCollectionByNamespace(in.Namespace)
	var role roleInMongo
	err := collection.FindOne(
		ctx,
		bson.M{
			"managed._managementType": role_managed_builtin,
			"managed.builtInType":     builtInType,
		},
		options.FindOne().SetHint(fast_search_built_in_index),
	).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Role wasnt founded.")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while searhing for role in database. "+err.Error())
	}
	return &nativeIAmRoleGRPC.GetBuiltInRoleResponse{Role: role.ToGRPCRole(in.Namespace)}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) AddPolicy(ctx context.Context, in *nativeIAmRoleGRPC.AddPolicyRequest) (*nativeIAmRoleGRPC.AddPolicyResponse, error) {
	roleUUID, err := primitive.ObjectIDFromHex(in.RoleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Role UUID has bad format.")
	}

	// Check if policy exists. This is not atomic but helps to avoid some errors
	existResponse, err := s.policyServer.Exist(ctx, &policy.ExistPolicyRequest{
		Namespace: in.PolicyNamespace,
		Uuid:      in.PolicyUUID,
		UseCache:  true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to check existance of policy: "+err.Error())
	}
	if !existResponse.Exist {
		return nil, status.Error(codes.FailedPrecondition, "Policy doesnt exist")
	}

	collection := s.getCollectionByNamespace(in.RoleNamespace)
	var role roleInMongo
	policy := &assignedPolicyInMongo{
		Namespace: in.PolicyNamespace,
		UUID:      in.PolicyUUID,
	}
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": roleUUID},
		bson.M{
			"$addToSet":    bson.M{"policies": policy},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Role wasnt founded.")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while searhing for role in database. "+err.Error())
	}

	return &nativeIAmRoleGRPC.AddPolicyResponse{Role: role.ToGRPCRole(in.RoleNamespace)}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) RemovePolicy(ctx context.Context, in *nativeIAmRoleGRPC.RemovePolicyRequest) (*nativeIAmRoleGRPC.RemovePolicyResponse, error) {
	roleUUID, err := primitive.ObjectIDFromHex(in.RoleUUID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Role UUID has bad format.")
	}

	collection := s.getCollectionByNamespace(in.RoleNamespace)
	var role roleInMongo
	policy := &assignedPolicyInMongo{
		Namespace: in.PolicyNamespace,
		UUID:      in.PolicyUUID,
	}
	err = collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": roleUUID},
		bson.M{
			"$pull":        bson.M{"policies": policy},
			"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
			"$inc":         bson.M{"version": 1},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Role wasnt founded.")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(codes.NotFound, "Role wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(codes.Internal, "Error while searhing for role in database. "+err.Error())
	}

	return &nativeIAmRoleGRPC.RemovePolicyResponse{Role: role.ToGRPCRole(in.RoleNamespace)}, status.Error(codes.OK, "")
}
func (s *IAMRoleServer) Exist(ctx context.Context, in *nativeIAmRoleGRPC.ExistRoleRequest) (*nativeIAmRoleGRPC.ExistRoleResponse, error) {
	uuid, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Role UUID has bad format.")
	}

	collection := s.getCollectionByNamespace(in.Namespace)
	count, err := collection.CountDocuments(ctx, bson.M{"_id": uuid}, options.Count().SetLimit(1))
	if err != nil {
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return &nativeIAmRoleGRPC.ExistRoleResponse{Exist: false}, status.Error(codes.OK, "")
			}
		}
		return nil, status.Error(codes.Internal, "Error while searhing for role in database. "+err.Error())
	}

	return &nativeIAmRoleGRPC.ExistRoleResponse{Exist: count != 0}, status.Error(codes.OK, "")
}
