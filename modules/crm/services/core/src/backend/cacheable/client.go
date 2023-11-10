package cacheable

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const (
	clientDataCacheKeyPrefix = "crm_client_data_"
	clientDataCacheTTL       = time.Second * 30

	allClientsCacheKeyPrefix = "crm_client_all_"
	allClientsCacheTTL       = time.Second * 30

	allClientContactPersonsCacheKeyPrefix = "crm_client_contactperson_all_"
	allClientContactPersonsCacheTTL       = time.Second * 30
)

func MakeClientDataCacheKey(namespace string, uuid string) string {
	return fmt.Sprintf("%s_%s_%s", clientDataCacheKeyPrefix, namespace, uuid)
}

func GetClientFromCache(ctx context.Context, systemStub *system.SystemStub, namespace string, uuid string) (*models.Client, error) {
	cacheKey := MakeClientDataCacheKey(namespace, uuid)
	cacheBytes, _ := systemStub.Cache.Get(ctx, cacheKey)
	if cacheBytes != nil {
		var client models.Client
		err := json.Unmarshal(cacheBytes, &client)
		if err != nil {
			return nil, err
		}

		return &client, nil
	}

	return nil, nil
}

func PutClientInCache(ctx context.Context, systemStub *system.SystemStub, client *models.Client) error {
	cacheKey := MakeClientDataCacheKey(client.Namespace, client.UUID)
	cacheBytes, err := json.Marshal(client)
	if err != nil {
		return errors.Join(errors.New("failed to marshal client to cache"), err)
	}

	err = systemStub.Cache.Set(ctx, cacheKey, cacheBytes, clientDataCacheTTL)
	if err != nil {
		return errors.Join(errors.New("failed to set client to cache"), err)
	}
	return nil
}

func RemoveClientFromCache(ctx context.Context, systemStub *system.SystemStub, namespace string, uuid string) error {
	clientCacheKey := MakeClientDataCacheKey(namespace, uuid)
	allClientsCacheKey := MakeAllClientsCacheKey(namespace)
	allContactPersonsCacheKey := MakeAllClientContactPersonsCacheKey(namespace, uuid)
	err := systemStub.Cache.Remove(ctx, clientCacheKey, allClientsCacheKey, allContactPersonsCacheKey)
	if err != nil {
		return errors.Join(errors.New("failed to remove client from cache"), err)
	}
	return nil
}

func MakeAllClientsCacheKey(namespace string) string {
	return fmt.Sprintf("%s_%s", allClientsCacheKeyPrefix, namespace)
}

func GetAllClientsFromCache(ctx context.Context, systemStub *system.SystemStub, namespace string) ([]models.Client, error) {
	cacheKey := MakeAllClientsCacheKey(namespace)
	cacheBytes, _ := systemStub.Cache.Get(ctx, cacheKey)
	if cacheBytes != nil {
		var clients []models.Client
		err := json.Unmarshal(cacheBytes, &clients)
		if err != nil {
			return nil, err
		}

		return clients, nil
	}

	return nil, nil
}

func PutAllClientsInCache(ctx context.Context, systemStub *system.SystemStub, namespace string, clients []models.Client) error {
	cacheKey := MakeAllClientsCacheKey(namespace)
	cacheBytes, err := json.Marshal(clients)
	if err != nil {
		return errors.Join(errors.New("failed to marshal clients to cache"), err)
	}

	err = systemStub.Cache.Set(ctx, cacheKey, cacheBytes, allClientsCacheTTL)
	if err != nil {
		return errors.Join(errors.New("failed to set clients to cache"), err)
	}
	return nil
}

func RemoveAllClientsFromCache(ctx context.Context, systemStub *system.SystemStub, namespace string) error {
	cacheKey := MakeAllClientsCacheKey(namespace)
	err := systemStub.Cache.Remove(ctx, cacheKey)
	if err != nil {
		return errors.Join(errors.New("failed to remove clients from cache"), err)
	}
	return nil
}

func MakeAllClientContactPersonsCacheKey(namespace string, clientUUID string) string {
	return fmt.Sprintf("%s_%s_%s", allClientContactPersonsCacheKeyPrefix, namespace, clientUUID)
}

func GetAllClientContactPersonsFromCache(ctx context.Context, systemStub *system.SystemStub, namespace string, clientUUID string) ([]models.ContactPerson, error) {
	cacheKey := MakeAllClientContactPersonsCacheKey(namespace, clientUUID)
	cacheBytes, _ := systemStub.Cache.Get(ctx, cacheKey)
	if cacheBytes != nil {
		var contactPersons []models.ContactPerson
		err := json.Unmarshal(cacheBytes, &contactPersons)
		if err != nil {
			return nil, err
		}

		return contactPersons, nil
	}

	return nil, nil
}

func PutAllClientContactPersonsInCache(ctx context.Context, systemStub *system.SystemStub, namespace string, contactPersons []models.ContactPerson, clientUUID string) error {
	cacheKey := MakeAllClientContactPersonsCacheKey(namespace, clientUUID)
	cacheBytes, err := json.Marshal(contactPersons)
	if err != nil {
		return errors.Join(errors.New("failed to marshal contact persons to cache"), err)
	}

	err = systemStub.Cache.Set(ctx, cacheKey, cacheBytes, allClientContactPersonsCacheTTL)
	if err != nil {
		return errors.Join(errors.New("failed to set contact persons to cache"), err)
	}

	return nil
}

func RemoveAllClientContactPersonsFromCache(ctx context.Context, systemStub *system.SystemStub, namespace string, clientUUID string) error {
	cacheKey := MakeAllClientContactPersonsCacheKey(namespace, clientUUID)
	err := systemStub.Cache.Remove(ctx, cacheKey)
	if err != nil {
		return errors.Join(errors.New("failed to remove contact persons from cache"), err)
	}
	return nil
}

type clientRepository struct {
	wrapedRespository models.ClientRepository
	logger            *slog.Logger
	namespace         string
	systemStub        *system.SystemStub
}

func (r *clientRepository) Create(ctx context.Context, name string) (*models.Client, error) {
	client, err := r.wrapedRespository.Create(ctx, name)
	if err != nil {
		RemoveAllClientsFromCache(ctx, r.systemStub, r.namespace)
	}
	return client, err
}
func (r *clientRepository) Get(ctx context.Context, uuid string, useCache bool) (*models.Client, error) {
	if useCache {
		client, err := GetClientFromCache(ctx, r.systemStub, r.namespace, uuid)
		if err != nil {
			err = errors.Join(errors.New("failed to get client from cache"), err)
			r.logger.Warn(err.Error())
		} else if client != nil {
			return client, nil
		}
	}

	client, err := r.wrapedRespository.Get(ctx, uuid, useCache)
	if err != nil {
		return nil, err
	}

	if useCache {
		err = PutClientInCache(ctx, r.systemStub, client)
		if err != nil {
			err = errors.Join(errors.New("failed to put client to cache"), err)
			r.logger.Warn(err.Error())
		}
	}

	return client, nil
}
func (r *clientRepository) GetAll(ctx context.Context, useCache bool) ([]models.Client, error) {
	if useCache {
		clients, err := GetAllClientsFromCache(ctx, r.systemStub, r.namespace)
		if err != nil {
			err = errors.Join(errors.New("failed to get clients from cache"), err)
			r.logger.Warn(err.Error())
		} else if clients != nil {
			return clients, nil
		}
	}

	clients, err := r.wrapedRespository.GetAll(ctx, useCache)
	if err != nil {
		return nil, err
	}

	if useCache {
		err = PutAllClientsInCache(ctx, r.systemStub, r.namespace, clients)
		if err != nil {
			err = errors.Join(errors.New("failed to put clients to cache"), err)
			r.logger.Warn(err.Error())
		}
	}

	return clients, nil
}
func (r *clientRepository) Update(ctx context.Context, uuid string, name string) (*models.Client, error) {
	client, err := r.wrapedRespository.Update(ctx, uuid, name)
	if err != nil {
		cacheErr := RemoveClientFromCache(ctx, r.systemStub, r.namespace, uuid)
		if cacheErr != nil {
			cacheErr = errors.Join(errors.New("failed to remove client from cache"), cacheErr)
			r.logger.Warn(cacheErr.Error())
		}
	}
	return client, err
}
func (r *clientRepository) Delete(ctx context.Context, uuid string) (*models.Client, error) {
	client, err := r.wrapedRespository.Delete(ctx, uuid)
	if err != nil {
		cacheErr := RemoveClientFromCache(ctx, r.systemStub, r.namespace, uuid)
		if cacheErr != nil {
			cacheErr = errors.Join(errors.New("failed to remove client from cache"), cacheErr)
			r.logger.Warn(cacheErr.Error())
		}
	}
	return client, err
}

func (r *clientRepository) AddContactPerson(ctx context.Context, clientUUID string, name string, email string, phone []string, comment string) (*models.ContactPerson, error) {
	contactPerson, err := r.wrapedRespository.AddContactPerson(ctx, clientUUID, name, email, phone, comment)
	if err != nil {
		cacheErr := RemoveClientFromCache(ctx, r.systemStub, r.namespace, clientUUID)
		if cacheErr != nil {
			cacheErr = errors.Join(errors.New("failed to remove client from cache while adding contact person"), cacheErr)
			r.logger.Warn(cacheErr.Error())
		}
	}
	return contactPerson, err
}
func (r *clientRepository) UpdateContactPerson(ctx context.Context, clientUUID string, contactPersonUUID string, name string, email string, phone []string, notRelevant bool, comment string) (*models.ContactPerson, error) {
	contactPerson, err := r.wrapedRespository.UpdateContactPerson(ctx, clientUUID, contactPersonUUID, name, email, phone, notRelevant, comment)
	if err != nil {
		cacheErr := RemoveClientFromCache(ctx, r.systemStub, r.namespace, clientUUID)
		if cacheErr != nil {
			cacheErr = errors.Join(errors.New("failed to remove client contact person from cache"), cacheErr)
			r.logger.Warn(cacheErr.Error())
		}
	}
	return contactPerson, err
}
func (r *clientRepository) DeleteContactPerson(ctx context.Context, contactPersonUUID string) (*models.ContactPerson, error) {
	contactPerson, err := r.wrapedRespository.DeleteContactPerson(ctx, contactPersonUUID)
	if err != nil {
		cacheErr := RemoveClientFromCache(ctx, r.systemStub, r.namespace, contactPerson.ClientUUID)
		if cacheErr != nil {
			cacheErr = errors.Join(errors.New("failed to remove client contact person from cache"), cacheErr)
			r.logger.Warn(cacheErr.Error())
		}
	}
	return contactPerson, err
}
func (r *clientRepository) GetContactPersonsForClient(ctx context.Context, clientUUID string, useCache bool) ([]models.ContactPerson, error) {
	var cacheKey string
	if useCache {
		cacheKey = MakeAllClientContactPersonsCacheKey(r.namespace, clientUUID)
		cacheBytes, _ := r.systemStub.Cache.Get(ctx, cacheKey)
		if cacheBytes != nil {
			var contactPersons []models.ContactPerson
			err := json.Unmarshal(cacheBytes, &contactPersons)
			if err != nil {
				err = errors.Join(errors.New("failed to unmarshal contact persons from cache"), err)
				r.logger.Warn(err.Error())
			} else {
				return contactPersons, nil
			}
		}
	}

	contactPersons, err := r.wrapedRespository.GetContactPersonsForClient(ctx, clientUUID, useCache)
	if err != nil {
		return nil, err
	}

	if useCache {
		err = PutAllClientContactPersonsInCache(ctx, r.systemStub, r.namespace, contactPersons, clientUUID)
		if err != nil {
			err = errors.Join(errors.New("failed to put contact persons to cache"), err)
			r.logger.Warn(err.Error())
		}
	}

	return contactPersons, nil
}
