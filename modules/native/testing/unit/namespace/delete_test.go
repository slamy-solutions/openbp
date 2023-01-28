package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	tools "github.com/slamy-solutions/openbp/modules/native/testing/tools"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type DeleteNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *DeleteNamespaceTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
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
func (suite *DeleteNamespaceTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub.Close()
	suite.systemStub.Close(ctx)
}
func TestDeleteNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteNamespaceTestSuite))
}

func (s *DeleteNamespaceTestSuite) TestDeletesDBOnNamespaceDeletetion() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	name := tools.GetRandomString(20)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	require.Nil(s.T(), err)

	// Create some collection in namespace
	collectionName := "testcoll_" + tools.GetRandomString(20)
	s.systemStub.DB.Database("openbp_namespace_"+name).Collection(collectionName).InsertOne(ctx, bson.M{})

	dbs, err := s.systemStub.DB.ListDatabaseNames(ctx, bson.M{})
	require.Nil(s.T(), err)
	require.Contains(s.T(), dbs, "openbp_namespace_"+name)

	_, err = s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)

	dbs, err = s.systemStub.DB.ListDatabaseNames(ctx, bson.M{})
	require.Nil(s.T(), err)

	require.NotContains(s.T(), dbs, "openbp_namespace_"+name)
}
