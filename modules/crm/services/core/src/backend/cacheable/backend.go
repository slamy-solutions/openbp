package cacheable

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type cacheableBackend struct {
	outerBackend        models.Backend
	clientRepository    models.ClientRepository
	performerRepository models.PerformerRepository
}

func WrapBackendIntoCachable(outerBackend models.Backend, logger *slog.Logger, namespace string, systemStub *system.SystemStub) models.Backend {
	return &cacheableBackend{
		outerBackend: outerBackend,
		clientRepository: &clientRepository{
			wrapedRespository: outerBackend.ClientRepository(),
			logger:            logger.With(slog.String("repository", "client"), slog.String("namespace", namespace)),
			namespace:         namespace,
			systemStub:        systemStub,
		},
		performerRepository: &performerRepository{
			wrapedRespository: outerBackend.PerformerRepository(),
			logger:            logger.With(slog.String("repository", "performer"), slog.String("namespace", namespace)),
			namespace:         namespace,
			systemStub:        systemStub,
		},
	}
}

func (b *cacheableBackend) GetType() models.BackendType {
	return b.outerBackend.GetType()
}

func (b *cacheableBackend) ClientRepository() models.ClientRepository {
	return b.clientRepository
}

func (b *cacheableBackend) PerformerRepository() models.PerformerRepository {
	return b.performerRepository
}
