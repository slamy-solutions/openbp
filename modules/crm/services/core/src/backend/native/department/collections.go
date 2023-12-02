package department

import (
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.mongodb.org/mongo-driver/mongo"
)

var departmentCollectionName = "crm_native_department"

func GetDepartmentCollection(systemStub *system.SystemStub, namespace string) *mongo.Collection {
	dbName := "openbp_global"
	if namespace != "" {
		dbName = "openbp_" + namespace
	}

	return systemStub.DB.Database(dbName).Collection(departmentCollectionName)
}
