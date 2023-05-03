package password

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/system/testing/tools"
)

type AuthenticateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *AuthenticateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService().WithIAMAuthenticationService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *AuthenticateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestAuthenticateTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticateTestSuite))
}

func (s *AuthenticateTestSuite) TestAuthenticateInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	pwd := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), authResponse.Authenticated)
}

func (s *AuthenticateTestSuite) TestAuthenticateInNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	pwd := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), authResponse.Authenticated)
}

func (s *AuthenticateTestSuite) TestAuthenticateForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: tools.GetRandomString(20),
		Identity:  primitive.NewObjectID().Hex(),
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)
}

func (s *AuthenticateTestSuite) TestAuthenticateForNonExistingIdentityInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  primitive.NewObjectID().Hex(),
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)
}

func (s *AuthenticateTestSuite) TestAuthenticateNonExistingIdentityForNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	pwd := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: namespaceName,
		Identity:  primitive.NewObjectID().Hex(),
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)
}

func (s *AuthenticateTestSuite) TestAuthenticateBadPasswordInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	pwd := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)
}

func (s *AuthenticateTestSuite) TestAuthenticateBadPasswordForNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	pwd := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)
}
