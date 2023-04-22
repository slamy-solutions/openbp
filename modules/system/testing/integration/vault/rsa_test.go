package rsa

import (
	"context"
	"testing"
	"time"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	tools "github.com/slamy-solutions/openbp/modules/system/testing/tools"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RSATestSuite struct {
	suite.Suite

	systemStub *system.SystemStub
}

func (suite *RSATestSuite) SetupSuite() {
	suite.systemStub = system.NewSystemStub(system.NewSystemStubConfig().WithVault())
	err := suite.systemStub.Connect(context.Background())
	if err != nil {
		panic(err)
	}
}
func (suite *RSATestSuite) TearDownSuite() {
	suite.systemStub.Close(context.Background())
}
func TestRSATestSuite(t *testing.T) {
	suite.Run(t, new(RSATestSuite))
}

func (s *RSATestSuite) TestSignAndVerify() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	rsaKeyName := tools.GetRandomString(20)
	_, err := s.systemStub.Vault.EnsureRSAKeyPair(ctx, &vault.EnsureRSAKeyPairRequest{
		KeyName: rsaKeyName,
	})
	require.Nil(s.T(), err)
	//TODO: delete RSA key after the test

	data := tools.GetRandomBytes(100000)

	signChannel, err := s.systemStub.Vault.RSASign(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 1379 {
		endByte := i + 1379
		if len(data) < endByte {
			endByte = len(data)
		}
		err := signChannel.Send(&vault.RSASignRequest{
			KeyName: rsaKeyName,
			Data:    data[i:endByte],
		})
		require.Nil(s.T(), err)
	}

	err = signChannel.CloseSend()
	require.Nil(s.T(), err)

	signResponse, err := signChannel.CloseAndRecv()
	require.Nil(s.T(), err)

	verifyChannel, err := s.systemStub.Vault.RSAVerify(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 1234 {
		endByte := i + 1234
		if len(data) < endByte {
			endByte = len(data)
		}
		err := verifyChannel.Send(&vault.RSAVerifyRequest{
			KeyName:   rsaKeyName,
			Data:      data[i:endByte],
			Signature: signResponse.Signature,
		})
		require.Nil(s.T(), err)
	}
	verifyResponse, err := verifyChannel.CloseAndRecv()
	require.Nil(s.T(), err)
	require.True(s.T(), verifyResponse.Valid)
}

func (s *RSATestSuite) TestVerifyBadSignature() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rsaKeyName := tools.GetRandomString(20)
	_, err := s.systemStub.Vault.EnsureRSAKeyPair(ctx, &vault.EnsureRSAKeyPairRequest{
		KeyName: rsaKeyName,
	})
	require.Nil(s.T(), err)

	data := tools.GetRandomBytes(100000)

	signChannel, err := s.systemStub.Vault.RSASign(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 7651 {
		endByte := i + 7651
		if len(data) < endByte {
			endByte = len(data)
		}
		err := signChannel.Send(&vault.RSASignRequest{
			KeyName: rsaKeyName,
			Data:    data[i:endByte],
		})
		require.Nil(s.T(), err)
	}
	signResponse, err := signChannel.CloseAndRecv()
	require.Nil(s.T(), err)

	badSignature := signResponse.Signature
	badSignature[0] = ^badSignature[0]

	verifyChannel, err := s.systemStub.Vault.RSAVerify(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 1234 {
		endByte := i + 1234
		if len(data) < endByte {
			endByte = len(data)
		}
		err := verifyChannel.Send(&vault.RSAVerifyRequest{
			KeyName:   rsaKeyName,
			Data:      data[i:endByte],
			Signature: badSignature,
		})
		require.Nil(s.T(), err)
	}
	verifyResponse, err := verifyChannel.CloseAndRecv()
	require.Nil(s.T(), err)
	require.False(s.T(), verifyResponse.Valid)
}
