package services

import (
	"context"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"google.golang.org/grpc"
)

type NativeConnectionHandler struct {
	Namespace namespace.NamespaceServiceClient

	IAMOAuth  oauth.IAMOAuthServiceClient
	ActorUser user.ActorUserServiceClient

	dials []*grpc.ClientConn
}

func ConnectToNativeServices(ctx context.Context) (*NativeConnectionHandler, error) {
	NATIVE_NAMESPACE_URL := getConfigEnv("NATIVE_NAMESPACE_URL", "native_namespace:80")

	NATIVE_ACTOR_USER_URL := getConfigEnv("NATIVE_ACTOR_USER_URL", "native_actor_user:80")

	NATIVE_IAM_OAUTH_URL := getConfigEnv("NATIVE_IAM_OAUTH_URL", "native_iam_oauth:80")

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

	actorUserDial, actorUserConnection, err := native.NewActorUserConnection(NATIVE_ACTOR_USER_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, actorUserDial)

	iamOAuthDial, iamOAuthConnection, err := native.NewIAMOAuthConnection(NATIVE_IAM_OAUTH_URL)
	if err != nil {
		clear()
		return nil, err
	}
	dials = append(dials, iamOAuthDial)

	return &NativeConnectionHandler{
		Namespace: namespaceConnection,
		IAMOAuth:  iamOAuthConnection,
		ActorUser: actorUserConnection,
		dials:     dials,
	}, nil
}

func (h *NativeConnectionHandler) Shutdown(ctx context.Context) {
	for _, dial := range h.dials {
		dial.Close()
	}
}
