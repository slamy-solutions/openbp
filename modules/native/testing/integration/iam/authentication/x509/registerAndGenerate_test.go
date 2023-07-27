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
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/tools/testing/tools"
)

type RegisterAndGenerateInfoTestSuite struct {
	suite.Suite

	nativeStub *native.NativeStub
}

func (suite *RegisterAndGenerateInfoTestSuite) SetupSuite() {
	suite.nativeStub = native.NewNativeStub(native.NewStubConfig().WithIAMService().WithNamespaceService())
	err := suite.nativeStub.Connect()
	if err != nil {
		panic(err)
	}
}
func (suite *RegisterAndGenerateInfoTestSuite) TearDownSuite() {
	suite.nativeStub.Close()
}
func TestRegisterAndGenerateInfoTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterAndGenerateInfoTestSuite))
}

func (s *RegisterAndGenerateInfoTestSuite) TestRegisterAndGenerate() {
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

	_, err = cryptoX509.ParseCertificate(response.Raw)
	require.Nil(s.T(), err)
}

func (s *RegisterAndGenerateInfoTestSuite) TestRegisterAndGenerateInNamespace() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	namespaceName := tools.GetRandomString(32)
	_, err := s.nativeStub.Services.Namespace.Create(ctx, &namespace.CreateNamespaceRequest{
		Name:        namespaceName,
		FullName:    tools.GetRandomString(32),
		Description: tools.GetRandomString(32),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.Namespace.Delete(context.Background(), &namespace.DeleteNamespaceRequest{Name: namespaceName})

	createIdentityResponse, err := s.nativeStub.Services.IAM.Identity.Create(ctx, &identity.CreateIdentityRequest{
		Namespace:       namespaceName,
		Name:            tools.GetRandomString(20),
		InitiallyActive: true,
		Managed:         &identity.CreateIdentityRequest_No{No: &identity.NotManagedData{}},
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Identity.Delete(context.Background(), &identity.DeleteIdentityRequest{Namespace: namespaceName, Uuid: createIdentityResponse.Identity.Uuid})

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(s.T(), err)

	response, err := s.nativeStub.Services.IAM.Authentication.X509.RegisterAndGenerate(ctx, &x509.RegisterAndGenerateRequest{
		Namespace: namespaceName,
		Identity:  createIdentityResponse.Identity.Uuid,
		PublicKey: cryptoX509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	})
	require.Nil(s.T(), err)
	defer s.nativeStub.Services.IAM.Authentication.X509.Delete(context.Background(), &x509.DeleteRequest{Namespace: namespaceName, Uuid: response.Info.Uuid})

	_, err = cryptoX509.ParseCertificate(response.Raw)
	require.Nil(s.T(), err)
}
