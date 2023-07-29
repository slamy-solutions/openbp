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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
)

type CountTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *CountTestSuite) SetupSuite() {
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
func (suite *CountTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCountTestSuite(t *testing.T) {
	suite.Run(t, new(CountTestSuite))
}

func (s *CountTestSuite) TestCount() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	devicesCount := 10
	devices := make([]string, 0, devicesCount)
	defer func() {
		for i := 0; i < devicesCount; i++ {
			s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: "", Uuid: devices[i]})
		}
	}()

	for i := 0; i < 10; i++ {
		createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
			Namespace:   "",
			Name:        tools.GetRandomString(32),
			Description: tools.GetRandomString(32),
		})
		require.Nil(s.T(), err)
		devices = append(devices, createResponse.Device.Uuid)
	}

	countResponse, err := s.iotStub.Core.Device.Count(ctx, &device.CountRequest{Namespace: ""})
	require.Nil(s.T(), err)
	require.True(s.T(), countResponse.Count >= uint64(devicesCount))
}

func (s *CountTestSuite) TestCountInNamespace() {
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

	devicesCount := 10
	devices := make([]string, 0, devicesCount)
	defer func() {
		for i := 0; i < devicesCount; i++ {
			s.iotStub.Core.Device.Delete(context.Background(), &device.DeleteRequest{Namespace: namespaceName, Uuid: devices[i]})
		}
	}()

	for i := 0; i < 10; i++ {
		createResponse, err := s.iotStub.Core.Device.Create(ctx, &device.CreateRequest{
			Namespace:   namespaceName,
			Name:        tools.GetRandomString(32),
			Description: tools.GetRandomString(32),
		})
		require.Nil(s.T(), err)
		devices = append(devices, createResponse.Device.Uuid)
	}

	countResponse, err := s.iotStub.Core.Device.Count(ctx, &device.CountRequest{Namespace: namespaceName})
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), devicesCount, countResponse.Count)
}

func (s *CountTestSuite) TestCountInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	countResponse, err := s.iotStub.Core.Device.Count(ctx, &device.CountRequest{Namespace: tools.GetRandomString(32)})
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), 0, countResponse.Count)
}
