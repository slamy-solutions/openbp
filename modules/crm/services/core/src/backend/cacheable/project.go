package cacheable

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const (
	projectDataCacheKeyPrefix = "crm_project_data_"
	projectDataCacheTTL       = time.Second * 30

	allProjectsCacheKeyPrefix = "crm_project_all_"
	allProjectsCacheTTL       = time.Second * 30
)

func MakeProjectDataCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("%s_%s_%s", projectDataCacheKeyPrefix, namespace, uuid)
}

func MakeAllProjectsDataCacheKey(namespace string, clientUUID string, departmentUUID string) string {
	return fmt.Sprintf("%s_%s_%s_%s", allProjectsCacheKeyPrefix, namespace, clientUUID, departmentUUID)
}

type projectRepository struct {
	wrapedRespository models.ProjectRepository
	logger            *slog.Logger
	namespace         string
	systemStub        *system.SystemStub
}

func (r *projectRepository) Create(ctx context.Context, name string, clientUUID string, contactUUID string, departmentUUID string) (*models.Project, error) {
	project, err := r.wrapedRespository.Create(ctx, name, clientUUID, contactUUID, departmentUUID)
	if err == nil {
		r.systemStub.Cache.Remove(ctx, MakeAllProjectsDataCacheKey(r.namespace, clientUUID, departmentUUID))
	}
	return project, err
}
func (r *projectRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Project, error) {
	var cacheKey string
	if useCache {
		cacheKey := MakeProjectDataCacheKey(r.namespace, uuid)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var project models.Project
			err := json.Unmarshal(cacheBytes, &project)
			if err != nil {
				r.logger.Warn("failed to unmarshal project from cache", "error", err.Error())
			} else {
				return &project, nil
			}
		}
	}

	project, err := r.wrapedRespository.Get(ctx, uuid, useCache)

	if err == nil && useCache {
		cacheBytes, err := json.Marshal(project)
		if err != nil {
			r.logger.Warn("failed to marshal project for cache", "error", err.Error())
			return project, err
		}

		err = r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, projectDataCacheTTL)
		if err != nil {
			r.logger.Warn("failed to set project to cache", "error", err.Error())
		}
	}

	return project, err
}
func (r *projectRepository) GetAll(ctx context.Context, useCache bool, clientUUID string, departmentUUID string) ([]models.Project, error) {
	var cacheKey string
	if useCache {
		cacheKey := MakeAllProjectsDataCacheKey(r.namespace, clientUUID, departmentUUID)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var projects []models.Project
			err := json.Unmarshal(cacheBytes, &projects)
			if err != nil {
				r.logger.Warn("failed to unmarshal projects from cache", "error", err.Error())
			} else {
				return projects, nil
			}
		}
	}

	projects, err := r.wrapedRespository.GetAll(ctx, useCache, clientUUID, departmentUUID)

	if err == nil && useCache {
		cacheBytes, err := json.Marshal(projects)
		if err != nil {
			r.logger.Warn("failed to marshal projects for cache", "error", err.Error())
			return projects, err
		}

		err = r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, allProjectsCacheTTL)
		if err != nil {
			r.logger.Warn("failed to set projects to cache", "error", err.Error())
		}
	}

	return projects, err
}
func (r *projectRepository) Update(ctx context.Context, uuid string, name string, clientUUID string, contactUUID string, departmentUUID string, notRelevant bool) (*models.Project, error) {
	project, err := r.wrapedRespository.Update(ctx, uuid, name, clientUUID, contactUUID, departmentUUID, notRelevant)
	if err == nil {
		err = r.systemStub.Cache.Remove(ctx, MakeAllProjectsDataCacheKey(r.namespace, clientUUID, departmentUUID), MakeProjectDataCacheKey(r.namespace, uuid))
		if err != nil {
			r.logger.Warn("failed to remove project from cache", "error", err.Error())
		}
	}
	return project, err
}
func (r *projectRepository) Delete(ctx context.Context, uuid string) (*models.Project, error) {
	project, err := r.wrapedRespository.Delete(ctx, uuid)
	if err == nil {
		err = r.systemStub.Cache.Remove(ctx, MakeAllProjectsDataCacheKey(r.namespace, project.ClientUUID, project.DepartmentUUID), MakeProjectDataCacheKey(r.namespace, uuid))
		if err != nil {
			r.logger.Warn("failed to remove project from cache", "error", err.Error())
		}
	}
	return project, err
}
