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

	native_storage_bucket_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/bucket"
	native_storage_fs_grpc "github.com/slamy-solutions/openbp/modules/native/libs/golang/storage/fs"
	"github.com/slamy-solutions/openbp/modules/native/services/storage/src/services/bucket"
	"github.com/slamy-solutions/openbp/modules/native/services/storage/src/services/fs"
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
		system.NewSystemStubConfig().
			WithOTel(system.NewOTelConfig("native", "storage", VERSION, getHostname())).
			WithCache().
			WithDB().
			WithNats(),
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
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: time.Minute * 5,
		}),
	)

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	bucketRepository, err := bucket.NewBucketRepository(context.Background(), systemStub, logger)
	if err != nil {
		panic("Failed to create bucket repository: " + err.Error())
	}
	bucketService := bucket.NewService(bucketRepository, logger)
	native_storage_bucket_grpc.RegisterBucketServiceServer(grpcServer, bucketService)

	fileRepository, err := fs.NewFSRepository(systemStub, logger)
	if err != nil {
		panic("Failed to create file repository: " + err.Error())
	}
	fileService := fs.NewService(fileRepository, logger)
	native_storage_fs_grpc.RegisterFSServiceServer(grpcServer, fileService)

	eventHandler, err := NewEventHandlerService(systemStub, logger)
	if err != nil {
		panic("failed to setup event hanle service: " + err.Error())
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
