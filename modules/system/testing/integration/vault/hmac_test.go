package vault

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

type HMACTestSuite struct {
	suite.Suite

	systemStub *system.SystemStub
}

func (suite *HMACTestSuite) SetupSuite() {
	suite.systemStub = system.NewSystemStub(system.NewSystemStubConfig().WithVault())
	err := suite.systemStub.Connect(context.Background())
	if err != nil {
		panic(err)
	}
}
func (suite *HMACTestSuite) TearDownSuite() {
	suite.systemStub.Close(context.Background())
}
func TestHMACTestSuite(t *testing.T) {
	suite.Run(t, new(HMACTestSuite))
}

func (s *HMACTestSuite) TestSignAndVerifyStream() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	data := tools.GetRandomBytes(100000)

	signChannel, err := s.systemStub.Vault.HMACSignStream(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 1379 {
		endByte := i + 1379
		if len(data) < endByte {
			endByte = len(data)
		}
		err := signChannel.Send(&vault.HMACSignStreamRequest{
			Data: data[i:endByte],
		})
		require.Nil(s.T(), err)
	}

	signResponse, err := signChannel.CloseAndRecv()
	require.Nil(s.T(), err)

	verifyChannel, err := s.systemStub.Vault.HMACVerifyStream(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 1234 {
		endByte := i + 1234
		if len(data) < endByte {
			endByte = len(data)
		}
		err := verifyChannel.Send(&vault.HMACVerifyStreamRequest{
			Data:      data[i:endByte],
			Signature: signResponse.Signature,
		})
		require.Nil(s.T(), err)
	}
	verifyResponse, err := verifyChannel.CloseAndRecv()
	require.Nil(s.T(), err)
	require.True(s.T(), verifyResponse.Valid)
}

func (s *HMACTestSuite) TestVerifyStreamBadSignature() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data := tools.GetRandomBytes(100000)

	signChannel, err := s.systemStub.Vault.HMACSignStream(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 7651 {
		endByte := i + 7651
		if len(data) < endByte {
			endByte = len(data)
		}
		err := signChannel.Send(&vault.HMACSignStreamRequest{
			Data: data[i:endByte],
		})
		require.Nil(s.T(), err)
	}
	signResponse, err := signChannel.CloseAndRecv()
	require.Nil(s.T(), err)

	badSignature := signResponse.Signature
	badSignature[0] = ^badSignature[0]

	verifyChannel, err := s.systemStub.Vault.HMACVerifyStream(ctx)
	require.Nil(s.T(), err)
	for i := 0; i < len(data); i += 1234 {
		endByte := i + 1234
		if len(data) < endByte {
			endByte = len(data)
		}
		err := verifyChannel.Send(&vault.HMACVerifyStreamRequest{
			Data:      data[i:endByte],
			Signature: badSignature,
		})
		require.Nil(s.T(), err)
	}
	verifyResponse, err := verifyChannel.CloseAndRecv()
	require.Nil(s.T(), err)
	require.False(s.T(), verifyResponse.Valid)
}

func (s *HMACTestSuite) TestSignAndVerify() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	data := tools.GetRandomBytes(100000)

	signResponse, err := s.systemStub.Vault.HMACSign(ctx, &vault.HMACSignRequest{
		Data: data,
	})
	require.Nil(s.T(), err)

	verifyResponse, err := s.systemStub.Vault.HMACVerify(ctx, &vault.HMACVerifyRequest{
		Data:      data,
		Signature: signResponse.Signature,
	})
	require.Nil(s.T(), err)
	require.True(s.T(), verifyResponse.Valid)
}

func (s *HMACTestSuite) TestVerifyBadSignature() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data := tools.GetRandomBytes(100000)

	signResponse, err := s.systemStub.Vault.HMACSign(ctx, &vault.HMACSignRequest{
		Data: data,
	})
	require.Nil(s.T(), err)

	badSignature := signResponse.Signature
	badSignature[0] = ^badSignature[0]

	verifyResponse, err := s.systemStub.Vault.HMACVerify(ctx, &vault.HMACVerifyRequest{
		Data:      data,
		Signature: badSignature,
	})
	require.Nil(s.T(), err)
	require.False(s.T(), verifyResponse.Valid)
}
