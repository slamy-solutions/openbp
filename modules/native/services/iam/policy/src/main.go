package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/open-erp/modules/system/libs/go/cache"
	"github.com/slamy-solutions/open-erp/modules/system/libs/go/mongodb"
	"github.com/slamy-solutions/open-erp/modules/system/libs/go/telemetry"

	native_iam_policy_grpc "github.com/slamy-solutions/open-erp/modules/native/services/iam/policy/src/grpc/native_iam_policy"
	native_namespace_grpc "github.com/slamy-solutions/open-erp/modules/native/services/iam/policy/src/grpc/native_namespace"
	"github.com/slamy-solutions/open-erp/modules/native/services/iam/policy/src/services"
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
	SYSTEM_DB_PREFIX := getConfigEnv("SYSTEM_DB_PREFIX", "openerp_")
	SYSTEM_TELEMETRY_EXPORTER_ENDPOINT := getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680")
	SYSTEM_CACHE_URL := getConfigEnv("SYSTEM_CACHE_URL", "redis://system_cache")

	NATIVE_NAMESPACE_URL := getConfigEnv("NATIVE_NAMESPACE_URL", "native_namespace:80")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := telemetry.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "iam.policy", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up DB
	dbClient, err := mongodb.Connect(SYSTEM_DB_URL)
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
	nativeNamespaceConnection, err := grpc.Dial(
		NATIVE_NAMESPACE_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeNamespaceConnection.Close()
	nativeNamespaceClient := native_namespace_grpc.NewNamespaceServiceClient(nativeNamespaceConnection)
	fmt.Println("Initialized native_namespace connection")

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	iamPolicyServer := services.NewIAMPolicyServer(dbClient, SYSTEM_DB_PREFIX, cacheClient, nativeNamespaceClient)
	native_iam_policy_grpc.RegisterIAMPolicyServiceServer(grpcServer, iamPolicyServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
