package repositories

import (
	"context"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	catalogCollectionName            = "catalogs"
	catalogEntryCollectionNamePrefix = "catalog_"
)

func getDadatabeByNamespace(system *system.SystemStub, namespace string) *mongo.Database {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return system.DB.Database(dbName)
}

func getCatalogCollection(system *system.SystemStub, namespace string) *mongo.Collection {
	db := getDadatabeByNamespace(system, namespace)
	return db.Collection(catalogCollectionName)
}

func createEntryCollection(ctx context.Context, system *system.SystemStub, namespace string, catalogName string) error {
	db := getDadatabeByNamespace(system, namespace)
	return db.CreateCollection(ctx, catalogEntryCollectionNamePrefix+catalogName)
}

func dropEntryCollection(ctx context.Context, system *system.SystemStub, namespace string, catalogName string) error {
	db := getDadatabeByNamespace(system, namespace)
	return db.Collection(catalogEntryCollectionNamePrefix + catalogName).Drop(ctx)
}

func getCatalogEntryCollection(system *system.SystemStub, namespace string, catalogName string) *mongo.Collection {
	db := getDadatabeByNamespace(system, namespace)
	return db.Collection(catalogEntryCollectionNamePrefix + catalogName)
}
