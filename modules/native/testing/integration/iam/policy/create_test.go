package policy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type CreateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *CreateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *CreateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCreateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}

func (s *CreateTestSuite) TestParamsValidation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for _, tc := range writeParamsValidationTestCases {
		s.Run(tc.testName, func() {
			r, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
				Namespace:            "",
				Name:                 tc.name,
				Description:          tc.description,
				Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
				NamespaceIndependent: false,
				Resources:            tc.resources,
				Actions:              tc.actions,
			})
			if err == nil {
				defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: r.Policy.Uuid})
			}

			assert.Equal(s.T(), tc.ok, err == nil)
			if err != nil {
				e, ok := status.FromError(err)
				require.True(s.T(), ok)
				require.Equal(s.T(), e.Code(), codes.InvalidArgument)
			}
		})
	}
}

func (s *CreateTestSuite) TestReturnsDataInReponseToCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	description := tools.GetRandomString(20)
	resources := []string{tools.GetRandomString(20)}
	actions := []string{tools.GetRandomString(20)}

	r, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 name,
		Description:          description,
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            resources,
		Actions:              actions,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: r.Policy.Uuid})

	require.Equal(s.T(), name, r.Policy.Name)
	require.Equal(s.T(), description, r.Policy.Description)
	require.Len(s.T(), r.Policy.Actions, 1)
	require.Equal(s.T(), actions[0], r.Policy.Actions[0])
	require.Len(s.T(), r.Policy.Resources, 1)
	require.Equal(s.T(), resources[0], r.Policy.Resources[0])
}

func (s *CreateTestSuite) TestAvailableAfterCreationInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            "",
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: createResponse.Policy.Uuid})

	existResponse, err := s.nativeStub.Services.IAM.Policy.Exist(ctx, &policy.ExistPolicyRequest{
		Namespace: "",
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exist)
}

func (s *CreateTestSuite) TestAvailableAfterCreationInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
		Namespace:            namespaceName,
		Name:                 tools.GetRandomString(20),
		Description:          tools.GetRandomString(20),
		Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		NamespaceIndependent: false,
		Resources:            []string{},
		Actions:              []string{},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: createResponse.Policy.Uuid})

	existResponse, err := s.nativeStub.Services.IAM.Policy.Exist(ctx, &policy.ExistPolicyRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Policy.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exist)
}
