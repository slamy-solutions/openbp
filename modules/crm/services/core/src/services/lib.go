package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getBackend(ctx context.Context, factory backend.BackendFactory, namespace string, logger *slog.Logger) (models.Backend, error) {
	bkd, err := factory.BuildBackendForNamespace(ctx, namespace)
	if err != nil {
		if errors.Is(err, backend.ErrVaultIsSealed) {
			return nil, status.Errorf(codes.FailedPrecondition, "vault is sealed")
		}

		err := errors.Join(errors.New("failed to build backend"), err)
		logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to build backend: %s", err.Error())
	}

	return bkd, nil
}
