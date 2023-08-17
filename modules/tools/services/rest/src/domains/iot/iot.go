package iot

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/iot/integration"
)

func FillRouterGroup(logger *logrus.Entry, group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub, iotStub *iot.IOTStub) {
	deviceRouter := NewDeviceRouter(logger.WithField("domain.service", "device"), nativeStub, iotStub)

	group.POST("/devices/device", deviceRouter.Create)
	group.GET("/devices", deviceRouter.List)
	group.GET("/devices/device", deviceRouter.Get)
	group.PATCH("/devices/device", deviceRouter.Update)
	group.DELETE("/devices/device", deviceRouter.Delete)

	fleetRouter := NewFleetRouter(logger.WithField("domain.service", "fleet"), nativeStub, iotStub)
	group.POST("/fleets/fleet", fleetRouter.Create)
	group.GET("/fleets/fleet", fleetRouter.Get)
	group.GET("/fleets", fleetRouter.List)
	group.PATCH("/fleets/fleet", fleetRouter.Update)
	group.DELETE("/fleets/fleet", fleetRouter.Delete)
	group.GET("/fleets/fleet/devices", fleetRouter.ListDevices)
	group.PUT("/fleets/fleet/devices/device", fleetRouter.AddDevice)
	group.DELETE("/fleets/fleet/devices/device", fleetRouter.RemoveDevice)

	telemetryRouter := NewTelemetryRouter(logger.WithField("domain.service", "telemetry"), systemStub, nativeStub, iotStub)
	group.GET("/telemetry/listen", telemetryRouter.Listen)

	integration.FillRouterGroup(logger.WithField("domain.service", "integration"), group.Group("/integration"), systemStub, nativeStub, iotStub)
}
