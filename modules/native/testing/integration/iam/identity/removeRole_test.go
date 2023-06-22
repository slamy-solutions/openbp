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

type RemoveRoleTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *RemoveRoleTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *RemoveRoleTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestRemoveRoleTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveRoleTestSuite))
}

func (s *RemoveRoleTestSuite) TestReturnsActualDataAfterRemoveFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	removeRoleResponse, err := s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)
	require.Len(s.T(), removeRoleResponse.Identity.Roles, 0)

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 0)
}

func (s *RemoveRoleTestSuite) TestReturnsActualDataAfterAddingInNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     namespaceName,
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	removeRoleResponse, err := s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     namespaceName,
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)
	require.Len(s.T(), removeRoleResponse.Identity.Roles, 0)

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 0)
}

func (s *RemoveRoleTestSuite) TestMultipleRemoveInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			RoleNamespace:     "",
			RoleUUID:          createRoleResponse.Role.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 0)
}

func (s *RemoveRoleTestSuite) TestMultipleRemoveInNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     namespaceName,
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
			IdentityNamespace: namespaceName,
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			RoleNamespace:     namespaceName,
			RoleUUID:          createRoleResponse.Role.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 0)
}

func (s *RemoveRoleTestSuite) TestRemoveNonExistingRoleIsOkForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     "",
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
		RoleNamespace:     "",
		RoleUUID:          primitive.NewObjectID().Hex(),
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, getIdentityResponse.Identity.Roles[0].Uuid)
}

func (s *RemoveRoleTestSuite) TestRemoveForNonExistingRoleIsOkForNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddRole(ctx, &identity.AddRoleRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		RoleNamespace:     namespaceName,
		RoleUUID:          createRoleResponse.Role.Uuid,
	})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
		RoleNamespace:     namespaceName,
		RoleUUID:          primitive.NewObjectID().Hex(),
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Roles, 1)
	require.Equal(s.T(), createRoleResponse.Role.Uuid, getIdentityResponse.Identity.Roles[0].Uuid)
}

func (s *RemoveRoleTestSuite) TestRemoveForNonExistingIdentityFailsWithNotFoundErrorForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
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

func (s *RemoveRoleTestSuite) TestRemoveForNonExistingIdentityFailsWithNotFoundErrorForNamespace() {
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

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
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

func (s *RemoveRoleTestSuite) TestRemoveForNonExistingNamespaceFailsWithNotFoundError() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createRoleResponse, err := s.nativeStub.Services.IAM.Role.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: tools.GetRandomString(20),
		Managed:     &role.CreateRoleRequest_No{No: &role.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Role.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createRoleResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.RemoveRole(ctx, &identity.RemoveRoleRequest{
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
