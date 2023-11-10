package client

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clientInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	LastUpdateTime time.Time `bson:"lastUpdateTime"`
	CreationTime   time.Time `bson:"creationTime"`
	Version        int64     `bson:"version"`
}

type contactPersonInMongo struct {
	UUID        primitive.ObjectID `bson:"_id,omitempty"`
	ClientUUID  primitive.ObjectID `bson:"clientUUID"`
	Name        string             `bson:"name"`
	Email       string             `bson:"email"`
	Phone       []string           `bson:"phone"`
	NotRelevant bool               `bson:"notRelevant"`
	Comment     string             `bson:"comment"`
}
