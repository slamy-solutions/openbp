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

type ExistsTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ExistsTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ExistsTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestExistsTestSuite(t *testing.T) {
	suite.Run(t, new(ExistsTestSuite))
}

func (s *ExistsTestSuite) TestExistsInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	existsResponse, err := s.nativeStub.Services.IAM.Authentication.Password.Exists(ctx, &password.ExistsRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), existsResponse.Exists)

	_, err = s.nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Authentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	existsResponse, err = s.nativeStub.Services.IAM.Authentication.Password.Exists(ctx, &password.ExistsRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), existsResponse.Exists)
}

func (s *ExistsTestSuite) TestExistsInNamespace() {
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
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
	})

	existsResponse, err := s.nativeStub.Services.IAM.Authentication.Password.Exists(ctx, &password.ExistsRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), existsResponse.Exists)

	pwd := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx, &password.CreateOrUpdateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Password:  pwd,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Authentication.Password.Delete(context.Background(), &password.DeleteRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})

	existsResponse, err = s.nativeStub.Services.IAM.Authentication.Password.Exists(ctx, &password.ExistsRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), existsResponse.Exists)
}

func (s *ExistsTestSuite) TestDoesntExistForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existsResponse, err := s.nativeStub.Services.IAM.Authentication.Password.Exists(ctx, &password.ExistsRequest{
		Namespace: tools.GetRandomString(20),
		Identity:  primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), existsResponse.Exists)
}
