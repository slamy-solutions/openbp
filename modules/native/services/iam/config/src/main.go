package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"

	native_iam_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/config"
	native_keyvaluestorage_grpc "github.com/slamy-solutions/openbp/modules/native/services/iam/config/src/grpc/native_keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/config/src/services"
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
	SYSTEM_CACHE_URL := getConfigEnv("SYSTEM_CACHE_URL", "redis://system_cache")

	NATIVE_KEYVALUESTORAGE_URL := getConfigEnv("NATIVE_KEYVALUESTORAGE_URL", "native_keyvaluestorage:80")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "iam.config", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up Cache
	cacheClient, err := cache.New(SYSTEM_CACHE_URL)
	if err != nil {
		panic(err)
	}
	defer cacheClient.Shutdown(ctx)
	fmt.Println("Initialized cache")

	// Setting up native_keyvaluestorage connection
	nativeKeyValueStorageConnection, err := grpc.Dial(
		NATIVE_KEYVALUESTORAGE_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeKeyValueStorageConnection.Close()
	nativeKeyValueStorageClient := native_keyvaluestorage_grpc.NewKeyValueStorageServiceClient(nativeKeyValueStorageConnection)
	fmt.Println("Initialized native_keyvaluestorage connection")

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	iamConfigServer := services.NewIAMConfigServer(cacheClient, nativeKeyValueStorageClient)
	native_iam_grpc.RegisterIAMConfigServiceServer(grpcServer, iamConfigServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
