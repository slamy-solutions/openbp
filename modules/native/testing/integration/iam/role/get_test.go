package role

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type GetTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMRoleService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetTestSuite(t *testing.T) {
	suite.Run(t, new(GetTestSuite))
}

func (s *GetTestSuite) TestGetsDataInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        name,
		Description: description,
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	getResponse, err := s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      createResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	assert.Equal(s.T(), name, getResponse.Role.Name)
	assert.Equal(s.T(), description, getResponse.Role.Description)
	require.Len(s.T(), getResponse.Role.Policies, 0)
}

func (s *GetTestSuite) TestGetsDataInNamespace() {
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

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        name,
		Description: description,
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	getResponse, err := s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	assert.Equal(s.T(), name, getResponse.Role.Name)
	assert.Equal(s.T(), description, getResponse.Role.Description)
	require.Len(s.T(), getResponse.Role.Policies, 0)
}

func (s *GetTestSuite) TestFailsWithNotFoundErrorWhenRoleDoesntExistInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *GetTestSuite) TestFailsWithNotFoundErrorWhenRoleDoesntExistInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *GetTestSuite) TestFailsWithNotFoundErrorWhenNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
