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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type AddPolicyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *AddPolicyTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *AddPolicyTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestAddPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(AddPolicyTestSuite))
}

func (s *AddPolicyTestSuite) TestReturnsActualDataAfterAddingInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleCreateResponse.Role.Uuid})

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	addPolicyResponse, err := s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		RoleNamespace:   "",
		RoleUUID:        roleCreateResponse.Role.Uuid,
		PolicyNamespace: "",
		PolicyUUID:      createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), addPolicyResponse.Role.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, addPolicyResponse.Role.Policies[0].Uuid)

	getRoleResponse, err := s.nativeStub.Services.IAM.Role.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      roleCreateResponse.Role.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getRoleResponse.Role.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getRoleResponse.Role.Policies[0].Uuid)
}

func (s *AddPolicyTestSuite) TestReturnsActualDataAfterAddingInNamespace() {
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

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: roleCreateResponse.Role.Uuid})

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	addPolicyResponse, err := s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		RoleNamespace:   namespaceName,
		RoleUUID:        roleCreateResponse.Role.Uuid,
		PolicyNamespace: namespaceName,
		PolicyUUID:      createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), addPolicyResponse.Role.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, addPolicyResponse.Role.Policies[0].Uuid)

	getRoleResponse, err := s.nativeStub.Services.IAM.Role.Get(ctx, &role.GetRoleRequest{
		Namespace: namespaceName,
		Uuid:      roleCreateResponse.Role.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getRoleResponse.Role.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getRoleResponse.Role.Policies[0].Uuid)
}

func (s *AddPolicyTestSuite) TestMultipleAddingInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleCreateResponse.Role.Uuid})

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
			RoleNamespace:   "",
			RoleUUID:        roleCreateResponse.Role.Uuid,
			PolicyNamespace: "",
			PolicyUUID:      createPolicyResponse.Policy.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getRoleResponse, err := s.nativeStub.Services.IAM.Role.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      roleCreateResponse.Role.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getRoleResponse.Role.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getRoleResponse.Role.Policies[0].Uuid)
}

func (s *AddPolicyTestSuite) TestMultipleAddingInNamespace() {
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

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: roleCreateResponse.Role.Uuid})

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
			RoleNamespace:   namespaceName,
			RoleUUID:        roleCreateResponse.Role.Uuid,
			PolicyNamespace: namespaceName,
			PolicyUUID:      createPolicyResponse.Policy.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getRoleResponse, err := s.nativeStub.Services.IAM.Role.Get(ctx, &role.GetRoleRequest{
		Namespace: namespaceName,
		Uuid:      roleCreateResponse.Role.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getRoleResponse.Role.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getRoleResponse.Role.Policies[0].Uuid)
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingPolicyFailsWithFailedPreconditionErrorForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleCreateResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		PolicyNamespace: "",
		PolicyUUID:      primitive.NewObjectID().Hex(),
		RoleNamespace:   "",
		RoleUUID:        roleCreateResponse.Role.Uuid,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingPolicyFailsWithFailedPreconditionErrorForNamespace() {
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

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: roleCreateResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		PolicyNamespace: namespaceName,
		PolicyUUID:      primitive.NewObjectID().Hex(),
		RoleNamespace:   namespaceName,
		RoleUUID:        roleCreateResponse.Role.Uuid,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingRoleFailsWithNotFoundErrorForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		PolicyNamespace: "",
		PolicyUUID:      createPolicyResponse.Policy.Uuid,
		RoleNamespace:   "",
		RoleUUID:        primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingRoleFailsWithNotFoundErrorForNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		PolicyNamespace: namespaceName,
		PolicyUUID:      createPolicyResponse.Policy.Uuid,
		RoleNamespace:   namespaceName,
		RoleUUID:        primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingNamespaceFailsWithNotFoundError() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	roleCreateResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: roleCreateResponse.Role.Uuid})

	createPolicyResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IAM.Role.AddPolicy(ctx, &role.AddPolicyRequest{
		PolicyNamespace: "",
		PolicyUUID:      createPolicyResponse.Policy.Uuid,
		RoleNamespace:   tools.GetRandomString(20),
		RoleUUID:        primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
