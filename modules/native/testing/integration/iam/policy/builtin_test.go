package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type BuiltInTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *BuiltInTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *BuiltInTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestBuiltInTestSuite(t *testing.T) {
	suite.Run(t, new(BuiltInTestSuite))
}

func (s *BuiltInTestSuite) TestGetForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
		Namespace: "",
		Type:      policy.BuiltInPolicyType_GLOBAL_ROOT,
	})
	require.Nil(s.T(), err)
}

func (s *BuiltInTestSuite) TestGetForNamespaceRoot() {
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

	time.Sleep(time.Millisecond * 100)

	_, err = s.nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
		Namespace: namespaceName,
		Type:      policy.BuiltInPolicyType_NAMESPACE_ROOT,
	})
	require.Nil(s.T(), err)
}

func (s *BuiltInTestSuite) TestGetForNamespaceEmpty() {
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

	time.Sleep(time.Millisecond * 100)

	_, err = s.nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
		Namespace: namespaceName,
		Type:      policy.BuiltInPolicyType_EMPTY,
	})
	require.Nil(s.T(), err)
}

func (s *BuiltInTestSuite) TestFailsWithNotFoundErrorForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IamPolicy.GetBuiltInPolicy(ctx, &policy.GetBuiltInPolicyRequest{
		Namespace: tools.GetRandomString(20),
		Type:      policy.BuiltInPolicyType_EMPTY,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt exist")
	}
}
