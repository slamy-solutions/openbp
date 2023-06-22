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

type GetTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetTestSuite(t *testing.T) {
	suite.Run(t, new(GetTestSuite))
}

func (s *GetTestSuite) TestGetsFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	metadata1 := tools.GetRandomString(20)
	metadata2 := tools.GetRandomString(20)

	tokenCreateResponse1, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  metadata1,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse1.TokenData.Uuid})

	tokenCreateResponse2, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  metadata2,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse2.TokenData.Uuid})

	tokenGetResponse1, err := s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse1.TokenData.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), "", tokenGetResponse1.TokenData.Namespace)
	require.Equal(s.T(), metadata1, tokenGetResponse1.TokenData.CreationMetadata)

	tokenGetResponse2, err := s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse1.TokenData.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), "", tokenGetResponse2.TokenData.Namespace)
	require.Equal(s.T(), metadata1, tokenGetResponse2.TokenData.CreationMetadata)
}

func (s *GetTestSuite) TestGetsFromNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	metadata1 := tools.GetRandomString(20)
	metadata2 := tools.GetRandomString(20)

	tokenCreateResponse1, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  metadata1,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse1.TokenData.Uuid})

	tokenCreateResponse2, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  metadata2,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse2.TokenData.Uuid})

	tokenGetResponse1, err := s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse1.TokenData.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), namespaceName, tokenGetResponse1.TokenData.Namespace)
	require.Equal(s.T(), metadata1, tokenGetResponse1.TokenData.CreationMetadata)

	tokenGetResponse2, err := s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse1.TokenData.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), namespaceName, tokenGetResponse2.TokenData.Namespace)
	require.Equal(s.T(), metadata1, tokenGetResponse2.TokenData.CreationMetadata)
}

func (s *GetTestSuite) TestFailsIfTokenDoesntExistForGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  "{}",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	_, err = s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: "",
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  true,
	})

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *GetTestSuite) TestFailsIfTokenDoesntExistForNamespace() {
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

	identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: identityCreateResponse.Identity.Uuid})

	tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  "{}",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Token.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse.TokenData.Uuid})

	_, err = s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: namespaceName,
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  false,
	})

	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}

func (s *GetTestSuite) TestFailsIfNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IAM.Token.Get(ctx, &token.GetRequest{
		Namespace: tools.GetRandomString(20),
		Uuid:      primitive.NewObjectID().Hex(),
		UseCache:  false,
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.NotFound, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
