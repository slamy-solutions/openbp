package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeNamespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type IAMPolicyServer struct {
	nativeIAmPolicyGRPC.UnimplementedIAMPolicyServiceServer

	mongoGlobalCollection *mongo.Collection

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

const (
	POLICY_CACHE_TIMEOUT         = time.Second * 30
	BUILTIN_POLICY_CACHE_TIMEOUT = time.Second * 30
)

func (s *IAMPolicyServer) collectionByNamespace(namespace string) *mongo.Collection {
	if namespace == "" {
		return s.mongoGlobalCollection
	} else {
		db := s.systemStub.DB.Database(fmt.Sprintf("openbp_namespace_%s", namespace))
		return db.Collection("native_iam_policy")
	}
}

func makePolicyCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("native_iam_policy_data_%s_%s", namespace, uuid)
}
func makeBuiltInPolicyCacheKey(builtInType string) string {
	return fmt.Sprintf("native_iam_policy_builtin_%s", builtInType)
}

func NewIAMPolicyServer(ctx context.Context, systemStub *system.SystemStub, nativeStub *native.NativeStub) (*IAMPolicyServer, error) {
	err := ensureIndexesForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("Failed to ensure indexes for global namespace. " + err.Error())
	}

	err = ensureBuiltInsForNamespace(ctx, "", systemStub)
	if err != nil {
		return nil, errors.New("Failed to initialize builtin policies for global namespace. " + err.Error())
	}

	mongoGlobalCollection := systemStub.DB.Database("openbp_global").Collection("native_iam_policy")
	return &IAMPolicyServer{
		mongoGlobalCollection: mongoGlobalCollection,
		systemStub:            systemStub,
		nativeStub:            nativeStub,
	}, nil
}

func (s *IAMPolicyServer) Create(ctx context.Context, in *nativeIAmPolicyGRPC.CreatePolicyRequest) (*nativeIAmPolicyGRPC.CreatePolicyResponse, error) {
	err := validatePolicyData(in.Name, in.Description, in.Actions, in.Resources)
	if err != nil {
		return nil, err
	}

	// Check if namespace really exist
	if in.Namespace != "" {
		r, err := s.nativeStub.Services.Namespace.Exists(ctx, &nativeNamespaceGRPC.IsNamespaceExistRequest{Name: in.Namespace, UseCache: true})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		if !r.Exist {
			return nil, status.Error(grpccodes.FailedPrecondition, "Namespace doesnt exist")
		}
	}

	var managed managementTypeInMongo
	managed.ManagementType = policy_managed_no
	if t, ok := in.Managed.(*nativeIAmPolicyGRPC.CreatePolicyRequest_Identity); ok {
		managed.ManagementType = policy_managed_identity
		managed.IdentityNamespace = t.Identity.IdentityNamespace
		managed.IdentityUUID = t.Identity.IdentityUUID
	}
	if t, ok := in.Managed.(*nativeIAmPolicyGRPC.CreatePolicyRequest_Service); ok {
		managed.ManagementType = policy_managed_service
		managed.ServiceName = t.Service.Service
		managed.ServiceReason = t.Service.Reason
		managed.ServiceManagementID = t.Service.ManagementId
	}

	creationTime := time.Now().UTC()

	policy := policyInMongo{
		Name:        in.Name,
		Description: in.Description,

		Managed: &managed,

		NamespaceIndependent: in.NamespaceIndependent,
		Actions:              in.Actions,
		Resources:            in.Resources,

		Tags: []string{},

		Created: creationTime,
		Updated: creationTime,
		Version: 0,
	}

	collection := s.collectionByNamespace(in.Namespace)

	if policy.Managed.ManagementType == policy_managed_service && policy.Managed.ServiceManagementID != "" {
		r, err := collection.UpdateOne(
			ctx,
			bson.M{"managed": bson.M{
				"serviceName":         policy.Managed.ServiceName,
				"serviceManagementId": policy.Managed.ServiceManagementID,
			}},
			bson.M{
				"$setOnInsert": policy,
			},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while adding policy to the database. "+err.Error())
		}
		if r.UpsertedCount == 0 {
			return nil, status.Error(grpccodes.AlreadyExists, "Policy managed by this service with same managementId already exists.")
		}
		policy.UUID = r.UpsertedID.(primitive.ObjectID)
	} else {
		insertResponse, err := collection.InsertOne(ctx, policy)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to insert policy to the database. "+err.Error())
		}
		policy.UUID = insertResponse.InsertedID.(primitive.ObjectID)
	}

	return &nativeIAmPolicyGRPC.CreatePolicyResponse{
		Policy: policy.ToGRPCPolicy(in.Namespace),
	}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) Get(ctx context.Context, in *nativeIAmPolicyGRPC.GetPolicyRequest) (*nativeIAmPolicyGRPC.GetPolicyResponse, error) {
	var cacheKey string
	if in.UseCache {
		cacheKey = makePolicyCacheKey(in.Namespace, in.Uuid)
		byteData, _ := s.systemStub.Cache.Get(ctx, cacheKey)
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

	collection := s.collectionByNamespace(in.Namespace)
	var mongoPolicy policyInMongo
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mongoPolicy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Policy with provided UUID and namespace not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "Policy wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	policy := mongoPolicy.ToGRPCPolicy(in.Namespace)

	if in.UseCache {
		policyBytes, err := proto.Marshal(policy)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling policy to cache: "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, policyBytes, POLICY_CACHE_TIMEOUT)
	}

	return &nativeIAmPolicyGRPC.GetPolicyResponse{Policy: policy}, status.Error(grpccodes.OK, "")
}
func (s *IAMPolicyServer) GetMultiple(in *nativeIAmPolicyGRPC.GetMultiplePoliciesRequest, out nativeIAmPolicyGRPC.IAMPolicyService_GetMultipleServer) error {
	ctx := out.Context()

	// Split all the policies into namespaces
	namespaces := make(map[string][]*nativeIAmPolicyGRPC.GetMultiplePoliciesRequest_RequestedPolicy)
	for _, policy := range in.Policies {
		policies, ok := namespaces[policy.Namespace]
		if !ok {
			policies = make([]*nativeIAmPolicyGRPC.GetMultiplePoliciesRequest_RequestedPolicy, 0, 1)
		}

		policies = append(policies, policy)
		namespaces[policy.Namespace] = policies
	}

	// For each namespace run mongo request
	var requestsWaitGroup sync.WaitGroup
	requestsWaitGroup.Add(len(namespaces))
	results := make(chan struct {
		err       error
		namespace string
		policy    policyInMongo
	})
	for namespace, policies := range namespaces {
		go func(namespace string, policies []*nativeIAmPolicyGRPC.GetMultiplePoliciesRequest_RequestedPolicy) {
			defer requestsWaitGroup.Done()

			uuids := bson.A{}
			for _, policy := range policies {
				id, err := primitive.ObjectIDFromHex(policy.Uuid)
				if err == nil {
					uuids = append(uuids, id)
				}
			}

			collection := s.collectionByNamespace(namespace)
			cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$in": uuids}})
			if err != nil {
				if err, ok := err.(mongo.WriteException); ok {
					if err.HasErrorLabel("InvalidNamespace") {
						return //Its ok. Policy not founded
					}
				}
				results <- struct {
					err       error
					namespace string
					policy    policyInMongo
				}{err: err, namespace: "", policy: policyInMongo{}}
				return
			}
			defer cursor.Close(ctx)

			for cursor.Next(ctx) {
				var policy policyInMongo
				if err := cursor.Decode(&policy); err != nil {
					results <- struct {
						err       error
						namespace string
						policy    policyInMongo
					}{err: err, namespace: "", policy: policyInMongo{}}
					return
				}
				results <- struct {
					err       error
					namespace string
					policy    policyInMongo
				}{err: nil, namespace: namespace, policy: policy}
			}

			if err := cursor.Err(); err != nil {
				results <- struct {
					err       error
					namespace string
					policy    policyInMongo
				}{err: err, namespace: "", policy: policyInMongo{}}
			}
		}(namespace, policies)
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

		err = out.Send(&nativeIAmPolicyGRPC.GetMultiplePoliciesResponse{
			Policy: result.policy.ToGRPCPolicy(result.namespace),
		})
	}
	if err != nil {
		return status.Error(grpccodes.Internal, "Error while fetching policies from database. "+err.Error())
	}
	return status.Error(grpccodes.OK, "")
}
func (s *IAMPolicyServer) Exist(ctx context.Context, in *nativeIAmPolicyGRPC.ExistPolicyRequest) (*nativeIAmPolicyGRPC.ExistPolicyResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Policy UUID has bad format")
	}

	var cacheKey string
	if in.UseCache {
		cacheKey = makePolicyCacheKey(in.Namespace, in.Uuid)
		exist, err := s.systemStub.Cache.Has(ctx, cacheKey)
		if err == nil && exist {
			return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: true}, status.Error(grpccodes.OK, "")
		}
	}

	collection := s.collectionByNamespace(in.Namespace)
	if in.UseCache {
		var mongoPolicy policyInMongo
		err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&mongoPolicy)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: false}, status.Error(grpccodes.OK, "")
			}
			if err, ok := err.(mongo.WriteException); ok {
				if err.HasErrorLabel("InvalidNamespace") {
					return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: false}, status.Error(grpccodes.OK, "Namespace doesnt exist.")
				}
			}
			return nil, status.Error(grpccodes.Internal, err.Error())
		}
		policy := mongoPolicy.ToGRPCPolicy(in.Namespace)
		policyBytes, err := proto.Marshal(policy)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while marshaling policy to cache: "+err.Error())
		}
		s.systemStub.Cache.Set(ctx, cacheKey, policyBytes, POLICY_CACHE_TIMEOUT)
		return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: true}, status.Error(grpccodes.OK, "")
	} else {
		// Fast check in mongo if exist without getting data
		count, err := collection.CountDocuments(ctx, bson.M{"_id": id}, options.Count().SetLimit(1))
		if err != nil {
			if err, ok := err.(mongo.WriteException); ok {
				if err.HasErrorLabel("InvalidNamespace") {
					return &nativeIAmPolicyGRPC.ExistPolicyResponse{Exist: false}, status.Error(grpccodes.OK, "Namespace doesnt exist.")
				}
			}
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

	collection := s.collectionByNamespace(in.Namespace)
	updateQuery := bson.M{
		"$set": bson.M{
			"name":                 in.Name,
			"description":          in.Description,
			"namespaceIndependent": in.NamespaceIndependent,
			"resources":            in.Resources,
			"actions":              in.Actions,
		},
		"$currentDate": bson.M{"updated": bson.M{"$type": "timestamp"}},
		"$inc":         bson.M{"version": 1},
	}

	var policy policyInMongo
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, updateQuery).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Policy with specified uuid not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "Policy with specified namespace not found")
			}
		}
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	s.systemStub.Cache.Remove(ctx, makePolicyCacheKey(in.Namespace, in.Uuid))

	return &nativeIAmPolicyGRPC.UpdatePolicyResponse{
		Policy: policy.ToGRPCPolicy(in.Namespace),
	}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) Delete(ctx context.Context, in *nativeIAmPolicyGRPC.DeletePolicyRequest) (*nativeIAmPolicyGRPC.DeletePolicyResponse, error) {
	id, err := primitive.ObjectIDFromHex(in.Uuid)
	if err != nil {
		return nil, status.Error(grpccodes.InvalidArgument, "Bad UUID format")
	}

	collection := s.collectionByNamespace(in.Namespace)
	r, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, err.Error())
	}

	if r.DeletedCount != 0 {
		s.systemStub.Cache.Remove(ctx, makePolicyCacheKey(in.Namespace, in.Uuid))
	}

	return &nativeIAmPolicyGRPC.DeletePolicyResponse{Existed: r.DeletedCount != 0}, status.Error(grpccodes.OK, "")
}

func (s *IAMPolicyServer) List(in *nativeIAmPolicyGRPC.ListPoliciesRequest, out nativeIAmPolicyGRPC.IAMPolicyService_ListServer) error {
	ctx := out.Context()

	collection := s.collectionByNamespace(in.Namespace)
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
			Policy: result.ToGRPCPolicy(in.Namespace),
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
func (s *IAMPolicyServer) GetServiceManagedPolicy(ctx context.Context, in *nativeIAmPolicyGRPC.GetServiceManagedPolicyRequest) (*nativeIAmPolicyGRPC.GetServiceManagedPolicyResponse, error) {
	// TODO: Add cache. Probably need to hash "service" and "managedid" to create index cause theirs values can be in any format

	collection := s.collectionByNamespace(in.Namespace)
	var policy policyInMongo
	err := collection.FindOne(ctx, bson.M{
		"managed._managementType":     policy_managed_service,
		"managed.serviceName":         in.Service,
		"managed.serviceManagementId": in.ManagedId,
	}, options.FindOne().SetHint(fast_search_service_index)).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Policy with provided service name, id and namespace not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "Policy wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Error while searching for policy in database. "+err.Error())
	}

	return &nativeIAmPolicyGRPC.GetServiceManagedPolicyResponse{Policy: policy.ToGRPCPolicy(in.Namespace)}, status.Error(grpccodes.OK, "")
}
func (s *IAMPolicyServer) GetBuiltInPolicy(ctx context.Context, in *nativeIAmPolicyGRPC.GetBuiltInPolicyRequest) (*nativeIAmPolicyGRPC.GetBuiltInPolicyResponse, error) {
	var builtInType string
	switch in.Type {
	case nativeIAmPolicyGRPC.BuiltInPolicyType_NAMESPACE_ROOT:
		builtInType = policy_builtin_namespace_root
	case nativeIAmPolicyGRPC.BuiltInPolicyType_GLOBAL_ROOT:
		builtInType = policy_builtin_global_root
	default:
		builtInType = policy_builtin_empty
	}

	cacheKey := makeBuiltInPolicyCacheKey(builtInType)

	// Try to get policy from the cache
	cacheByteData, _ := s.systemStub.Cache.Get(ctx, cacheKey)
	if cacheByteData != nil {
		var policy nativeIAmPolicyGRPC.Policy
		err := proto.Unmarshal(cacheByteData, &policy)
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Error while unmarshaling policy from cache: "+err.Error())
		}
		return &nativeIAmPolicyGRPC.GetBuiltInPolicyResponse{Policy: &policy}, status.Error(grpccodes.OK, "")
	}

	collection := s.collectionByNamespace(in.Namespace)

	// Get actual policy from the database
	var policy policyInMongo
	err := collection.FindOne(ctx, bson.M{
		"managed._managementType": policy_managed_builtin,
		"managed.builtInType":     builtInType,
	}, options.FindOne().SetHint(fast_search_built_in_index)).Decode(&policy)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(grpccodes.NotFound, "Policy with provided BuiltIn type and namespace not found")
		}
		if err, ok := err.(mongo.WriteException); ok {
			if err.HasErrorLabel("InvalidNamespace") {
				return nil, status.Error(grpccodes.NotFound, "Policy wasnt founded. Probably namespace doesnt exist")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Error while searching for policy in database. "+err.Error())
	}

	// Put policy to the cache
	policyGRPC := policy.ToGRPCPolicy(in.Namespace)
	policyBytes, err := proto.Marshal(policyGRPC)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error while marshaling policy to cache: "+err.Error())
	}
	s.systemStub.Cache.Set(ctx, cacheKey, policyBytes, BUILTIN_POLICY_CACHE_TIMEOUT)

	return &nativeIAmPolicyGRPC.GetBuiltInPolicyResponse{Policy: policyGRPC}, status.Error(grpccodes.OK, "")
}
