package models

import (
	"context"
	"errors"
	"time"

	client "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ContactPerson struct {
	Namespace   string   `json:"namespace"`
	UUID        string   `json:"uuid"`
	ClientUUID  string   `json:"clientUUID"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Phone       []string `json:"phone"`
	NotRelevant bool     `json:"notRelevant"`
	Comment     string   `json:"comment"`
}

type Client struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`

	ContactPersons []ContactPerson `json:"contactPersons"`

	LastUpdateTime time.Time `json:"lastUpdateTime"`
	CreationTime   time.Time `json:"creationTime"`
	Version        int64     `json:"version"`
}

var ErrClientNotFound = errors.New("client not found")
var ErrClientContactPersonNotFound = errors.New("client contact person not found")

var ErrClientUUIDInvalid = errors.New("client uuid is invalid")
var ErrClientContactPersonUUIDInvalid = errors.New("client contact person uuid is invalid")

type ClientRepository interface {
	Create(ctx context.Context, name string) (*Client, error)
	Get(ctx context.Context, uuid string, useCache bool) (*Client, error)
	GetAll(ctx context.Context, useCache bool) ([]Client, error)
	Update(ctx context.Context, uuid string, name string) (*Client, error)
	Delete(ctx context.Context, uuid string) (*Client, error)

	AddContactPerson(ctx context.Context, clientUUID string, name string, email string, phone []string, comment string) (*ContactPerson, error)
	UpdateContactPerson(ctx context.Context, clientUUID string, contactPersonUUID string, name string, email string, phone []string, notRelevant bool, comment string) (*ContactPerson, error)
	DeleteContactPerson(ctx context.Context, contactPersonUUID string) (*ContactPerson, error)
	GetContactPersonsForClient(ctx context.Context, clientUUID string, useCache bool) ([]ContactPerson, error)
}

func (c *ContactPerson) ToGRPC() *client.ContactPerson {
	return &client.ContactPerson{
		Uuid:        c.UUID,
		Name:        c.Name,
		Email:       c.Email,
		Phone:       c.Phone,
		NotRelevant: c.NotRelevant,
		Comment:     c.Comment,
	}
}

func (c *Client) ToGRPC() *client.Client {
	var contactPersons []*client.ContactPerson = make([]*client.ContactPerson, len(c.ContactPersons))
	for i, contactPerson := range c.ContactPersons {
		contactPersons[i] = contactPerson.ToGRPC()
	}

	return &client.Client{
		Namespace:      c.Namespace,
		CreatedAt:      timestamppb.New(c.CreationTime),
		UpdatedAt:      timestamppb.New(c.LastUpdateTime),
		Version:        c.Version,
		Uuid:           c.UUID,
		Name:           c.Name,
		ContactPersons: contactPersons,
	}
}
