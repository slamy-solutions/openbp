package models

import (
	"context"
	"errors"

	performer "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/performer"
)

type Performer struct {
	Namespace      string `json:"namespace"`
	UUID           string `json:"uuid"`
	DepartmentUUID string `json:"departmentUUID"`

	// Associated native_iam_user UUID
	UserUUID string `json:"userUUID"`

	// Name extracted from user public name
	Name string `json:"name"`
	// AvatarURL extracted from user
	AvatarURL string `json:"avatarURL"`
}

var ErrPerformerNotFound = errors.New("perofrmer not found")
var ErrPerformerAlreadyExists = errors.New("perofrmer already exists. Use diferrent user to create performer")
var ErrPerformerUUIDInvalid = errors.New("perofrmer uuid is invalid")

var ErrPerformerBadDepartmentUUID = errors.New("bad department UUID")
var ErrPerformerBadUserUUID = errors.New("bad user UUID")
var ErrPerformerUserNotFound = errors.New("user not found")

type PerformerRepository interface {
	Create(ctx context.Context, departmentUUID string, userUUID string) (*Performer, error)
	Get(ctx context.Context, uuid string, useCache bool) (*Performer, error)
	GetAll(ctx context.Context, useCache bool) ([]Performer, error)
	Update(ctx context.Context, uuid string, departmentUUID string) (*Performer, error)
	Delete(ctx context.Context, uuid string) (*Performer, error)
}

func (p *Performer) ToGRPC() *performer.Performer {
	return &performer.Performer{
		UUID:           p.UUID,
		Namespace:      p.Namespace,
		DepartmentUUID: p.DepartmentUUID,
		UserUUID:       p.UserUUID,
		Name:           p.Name,
		AvatarUrl:      p.AvatarURL,
	}
}
