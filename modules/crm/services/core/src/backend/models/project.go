package models

import (
	"context"
	"errors"

	project "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/project"
)

type Project struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`

	ClientUUID     string `json:"clientUUID"`
	ContactUUID    string `json:"contactUUID"`
	DepartmentUUID string `json:"departmentUUID"`

	NotRelevant bool `json:"notRelevant"`
}

var ErrProjectNotFound = errors.New("project not found")
var ErrProjectUUIDInvalid = errors.New("project uuid is invalid")

var ErrProjectBadClientUUID = errors.New("bad client UUID")
var ErrProjectBadContactUUID = errors.New("bad contact UUID")
var ErrProjectBadDepartmentUUID = errors.New("bad department UUID")

type ProjectRepository interface {
	Create(ctx context.Context, name string, clientUUID string, contactUUID string, departmentUUID string) (*Project, error)
	Get(ctx context.Context, uuid string, useCache bool) (*Project, error)
	GetAll(ctx context.Context, useCache bool, clientUUID string, departmentUUID string) ([]Project, error)
	Update(ctx context.Context, uuid string, name string, clientUUID string, contactUUID string, departmentUUID string, notRelevant bool) (*Project, error)
	Delete(ctx context.Context, uuid string) (*Project, error)
}

func (p *Project) ToGRPC() *project.Project {
	return &project.Project{
		Namespace: p.Namespace,
		Uuid:      p.UUID,
		Name:      p.Name,

		ClientUUID:     p.ClientUUID,
		ContactUUID:    p.ContactUUID,
		DepartmentUUID: p.DepartmentUUID,

		NotRelevant: p.NotRelevant,
	}
}
