package token

import (
	"context"
	"reflect"
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

type CreateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *CreateTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMTokenService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *CreateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestCreateTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTestSuite))
}

func (s *CreateTestSuite) TestCreatesInGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identityCreateResponse, err := s.nativeStub.Services.IamIdentity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamIdentity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: identityCreateResponse.Identity.Uuid})

	metadata := tools.GetRandomString(20)
	scopes := []*token.Scope{
		{
			Namespace: "",
			Resources: []string{tools.GetRandomString(20), tools.GetRandomString(20)},
			Actions:   []string{tools.GetRandomString(20)},
		},
		{
			Namespace: "",
			Resources: []string{tools.GetRandomString(20)},
			Actions:   []string{tools.GetRandomString(20), tools.GetRandomString(20)},
		},
		{
			Namespace: tools.GetRandomString(20),
			Resources: []string{tools.GetRandomString(20)},
			Actions:   []string{tools.GetRandomString(20)},
		},
	}

	tokenCreateResponse, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    scopes,
		Metadata:  metadata,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamToken.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	require.Equal(s.T(), "", tokenCreateResponse.TokenData.Namespace)
	require.Equal(s.T(), metadata, tokenCreateResponse.TokenData.CreationMetadata)

	for _, scope := range scopes {
		founded := false
		for _, receivedScope := range tokenCreateResponse.TokenData.Scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}

	for _, receivedScope := range tokenCreateResponse.TokenData.Scopes {
		founded := false
		for _, scope := range scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}

	tokenGetResponse, err := s.nativeStub.Services.IamToken.Get(ctx, &token.GetRequest{
		Namespace: "",
		Uuid:      tokenCreateResponse.TokenData.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), "", tokenGetResponse.TokenData.Namespace)
	require.Equal(s.T(), metadata, tokenGetResponse.TokenData.CreationMetadata)

	for _, scope := range scopes {
		founded := false
		for _, receivedScope := range tokenGetResponse.TokenData.Scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}

	for _, receivedScope := range tokenGetResponse.TokenData.Scopes {
		founded := false
		for _, scope := range scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}
}

func (s *CreateTestSuite) TestCreatesInNamespace() {
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

	metadata := tools.GetRandomString(20)
	scopes := []*token.Scope{
		{
			Namespace: "",
			Resources: []string{tools.GetRandomString(20), tools.GetRandomString(20)},
			Actions:   []string{tools.GetRandomString(20)},
		},
		{
			Namespace: "",
			Resources: []string{tools.GetRandomString(20)},
			Actions:   []string{tools.GetRandomString(20), tools.GetRandomString(20)},
		},
		{
			Namespace: tools.GetRandomString(20),
			Resources: []string{tools.GetRandomString(20)},
			Actions:   []string{tools.GetRandomString(20)},
		},
	}

	tokenCreateResponse, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: namespaceName,
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    scopes,
		Metadata:  metadata,
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamToken.Delete(context.Background(), &token.DeleteRequest{Namespace: namespaceName, Uuid: tokenCreateResponse.TokenData.Uuid})

	require.Equal(s.T(), namespaceName, tokenCreateResponse.TokenData.Namespace)
	require.Equal(s.T(), metadata, tokenCreateResponse.TokenData.CreationMetadata)

	for _, scope := range scopes {
		founded := false
		for _, receivedScope := range tokenCreateResponse.TokenData.Scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}

	for _, receivedScope := range tokenCreateResponse.TokenData.Scopes {
		founded := false
		for _, scope := range scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}

	tokenGetResponse, err := s.nativeStub.Services.IamToken.Get(ctx, &token.GetRequest{
		Namespace: namespaceName,
		Uuid:      tokenCreateResponse.TokenData.Uuid,
		UseCache:  true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), namespaceName, tokenGetResponse.TokenData.Namespace)
	require.Equal(s.T(), metadata, tokenGetResponse.TokenData.CreationMetadata)

	for _, scope := range scopes {
		founded := false
		for _, receivedScope := range tokenGetResponse.TokenData.Scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}

	for _, receivedScope := range tokenGetResponse.TokenData.Scopes {
		founded := false
		for _, scope := range scopes {
			if receivedScope.Namespace != scope.Namespace {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Resources, scope.Resources) {
				continue
			}
			if !reflect.DeepEqual(receivedScope.Actions, scope.Actions) {
				continue
			}

			founded = true
			break
		}
		require.True(s.T(), founded)
	}
}

func (s *CreateTestSuite) TestFailsToCreateWhenNamespaceDoesntExist() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: tools.GetRandomString(20),
		Identity:  primitive.NewObjectID().Hex(),
		Scopes:    []*token.Scope{},
		Metadata:  tools.GetRandomString(20),
	})
	if st, ok := status.FromError(err); ok {
		require.Equal(s.T(), codes.FailedPrecondition, st.Code())
	} else {
		require.Fail(s.T(), "Error wasnt received")
	}
}
