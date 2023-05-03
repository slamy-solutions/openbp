package password

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/system/testing/tools"
)

type CreateOrUpdateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *CreateOrUpdateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService().WithIAMAuthenticationService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *CreateOrUpdateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCreateOrUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrUpdateTestSuite))
}

func (s *CreateOrUpdateTestSuite) TestCreateAndUpdateInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
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

	pwd1 := tools.GetRandomString(20)
	passwordResponse, err := s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd1,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.True(s.T(), passwordResponse.Created)

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd1,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), authResponse.Authenticated)

	pwd2 := tools.GetRandomString(20)
	passwordResponse, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd2,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.False(s.T(), passwordResponse.Created)

	authResponse, err = s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd1,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)

	authResponse, err = s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd2,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), authResponse.Authenticated)
}

func (s *CreateOrUpdateTestSuite) TestCreateAndUpdateInNamespace() {
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

	pwd1 := tools.GetRandomString(20)
	passwordResponse, err := s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd1,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.True(s.T(), passwordResponse.Created)

	authResponse, err := s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd1,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), authResponse.Authenticated)

	pwd2 := tools.GetRandomString(20)
	passwordResponse, err = s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd2,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.False(s.T(), passwordResponse.Created)

	authResponse, err = s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd1,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), authResponse.Authenticated)

	authResponse, err = s.nativeStub.Services.IamAuthentication.Password.Authenticate(ctx, &password.AuthenticateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd2,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), authResponse.Authenticated)
}

func (s *CreateOrUpdateTestSuite) TestCreateInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pwd := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: tools.GetRandomString(20),
		Identity:  primitive.NewObjectID().Hex(),
		Password:  pwd,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *CreateOrUpdateTestSuite) TestCreateForNonExistingIdentityIsOK() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityUUID := primitive.NewObjectID().Hex()

	passwordResponse, err := s.nativeStub.Services.IamAuthentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityUUID,
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamAuthentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: "",
		Identity:  identityUUID,
	})
	require.True(s.T(), passwordResponse.Created)
}
