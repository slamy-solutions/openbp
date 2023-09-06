package servers

import (
	"log/slog"

	catalogGRPC "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/repositories"
	"google.golang.org/grpc"
)

func RegisterGRPCServers(server *grpc.Server, logger *slog.Logger, repository *repositories.CompositeCatalogRepository) {
	catalogServer := NewCatalogServer(logger.With("server", "catalog"), repository)
	catalogGRPC.RegisterCatalogServiceServer(server, catalogServer)
	catalogEntryServer := NewCatalogEntryServer(logger.With("server", "entry"), repository)
	catalogGRPC.RegisterCatalogEntryServiceServer(server, catalogEntryServer)
	catalogIdexServer := NewCatalogIndexServer(logger.With("server", "index"), repository)
	catalogGRPC.RegisterCatalogIndexServiceServer(server, catalogIdexServer)
}
