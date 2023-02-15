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

type IndexPolicyTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *IndexPolicyTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService())
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
func (suite *IndexPolicyTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub.Close()
	suite.systemStub.Close(ctx)
}
func TestIndexPolicyTestSuite(t *testing.T) {
	suite.Run(t, new(IndexPolicyTestSuite))
}

func (s *IndexPolicyTestSuite) TestIndexCreationOnNewNamespace() {
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
		specs, err := s.systemStub.DB.Database("openbp_namespace_" + namespaceName).Collection("native_iam_policy").Indexes().ListSpecifications(ctx)
		require.Nil(s.T(), err)

		serviceFounded := false
		builtInFounded := false
		for _, index := range specs {
			serviceFounded = serviceFounded || (index.Name == "fast_search_service")
			builtInFounded = builtInFounded || (index.Name == "fast_search_built_in")
		}

		if serviceFounded && builtInFounded {
			return
		}

		time.Sleep(time.Second)
	}

	require.Fail(s.T(), "Index wasnt created")
}
