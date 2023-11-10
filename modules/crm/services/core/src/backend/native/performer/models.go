package performer

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type performerInMongo struct {
	UUID           primitive.ObjectID `bson:"_id,omitempty"`
	DepartmentUUID primitive.ObjectID `bson:"departmentUUID"`
	UserUUID       primitive.ObjectID `bson:"userUUID"`
}
