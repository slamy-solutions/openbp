package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/system/libs/golang/cache"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/db"
	"github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	native_actor_user_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/services/actor/user/src/services"
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

	NATIVE_IAM_IDENTITY_URL := getConfigEnv("NATIVE_IAM_IDENTITY_URL", "native_iam_identity:80")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "actor.user", VERSION, "1")
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

	actorUserServer, err := services.NewActorUserServer(ctx, dbClient, cacheClient, nativeIAmIdentityClient)
	native_actor_user_grpc.RegisterActorUserServiceServer(grpcServer, actorUserServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
