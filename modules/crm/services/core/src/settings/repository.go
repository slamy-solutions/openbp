package settings

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/models"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const settingsKey = "crm_settings"

var ErrVaultIsSealed = errors.New("vault is sealed")

type settingsRepository struct {
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

type SettingsRepository interface {
	Get(ctx context.Context, namespace string, useCache bool) (*Settings, error)
	Set(ctx context.Context, settings *Settings) error
}

func NewSettingsRepository(systemStub *system.SystemStub, nativeStub *native.NativeStub) SettingsRepository {
	return &settingsRepository{
		systemStub: systemStub,
		nativeStub: nativeStub,
	}
}

func (r *settingsRepository) Get(ctx context.Context, namespace string, useCache bool) (*Settings, error) {

	//TODO: additional internal cache to ommit unnecessary requests to vault

	getKeyResponse, err := r.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: namespace,
		Key:       settingsKey,
		UseCache:  useCache,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			if status.Code() == codes.NotFound {
				return &Settings{
					Namespace:   namespace,
					BackendType: models.BackendTypeNative,
					OneCData:    nil,
				}, nil
			}

			if status.Code() == codes.FailedPrecondition {
				return nil, ErrVaultIsSealed
			}
		}

		err := errors.Join(errors.New("failed to get settings from vault"), err)
		return nil, err
	}

	settingsBytes := getKeyResponse.Value
	decryptResponse, err := r.systemStub.Vault.Decrypt(ctx, &vault.DecryptRequest{
		EncryptedData: settingsBytes,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			if status.Code() == codes.FailedPrecondition {
				return nil, ErrVaultIsSealed
			}
		}

		return nil, errors.Join(errors.New("failed to decrypt settings"), err)
	}

	var settings Settings
	err = json.Unmarshal(decryptResponse.PlainData, &settings)
	if err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal settings"), err)
	}

	return &settings, nil
}

func (r *settingsRepository) Set(ctx context.Context, settings *Settings) error {
	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return errors.Join(errors.New("failed to marshal settings"), err)
	}

	encryptResponse, err := r.systemStub.Vault.Encrypt(ctx, &vault.EncryptRequest{
		PlainData: settingsBytes,
	})
	if err != nil {
		if status, ok := status.FromError(err); ok {
			if status.Code() == codes.FailedPrecondition {
				return ErrVaultIsSealed
			}
		}

		return errors.Join(errors.New("failed to encrypt settings"), err)
	}

	_, err = r.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: settings.Namespace,
		Key:       settingsKey,
		Value:     encryptResponse.EncryptedData,
	})
	if err != nil {
		return errors.Join(errors.New("failed to set settings in vault"), err)
	}

	return nil
}
