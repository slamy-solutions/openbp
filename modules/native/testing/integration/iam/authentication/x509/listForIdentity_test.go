package x509

import (
	"context"
	rand "crypto/rand"
	"crypto/rsa"
	cryptoX509 "crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/tools/testing/tools"
)

type ListForIdentityTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *ListForIdentityTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithIAMService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *ListForIdentityTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestListForIdentityTestSuite(t *testing.T) {
	suite.Run(t, new(ListForIdentityTestSuite))
}

func (s *ListForIdentityTestSuite) TestRegisterAndList() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	createIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       "",
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: "", Uuid: createIdentityResponse.Identity.Uuid})

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(s.T(), err)

	response, err := s.nativeStub.Services.IAM.Authentication.X509.RegisterAndGenerate(ctx, &x509.RegisterAndGenerateRequest{
		Namespace: "",
		Identity:  createIdentityResponse.Identity.Uuid,
		PublicKey: cryptoX509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Authentication.X509.Delete(context.Background(), &x509.DeleteRequest{Namespace: "", Uuid: response.Info.Uuid})

	listResponseStream, err := s.nativeStub.Services.IAM.Authentication.X509.ListForIdentity(ctx, &x509.ListForIdentityRequest{
		Namespace: "",
		Identity:  createIdentityResponse.Identity.Uuid,
		Skip:      0,
		Limit:     0,
	})
	require.Nil(s.T(), err)

	r, err := listResponseStream.Recv()
	require.Nil(s.T(), err)
	require.Equal(s.T(), response.Info.Uuid, r.Certificate.Uuid)
}
