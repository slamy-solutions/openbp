package services

import (
	"context"
	"errors"
	"log/slog"

	project "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/project"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer

	backend backend.BackendFactory
	logger  *slog.Logger
}

func NewProjectServer(backend backend.BackendFactory, logger *slog.Logger) *ProjectService {
	return &ProjectService{
		backend: backend,
		logger:  logger,
	}
}

func (s *ProjectService) Create(ctx context.Context, in *project.CreateProjectRequest) (*project.CreateProjectResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Create")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.ProjectRepository().Create(ctx, in.Name, in.ClientUUID, in.ContactUUID, in.DepartmentUUID)
	if err != nil {
		err := errors.Join(errors.New("failed to create project"), err)
		s.logger.With(slog.String("route", "Create")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create project: %s", err.Error())
	}

	return &project.CreateProjectResponse{
		Project: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ProjectService) Get(ctx context.Context, in *project.GetProjectRequest) (*project.GetProjectResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Get")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.ProjectRepository().Get(ctx, in.Uuid, in.UseCache)
	if err != nil {
		if errors.Is(err, models.ErrProjectNotFound) {
			return nil, status.Errorf(codes.NotFound, "project not found")
		}

		err := errors.Join(errors.New("failed to get project"), err)
		s.logger.With(slog.String("route", "Get")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get project: %s", err.Error())
	}

	return &project.GetProjectResponse{
		Project: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ProjectService) GetAll(ctx context.Context, in *project.GetAllProjectsRequest) (*project.GetAllProjectsResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	ps, err := bkd.ProjectRepository().GetAll(ctx, in.UseCache, in.ClientUUID, in.DepartmentUUID)
	if err != nil {
		err := errors.Join(errors.New("failed to get all projects"), err)
		s.logger.With(slog.String("route", "GetAll")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get all projects: %s", err.Error())
	}

	projects := make([]*project.Project, 0, len(ps))
	for _, p := range ps {
		projects = append(projects, p.ToGRPC())
	}

	return &project.GetAllProjectsResponse{
		Projects: projects,
	}, status.Error(codes.OK, "")
}
func (s *ProjectService) Update(ctx context.Context, in *project.UpdateProjectRequest) (*project.UpdateProjectResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Update")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.ProjectRepository().Update(ctx, in.Uuid, in.Name, in.ClientUUID, in.ContactUUID, in.DepartmentUUID, in.NotRelevant)
	if err != nil {
		if err == models.ErrProjectNotFound {
			return nil, status.Errorf(codes.NotFound, "project not found")
		}

		err := errors.Join(errors.New("failed to update project"), err)
		s.logger.With(slog.String("route", "Update")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update project: %s", err.Error())
	}

	return &project.UpdateProjectResponse{
		Project: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *ProjectService) Delete(ctx context.Context, in *project.DeleteProjectRequest) (*project.DeleteProjectResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "Delete")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.ProjectRepository().Delete(ctx, in.Uuid)
	if err != nil {
		if err == models.ErrProjectNotFound {
			return nil, status.Errorf(codes.NotFound, "project not found")
		}

		err := errors.Join(errors.New("failed to delete project"), err)
		s.logger.With(slog.String("route", "Delete")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete project: %s", err.Error())
	}

	return &project.DeleteProjectResponse{
		Project: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
