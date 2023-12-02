package project

import (
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
)

var projectCollectionName = "crm_native_project"

func GetProjectCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(projectCollectionName)
}
