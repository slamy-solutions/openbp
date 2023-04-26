package policy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type IdentityManagedTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *IdentityManagedTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *IdentityManagedTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestIdentityManagedTestSuite(t *testing.T) {
	suite.Run(t, new(IdentityManagedTestSuite))
}

func (s *IdentityManagedTestSuite) TestIdentityManagedForGlobalNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &policy.CreatePolicyRequest_Identity{Identity: &policy.IdentityManagedData{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
		}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, createResponse.Policy.Managed.(*policy.Policy_Identity).Identity.IdentityUUID)

	getResponse, err := s.nativeStub.Services.IamPolicy.Get(ctx, &policy.GetPolicyRequest{
		Namespace: "",
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, getResponse.Policy.Managed.(*policy.Policy_Identity).Identity.IdentityUUID)
}

func (s *IdentityManagedTestSuite) TestIdentityManagedForNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(20),
		Description: "",
		Managed: &policy.CreatePolicyRequest_Identity{Identity: &policy.IdentityManagedData{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
		}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, createResponse.Policy.Managed.(*policy.Policy_Identity).Identity.IdentityUUID)

	getResponse, err := s.nativeStub.Services.IamPolicy.Get(ctx, &policy.GetPolicyRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), identityCreateResponse.Identity.Uuid, getResponse.Policy.Managed.(*policy.Policy_Identity).Identity.IdentityUUID)
}
