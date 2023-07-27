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

	"github.com/sirupsen/logrus"
	deviceGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	fleetGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	telemetryGRPC "github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/telemetry"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/device"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/fleet"
	"github.com/slamy-solutions/openbp/modules/iot/services/core/src/services/telemetry"
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
		system.NewSystemStubConfig().WithCache().WithDB().WithNats().WithOTel(system.NewOTelConfig("iot", "core", VERSION, getHostname())),
	)
	systemConnectionContext, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := systemStub.Connect(systemConnectionContext)
	if err != nil {
		panic("Failed to connect to the system services: " + err.Error())
	}
	defer systemStub.Close(context.Background())

	// Connect to the "native" module services
	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
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

	logger := logrus.StandardLogger()

	deviceServer, err := device.NewDeviceServer(context.Background(), logger.WithField("service", "device"), systemStub, nativeStub)
	if err != nil {
		panic("Failed to setup device server: " + err.Error())
	}
	deviceGRPC.RegisterDeviceServiceServer(grpcServer, deviceServer)

	fleetServer, err := fleet.NewFleetServer(context.Background(), logger.WithField("service", "fleet"), systemStub, nativeStub, deviceServer)
	if err != nil {
		panic("Failed to setup fleet server: " + err.Error())
	}
	fleetGRPC.RegisterFleetServiceServer(grpcServer, fleetServer)

	telemetryServer, err := telemetry.NewTelemetryServer(context.Background(), logger.WithField("service", "telemetry"), systemStub, nativeStub, deviceServer)
	if err != nil {
		panic("Failed to setup telemetry server: " + err.Error())
	}
	telemetryGRPC.RegisterTelemetryServiceServer(grpcServer, telemetryServer)

	eventHandler, err := NewEventHandler(logger.WithField("service", "event_handler"), systemStub)
	if err != nil {
		panic("Failed to start event handler: " + err.Error())
	}
	defer eventHandler.Close()

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
