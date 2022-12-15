package user

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	tests "github.com/slamy-solutions/openbp/modules/native/testing"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	nativeActorUserGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
)

type ExampleTestSuite struct {
	suite.Suite

	nativeActorUserClient nativeActorUserGRPC.ActorUserServiceClient

	grpcConnections []*grpc.ClientConn
}

func (suite *ExampleTestSuite) SetupSuite() {
	connection, client, err := native.NewActorUserConnection(tests.NATIVE_ACTOR_USER_URL)
	if err != nil {
		panic("Failed to connect to the native_actor_user service: " + err.Error())
	}
	suite.nativeActorUserClient = client
	suite.grpcConnections = append(suite.grpcConnections, connection)
}
func (suite *ExampleTestSuite) TearDownSuite() {
	for _, con := range suite.grpcConnections {
		con.Close()
	}
}

func (suite *ExampleTestSuite) TestSearchLoginMatch() {
	// Create 10 users
	const USERS_COUNT = 10
	createdUsers := make([]*nativeActorUserGRPC.User, USERS_COUNT)
	for i := 0; i < USERS_COUNT; i++ {
		createResponse, err := suite.nativeActorUserClient.Create(context.Background(), &nativeActorUserGRPC.CreateRequest{
			Login:    tests.RandomString(20),
			FullName: tests.RandomString(20),
			Avatar:   "",
			Email:    tests.RandomString(120) + "@" + tests.RandomString(3) + ".com",
		})
		require.Nil(suite.T(), err, "Error on creating users")
		createdUsers = append(createdUsers, createResponse.User)
		defer suite.nativeActorUserClient.Delete(context.Background(), &nativeActorUserGRPC.DeleteRequest{Uuid: createResponse.User.Uuid})
	}

	// Search my matched user
	searchedUser := createdUsers[len(createdUsers)/2]
	search, err := suite.nativeActorUserClient.Search(context.Background(), &nativeActorUserGRPC.SearchRequest{
		Limit: 1,
		Match: searchedUser.Login,
	})
	defer search.CloseSend()
	require.Nil(suite.T(), err, "Error on searching users")
	searchedItem, err := search.Recv()
	require.Nil(suite.T(), err, "Error on receiving response")
	require.Equal(suite.T(), searchedItem.User.Uuid, searchedUser.Uuid)

	_, err = search.Recv()
	require.ErrorIs(suite.T(), err, io.EOF)
}

func (suite *ExampleTestSuite) TestSearchFullNameMatch() {
	// Create 10 users
	const USERS_COUNT = 10
	createdUsers := make([]*nativeActorUserGRPC.User, USERS_COUNT)
	for i := 0; i < USERS_COUNT; i++ {
		createResponse, err := suite.nativeActorUserClient.Create(context.Background(), &nativeActorUserGRPC.CreateRequest{
			Login:    tests.RandomString(20),
			FullName: tests.RandomString(20),
			Avatar:   "",
			Email:    tests.RandomString(120) + "@" + tests.RandomString(3) + ".com",
		})
		require.Nil(suite.T(), err, "Error on creating users")
		createdUsers = append(createdUsers, createResponse.User)
		defer suite.nativeActorUserClient.Delete(context.Background(), &nativeActorUserGRPC.DeleteRequest{Uuid: createResponse.User.Uuid})
	}

	// Search my matched user
	searchedUser := createdUsers[len(createdUsers)/2]
	search, err := suite.nativeActorUserClient.Search(context.Background(), &nativeActorUserGRPC.SearchRequest{
		Limit: 1,
		Match: searchedUser.FullName,
	})
	defer search.CloseSend()
	require.Nil(suite.T(), err, "Error on searching users")
	searchedItem, err := search.Recv()
	require.Nil(suite.T(), err, "Error on receiving response")
	require.Equal(suite.T(), searchedItem.User.Uuid, searchedUser.Uuid)

	_, err = search.Recv()
	require.ErrorIs(suite.T(), err, io.EOF)
}

func (suite *ExampleTestSuite) TestSearchEmailMatch() {
	// Create 10 users
	const USERS_COUNT = 10
	createdUsers := make([]*nativeActorUserGRPC.User, USERS_COUNT)
	for i := 0; i < USERS_COUNT; i++ {
		createResponse, err := suite.nativeActorUserClient.Create(context.Background(), &nativeActorUserGRPC.CreateRequest{
			Login:    tests.RandomString(20),
			FullName: tests.RandomString(20),
			Avatar:   "",
			Email:    tests.RandomString(120) + "@" + tests.RandomString(3) + ".com",
		})
		require.Nil(suite.T(), err, "Error on creating users")
		createdUsers = append(createdUsers, createResponse.User)
		defer suite.nativeActorUserClient.Delete(context.Background(), &nativeActorUserGRPC.DeleteRequest{Uuid: createResponse.User.Uuid})
	}

	// Search my matched user
	searchedUser := createdUsers[len(createdUsers)/2]
	search, err := suite.nativeActorUserClient.Search(context.Background(), &nativeActorUserGRPC.SearchRequest{
		Limit: 1,
		Match: searchedUser.Email,
	})
	require.Nil(suite.T(), err, "Error on searching users")
	defer search.CloseSend()
	searchedItem, err := search.Recv()
	require.Nil(suite.T(), err, "Error on receiving response")
	require.Equal(suite.T(), searchedItem.User.Uuid, searchedUser.Uuid)

	_, err = search.Recv()
	require.ErrorIs(suite.T(), err, io.EOF)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
