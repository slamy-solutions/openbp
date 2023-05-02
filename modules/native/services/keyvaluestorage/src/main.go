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
	native_keyvaluestorage_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/native/services/keyvaluestorage/src/services"
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
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("native", "keyvaluestorage", VERSION, getHostname())),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}
	defer systemStub.Close(context.Background())

	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to native services: " + err.Error())
	}
	defer nativeStub.Close()

	// Listen for events from other services
	eventHandler, err := services.NewEventHandlerService(systemStub.DB, systemStub.Nats)
	if err != nil {
		panic("Failed to setup event handler: " + err.Error())
	}
	defer eventHandler.Close()

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.MaxRecvMsgSize(1024*1024*16), // 16 megabytes
		grpc.MaxSendMsgSize(1024*1024*16), // 16 megabytes
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	storageServer := services.NewKeyValueStorageServer(systemStub.DB, systemStub.Cache, nativeStub.Services.Namespace)
	native_keyvaluestorage_grpc.RegisterKeyValueStorageServiceServer(grpcServer, storageServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
