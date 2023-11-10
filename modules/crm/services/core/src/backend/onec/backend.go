package onec

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/connector"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/performer"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type OneCBackendFactory = func(namespace string, serverToken string, serverURL string) models.Backend

type OneCBackend struct {
	Namespace string

	connector *connector.OneCConnector

	clientRepository    *clientRepository
	performerRepository performer.PerformerRepository
}

// Factory for OneCBackends. You should use only one factor instance per entire application. Factory makes sure, that backend shares same connections pool between all the servers thus greatly increasing performance.
func NewOneCBackendFactory(systemStub *system.SystemStub, nativeStub *native.NativeStub) OneCBackendFactory {
	logger := slog.Default()

	return func(namespace string, serverToken string, serverURL string) models.Backend {
		backend := OneCBackend{
			Namespace: namespace,
			connector: connector.NewOneCConnector(serverURL, serverToken),
		}

		backend.clientRepository = &clientRepository{
			namespace: namespace,
			connector: backend.connector,
			logger:    logger.With("repository", "client"),
		}

		backend.performerRepository = performer.NewPerformerRepository(logger.With("repository", "performer"), namespace, systemStub, nativeStub, backend.connector)

		return &backend
	}
}

func (b *OneCBackend) GetType() models.BackendType {
	return models.BackendType1C
}

func (b *OneCBackend) ClientRepository() models.ClientRepository {
	return b.clientRepository
}

func (b *OneCBackend) PerformerRepository() models.PerformerRepository {
	return &b.performerRepository
}
