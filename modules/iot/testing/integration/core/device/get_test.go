package device

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/testing/tools"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type GetTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *GetTestSuite) SetupSuite() {
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
func (suite *GetTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetTestSuite(t *testing.T) {
	suite.Run(t, new(GetTestSuite))
}

func (s *GetTestSuite) TestGet() {
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

	getResponse, err := s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: "",
		Uuid:      createResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)

	require.Equal(s.T(), name, getResponse.Device.Name)
	require.Equal(s.T(), description, getResponse.Device.Description)
	require.Equal(s.T(), createResponse.Device.Uuid, getResponse.Device.Uuid)
}

func (s *GetTestSuite) TestGetFromNamespace() {
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

	getResponse, err := s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)

	require.Equal(s.T(), name, getResponse.Device.Name)
	require.Equal(s.T(), description, getResponse.Device.Description)
	require.Equal(s.T(), createResponse.Device.Uuid, getResponse.Device.Uuid)
}

func (s *GetTestSuite) TestGetNonExisting() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createResponse.Device.Uuid})
	require.Nil(s.T(), err)

	_, err = s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: "",
		Uuid:      createResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *GetTestSuite) TestGetNonExistingFromNamespace() {
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

	createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: createResponse.Device.Uuid})
	require.Nil(s.T(), err)

	_, err = s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *GetTestSuite) TestGetFromNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createResponse.Device.Uuid})
	require.Nil(s.T(), err)

	_, err = s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: tools.GetRandomString(32),
		Uuid:      createResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}
