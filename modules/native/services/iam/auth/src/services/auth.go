package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/status"

	grpccodes "google.golang.org/grpc/codes"

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"

	nativeIAmAuthGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_auth"
	nativeIAmAuthenticationPasswordGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_authentication_password"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_identity"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_policy"
	nativeNamespaceGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_namespace"
)

type IAmAuthServer struct {
	nativeIAmAuthGRPC.UnimplementedIAMAuthServiceServer

	mongoClient                           *mongo.Client
	mongoDbPrefix                         string
	cacheClient                           cache.Cache
	nativeNamespaceClient                 nativeNamespaceGRPC.NamespaceServiceClient
	nativeIAmPolicyClient                 nativeIAmPolicyGRPC.IAMPolicyServiceClient
	nativeIAmIdentityClient               nativeIAmIdentityGRPC.IAMIdentityServiceClient
	nativeIAmAuthenticationPasswordClient nativeIAmAuthenticationPasswordGRPC.IAMAuthenticationPasswordServiceClient
}

func (s *IAmAuthServer) CreateTokenWithPassword(ctx context.Context, in *nativeIAmAuthGRPC.CreateTokenWithPasswordRequest) (*nativeIAmAuthGRPC.CreateTokenWithPasswordResponse, error) {
	// Authenticate
	authenticateResponse, err := s.nativeIAmAuthenticationPasswordClient.Authenticate(ctx, &nativeIAmAuthenticationPasswordGRPC.AuthenticateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Password:  in.Password,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error while performing authentication: "+err.Error())
	}
	if !authenticateResponse.Authenticated {
		return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHENTICATED, AccessToken: "", RefreshToken: ""}, nil
	}

	// Check if identity is not disabled
	identityResponse, err := s.nativeIAmIdentityClient.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{Namespace: in.Namespace, Uuid: in.Identity, UseCache: true})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == grpccodes.NotFound {
				return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHENTICATED, AccessToken: "", RefreshToken: ""}, nil
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get identity information: "+err.Error())
	}
	if !identityResponse.Identity.Active {
		return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHORIZED, AccessToken: "", RefreshToken: ""}, nil
	}

	// Generate new login token

	return nil, status.Errorf(grpccodes.Unimplemented, "method CreateTokenWithPassword not implemented")
}
func (s *IAmAuthServer) RefreshToken(ctx context.Context, in *nativeIAmAuthGRPC.RefreshTokenRequest) (*nativeIAmAuthGRPC.RefreshTokenResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method RefreshToken not implemented")
}
func (s *IAmAuthServer) InvalidateToken(ctx context.Context, in *nativeIAmAuthGRPC.InvalidateTokenRequest) (*nativeIAmAuthGRPC.InvalidateTokenResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method InvalidateToken not implemented")
}
func (s *IAmAuthServer) VerifyTokenAccess(ctx context.Context, in *nativeIAmAuthGRPC.VerifyTokenAccessRequest) (*nativeIAmAuthGRPC.VerifyTokenAccessResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method VerifyTokenAccess not implemented")
}
