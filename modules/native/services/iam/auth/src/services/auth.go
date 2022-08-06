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
	nativeIAmTokenGRPC "github.com/slamy-solutions/open-erp/modules/native/services/iam/auth/src/grpc/native_iam_token"
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
	nativeIAmTokenClient                  nativeIAmTokenGRPC.IAMTokenServiceClient
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

	// Create list of scopes for token
	// TODO: scopes from the request. At this point, system get all possible scopes for identity
	// TODO: do this in parallel

	scopes := make([]*nativeIAmTokenGRPC.Scope, len(identityResponse.Identity.Policies))
	for i, policy := range identityResponse.Identity.Policies {
		policyResponse, err := s.nativeIAmPolicyClient.Get(ctx, &nativeIAmPolicyGRPC.GetPolicyRequest{
			Namespace: policy.Namespace,
			Uuid:      policy.Uuid,
			UseCache:  false, // CreateTokenWithPassword is very rare operation. Invalid cache will result in token to be valid for a very long period of time.
		})
		if err != nil {
			return nil, status.Error(grpccodes.Internal, "Failed to get policy information: "+err.Error())
		}
		scopes[i] = &nativeIAmTokenGRPC.Scope{
			Namespace: policy.Namespace,
			Resources: policyResponse.Policy.Resources,
			Actions:   policyResponse.Policy.Actions,
		}
	}

	// Generate new login token
	tokenResponse, err := s.nativeIAmTokenClient.Create(ctx, &nativeIAmTokenGRPC.CreateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Scopes:    scopes,
		Metadata:  in.Metadata,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to create token. "+err.Error())
	}

	return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{
		Status:       nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_OK,
		AccessToken:  tokenResponse.Token,
		RefreshToken: tokenResponse.RefreshToken,
	}, status.Error(grpccodes.OK, "")
}
func (s *IAmAuthServer) RefreshToken(ctx context.Context, in *nativeIAmAuthGRPC.RefreshTokenRequest) (*nativeIAmAuthGRPC.RefreshTokenResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method RefreshToken not implemented")
}
func (s *IAmAuthServer) InvalidateToken(ctx context.Context, in *nativeIAmAuthGRPC.InvalidateTokenRequest) (*nativeIAmAuthGRPC.InvalidateTokenResponse, error) {
	return nil, status.Errorf(grpccodes.Unimplemented, "method InvalidateToken not implemented")
}
func (s *IAmAuthServer) VerifyTokenAccess(ctx context.Context, in *nativeIAmAuthGRPC.VerifyTokenAccessRequest) (*nativeIAmAuthGRPC.VerifyTokenAccessResponse, error) {
	tokenResponse, err := s.nativeIAmTokenClient.Authorize(ctx, &nativeIAmTokenGRPC.AuthorizeRequest{
		Token:    in.AccessToken,
		UseCache: true,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to athorize token. "+err.Error())
	}
	if tokenResponse.Status != nativeIAmTokenGRPC.AuthorizeResponse_OK {
		return &nativeIAmAuthGRPC.VerifyTokenAccessResponse{
			HasAccess: false,
		}, status.Error(grpccodes.OK, "")
	}

	hasAccess := false
	for _, scope := range tokenResponse.TokenData.Scopes {
		if scope.Namespace != in.Namespace {
			continue
		}
		hasAccess = RequestScope(scope.Resources, scope.Actions, in.Resources, in.Actions)
		if hasAccess {
			break
		}
	}

	return &nativeIAmAuthGRPC.VerifyTokenAccessResponse{
		HasAccess: hasAccess,
	}, status.Error(grpccodes.OK, "")
}
