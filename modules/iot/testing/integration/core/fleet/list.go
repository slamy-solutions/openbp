package fleet

import (
	"context"
	"io"
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

type ListTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *ListTestSuite) SetupSuite() {
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
func (suite *ListTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}

func (s *ListTestSuite) TestList() {
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

	listStream, err := s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: "", Skip: 0, Limit: 0})
	require.Nil(s.T(), err)
	receivedIds := make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets, receivedIds)

	listStream, err = s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: "", Skip: 0, Limit: 3})
	require.Nil(s.T(), err)
	receivedIds = make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets[:3], receivedIds)

	listStream, err = s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: "", Skip: 3, Limit: 0})
	require.Nil(s.T(), err)
	receivedIds = make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets[3:], receivedIds)

	listStream, err = s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: "", Skip: 3, Limit: 4})
	require.Nil(s.T(), err)
	receivedIds = make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets[3:7], receivedIds)
}

func (s *ListTestSuite) TestListInNamespace() {
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

	listStream, err := s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: namespaceName, Skip: 0, Limit: 0})
	require.Nil(s.T(), err)
	receivedIds := make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets, receivedIds)

	listStream, err = s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: namespaceName, Skip: 0, Limit: 3})
	require.Nil(s.T(), err)
	receivedIds = make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets[:3], receivedIds)

	listStream, err = s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: namespaceName, Skip: 3, Limit: 0})
	require.Nil(s.T(), err)
	receivedIds = make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets[3:], receivedIds)

	listStream, err = s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: namespaceName, Skip: 3, Limit: 4})
	require.Nil(s.T(), err)
	receivedIds = make([]string, 0, fleetsCount)
	for {
		d, err := listStream.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedIds = append(receivedIds, d.Fleet.Uuid)
	}
	require.Equal(s.T(), fleets[3:7], receivedIds)
}

func (s *ListTestSuite) TestListInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	listStream, err := s.iotStub.Core.Fleet.List(ctx, &fleet.ListRequest{Namespace: tools.GetRandomString(32), Skip: 0, Limit: 0})
	require.Nil(s.T(), err)
	_, err = listStream.Recv()
	require.Equal(s.T(), io.EOF, err)
}
