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
	performerDataCacheKeyPrefix = "crm_performer_data_"
	performerDataCacheTTL       = time.Second * 30

	allPerformersCacheKeyPrefix = "crm_performer_all_"
	allPerformersCacheTTL       = time.Second * 30
)

func MakePerformerDataCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("%s_%s_%s", performerDataCacheKeyPrefix, namespace, uuid)
}

func MakeAllPerformersDataCacheKey(namespace string) string {
	return fmt.Sprintf("%s_%s", allPerformersCacheKeyPrefix, namespace)
}

type performerRepository struct {
	wrapedRespository models.PerformerRepository
	logger            *slog.Logger
	namespace         string
	systemStub        *system.SystemStub
}

func (r *performerRepository) Create(ctx context.Context, departmentUUID string, userUUID string) (*models.Performer, error) {
	performer, err := r.wrapedRespository.Create(ctx, departmentUUID, userUUID)
	if err == nil {
		r.systemStub.Cache.Remove(ctx, MakeAllPerformersDataCacheKey(r.namespace))
	}
	return performer, err
}
func (r *performerRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Performer, error) {
	var cacheKey string
	if useCache {
		cacheKey := MakePerformerDataCacheKey(r.namespace, uuid)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var performer models.Performer
			err := json.Unmarshal(cacheBytes, &performer)
			if err != nil {
				r.logger.Warn("failed to unmarshal performer from cache", "error", err.Error())
			} else {
				return &performer, nil
			}
		}
	}

	performer, err := r.wrapedRespository.Get(ctx, uuid, useCache)

	if err == nil && useCache {
		cacheBytes, err := json.Marshal(performer)
		if err != nil {
			r.logger.Warn("failed to marshal performer for cache", "error", err.Error())
			return performer, err
		}

		err = r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, performerDataCacheTTL)
		if err != nil {
			r.logger.Warn("failed to set performer to cache", "error", err.Error())
		}
	}

	return performer, err
}
func (r *performerRepository) GetAll(ctx context.Context, useCache bool) ([]models.Performer, error) {
	var cacheKey string
	if useCache {
		cacheKey := MakeAllPerformersDataCacheKey(r.namespace)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var performers []models.Performer
			err := json.Unmarshal(cacheBytes, &performers)
			if err != nil {
				r.logger.Warn("failed to unmarshal performers from cache", "error", err.Error())
			} else {
				return performers, nil
			}
		}
	}

	performers, err := r.wrapedRespository.GetAll(ctx, useCache)

	if err == nil && useCache {
		cacheBytes, err := json.Marshal(performers)
		if err != nil {
			r.logger.Warn("failed to marshal performers for cache", "error", err.Error())
			return performers, err
		}

		err = r.systemStub.Cache.Set(ctx, cacheKey, cacheBytes, allPerformersCacheTTL)
		if err != nil {
			r.logger.Warn("failed to set performers to cache", "error", err.Error())
		}
	}

	return performers, err
}
func (r *performerRepository) Update(ctx context.Context, uuid string, departmentUUID string) (*models.Performer, error) {
	performer, err := r.wrapedRespository.Update(ctx, uuid, departmentUUID)
	if err == nil {
		err = r.systemStub.Cache.Remove(ctx, MakeAllPerformersDataCacheKey(r.namespace), MakePerformerDataCacheKey(r.namespace, uuid))
		if err != nil {
			r.logger.Warn("failed to remove performer from cache", "error", err.Error())
		}
	}
	return performer, err
}
func (r *performerRepository) Delete(ctx context.Context, uuid string) (*models.Performer, error) {
	performer, err := r.wrapedRespository.Delete(ctx, uuid)
	if err == nil {
		err = r.systemStub.Cache.Remove(ctx, MakeAllPerformersDataCacheKey(r.namespace), MakePerformerDataCacheKey(r.namespace, uuid))
		if err != nil {
			r.logger.Warn("failed to remove performer from cache", "error", err.Error())
		}
	}
	return performer, err
}
