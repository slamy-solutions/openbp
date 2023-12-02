package services

import (
	"context"
	"errors"
	"log/slog"

	department "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/department"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DepartmentService struct {
	department.UnimplementedDepartmentServiceServer

	backend backend.BackendFactory
	logger  *slog.Logger
}

func NewDepartmentServer(backend backend.BackendFactory, logger *slog.Logger) *DepartmentService {
	return &DepartmentService{
		backend: backend,
		logger:  logger,
	}
}

func (s *DepartmentService) Create(ctx context.Context, in *department.CreateDepartmentRequest) (*department.CreateDepartmentResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	d, err := bkd.DepartmentRepository().Create(ctx, in.Name)
	if err != nil {
		err := errors.Join(errors.New("failed to create department"), err)
		s.logger.With(slog.String("route", "Create")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create department: %s", err.Error())
	}

	return &department.CreateDepartmentResponse{
		Department: d.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *DepartmentService) Get(ctx context.Context, in *department.GetDepartmentRequest) (*department.GetDepartmentResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	d, err := bkd.DepartmentRepository().Get(ctx, in.Uuid, in.UseCache)
	if err != nil {
		if errors.Is(err, models.ErrDepartmentNotFound) {
			return nil, status.Errorf(codes.NotFound, "department not found")
		}

		err := errors.Join(errors.New("failed to get department"), err)
		s.logger.With(slog.String("route", "Get")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get department: %s", err.Error())
	}

	return &department.GetDepartmentResponse{
		Department: d.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *DepartmentService) GetAll(ctx context.Context, in *department.GetAllDepartmentsRequest) (*department.GetAllDepartmentsResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	departments, err := bkd.DepartmentRepository().GetAll(ctx, in.UseCache)
	if err != nil {
		err := errors.Join(errors.New("failed to get departments"), err)
		s.logger.With(slog.String("route", "GetAll")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get departments: %s", err.Error())
	}

	var responseDepartments []*department.Department = make([]*department.Department, len(departments))
	for i, department := range departments {
		responseDepartments[i] = department.ToGRPC()
	}

	return &department.GetAllDepartmentsResponse{
		Departments: responseDepartments,
	}, status.Error(codes.OK, "")
}
func (s *DepartmentService) Update(ctx context.Context, in *department.UpdateDepartmentRequest) (*department.UpdateDepartmentResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	d, err := bkd.DepartmentRepository().Update(ctx, in.Uuid, in.Name)
	if err != nil {
		if errors.Is(err, models.ErrDepartmentNotFound) {
			return nil, status.Errorf(codes.NotFound, "department not found")
		}

		err := errors.Join(errors.New("failed to update department"), err)
		s.logger.With(slog.String("route", "Update")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to update department: %s", err.Error())
	}

	return &department.UpdateDepartmentResponse{
		Department: d.ToGRPC(),
	}, status.Error(codes.OK, "")
}
func (s *DepartmentService) Delete(ctx context.Context, in *department.DeleteDepartmentRequest) (*department.DeleteDepartmentResponse, error) {
	bkd, err := getBackend(ctx, s.backend, in.Namespace, s.logger.With(slog.String("route", "GetAll")))
	if err != nil {
		return nil, err
	}

	d, err := bkd.DepartmentRepository().Delete(ctx, in.Uuid)
	if err != nil {
		if errors.Is(err, models.ErrDepartmentNotFound) {
			return nil, status.Errorf(codes.NotFound, "department not found")
		}

		err := errors.Join(errors.New("failed to delete department"), err)
		s.logger.With(slog.String("route", "Delete")).Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to delete department: %s", err.Error())
	}

	return &department.DeleteDepartmentResponse{
		Department: d.ToGRPC(),
	}, status.Error(codes.OK, "")
}
