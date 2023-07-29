package fleet

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	"github.com/slamy-solutions/openbp/modules/iot/testing/tools"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type CreateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *CreateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}

	suite.iotStub = iot.NewIOTStub(iot.NewStubConfig().WithCoreService())
	err = suite.iotStub.Connect()
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

func (s *CreateTestSuite) TestCreate() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(32)
	description := tools.GetRandomString(32)

	createResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   "",
		Name:        name,
		Description: description,
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: createResponse.Fleet.Uuid})

	require.Equal(s.T(), name, createResponse.Fleet.Name)
	require.Equal(s.T(), description, createResponse.Fleet.Description)

	existResponse, err := s.iotStub.Core.Fleet.Exists(ctx, &fleet.ExistsRequest{Namespace: "", Uuid: createResponse.Fleet.Uuid})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exists)
}

func (s *CreateTestSuite) TestCreateInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(32)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    tools.GetRandomString(10),
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	name := tools.GetRandomString(32)
	description := tools.GetRandomString(32)

	createResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   namespaceName,
		Name:        name,
		Description: description,
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: namespaceName, Uuid: createResponse.Fleet.Uuid})

	require.Equal(s.T(), name, createResponse.Fleet.Name)
	require.Equal(s.T(), description, createResponse.Fleet.Description)

	existResponse, err := s.iotStub.Core.Fleet.Exists(ctx, &fleet.ExistsRequest{Namespace: namespaceName, Uuid: createResponse.Fleet.Uuid})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exists)
}
