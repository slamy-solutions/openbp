package smoke

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

func TestIAMToken(t *testing.T) {
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithIAMTokenService().WithIAMIdentityService())
	err := nativeStub.Connect()
	require.Nil(t, err)
	defer nativeStub.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	metadata := tools.GetRandomString(20)

	tokenCreateResponse, err := nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  metadata,
	})
	require.Nil(t, err)
	defer nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})

	tokenValidateResponse, err := nativeStub.Services.IamToken.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(t, err)
	require.Equal(t, token.ValidateResponse_OK, tokenValidateResponse.Status)
	assert.Equal(t, metadata, tokenValidateResponse.TokenData.CreationMetadata)
	assert.Equal(t, tokenCreateResponse.TokenData.Uuid, tokenValidateResponse.TokenData.Uuid)

	tokenDeleteResponse, err := nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(t, err)
	assert.True(t, tokenDeleteResponse.Existed)
}
