package namespace

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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type GetNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetNamespaceTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetNamespaceTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(GetNamespaceTestSuite))
}

func (s *GetNamespaceTestSuite) TestGetInfoAboutNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	fullName := tools.GetRandomString(30)
	description := tools.GetRandomString(30)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    fullName,
		Description: description,
	})
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})

	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: name})
	require.Nil(s.T(), err)

	assert.Equal(s.T(), r.Namespace.Name, name)
	assert.Equal(s.T(), r.Namespace.FullName, fullName)
	assert.Equal(s.T(), r.Namespace.Description, description)
}

func (s *GetNamespaceTestSuite) TestFailsWhileGettingNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: tools.GetRandomString(20)})
	require.NotNil(s.T(), err)
	e, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), e.Code(), codes.NotFound)
}
