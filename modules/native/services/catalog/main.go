package main

import (
	context "context"
	"fmt"
	"net"
	"time"

	native_catalog_grpc "slamy/opencrm/native/catalog/grpc/native_catalog_grpc"
	services "slamy/opencrm/native/catalog/services"

	"slamy/opencrm/native/lib/telemetry"

	"slamy/opencrm/native/lib/cache"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"google.golang.org/grpc"
)

func main() {
	// Starting OpenTelemetry
	ctx := context.Background()

	traceProvider, err := telemetry.Register(ctx)
	if err != nil {
		panic(err)
	}
	defer traceProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	// Starting cache service
	cacheConfig, err := cache.ConfigFromEnv()
	if err != nil {
		panic(err)
	}
	cache := cache.NewCache(cacheConfig, nil)
	defer cache.Shutdown(ctx)
	fmt.Println("Initialized cache")

	// Conencting to the DB
	dbOptions := options.Client()
	dbOptions.Monitor = otelmongo.NewMonitor(otelmongo.WithTracerProvider(traceProvider), otelmongo.WithCommandAttributeDisabled(true))
	dbOptions.ApplyURI("mongodb://root:example@system_db/admin")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbClient, err := mongo.Connect(timeoutCtx, dbOptions)
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect(ctx)
	db := dbClient.Database("admin")
	fmt.Println("Initialized DB connection")

	// Connecting to other microservices

	// Initializing

	// Starting gRPC servers
	grpcServer := grpc.NewServer()

	catalogServer := services.CatalogServer{DB: db, Cache: &cache, Tracer: traceProvider.Tracer("native_catalog")}

	native_catalog_grpc.RegisterCatalogServiceServer(grpcServer, &catalogServer)

	fmt.Println("Start listening for gRPC connections")
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
