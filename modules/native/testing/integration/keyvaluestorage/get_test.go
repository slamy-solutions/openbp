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
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type GetKeyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetKeyTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetKeyTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetKeyTestSuite(t *testing.T) {
	suite.Run(t, new(GetKeyTestSuite))
}

func (s *GetKeyTestSuite) TestGetsValue() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(30)
	value := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(value),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})
	assert.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), []byte(value), r.Value)
}

func (s *GetKeyTestSuite) TestGetsValueInNamespace() {
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

	key := tools.GetRandomString(30)
	value := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: namespaceName,
		Key:       key,
		Value:     []byte(value),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespaceName, Key: key})
	assert.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: namespaceName,
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), []byte(value), r.Value)
}

func (s *GetKeyTestSuite) TestFailsForNonExistingValue() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	existingKey := tools.GetRandomString(20)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       existingKey,
		Value:     []byte(tools.GetRandomString(20)),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: existingKey})

	_, err = s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       tools.GetRandomString(20),
		UseCache:  true,
	})
	require.NotNil(s.T(), err)

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *GetKeyTestSuite) TestFailsForNonExistingValueInNamespace() {
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

	existingKey := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: namespaceName,
		Key:       existingKey,
		Value:     []byte(tools.GetRandomString(20)),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespaceName, Key: existingKey})

	_, err = s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: namespaceName,
		Key:       tools.GetRandomString(20),
		UseCache:  true,
	})
	require.NotNil(s.T(), err)

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *GetKeyTestSuite) TestFailsForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.Keyvaluestorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: tools.GetRandomString(20),
		Key:       tools.GetRandomString(20),
		UseCache:  true,
	})
	require.NotNil(s.T(), err)

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
