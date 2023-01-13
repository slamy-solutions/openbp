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

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/db"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	native_iam_identity_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/identity/src/services"
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
	SYSTEM_DB_URL := getConfigEnv("SYSTEM_DB_URL", "mongodb://root:example@system_db/admin")
	SYSTEM_TELEMETRY_EXPORTER_ENDPOINT := getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680")
	SYSTEM_CACHE_URL := getConfigEnv("SYSTEM_CACHE_URL", "redis://system_cache")

	NATIVE_NAMESPACE_URL := getConfigEnv("NATIVE_NAMESPACE_URL", "native_namespace:80")
	NATIVE_IAM_POLICY_URL := getConfigEnv("NATIVE_IAM_POLICY_URL", "native_iam_policy:80")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "iam.identity", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up DB
	dbClient, err := db.Connect(SYSTEM_DB_URL)
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect(ctx)
	fmt.Println("Initialized DB")

	// Setting up Cache
	cacheClient, err := cache.New(SYSTEM_CACHE_URL)
	if err != nil {
		panic(err)
	}
	defer cacheClient.Shutdown(ctx)
	fmt.Println("Initialized cache")

	// Setting up native_namespace connection
	nativeNamespaceConnection, nativeNamespaceClient, err := native.NewNamespaceConnection(NATIVE_NAMESPACE_URL)
	if err != nil {
		panic(err)
	}
	defer nativeNamespaceConnection.Close()
	fmt.Println("Initialized native_namespace connection")

	// Setting up native_iam_policy connection
	nativeIAmPolicyConnection, nativeIAmPolicyClient, err := native.NewIAMPolicyConnection(NATIVE_IAM_POLICY_URL)
	if err != nil {
		panic(err)
	}
	defer nativeIAmPolicyConnection.Close()
	fmt.Println("Initialized native_iam_policy connection")

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	iamIdentityServer := services.NewIAmIdentityServer(dbClient, cacheClient, nativeNamespaceClient, nativeIAmPolicyClient)
	native_iam_identity_grpc.RegisterIAMIdentityServiceServer(grpcServer, iamIdentityServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
