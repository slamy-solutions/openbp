package identity

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type CreateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *CreateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *CreateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCreateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}

func (s *CreateTestSuite) TestReturnsDataInReponseToCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)

	r, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            name,
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: r.Identity.Uuid})

	require.Equal(s.T(), name, r.Identity.Name)
	require.True(s.T(), r.Identity.Active)
	require.Len(s.T(), r.Identity.Policies, 0)
	require.Len(s.T(), r.Identity.Roles, 0)
}

func (s *CreateTestSuite) TestAvailableAfterCreationInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	existResponse, err := s.nativeStub.Services.IamIdentity.Exists(ctx, &identity.ExistsIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exists)
}

func (s *CreateTestSuite) TestAvailableAfterCreationInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	existResponse, err := s.nativeStub.Services.IamIdentity.Exists(ctx, &identity.ExistsIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exists)
}
