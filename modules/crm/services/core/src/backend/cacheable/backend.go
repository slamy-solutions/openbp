package cacheable

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type cacheableBackend struct {
	outerBackend         models.Backend
	clientRepository     models.ClientRepository
	performerRepository  models.PerformerRepository
	departmentRepository models.DepartmentRepository
	projectRepository    models.ProjectRepository
	kanbanRepository     models.KanbanRepository
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
		departmentRepository: &departmentRepository{
			wrapedRespository: outerBackend.DepartmentRepository(),
			logger:            logger.With(slog.String("repository", "department"), slog.String("namespace", namespace)),
			namespace:         namespace,
			systemStub:        systemStub,
		},
		projectRepository: &projectRepository{
			wrapedRespository: outerBackend.ProjectRepository(),
			logger:            logger.With(slog.String("repository", "project"), slog.String("namespace", namespace)),
			namespace:         namespace,
			systemStub:        systemStub,
		},
		kanbanRepository: &kanbanRepository{
			wrapedRespository: outerBackend.KanbanRepository(),
			logger:            logger.With(slog.String("repository", "kanban"), slog.String("namespace", namespace)),
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

func (b *cacheableBackend) DepartmentRepository() models.DepartmentRepository {
	return b.departmentRepository
}

func (b *cacheableBackend) ProjectRepository() models.ProjectRepository {
	return b.projectRepository
}
func (b *cacheableBackend) KanbanRepository() models.KanbanRepository {
	return b.kanbanRepository
}
