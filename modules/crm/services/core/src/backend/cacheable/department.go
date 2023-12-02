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
	departmentDataCacheKeyPrefix = "crm_department_data_"
	departmentDataCacheTTL       = time.Minute * 5

	allDepartmentsCacheKeyPrefix = "crm_department_all_"
	allDepartmentsCacheTTL       = time.Minute * 5
)

func MakeDepartmentDataCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("%s_%s_%s", departmentDataCacheKeyPrefix, namespace, uuid)
}

func MakeAllDepartmentsDataCacheKey(namespace string) string {
	return fmt.Sprintf("%s_%s", allDepartmentsCacheKeyPrefix, namespace)
}

type departmentRepository struct {
	wrapedRespository models.DepartmentRepository
	logger            *slog.Logger
	namespace         string
	systemStub        *system.SystemStub
}

func (r *departmentRepository) Create(ctx context.Context, name string) (*models.Department, error) {
	department, err := r.wrapedRespository.Create(ctx, name)
	if err == nil {
		r.systemStub.Cache.Remove(ctx, MakeAllDepartmentsDataCacheKey(r.namespace))
	}
	return department, err
}
func (r *departmentRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Department, error) {
	var cacheKey string
	if useCache {
		cacheKey := MakeDepartmentDataCacheKey(r.namespace, uuid)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var department models.Department
			err := json.Unmarshal(cacheBytes, &department)
			if err != nil {
				r.logger.Warn("failed to unmarshal department from cache", "error", err.Error())
			} else {
				return &department, nil
			}
		}
	}

	department, err := r.wrapedRespository.Get(ctx, uuid, useCache)

	if err == nil && useCache {
		cacheBytes, err := json.Marshal(department)
		if err != nil {
			r.logger.Warn("failed to marshal department for cache", "error", err.Error())
			return department, err
		}

		err = r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, departmentDataCacheTTL)
		if err != nil {
			r.logger.Warn("failed to set department to cache", "error", err.Error())
		}
	}

	return department, err
}
func (r *departmentRepository) GetAll(ctx context.Context, useCache bool) ([]models.Department, error) {
	var cacheKey string
	if useCache {
		cacheKey := MakeAllDepartmentsDataCacheKey(r.namespace)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var departments []models.Department
			err := json.Unmarshal(cacheBytes, &departments)
			if err != nil {
				r.logger.Warn("failed to unmarshal departments from cache", "error", err.Error())
			} else {
				return departments, nil
			}
		}
	}

	departments, err := r.wrapedRespository.GetAll(ctx, useCache)

	if err == nil && useCache {
		cacheBytes, err := json.Marshal(departments)
		if err != nil {
			r.logger.Warn("failed to marshal departments for cache", "error", err.Error())
			return departments, err
		}

		err = r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, allDepartmentsCacheTTL)
		if err != nil {
			r.logger.Warn("failed to set departments to cache", "error", err.Error())
		}
	}

	return departments, err
}
func (r *departmentRepository) Update(ctx context.Context, uuid string, name string) (*models.Department, error) {
	department, err := r.wrapedRespository.Update(ctx, uuid, name)
	if err == nil {
		r.systemStub.Cache.Remove(ctx, MakeAllDepartmentsDataCacheKey(r.namespace), MakeDepartmentDataCacheKey(r.namespace, uuid))
	}
	return department, err
}
func (r *departmentRepository) Delete(ctx context.Context, uuid string) (*models.Department, error) {
	department, err := r.wrapedRespository.Delete(ctx, uuid)
	if err == nil {
		r.systemStub.Cache.Remove(ctx, MakeAllDepartmentsDataCacheKey(r.namespace), MakeDepartmentDataCacheKey(r.namespace, uuid))
	}
	return department, err
}
