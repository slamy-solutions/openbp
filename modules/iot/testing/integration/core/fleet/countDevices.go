package fleet

import (
	"context"
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

type CountDevicesTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *CountDevicesTestSuite) SetupSuite() {
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
func (suite *CountDevicesTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCountDevicesTestSuite(t *testing.T) {
	suite.Run(t, new(CountDevicesTestSuite))
}

func (s *CountTestSuite) TestCountDevices() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fleetsCount := 5
	devicePerFleetCount := 6
	fleets := make([]string, 0, fleetsCount)
	defer func() {
		for i := 0; i < fleetsCount; i++ {
			s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: fleets[i]})
		}
	}()
	devices := make([]string, 0, fleetsCount*devicePerFleetCount)
	defer func() {
		for i := 0; i < fleetsCount*devicePerFleetCount; i++ {
			s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: devices[i]})
		}
	}()

	for fleetId := 0; fleetId < fleetsCount; fleetId++ {
		fleetCreateResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
			Namespace:   "",
			Name:        tools.GetRandomString(32),
			Description: tools.GetRandomString(32),
		})
		require.Nil(s.T(), err)
		fleets = append(fleets, fleetCreateResponse.Fleet.Uuid)

		for deviceId := 0; deviceId < devicePerFleetCount; deviceId++ {
			deviceCreateResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
				Namespace:   "",
				Name:        tools.GetRandomString(32),
				Description: tools.GetRandomString(32),
			})
			require.Nil(s.T(), err)
			devices = append(devices, deviceCreateResponse.Device.Uuid)
		}
	}

	defer func() {
		for fleetId := 0; fleetId < fleetsCount; fleetId++ {
			for deviceId := 0; deviceId < fleetId; deviceId++ {
				s.iotStub.Core.Fleet.RemoveDevice(context.Background(), &fleet.RemoveDeviceRequest{
					Namespace:  "",
					FleetUUID:  fleets[fleetId],
					DeviceUUID: devices[fleetId*devicePerFleetCount+deviceId],
				})
			}
		}
	}()

	for fleetId := 0; fleetId < fleetsCount; fleetId++ {
		for deviceId := 0; deviceId < fleetId; deviceId++ {
			_, err := s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
				Namespace:  "",
				FleetUUID:  fleets[fleetId],
				DeviceUUID: devices[fleetId*devicePerFleetCount+deviceId],
			})
			require.Nil(s.T(), err)
		}
	}

	countResponse, err := s.iotStub.Core.Fleet.CountDevices(ctx, &fleet.CountDevicesRequest{
		Namespace: "",
		Uuid:      fleets[3],
	})
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), 3, countResponse.Count)
}

func (s *CountTestSuite) TestCountDevicesInNamespace() {
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

	fleetsCount := 5
	devicePerFleetCount := 6
	fleets := make([]string, 0, fleetsCount)
	defer func() {
		for i := 0; i < fleetsCount; i++ {
			s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: namespaceName, Uuid: fleets[i]})
		}
	}()
	devices := make([]string, 0, fleetsCount*devicePerFleetCount)
	defer func() {
		for i := 0; i < fleetsCount*devicePerFleetCount; i++ {
			s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: devices[i]})
		}
	}()

	for fleetId := 0; fleetId < fleetsCount; fleetId++ {
		fleetCreateResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
			Namespace:   namespaceName,
			Name:        tools.GetRandomString(32),
			Description: tools.GetRandomString(32),
		})
		require.Nil(s.T(), err)
		fleets = append(fleets, fleetCreateResponse.Fleet.Uuid)

		for deviceId := 0; deviceId < devicePerFleetCount; deviceId++ {
			deviceCreateResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
				Namespace:   namespaceName,
				Name:        tools.GetRandomString(32),
				Description: tools.GetRandomString(32),
			})
			require.Nil(s.T(), err)
			devices = append(devices, deviceCreateResponse.Device.Uuid)
		}
	}

	defer func() {
		for fleetId := 0; fleetId < fleetsCount; fleetId++ {
			for deviceId := 0; deviceId < fleetId; deviceId++ {
				s.iotStub.Core.Fleet.RemoveDevice(context.Background(), &fleet.RemoveDeviceRequest{
					Namespace:  namespaceName,
					FleetUUID:  fleets[fleetId],
					DeviceUUID: devices[fleetId*devicePerFleetCount+deviceId],
				})
			}
		}
	}()

	for fleetId := 0; fleetId < fleetsCount; fleetId++ {
		for deviceId := 0; deviceId < fleetId; deviceId++ {
			_, err := s.iotStub.Core.Fleet.AddDevice(ctx, &fleet.AddDeviceRequest{
				Namespace:  namespaceName,
				FleetUUID:  fleets[fleetId],
				DeviceUUID: devices[fleetId*devicePerFleetCount+deviceId],
			})
			require.Nil(s.T(), err)
		}
	}

	countResponse, err := s.iotStub.Core.Fleet.CountDevices(ctx, &fleet.CountDevicesRequest{
		Namespace: namespaceName,
		Uuid:      fleets[3],
	})
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), 3, countResponse.Count)
}

func (s *CountTestSuite) TestCountDevicesForNonExistingFleet() {
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

	countResponse, err := s.iotStub.Core.Fleet.CountDevices(ctx, &fleet.CountDevicesRequest{
		Namespace: "",
		Uuid:      createFleetResponse.Fleet.Uuid,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 0, countResponse.Count)
}

func (s *CountTestSuite) TestCountDevicesForNonExistingFleetInNamespace() {
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

	countResponse, err := s.iotStub.Core.Fleet.CountDevices(ctx, &fleet.CountDevicesRequest{
		Namespace: namespaceName,
		Uuid:      createFleetResponse.Fleet.Uuid,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 0, countResponse.Count)
}

func (s *CountTestSuite) TestCountDevicesForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createFleetResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
		Namespace:   "",
		Name:        tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	_, err = s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: createFleetResponse.Fleet.Uuid})
	require.Nil(s.T(), err)

	countResponse, err := s.iotStub.Core.Fleet.CountDevices(ctx, &fleet.CountDevicesRequest{
		Namespace: tools.GetRandomString(32),
		Uuid:      createFleetResponse.Fleet.Uuid,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), 0, countResponse.Count)
}
