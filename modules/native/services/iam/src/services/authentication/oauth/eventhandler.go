package oauth

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func HandleNamespaceCreationEvent(ctx context.Context, logger *log.Entry, namespace *namespaceGRPC.Namespace, systemStub *system.SystemStub) error {
	err := EnsureIndexesForNamespace(ctx, namespace.Name, systemStub)
	if err != nil {
		logger.Error("failed to create indexes: " + err.Error())
		return errors.New("failed to create indexes: " + err.Error())
	}

	logger.Info("Successfully handled namespace creation event.")
	return nil
}
