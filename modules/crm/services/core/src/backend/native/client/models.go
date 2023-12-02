package client

import (
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientInMongo struct {
	UUID primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`

	LastUpdateTime time.Time `bson:"lastUpdateTime"`
	CreationTime   time.Time `bson:"creationTime"`
	Version        int64     `bson:"version"`
}

func (c *ClientInMongo) ToBackendModel(namespace string, contactPersons []models.ContactPerson) *models.Client {
	return &models.Client{
		Namespace: namespace,
		UUID:      c.UUID.Hex(),
		Name:      c.Name,

		ContactPersons: contactPersons,

		LastUpdateTime: c.LastUpdateTime,
		CreationTime:   c.CreationTime,
		Version:        c.Version,
	}
}

type ContactPersonInMongo struct {
	UUID        primitive.ObjectID `bson:"_id,omitempty"`
	ClientUUID  primitive.ObjectID `bson:"clientUUID"`
	Name        string             `bson:"name"`
	Email       string             `bson:"email"`
	Phone       []string           `bson:"phone"`
	NotRelevant bool               `bson:"notRelevant"`
	Comment     string             `bson:"comment"`
}

func (c *ContactPersonInMongo) ToBackendModel(namespace string) *models.ContactPerson {
	return &models.ContactPerson{
		Namespace:   namespace,
		UUID:        c.UUID.Hex(),
		ClientUUID:  c.ClientUUID.Hex(),
		Name:        c.Name,
		Email:       c.Email,
		Phone:       c.Phone,
		NotRelevant: c.NotRelevant,
		Comment:     c.Comment,
	}
}
