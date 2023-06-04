package role

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ListTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ListTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMRoleService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ListTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}

func (s *ListTestSuite) TestListsDataInGlobalNamespace() {
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

	listStream, err := s.nativeStub.Services.IamRole.List(ctx, &role.ListRolesRequest{
		Namespace: "",
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)

	listed := map[string]struct{}{}
	for {
		response, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}

		listed[response.Role.Uuid] = struct{}{}
	}

	for _, roleUUID := range created {
		_, ok := listed[roleUUID]
		require.True(s.T(), ok)
	}
}

func (s *ListTestSuite) TestListsDataInNamespace() {
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

	listStream, err := s.nativeStub.Services.IamRole.List(ctx, &role.ListRolesRequest{
		Namespace: namespaceName,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)

	listed := map[string]struct{}{}
	for {
		response, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}

		listed[response.Role.Uuid] = struct{}{}
	}

	for _, roleUUID := range created {
		_, ok := listed[roleUUID]
		require.True(s.T(), ok)
	}
}

func (s *ListTestSuite) TestListSkipAndLimit() {
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

	listStream, err := s.nativeStub.Services.IamRole.List(ctx, &role.ListRolesRequest{
		Namespace: namespaceName,
		Skip:      3,
		Limit:     2,
	})
	require.Nil(s.T(), err)

	listed := []string{}
	for {
		response, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}

		listed = append(listed, response.Role.Uuid)
	}

	require.Len(s.T(), listed, 2)
	require.Equal(s.T(), created[3], listed[0])
	require.Equal(s.T(), created[4], listed[1])
}

func (s *ListTestSuite) TestListInNonExistingNamespaceIsOk() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	listStream, err := s.nativeStub.Services.IamRole.List(ctx, &role.ListRolesRequest{
		Namespace: tools.GetRandomString(20),
		Skip:      3,
		Limit:     2,
	})
	require.Nil(s.T(), err)

	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}
