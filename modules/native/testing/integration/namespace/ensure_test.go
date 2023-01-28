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

type EnsureNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *EnsureNamespaceTestSuite) SetupSuite() {
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
func (suite *EnsureNamespaceTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub.Close()
	suite.systemStub.Close(ctx)
}
func TestEnsureNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(EnsureNamespaceTestSuite))
}

func (s *EnsureNamespaceTestSuite) TestParamsValidation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	for _, tc := range createParamsValidationTestCases {
		s.Run(tc.testName, func() {
			_, err := s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
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

func (s *EnsureNamespaceTestSuite) TestNamespaceAvailableAfterCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	fullName := tools.GetRandomString(30)
	description := tools.GetRandomString(30)

	_, err := s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
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

func (s *EnsureNamespaceTestSuite) TestRaisesEventOnCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	js, _ := s.systemStub.Nats.JetStream()
	sub, err := js.PullSubscribe("native.namespace.event.created", "", nats.BindStream("native_namespace_event"))
	require.Nil(s.T(), err)
	defer sub.Unsubscribe()

	namespaceName := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{Name: namespaceName, FullName: "", Description: ""})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	waitStarted := time.Now().Unix()

	for time.Now().Unix()-waitStarted < 10 {
		events, err := sub.Fetch(1024)
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			} else {
				require.Fail(s.T(), "Error while fetching events from nats")
			}
		}

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

func (s *EnsureNamespaceTestSuite) TestIndicatesWhenCreatedAndWhenNot() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)

	r1, err := s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), r1.Created)

	r2, err := s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
		Name:        name,
		FullName:    tools.GetRandomString(30),
		Description: tools.GetRandomString(30),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), r2.Created)
}

func (s *EnsureNamespaceTestSuite) TestDoesntOverrideDataWhileExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	fullName1 := tools.GetRandomString(30)
	description1 := tools.GetRandomString(30)

	_, err := s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
		Name:        name,
		FullName:    fullName1,
		Description: description1,
	})
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{
		Name: name,
	})
	require.Nil(s.T(), err)

	fullName2 := tools.GetRandomString(30)
	description2 := tools.GetRandomString(30)

	_, err = s.nativeStub.Services.Namespace.Ensure(ctx, &namespace.EnsureNamespaceRequest{
		Name:        name,
		FullName:    fullName2,
		Description: description2,
	})
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: name})
	require.Nil(s.T(), err)

	assert.Equal(s.T(), r.Namespace.Name, name)
	assert.Equal(s.T(), r.Namespace.FullName, fullName1)
	assert.Equal(s.T(), r.Namespace.Description, description1)
}
