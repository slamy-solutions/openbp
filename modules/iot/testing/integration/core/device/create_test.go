package device

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/testing/tools"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type CreateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *CreateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
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

	createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        name,
		Description: description,
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createResponse.Device.Uuid})

	require.Equal(s.T(), name, createResponse.Device.Name)
	require.Equal(s.T(), description, createResponse.Device.Description)

	existResponse, err := s.iotStub.Core.Device.Exists(ctx, &device.ExistsRequest{Namespace: "", Uuid: createResponse.Device.Uuid})
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

	createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   namespaceName,
		Name:        name,
		Description: description,
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: createResponse.Device.Uuid})

	require.Equal(s.T(), name, createResponse.Device.Name)
	require.Equal(s.T(), description, createResponse.Device.Description)

	existResponse, err := s.iotStub.Core.Device.Exists(ctx, &device.ExistsRequest{Namespace: namespaceName, Uuid: createResponse.Device.Uuid})
	require.Nil(s.T(), err)
	require.True(s.T(), existResponse.Exists)
}

func (s *CreateTestSuite) TestCreatesValidIdentity() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createResponse.Device.Uuid})

	identityExistResponse, err := s.nativeStub.Services.IAM.Identity.Exists(ctx, &identity.ExistsIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Device.Identity,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), identityExistResponse.Exists)
}
