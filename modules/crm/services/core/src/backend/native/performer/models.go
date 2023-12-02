package performer

import (
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PerformerInMongo struct {
	UUID           primitive.ObjectID `bson:"_id,omitempty"`
	DepartmentUUID primitive.ObjectID `bson:"departmentUUID"`
	UserUUID       primitive.ObjectID `bson:"userUUID"`
}

func (p *PerformerInMongo) ToBackendModel(namespace string, name string, avatarURL string) *models.Performer {
	return &models.Performer{
		Namespace:      namespace,
		UUID:           p.UUID.Hex(),
		Name:           name,
		AvatarURL:      avatarURL,
		DepartmentUUID: p.DepartmentUUID.Hex(),
		UserUUID:       p.UserUUID.Hex(),
	}
}
