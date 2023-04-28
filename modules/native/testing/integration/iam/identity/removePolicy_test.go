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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type RemovePolicyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *RemovePolicyTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService().WithIAMPolicyService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *RemovePolicyTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestRemovePolicyTestSuite(t *testing.T) {
	suite.Run(t, new(RemovePolicyTestSuite))
}

func (s *RemovePolicyTestSuite) TestReturnsActualDataAfterRemoveFromGlobalNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	removePolicyResponse, err := s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)
	require.Len(s.T(), removePolicyResponse.Identity.Policies, 0)

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 0)
}

func (s *RemovePolicyTestSuite) TestReturnsActualDataAfterAddingInNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   namespaceName,
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	removePolicyResponse, err := s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   namespaceName,
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)
	require.Len(s.T(), removePolicyResponse.Identity.Policies, 0)

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 0)
}

func (s *RemovePolicyTestSuite) TestMultipleRemoveInGlobalNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			PolicyNamespace:   "",
			PolicyUUID:        createPolicyResponse.Policy.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 0)
}

func (s *RemovePolicyTestSuite) TestMultipleRemoveInNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   namespaceName,
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	for i := 0; i < 5; i++ {
		_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
			IdentityNamespace: namespaceName,
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			PolicyNamespace:   namespaceName,
			PolicyUUID:        createPolicyResponse.Policy.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 0)
}

func (s *RemovePolicyTestSuite) TestRemoveNonExistingPolicyIsOkForGlobalNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		PolicyNamespace:   "",
		PolicyUUID:        primitive.NewObjectID().Hex(),
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getIdentityResponse.Identity.Policies[0].Uuid)
}

func (s *RemovePolicyTestSuite) TestRemoveForNonExistingPolicyIsOkForNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   namespaceName,
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		PolicyNamespace:   namespaceName,
		PolicyUUID:        primitive.NewObjectID().Hex(),
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	require.Nil(s.T(), err)

	getIdentityResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getIdentityResponse.Identity.Policies[0].Uuid)
}

func (s *RemovePolicyTestSuite) TestRemoveForNonExistingIdentityFailsWithNotFoundErrorForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
		IdentityNamespace: "",
		IdentityUUID:      primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *RemovePolicyTestSuite) TestRemoveForNonExistingIdentityFailsWithNotFoundErrorForNamespace() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		PolicyNamespace:   namespaceName,
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
		IdentityNamespace: namespaceName,
		IdentityUUID:      primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *RemovePolicyTestSuite) TestRemoveForNonExistingNamespaceFailsWithNotFoundError() {
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

	createPolicyResponse, err := s.nativeStub.Services.IamPolicy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamPolicy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createPolicyResponse.Policy.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.RemovePolicy(ctx, &identity.RemovePolicyRequest{
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
		IdentityNamespace: tools.GetRandomString(20),
		IdentityUUID:      primitive.NewObjectID().Hex(),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
