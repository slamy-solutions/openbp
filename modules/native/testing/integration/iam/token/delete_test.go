package token

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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type DeleteTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *DeleteTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMTokenService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *DeleteTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}

func (s *DeleteTestSuite) TestDeleteFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  "{}",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamToken.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	deleteResponse, err := s.nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), deleteResponse.Existed)

	_, err = s.nativeStub.Services.IamToken.Get(ctx, &token.GetRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
		UseCache:  false,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *DeleteTestSuite) TestDeleteFromNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  "{}",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamToken.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse.TokenData.Uuid})

	deleteResponse, err := s.nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse.TokenData.Uuid,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), deleteResponse.Existed)

	_, err = s.nativeStub.Services.IamToken.Get(ctx, &token.GetRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse.TokenData.Uuid,
		UseCache:  false,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *DeleteTestSuite) TestDeleteNonExistingTokenFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	deleteResponse, err := s.nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingTokenFromNamespace() {
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

	deleteResponse, err := s.nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}

func (s *DeleteTestSuite) TestDeleteNonExistingTokenFromNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	deleteResponse, err := s.nativeStub.Services.IamToken.Delete(ctx, &token.DeleteRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
	})
	require.Nil(s.T(), err)
	require.False(s.T(), deleteResponse.Existed)
}
