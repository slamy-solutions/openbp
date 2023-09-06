package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/slamy-solutions/openbp/modules/erp/services/core/src/services/catalog"
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
	// Connect to the "system" module services
	systemStub := system.NewSystemStub(
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("erp", "core", VERSION, getHostname())).WithVault(),
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

	catalog.RegisterCatalogModule(logger.With(slog.String("module", "catalog")), grpcServer, systemStub)

	/*eventHandler, err := NewEventHandler(logger.WithField("service", "event_handler"), systemStub)
	if err != nil {
		panic("Failed to start event handler: " + err.Error())
	}
	defer eventHandler.Close()*/

	logger.Info("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
