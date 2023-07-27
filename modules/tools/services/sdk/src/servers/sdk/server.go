package sdk

import (
	"context"

	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	sdk "github.com/slamy-solutions/openbp/modules/tools/libs/golang/sdk/sdk"
	sdkTools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SDKServer struct {
	sdk.UnimplementedSDKServiceServer

	nativeStub *native.NativeStub

	authHandler sdkTools.AuthHandler
	logger      *logrus.Entry
}

func NewSDKServer(authHandler sdkTools.AuthHandler, logger *logrus.Entry, nativeStub *native.NativeStub) *SDKServer {
	return &SDKServer{
		authHandler: authHandler,
		logger:      logger,
		nativeStub:  nativeStub,
	}
}

func (s *SDKServer) loggerForEndpoint(endpointName string) *logrus.Entry {
	return s.logger.WithField("endpoint", endpointName)
}

func (s *SDKServer) DownloadTLSResources(ctx context.Context, in *sdk.DownloadTLSResourcesRequest) (*sdk.DownloadTLSResourcesResponse, error) {
	_, err := s.authHandler.Authorize(ctx, &[]sdkTools.Scope{})
	if err != nil {
		st := sdkTools.AuthHandlerErrorToGRPCStatus(err)
		if st.Code() == codes.Internal {
			s.loggerForEndpoint("DownloadTLSResources").Error("Error while handling request authorization: " + err.Error())
		}
		return nil, st.Err()
	}

	return nil, status.Errorf(codes.Unimplemented, "method DownloadTLSResources not implemented")
}
func (s *SDKServer) Ping(ctx context.Context, in *sdk.PingRequest) (*sdk.PingResponse, error) {
	_, err := s.authHandler.Authorize(ctx, &[]sdkTools.Scope{})
	if err != nil {
		st := sdkTools.AuthHandlerErrorToGRPCStatus(err)
		if st.Code() == codes.Internal {
			s.loggerForEndpoint("Ping").Error("Error while handling request authorization: " + err.Error())
		}
		return nil, st.Err()
	}

	return &sdk.PingResponse{}, status.Error(codes.OK, "")
}

func (s *SDKServer) RegisterPublicKeyAsUser(ctx context.Context, in *sdk.RegisterPublicKeyAsUserRequest) (*sdk.RegisterPublicKeyAsUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
	/*
		getByLoginResponse, err := s.nativeStub.Services.ActorUser.GetByLogin(ctx, &user.GetByLoginRequest{
			Namespace: in.Namespace,
			Login:     in.User,
			UseCache:  false,
		})

		if err != nil {
			st := sdkTools.AuthHandlerErrorToGRPCStatus(err)
			if st.Code() == codes.NotFound {
				return nil, status.Error(codes.Unauthenticated, "")
			}

			s.loggerForEndpoint("RegisterPublicKeyAsUser").Error("Error while searching for user: " + err.Error())
			return nil, status.Error(codes.Internal, "")
		}

		_, err = s.nativeStub.Services.IAM.Auth.CheckAccessWithPassword(ctx, &auth.CheckAccessWithPasswordRequest{
			Namespace: in.Namespace,
			Identity:  getByLoginResponse.User.Identity,
			Password:  in.Password,
			Metadata:  "{}",
			Scopes:    []*auth.Scope{},
		})

		if err != nil {
			s.loggerForEndpoint("RegisterPublicKeyAsUser").Error("Error while checking access for user: " + err.Error())
			return nil, status.Error(codes.Internal, "")
		}

		createCertificateResponse, err := s.nativeStub.Services.IAM.Authentication.X509.RegisterAndGenerate(ctx, &x509.RegisterAndGenerateRequest{
			Namespace:   in.Namespace,
			Identity:    getByLoginResponse.User.Identity,
			PublicKey:   in.PublicKey,
			Description: "",
		})

		if err != nil {
			s.loggerForEndpoint("RegisterPublicKeyAsUser").Error("Error while generating certificate: " + err.Error())
			return nil, status.Error(codes.Internal, "")
		}

		return &sdk.RegisterPublicKeyAsUserResponse{
			Certificate: createCertificateResponse.Raw,
		}, status.Error(codes.OK, "")
	*/
}
