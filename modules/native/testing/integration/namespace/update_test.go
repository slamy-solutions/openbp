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

type UpdateNamespaceTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *UpdateNamespaceTestSuite) SetupSuite() {
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
func (suite *UpdateNamespaceTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	suite.nativeStub.Close()
	suite.systemStub.Close(ctx)
}
func TestUpdateNamespaceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateNamespaceTestSuite))
}

func (s *UpdateNamespaceTestSuite) TestParamsValidation() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	testCases := []struct {
		testName    string
		fullName    string
		description string
		ok          bool
	}{
		{"Empty fullName is ok", "", "123", true},
		{"Max fullName length is 128", tools.GetRandomString(130), "123", false},
		{"Empty description is ok", "alalal", "", true},
		{"Max description length is 512", tools.GetRandomString(10), tools.GetRandomString(530), false},
	}

	for _, tc := range testCases {
		s.Run(tc.testName, func() {
			name := tools.GetRandomString(20)
			_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
				Name:        name,
				FullName:    tools.GetRandomString(10),
				Description: tools.GetRandomString(10),
			})
			defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: name})
			require.Nil(s.T(), err)

			_, err = s.nativeStub.Services.Namespace.Update(ctx, &namespace.UpdateNamespaceRequest{
				Name:        name,
				FullName:    tc.fullName,
				Description: tc.description,
			})

			assert.Equal(s.T(), tc.ok, err == nil)
			if err != nil {
				e, ok := status.FromError(err)
				require.True(s.T(), ok)
				require.Equal(s.T(), e.Code(), codes.InvalidArgument)
			}
		})
	}
}

func (s *UpdateNamespaceTestSuite) TestChangesDataOnUpdate() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	name := tools.GetRandomString(20)
	fullName1 := tools.GetRandomString(30)
	description1 := tools.GetRandomString(30)

	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
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

	_, err = s.nativeStub.Services.Namespace.Update(ctx, &namespace.UpdateNamespaceRequest{Name: name, FullName: fullName2, Description: description2})
	require.Nil(s.T(), err)

	r, err := s.nativeStub.Services.Namespace.Get(ctx, &namespace.GetNamespaceRequest{Name: name, UseCache: false})
	require.Nil(s.T(), err)

	assert.Equal(s.T(), r.Namespace.Name, name)
	assert.Equal(s.T(), r.Namespace.FullName, fullName2)
	assert.Equal(s.T(), r.Namespace.Description, description2)
}

func (s *UpdateNamespaceTestSuite) TestRaisesEventOnUpdate() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	js, _ := s.systemStub.Nats.JetStream()
	sub, err := js.PullSubscribe("native.namespace.event.updated", "", nats.BindStream("native_namespace_event"))
	require.Nil(s.T(), err)
	defer sub.Unsubscribe()

	namespaceName := tools.GetRandomString(20)
	_, err = s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{Name: namespaceName, FullName: "", Description: ""})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	_, err = s.nativeStub.Services.Namespace.Update(ctx, &namespace.UpdateNamespaceRequest{Name: namespaceName, FullName: "updated", Description: "updated"})
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

func (s *UpdateNamespaceTestSuite) TestFailsWhileUpdatingOnNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := s.nativeStub.Services.Namespace.Update(ctx, &namespace.UpdateNamespaceRequest{Name: tools.GetRandomString(20), FullName: "updated", Description: "updated"})
	require.NotNil(s.T(), err)
	e, ok := status.FromError(err)
	require.True(s.T(), ok)
	require.Equal(s.T(), e.Code(), codes.NotFound)
}
