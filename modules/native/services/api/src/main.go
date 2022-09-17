package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/slamy-solutions/openbp/modules/system/libs/go/telemetry"

	"github.com/slamy-solutions/openbp/modules/native/services/api/src/routes"
	"github.com/slamy-solutions/openbp/modules/native/services/api/src/tools"
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

	// Setting up Telemetry
	ctx := context.Background()
	telemetryProvider, err := telemetry.Register(ctx, SYSTEM_TELEMETRY_EXPORTER_ENDPOINT, "native", "api", VERSION, "1")
	if err != nil {
		panic(err)
	}
	defer telemetryProvider.Shutdown(ctx)
	fmt.Println("Initialized telemetry")

	tools, err := tools.NewConnectorTools()
	if err != nil {
		panic(err)
	}
	defer tools.Close()

	r := gin.Default()
	r.Use(otelgin.Middleware("native_api"))

	authGroup := r.Group("/api/bff/native/auth")
	routes.RegisterAuthRoutes(authGroup, tools.IAmAuth, tools.ActorUser, tools.IAmAuthenticationPassword)

	r.Run("0.0.0.0:80")
}
