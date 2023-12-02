package fs

import (
	"context"
	"errors"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func HandleNamespaceCreationEvent(ctx context.Context, namespace *namespaceGRPC.Namespace, systemStub *system.SystemStub) error {
	err := prepareCollections(ctx, systemStub, namespace.Name)
	if err != nil {
		return errors.Join(errors.New("error while preparing files collection"), err)
	}

	return nil
}
