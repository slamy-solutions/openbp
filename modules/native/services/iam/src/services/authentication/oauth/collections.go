package oauth

import (
	"context"
	"errors"
	"fmt"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func configCollectionByNamespace(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	if namespace == "" {
		return systemStub.DB.Database("openbp_global").Collection("native_iam_authentication_oauth2_config")
	} else {
		dbName := fmt.Sprintf("openbp_namespace_%s", namespace)
		return systemStub.DB.Database(dbName).Collection("native_iam_authentication_oauth2_config")
	}
}

func registrationCollectionByNamespace(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	if namespace == "" {
		return systemStub.DB.Database("openbp_global").Collection("native_iam_authentication_oauth2_registration")
	} else {
		dbName := fmt.Sprintf("openbp_namespace_%s", namespace)
		return systemStub.DB.Database(dbName).Collection("native_iam_authentication_oauth2_registration")
	}
}

func EnsureIndexesForNamespace(ctx context.Context, namespace string, systemStub *system.SystemStub) error {
	collection := registrationCollectionByNamespace(systemStub, namespace)

	// TODO: sparse index for provider type

	indexModels := make([]mongo.IndexModel, 0, len(ProviderTypes))
	for _, providerType := range ProviderTypes {
		indexModels = append(indexModels, mongo.IndexModel{
			Keys:    bson.D{bson.E{Key: fmt.Sprintf("%s.id", providerType), Value: 1}},
			Options: options.Index().SetName(fmt.Sprintf("unique_%s_user_id", providerType)).SetUnique(true),
		})
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return errors.Join(errors.New("failed to ensure indexes for oauth2 registration"), err)
	}

	return nil
}
