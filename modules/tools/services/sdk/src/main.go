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

	sdkTools "github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/tools"

	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	"github.com/slamy-solutions/openbp/modules/tools/services/sdk/src/servers"
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
	// --- Setting up connection to the system services
	systemStub := system.NewSystemStub(
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("tools", "sdk", VERSION, getHostname())),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}
	defer systemStub.Close(context.Background())

	// --- Setting up connection to the native services
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to the native services: " + err.Error())
	}
	defer nativeStub.Close()
	fmt.Println("Connected to the native services")

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Minute * 5,
			MaxConnectionAge:  time.Minute * 30,
		}),
	)

	authHandler, err := sdkTools.NewAuthHandler(nativeStub, sdkTools.NewTraefikCertificateExtractor())
	if err != nil {
		panic("Failed to initialize  auth handler: " + err.Error())
	}

	logger := logrus.StandardLogger()

	err = servers.RegisterGRPCServers(grpcServer, authHandler, logger.WithField("module", "server"), systemStub, nativeStub)
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
