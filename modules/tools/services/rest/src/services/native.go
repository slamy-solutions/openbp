package services

import (
	"context"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"google.golang.org/grpc"
)

type NativeConnectionHandler struct {
	Namespace namespace.NamespaceServiceClient

	IAMAuth                   auth.IAMAuthServiceClient
	IAMPolicy                 policy.IAMPolicyServiceClient
	IAMIdentity               identity.IAMIdentityServiceClient
	IAMAuthenticationPassword password.IAMAuthenticationPasswordServiceClient
	ActorUser                 user.ActorUserServiceClient
	KeyValueStorage           keyvaluestorage.KeyValueStorageServiceClient

	dials []*grpc.ClientConn
}

func ConnectToNativeServices(ctx context.Context) (*NativeConnectionHandler, error) {
	NATIVE_NAMESPACE_URL := getConfigEnv("NATIVE_NAMESPACE_URL", "native_namespace:80")
	NATIVE_KEYVALUESTORAGE_URL := getConfigEnv("NATIVE_KEYVALUESTORAGE_URL", "native_keyvaluestorage:80")

	NATIVE_ACTOR_USER_URL := getConfigEnv("NATIVE_ACTOR_USER_URL", "native_actor_user:80")

	NATIVE_IAM_AUTH_URL := getConfigEnv("NATIVE_IAM_AUTH_URL", "native_iam_auth:80")
	NATIVE_IAM_POLICY_URL := getConfigEnv("NATIVE_IAM_POLICY_URL", "native_iam_policy:80")
	NATIVE_IAM_IDENTITY_URL := getConfigEnv("NATIVE_IAM_IDENTITY_URL", "native_iam_identity:80")
	NATIVE_IAM_AUTHENTICATION_PASSWORD_URL := getConfigEnv("NATIVE_IAM_AUTHENTICATION_PASSWORD_URL", "native_iam_authentication_password:80")

	dials := make([]*grpc.ClientConn, 1)
	clear := func() {
		for _, dial := range dials {
			dial.Close()
		}
	}

	namespaceDial, namespaceConnection, err := native.NewNamespaceConnection(NATIVE_NAMESPACE_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, namespaceDial)

	keyvaluestorageDial, keyvaluestorageConnection, err := native.NewKeyValueStorageConnection(NATIVE_KEYVALUESTORAGE_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, keyvaluestorageDial)

	actorUserDial, actorUserConnection, err := native.NewActorUserConnection(NATIVE_ACTOR_USER_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, actorUserDial)

	iamAuthDial, iamAuthConnection, err := native.NewIAMAuthConnection(NATIVE_IAM_AUTH_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, iamAuthDial)

	iamPolicyDial, iamPolicyConnection, err := native.NewIAMPolicyConnection(NATIVE_IAM_POLICY_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, iamPolicyDial)

	iamIdentityDial, iamIdentityConnection, err := native.NewIAMIdentityConnection(NATIVE_IAM_IDENTITY_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, iamIdentityDial)

	iamAuthenticationPasswordDial, iamAuthenticationPasswordConnection, err := native.NewIAMAuthenticationPasswordConnection(NATIVE_IAM_AUTHENTICATION_PASSWORD_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, iamAuthenticationPasswordDial)

	return &NativeConnectionHandler{
		Namespace:                 namespaceConnection,
		KeyValueStorage:           keyvaluestorageConnection,
		IAMAuth:                   iamAuthConnection,
		IAMPolicy:                 iamPolicyConnection,
		IAMIdentity:               iamIdentityConnection,
		IAMAuthenticationPassword: iamAuthenticationPasswordConnection,
		ActorUser:                 actorUserConnection,
		dials:                     dials,
	}, nil
}

func (h *NativeConnectionHandler) Shutdown(ctx context.Context) {
	for _, dial := range h.dials {
		dial.Close()
	}
}
