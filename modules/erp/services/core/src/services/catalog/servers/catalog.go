package servers

import (
	"context"
	"errors"
	"log/slog"

	grpc "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/models"
	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/repositories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogServer struct {
	grpc.UnimplementedCatalogServiceServer

	logger     *slog.Logger
	repository *repositories.CompositeCatalogRepository
}

func NewCatalogServer(logger *slog.Logger, repository *repositories.CompositeCatalogRepository) *CatalogServer {
	return &CatalogServer{
		logger:     logger,
		repository: repository,
	}
}

func (s *CatalogServer) Create(ctx context.Context, in *grpc.CreateCatalogRequest) (*grpc.CreateCatalogResponse, error) {
	logger := s.logger.With("method", "Create")

	var fields []interface{}
	for _, field := range in.Fields {
		field, err := models.CatalogFieldFromGRPCFieldSchema(field)
		if err != nil {
			err = errors.Join(errors.New("failed to convert field schema"), err)
			logger.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		fields = append(fields, field)
	}

	catalog, err := s.repository.Catalog.Create(ctx, models.Catalog{
		Namespace:  in.Namespace,
		Name:       in.Name,
		PublicName: in.PublicName,
		Fields:     fields,
	})
	if err != nil {
		err = errors.Join(errors.New("failed to create catalog"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Catalog created", catalog.ToLog("catalog"))
	return &grpc.CreateCatalogResponse{
		Catalog: catalog.ToGRPCatalog(),
	}, status.Error(codes.OK, "")
}
func (s *CatalogServer) Delete(ctx context.Context, in *grpc.DeleteCatalogRequest) (*grpc.DeleteCatalogResponse, error) {
	logger := s.logger.With("method", "Delete")

	catalog, err := s.repository.Catalog.Delete(ctx, in.Namespace, in.Name)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		err = errors.Join(errors.New("failed to delete catalog"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Catalog deleted", catalog.ToLog("catalog"))
	return &grpc.DeleteCatalogResponse{}, status.Error(codes.OK, "")
}
func (s *CatalogServer) Update(ctx context.Context, in *grpc.UpdateCatalogRequest) (*grpc.UpdateCatalogReponse, error) {
	logger := s.logger.With("method", "Update")

	var fields []interface{}
	for _, field := range in.Fields {
		field, err := models.CatalogFieldFromGRPCFieldSchema(field)
		if err != nil {
			err = errors.Join(errors.New("failed to convert field schema"), err)
			logger.Error(err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		fields = append(fields, field)
	}

	catalog, err := s.repository.Catalog.Update(ctx, models.Catalog{
		Namespace:  in.Namespace,
		Name:       in.Name,
		PublicName: in.PublicName,
		Fields:     fields,
	})
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		err = errors.Join(errors.New("failed to update catalog"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Catalog updated", catalog.ToLog("catalog"))
	return &grpc.UpdateCatalogReponse{
		Catalog: catalog.ToGRPCatalog(),
	}, status.Error(codes.OK, "")
}
func (s *CatalogServer) Get(ctx context.Context, in *grpc.GetCatalogRequest) (*grpc.GetCatalogResponse, error) {
	logger := s.logger.With("method", "Get")

	catalog, err := s.repository.Catalog.Get(ctx, in.Namespace, in.Name)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		err = errors.Join(errors.New("failed to get catalog"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.GetCatalogResponse{
		Catalog: catalog.ToGRPCatalog(),
	}, status.Error(codes.OK, "")
}
func (s *CatalogServer) GetAll(ctx context.Context, in *grpc.GetAllCatalogsRequest) (*grpc.GetAllCatalogsResponse, error) {
	logger := s.logger.With("method", "GetAll")

	catalogs, err := s.repository.Catalog.GetAllInNamespace(ctx, in.Namespace)
	if err != nil {
		err = errors.Join(errors.New("failed to get catalogs"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	var grpcCatalogs []*grpc.Catalog
	for _, catalog := range catalogs {
		grpcCatalogs = append(grpcCatalogs, catalog.ToGRPCatalog())
	}

	return &grpc.GetAllCatalogsResponse{
		Catalogs: grpcCatalogs,
	}, status.Error(codes.OK, "")
}
