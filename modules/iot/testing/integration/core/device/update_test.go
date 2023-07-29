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

type UpdateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
}

func (suite *UpdateTestSuite) SetupSuite() {
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
func (suite *UpdateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}

func (s *UpdateTestSuite) TestUpdate() {
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

	timeBeforeUpdate := time.Now().UTC()
	require.Less(s.T(), createResponse.Device.Created.AsTime(), timeBeforeUpdate)

	newDescription := tools.GetRandomString(32)

	updateResponse, err := s.iotStub.Core.Device.Update(ctx, &device.UpdateRequest{
		Namespace:   "",
		Uuid:        createResponse.Device.Uuid,
		Description: newDescription,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newDescription, updateResponse.Device.Description)
	require.Equal(s.T(), createResponse.Device.Uuid, updateResponse.Device.Uuid)
	require.Equal(s.T(), createResponse.Device.Version+1, updateResponse.Device.Version)

	getResponse, err := s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: "",
		Uuid:      createResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)

	require.Equal(s.T(), newDescription, getResponse.Device.Description)
	require.Equal(s.T(), createResponse.Device.Uuid, getResponse.Device.Uuid)
	require.Equal(s.T(), createResponse.Device.Version+1, getResponse.Device.Version)
}

func (s *UpdateTestSuite) TestUpdateInNamespace() {
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

	timeBeforeUpdate := time.Now().UTC()
	require.Less(s.T(), createResponse.Device.Created.AsTime(), timeBeforeUpdate)

	newDescription := tools.GetRandomString(32)

	updateResponse, err := s.iotStub.Core.Device.Update(ctx, &device.UpdateRequest{
		Namespace:   namespaceName,
		Uuid:        createResponse.Device.Uuid,
		Description: newDescription,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), newDescription, updateResponse.Device.Description)
	require.Equal(s.T(), createResponse.Device.Uuid, updateResponse.Device.Uuid)
	require.Equal(s.T(), createResponse.Device.Version+1, updateResponse.Device.Version)

	getResponse, err := s.iotStub.Core.Device.Get(ctx, &device.GetRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Device.Uuid,
	})
	require.Nil(s.T(), err)

	require.Equal(s.T(), newDescription, getResponse.Device.Description)
	require.Equal(s.T(), createResponse.Device.Uuid, getResponse.Device.Uuid)
	require.Equal(s.T(), createResponse.Device.Version+1, getResponse.Device.Version)
}

func (s *UpdateTestSuite) TestUpdateNonExisting() {
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

	_, err = s.iotStub.Core.Device.Update(ctx, &device.UpdateRequest{
		Namespace:   "",
		Uuid:        createResponse.Device.Uuid,
		Description: tools.GetRandomString(32),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *UpdateTestSuite) TestUpdateNonExistingInNamespace() {
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

	_, err = s.iotStub.Core.Device.Update(ctx, &device.UpdateRequest{
		Namespace:   namespaceName,
		Uuid:        createResponse.Device.Uuid,
		Description: tools.GetRandomString(32),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}

func (s *UpdateTestSuite) TestUpdateInNonExistingNamespace() {
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

	_, err = s.iotStub.Core.Device.Update(ctx, &device.UpdateRequest{
		Namespace:   tools.GetRandomString(32),
		Uuid:        createResponse.Device.Uuid,
		Description: tools.GetRandomString(32),
	})
	require.NotNil(s.T(), err)
	st, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), codes.NotFound, st.Code())
}
