package tools

import (
	"context"
	"crypto/x509"
	"errors"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotAuthenticated = errors.New("not authenticated")
var ErrNotAuthorized = errors.New("not authorized")
var ErrVaultIsSealed = errors.New("the vault is sealed")

type Scope struct {
	Namespace string
	Resources []string
	Actions   []string
}

type AuthInfo struct {
	Namespace string
	Identity  string
}

type AuthorizedContext interface {
	Authorize(ctx context.Context, scopes *[]Scope) error
	GetAuthInfo() AuthInfo
}

type AuthHandler interface {
	/*
		If you want to authorize multiple times this context will reduce the number of API call and theirs complexity.

		In the first call it will validate the certificate and raise error if it is not valid / doesnt exist / disabled ...
		In consequitive calls it will only check to `Authorize` it will only check the scopes of the identity without validating the certificate itself.
		Make sure to use it for short period of times, because in the meantime certificate may become not valid.
	*/
	StartAuthorizedContext(ctx context.Context) (AuthorizedContext, error)

	// Validates the certificate and checks identity for scopes
	Authorize(ctx context.Context, scopes *[]Scope) (*AuthInfo, error)
}

type CertificateExtractor interface {
	Exctract(ctx context.Context) (*x509.Certificate, error)
	/*Start() error
	Close() error*/
}

type authHandler struct {
	nativeStub *native.NativeStub
	extractor  CertificateExtractor
}

type authorizedContext struct {
	authHandler

	authInfo AuthInfo
}

// Create authorization instance based on communication with Traefik reverse proxy.
func NewAuthHandler(nativeStub *native.NativeStub, extractor CertificateExtractor) (AuthHandler, error) {
	return &authHandler{
		nativeStub: nativeStub,
		extractor:  extractor,
	}, nil
}

// Converts handler error to GRPC status. If error is internal, the message will be empty.
func AuthHandlerErrorToGRPCStatus(err error) *status.Status {
	if err == ErrNotAuthorized {
		return status.New(codes.PermissionDenied, "auth error: "+err.Error())
	}
	if err == ErrNotAuthenticated {
		return status.New(codes.Unauthenticated, "auth error: "+err.Error())
	}
	if err == ErrVaultIsSealed {
		return status.New(codes.FailedPrecondition, "auth error: "+err.Error())
	}

	return status.New(codes.Internal, "")
}

func (a *authHandler) StartAuthorizedContext(ctx context.Context) (AuthorizedContext, error) {
	cert, err := a.extractor.Exctract(ctx)
	if err != nil {
		return nil, errors.New("failed to extract certificate: " + err.Error())
	}

	checkResponse, err := a.nativeStub.Services.IAM.Auth.CheckAccessWithX509(ctx, &auth.CheckAccessWithX509Request{
		Certificate: cert.Raw,
		Scopes:      []*auth.Scope{},
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return nil, ErrVaultIsSealed
			}
		}

		return nil, errors.New("failed to check access with X509 certificate: " + err.Error())
	}

	switch checkResponse.Status {
	case auth.CheckAccessWithX509Response_OK:
		return &authorizedContext{
			authHandler: *a,
			authInfo: AuthInfo{
				Namespace: checkResponse.CertificateInfo.Namespace,
				Identity:  checkResponse.CertificateInfo.Identity,
			},
		}, nil
	case auth.CheckAccessWithX509Response_UNAUTHORIZED:
		return nil, ErrNotAuthorized
	default:
		return nil, ErrNotAuthenticated
	}
}

func (a *authHandler) Authorize(ctx context.Context, scopes *[]Scope) (*AuthInfo, error) {
	cert, err := a.extractor.Exctract(ctx)
	if err != nil {
		return nil, errors.New("failed to extract certificate: " + err.Error())
	}

	requestedScopes := make([]*auth.Scope, 0, len(*scopes))
	for _, scope := range *scopes {
		requestedScopes = append(requestedScopes, &auth.Scope{
			Namespace: scope.Namespace,
			Resources: scope.Resources,
			Actions:   scope.Actions,
		})
	}

	checkResponse, err := a.nativeStub.Services.IAM.Auth.CheckAccessWithX509(ctx, &auth.CheckAccessWithX509Request{
		Certificate: cert.Raw,
		Scopes:      requestedScopes,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return nil, ErrVaultIsSealed
			}
		}

		return nil, errors.New("failed to check access with X509 certificate: " + err.Error())
	}

	switch checkResponse.Status {
	case auth.CheckAccessWithX509Response_OK:
		return &AuthInfo{Namespace: checkResponse.CertificateInfo.Namespace, Identity: checkResponse.CertificateInfo.Identity}, nil
	case auth.CheckAccessWithX509Response_UNAUTHORIZED:
		return &AuthInfo{Namespace: checkResponse.CertificateInfo.Namespace, Identity: checkResponse.CertificateInfo.Identity}, ErrNotAuthorized
	default:
		return nil, ErrNotAuthenticated
	}
}

func (a *authorizedContext) Authorize(ctx context.Context, scopes *[]Scope) error {
	requestedScopes := make([]*auth.Scope, 0, len(*scopes))
	for _, scope := range *scopes {
		requestedScopes = append(requestedScopes, &auth.Scope{
			Namespace: scope.Namespace,
			Resources: scope.Resources,
			Actions:   scope.Actions,
		})
	}

	checkResponse, err := a.nativeStub.Services.IAM.Auth.CheckAccess(ctx, &auth.CheckAccessRequest{
		Namespace: a.authInfo.Namespace,
		Identity:  a.authInfo.Identity,
		Scopes:    requestedScopes,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				return ErrVaultIsSealed
			}
		}

		return errors.New("failed to check access with X509 certificate: " + err.Error())
	}

	switch checkResponse.Status {
	case auth.CheckAccessResponse_OK:
		return nil
	case auth.CheckAccessResponse_UNAUTHORIZED:
		return ErrNotAuthorized
	default:
		return ErrNotAuthenticated
	}
}
func (a *authorizedContext) GetAuthInfo() AuthInfo {
	return AuthInfo{
		Namespace: "",
		Identity:  "",
	}
}
