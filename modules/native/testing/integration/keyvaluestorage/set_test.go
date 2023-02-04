package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type SetKeyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *SetKeyTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *SetKeyTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestSetKeyTestSuite(t *testing.T) {
	suite.Run(t, new(SetKeyTestSuite))
}

func (s *SetKeyTestSuite) TestFailsToSetInNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespace := tools.GetRandomString(20)
	key := tools.GetRandomString(5)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: namespace,
		Key:       key,
		Value:     []byte("nothing"),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespace, Key: key})

	require.NotNil(s.T(), err)
	e, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), e.Code(), codes.FailedPrecondition)
}

func (s *SetKeyTestSuite) TestSetOverwritesValue() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(30)
	value1 := tools.GetRandomString(30)
	value2 := tools.GetRandomString(30)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value1),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value2),
	})
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       key,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), value2, string(r.Value))
}

func (s *SetKeyTestSuite) TestKeyExistAfterSetInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       key,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), value, string(r.Value))
}

func (s *SetKeyTestSuite) TestKeyExistAfterSetInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create namespace for tests
	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: namespaceName,
	})
	require.Nil(s.T(), err)

	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: namespaceName,
		Key:       key,
		Value:     []byte(value),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespaceName, Key: key})
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: namespaceName,
		Key:       key,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), value, string(r.Value))
}
