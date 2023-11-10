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

	clientGRPC "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/client"
	onecSyncGRPC "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/onecsync"
	performerGRPC "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/performer"
	settingsRGPC "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/settings"

	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/backend/onec/sync"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/services"
	"github.com/slamy-solutions/openbp/modules/crm/services/core/src/settings"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
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
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("crm", "core", VERSION, getHostname())).WithVault(),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}
	defer systemStub.Close(context.Background())

	// Connect to the "native" module services
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithKeyValueStorageService().WithIAMService())
	err = nativeStub.Connect()
	if err != nil {
		panic("Failed to connect to native services: " + err.Error())
	}
	defer nativeStub.Close()

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

	backendFactory := backend.NewBackendFactory(logger, systemStub, nativeStub)

	clientService := services.NewClientServer(backendFactory, logger.With(slog.String("service", "client")))
	clientGRPC.RegisterClientServiceServer(grpcServer, clientService)

	performerService := services.NewPerformerServer(backendFactory, logger.With(slog.String("service", "performer")))
	performerGRPC.RegisterPerformerServiceServer(grpcServer, performerService)

	settingsRepository := settings.NewSettingsRepository(systemStub, nativeStub)
	settingsService := services.NewSettingsServer(settingsRepository, logger.With(slog.String("service", "settings")))
	settingsRGPC.RegisterSettingsServiceServer(grpcServer, settingsService)

	onecSyncEngine := sync.NewSyncEngine(logger.With(slog.String("service", "onec_sync_engine")), systemStub, nativeStub, settingsRepository)
	onecSyncEngine.Start()
	defer onecSyncEngine.Stop()

	onecService := services.NewOneCSyncServer(onecSyncEngine, logger.With(slog.String("service", "onec_sync")))
	onecSyncGRPC.RegisterOneCSyncServiceServer(grpcServer, onecService)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
