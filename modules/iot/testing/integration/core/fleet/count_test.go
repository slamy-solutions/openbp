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

	fleetsCount := 10
	fleets := make([]string, 0, fleetsCount)
	defer func() {
		for i := 0; i < fleetsCount; i++ {
			s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: "", Uuid: fleets[i]})
		}
	}()

	for i := 0; i < 10; i++ {
		createResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
			Namespace:   "",
			Name:        tools.GetRandomString(32),
			Description: tools.GetRandomString(32),
		})
		require.Nil(s.T(), err)
		fleets = append(fleets, createResponse.Fleet.Uuid)
	}

	countResponse, err := s.iotStub.Core.Fleet.Count(ctx, &fleet.CountRequest{Namespace: ""})
	require.Nil(s.T(), err)
	require.True(s.T(), fleetsCount <= int(countResponse.Count))
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

	fleetsCount := 10
	fleets := make([]string, 0, fleetsCount)
	defer func() {
		for i := 0; i < fleetsCount; i++ {
			s.iotStub.Core.Fleet.Delete(context.Background(), &fleet.DeleteRequest{Namespace: namespaceName, Uuid: fleets[i]})
		}
	}()

	for i := 0; i < 10; i++ {
		createResponse, err := s.iotStub.Core.Fleet.Create(ctx, &fleet.CreateRequest{
			Namespace:   namespaceName,
			Name:        tools.GetRandomString(32),
			Description: tools.GetRandomString(32),
		})
		require.Nil(s.T(), err)
		fleets = append(fleets, createResponse.Fleet.Uuid)
	}

	countResponse, err := s.iotStub.Core.Fleet.Count(ctx, &fleet.CountRequest{Namespace: namespaceName})
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), fleetsCount, countResponse.Count)
}

func (s *CountTestSuite) TestCountInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	countResponse, err := s.iotStub.Core.Fleet.Count(ctx, &fleet.CountRequest{Namespace: tools.GetRandomString(32)})
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), 0, countResponse.Count)
}
