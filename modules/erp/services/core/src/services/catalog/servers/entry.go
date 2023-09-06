package servers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"

	grpc "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/models"
	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogEntryServer struct {
	grpc.UnimplementedCatalogEntryServiceServer

	logger     *slog.Logger
	repository *repositories.CompositeCatalogRepository
}

func NewCatalogEntryServer(logger *slog.Logger, repository *repositories.CompositeCatalogRepository) *CatalogEntryServer {
	return &CatalogEntryServer{
		logger:     logger,
		repository: repository,
	}
}

func (s *CatalogEntryServer) Create(ctx context.Context, in *grpc.CreateCatalogEntryRequest) (*grpc.CreateCatalogEntryResponse, error) {
	logger := s.logger.With("method", "Create")

	var data map[string]interface{}
	err := json.Unmarshal(in.Data, &data)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid data. cant parse as JSON")
	}

	catalogEntry := models.CatalogEntry{
		Namespace: in.Namespace,
		Catalog:   in.Catalog,
		Data:      data,
	}

	newCatalogEntry, err := s.repository.Entry.Create(ctx, catalogEntry)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return nil, status.Error(codes.NotFound, "catalog not found")
		}

		if errors.Is(err, repositories.ErrCatalogDataSchemeInvalid) {
			return nil, status.Error(codes.InvalidArgument, "invalid data scheme"+err.Error())
		}

		err = errors.Join(errors.New("failed to create catalog entry"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Catalog entry created", newCatalogEntry.ToLog("catalogEntry"))
	return &grpc.CreateCatalogEntryResponse{
		Entry: newCatalogEntry.ToGRPCCatalogEntry(),
	}, status.Error(codes.OK, "")
}
func (s *CatalogEntryServer) Delete(ctx context.Context, in *grpc.DeleteCatalogEntryRequest) (*grpc.DeleteCatalogEntryResponse, error) {
	logger := s.logger.With("method", "Delete")

	entryUUID, err := primitive.ObjectIDFromHex(in.Entry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid entry UUID")
	}

	catalogEntry, err := s.repository.Entry.Delete(ctx, in.Namespace, in.Catalog, entryUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogEntryNotFound) {
			return nil, status.Error(codes.NotFound, "catalog entry not found")
		}

		err = errors.Join(errors.New("failed to delete catalog entry"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Catalog entry deleted", catalogEntry.ToLog("catalogEntry"))
	return &grpc.DeleteCatalogEntryResponse{}, status.Error(codes.OK, "")
}
func (s *CatalogEntryServer) Update(ctx context.Context, in *grpc.UpdateCatalogEntryRequest) (*grpc.UpdateCatalogEntryResponse, error) {
	logger := s.logger.With("method", "Update")

	entryUUID, err := primitive.ObjectIDFromHex(in.Entry.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid entry UUID")
	}

	var data map[string]interface{}
	err = json.Unmarshal(in.Entry.Data, &data)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid data. cant parse as JSON")
	}

	catalogEntry := models.CatalogEntry{
		Namespace: in.Entry.Namespace,
		Catalog:   in.Entry.Catalog,
		UUID:      entryUUID,
		Data:      data,
	}

	newCatalogEntry, err := s.repository.Entry.Update(ctx, catalogEntry)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogEntryNotFound) {
			return nil, status.Error(codes.NotFound, "catalog entry not found")
		}

		if errors.Is(err, repositories.ErrCatalogDataSchemeInvalid) {
			return nil, status.Error(codes.InvalidArgument, "invalid data scheme"+err.Error())
		}

		err = errors.Join(errors.New("failed to update catalog entry"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.Info("Catalog entry updated", newCatalogEntry.ToLog("catalogEntry"))
	return &grpc.UpdateCatalogEntryResponse{
		Entry: newCatalogEntry.ToGRPCCatalogEntry(),
	}, status.Error(codes.OK, "")
}

func (s *CatalogEntryServer) Get(ctx context.Context, in *grpc.GetCatalogEntryRequest) (*grpc.GetCatalogEntryResponse, error) {
	logger := s.logger.With("method", "Get")

	entryUUID, err := primitive.ObjectIDFromHex(in.Entry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid entry UUID")
	}

	catalogEntry, err := s.repository.Entry.Get(ctx, in.Namespace, in.Catalog, entryUUID)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogEntryNotFound) {
			return nil, status.Error(codes.NotFound, "catalog entry not found")
		}

		err = errors.Join(errors.New("failed to get catalog entry"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.GetCatalogEntryResponse{
		Entry: catalogEntry.ToGRPCCatalogEntry(),
	}, status.Error(codes.OK, "")
}

func (s *CatalogEntryServer) List(in *grpc.ListCatalogEntriesRequest, out grpc.CatalogEntryService_ListServer) error {
	ctx := out.Context()
	logger := s.logger.With("method", "List")

	stream, err := s.repository.Entry.List(ctx, in.Namespace, in.Catalog, in.Limit, in.Limit)
	if err != nil {
		if errors.Is(err, repositories.ErrCatalogNotFound) {
			return status.Error(codes.NotFound, "catalog not found")
		}

		err = errors.Join(errors.New("failed to list catalog entries"), err)
		logger.Error(err.Error())
		return status.Error(codes.Internal, err.Error())
	}
	defer stream.Close(context.Background())

	for {
		catalogEntry, err := stream.Next(ctx)
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.Join(errors.New("failed to list catalog entries from stream"), err)
			logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

		err = out.Send(&grpc.ListCatalogEntriesResponse{
			Entry: catalogEntry.ToGRPCCatalogEntry(),
		})
		if err != nil {
			err = errors.Join(errors.New("failed to send catalog entry to stream"), err)
			logger.Error(err.Error())
			return status.Error(codes.Internal, err.Error())
		}

	}

	return status.Error(codes.OK, "")
}

func (s *CatalogEntryServer) Count(ctx context.Context, in *grpc.CountCatalogEntriesRequest) (*grpc.CountCatalogEntriesResponse, error) {
	logger := s.logger.With("method", "Count")

	count, err := s.repository.Entry.Count(ctx, in.Namespace, in.Catalog)
	if err != nil {
		err = errors.Join(errors.New("failed to count catalog entries"), err)
		logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.CountCatalogEntriesResponse{
		Count: count,
	}, status.Error(codes.OK, "")
}
