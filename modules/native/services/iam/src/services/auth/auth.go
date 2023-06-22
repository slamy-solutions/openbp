package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/status"

	grpccodes "google.golang.org/grpc/codes"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	nativeIAmAuthGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	nativeIAmAuthenticationPasswordGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmRoleGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	nativeIAmTokenGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"

	authentication_password "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication"
	identity_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/identity"
	policy_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/policy"
	role_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/role"
	token_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/token"
)

type IAmAuthServer struct {
	nativeIAmAuthGRPC.UnimplementedIAMAuthServiceServer

	systemStub *system.SystemStub

	authenticationPasswordServer *authentication_password.PasswordIdentificationService
	identityServer               *identity_server.IAmIdentityServer
	policyServer                 *policy_server.IAMPolicyServer
	roleServer                   *role_server.IAMRoleServer
	tokenServer                  *token_server.IAmTokenServer
}

func NewIAmAuthServer(
	systemStub *system.SystemStub,
	authenticationPasswordServer *authentication_password.PasswordIdentificationService,
	identityServer *identity_server.IAmIdentityServer,
	policyServer *policy_server.IAMPolicyServer,
	roleServer *role_server.IAMRoleServer,
	tokenServer *token_server.IAmTokenServer,
) *IAmAuthServer {
	return &IAmAuthServer{
		systemStub:                   systemStub,
		authenticationPasswordServer: authenticationPasswordServer,
		identityServer:               identityServer,
		policyServer:                 policyServer,
		roleServer:                   roleServer,
		tokenServer:                  tokenServer,
	}
}

func (s *IAmAuthServer) fetchIdentityPolicies(ctx context.Context, identity *nativeIAmIdentityGRPC.Identity) ([]*nativeIAmPolicyGRPC.Policy, error) {
	searchedPolicies := make([]*nativeIAmPolicyGRPC.GetMultiplePoliciesRequest_RequestedPolicy, len(identity.Policies))
	for index, policy := range identity.Policies {
		searchedPolicies[index] = &nativeIAmPolicyGRPC.GetMultiplePoliciesRequest_RequestedPolicy{
			Namespace: policy.Namespace,
			Uuid:      policy.Uuid,
		}
	}

	// Get all the identity roles. Get list of policies assigned for each role
	if len(identity.Roles) > 0 {
		requestedIdentityRoles := make([]*nativeIAmRoleGRPC.GetMultipleRolesRequest_RequestedRole, len(identity.Roles))
		for index, role := range identity.Roles {
			requestedIdentityRoles[index] = &nativeIAmRoleGRPC.GetMultipleRolesRequest_RequestedRole{
				Namespace: role.Namespace,
				Uuid:      role.Uuid,
			}
		}
		responseStream := s.roleServer.OpenGetMultipleChannel(ctx, &nativeIAmRoleGRPC.GetMultipleRolesRequest{Roles: requestedIdentityRoles})

		var err error = nil
		for role := range responseStream {
			if err != nil {
				continue
			}
			if role.Err != nil {
				err = role.Err
				continue
			}

			for _, policy := range role.Role.Policies {
				searchedPolicies = append(searchedPolicies, &nativeIAmPolicyGRPC.GetMultiplePoliciesRequest_RequestedPolicy{
					Namespace: policy.Namespace,
					Uuid:      policy.UUID,
				})
			}
		}

		// TODO: make sure policies are unique

		if err != nil {
			return nil, errors.New("error while fetching all roles for the identity: " + err.Error())
		}
	}

	policies := make([]*nativeIAmPolicyGRPC.Policy, 0, len(searchedPolicies))
	policiesStream := s.policyServer.OpenGetMultipleChannel(ctx, &nativeIAmPolicyGRPC.GetMultiplePoliciesRequest{
		Policies: searchedPolicies,
	})

	var err error = nil
	for policy := range policiesStream {
		if err != nil {
			continue
		}
		if policy.Err != nil {
			err = policy.Err
			continue
		}

		policies = append(policies, policy.Policy.ToGRPCPolicy(policy.Namespace))
	}

	if err != nil {
		return nil, errors.New("error while fetching all policies for the identity: " + err.Error())
	}

	return policies, nil
}

func (s *IAmAuthServer) CreateTokenWithPassword(ctx context.Context, in *nativeIAmAuthGRPC.CreateTokenWithPasswordRequest) (*nativeIAmAuthGRPC.CreateTokenWithPasswordResponse, error) {

	// TODO: refactor this function. Refactor usage of the pointers

	// Authenticate
	authenticateResponse, err := s.authenticationPasswordServer.Authenticate(ctx, &nativeIAmAuthenticationPasswordGRPC.AuthenticateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Password:  in.Password,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error while performing authentication: "+err.Error())
	}
	if !authenticateResponse.Authenticated {
		return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_CREDENTIALS_INVALID, AccessToken: "", RefreshToken: ""}, nil
	}

	// Check if identity is not disabled. Dont use cache, because it is rare operation and invalid cache will result in allowing access for very long period of time
	identityResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{Namespace: in.Namespace, Uuid: in.Identity, UseCache: false})
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
		return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_IDENTITY_NOT_ACTIVE, AccessToken: "", RefreshToken: ""}, nil
	}

	// Get all policies for identity
	policies, err := s.fetchIdentityPolicies(ctx, identityResponse.Identity)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, in.Scopes) {
		return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHORIZED, AccessToken: "", RefreshToken: ""}, nil
	}

	// Get all policies as scopes if list of scopes in request is empty
	var scopesToAssign []*nativeIAmTokenGRPC.Scope
	if len(in.Scopes) == 0 {
		scopesToAssign = make([]*nativeIAmTokenGRPC.Scope, len(policies))
		for i, policy := range policies {
			scopesToAssign[i] = &nativeIAmTokenGRPC.Scope{
				Namespace:            policy.Namespace,
				Resources:            policy.Resources,
				Actions:              policy.Actions,
				NamespaceIndependent: policy.NamespaceIndependent,
			}
		}
	} else {
		scopesToAssign = make([]*nativeIAmTokenGRPC.Scope, len(in.Scopes))
		for i, scope := range in.Scopes {
			scopesToAssign[i] = &nativeIAmTokenGRPC.Scope{
				Namespace:            scope.Namespace,
				Resources:            scope.Resources,
				Actions:              scope.Actions,
				NamespaceIndependent: scope.NamespaceIndependent,
			}
		}
	}

	// Generate new login token
	tokenResponse, err := s.tokenServer.Create(ctx, &nativeIAmTokenGRPC.CreateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Scopes:    scopesToAssign,
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
	refershResponse, err := s.tokenServer.Refresh(ctx, &nativeIAmTokenGRPC.RefreshRequest{
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to refresh token token. "+err.Error())
	}
	switch refershResponse.Status {
	case nativeIAmTokenGRPC.RefreshResponse_OK:
		break
	case nativeIAmTokenGRPC.RefreshResponse_DISABLED:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_DISABLED}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_EXPIRED:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_EXPIRED}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_INVALID:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_INVALID}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_NOT_FOUND:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_NOT_FOUND}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_NOT_REFRESH_TOKEN:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_IS_NOT_REFRESH_TOKEN}, status.Error(grpccodes.OK, "")
	default:
		return nil, status.Error(grpccodes.Internal, "Unknow refresh response from the native_iam_token service. Received status: "+refershResponse.Status.String())
	}

	identityResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{
		Namespace: refershResponse.TokenData.Namespace,
		Uuid:      refershResponse.TokenData.Identity,
		UseCache:  false, //Token refresh is not frequent operation. Using invalid data will allow to use token for very long period of time.
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// This error must not occure in healthy system
			if st.Code() == grpccodes.NotFound {
				return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_IDENTITY_NOT_FOUND}, status.Error(grpccodes.OK, "")
			}
		}
		return nil, status.Error(grpccodes.Internal, "Failed to get identity of the token. Error: "+err.Error())
	}

	if !identityResponse.Identity.Active {
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_IDENTITY_NOT_ACTIVE}, status.Error(grpccodes.OK, "")
	}

	policies, err := s.fetchIdentityPolicies(ctx, identityResponse.Identity)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !areTokenScopesValidForIdentityScopes(policies, refershResponse.TokenData.Scopes) {
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_IDENTITY_UNAUTHENTICATED}, status.Error(grpccodes.OK, "")
	}

	return &nativeIAmAuthGRPC.RefreshTokenResponse{
		Status:      nativeIAmAuthGRPC.RefreshTokenResponse_OK,
		AccessToken: refershResponse.Token,
	}, status.Error(grpccodes.OK, "")
}

func (s *IAmAuthServer) CheckAccessWithToken(ctx context.Context, in *nativeIAmAuthGRPC.CheckAccessWithTokenRequest) (*nativeIAmAuthGRPC.CheckAccessWithTokenResponse, error) {
	tokenResponse, err := s.tokenServer.Validate(ctx, &nativeIAmTokenGRPC.ValidateRequest{
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
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_EXPIRED, Message: "Token expired"}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_DISABLED:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_DISABLED, Message: "Token was manually disabled"}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_INVALID:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_INVALID, Message: "Token invalid. Maybe it has bad structure or signature"}, status.Error(grpccodes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_NOT_FOUND:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_NOT_FOUND, Message: "Token not found. Most probably it was deteled and cant be used."}, status.Error(grpccodes.OK, "")
	}

	if !areTokenScopesAllowAccess(tokenResponse.TokenData.Scopes, in.Scopes) {
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{
			Status:  nativeIAmAuthGRPC.CheckAccessWithTokenResponse_UNAUTHORIZED,
			Message: "Token doesnt have enought privileges to access provided scopes", //TODO: add here information about additional required policies
		}, status.Error(grpccodes.OK, "")
	}

	return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{
		Status:       nativeIAmAuthGRPC.CheckAccessWithTokenResponse_OK,
		Message:      "",
		Namespace:    tokenResponse.TokenData.Namespace,
		TokenUUID:    tokenResponse.TokenData.Uuid,
		IdentityUUID: tokenResponse.TokenData.Identity,
	}, status.Error(grpccodes.OK, "")
}

func (s *IAmAuthServer) CheckAccessWithPassword(ctx context.Context, in *nativeIAmAuthGRPC.CheckAccessWithPasswordRequest) (*nativeIAmAuthGRPC.CheckAccessWithPasswordResponse, error) {
	//TODO: use provided metadata

	authenticateResponse, err := s.authenticationPasswordServer.Authenticate(ctx, &nativeIAmAuthenticationPasswordGRPC.AuthenticateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Password:  in.Password,
	})
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Error while authorizing using password. "+err.Error())
	}
	if !authenticateResponse.Authenticated {
		return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHENTICATED, Message: "Identity or password doesnt match."}, status.Error(grpccodes.OK, "")
	}

	// Find identity and its policies
	identityGetResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{
		Namespace: in.Namespace,
		Uuid:      in.Identity,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// This error must not occure in healthy system
			if st.Code() == grpccodes.NotFound {
				return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHENTICATED, Message: "Identity or password doesnt match."}, status.Error(grpccodes.OK, "")
			}
		}

		return nil, status.Error(grpccodes.Internal, "Error while searching for identity. "+err.Error())
	}

	if !identityGetResponse.Identity.Active {
		return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHORIZED, Message: "Identity is not active."}, status.Error(grpccodes.OK, "")
	}

	//TODO: Cache this
	// Get all policies for identity
	policies, err := s.fetchIdentityPolicies(ctx, identityGetResponse.Identity)
	if err != nil {
		return nil, status.Error(grpccodes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, in.Scopes) {
		return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHORIZED, Message: "Not enought privileges"}, status.Error(grpccodes.OK, "")
	}

	return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_OK, Message: ""}, status.Error(grpccodes.OK, "")
}
