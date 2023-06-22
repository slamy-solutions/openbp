package policy

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newEmptyPolicy() *policyInMongo {
	creationTime := time.Now().UTC()
	return &policyInMongo{
		Name:                 "Empty",
		Description:          "Empty policy with zero access",
		Managed:              &managementTypeInMongo{ManagementType: policy_managed_builtin, BuiltInType: policy_builtin_empty},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
		Tags:                 []string{},
		Created:              creationTime,
		Updated:              creationTime,
		Version:              0,
	}
}

func newNamespaceRootPolicy() *policyInMongo {
	creationTime := time.Now().UTC()
	return &policyInMongo{
		Name:                 "Namespace root",
		Description:          "Root policy with full access to the namespace where it is defined",
		Managed:              &managementTypeInMongo{ManagementType: policy_managed_builtin, BuiltInType: policy_builtin_namespace_root},
		NamespaceIndependent: false,
		Resources:            []string{"*"},
		Actions:              []string{"*"},
		Tags:                 []string{},
		Created:              creationTime,
		Updated:              creationTime,
		Version:              0,
	}
}

func newGlobalRootPolicy() *policyInMongo {
	creationTime := time.Now().UTC()
	return &policyInMongo{
		Name:                 "Global Root",
		Description:          "Root policy with full access to everything",
		Managed:              &managementTypeInMongo{ManagementType: policy_managed_builtin, BuiltInType: policy_builtin_global_root},
		NamespaceIndependent: true,
		Resources:            []string{"*"},
		Actions:              []string{"*"},
		Tags:                 []string{},
		Created:              creationTime,
		Updated:              creationTime,
		Version:              0,
	}
}

func ensureBuiltInsForNamespace(ctx context.Context, namespace string, systemStub *system.SystemStub) error {
	builtInPolicyCreationContext, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	collection := systemStub.DB.Database("openbp_global").Collection("native_iam_policy")
	if namespace != "" {
		collection = systemStub.DB.Database("openbp_namespace_" + namespace).Collection("native_iam_policy")
	}

	r, err := collection.UpdateOne(
		builtInPolicyCreationContext,
		bson.M{"managed._managementType": policy_managed_builtin, "managed.builtInType": policy_builtin_empty},
		bson.M{"$setOnInsert": newEmptyPolicy()},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.New("Failed to create \"empty\" policy for [" + namespace + "] namespace: " + err.Error())
	}
	if r.UpsertedCount != 0 {
		log.Info("Created \"empty\" builtIn policy for the [" + namespace + "] namespace.")
	}

	r, err = collection.UpdateOne(
		builtInPolicyCreationContext,
		bson.M{"managed._managementType": policy_managed_builtin, "managed.builtInType": policy_builtin_namespace_root},
		bson.M{"$setOnInsert": newNamespaceRootPolicy()},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return errors.New("Failed to create \"namespace root\" policy for [" + namespace + "] namespace: " + err.Error())
	}
	if r.UpsertedCount != 0 {
		log.Info("Created \"namespace root\" builtIn policy for the [" + namespace + "] namespace.")
	}

	if namespace == "" {
		r, err = collection.UpdateOne(
			builtInPolicyCreationContext,
			bson.M{"managed._managementType": policy_managed_builtin, "managed.builtInType": policy_builtin_global_root},
			bson.M{"$setOnInsert": newGlobalRootPolicy()},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return errors.New("Failed to create \"global root\" policy for global namespace: " + err.Error())
		}
		if r.UpsertedCount != 0 {
			log.Info("Created \"global root\" builtIn policy for the global namespace.")
		}
	}

	return nil
}
