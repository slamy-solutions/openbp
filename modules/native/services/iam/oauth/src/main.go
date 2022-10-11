package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/system/libs/go/telemetry"

	native_iam_authentication_password_grpc "github.com/slamy-solutions/openbp/modules/native/services/iam/oauth/src/grpc/native_iam_authentication_password"
	native_iam_identity_grpc "github.com/slamy-solutions/openbp/modules/native/services/iam/oauth/src/grpc/native_iam_identity"
	native_iam_oauth_grpc "github.com/slamy-solutions/openbp/modules/native/services/iam/oauth/src/grpc/native_iam_oauth"
	native_iam_policy_grpc "github.com/slamy-solutions/openbp/modules/native/services/iam/oauth/src/grpc/native_iam_policy"
	native_iam_token_grpc "github.com/slamy-solutions/openbp/modules/native/services/iam/oauth/src/grpc/native_iam_token"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/oauth/src/services"
)

const (
	VERSION = "1.0.0"
)

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	SYSTEM_TELEMETRY_EXPORTER_ENDPOINT := getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680")

	NATIVE_IAM_POLICY_URL := getConfigEnv("NATIVE_IAM_POLICY_URL", "native_iam_policy:80")
	NATIVE_IAM_TOKEN_URL := getConfigEnv("NATIVE_IAM_TOKEN_URL", "native_iam_token:80")
	NATIVE_IAM_IDENTITY_URL := getConfigEnv("NATIVE_IAM_IDENTITY_URL", "native_iam_identity:80")
	NATIVE_IAM_AUTHENTICATION_PASSWORD_URL := getConfigEnv("NATIVE_IAM_AUTHENTICATION_PASSWORD_URL", "native_iam_authentication_password:80")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := telemetry.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "iam.oauth", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up native_iam_policy connection
	nativeIAmPolicyConnection, err := grpc.Dial(
		NATIVE_IAM_POLICY_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeIAmPolicyConnection.Close()
	nativeIAmPolicyClient := native_iam_policy_grpc.NewIAMPolicyServiceClient(nativeIAmPolicyConnection)
	fmt.Println("Initialized native_iam_policy connection")

	// Setting up native_iam_token connection
	nativeIAmTokenConnection, err := grpc.Dial(
		NATIVE_IAM_TOKEN_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeIAmTokenConnection.Close()
	nativeIAmTokenClient := native_iam_token_grpc.NewIAMTokenServiceClient(nativeIAmTokenConnection)
	fmt.Println("Initialized native_iam_token connection")

	// Setting up native_iam_authentication_password connection
	nativeIAmAuthenticationPasswordConnection, err := grpc.Dial(
		NATIVE_IAM_AUTHENTICATION_PASSWORD_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeIAmAuthenticationPasswordConnection.Close()
	nativeIAmAuthenticationPasswordClient := native_iam_authentication_password_grpc.NewIAMAuthenticationPasswordServiceClient(nativeIAmAuthenticationPasswordConnection)
	fmt.Println("Initialized native_iam_authentication_password connection")

	// Setting up native_iam_identity connection
	nativeIAmIdentityConnection, err := grpc.Dial(
		NATIVE_IAM_IDENTITY_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeIAmIdentityConnection.Close()
	nativeIAmIdentityClient := native_iam_identity_grpc.NewIAMIdentityServiceClient(nativeIAmIdentityConnection)
	fmt.Println("Initialized native_iam_identity connection")

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	iamOAuthServer := services.NewIAmAuthServer(
		nativeIAmPolicyClient,
		nativeIAmIdentityClient,
		nativeIAmTokenClient,
		nativeIAmAuthenticationPasswordClient,
	)
	native_iam_oauth_grpc.RegisterIAMOAuthServiceServer(grpcServer, iamOAuthServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
