package role

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type CountTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *CountTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMRoleService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *CountTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCountTestSuite(t *testing.T) {
	suite.Run(t, new(CountTestSuite))
}

func (s *CountTestSuite) TestCountsDataInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	created := []string{}
	defer func() {
		for _, roleUUID := range created {
			s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleUUID})
		}
	}()
	for i := 0; i < 10; i++ {
		r, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
			Namespace:   "",
			Name:        tools.GetRandomString(20),
			Description: "",
			Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		created = append(created, r.Role.Uuid)
	}

	response, err := s.nativeStub.Services.IamRole.Count(ctx, &role.CountRolesRequest{Namespace: "", UseCache: false})
	require.Nil(s.T(), err)
	require.GreaterOrEqual(s.T(), response.Count, 10)
}

func (s *CountTestSuite) TestCountsDataInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    tools.GetRandomString(10),
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	response, err := s.nativeStub.Services.IamRole.Count(ctx, &role.CountRolesRequest{Namespace: namespaceName, UseCache: true})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 0, response.Count)

	created := []string{}
	defer func() {
		for _, roleUUID := range created {
			s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: roleUUID})
		}
	}()
	for i := 0; i < 10; i++ {
		r, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
			Namespace:   namespaceName,
			Name:        tools.GetRandomString(20),
			Description: "",
			Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		created = append(created, r.Role.Uuid)
	}

	response, err = s.nativeStub.Services.IamRole.Count(ctx, &role.CountRolesRequest{Namespace: namespaceName, UseCache: true})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 10, response.Count)

	for i := 0; i < 7; i++ {
		r, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
			Namespace:   namespaceName,
			Name:        tools.GetRandomString(20),
			Description: "",
			Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		created = append(created, r.Role.Uuid)
	}

	response, err = s.nativeStub.Services.IamRole.Count(ctx, &role.CountRolesRequest{Namespace: namespaceName, UseCache: true})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 17, response.Count)
}

func (s *CountTestSuite) TestCountInNonExistingNamespaceIsOk() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	response, err := s.nativeStub.Services.IamRole.Count(ctx, &role.CountRolesRequest{Namespace: tools.GetRandomString(20), UseCache: true})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 0, response.Count)
}
