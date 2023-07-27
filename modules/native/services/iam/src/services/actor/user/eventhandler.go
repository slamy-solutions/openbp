package user

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

func HandleNamespaceCreationEvent(ctx context.Context, logger *log.Entry, namespace *namespaceGRPC.Namespace, systemStub *system.SystemStub) error {
	err := ensureIndexesForNamespace(ctx, namespace.Name, systemStub)
	if err != nil {
		logger.Error("failed to create indexes: " + err.Error())
		return errors.New("failed to create indexes: " + err.Error())
	}

	logger.Info("Successfully handled namespace creation event.")
	return nil
}
