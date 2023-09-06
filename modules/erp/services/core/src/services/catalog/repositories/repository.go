package repositories

import (
	"log/slog"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type CompositeCatalogRepository struct {
	Catalog CatalogRepository
	Entry   CatalogEntryRepository
	Inxes   CatalogIndexRepository
}

func NewCompositeCatalogRepository(logger *slog.Logger, systemStub *system.SystemStub) *CompositeCatalogRepository {
	catalog := NewCatalogRepository(logger.With("repository", "catalog"), systemStub)
	entry := NewCatalogEntryRepository(logger.With("repository", "entry"), systemStub)
	indexes := NewCatalogIndexRepository(logger.With("repository", "index"), systemStub)

	return &CompositeCatalogRepository{
		Catalog: catalog,
		Entry:   entry,
		Inxes:   indexes,
	}
}
