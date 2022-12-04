package main

import (
	"context"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/auth"
)

func main() {
	// Connect to internal services
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	servicesHandler, err := services.ConnectToServices(ctx)
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(otelgin.Middleware("tools_rest"))

	auth.FillRouterGroup(r.Group("/auth"), servicesHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
