package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/nats-io/nats.go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/db"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"

	native_namespace_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/namespace"
	"github.com/slamy-solutions/openbp/modules/native/services/namespace/src/services"
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
	SYSTEM_CACHE_URL := getConfigEnv("SYSTEM_CACHE_URL", "redis://system_cache")
	SYSTEM_NATS_URL := getConfigEnv("SYSTEM_NATS_URL", "nats://system_nats:4222")
	SYSTEM_TELEMETRY_EXPORTER_ENDPOINT := getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "namespace", VERSION, "1")
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

	// Setting up DB
	dbClient, err := db.Connect(SYSTEM_DB_URL)
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect(ctx)
	fmt.Println("Initialized DB")

	// Setting up Nats
	nc, err := nats.Connect(SYSTEM_NATS_URL)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	namespaceServer, err := services.New(dbClient, cacheClient, js)
	if err != nil {
		panic(err)
	}
	native_namespace_grpc.RegisterNamespaceServiceServer(grpcServer, namespaceServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
