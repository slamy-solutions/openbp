package namespace

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
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

type CreateNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *CreateNamespaceTestSuite) SetupSuite() {
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
func (suite *CreateNamespaceTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub.Close()
	suite.systemStub.Close(ctx)
}
func TestCreateNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(CreateNamespaceTestSuite))
}

func (s *CreateNamespaceTestSuite) TestParamsValidation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for _, tc := range createParamsValidationTestCases {
		s.Run(tc.testName, func() {
			_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
				Name:        tc.name,
				FullName:    tc.fullName,
				Description: tc.description,
			})
			defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: tc.name})

			assert.Equal(s.T(), tc.ok, err == nil)
			if err != nil {
				e, ok := status.FromError(err)
				require.True(s.T(), ok)
				require.Equal(s.T(), e.Code(), codes.InvalidArgument)
			}
		})
	}
}

func (s *CreateNamespaceTestSuite) TestNamespaceAvailableAfterCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	fullName := tools.GetRandomString(30)
	description := tools.GetRandomString(30)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    fullName,
		Description: description,
	})
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})

	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: name})
	require.Nil(s.T(), err)

	assert.Equal(s.T(), r.Namespace.Name, name)
	assert.Equal(s.T(), r.Namespace.FullName, fullName)
	assert.Equal(s.T(), r.Namespace.Description, description)
}

func (s *CreateNamespaceTestSuite) TestRaisesEventOnCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	js, _ := s.systemStub.Nats.JetStream()
	sub, err := js.PullSubscribe("native.namespace.event.created", "", nats.BindStream("native_namespace_event"))
	require.Nil(s.T(), err)
	defer sub.Unsubscribe()

	namespaceName := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{Name: namespaceName, FullName: "", Description: ""})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

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

func (s *CreateNamespaceTestSuite) TestCantCreateWithSameName() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)

	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	require.NotNil(s.T(), err)
	e, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), e.Code(), codes.AlreadyExists)
}
