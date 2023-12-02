package onec

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
)

type department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type departmentRepository struct {
	logger    *slog.Logger
	connector *connector.OneCConnector
	namespace string
}

func (r *departmentRepository) Create(ctx context.Context, name string) (*models.Department, error) {
	departmentCreateRequest := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}

	response, statusCode, err := connector.MakeRequest[department](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/department/%s", r.connector.ServerURL, r.connector.ServerToken),
		departmentCreateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create department"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to create department. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	r.logger.Info("Created department", slog.Group("department", "name", name, "uuid", response.ID))
	return &models.Department{
		Namespace: r.namespace,
		UUID:      response.ID,
		Name:      response.Name,
	}, nil
}
func (r *departmentRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Department, error) {
	//TODO: implement separate request

	departments, err := r.GetAll(ctx, useCache)
	if err != nil {
		return nil, err
	}

	for _, department := range departments {
		if department.UUID == uuid {
			return &department, nil
		}
	}

	return nil, models.ErrDepartmentNotFound
}
func (r *departmentRepository) GetAll(ctx context.Context, useCache bool) ([]models.Department, error) {
	response, statusCode, err := connector.MakeRequest[[]department](
		ctx,
		r.connector,
		"GET",
		fmt.Sprintf("%s/department/%s", r.connector.ServerURL, r.connector.ServerToken),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get all departments"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to create department. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	departments := make([]models.Department, len(*response))
	for i, department := range *response {
		departments[i] = models.Department{
			Namespace: r.namespace,
			UUID:      department.ID,
			Name:      department.Name,
		}
	}

	return departments, nil
}
func (r *departmentRepository) Update(ctx context.Context, uuid string, name string) (*models.Department, error) {
	departmentUpdateRequest := department{
		ID:   uuid,
		Name: name,
	}

	_, statusCode, err := connector.MakeRequest[department](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/department/%s", r.connector.ServerURL, r.connector.ServerToken),
		departmentUpdateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update department"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to update department. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	r.logger.Info("Updated department", slog.Group("department", "name", name, "uuid", uuid))
	return &models.Department{
		Namespace: r.namespace,
		UUID:      uuid,
		Name:      name,
	}, nil
}
func (r *departmentRepository) Delete(ctx context.Context, uuid string) (*models.Department, error) {
	existingDepartment, err := r.Get(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	departmentDeleteRequest := struct {
		ID string `json:"id"`
	}{
		ID: uuid,
	}

	_, statusCode, err := connector.MakeRequest[department](
		ctx,
		r.connector,
		"DELETE",
		fmt.Sprintf("%s/department/%s", r.connector.ServerURL, r.connector.ServerToken),
		departmentDeleteRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete department"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to delete department. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	r.logger.Info("Deleted department", slog.Group("department", "name", existingDepartment.Name, "uuid", uuid))
	return &models.Department{
		Namespace: r.namespace,
		UUID:      uuid,
		Name:      existingDepartment.Name,
	}, nil
}
