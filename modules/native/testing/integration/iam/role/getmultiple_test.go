package role

import (
	"context"
	"io"
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

type GetMultipleTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetMultipleTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMRoleService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetMultipleTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetMultipleTestSuite(t *testing.T) {
	suite.Run(t, new(GetMultipleTestSuite))
}

func (s *GetMultipleTestSuite) TestGetMultiple() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName1 := tools.GetRandomString(20)
	namespaceName2 := tools.GetRandomString(20)

	var roles = []*struct {
		namespace string
		create    bool
		search    bool
		uuid      string
	}{
		{"", true, true, primitive.NewObjectID().Hex()},
		{"", true, true, primitive.NewObjectID().Hex()},
		{"", true, false, primitive.NewObjectID().Hex()},
		{"", true, false, primitive.NewObjectID().Hex()},
		{"", false, true, primitive.NewObjectID().Hex()},
		{"", false, true, primitive.NewObjectID().Hex()},

		{namespaceName1, true, true, primitive.NewObjectID().Hex()},
		{namespaceName1, true, true, primitive.NewObjectID().Hex()},
		{namespaceName1, true, false, primitive.NewObjectID().Hex()},
		{namespaceName1, true, false, primitive.NewObjectID().Hex()},
		{namespaceName1, false, true, primitive.NewObjectID().Hex()},
		{namespaceName1, false, true, primitive.NewObjectID().Hex()},

		{namespaceName2, true, true, primitive.NewObjectID().Hex()},
		{namespaceName2, true, true, primitive.NewObjectID().Hex()},
		{namespaceName2, true, false, primitive.NewObjectID().Hex()},
		{namespaceName2, true, false, primitive.NewObjectID().Hex()},
		{namespaceName2, false, true, primitive.NewObjectID().Hex()},
		{namespaceName2, false, true, primitive.NewObjectID().Hex()},

		{tools.GetRandomString(20), false, true, primitive.NewObjectID().Hex()},
		{tools.GetRandomString(20), false, true, primitive.NewObjectID().Hex()},
	}

	//Create required namespaces
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName1,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(ctx, &namespace.DeleteNamespaceRequest{Name: namespaceName1})

	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName2,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(ctx, &namespace.DeleteNamespaceRequest{Name: namespaceName2})

	// Make sure at the end all roles deleted
	defer func() {
		for _, r := range roles {
			if r.create {
				s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{
					Namespace: r.namespace,
					Uuid:      r.uuid,
				})
			}
		}
	}()

	// Create all the roles
	for _, roleData := range roles {
		if roleData.create {
			r, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
				Namespace:   roleData.namespace,
				Name:        tools.GetRandomString(20),
				Description: tools.GetRandomString(20),
				Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
			})
			require.Nil(s.T(), err)
			roleData.uuid = r.Role.Uuid
		}
	}

	rolesToSearch := make([]*role.GetMultipleRolesRequest_RequestedRole, 0, 12)
	for _, roleData := range roles {
		if roleData.search {
			rolesToSearch = append(rolesToSearch, &role.GetMultipleRolesRequest_RequestedRole{
				Namespace: roleData.namespace,
				Uuid:      roleData.uuid,
			})
		}
	}
	require.Len(s.T(), rolesToSearch, 14)

	r, err := s.nativeStub.Services.IamRole.GetMultiple(ctx, &role.GetMultipleRolesRequest{
		Roles: rolesToSearch,
	})
	require.Nil(s.T(), err)

	receivedRoles := make([]*role.Role, 0, 14)
	for {
		chunk, err := r.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}

		receivedRoles = append(receivedRoles, chunk.Role)
	}
	require.Len(s.T(), receivedRoles, 6)

	existInResponse := func(namespaceName string, uuid string) bool {
		for _, r := range receivedRoles {
			if r.Namespace == namespaceName && r.Uuid == uuid {
				return true
			}
		}
		return false
	}

	existInRequest := func(namespaceName string, uuid string) bool {
		for _, r := range rolesToSearch {
			if r.Namespace == namespaceName && r.Uuid == uuid {
				return true
			}
		}
		return false
	}

	// Validating if request and result are same
	for _, r := range receivedRoles {
		require.True(s.T(), existInRequest(r.Namespace, r.Uuid))
	}
	for _, r := range roles {
		if r.create && r.search {
			require.True(s.T(), existInResponse(r.namespace, r.uuid))
		}
	}
}
