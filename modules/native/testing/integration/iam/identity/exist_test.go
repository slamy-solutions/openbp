package identity

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ExistTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ExistTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ExistTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestExistTestSuite(t *testing.T) {
	suite.Run(t, new(ExistTestSuite))
}

func (s *ExistTestSuite) TestExistInGlobalNamespace() {
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

func (s *ExistTestSuite) TestDoesntExistInGlobalNamespace() {
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
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.False(s.T(), existResponse.Exists)
}

func (s *ExistTestSuite) TestExistInNamespace() {
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

func (s *ExistTestSuite) TestDoesntExistInNamespace() {
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
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.False(s.T(), existResponse.Exists)
}

func (s *ExistTestSuite) TestDoesntExistWhenNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existResponse, err := s.nativeStub.Services.IamIdentity.Exists(ctx, &identity.ExistsIdentityRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.False(s.T(), existResponse.Exists)
}
