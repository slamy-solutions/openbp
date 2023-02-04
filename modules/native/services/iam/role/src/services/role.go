package services

import (
	"context"
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/status"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"google.golang.org/grpc/codes"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmRoleGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
)

type IAMRoleServer struct {
	nativeIAmRoleGRPC.UnimplementedIAMRoleServiceServer

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
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

func NewIAMRoleServer(ctx context.Context, systemStub *system.SystemStub, nativeStub *native.NativeStub) (*IAMRoleServer, error) {

	// Ensure indexes for global namespace
	err := ensureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("Failed to ensure indexes for global namespace: " + err.Error())
	}

	// Ensure BuiltIn roles for global namespace
	err = ensureBuiltInsForNamespace(ctx, "", systemStub, nativeStub)
	if err != nil {
		return nil, errors.New("Failed to ensure builtIns for global namespace: " + err.Error())
	}

	return &IAMRoleServer{
		systemStub: systemStub,
		nativeStub: nativeStub,
	}, nil
}

//TODO: add cache

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
			bson.M{"managed": bson.M{
				"serviceName":         role.Managed.ServiceName,
				"serviceManagementId": role.Managed.ServiceManagementID,
			}},
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

	log.Info("Created new role with UUID [%s] and name [%s] in the [%s] namespace", role.ID.Hex(), in.Name, in.Namespace)

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
func (s *IAMRoleServer) GetMultiple(in *nativeIAmRoleGRPC.GetMultipleRolesRequest, out nativeIAmRoleGRPC.IAMRoleService_GetMultipleServer) error {
	ctx := out.Context()

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
		err       error
		namespace string
		role      roleInMongo
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
					err       error
					namespace string
					role      roleInMongo
				}{err: err, namespace: "", role: roleInMongo{}}
				return
			}
			defer cursor.Close(ctx)

			for cursor.Next(ctx) {
				var role roleInMongo
				if err := cursor.Decode(role); err != nil {
					results <- struct {
						err       error
						namespace string
						role      roleInMongo
					}{err: err, namespace: "", role: roleInMongo{}}
					return
				}
				results <- struct {
					err       error
					namespace string
					role      roleInMongo
				}{err: nil, namespace: namespace, role: role}
			}

			if err := cursor.Err(); err != nil {
				results <- struct {
					err       error
					namespace string
					role      roleInMongo
				}{err: err, namespace: "", role: roleInMongo{}}
			}
		}(namespace, roles)
	}

	// Wait for all requests to finish and close results channel
	go func() {
		requestsWaitGroup.Wait()
		close(results)
	}()

	var err error = nil
	for result := range results {
		if err != nil {
			// We have error. We dont need other data
			continue
		}
		if result.err != nil {
			err = result.err
			continue
		}

		err = out.Send(&nativeIAmRoleGRPC.GetMultipleRolesResponse{
			Role: result.role.ToGRPCRole(result.namespace),
		})
	}
	if err != nil {
		return status.Error(codes.Internal, "Error while fetching roles from database. "+err.Error())
	}
	return status.Error(codes.OK, "")
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
		log.Info("Deleted role with UUID [%s] from the [%s] namespace", in.Uuid, in.Namespace)
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
	existResponse, err := s.nativeStub.Services.IamPolicy.Exist(ctx, &policy.ExistPolicyRequest{
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
