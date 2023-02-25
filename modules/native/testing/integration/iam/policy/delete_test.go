package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type DeleteTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *DeleteTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService())
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

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})

	deleteResponse, err := s.nativeStub.Services.IamPolicy.Delete(ctx, &policy.DeletePolicyRequest{
		Namespace: "",
		Uuid:      createResponse.Policy.Uuid,
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

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})

	deleteResponse, err := s.nativeStub.Services.IamPolicy.Delete(ctx, &policy.DeletePolicyRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingPolicyFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})

	deleteResponse, err := s.nativeStub.Services.IamPolicy.Delete(ctx, &policy.DeletePolicyRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingPolicyFromNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})

	deleteResponse, err := s.nativeStub.Services.IamPolicy.Delete(ctx, &policy.DeletePolicyRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingPolicyFromNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	deleteResponse, err := s.nativeStub.Services.IamPolicy.Delete(ctx, &policy.DeletePolicyRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}
