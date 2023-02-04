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
	unique_login_index         = "unique_login"
	fast_identity_search_index = "fast_identity_search"
)

func ensureIndexesForNamespace(ctx context.Context, namespace string, systemStub *system.SystemStub) error {
	collection := systemStub.DB.Database("openbp_global").Collection("native_actor_user")
	if namespace != "" {
		collection = systemStub.DB.Database("openbp_namespace_" + namespace).Collection("native_actor_user")
	}

	_, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Options: options.Index().SetName(unique_login_index).SetUnique(true),
			Keys:    bson.D{bson.E{Key: "login", Value: 1}},
		},
		{
			Options: options.Index().SetName(fast_identity_search_index),
			Keys:    bson.D{bson.E{Key: "identity", Value: "hashed"}},
		},
	})

	if err == nil {
		log.Info("Ensured indexes for the [" + namespace + "] namespace.")
	}

	return err
}
