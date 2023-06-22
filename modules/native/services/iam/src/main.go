package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	native_iam_auth_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	native_iam_authentication_password_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	native_iam_identity_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	native_iam_policy_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	native_iam_role_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	native_iam_token_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/token"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/auth"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/authentication"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/identity"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/policy"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/role"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/src/services/token"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const (
	VERSION = "1.0.0"
)

func getHostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}

func main() {
	systemStub := system.NewSystemStub(
		system.NewSystemStubConfig().
			WithOTel(system.NewOTelConfig("native", "iam", VERSION, getHostname())).
			WithCache().
			WithDB().
			WithNats().
			WithVault(),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}

	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to native services: " + err.Error())
	}
	defer nativeStub.Close()

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	policyServer, err := policy.NewIAMPolicyServer(context.Background(), systemStub, nativeStub)
	if err != nil {
		panic("Failed to startup policy server: " + err.Error())
	}
	native_iam_policy_grpc.RegisterIAMPolicyServiceServer(grpcServer, policyServer)

	roleServer, err := role.NewIAMRoleServer(context.Background(), systemStub, nativeStub, policyServer)
	if err != nil {
		panic("Failed to startup role server: " + err.Error())
	}
	native_iam_role_grpc.RegisterIAMRoleServiceServer(grpcServer, roleServer)

	identityServer, err := identity.NewIAmIdentityServer(context.Background(), systemStub, nativeStub, policyServer, roleServer)
	if err != nil {
		panic("Failed to startup identity server: " + err.Error())
	}
	native_iam_identity_grpc.RegisterIAMIdentityServiceServer(grpcServer, identityServer)

	tokenServer := token.NewIAmTokenServer(systemStub, nativeStub)
	native_iam_token_grpc.RegisterIAMTokenServiceServer(grpcServer, tokenServer)

	authenticationPasswordServer := authentication.NewPasswordIdentificationService(systemStub, nativeStub)
	native_iam_authentication_password_grpc.RegisterIAMAuthenticationPasswordServiceServer(grpcServer, authenticationPasswordServer)

	iamAuthServer := auth.NewIAmAuthServer(systemStub, authenticationPasswordServer, identityServer, policyServer, roleServer, tokenServer)
	native_iam_auth_grpc.RegisterIAMAuthServiceServer(grpcServer, iamAuthServer)

	eventHandler, err := NewEventHandlerService(systemStub, nativeStub, policyServer, roleServer)
	if err != nil {
		panic("failed to setup event hanle service: " + err.Error())
	}
	defer eventHandler.Close()

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
