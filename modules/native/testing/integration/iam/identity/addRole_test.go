package identity

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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type AddRoleTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *AddRoleTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService().WithIAMRoleService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *AddRoleTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestAddRoleTestSuite(t *testing.T) {
	suite.Run(t, new(AddRoleTestSuite))
}

func (s *AddRoleTestSuite) TestReturnsActualDataAfterAddingInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	addRoleResponse, err := s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), addRoleResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, addRoleResponse.Identity.Roles[0].Uuid)

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, getIdentityResponse.Identity.Roles[0].Uuid)
}

func (s *AddRoleTestSuite) TestReturnsActualDataAfterAddingInNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	addRoleResponse, err := s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     namespaceName,
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), addRoleResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, addRoleResponse.Identity.Roles[0].Uuid)

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, getIdentityResponse.Identity.Roles[0].Uuid)
}

func (s *AddRoleTestSuite) TestMultipleAddingInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			RoleNamespace:     "",
			RoleUUID:          createRoleResponse.Role.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, getIdentityResponse.Identity.Roles[0].Uuid)
}

func (s *AddRoleTestSuite) TestMultipleAddingInNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
			IdentityNamespace: namespaceName,
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			RoleNamespace:     namespaceName,
			RoleUUID:          createRoleResponse.Role.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, getIdentityResponse.Identity.Roles[0].Uuid)
}

func (s *AddRoleTestSuite) TestAddingForNonExistingRoleFailsWithFailedPreconditionErrorForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		RoleNamespace:     "",
		RoleUUID:          primitive.NewObjectID().Hex(),
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddRoleTestSuite) TestAddingForNonExistingRoleFailsWithFailedPreconditionErrorForNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		RoleNamespace:     namespaceName,
		RoleUUID:          primitive.NewObjectID().Hex(),
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddRoleTestSuite) TestAddingForNonExistingIdentityFailsWithNotFoundErrorForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
		IdentityNamespace: "",
		IdentityUUID:      primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddRoleTestSuite) TestAddingForNonExistingIdentityFailsWithNotFoundErrorForNamespace() {
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

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		RoleNamespace:     namespaceName,
		RoleUUID:          createRoleResponse.Role.Uuid,
		IdentityNamespace: namespaceName,
		IdentityUUID:      primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddRoleTestSuite) TestAddingForNonExistingNamespaceFailsWithNotFoundError() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddRole(ctx, &identity.AddRoleRequest{
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
		IdentityNamespace: tools.GetRandomString(20),
		IdentityUUID:      primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
