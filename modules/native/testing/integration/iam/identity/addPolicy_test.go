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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

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

	addPolicyResponse, err := s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   "",
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), addPolicyResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, addPolicyResponse.Identity.Policies[0].Uuid)

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getIdentityResponse.Identity.Policies[0].Uuid)
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

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

	addPolicyResponse, err := s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
		PolicyNamespace:   namespaceName,
		PolicyUUID:        createPolicyResponse.Policy.Uuid,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), addPolicyResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, addPolicyResponse.Identity.Policies[0].Uuid)

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getIdentityResponse.Identity.Policies[0].Uuid)
}

func (s *AddPolicyTestSuite) TestMultipleAddingInGlobalNamespace() {
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
		_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
			IdentityNamespace: "",
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			PolicyNamespace:   "",
			PolicyUUID:        createPolicyResponse.Policy.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getIdentityResponse.Identity.Policies[0].Uuid)
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

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
		_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
			IdentityNamespace: namespaceName,
			IdentityUUID:      identityCreateResponse.Identity.Uuid,
			PolicyNamespace:   namespaceName,
			PolicyUUID:        createPolicyResponse.Policy.Uuid,
		})
		require.Nil(s.T(), err)
	}

	getIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      identityCreateResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)

	require.Len(s.T(), getIdentityResponse.Identity.Policies, 1)
	require.Equal(s.T(), createPolicyResponse.Policy.Uuid, getIdentityResponse.Identity.Policies[0].Uuid)
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingPolicyFailsWithFailedPreconditionErrorForGlobalNamespace() {
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

	_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
		PolicyNamespace:   "",
		PolicyUUID:        primitive.NewObjectID().Hex(),
		IdentityNamespace: "",
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
		PolicyNamespace:   namespaceName,
		PolicyUUID:        primitive.NewObjectID().Hex(),
		IdentityNamespace: namespaceName,
		IdentityUUID:      identityCreateResponse.Identity.Uuid,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *AddPolicyTestSuite) TestAddingForNonExistingIdentityFailsWithNotFoundErrorForGlobalNamespace() {
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

	_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
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

func (s *AddPolicyTestSuite) TestAddingForNonExistingIdentityFailsWithNotFoundErrorForNamespace() {
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

	_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
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

func (s *AddPolicyTestSuite) TestAddingForNonExistingNamespaceFailsWithNotFoundError() {
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

	_, err = s.nativeStub.Services.IAM.Identity.AddPolicy(ctx, &identity.AddPolicyRequest{
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
