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

type RefreshTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *RefreshTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *RefreshTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestRefreshTestSuite(t *testing.T) {
	suite.Run(t, new(RefreshTestSuite))
}

func (s *RefreshTestSuite) TestRefreshTokenFromGlobalNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.RefreshToken,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_OK, refreshResponse.Status)

	validateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    refreshResponse.Token,
		UseCache: false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, validateResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshTokenFromNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.RefreshToken,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_OK, refreshResponse.Status)

	validateResponse, err := s.nativeStub.Services.IAM.Token.Validate(ctx, &token.ValidateRequest{
		Token:    refreshResponse.Token,
		UseCache: false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, validateResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshDisabledTokenInGlobalNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.RefreshToken,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_DISABLED, refreshResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshDisabledTokenInNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.RefreshToken,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_DISABLED, refreshResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshDeletedTokenInGlobalNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.RefreshToken,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_NOT_FOUND, refreshResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshDeletedTokenInNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.RefreshToken,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_NOT_FOUND, refreshResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshNotRefreshTokenFromGlobalNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.Token,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_NOT_REFRESH_TOKEN, refreshResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshNotRefreshTokenFromNamespace() {
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

	refreshResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tokenCreateResponse.Token,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_NOT_REFRESH_TOKEN, refreshResponse.Status)
}

func (s *RefreshTestSuite) TestRefreshInvalidTokenString() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tokenGetResponse, err := s.nativeStub.Services.IAM.Token.Refresh(ctx, &token.RefreshRequest{
		RefreshToken: tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.RefreshResponse_INVALID, tokenGetResponse.Status)
}
