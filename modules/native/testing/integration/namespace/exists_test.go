package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ExistsNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ExistsNamespaceTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ExistsNamespaceTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestExistsNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(ExistsNamespaceTestSuite))
}

func (s *ExistsNamespaceTestSuite) TestChecksIfNamespaceExists() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)

	r, err := s.nativeStub.Services.Namespace.Exists(ctx, &namespace.IsNamespaceExistRequest{Name: name})
	require.Nil(s.T(), err)
	require.False(s.T(), r.Exist)

	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})

	r, err = s.nativeStub.Services.Namespace.Exists(ctx, &namespace.IsNamespaceExistRequest{Name: name})
	require.Nil(s.T(), err)
	require.True(s.T(), r.Exist)

	_, err = s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)

	r, err = s.nativeStub.Services.Namespace.Exists(ctx, &namespace.IsNamespaceExistRequest{Name: name})
	require.Nil(s.T(), err)
	require.False(s.T(), r.Exist)
}
