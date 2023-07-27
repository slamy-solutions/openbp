package fleet

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	"github.com/slamy-solutions/openbp/modules/iot/testing/tools"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type AddDeviceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *AddDeviceTestSuite) SetupSuite() {
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
func (suite *AddDeviceTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestAddDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(AddDeviceTestSuite))
}

func (s *AddDeviceTestSuite) TestAdd() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createDeviceResponse.Device.Uuid})

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: createFleetResponse.Fleet.Uuid})

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  "",
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)

	listStream, err := s.iotStub.Core.Fleet.ListDevices(ctx, &fleet.ListDevicesRequest{
		Namespace: "",
		Uuid:      createFleetResponse.Fleet.Uuid,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)
	listResponseEntry, err := listStream.Recv()
	require.Nil(s.T(), err)

	require.Equal(s.T(), createDeviceResponse.Device.Uuid, listResponseEntry.Device.Device.Uuid)

	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}

func (s *AddDeviceTestSuite) TestAddInNamespace() {
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

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: createDeviceResponse.Device.Uuid})

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: namespaceName, Uuid: createFleetResponse.Fleet.Uuid})

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  namespaceName,
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)

	listStream, err := s.iotStub.Core.Fleet.ListDevices(ctx, &fleet.ListDevicesRequest{
		Namespace: namespaceName,
		Uuid:      createFleetResponse.Fleet.Uuid,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)
	listResponseEntry, err := listStream.Recv()
	require.Nil(s.T(), err)

	require.Equal(s.T(), createDeviceResponse.Device.Uuid, listResponseEntry.Device.Device.Uuid)

	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}

func (s *AddDeviceTestSuite) TestAddNonExistingDevice() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createDeviceResponse.Device.Uuid})
	require.Nil(s.T(), err)

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: createFleetResponse.Fleet.Uuid})

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  "",
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())

	listStream, err := s.iotStub.Core.Fleet.ListDevices(ctx, &fleet.ListDevicesRequest{
		Namespace: "",
		Uuid:      createFleetResponse.Fleet.Uuid,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)
	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}

func (s *AddDeviceTestSuite) TestAddToNonExistingFleet() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createDeviceResponse.Device.Uuid})

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: createFleetResponse.Fleet.Uuid})
	require.Nil(s.T(), err)

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  "",
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())

	listStream, err := s.iotStub.Core.Fleet.ListDevices(ctx, &fleet.ListDevicesRequest{
		Namespace: "",
		Uuid:      createFleetResponse.Fleet.Uuid,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)
	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}

func (s *AddDeviceTestSuite) TestAddNonExistingDeviceInNamespace() {
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

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: createDeviceResponse.Device.Uuid})
	require.Nil(s.T(), err)

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: namespaceName, Uuid: createFleetResponse.Fleet.Uuid})

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  namespaceName,
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())

	listStream, err := s.iotStub.Core.Fleet.ListDevices(ctx, &fleet.ListDevicesRequest{
		Namespace: namespaceName,
		Uuid:      createFleetResponse.Fleet.Uuid,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)
	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}

func (s *AddDeviceTestSuite) TestAddToNonExistingFleetInNamespace() {
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

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: createDeviceResponse.Device.Uuid})

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   namespaceName,
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: namespaceName, Uuid: createFleetResponse.Fleet.Uuid})
	require.Nil(s.T(), err)

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  namespaceName,
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *AddDeviceTestSuite) TestAddInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createDeviceResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: createDeviceResponse.Device.Uuid})

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: createFleetResponse.Fleet.Uuid})

	namespaceName := tools.GetRandomString(32)

	_, err = s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
		Namespace:  namespaceName,
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}
