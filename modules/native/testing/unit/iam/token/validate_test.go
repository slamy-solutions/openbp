package token

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/testing/tools"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	token_jwt "github.com/slamy-solutions/openbp/modules/native/services/iam/token/src/jwt"
)

type ValidateTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

func (suite *ValidateTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMTokenService().WithIAMIdentityService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}

	suite.systemStub = system.NewSystemStub(system.NewSystemStubConfig().WithVault())
	err = suite.systemStub.Connect(ctx)
	if err != nil {
		panic(err)
	}
}
func (suite *ValidateTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestValidateTestSuite(t *testing.T) {
	suite.Run(t, new(ValidateTestSuite))
}

func (s *ValidateTestSuite) TestValidateTokenExpired() {
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

	tokenCreateResponse, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  "{}",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamToken.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	tokenValidateResponse, err := s.nativeStub.Services.IamToken.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, tokenValidateResponse.Status)

	tokenData := token_jwt.NewJWTData(tokenCreateResponse.TokenData.Uuid, "", identityCreateResponse.Identity.Uuid, []token_jwt.JWTScope{}, false, time.Now().UTC().Add(time.Second))
	jwtService := token_jwt.NewJWTService(s.systemStub)
	signedTokenString, err := jwtService.JWTDataToSignedString(ctx, tokenData)
	require.Nil(s.T(), err)

	time.Sleep(time.Second * 2)

	tokenValidateResponse, err = s.nativeStub.Services.IamToken.Validate(ctx, &token.ValidateRequest{
		Token:    signedTokenString,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_EXPIRED, tokenValidateResponse.Status)
}

func (s *ValidateTestSuite) TestValidateTokenWrongRSASignature() {
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

	tokenCreateResponse, err := s.nativeStub.Services.IamToken.Create(ctx, &token.CreateRequest{
		Namespace: "",
		Identity:  identityCreateResponse.Identity.Uuid,
		Scopes:    []*token.Scope{},
		Metadata:  "{}",
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IamToken.Delete(context.Background(), &token.DeleteRequest{Namespace: "", Uuid: tokenCreateResponse.TokenData.Uuid})

	tokenValidateResponse, err := s.nativeStub.Services.IamToken.Validate(ctx, &token.ValidateRequest{
		Token:    tokenCreateResponse.Token,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_OK, tokenValidateResponse.Status)

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(s.T(), err)

	tokenData := token_jwt.NewJWTData(tokenCreateResponse.TokenData.Uuid, "", identityCreateResponse.Identity.Uuid, []token_jwt.JWTScope{}, false, time.Now().UTC().Add(time.Minute))
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS512, tokenData)
	tokenString, err := jwtToken.SignedString(pk)
	require.Nil(s.T(), err)

	tokenValidateResponse, err = s.nativeStub.Services.IamToken.Validate(ctx, &token.ValidateRequest{
		Token:    tokenString,
		UseCache: true,
	})
	require.Nil(s.T(), err)
	require.Equal(s.T(), token.ValidateResponse_INVALID, tokenValidateResponse.Status)
}
