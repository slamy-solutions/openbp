package policy

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ListTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ListTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
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
		for _, policyUUID := range created {
			s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: "", Uuid: policyUUID})
		}
	}()
	for i := 0; i < 10; i++ {
		r, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
			Namespace:            "",
			Name:                 tools.GetRandomString(20),
			Description:          "",
			NamespaceIndependent: false,
			Resources:            []string{},
			Actions:              []string{},
			Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		created = append(created, r.Policy.Uuid)
	}

	listStream, err := s.nativeStub.Services.IAM.Policy.List(ctx, &policy.ListPoliciesRequest{
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

		listed[response.Policy.Uuid] = struct{}{}
	}

	for _, policyUUID := range created {
		_, ok := listed[policyUUID]
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
		for _, policyUUID := range created {
			s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: policyUUID})
		}
	}()
	for i := 0; i < 10; i++ {
		r, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
			Namespace:            namespaceName,
			Name:                 tools.GetRandomString(20),
			Description:          "",
			NamespaceIndependent: false,
			Resources:            []string{},
			Actions:              []string{},
			Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		created = append(created, r.Policy.Uuid)
	}

	listStream, err := s.nativeStub.Services.IAM.Policy.List(ctx, &policy.ListPoliciesRequest{
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

		listed[response.Policy.Uuid] = struct{}{}
	}

	for _, policyUUID := range created {
		_, ok := listed[policyUUID]
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
		for _, policyUUID := range created {
			s.nativeStub.Services.IAM.Policy.Delete(context.Background(), &policy.DeletePolicyRequest{Namespace: namespaceName, Uuid: policyUUID})
		}
	}()
	for i := 0; i < 10; i++ {
		r, err := s.nativeStub.Services.IAM.Policy.Create(ctx, &policy.CreatePolicyRequest{
			Namespace:            namespaceName,
			Name:                 tools.GetRandomString(20),
			Description:          "",
			NamespaceIndependent: false,
			Resources:            []string{},
			Actions:              []string{},
			Managed:              &policy.CreatePolicyRequest_No{No: &policy.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		created = append(created, r.Policy.Uuid)
	}

	listStream, err := s.nativeStub.Services.IAM.Policy.List(ctx, &policy.ListPoliciesRequest{
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

		listed = append(listed, response.Policy.Uuid)
	}

	require.Len(s.T(), listed, 2)
	require.Equal(s.T(), created[3], listed[0])
	require.Equal(s.T(), created[4], listed[1])
}

func (s *ListTestSuite) TestListInNonExistingNamespaceIsOk() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	listStream, err := s.nativeStub.Services.IAM.Policy.List(ctx, &policy.ListPoliciesRequest{
		Namespace: tools.GetRandomString(20),
		Skip:      3,
		Limit:     2,
	})
	require.Nil(s.T(), err)

	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}
