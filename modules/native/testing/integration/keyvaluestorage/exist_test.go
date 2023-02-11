package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ExistTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ExistTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ExistTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestExistTestSuite(t *testing.T) {
	suite.Run(t, new(ExistTestSuite))
}

func (s *GetKeyTestSuite) TestExist() {
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
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: "",
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), r.Exist)
}

func (s *GetKeyTestSuite) TestExistInNamespace() {
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
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: namespaceName,
		Key:       key,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), r.Exist)
}

func (s *GetKeyTestSuite) DoesntExist() {
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
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: "",
		Key:       tools.GetRandomString(30),
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r.Exist)
}

func (s *GetKeyTestSuite) TestDoesntExistInNamespace() {
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
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: namespaceName,
		Key:       tools.GetRandomString(20),
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r.Exist)
}

func (s *GetKeyTestSuite) TestDoesntExistForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, err := s.nativeStub.Services.Keyvaluestorage.Exist(ctx, &keyvaluestorage.ExistRequest{
		Namespace: tools.GetRandomString(20),
		Key:       tools.GetRandomString(20),
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r.Exist)
}
