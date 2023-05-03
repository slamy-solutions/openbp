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
	native_iam_authentication_password_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/authentication/src/eventHandling"
	"github.com/slamy-solutions/openbp/modules/native/services/iam/authentication/src/services"
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
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithVault().WithOTel(system.NewOTelConfig("native", "iam_authentication", VERSION, getHostname())),
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

	ensureIndexesContext, cancelEnsureIndexesContext := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelEnsureIndexesContext()
	err = services.EnsureIndexesForNamespace(ensureIndexesContext, "", systemStub)
	if err != nil {
		panic("Failed to ensure indexes for global namespace " + err.Error())
	}

	eventsHandler, err := eventHandling.NewEventHandlerService(systemStub)
	if err != nil {
		panic("Failed to create events handler. " + err.Error())
	}
	defer eventsHandler.Close()

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	iamAuthenticationPasswordServer := services.NewPasswordIdentificationService(systemStub, nativeStub)
	native_iam_authentication_password_grpc.RegisterIAMAuthenticationPasswordServiceServer(grpcServer, iamAuthenticationPasswordServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
