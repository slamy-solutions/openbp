package identity

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type ActiveTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ActiveTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ActiveTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestActiveTestSuite(t *testing.T) {
	suite.Run(t, new(ActiveTestSuite))
}

func (s *ActiveTestSuite) TestInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		Active:    false,
	})
	require.Nil(s.T(), err)
	getResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), getResponse.Identity.Active)

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		Active:    true,
	})
	require.Nil(s.T(), err)
	getResponse, err = s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), getResponse.Identity.Active)

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		Active:    false,
	})
	require.Nil(s.T(), err)
	getResponse, err = s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: "",
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), getResponse.Identity.Active)
}

func (s *ActiveTestSuite) TestInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		Active:    false,
	})
	require.Nil(s.T(), err)
	getResponse, err := s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), getResponse.Identity.Active)

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		Active:    true,
	})
	require.Nil(s.T(), err)
	getResponse, err = s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), getResponse.Identity.Active)

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		Active:    false,
	})
	require.Nil(s.T(), err)
	getResponse, err = s.nativeStub.Services.IamIdentity.Get(ctx, &identity.GetIdentityRequest{
		Namespace: namespaceName,
		Uuid:      createResponse.Identity.Uuid,
		UseCache:  false,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), getResponse.Identity.Active)
}

func (s *ActiveTestSuite) TestFailsWithNotFoundErrorWhenIdentityDoesntExistInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
		Active:    false,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ActiveTestSuite) TestFailsWithNotFoundErrorWhenIdentityDoesntExistInNamespace() {
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

	createResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		InitiallyActive: true,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createResponse.Identity.Uuid})

	_, err = s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
		Active:    false,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *ActiveTestSuite) TestFailsWithNotFoundErrorWhenNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IamIdentity.SetActive(ctx, &identity.SetIdentityActiveRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
		Active:    false,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
