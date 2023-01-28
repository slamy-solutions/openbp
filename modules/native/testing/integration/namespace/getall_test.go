package namespace

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type GetAllNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetAllNamespaceTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetAllNamespaceTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetAllNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(GetAllNamespaceTestSuite))
}

func (s *GetAllNamespaceTestSuite) TestGetCreatesNamespaces() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	namespaceCount := 10
	names := make([]string, namespaceCount)
	for i := 1; i < namespaceCount; i++ {
		names[i] = tools.GetRandomString(20)
	}

	// Delete namespaces after tests
	defer func() {
		for i := 1; i < namespaceCount; i++ {
			s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: names[i]})
		}
	}()

	// Create namespaces
	for i := 1; i < namespaceCount; i++ {
		_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{Name: names[i], FullName: "", Description: ""})
		require.Nil(s.T(), err)
	}

	// Search for namespaces
	findenNames := map[string]struct{}{}
	stream, err := s.nativeStub.Services.Namespace.GetAll(ctx, &namespace.GetAllNamespacesRequest{})
	require.Nil(s.T(), err)
	for {
		data, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.Fail(s.T(), "Failed to get namespaces from service")
		}
		findenNames[data.Namespace.Name] = struct{}{}
	}

	for i := 1; i < namespaceCount; i++ {
		_, ok := findenNames[names[i]]
		require.True(s.T(), ok)
	}
}
