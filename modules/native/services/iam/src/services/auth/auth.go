package auth

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	nativeIAmAuthGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	nativeIAmAuthenticationPasswordGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/x509"
	nativeIAmIdentityGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	nativeIAmPolicyGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	nativeIAmRoleGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	nativeIAmTokenGRPC "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"

	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/oauth"
	authentication_OAuth "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/oauth"
	authentication_password "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/password"
	authentication_x509 "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication/x509"
	identity_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/identity"
	policy_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/policy"
	role_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/role"
	token_server "github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/token"
)

type IAmAuthServer struct {
	nativeIAmAuthGRPC.UnimplementedIAMAuthServiceServer

	systemStub *system.SystemStub

	authenticationPasswordServer *authentication_password.PasswordIdentificationService
	authenticationX509Server     *authentication_x509.X509IdentificationServer
	authenticationOAuthServer    *authentication_OAuth.OAuthServer
	identityServer               *identity_server.IAmIdentityServer
	policyServer                 *policy_server.IAMPolicyServer
	roleServer                   *role_server.IAMRoleServer
	tokenServer                  *token_server.IAmTokenServer
}

func NewIAmAuthServer(
	systemStub *system.SystemStub,
	authenticationPasswordServer *authentication_password.PasswordIdentificationService,
	authenticationX509Server *authentication_x509.X509IdentificationServer,
	authenticationOAuthServer *authentication_OAuth.OAuthServer,
	identityServer *identity_server.IAmIdentityServer,
	policyServer *policy_server.IAMPolicyServer,
	roleServer *role_server.IAMRoleServer,
	tokenServer *token_server.IAmTokenServer,
) *IAmAuthServer {
	return &IAmAuthServer{
		systemStub:                   systemStub,
		authenticationPasswordServer: authenticationPasswordServer,
		authenticationX509Server:     authenticationX509Server,
		authenticationOAuthServer:    authenticationOAuthServer,
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

var errCreateTokenIdentityNotActive = errors.New("identity not active")
var errCreateTokenUnauthorized = errors.New("unauthorized")

func (s *IAmAuthServer) createTokenForIdentity(ctx context.Context, namespace string, identity string, scopes []*nativeIAmAuthGRPC.Scope, metadata string) (string, string, error) {
	// Check if identity is not disabled. Dont use cache, because it is rare operation and invalid cache will result in allowing access for very long period of time
	identityResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{Namespace: namespace, Uuid: identity, UseCache: false})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// This error must not occure in healthy system
			if st.Code() == codes.NotFound {
				return "", "", status.Error(codes.Internal, "Failed to get identity information. Identity not found.")
			}
		}

		return "", "", status.Error(codes.Internal, "Failed to get identity information: "+err.Error())
	}
	if !identityResponse.Identity.Active {
		return "", "", errCreateTokenIdentityNotActive
	}

	// Get all policies for identity
	policies, err := s.fetchIdentityPolicies(ctx, identityResponse.Identity)
	if err != nil {
		return "", "", status.Error(codes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, scopes) {
		return "", "", errCreateTokenUnauthorized
	}

	// Get all policies as scopes if list of scopes in request is empty
	var scopesToAssign []*nativeIAmTokenGRPC.Scope
	if len(scopes) == 0 {
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
		scopesToAssign = make([]*nativeIAmTokenGRPC.Scope, len(scopes))
		for i, scope := range scopes {
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
		Namespace: namespace,
		Identity:  identity,
		Scopes:    scopesToAssign,
		Metadata:  metadata,
	})
	if err != nil {
		return "", "", status.Error(codes.Internal, "Failed to create token. "+err.Error())
	}

	return tokenResponse.Token, tokenResponse.RefreshToken, nil
}

func (s *IAmAuthServer) CreateTokenWithPassword(ctx context.Context, in *nativeIAmAuthGRPC.CreateTokenWithPasswordRequest) (*nativeIAmAuthGRPC.CreateTokenWithPasswordResponse, error) {
	authenticateResponse, err := s.authenticationPasswordServer.Authenticate(ctx, &nativeIAmAuthenticationPasswordGRPC.AuthenticateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Password:  in.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while performing authentication: "+err.Error())
	}
	if !authenticateResponse.Authenticated {
		return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_CREDENTIALS_INVALID}, nil
	}

	accessToken, refreshToken, err := s.createTokenForIdentity(ctx, in.Namespace, in.Identity, in.Scopes, in.Metadata)
	if err != nil {
		if err == errCreateTokenIdentityNotActive {
			return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_IDENTITY_NOT_ACTIVE}, nil
		}
		if err == errCreateTokenUnauthorized {
			return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{Status: nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHORIZED}, nil
		}

		return nil, err
	}

	return &nativeIAmAuthGRPC.CreateTokenWithPasswordResponse{
		Status:       nativeIAmAuthGRPC.CreateTokenWithPasswordResponse_OK,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, status.Error(codes.OK, "")
}

func (s *IAmAuthServer) CreateTokenWithOAuth2(ctx context.Context, in *nativeIAmAuthGRPC.CreateTokenWithOAuth2Request) (*nativeIAmAuthGRPC.CreateTokenWithOAuth2Response, error) {
	providerType, ok := oauth.ProviderStringTypeToGRPC[in.Provider]
	if !ok {
		return nil, status.Errorf(codes.Unimplemented, "Uninplemented provider type: %s", in.Provider)
	}

	authenticateResponse, err := s.authenticationOAuthServer.Authenticate(ctx, &oauth2.AuthenticateRequest{
		Namespace:    in.Namespace,
		Provider:     providerType,
		Code:         in.Code,
		CodeVerifier: in.CodeVerifier,
		RedirectUrl:  in.RedirectURL,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			return nil, status.Error(codes.FailedPrecondition, "failed to authenticate. most probably vault is sealed: "+err.Error())
		}

		return nil, status.Error(codes.Internal, "failed to authenticate: "+err.Error())
	}

	switch authenticateResponse.Status {
	case oauth2.AuthenticateResponse_OK:
		break
	case oauth2.AuthenticateResponse_UNAUTHENTICATED, oauth2.AuthenticateResponse_ERROR_WHILE_FETCHING_USER_DETAILS, oauth2.AuthenticateResponse_ERROR_WHILE_RETRIEVING_AUTH_TOKEN:
		return &nativeIAmAuthGRPC.CreateTokenWithOAuth2Response{Status: nativeIAmAuthGRPC.CreateTokenWithOAuth2Response_UNAUTHENTICATED}, nil
	default:
		return nil, status.Errorf(codes.Internal, "Unknown authenticate response status: %s", authenticateResponse.Status.String())

	}

	accessToken, refreshToken, err := s.createTokenForIdentity(ctx, in.Namespace, authenticateResponse.Identity, in.Scopes, in.Metadata)
	if err != nil {
		if err == errCreateTokenIdentityNotActive {
			return &nativeIAmAuthGRPC.CreateTokenWithOAuth2Response{Status: nativeIAmAuthGRPC.CreateTokenWithOAuth2Response_IDENTITY_NOT_ACTIVE}, nil
		}
		if err == errCreateTokenUnauthorized {
			return &nativeIAmAuthGRPC.CreateTokenWithOAuth2Response{Status: nativeIAmAuthGRPC.CreateTokenWithOAuth2Response_UNAUTHORIZED}, nil
		}

		return nil, err
	}

	return &nativeIAmAuthGRPC.CreateTokenWithOAuth2Response{
		Status:       nativeIAmAuthGRPC.CreateTokenWithOAuth2Response_OK,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, status.Error(codes.OK, "")
}

func (s *IAmAuthServer) RefreshToken(ctx context.Context, in *nativeIAmAuthGRPC.RefreshTokenRequest) (*nativeIAmAuthGRPC.RefreshTokenResponse, error) {
	refershResponse, err := s.tokenServer.Refresh(ctx, &nativeIAmTokenGRPC.RefreshRequest{
		RefreshToken: in.RefreshToken,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to refresh token token. "+err.Error())
	}
	switch refershResponse.Status {
	case nativeIAmTokenGRPC.RefreshResponse_OK:
		break
	case nativeIAmTokenGRPC.RefreshResponse_DISABLED:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_DISABLED}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_EXPIRED:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_EXPIRED}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_INVALID:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_INVALID}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_NOT_FOUND:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_NOT_FOUND}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.RefreshResponse_NOT_REFRESH_TOKEN:
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_TOKEN_IS_NOT_REFRESH_TOKEN}, status.Error(codes.OK, "")
	default:
		return nil, status.Error(codes.Internal, "Unknow refresh response from the native_iam_token service. Received status: "+refershResponse.Status.String())
	}

	identityResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{
		Namespace: refershResponse.TokenData.Namespace,
		Uuid:      refershResponse.TokenData.Identity,
		UseCache:  false, //Token refresh is not frequent operation. Using invalid data will allow to use token for very long period of time.
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// This error must not occure in healthy system
			if st.Code() == codes.NotFound {
				return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_IDENTITY_NOT_FOUND}, status.Error(codes.OK, "")
			}
		}
		return nil, status.Error(codes.Internal, "Failed to get identity of the token. Error: "+err.Error())
	}

	if !identityResponse.Identity.Active {
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_IDENTITY_NOT_ACTIVE}, status.Error(codes.OK, "")
	}

	policies, err := s.fetchIdentityPolicies(ctx, identityResponse.Identity)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !areTokenScopesValidForIdentityScopes(policies, refershResponse.TokenData.Scopes) {
		return &nativeIAmAuthGRPC.RefreshTokenResponse{Status: nativeIAmAuthGRPC.RefreshTokenResponse_IDENTITY_UNAUTHENTICATED}, status.Error(codes.OK, "")
	}

	return &nativeIAmAuthGRPC.RefreshTokenResponse{
		Status:      nativeIAmAuthGRPC.RefreshTokenResponse_OK,
		AccessToken: refershResponse.Token,
	}, status.Error(codes.OK, "")
}

func (s *IAmAuthServer) CheckAccessWithToken(ctx context.Context, in *nativeIAmAuthGRPC.CheckAccessWithTokenRequest) (*nativeIAmAuthGRPC.CheckAccessWithTokenResponse, error) {
	tokenResponse, err := s.tokenServer.Validate(ctx, &nativeIAmTokenGRPC.ValidateRequest{
		Token:    in.AccessToken,
		UseCache: true,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to validate token. "+err.Error())
	}

	switch tokenResponse.Status {
	case nativeIAmTokenGRPC.ValidateResponse_OK:
		break
	case nativeIAmTokenGRPC.ValidateResponse_EXPIRED:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_EXPIRED, Message: "Token expired"}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_DISABLED:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_DISABLED, Message: "Token was manually disabled"}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_INVALID:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_INVALID, Message: "Token invalid. Maybe it has bad structure or signature"}, status.Error(codes.OK, "")
	case nativeIAmTokenGRPC.ValidateResponse_NOT_FOUND:
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{Status: nativeIAmAuthGRPC.CheckAccessWithTokenResponse_TOKEN_NOT_FOUND, Message: "Token not found. Most probably it was deteled and cant be used."}, status.Error(codes.OK, "")
	}

	if !areTokenScopesAllowAccess(tokenResponse.TokenData.Scopes, in.Scopes) {
		return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{
			Status:  nativeIAmAuthGRPC.CheckAccessWithTokenResponse_UNAUTHORIZED,
			Message: "Token doesnt have enought privileges to access provided scopes", //TODO: add here information about additional required policies
		}, status.Error(codes.OK, "")
	}

	return &nativeIAmAuthGRPC.CheckAccessWithTokenResponse{
		Status:       nativeIAmAuthGRPC.CheckAccessWithTokenResponse_OK,
		Message:      "",
		Namespace:    tokenResponse.TokenData.Namespace,
		TokenUUID:    tokenResponse.TokenData.Uuid,
		IdentityUUID: tokenResponse.TokenData.Identity,
	}, status.Error(codes.OK, "")
}

func (s *IAmAuthServer) CheckAccessWithPassword(ctx context.Context, in *nativeIAmAuthGRPC.CheckAccessWithPasswordRequest) (*nativeIAmAuthGRPC.CheckAccessWithPasswordResponse, error) {
	//TODO: use provided metadata

	authenticateResponse, err := s.authenticationPasswordServer.Authenticate(ctx, &nativeIAmAuthenticationPasswordGRPC.AuthenticateRequest{
		Namespace: in.Namespace,
		Identity:  in.Identity,
		Password:  in.Password,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while authorizing using password. "+err.Error())
	}
	if !authenticateResponse.Authenticated {
		return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHENTICATED, Message: "Identity or password doesnt match."}, status.Error(codes.OK, "")
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
			if st.Code() == codes.NotFound {
				return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHENTICATED, Message: "Identity or password doesnt match."}, status.Error(codes.OK, "")
			}
		}

		return nil, status.Error(codes.Internal, "Error while searching for identity. "+err.Error())
	}

	if !identityGetResponse.Identity.Active {
		return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHORIZED, Message: "Identity is not active."}, status.Error(codes.OK, "")
	}

	//TODO: Cache this
	// Get all policies for identity
	policies, err := s.fetchIdentityPolicies(ctx, identityGetResponse.Identity)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, in.Scopes) {
		return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_UNAUTHORIZED, Message: "Not enought privileges"}, status.Error(codes.OK, "")
	}

	return &nativeIAmAuthGRPC.CheckAccessWithPasswordResponse{Status: nativeIAmAuthGRPC.CheckAccessWithPasswordResponse_OK, Message: ""}, status.Error(codes.OK, "")
}

func (s *IAmAuthServer) CheckAccessWithX509(ctx context.Context, in *nativeIAmAuthGRPC.CheckAccessWithX509Request) (*nativeIAmAuthGRPC.CheckAccessWithX509Response, error) {
	certificateValidateResponse, err := s.authenticationX509Server.ValidateAndGetFromRawX509(ctx, &x509.ValidateAndGetFromRawX509Request{
		Raw: in.Certificate,
	})
	if err != nil {
		st, ok := status.FromError(err)

		if ok {
			switch st.Code() {
			case codes.OK:
				break
			case codes.FailedPrecondition:
				return nil, err
			default:
				return nil, status.Error(codes.Internal, "error while validating certificate: "+err.Error())
			}
		} else {
			return nil, status.Error(codes.Internal, "error while validating certificate: "+err.Error())
		}
	}

	switch certificateValidateResponse.Status {
	case x509.ValidateAndGetFromRawX509Response_OK:
		break
	case x509.ValidateAndGetFromRawX509Response_INVALID_FORMAT:
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_CERTIFICATE_INVALID_FORMAT, Message: ""}, status.Error(codes.OK, "")
	case x509.ValidateAndGetFromRawX509Response_SIGNATURE_INVALID:
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_CERTIFICATE_INVALID, Message: "Invalid signature or other fields of the certificate"}, status.Error(codes.OK, "")
	case x509.ValidateAndGetFromRawX509Response_NOT_FOUND:
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_CERTIFICATE_NOT_FOUND, Message: "Certificate has valid signature, but doesnt exist in the system"}, status.Error(codes.OK, "")
	default:
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_CERTIFICATE_INVALID, Message: "Certificate invalid for unknown reason"}, status.Error(codes.OK, "")
	}

	certificateInfo := &nativeIAmAuthGRPC.CheckAccessWithX509Response_CertificateInfo{
		Namespace: certificateValidateResponse.Certificate.Namespace,
		Uuid:      certificateValidateResponse.Certificate.Uuid,
		Identity:  certificateValidateResponse.Certificate.Identity,
	}

	if certificateValidateResponse.Certificate.Disabled {
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{
			Status:          nativeIAmAuthGRPC.CheckAccessWithX509Response_CERTIFICATE_DISABLED,
			Message:         "Certificate was manually disabled",
			CertificateInfo: certificateInfo,
		}, status.Error(codes.OK, "")
	}

	// Find identity
	identityGetResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{
		Namespace: certificateInfo.Namespace,
		Uuid:      certificateInfo.Identity,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_IDENTITY_NOT_FOUND, Message: "Identity not found.", CertificateInfo: certificateInfo}, status.Error(codes.OK, "")
			}
		}

		return nil, status.Error(codes.Internal, "Error while searching for identity. "+err.Error())
	}

	if !identityGetResponse.Identity.Active {
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_IDENTITY_NOT_ACTIVE, Message: "Identity is not active.", CertificateInfo: certificateInfo}, status.Error(codes.OK, "")
	}

	//TODO: Cache this
	// Get all policies for identity
	policies, err := s.fetchIdentityPolicies(ctx, identityGetResponse.Identity)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, in.Scopes) {
		return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_UNAUTHORIZED, Message: "Not enought privileges", CertificateInfo: certificateInfo}, status.Error(codes.OK, "")
	}

	return &nativeIAmAuthGRPC.CheckAccessWithX509Response{Status: nativeIAmAuthGRPC.CheckAccessWithX509Response_OK, Message: "", CertificateInfo: certificateInfo}, status.Error(codes.OK, "")
}

func (s *IAmAuthServer) CheckAccess(ctx context.Context, in *nativeIAmAuthGRPC.CheckAccessRequest) (*nativeIAmAuthGRPC.CheckAccessResponse, error) {
	// Find identity
	identityGetResponse, err := s.identityServer.Get(ctx, &nativeIAmIdentityGRPC.GetIdentityRequest{
		Namespace: in.Namespace,
		Uuid:      in.Identity,
		UseCache:  false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return &nativeIAmAuthGRPC.CheckAccessResponse{Status: nativeIAmAuthGRPC.CheckAccessResponse_IDENTITY_NOT_FOUND, Message: "Identity not found."}, status.Error(codes.OK, "")
			}
		}

		return nil, status.Error(codes.Internal, "Error while searching for identity. "+err.Error())
	}

	if !identityGetResponse.Identity.Active {
		return &nativeIAmAuthGRPC.CheckAccessResponse{Status: nativeIAmAuthGRPC.CheckAccessResponse_IDENTITY_NOT_ACTIVE, Message: "Identity is not active."}, status.Error(codes.OK, "")
	}

	//TODO: Cache this
	// Get all policies for identity
	policies, err := s.fetchIdentityPolicies(ctx, identityGetResponse.Identity)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get policy information for identity: "+err.Error())
	}

	if !arePoliciesAllowScopes(policies, in.Scopes) {
		return &nativeIAmAuthGRPC.CheckAccessResponse{Status: nativeIAmAuthGRPC.CheckAccessResponse_UNAUTHORIZED, Message: "Not enought privileges"}, status.Error(codes.OK, "")
	}

	return &nativeIAmAuthGRPC.CheckAccessResponse{Status: nativeIAmAuthGRPC.CheckAccessResponse_OK, Message: ""}, status.Error(codes.OK, "")
}
