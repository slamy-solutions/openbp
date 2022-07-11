package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/olebedev/emitter"
	"github.com/slamy-solutions/open-erp/modules/system/libs/go/telemetry"

	native_lambda_grpc "github.com/slamy-solutions/open-erp/modules/native/services/lambda/entrypoint/src/grpc/native_lambda"
	"github.com/slamy-solutions/open-erp/modules/native/services/lambda/entrypoint/src/services"
)

const (
	VERSION = "1.0.0"
)

func getConfigEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	SYSTEM_TELEMETRY_EXPORTER_ENDPOINT := getConfigEnv("SYSTEM_TELEMETRY_EXPORTER_ENDPOINT", "system_telemetry:55680")
	SYSTEM_RABBITMQ_URL := getConfigEnv("SYSTEM_RABBITMQ_URL", "amqp://system_rabbitmq:5672")

	NATIVE_LAMBDA_MANAGER_URL := getConfigEnv("NATIVE_LAMBDA_MANAGER_URL", "native_lambda_manager:80")

	ctx := context.Background()

	// Setting up Telemetry
	telemetryProvider, err := telemetry.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "lambda.entrypoint", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Setting up native_lambda_manager connection
	nativeLambdaManagerConnection, err := grpc.Dial(
		NATIVE_LAMBDA_MANAGER_URL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	defer nativeLambdaManagerConnection.Close()
	nativeLambdaManagerClient := native_lambda_grpc.NewLambdaManagerServiceClient(nativeLambdaManagerConnection)
	fmt.Println("Initialized native_namespace connection")

	// Setting up rabbitmq connection
	amqpConnection, err := amqp.Dial(SYSTEM_RABBITMQ_URL)
	if err != nil {
		panic(err)
	}
	defer amqpConnection.Close()

	amqpResponseChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}
	defer amqpResponseChannel.Close()

	amqpRequestChannel, err := amqpConnection.Channel()
	if err != nil {
		panic(err)
	}
	defer amqpRequestChannel.Close()

	// Setting up response event emitter
	emmiter := emitter.New(2048)
	responseBindKey, closeResponseListener, err := services.ListenForResponces(ctx, amqpResponseChannel, emmiter)
	if err != nil {
		panic(err)
	}
	defer closeResponseListener()

	// Creating grpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	lambdaManagerServer, err := services.NewLambdaEntrypointServer(nativeLambdaManagerClient, amqpRequestChannel, emmiter, responseBindKey)
	if err != nil {
		panic(err)
	}
	native_lambda_grpc.RegisterLambdaEntrypointServiceServer(grpcServer, lambdaManagerServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
