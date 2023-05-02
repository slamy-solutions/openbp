package role

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ServiceManagedTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ServiceManagedTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMRoleService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ServiceManagedTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestServiceManagedTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceManagedTestSuite))
}

func (s *ServiceManagedTestSuite) TestServiceManagedForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	require.Equal(s.T(), managedService, createResponse.Role.Managed.(*role.Role_Service).Service.Service)
	require.Equal(s.T(), managedReason, createResponse.Role.Managed.(*role.Role_Service).Service.Reason)
	require.Equal(s.T(), managedId, createResponse.Role.Managed.(*role.Role_Service).Service.ManagementId)

	getResponse, err := s.nativeStub.Services.IamRole.GetServiceManagedRole(ctx, &role.GetServiceManagedRoleRequest{
		Namespace: "",
		Service:   managedService,
		ManagedId: managedId,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), createResponse.Role.Uuid, getResponse.Role.Uuid)
	require.Equal(s.T(), managedService, getResponse.Role.Managed.(*role.Role_Service).Service.Service)
	require.Equal(s.T(), managedReason, getResponse.Role.Managed.(*role.Role_Service).Service.Reason)
	require.Equal(s.T(), managedId, getResponse.Role.Managed.(*role.Role_Service).Service.ManagementId)
}

func (s *ServiceManagedTestSuite) TestServiceManagedForNamespace() {
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

	time.Sleep(time.Millisecond * 100) // Wait for indexes to be created

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	require.Equal(s.T(), managedService, createResponse.Role.Managed.(*role.Role_Service).Service.Service)
	require.Equal(s.T(), managedReason, createResponse.Role.Managed.(*role.Role_Service).Service.Reason)
	require.Equal(s.T(), managedId, createResponse.Role.Managed.(*role.Role_Service).Service.ManagementId)

	time.Sleep(time.Second) // You need time after namespace creation for this functionality to start working.

	getResponse, err := s.nativeStub.Services.IamRole.GetServiceManagedRole(ctx, &role.GetServiceManagedRoleRequest{
		Namespace: namespaceName,
		Service:   managedService,
		ManagedId: managedId,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), createResponse.Role.Uuid, getResponse.Role.Uuid)
	require.Equal(s.T(), managedService, getResponse.Role.Managed.(*role.Role_Service).Service.Service)
	require.Equal(s.T(), managedReason, getResponse.Role.Managed.(*role.Role_Service).Service.Reason)
	require.Equal(s.T(), managedId, getResponse.Role.Managed.(*role.Role_Service).Service.ManagementId)
}

func (s *ServiceManagedTestSuite) TestServiceManagedFailsWithAlreadyExistForSameMnagementIdForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	createResponse2, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	if err != nil {
		s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse2.Role.Uuid})
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ServiceManagedTestSuite) TestServiceManagedFailsWithAlreadyExistForSameMnagementIdForNamespace() {
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

	time.Sleep(time.Millisecond * 100) // Wait for indexes to be created

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	createResponse2, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	if err != nil {
		s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse2.Role.Uuid})
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ServiceManagedTestSuite) TestServiceManagedFailWithNotFoundForBabManagedIdForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamRole.GetServiceManagedRole(ctx, &role.GetServiceManagedRoleRequest{
		Namespace: "",
		Service:   managedService,
		ManagedId: tools.GetRandomString(20),
		UseCache:  true,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ServiceManagedTestSuite) TestServiceManagedFailWithNotFoundForBabManagedIdForNamespace() {
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

	time.Sleep(time.Millisecond * 100) // Wait for indexes to be created

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	time.Sleep(time.Second) // You need time after namespace creation for this functionality to start working.

	_, err = s.nativeStub.Services.IamRole.GetServiceManagedRole(ctx, &role.GetServiceManagedRoleRequest{
		Namespace: namespaceName,
		Service:   managedService,
		ManagedId: tools.GetRandomString(20),
		UseCache:  true,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ServiceManagedTestSuite) TestServiceManagedFailWithNotFoundForNonExistingNamespace() {
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

	managedService := tools.GetRandomString(20)
	managedReason := tools.GetRandomString(20)
	managedId := tools.GetRandomString(20)

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Service{Service: &role.ServiceManagedData{
			Service:      managedService,
			Reason:       managedReason,
			ManagementId: managedId,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})

	_, err = s.nativeStub.Services.IamRole.GetServiceManagedRole(ctx, &role.GetServiceManagedRoleRequest{
		Namespace: tools.GetRandomString(20),
		Service:   managedService,
		ManagedId: tools.GetRandomString(20),
		UseCache:  true,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ServiceManagedTestSuite) TestIdentityManagedForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: false,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Identity{Identity: &role.IdentityManagedData{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: "", Uuid: createResponse.Role.Uuid})
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, createResponse.Role.Managed.(*role.Role_Identity).Identity.IdentityUUID)

	getResponse, err := s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: "",
		Uuid:      createResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, getResponse.Role.Managed.(*role.Role_Identity).Identity.IdentityUUID)
}

func (s *ServiceManagedTestSuite) TestIdentityManagedForNamespace() {
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
		InitiallyActive: false,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	createResponse, err := s.nativeStub.Services.IamRole.Create(ctx, &role.CreateRoleRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &role.CreateRoleRequest_Identity{Identity: &role.IdentityManagedData{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
		}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamRole.Delete(context.Background(), &role.DeleteRoleRequest{Namespace: namespaceName, Uuid: createResponse.Role.Uuid})
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, createResponse.Role.Managed.(*role.Role_Identity).Identity.IdentityUUID)

	getResponse, err := s.nativeStub.Services.IamRole.Get(ctx, &role.GetRoleRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Role.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, getResponse.Role.Managed.(*role.Role_Identity).Identity.IdentityUUID)
}
