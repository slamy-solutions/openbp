package services

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newEmptyRole() *roleInMongo {
	creationTime := time.Now().UTC()
	return &roleInMongo{
		Name:        "Empty",
		Description: "Empty role with zero policies",
		Managed:     &managementTypeInMongo{ManagementType: role_managed_builtin, BuiltInType: role_builtin_empty},
		Policies:    []*assignedPolicyInMongo{},
		Tags:        []string{},
		Created:     creationTime,
		Updated:     creationTime,
		Version:     0,
	}
}

func newNamespaceRootRole(namespace string, rootPolicyUUID string) *roleInMongo {
	creationTime := time.Now().UTC()
	return &roleInMongo{
		Name:        "Namespace Root",
		Description: "Root role with full access to the namespace where it is defined",
		Managed:     &managementTypeInMongo{ManagementType: role_managed_builtin, BuiltInType: role_builtin_namespace_root},
		Policies: []*assignedPolicyInMongo{
			{
				Namespace: namespace,
				UUID:      rootPolicyUUID,
			},
		},
		Tags:    []string{},
		Created: creationTime,
		Updated: creationTime,
		Version: 0,
	}
}

func newGlobalRootRole(namespace string, rootPolicyUUID string) *roleInMongo {
	creationTime := time.Now().UTC()
	return &roleInMongo{
		Name:        "Global Root",
		Description: "Root role with full access to everything",
		Managed:     &managementTypeInMongo{ManagementType: role_managed_builtin, BuiltInType: role_builtin_global_root},
		Policies: []*assignedPolicyInMongo{
			{
				Namespace: namespace,
				UUID:      rootPolicyUUID,
			},
		},
		Tags:    []string{},
		Created: creationTime,
		Updated: creationTime,
		Version: 0,
	}
}

func ensureBuiltInsForNamespace(ctx context.Context, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub) error {
	builtInRolesCreationContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	collection := systemStub.DB.Database("openbp_global").Collection("native_iam_role")
	if namespace != "" {
		collection = systemStub.DB.Database("openbp_namespace_" + namespace).Collection("native_iam_role")
	}

	// EMPTY
	emptyInsertResult, err := collection.UpdateOne(
		builtInRolesCreationContext,
		bson.M{"managed": bson.M{"_managementType": role_managed_builtin, "builtInType": role_builtin_empty}},
		bson.M{"$setOnInsert": newEmptyRole()},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.New("Failed to create \"empty\" role for the [" + namespace + "] namespace. Failed to insert role to the database. " + err.Error())
	}
	if emptyInsertResult.UpsertedCount != 0 {
		log.Info("Created \"empty\" builtIn role for the [" + namespace + "] namespace.")
	}

	//NAMESPACE_ROOT
	// find built-in policy with root namespace access
	rootPolicyResponse, err := nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
		Namespace: namespace,
		Type:      policy.BuiltInPolicyType_NAMESPACE_ROOT,
	})
	if err != nil {
		return errors.New("Failed to create \"namespace root\" role for the [" + namespace + "] namespace. Failed to search for \"namespace root\" policy. " + err.Error())
	}
	namespaceRootInsertResult, err := collection.UpdateOne(
		builtInRolesCreationContext,
		bson.M{"managed": bson.M{"_managementType": role_managed_builtin, "builtInType": role_builtin_namespace_root}},
		bson.M{"$setOnInsert": newNamespaceRootRole(namespace, rootPolicyResponse.Policy.Uuid)},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.New("Failed to create \"namespace root\" role for the [" + namespace + "] namespace. Failed to insert role to the database. " + err.Error())
	}
	if namespaceRootInsertResult.UpsertedCount != 0 {
		log.Info("Created \"namespace root\" builtIn role for the [" + namespace + "] namespace.")
	}

	//GLOBAL_ROOT
	if namespace == "" {
		rootPolicyResponse, err := nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
			Namespace: namespace,
			Type:      policy.BuiltInPolicyType_GLOBAL_ROOT,
		})
		if err != nil {
			return errors.New("Failed to create \"global root\" role for the [" + namespace + "] namespace. Failed to search for \"global root\" policy. " + err.Error())
		}
		globalRootInsertResult, err := collection.UpdateOne(
			builtInRolesCreationContext,
			bson.M{"managed": bson.M{"_managementType": role_managed_builtin, "builtInType": role_builtin_global_root}},
			bson.M{"$setOnInsert": newGlobalRootRole(namespace, rootPolicyResponse.Policy.Uuid)},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return errors.New("Failed to create \"global root\" role for the [" + namespace + "] namespace. Failed to insert role to the database. " + err.Error())
		}
		if globalRootInsertResult.UpsertedCount != 0 {
			log.Info("Created \"global root\" builtIn role for the [" + namespace + "] namespace.")
		}
	}

	return nil
}
