package kanban

import (
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
)

var stageCollectionName = "crm_native_kanban_stage"
var ticketCollectionName = "crm_native_kanban_ticket"

func GetStageCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(stageCollectionName)
}

func GetTicketCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(ticketCollectionName)
}
