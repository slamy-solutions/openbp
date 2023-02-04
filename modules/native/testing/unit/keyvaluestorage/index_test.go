package keyvaluestorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type IndexKeyValueStorageTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *IndexKeyValueStorageTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithKeyValueStorageService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}

	suite.systemStub = system.NewSystemStub(system.NewSystemStubConfig().WithDB())
	err = suite.systemStub.Connect(ctx)
	if err != nil {
		panic(err)
	}
}
func (suite *IndexKeyValueStorageTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub.Close()
	suite.systemStub.Close(ctx)
}
func TestIndexKeyValueStorageTestSuite(t *testing.T) {
	suite.Run(t, new(IndexKeyValueStorageTestSuite))
}

func (s *IndexKeyValueStorageTestSuite) TestIndexCreationOnNewNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Create namespace
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

	waitStartTimestamp := time.Now().UTC().Unix()
	for time.Now().UTC().Unix()-waitStartTimestamp < 20 {
		specs, err := s.systemStub.DB.Database("openbp_namespace_" + namespaceName).Collection("native_keyvaluestorage").Indexes().ListSpecifications(ctx)
		require.Nil(s.T(), err)

		for _, index := range specs {
			if index.Name == "key_hashed" { // Search for specific index name
				return
			}
		}

		time.Sleep(time.Second)
	}

	require.Fail(s.T(), "Index wasnt created")
}
