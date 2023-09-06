package catalog

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/repositories"
	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog/servers"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"google.golang.org/grpc"
)

type CatalogModule struct {
	Repository *repositories.CompositeCatalogRepository
}

func RegisterCatalogModule(logger *slog.Logger, server *grpc.Server, systemStub *system.SystemStub) *CatalogModule {
	repository := repositories.NewCompositeCatalogRepository(logger, systemStub)
	servers.RegisterGRPCServers(server, logger, repository)

	return &CatalogModule{
		Repository: repository,
	}
}
