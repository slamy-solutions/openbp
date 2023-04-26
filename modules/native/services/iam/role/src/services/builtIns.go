package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
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

func ensureBuiltInsForNamespace(ctx context.Context, namespaceName string, systemStub *system.SystemStub, nativeStub *native.NativeStub) error {
	builtInRolesCreationContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	logger := log.WithFields(log.Fields{"namespace": namespaceName, "component": "builtins"})

	collection := systemStub.DB.Database("openbp_global").Collection("native_iam_role")
	if namespaceName != "" {
		collection = systemStub.DB.Database("openbp_namespace_" + namespaceName).Collection("native_iam_role")
	}

	// EMPTY
	emptyInsertResult, err := collection.UpdateOne(
		builtInRolesCreationContext,
		bson.M{"managed": bson.M{"_managementType": role_managed_builtin, "builtInType": role_builtin_empty}},
		bson.M{"$setOnInsert": newEmptyRole()},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.New("Failed to create \"empty\" role for the [" + namespaceName + "] namespace. Failed to insert role to the database. " + err.Error())
	}
	if emptyInsertResult.UpsertedCount != 0 {
		logger.Info("Created \"empty\" builtIn role.")
	}

	//NAMESPACE_ROOT
	// find built-in policy with root namespace access
	rootPolicyUUID := ""
	searchRetries := 10
	for {
		rootPolicyResponse, err := nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
			Namespace: namespaceName,
			Type:      policy.BuiltInPolicyType_NAMESPACE_ROOT,
		})
		if err != nil {
			if searchRetries >= 0 {
				// Check if namespace still exist. If not - we can fail.
				namespaceExistRespose, namespaceErr := nativeStub.Services.Namespace.Exists(ctx, &namespace.IsNamespaceExistRequest{
					Name:     namespaceName,
					UseCache: false,
				})
				if namespaceErr != nil {
					return errors.New("Failed to create \"namespace root\" role for the [" + namespaceName + "] namespace. Failed to search for \"namespace root\" policy. Failed to get information about namespace: " + namespaceErr.Error())
				}

				if namespaceExistRespose.Exist {
					logger.Warning("Failed to search for \"namespace root\" builtIn policy. Namespace still exist. Retry.")
					time.Sleep(time.Second)
					searchRetries -= 1
					continue
				} else {
					return errors.New("Failed to create \"namespace root\" role for the [" + namespaceName + "] namespace. Failed to search for \"namespace root\" policy because namespace doesnt exist.")
				}
			}
			return errors.New("Failed to create \"namespace root\" role for the [" + namespaceName + "] namespace. Failed to search for \"namespace root\" policy. " + err.Error())
		}
		rootPolicyUUID = rootPolicyResponse.Policy.Uuid
		break
	}
	namespaceRootInsertResult, err := collection.UpdateOne(
		builtInRolesCreationContext,
		bson.M{"managed": bson.M{"_managementType": role_managed_builtin, "builtInType": role_builtin_namespace_root}},
		bson.M{"$setOnInsert": newNamespaceRootRole(namespaceName, rootPolicyUUID)},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.New("Failed to create \"namespace root\" role for the [" + namespaceName + "] namespace. Failed to insert role to the database. " + err.Error())
	}
	if namespaceRootInsertResult.UpsertedCount != 0 {
		logger.Info("Created \"namespace root\" builtIn role.")
	}

	//GLOBAL_ROOT
	if namespaceName == "" {
		rootPolicyResponse, err := nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
			Namespace: namespaceName,
			Type:      policy.BuiltInPolicyType_GLOBAL_ROOT,
		})
		if err != nil {
			return errors.New("Failed to create \"global root\" role for the [" + namespaceName + "] namespace. Failed to search for \"global root\" policy. " + err.Error())
		}
		globalRootInsertResult, err := collection.UpdateOne(
			builtInRolesCreationContext,
			bson.M{"managed": bson.M{"_managementType": role_managed_builtin, "builtInType": role_builtin_global_root}},
			bson.M{"$setOnInsert": newGlobalRootRole(namespaceName, rootPolicyResponse.Policy.Uuid)},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return errors.New("Failed to create \"global root\" role for the [" + namespaceName + "] namespace. Failed to insert role to the database. " + err.Error())
		}
		if globalRootInsertResult.UpsertedCount != 0 {
			logger.Info("Created \"global root\" builtIn role.")
		}
	}

	fmt.Println("Buildins ensured")

	return nil
}
