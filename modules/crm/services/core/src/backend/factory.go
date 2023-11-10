package backend

import (
	"context"
	"errors"
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/cacheable"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	nativeBackend "github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/native"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/settings"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

var ErrVaultIsSealed = errors.New("vault is sealed")

type backendFactory struct {
	logger             *slog.Logger
	systemStub         *system.SystemStub
	nativeStub         *native.NativeStub
	onecBackendFactory onec.OneCBackendFactory

	settingsRepository settings.SettingsRepository
}

type BackendFactory interface {
	BuildBackendForNamespace(ctx context.Context, namespace string) (models.Backend, error)

	GetOneCBackendFactory() onec.OneCBackendFactory
}

func NewBackendFactory(logger *slog.Logger, systemStub *system.SystemStub, nativeStub *native.NativeStub) BackendFactory {
	return &backendFactory{
		logger:             logger,
		systemStub:         systemStub,
		nativeStub:         nativeStub,
		onecBackendFactory: onec.NewOneCBackendFactory(systemStub, nativeStub),
	}
}

func (f *backendFactory) BuildBackendForNamespace(ctx context.Context, namespace string) (models.Backend, error) {
	settingsData, err := f.settingsRepository.Get(ctx, namespace, true)
	if err != nil {
		if errors.Is(err, settings.ErrVaultIsSealed) {
			err = errors.Join(ErrVaultIsSealed, err)
		}

		err = errors.Join(errors.New("failed to get settings for namespace"), err)
		return nil, err
	}

	var buildedBackend models.Backend

	if settingsData.BackendType == models.BackendType1C {
		buildedBackend = f.onecBackendFactory(namespace, settingsData.OneCData.Token, settingsData.OneCData.RemoteURL)
	} else if settingsData.BackendType == models.BackendTypeNative {
		buildedBackend = nativeBackend.NewNativeBackend(f.logger, namespace, f.systemStub, f.nativeStub)
	} else {
		panic("unknown backend type")
	}

	buildedBackend = cacheable.WrapBackendIntoCachable(
		buildedBackend,
		f.logger.With(
			slog.String("backend", "cacheable"),
			slog.String("innerBackend", string(settingsData.BackendType)),
		),
		namespace,
		f.systemStub,
	)

	return buildedBackend, nil
}

func (f *backendFactory) GetOneCBackendFactory() onec.OneCBackendFactory {
	return f.onecBackendFactory
}
