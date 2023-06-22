package token

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
)

type GetTokensForIdentityTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetTokensForIdentityTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetTokensForIdentityTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetTokensForIdentityTestSuite(t *testing.T) {
	suite.Run(t, new(GetTokensForIdentityTestSuite))
}

func (s *GetTokensForIdentityTestSuite) TestGetsFromGlobalNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	identities := []string{}
	defer func() {
		for _, uuid := range identities {
			s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: "", Uuid: uuid})
		}
	}()
	for i := 0; i < 5; i++ {
		identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
			Namespace:       "",
			Name:            tools.GetRandomString(20),
			InitiallyActive: true,
			Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		identities = append(identities, identityCreateResponse.Identity.Uuid)
	}

	tokens := []string{}
	defer func() {
		for _, uuid := range tokens {
			s.nativeStub.Services.IAM.Token.Delete(ctx, &token.DeleteRequest{Namespace: "", Uuid: uuid})
		}
	}()
	for _, identityUUID := range identities {
		for i := 0; i < 5; i++ {
			tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
				Namespace: "",
				Identity:  identityUUID,
				Scopes:    []*token.Scope{},
				Metadata:  "{}",
			})
			require.Nil(s.T(), err)
			tokens = append(tokens, tokenCreateResponse.TokenData.Uuid)
		}
	}

	tokensForIdentityResponse, err := s.nativeStub.Services.IAM.Token.GetTokensForIdentity(ctx, &token.GetTokensForIdentityRequest{
		Namespace:    "",
		Identity:     identities[1],
		ActiveFilter: token.GetTokensForIdentityRequest_ALL,
		Skip:         0,
		Limit:        0,
	})
	require.Nil(s.T(), err)

	receivedTokenUUIDs := []string{}
	for {
		chunk, err := tokensForIdentityResponse.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedTokenUUIDs = append(receivedTokenUUIDs, chunk.TokenData.Uuid)
	}
	require.ElementsMatch(s.T(), tokens[5:10], receivedTokenUUIDs)

	tokensForIdentityResponse, err = s.nativeStub.Services.IAM.Token.GetTokensForIdentity(ctx, &token.GetTokensForIdentityRequest{
		Namespace:    "",
		Identity:     identities[1],
		ActiveFilter: token.GetTokensForIdentityRequest_ALL,
		Skip:         1,
		Limit:        2,
	})
	require.Nil(s.T(), err)

	receivedTokenUUIDs = []string{}
	for {
		chunk, err := tokensForIdentityResponse.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedTokenUUIDs = append(receivedTokenUUIDs, chunk.TokenData.Uuid)
	}
	require.ElementsMatch(s.T(), tokens[7:9], receivedTokenUUIDs)
}

func (s *GetTokensForIdentityTestSuite) TestGetsFromNamespace() {
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

	identities := []string{}
	defer func() {
		for _, uuid := range identities {
			s.nativeStub.Services.IAM.Identity.Delete(ctx, &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: uuid})
		}
	}()
	for i := 0; i < 5; i++ {
		identityCreateResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
			Namespace:       namespaceName,
			Name:            tools.GetRandomString(20),
			InitiallyActive: true,
			Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
		})
		require.Nil(s.T(), err)
		identities = append(identities, identityCreateResponse.Identity.Uuid)
	}

	tokens := []string{}
	defer func() {
		for _, uuid := range tokens {
			s.nativeStub.Services.IAM.Token.Delete(ctx, &token.DeleteRequest{Namespace: namespaceName, Uuid: uuid})
		}
	}()
	for _, identityUUID := range identities {
		for i := 0; i < 5; i++ {
			tokenCreateResponse, err := s.nativeStub.Services.IAM.Token.Create(ctx, &token.CreateRequest{
				Namespace: namespaceName,
				Identity:  identityUUID,
				Scopes:    []*token.Scope{},
				Metadata:  "{}",
			})
			require.Nil(s.T(), err)
			tokens = append(tokens, tokenCreateResponse.TokenData.Uuid)
		}
	}

	tokensForIdentityResponse, err := s.nativeStub.Services.IAM.Token.GetTokensForIdentity(ctx, &token.GetTokensForIdentityRequest{
		Namespace:    namespaceName,
		Identity:     identities[1],
		ActiveFilter: token.GetTokensForIdentityRequest_ALL,
		Skip:         0,
		Limit:        0,
	})
	require.Nil(s.T(), err)

	receivedTokenUUIDs := []string{}
	for {
		chunk, err := tokensForIdentityResponse.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedTokenUUIDs = append(receivedTokenUUIDs, chunk.TokenData.Uuid)
	}
	require.ElementsMatch(s.T(), tokens[5:10], receivedTokenUUIDs)

	tokensForIdentityResponse, err = s.nativeStub.Services.IAM.Token.GetTokensForIdentity(ctx, &token.GetTokensForIdentityRequest{
		Namespace:    namespaceName,
		Identity:     identities[1],
		ActiveFilter: token.GetTokensForIdentityRequest_ALL,
		Skip:         1,
		Limit:        2,
	})
	require.Nil(s.T(), err)

	receivedTokenUUIDs = []string{}
	for {
		chunk, err := tokensForIdentityResponse.Recv()
		if err != nil {
			require.Equal(s.T(), io.EOF, err)
			break
		}
		receivedTokenUUIDs = append(receivedTokenUUIDs, chunk.TokenData.Uuid)
	}
	require.ElementsMatch(s.T(), tokens[7:9], receivedTokenUUIDs)
}

func (s *GetTokensForIdentityTestSuite) TestGetForNonExistingNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.nativeStub.Services.IAM.Token.GetTokensForIdentity(ctx, &token.GetTokensForIdentityRequest{
		Namespace:    tools.GetRandomString(20),
		Identity:     primitive.NewObjectID().Hex(),
		ActiveFilter: token.GetTokensForIdentityRequest_ALL,
		Skip:         0,
		Limit:        0,
	})
	require.Nil(s.T(), err)
}
