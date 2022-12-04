package services

import (
	"context"

	"google.golang.org/grpc/status"

	grpccodes "google.golang.org/grpc/codes"

	nativeIAmAuthenticationPasswordGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	nativeIAmOAuthGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmTokenGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
)

type IAmAuthServer struct {
	nativeIAmOAuthGRPC.UnimplementedIAMOAuthServiceServer

	nativeIAmPolicyClient                 nativeIAmPolicyGRPC.IAMPolicyServiceClient
	nativeIAmIdentityClient               nativeIAmIdentityGRPC.IAMIdentityServiceClient
	nativeIAmTokenClient                  nativeIAmTokenGRPC.IAMTokenServiceClient
	nativeIAmAuthenticationPasswordClient nativeIAmAuthenticationPasswordGRPC.IAMAuthenticationPasswordServiceClient
}

func NewIAmAuthServer(
	nativeIAmPolicyClient nativeIAmPolicyGRPC.IAMPolicyServiceClient,
	nativeIAmIdentityClient nativeIAmIdentityGRPC.IAMIdentityServiceClient,
	nativeIAmTokenClient nativeIAmTokenGRPC.IAMTokenServiceClient,
	nativeIAmAuthenticationPasswordClient nativeIAmAuthenticationPasswordGRPC.IAMAuthenticationPasswordServiceClient,
) *IAmAuthServer {
	return &IAmAuthServer{
		nativeIAmPolicyClient:                 nativeIAmPolicyClient,
		nativeIAmIdentityClient:               nativeIAmIdentityClient,
		nativeIAmTokenClient:                  nativeIAmTokenClient,
		nativeIAmAuthenticationPasswordClient: nativeIAmAuthenticationPasswordClient,
	}
}

func fetchIdentityPolicies(ctx context.Context, policyClient nativeIAmPolicyGRPC.IAMPolicyServiceClient, policyReferences []*nativeIAmIdentityGRPC.Identity_PolicyReference) ([]*nativeIAmPolicyGRPC.Policy, error) {
	// TODO: do this in parallel
	policies := make([]*nativeIAmPolicyGRPC.Policy, len(policyReferences))
	for i, policy := range policyReferences {
		policyResponse, err := policyClient.Get(ctx, &nativeIAmPolicyGRPC.GetPolicyRequest{
			Namespace: policy.Namespace,
			Uuid:      policy.Uuid,
			UseCache:  false, // CreateTokenWithPassword is very rare operation. Invalid cache will result in token to be valid for a very long period of time.
		})
		if err != nil {
			return nil, err
		}
		policies[i] = policyResponse.Policy
	}
	return policies, nil
}

func (s *IAmAuthServer) CreateTokenWithPassword(ctx context.Context, in *nativeIAmOAuthGRPC.CreateTokenWithPasswordRequest) (*nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse, error) {

	// TODO: refactor this function. Refactor usage of the pointers

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
		return &nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse_CREDENTIALS_INVALID, AccessToken: "", RefreshToken: ""}, nil
	}

	// Check if identity is not disabled. Dont use cache, because it is rare operation and invalid cache will result in allowing access for very long period of time
	identityResponse, err := s.nativeIAmIdentityClient.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{Namespace: in.Namespace, Uuid: in.Identity, UseCache: false})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// This error must not occure in healthy system
			if st.Code() == grpccodes.NotFound {
				return nil, status.Error(grpccodes.Internal, "Failed to get identity information. Identity not found.")
			}
		}

		return nil, status.Error(grpccodes.Internal, "Failed to get identity information: "+err.Error())
	}
	if !identityResponse.Identity.Active {
		return &nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse_IDENTITY_NOT_ACTIVE, AccessToken: "", RefreshToken: ""}, nil
	}

	// Get all policies for identity
	policies, err := fetchIdentityPolicies(ctx, s.nativeIAmPolicyClient, identityResponse.Identity.Policies)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, in.Scopes) {
		return &nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHORIZED, AccessToken: "", RefreshToken: ""}, nil
	}

	// Get all policies as scopes if list of scopes in request is empty
	var scopesToAssign []*nativeIAmTokenGRPC.Scope
	if len(in.Scopes) == 0 {
		scopesToAssign = make([]*nativeIAmTokenGRPC.Scope, len(policies))
		for i, policy := range policies {
			scopesToAssign[i] = &nativeIAmTokenGRPC.Scope{
				Namespace: policy.Namespace,
				Resources: policy.Resources,
				Actions:   policy.Actions,
			}
		}
	} else {
		scopesToAssign = make([]*nativeIAmTokenGRPC.Scope, len(in.Scopes))
		for i, scope := range in.Scopes {
			scopesToAssign[i] = &nativeIAmTokenGRPC.Scope{
				Namespace: scope.Namespace,
				Resources: scope.Resources,
				Actions:   scope.Actions,
			}
		}
	}

	// Generate new login token
	tokenResponse, err := s.nativeIAmTokenClient.Create(ctx, &nativeIAmTokenGRPC.CreateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Scopes:    scopesToAssign,
		Metadata:  in.Metadata,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to create token. "+err.Error())
	}

	return &nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse{
		Status:       nativeIAmOAuthGRPC.CreateTokenWithPasswordResponse_OK,
		AccessToken:  tokenResponse.Token,
		RefreshToken: tokenResponse.RefreshToken,
	}, status.Error(grpccodes.OK, "")
}

func (s *IAmAuthServer) RefreshToken(ctx context.Context, in *nativeIAmOAuthGRPC.RefreshTokenRequest) (*nativeIAmOAuthGRPC.RefreshTokenResponse, error) {
	refershResponse, err := s.nativeIAmTokenClient.Refresh(ctx, &nativeIAmTokenGRPC.RefreshRequest{
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to refresh token token. "+err.Error())
	}
	switch refershResponse.Status {
	case nativeIAmTokenGRPC.RefreshResponse_OK:
		break
	case nativeIAmTokenGRPC.RefreshResponse_DISABLED:
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_TOKEN_DISABLED}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_EXPIRED:
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_TOKEN_EXPIRED}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_INVALID:
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_TOKEN_INVALID}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_NOT_FOUND:
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_TOKEN_NOT_FOUND}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_NOT_REFRESH_TOKEN:
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_TOKEN_IS_NOT_REFRESH_TOKEN}, status.Error(grpccodes.OK, "")
	default:
		return nil, status.Error(grpccodes.Internal, "Unknow refresh response from the native_iam_token service. Received status: "+refershResponse.Status.String())
	}

	identityResponse, err := s.nativeIAmIdentityClient.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{
		Namespace: refershResponse.TokenData.Namespace,
		Uuid:      refershResponse.TokenData.Identity,
		UseCache:  false, //Token refresh is not frequent operation. Using invalid data will allow to use token for very long period of time.
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// This error must not occure in healthy system
			if st.Code() == grpccodes.NotFound {
				return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_IDENTITY_NOT_FOUND}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get identity of the token. Error: "+err.Error())
	}

	if !identityResponse.Identity.Active {
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_IDENTITY_NOT_ACTIVE}, status.Error(grpccodes.OK, "")
	}

	policies, err := fetchIdentityPolicies(ctx, s.nativeIAmPolicyClient, identityResponse.Identity.Policies)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !areTokenScopesValidForIdentityScopes(policies, refershResponse.TokenData.Scopes) {
		return &nativeIAmOAuthGRPC.RefreshTokenResponse{Status: nativeIAmOAuthGRPC.RefreshTokenResponse_IDENTITY_UNAUTHENTICATED}, status.Error(grpccodes.OK, "")
	}

	return &nativeIAmOAuthGRPC.RefreshTokenResponse{
		Status:      nativeIAmOAuthGRPC.RefreshTokenResponse_OK,
		AccessToken: refershResponse.Token,
	}, status.Error(grpccodes.OK, "")
}

func (s *IAmAuthServer) CheckAccess(ctx context.Context, in *nativeIAmOAuthGRPC.CheckAccessRequest) (*nativeIAmOAuthGRPC.CheckAccessResponse, error) {
	tokenResponse, err := s.nativeIAmTokenClient.Validate(ctx, &nativeIAmTokenGRPC.ValidateRequest{
		Token:    in.AccessToken,
		UseCache: true,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to validate token. "+err.Error())
	}

	switch tokenResponse.Status {
	case nativeIAmTokenGRPC.ValidateResponse_OK:
		break
	case nativeIAmTokenGRPC.ValidateResponse_EXPIRED:
		return &nativeIAmOAuthGRPC.CheckAccessResponse{Status: nativeIAmOAuthGRPC.CheckAccessResponse_TOKEN_EXPIRED, Message: "Token expired"}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_DISABLED:
		return &nativeIAmOAuthGRPC.CheckAccessResponse{Status: nativeIAmOAuthGRPC.CheckAccessResponse_TOKEN_DISABLED, Message: "Token was manually disabled"}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_INVALID:
		return &nativeIAmOAuthGRPC.CheckAccessResponse{Status: nativeIAmOAuthGRPC.CheckAccessResponse_TOKEN_INVALID, Message: "Token invalid. Maybe it has bad structure or signature"}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_NOT_FOUND:
		return &nativeIAmOAuthGRPC.CheckAccessResponse{Status: nativeIAmOAuthGRPC.CheckAccessResponse_TOKEN_NOT_FOUND, Message: "Token not found. Most probably it was deteled and cant be used."}, status.Error(grpccodes.OK, "")
	}

	if !areTokenScopesAllowAccess(tokenResponse.TokenData.Scopes, in.Scopes) {
		return &nativeIAmOAuthGRPC.CheckAccessResponse{
			Status:  nativeIAmOAuthGRPC.CheckAccessResponse_UNAUTHORIZED,
			Message: "Token doesnt have enought privileges to access provided scopes", //TODO: add here information about additional required policies
		}, status.Error(grpccodes.OK, "")
	}

	return &nativeIAmOAuthGRPC.CheckAccessResponse{
		Status:  nativeIAmOAuthGRPC.CheckAccessResponse_OK,
		Message: "",
	}, status.Error(grpccodes.OK, "")
}
