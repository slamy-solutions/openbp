package client

import (
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
)

var clientCollectionName = "crm_native_client"
var clientContactPersonCollectionName = "crm_native_client_contactperson"

func getClientCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(clientCollectionName)
}

func getClientContactPersonCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(clientContactPersonCollectionName)
}
