package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	native_actor_user_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/services/actor/user/src/services"
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
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("native", "actor_user", VERSION, getHostname())),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}

	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMIdentityService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to native services: " + err.Error())
	}
	defer nativeStub.Close()

	ctx := context.Background()

	eventHandler, err := services.NewEventHandlerService(systemStub)
	if err != nil {
		panic("Failed to initialize event handler. " + err.Error())
	}
	defer eventHandler.Close()

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	actorUserServer, err := services.NewActorUserServer(ctx, systemStub, nativeStub)
	if err != nil {
		panic("Failed to initialize user server. " + err.Error())
	}
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
