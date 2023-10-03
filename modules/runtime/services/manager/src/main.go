package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	rpcRPC "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/rpc"
	runtimeGRPC "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/runtime"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	rpcServer "github.com/slamy-solutions/openbp/modules/runtime/services/manager/src/services/rpc"
	runtimeServer "github.com/slamy-solutions/openbp/modules/runtime/services/manager/src/services/runtime"
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
	// Connect to the "system" module services
	systemStub := system.NewSystemStub(
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("runtime", "manager", VERSION, getHostname())).WithVault(),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}
	defer systemStub.Close(context.Background())

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

	logger := slog.Default()

	runtime := runtimeServer.NewManagerRuntimeServer(logger.With(slog.String("server", "runtime")), systemStub)
	runtimeGRPC.RegisterRuntimeServiceServer(grpcServer, runtime)

	rpc := rpcServer.NewManagerRPCServer(logger.With(slog.String("server", "rpc")), systemStub)
	rpcRPC.RegisterRPCServiceServer(grpcServer, rpc)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
