package models

import (
	"context"
	"errors"

	department "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/department"
)

type Department struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
}

var ErrDepartmentNotFound = errors.New("department not found")
var ErrDepartmentUUIDInvalid = errors.New("department uuid invalid")

type DepartmentRepository interface {
	Create(ctx context.Context, name string) (*Department, error)
	Get(ctx context.Context, uuid string, useCache bool) (*Department, error)
	GetAll(ctx context.Context, useCache bool) ([]Department, error)
	Update(ctx context.Context, uuid string, name string) (*Department, error)
	Delete(ctx context.Context, uuid string) (*Department, error)
}

func (p *Department) ToGRPC() *department.Department {
	return &department.Department{
		Namespace: p.Namespace,
		Uuid:      p.UUID,
		Name:      p.Name,
	}
}
