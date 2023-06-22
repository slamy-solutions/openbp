package role

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	policy_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/policy"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func HandleNamespaceCreationEvent(ctx context.Context, logger *log.Entry, namespace *namespaceGRPC.Namespace, systemStub *system.SystemStub, nativeStub *native.NativeStub, policyServer *policy_server.IAMPolicyServer) error {
	err := ensureIndexesForNamespace(ctx, namespace.Name, systemStub)
	if err != nil {
		logger.Error("failed to create indexes: " + err.Error())
		return errors.New("failed to create indexes: " + err.Error())
	}

	err = ensureBuiltInsForNamespace(ctx, namespace.Name, systemStub, nativeStub, policyServer)
	if err != nil {
		logger.Error("failed to ensure built-ins: " + err.Error())
		return errors.New("failed to ensure built-ins: " + err.Error())
	}

	logger.Info("Successfully handled namespace creation event.")
	return nil
}
