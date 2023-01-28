package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

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

	suite.systemStub = system.NewSystemStub(system.NewSystemStubConfig().WithNats())
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

func (s *DeleteNamespaceTestSuite) TestNamespaceUnavailableAfterDeletion() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	require.Nil(s.T(), err)
	_, err = s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: tools.GetRandomString(20)})
	require.NotNil(s.T(), err)
	e, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), e.Code(), codes.NotFound)
}

func (s *DeleteNamespaceTestSuite) TesIndicatesIfExistedBeforeDeletion() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	require.Nil(s.T(), err)

	r1, err := s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), r1.Existed)
	r2, err := s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r2.Existed)
}

func (s *DeleteNamespaceTestSuite) TestRaisesEventOnDeletion() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	js, _ := s.systemStub.Nats.JetStream()
	sub, err := js.PullSubscribe("native.namespace.event.deleted", "", nats.BindStream("native_namespace_event"))
	require.Nil(s.T(), err)
	defer sub.Unsubscribe()

	namespaceName := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{Name: namespaceName, FullName: "", Description: ""})
	require.Nil(s.T(), err)
	_, err = s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})
	require.Nil(s.T(), err)

	waitStarted := time.Now().Unix()

	for time.Now().Unix()-waitStarted < 10 {
		events, err := sub.Fetch(1024)
		require.Nil(s.T(), err)

		for _, event := range events {
			var nm namespace.Namespace
			err = proto.Unmarshal(event.Data, &nm)
			require.Nil(s.T(), err)
			if nm.Name == namespaceName {
				return
			}
		}
	}

	require.Fail(s.T(), "Event wasnt received")
}
