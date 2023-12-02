package project

import (
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	ClientUUID     primitive.ObjectID `bson:"clientUUID"`
	ContactUUID    primitive.ObjectID `bson:"contactUUID"`
	DepartmentUUID primitive.ObjectID `bson:"departmentUUID"`

	NotRelevant bool `bson:"notRelevant"`
}

func (p *ProjectInMongo) ToBackendModel(namespace string) *models.Project {
	return &models.Project{
		Namespace: namespace,
		UUID:      p.UUID.Hex(),
		Name:      p.Name,

		ClientUUID:     p.ClientUUID.Hex(),
		ContactUUID:    p.ContactUUID.Hex(),
		DepartmentUUID: p.DepartmentUUID.Hex(),

		NotRelevant: p.NotRelevant,
	}
}
