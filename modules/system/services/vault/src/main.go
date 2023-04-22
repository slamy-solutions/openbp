package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/slamy-solutions/openbp/modules/system/services/vault/src/pkcs"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	system_vault_grpc "github.com/slamy-solutions/openbp/modules/system/libs/golang/vault"
	"github.com/slamy-solutions/openbp/modules/system/services/vault/src/service"

	log "github.com/sirupsen/logrus"
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
	//Registering in the OTEL
	systemStub := system.NewSystemStub(system.NewSystemStubConfig().WithOTel(system.NewOTelConfig("system", "vault", VERSION, getHostname())))
	err := systemStub.Connect(context.Background())
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}
	defer systemStub.Close(context.Background())

	// Setup PKCS11
	pkcsCtx, err := pkcs.NewPKCSFromEnv()
	if err != nil {
		panic("Failed to load PKCS from env: " + err.Error())
	}
	log.Info("Using [" + pkcsCtx.GetProviderName() + "] HSM provider.")

	err = pkcsCtx.Initialize()
	if err != nil {
		panic("Failed to initialize PKCS: " + err.Error())
	}
	defer pkcsCtx.Close()

	sealer := pkcs.NewSealer(pkcsCtx)
	sealer.Start()
	defer sealer.Stop()

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	vaultService := service.NewVaultGRPCService(pkcsCtx, sealer)
	system_vault_grpc.RegisterVaultServiceServer(grpcServer, vaultService)

	log.Info("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic("Failed to create listener: " + err.Error())
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic("Error while serving gRPC connections: " + err.Error())
	}
}
