package services

import (
	"context"

	log "github.com/sirupsen/logrus"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	fast_search_built_in_index = "fast_search_built_in"
	fast_search_service_index  = "fast_search_service"
)

func ensureIndexesForNamespace(ctx context.Context, namespace string, systemStub *system.SystemStub) error {
	collection := systemStub.DB.Database("openbp_global").Collection("native_iam_policy")
	if namespace != "" {
		collection = systemStub.DB.Database("openbp_namespace_" + namespace).Collection("native_iam_policy")
	}

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		// Fast search built-in policies
		{
			Keys: bson.D{bson.E{Key: "managed.builtInType", Value: "hashed"}},
			Options: options.Index().
				SetName(fast_search_built_in_index).
				SetPartialFilterExpression(bson.M{
					"managed._managementType": policy_managed_builtin,
				}),
		},
		// Fast search service-managed policies
		{
			Keys: bson.D{
				bson.E{Key: "managed.serviceName", Value: "hashed"},
				bson.E{Key: "managed.serviceManagementId", Value: 1},
			},
			Options: options.Index().
				SetName(fast_search_service_index).
				SetPartialFilterExpression(bson.M{
					"managed._managementType": policy_managed_service,
				}),
		},
	})

	if err == nil {
		log.Info("Ensured indexes for the [" + namespace + "] namespace.")
	}

	return err
}
