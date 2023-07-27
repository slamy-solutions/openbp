package telemetry

import (
	"context"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const fleetDevicesUniqueIndexName = "unique_entries_search"
const fleetDevicesAddedIndex = "added_search"

func CreateIndexesForNamespace(ctx context.Context, systemStub *system.SystemStub, namespace string) error {

	return nil
}
