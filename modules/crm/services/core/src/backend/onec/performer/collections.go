package performer

import (
	"context"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var performerCollectionName = "crm_onec_performer"

func getPerformerCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(performerCollectionName)
}

func makeIndexes(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	performerCollection := getPerformerCollection(systemStub, namespace)

	_, err := performerCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: "hashed"}},
			Options: options.Index().SetName("uuid"),
		},
		{
			Keys:    bson.D{{Key: "userUUID", Value: 1}},
			Options: options.Index().SetName("unique_userUUID").SetUnique(true),
		},
	})
	return err
}
