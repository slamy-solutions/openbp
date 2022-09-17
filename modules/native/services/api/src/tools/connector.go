package tools

import (
	"fmt"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/system/libs/go/cache"

	actorUserGRPC "github.com/slamy-solutions/openbp/modules/native/services/api/src/grpc/native_actor_user"
	iamAuthGRPC "github.com/slamy-solutions/openbp/modules/native/services/api/src/grpc/native_iam_auth"
	iamAuthenticationPasswordGRPC "github.com/slamy-solutions/openbp/modules/native/services/api/src/grpc/native_iam_authentication_password"
	keyValueStorageGRPC "github.com/slamy-solutions/openbp/modules/native/services/api/src/grpc/native_keyvaluestorage"
	namespaceGRPC "github.com/slamy-solutions/openbp/modules/native/services/api/src/grpc/native_namespace"
)

type ConnectorTools struct {
	Cache                     cache.Cache
	Namespace                 namespaceGRPC.NamespaceServiceClient
	KeyValueStorage           keyValueStorageGRPC.KeyValueStorageServiceClient
	ActorUser                 actorUserGRPC.ActorUserServiceClient
	IAmAuth                   iamAuthGRPC.IAMAuthServiceClient
	IAmAuthenticationPassword iamAuthenticationPasswordGRPC.IAMAuthenticationPasswordServiceClient

	dials []*grpc.ClientConn
}

func NewConnectorTools() (*ConnectorTools, error) {
	tools := &ConnectorTools{}

	fmt.Println("Initializing native_namespace connection")
	namespaceDial, err := connectToGRPC("native_namespace:80")
	if err != nil {
		tools.Close()
		return nil, err
	}
	tools.dials = append(tools.dials, namespaceDial)
	tools.Namespace = namespaceGRPC.NewNamespaceServiceClient(namespaceDial)

	fmt.Println("Initializing native_keyvaluestorage connection")
	keyvaluestorageDial, err := connectToGRPC("native_keyvaluestorage:80")
	if err != nil {
		tools.Close()
		return nil, err
	}
	tools.dials = append(tools.dials, keyvaluestorageDial)
	tools.KeyValueStorage = keyValueStorageGRPC.NewKeyValueStorageServiceClient(keyvaluestorageDial)

	fmt.Println("Initializing native_actor_user connection")
	actorUserDial, err := connectToGRPC("native_actor_user:80")
	if err != nil {
		tools.Close()
		return nil, err
	}
	tools.dials = append(tools.dials, actorUserDial)
	tools.ActorUser = actorUserGRPC.NewActorUserServiceClient(actorUserDial)

	fmt.Println("Initializing native_iam_auth connection")
	iamAuthDial, err := connectToGRPC("native_iam_auth:80")
	if err != nil {
		tools.Close()
		return nil, err
	}
	tools.dials = append(tools.dials, iamAuthDial)
	tools.IAmAuth = iamAuthGRPC.NewIAMAuthServiceClient(iamAuthDial)

	fmt.Println("Initializing native_iam_authentication_password connection")
	iamAuthenticationPasswordDial, err := connectToGRPC("native_iam_authentication_password:80")
	if err != nil {
		tools.Close()
		return nil, err
	}
	tools.dials = append(tools.dials, iamAuthenticationPasswordDial)
	tools.IAmAuthenticationPassword = iamAuthenticationPasswordGRPC.NewIAMAuthenticationPasswordServiceClient(iamAuthenticationPasswordDial)

	return tools, nil
}

func (t *ConnectorTools) Close() {
	for _, dial := range t.dials {
		dial.Close()
	}
}

func connectToGRPC(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
}
