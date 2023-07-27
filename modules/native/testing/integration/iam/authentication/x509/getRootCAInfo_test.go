package x509

import (
	"context"
	cryptoX509 "crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
)

type GetRootCAInfoTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *GetRootCAInfoTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithIAMService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *GetRootCAInfoTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestGetRootCAInfoTestSuite(t *testing.T) {
	suite.Run(t, new(GetRootCAInfoTestSuite))
}

func (s *GetRootCAInfoTestSuite) TestGetInfo() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	response, err := s.nativeStub.Services.IAM.Authentication.X509.GetRootCAInfo(ctx, &x509.GetRootCAInfoRequest{})
	require.Nil(s.T(), err)

	_, err = cryptoX509.ParseCertificate(response.Certificate)
	require.Nil(s.T(), err)
}
