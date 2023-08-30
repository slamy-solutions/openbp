package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"

	accesscontrol "github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/accessControl"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/bootstrap"
	iotDomain "github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/iot"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/me"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/namespace"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	systemStub := system.NewSystemStub(system.NewSystemStubConfig().WithOTel(system.NewOTelConfig("tools", "rest", VERSION, getHostname())).WithVault())
	err := systemStub.Connect(ctx)
	if err != nil {
		panic(err)
	}

	nativeStub := native.NewNativeStub(native.NewStubConfig().WithNamespaceService().WithIAMService())
	err = nativeStub.Connect()
	if err != nil {
		panic(err)
	}

	iotStub := iot.NewIOTStub(iot.NewStubConfig().WithCoreService())
	err = iotStub.Connect()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"} //TODO: somehow handle this. Is this possible?
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(otelgin.Middleware("tools_rest"))

	logger := logrus.StandardLogger()

	auth.FillRouterGroup(r.Group("/api/auth"), systemStub, nativeStub)
	bootstrap.FillRouterGroup(r.Group("/api/bootstrap"), systemStub, nativeStub)
	namespace.FillRouterGroup(r.Group("/api/namespace"), nativeStub)
	accesscontrol.FillRouterGroup(r.Group("/api/accessControl"), nativeStub)
	iotDomain.FillRouterGroup(logger.WithField("domain.name", "iot"), r.Group("/api/iot"), systemStub, nativeStub, iotStub)
	me.FillRouterGroup(logger.WithField("domain.name", "me"), r.Group("/api/me"), systemStub, nativeStub, iotStub)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
