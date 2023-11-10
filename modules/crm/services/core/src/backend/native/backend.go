package native

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/client"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/performer"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type nativeBackend struct {
	clientRepository     client.ClientRepository
	performerRespository performer.PerformerRepository
}

func NewNativeBackend(logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub) models.Backend {
	return &nativeBackend{
		clientRepository:     client.NewClientRepository(logger.With(slog.String("backend", "native")), namespace, systemStub),
		performerRespository: performer.NewPerformerRepository(logger.With(slog.String("backend", "native")), namespace, systemStub, nativeStub),
	}
}

func (b *nativeBackend) GetType() models.BackendType {
	return models.BackendType1C
}

func (b *nativeBackend) ClientRepository() models.ClientRepository {
	return &b.clientRepository
}

func (b *nativeBackend) PerformerRepository() models.PerformerRepository {
	return &b.performerRespository
}
