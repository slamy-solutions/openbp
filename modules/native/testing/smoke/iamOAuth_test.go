package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestIAMAuth(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithIAMService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create identity
	identityCreateResponse, err := nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	//Add some password to the identity
	plainPassword := tools.GetRandomString(20)

	_, err = nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  plainPassword,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IAM.Authentication.Password.Delete(context.Background(), &password.DeleteRequest{Namespace: "", Identity: identityCreateResponse.Identity.Uuid})

	// Create Auth token
	createAuthTokenResponse, err := nativeStub.Services.IAM.Auth.CreateTokenWithPassword(ctx, &auth.CreateTokenWithPasswordRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  plainPassword,
		Metadata:  "",
		Scopes:    []*auth.Scope{},
	})
	require.Nil(t, err)

	getTokenDataResponse, err := nativeStub.Services.IAM.Token.RawGet(ctx, &token.RawGetRequest{
		Token:    createAuthTokenResponse.AccessToken,
		UseCache: true,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{
		Namespace: "",
		Uuid:      getTokenDataResponse.TokenData.Uuid,
	})

	// Try to authorize with Auth token
	authorizeResponse, err := nativeStub.Services.IAM.Auth.CheckAccessWithToken(ctx, &auth.CheckAccessWithTokenRequest{
		AccessToken: createAuthTokenResponse.AccessToken,
		Scopes:      []*auth.Scope{},
	})
	require.Nil(t, err)
	assert.Equal(t, auth.CheckAccessWithTokenResponse_OK, authorizeResponse.Status)
}
