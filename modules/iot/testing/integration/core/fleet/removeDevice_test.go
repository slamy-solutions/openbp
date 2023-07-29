package fleet

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	"github.com/slamy-solutions/openbp/modules/iot/testing/tools"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type RemoveDeviceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *RemoveDeviceTestSuite) SetupSuite() {
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
func (suite *RemoveDeviceTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestRemoveDeviceTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveDeviceTestSuite))
}

func (s *RemoveDeviceTestSuite) TestRemove() {
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

	for i := 0; i < 3; i++ {
		_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
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
		_, err = listStream.Recv()
		require.Equal(s.T(), io.EOF, err)
	}
}

func (s *RemoveDeviceTestSuite) TestRemoveInNamespace() {
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

	for i := 0; i < 3; i++ {
		_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
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
		_, err = listStream.Recv()
		require.Equal(s.T(), io.EOF, err)
	}
}

func (s *RemoveDeviceTestSuite) TestRemoveNonExistingDevice() {
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

	_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
		Namespace:  "",
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)
}

func (s *RemoveDeviceTestSuite) TestRemoveNonExistingDeviceInNamespace() {
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

	_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
		Namespace:  namespaceName,
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)
}

func (s *RemoveDeviceTestSuite) TestRemoveFromNonExistingFleet() {
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

	_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
		Namespace:  "",
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)
}

func (s *RemoveDeviceTestSuite) TestRemoveFromNonExistingFleetInNamespace() {
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

	_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
		Namespace:  namespaceName,
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)
}

func (s *RemoveDeviceTestSuite) TestRemoveInNonExistingNamespace() {
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

	_, err = s.iotStub.Core.Fleet.RemoveDevice(ctx, &fleet.RemoveDeviceRequest{
		Namespace:  tools.GetRandomString(32),
		FleetUUID:  createFleetResponse.Fleet.Uuid,
		DeviceUUID: createDeviceResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)
}
