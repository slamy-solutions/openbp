package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type SetKeyIfNotExistTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *SetKeyIfNotExistTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *SetKeyIfNotExistTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestSetKeyIfNotExistTestSuite(t *testing.T) {
	suite.Run(t, new(SetKeyIfNotExistTestSuite))
}

func (s *SetKeyIfNotExistTestSuite) TestValidatesInputs() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(20)
	bigValue := tools.GetRandomString(1024*1024*15 + 1)
	_, err := s.nativeStub.Services.Keyvaluestorage.SetIfNotExist(ctx, &keyvaluestorage.SetIfNotExistRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(bigValue),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})
	require.NotNil(s.T(), err)
	if st, ok := status.FromError(err); ok {
		assert.Equal(s.T(), codes.InvalidArgument, st.Code())
	} else {
		require.Fail(s.T(), "Error expected")
	}
}

func (s *SetKeyIfNotExistTestSuite) TestSetsValue() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Keyvaluestorage.SetIfNotExist(ctx, &keyvaluestorage.SetIfNotExistRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	assert.Equal(s.T(), []byte(value), r.Value)
}

func (s *SetKeyIfNotExistTestSuite) TestSetsValueInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    "",
		Description: "",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Keyvaluestorage.SetIfNotExist(ctx, &keyvaluestorage.SetIfNotExistRequest{
		Namespace: namespaceName,
		Key:       key,
		Value:     []byte(value),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespaceName, Key: key})

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: namespaceName,
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	assert.Equal(s.T(), []byte(value), r.Value)
}

func (s *SetKeyIfNotExistTestSuite) TestFailsIfNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(20)
	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Keyvaluestorage.SetIfNotExist(ctx, &keyvaluestorage.SetIfNotExistRequest{
		Namespace: namespaceName,
		Key:       key,
		Value:     []byte(value),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespaceName, Key: key})
	require.NotNil(s.T(), err)
	if st, ok := status.FromError(err); ok {
		assert.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error expected")
	}
}

func (s *SetKeyIfNotExistTestSuite) TestDoesntUpdateValueIfAlreadyExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(20)
	value := tools.GetRandomString(20)

	s.nativeStub.Services.Keyvaluestorage.SetIfNotExist(ctx, &keyvaluestorage.SetIfNotExistRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})

	s.nativeStub.Services.Keyvaluestorage.SetIfNotExist(ctx, &keyvaluestorage.SetIfNotExistRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(tools.GetRandomString(20)),
	})

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	assert.Equal(s.T(), []byte(value), r.Value)
}
