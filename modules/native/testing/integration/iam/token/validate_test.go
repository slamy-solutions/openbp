package token

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ValidateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ValidateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ValidateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestValidateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateTestSuite))
}

func (s *ValidateTestSuite) TestValidateTokenFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	validateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, validateResponse.Status)

	validateRefreshResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.RefreshToken,
		UseCache: false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, validateRefreshResponse.Status)
}

func (s *ValidateTestSuite) TestValidateTokenFromNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse.TokenData.Uuid})

	validateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, validateResponse.Status)

	validateRefreshResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.RefreshToken,
		UseCache: false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, validateRefreshResponse.Status)
}

func (s *ValidateTestSuite) TestValidateDisabledTokenInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	_, err = s.nativeStub.Services.IAM.Token.Disable(ctx, &token.DisableRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(s.T(), err)

	tokenValidateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_DISABLED, tokenValidateResponse.Status)

	tokenValidateResponse, err = s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.RefreshToken,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_DISABLED, tokenValidateResponse.Status)
}

func (s *ValidateTestSuite) TestValidateDisabledTokenInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse.TokenData.Uuid})

	_, err = s.nativeStub.Services.IAM.Token.Disable(ctx, &token.DisableRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(s.T(), err)

	tokenValidateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_DISABLED, tokenValidateResponse.Status)

	tokenValidateResponse, err = s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.RefreshToken,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_DISABLED, tokenValidateResponse.Status)
}

func (s *ValidateTestSuite) TestValidateDeletedTokenInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	_, err = s.nativeStub.Services.IAM.Token.Delete(ctx, &token.DeleteRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(s.T(), err)

	tokenValidateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_NOT_FOUND, tokenValidateResponse.Status)

	tokenValidateResponse, err = s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.RefreshToken,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_NOT_FOUND, tokenValidateResponse.Status)
}

func (s *ValidateTestSuite) TestValidateDeletedTokenInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse.TokenData.Uuid})

	_, err = s.nativeStub.Services.IAM.Token.Delete(ctx, &token.DeleteRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(s.T(), err)

	tokenValidateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_NOT_FOUND, tokenValidateResponse.Status)

	tokenValidateResponse, err = s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.RefreshToken,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_NOT_FOUND, tokenValidateResponse.Status)
}

func (s *ValidateTestSuite) TestValidateInvalidTokenString() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tokenGetResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    tools.GetRandomString(20),
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_INVALID, tokenGetResponse.Status)
}
