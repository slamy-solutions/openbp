package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	native_iam_oauth_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
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
	telemetryProvider, err := otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "iam.oauth", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up native_iam_policy connection
	nativeIAmPolicyConnection, nativeIAmPolicyClient, err := native.NewIAMPolicyConnection(NATIVE_IAM_POLICY_URL)
	if err != nil {
		panic(err)
	}
	defer nativeIAmPolicyConnection.Close()
	fmt.Println("Initialized native_iam_policy connection")

	// Setting up native_iam_token connection
	nativeIAmTokenConnection, nativeIAmTokenClient, err := native.NewIAMTokenConnection(NATIVE_IAM_TOKEN_URL)
	if err != nil {
		panic(err)
	}
	defer nativeIAmTokenConnection.Close()
	fmt.Println("Initialized native_iam_token connection")

	// Setting up native_iam_authentication_password connection
	nativeIAmAuthenticationPasswordConnection, nativeIAmAuthenticationPasswordClient, err := native.NewIAMAuthenticationPasswordConnection(NATIVE_IAM_AUTHENTICATION_PASSWORD_URL)
	if err != nil {
		panic(err)
	}
	defer nativeIAmAuthenticationPasswordConnection.Close()
	fmt.Println("Initialized native_iam_authentication_password connection")

	// Setting up native_iam_identity connection
	nativeIAmIdentityConnection, nativeIAmIdentityClient, err := native.NewIAMIdentityConnection(NATIVE_IAM_IDENTITY_URL)
	if err != nil {
		panic(err)
	}
	defer nativeIAmIdentityConnection.Close()
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
