package main

import (
	"context"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	native_iam_role_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	services "github.com/slamy-solutions/openbp/modules/native/services/iam/role/src/services"
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
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("native", "iam_role", VERSION, getHostname())),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}

	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMPolicyService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to native services: " + err.Error())
	}
	defer nativeStub.Close()

	eventsHandler, err := services.NewEventHandlerService(systemStub, nativeStub)
	if err != nil {
		panic("Failed to create events handler. " + err.Error())
	}
	defer eventsHandler.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	server, err := services.NewIAMRoleServer(context.Background(), systemStub, nativeStub)
	if err != nil {
		panic("Failed to create GRPC server: " + err.Error())
	}
	native_iam_role_grpc.RegisterIAMRoleServiceServer(grpcServer, server)

	log.Info("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
