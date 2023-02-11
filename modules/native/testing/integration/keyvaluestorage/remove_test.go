package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type RemoveKeyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *RemoveKeyTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *RemoveKeyTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestRemoveKeyTestSuite(t *testing.T) {
	suite.Run(t, new(RemoveKeyTestSuite))
}

func (s *RemoveKeyTestSuite) TestNotifiesIfValuesWasActuallyRemoved() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(30)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(tools.GetRandomString(30)),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})
	require.Nil(s.T(), err)

	r1, err := s.nativeStub.Services.Keyvaluestorage.Remove(ctx, &keyvaluestorage.RemoveRequest{
		Namespace: "",
		Key:       key,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), r1.Removed)

	r2, err := s.nativeStub.Services.Keyvaluestorage.Remove(ctx, &keyvaluestorage.RemoveRequest{
		Namespace: "",
		Key:       key,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r2.Removed)
}

func (s *RemoveKeyTestSuite) TestRemoveFromNonExistingNamespaceIsOk() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err := s.nativeStub.Services.Keyvaluestorage.Remove(ctx, &keyvaluestorage.RemoveRequest{
		Namespace: tools.GetRandomString(20),
		Key:       tools.GetRandomString(20),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r.Removed)
}

func (s *RemoveKeyTestSuite) TestDoesntExistAfterRemove() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	key := tools.GetRandomString(30)
	_, err := s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       key,
		Value:     []byte(tools.GetRandomString(30)),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: "", Key: key})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.Keyvaluestorage.Remove(ctx, &keyvaluestorage.RemoveRequest{
		Namespace: "",
		Key:       key,
	})
	require.Nil(s.T(), err)
	r, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: "",
		Key:       key,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	assert.False(s.T(), r.Exist)
}

func (s *RemoveKeyTestSuite) TestDoesntExistInNamespaceAfterRemove() {
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
	_, err = s.nativeStub.Services.Keyvaluestorage.Set(ctx, &keyvaluestorage.SetRequest{
		Namespace: namespaceName,
		Key:       key,
		Value:     []byte(tools.GetRandomString(30)),
	})
	defer s.nativeStub.Services.Keyvaluestorage.Remove(context.Background(), &keyvaluestorage.RemoveRequest{Namespace: namespaceName, Key: key})
	require.Nil(s.T(), err)

	r1, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: namespaceName,
		Key:       key,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	assert.False(s.T(), r1.Exist)

	_, err = s.nativeStub.Services.Keyvaluestorage.Remove(ctx, &keyvaluestorage.RemoveRequest{
		Namespace: namespaceName,
		Key:       key,
	})
	require.Nil(s.T(), err)
	r2, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: namespaceName,
		Key:       key,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	assert.False(s.T(), r2.Exist)
}
