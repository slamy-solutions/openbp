package onec

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
)

type project struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	ClientId   string `json:"clientId"`
	ContactId  string `json:"contactId"`
	Department struct {
		ID   string `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"department"`
	NotRelevant bool `json:"notRelevant"`
}

type projectRepository struct {
	logger    *slog.Logger
	connector *connector.OneCConnector
	namespace string
}

func (r *projectRepository) Create(ctx context.Context, name string, clientUUID string, contactUUID string, departmentUUID string) (*models.Project, error) {
	projectCreateRequest := project{
		Name:      name,
		ClientId:  clientUUID,
		ContactId: contactUUID,
		Department: struct {
			ID   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{
			ID:   departmentUUID,
			Name: "",
		},
		NotRelevant: false,
	}

	response, statusCode, err := connector.MakeRequest[project](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/project/%s", r.connector.ServerURL, r.connector.ServerToken),
		projectCreateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to create project"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to create project. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	r.logger.Info("Created project", slog.Group("project", "name", name, "uuid", response.ID))
	return &models.Project{
		Namespace:      r.namespace,
		UUID:           response.ID,
		Name:           response.Name,
		ClientUUID:     response.ClientId,
		ContactUUID:    response.ContactId,
		DepartmentUUID: response.Department.ID,
		NotRelevant:    response.NotRelevant,
	}, nil
}

func (r *projectRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Project, error) {
	return nil, errors.New("Not implemented")
}
func (r *projectRepository) GetAll(ctx context.Context, useCache bool, clientUUID string, departmentUUID string) ([]models.Project, error) {
	response, statusCode, err := connector.MakeRequest[[]project](
		ctx,
		r.connector,
		"GET",
		fmt.Sprintf("%s/project/%s?clientId=%s&departmentId=%s", r.connector.ServerURL, r.connector.ServerToken, url.QueryEscape(clientUUID), url.QueryEscape(departmentUUID)),
		nil,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to get all projects"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to create department. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	projects := make([]models.Project, len(*response))
	for i, p := range *response {
		projects[i] = models.Project{
			Namespace:      r.namespace,
			UUID:           p.ID,
			Name:           p.Name,
			ClientUUID:     p.ClientId,
			ContactUUID:    p.ContactId,
			DepartmentUUID: p.Department.ID,
			NotRelevant:    p.NotRelevant,
		}
	}

	return projects, nil
}
func (r *projectRepository) Update(ctx context.Context, uuid string, name string, clientUUID string, contactUUID string, departmentUUID string, notRelevant bool) (*models.Project, error) {
	projectCreateRequest := project{
		ID:        uuid,
		Name:      name,
		ClientId:  clientUUID,
		ContactId: contactUUID,
		Department: struct {
			ID   string `json:"id"`
			Name string `json:"name,omitempty"`
		}{
			ID:   departmentUUID,
			Name: "",
		},
		NotRelevant: false,
	}

	response, statusCode, err := connector.MakeRequest[project](
		ctx,
		r.connector,
		"POST",
		fmt.Sprintf("%s/project/%s", r.connector.ServerURL, r.connector.ServerToken),
		projectCreateRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to update project"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode == http.StatusNotFound {
		return nil, models.ErrProjectNotFound
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to update project. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	r.logger.Info("Updated project", slog.Group("project", "name", name, "uuid", response.ID))
	return &models.Project{
		Namespace:      r.namespace,
		UUID:           response.ID,
		Name:           response.Name,
		ClientUUID:     response.ClientId,
		ContactUUID:    response.ContactId,
		DepartmentUUID: response.Department.ID,
		NotRelevant:    response.NotRelevant,
	}, nil
}
func (r *projectRepository) Delete(ctx context.Context, uuid string) (*models.Project, error) {
	existingProject, err := r.Get(ctx, uuid, false)
	if err != nil {
		return nil, err
	}

	projectDeleteRequest := struct {
		ID string `json:"id"`
	}{
		ID: uuid,
	}

	_, statusCode, err := connector.MakeRequest[struct{}](
		ctx,
		r.connector,
		"DELETE",
		fmt.Sprintf("%s/project/%s", r.connector.ServerURL, r.connector.ServerToken),
		projectDeleteRequest,
	)
	if err != nil {
		err := errors.Join(errors.New("failed to delete project"), err)
		if !errors.Is(err, models.ErrBackendUnavailable) && !errors.Is(err, models.ErrBackendMissBehaviour) {
			r.logger.Error(err.Error())
		}
		return nil, err
	}

	if statusCode != http.StatusOK {
		err := errors.Join(models.ErrBackendMissBehaviour, errors.New(fmt.Sprintf("failed to delete project. Invalid status code from the backend: %d", statusCode)), err)
		return nil, err
	}

	existingProject.Namespace = r.namespace
	r.logger.Info("Deleted project", slog.Group("project", "name", existingProject.Name, "uuid", uuid))
	return existingProject, nil
}
