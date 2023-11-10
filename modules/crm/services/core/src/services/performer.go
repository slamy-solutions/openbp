package services

import (
	"context"
	"errors"
	"log/slog"

	performer "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/performer"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PerformerService struct {
	performer.UnimplementedPerformerServiceServer

	backend backend.BackendFactory
	logger  *slog.Logger
}

func NewPerformerServer(backend backend.BackendFactory, logger *slog.Logger) *PerformerService {
	return &PerformerService{
		backend: backend,
		logger:  logger,
	}
}

func (s *PerformerService) Create(ctx context.Context, in *performer.CreatePerformerRequest) (*performer.CreatePerformerResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.PerformerRepository().Create(ctx, in.DepartmentUUID, in.UserUUID)
	if err != nil {
		if err == models.ErrPerformerAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, "performer already exists")
		}

		if err == models.ErrPerformerUserNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		err := errors.Join(errors.New("failed to create performer"), err)
		s.logger.With(slog.String("route", "Create")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create performer: %s", err.Error())
	}

	return &performer.CreatePerformerResponse{
		Performer: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *PerformerService) Get(ctx context.Context, in *performer.GetPerformerRequest) (*performer.GetPerformerResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.PerformerRepository().Get(ctx, in.UUID, in.UseCache)
	if err != nil {
		if err == models.ErrPerformerNotFound {
			return nil, status.Errorf(codes.NotFound, "performer not found")
		}

		err := errors.Join(errors.New("failed to get performer"), err)
		s.logger.With(slog.String("route", "Get")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get performer: %s", err.Error())
	}

	return &performer.GetPerformerResponse{
		Performer: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *PerformerService) Update(ctx context.Context, in *performer.UpdatePerformerRequest) (*performer.UpdatePerformerResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.PerformerRepository().Update(ctx, in.UUID, in.DepartmentUUID)
	if err != nil {
		if err == models.ErrPerformerNotFound {
			return nil, status.Errorf(codes.NotFound, "performer not found")
		}

		err := errors.Join(errors.New("failed to update performer"), err)
		s.logger.With(slog.String("route", "Update")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update performer: %s", err.Error())
	}

	return &performer.UpdatePerformerResponse{
		Performer: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *PerformerService) Delete(ctx context.Context, in *performer.DeletePerformerRequest) (*performer.DeletePerformerResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	p, err := bkd.PerformerRepository().Delete(ctx, in.UUID)
	if err != nil {
		if err == models.ErrPerformerNotFound {
			return nil, status.Errorf(codes.NotFound, "performer not found")
		}

		err := errors.Join(errors.New("failed to delete performer"), err)
		s.logger.With(slog.String("route", "Delete")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete performer: %s", err.Error())
	}

	return &performer.DeletePerformerResponse{
		Performer: p.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *PerformerService) List(ctx context.Context, in *performer.ListPerformersRequest) (*performer.ListPerformersResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	ps, err := bkd.PerformerRepository().GetAll(ctx, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to list performers"), err)
		s.logger.With(slog.String("route", "List")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to list performers: %s", err.Error())
	}

	var performers []*performer.Performer = make([]*performer.Performer, len(ps))
	for i, p := range ps {
		performers[i] = p.ToGRPC()
	}

	return &performer.ListPerformersResponse{
		Performers: performers,
	}, status.Error(codes.OK, "")
}
