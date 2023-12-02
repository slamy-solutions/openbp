package department

import (
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DepartmentInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}

func (d *DepartmentInMongo) ToBackendModel(namespace string) *models.Department {
	return &models.Department{
		Namespace: namespace,
		UUID:      d.UUID.Hex(),
		Name:      d.Name,
	}
}
