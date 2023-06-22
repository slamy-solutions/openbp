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

type DeleteTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *DeleteTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *DeleteTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}

func (s *DeleteTestSuite) TestDeleteFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	deleteResponse, err := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteFromNamespace() {
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

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	deleteResponse, err := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingIdentityFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	deleteResponse, err := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingIdentityFromNamespace() {
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

	createResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	deleteResponse, err := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingIdentityFromNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	deleteResponse, err := s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}
