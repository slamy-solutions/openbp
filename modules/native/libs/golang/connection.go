package native

import (
	"errors"
	"time"

	actorUser "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	iamAuth "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	iamAuthenticationPassword "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	iamIdentity "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	iamPolicy "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	iamRole "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	iamToken "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func makeGrpcDial(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(
		[]grpc.DialOption{
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
			grpc.WithBlock(),
			grpc.WithTimeout(time.Second * 5),
		},
		opts...,
	)

	dial, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, errors.New("failed to establish grpc dial: " + err.Error())
	}

	return dial, nil
}

func makeGrpcClient[T interface{}](clientFunction func(grpc.ClientConnInterface) T, address string, opts ...grpc.DialOption) (*grpc.ClientConn, T, error) {
	dial, err := makeGrpcDial(address, opts...)
	if err != nil {
		var result T
		return nil, result, err
	}

	client := clientFunction(dial)
	return dial, client, nil
}

// Connect to Namespace service
func NewNamespaceConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, namespace.NamespaceServiceClient, error) {
	return makeGrpcClient(namespace.NewNamespaceServiceClient, address, opts...)
}

// Connect to KeyValueStorage service
func NewKeyValueStorageConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, keyvaluestorage.KeyValueStorageServiceClient, error) {
	return makeGrpcClient(keyvaluestorage.NewKeyValueStorageServiceClient, address, opts...)
}

// Connect to Actor_User service
func NewActorUserConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, actorUser.ActorUserServiceClient, error) {
	return makeGrpcClient(actorUser.NewActorUserServiceClient, address, opts...)
}

// Connect to IAM service
func NewIAMConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, *IAMService, error) {
	dial, err := makeGrpcDial(address, opts...)
	if err != nil {
		return nil, nil, errors.New("failed to connect to service: " + err.Error())
	}

	return dial, &IAMService{
		Authentication: &IamAuthenticationServices{
			Password: iamAuthenticationPassword.NewIAMAuthenticationPasswordServiceClient(dial),
		},
		Identity: iamIdentity.NewIAMIdentityServiceClient(dial),
		Auth:     iamAuth.NewIAMAuthServiceClient(dial),
		Policy:   iamPolicy.NewIAMPolicyServiceClient(dial),
		Role:     iamRole.NewIAMRoleServiceClient(dial),
		Token:    iamToken.NewIAMTokenServiceClient(dial),
	}, nil
}
