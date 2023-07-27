package bootstrap

import (
	"context"
	"errors"
	"os"
	"strconv"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
)

const (
	_ENV_BLOCK_ROOT_USER_INIT       = "BLOCK_ROOT_USER_INIT"
	_ENV_DEFAULT_ROOT_USER_LOGIN    = "DEFAULT_ROOT_USER_LOGIN"
	_ENV_DEFAULT_ROOT_USER_PASSWORD = "DEFAULT_ROOT_USER_PASSWORD"
)

var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")

func getenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, ErrEnvVarEmpty
	}
	return v, nil
}

func getenvBool(key string) bool {
	s, err := getenvStr(key)
	if err != nil {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}

func isVaultSealed(ctx context.Context, systemStub *system.SystemStub) (bool, error) {
	// Check if vault service is sealed
	vaultStatusResponse, err := systemStub.Vault.GetStatus(ctx, &vault.GetStatusRequest{})
	if err != nil {
		return false, err
	}
	return vaultStatusResponse.Sealed, nil
}

func isRootUserCreationBlocked() bool {
	return getenvBool(_ENV_BLOCK_ROOT_USER_INIT)
}

func isRootUserCreated(ctx context.Context, nativeStub *native.NativeStub) (bool, error) {
	countResponse, err := nativeStub.Services.IAM.Actor.User.Count(ctx, &user.CountRequest{Namespace: ""})
	if err != nil {
		return false, err
	}

	return countResponse.Count != 0, nil
}
