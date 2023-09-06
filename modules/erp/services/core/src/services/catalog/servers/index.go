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

type CatalogIdexServer struct {
	grpc.UnimplementedCatalogIndexServiceServer

	logger     *slog.Logger
	repository *repositories.CompositeCatalogRepository
}

func NewCatalogIndexServer(logger *slog.Logger, repository *repositories.CompositeCatalogRepository) *CatalogIdexServer {
	return &CatalogIdexServer{
		logger:     logger,
		repository: repository,
	}
}

func (s *CatalogIdexServer) ListIndexes(ctx context.Context, in *grpc.ListCatalogIndexesRequest) (*grpc.ListCatalogIndexesResponse, error) {
	logger := s.logger.With("method", "ListIndexes")

	indexes, err := s.repository.Inxes.GetAllForCatalog(ctx, in.Namespace, in.Catalog)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, "catalog not found")
		}

		err = errors.Join(errors.New("failed to get catalog indexes"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	var grpcIndexes []*grpc.CatalogIndex
	for _, index := range indexes {
		grpcIndexes = append(grpcIndexes, index.ToGRPCatalogIndex())
	}

	return &grpc.ListCatalogIndexesResponse{
		Indexes: grpcIndexes,
	}, status.Error(codes.OK, "")
}
func (s *CatalogIdexServer) EnsureIndex(ctx context.Context, in *grpc.EnsureCatalogIndexRequest) (*grpc.EnsureCatalogIndexResponse, error) {
	logger := s.logger.With("method", "EnsureIndex")

	index := models.CatalogIndexFromGRPC(in.Index)
	err := s.repository.Inxes.Ensure(ctx, index)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, "catalog not found")
		}

		err = errors.Join(errors.New("failed to ensure catalog index"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("catalog index ensured", index.ToLog("index"))
	return &grpc.EnsureCatalogIndexResponse{}, status.Error(codes.OK, "")
}
func (s *CatalogIdexServer) RemoveIndex(ctx context.Context, in *grpc.RemoveCatalogIndexRequest) (*grpc.RemoveCatalogIndexResponse, error) {
	logger := s.logger.With("method", "RemoveIndex")

	err := s.repository.Inxes.Delete(ctx, in.Namespace, in.Catalog, in.Index)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, "catalog not found")
		}

		err = errors.Join(errors.New("failed to remove catalog index"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("catalog index removed", slog.String("index", in.Index))
	return &grpc.RemoveCatalogIndexResponse{}, status.Error(codes.OK, "")
}
