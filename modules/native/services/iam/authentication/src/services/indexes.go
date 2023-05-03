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
	fast_search_password_identity_index = "fast_search_identity"
)

func EnsureIndexesForNamespace(ctx context.Context, namespace string, systemStub *system.SystemStub) error {
	collection := systemStub.DB.Database("openbp_global").Collection("native_iam_authentication_password")
	if namespace != "" {
		collection = systemStub.DB.Database("openbp_namespace_" + namespace).Collection("native_iam_authentication_password")
	}

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "identity", Value: "hashed"},
		},
		Options: options.Index().
			SetName(fast_search_password_identity_index),
	},
	)
	if err != nil {
		log.WithField("namespace", namespace).Error("Failed to ensure indexes for password: " + err.Error())
		return err
	}

	if err == nil {
		log.WithField("namespace", namespace).Info("Successfully ensured indexes for namespace.")
	}

	return err
}
