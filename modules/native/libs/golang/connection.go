package native

import (
	"time"

	actorUser "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	iamAuth "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	iamAuthenticationPassword "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	iamConfig "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/config"
	iamIdentity "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	iamPolicy "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	iamRole "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	iamToken "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func makeGrpcClient[T interface{}](clientFunction func(grpc.ClientConnInterface) T, address string, opts ...grpc.DialOption) (*grpc.ClientConn, T, error) {
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

// Connect to IAM_Token service
func NewIAMTokenConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamToken.IAMTokenServiceClient, error) {
	return makeGrpcClient(iamToken.NewIAMTokenServiceClient, address, opts...)
}

// Connect to IAM_Policy service
func NewIAMPolicyConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamPolicy.IAMPolicyServiceClient, error) {
	return makeGrpcClient(iamPolicy.NewIAMPolicyServiceClient, address, opts...)
}

// Connect to IAM_Role service
func NewIAMRoleConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamRole.IAMRoleServiceClient, error) {
	return makeGrpcClient(iamRole.NewIAMRoleServiceClient, address, opts...)
}

// Connect to IAM_Auth service
func NewIAMAuthConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamAuth.IAMAuthServiceClient, error) {
	return makeGrpcClient(iamAuth.NewIAMAuthServiceClient, address, opts...)
}

// Connect to IAM_Identity service
func NewIAMIdentityConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamIdentity.IAMIdentityServiceClient, error) {
	return makeGrpcClient(iamIdentity.NewIAMIdentityServiceClient, address, opts...)
}

// Connect to IAM_Config service
func NewIAMConfigConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamConfig.IAMConfigServiceClient, error) {
	return makeGrpcClient(iamConfig.NewIAMConfigServiceClient, address, opts...)
}

// Connect to IAM_Authentication_Password service
func NewIAMAuthenticationPasswordConnection(address string, opts ...grpc.DialOption) (*grpc.ClientConn, iamAuthenticationPassword.IAMAuthenticationPasswordServiceClient, error) {
	return makeGrpcClient(iamAuthenticationPassword.NewIAMAuthenticationPasswordServiceClient, address, opts...)
}
