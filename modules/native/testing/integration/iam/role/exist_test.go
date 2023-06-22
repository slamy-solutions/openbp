package role

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ExistTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ExistTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
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

	createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	existResponse, err := s.nativeStub.Services.IAM.Role.Exist(ctx, &role.ExistRoleRequest{
		Namespace: "",
		Uuid:      createResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.True(s.T(), existResponse.Exist)
}

func (s *ExistTestSuite) TestDoesntExistInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	existResponse, err := s.nativeStub.Services.IAM.Role.Exist(ctx, &role.ExistRoleRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.False(s.T(), existResponse.Exist)
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

	createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	existResponse, err := s.nativeStub.Services.IAM.Role.Exist(ctx, &role.ExistRoleRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.True(s.T(), existResponse.Exist)
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

	createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	existResponse, err := s.nativeStub.Services.IAM.Role.Exist(ctx, &role.ExistRoleRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.False(s.T(), existResponse.Exist)
}

func (s *ExistTestSuite) TestDoesntExistWhenNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existResponse, err := s.nativeStub.Services.IAM.Role.Exist(ctx, &role.ExistRoleRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	require.Nil(s.T(), err)

	require.False(s.T(), existResponse.Exist)
}
