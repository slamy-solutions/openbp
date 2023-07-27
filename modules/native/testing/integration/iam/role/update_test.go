package role

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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type UpdateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *UpdateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *UpdateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}

func (s *UpdateTestSuite) TestUpdateInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        name,
		Description: description,
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	newName := tools.GetRandomString(20)
	newDescription := tools.GetRandomString(20)

	updateResponse, err := s.nativeStub.Services.IAM.Role.Update(ctx, &role.UpdateRoleRequest{
		Namespace:      "",
		Uuid:           createResponse.Role.Uuid,
		NewName:        newName,
		NewDescription: newDescription,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, updateResponse.Role.Name)
	require.Equal(s.T(), newDescription, updateResponse.Role.Description)

	getResponse, err := s.nativeStub.Services.IAM.Role.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      createResponse.Role.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, getResponse.Role.Name)
	require.Equal(s.T(), newDescription, getResponse.Role.Description)
}

func (s *UpdateTestSuite) TestUpdateInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        name,
		Description: description,
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	newName := tools.GetRandomString(20)
	newDescription := tools.GetRandomString(20)

	updateResponse, err := s.nativeStub.Services.IAM.Role.Update(ctx, &role.UpdateRoleRequest{
		Namespace:      namespaceName,
		Uuid:           createResponse.Role.Uuid,
		NewName:        newName,
		NewDescription: newDescription,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, updateResponse.Role.Name)
	require.Equal(s.T(), newDescription, updateResponse.Role.Description)

	getResponse, err := s.nativeStub.Services.IAM.Role.Get(ctx, &role.GetRoleRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Role.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newName, getResponse.Role.Name)
	require.Equal(s.T(), newDescription, getResponse.Role.Description)
}

func (s *UpdateTestSuite) TestUpdateNonExistingRoleInGlobalNamespace() {
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

	_, err = s.nativeStub.Services.IAM.Role.Update(ctx, &role.UpdateRoleRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *UpdateTestSuite) TestUpdateNonExistingRoleInNamespace() {
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

	_, err = s.nativeStub.Services.IAM.Role.Update(ctx, &role.UpdateRoleRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *UpdateTestSuite) TestUpdateNonExistingRoleInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IAM.Role.Update(ctx, &role.UpdateRoleRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}
