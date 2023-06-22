package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestIAMAuthenticationPassword(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithIAMService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	identityCreateResponse, err := nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	plainPassword := tools.GetRandomString(20)

	_, err = nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  plainPassword,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IAM.Authentication.Password.Delete(context.Background(), &password.DeleteRequest{Namespace: "", Identity: identityCreateResponse.Identity.Uuid})

	authenticateResponse, err := nativeStub.Services.IAM.Authentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  plainPassword,
	})
	require.Nil(t, err)
	assert.True(t, authenticateResponse.Authenticated)

	passwordDeleteResponse, err := nativeStub.Services.IAM.Authentication.Password.Delete(ctx, &password.DeleteRequest{Namespace: "", Identity: identityCreateResponse.Identity.Uuid})
	require.Nil(t, err)
	assert.True(t, passwordDeleteResponse.Existed)
}
