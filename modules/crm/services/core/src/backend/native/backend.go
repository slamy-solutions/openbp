package native

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/client"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/department"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/kanban"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/performer"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native/project"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type nativeBackend struct {
	clientRepository     client.ClientRepository
	performerRespository performer.PerformerRepository
	departmentRepository department.DepartmentRepository
	projectRepository    project.ProjectRepository
	kanbanRepository     kanban.KanbanRepository
}

func NewNativeBackend(logger *slog.Logger, namespace string, systemStub *system.SystemStub, nativeStub *native.NativeStub) models.Backend {
	clientRepository := client.NewClientRepository(logger.With(slog.String("backend", "native"), slog.String("repository", "client")), namespace, systemStub)
	performerRespository := performer.NewPerformerRepository(logger.With(slog.String("backend", "native"), slog.String("repository", "performer")), namespace, systemStub, nativeStub)
	departmentRepository := department.NewDepartmentRepository(logger.With(slog.String("backend", "native"), slog.String("repository", "department")), namespace, systemStub, nativeStub)
	projectRepository := project.NewProjectRepository(logger.With(slog.String("backend", "native"), slog.String("repository", "project")), namespace, systemStub)
	kanbanRepository := kanban.NewKanbanRepository(logger.With(slog.String("backend", "native"), slog.String("repository", "kanban")), namespace, systemStub, nativeStub, &clientRepository, &performerRespository, &departmentRepository, &projectRepository)

	return &nativeBackend{
		clientRepository:     clientRepository,
		performerRespository: performerRespository,
		departmentRepository: departmentRepository,
		projectRepository:    projectRepository,
		kanbanRepository:     kanbanRepository,
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

func (b *nativeBackend) DepartmentRepository() models.DepartmentRepository {
	return &b.departmentRepository
}

func (b *nativeBackend) ProjectRepository() models.ProjectRepository {
	return &b.projectRepository
}

func (b *nativeBackend) KanbanRepository() models.KanbanRepository {
	return &b.kanbanRepository
}
