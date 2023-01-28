package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	otel "github.com/slamy-solutions/openbp/modules/system/libs/golang/otel"

	tools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"

	"github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/services"
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

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := otel.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "namespace", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up connection to the native services
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to the native services: " + err.Error())
	}
	defer nativeStub.Close()
	fmt.Println("Connected to the native services")

	modulesStub := &tools.ModulesStub{
		Native: nativeStub,
	}

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Minute * 5,
			MaxConnectionAge:  time.Minute * 30,
		}),
	)

	err = services.RegisterGRPCServices(modulesStub, grpcServer)
	if err != nil {
		panic("Failed register gRPC services: " + err.Error())
	}

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic("Failed to setup port: " + err.Error())
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic("Error on serving grpc: " + err.Error())
	}
}
