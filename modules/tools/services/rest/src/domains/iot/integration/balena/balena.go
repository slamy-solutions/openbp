package balena

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func FillRouterGroup(logger *logrus.Entry, group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub, iotStub *iot.IOTStub) {
	serversServer := NewServersServer(nativeStub, iotStub, logger.WithField("integration.balena.service", "servers"))

	group.GET("/servers", serversServer.ListServersInNamespace)
	group.POST("/servers/server", serversServer.CreateServer)
	group.PATCH("/servers/server", serversServer.UpdateServer)
	group.DELETE("/servers/server", serversServer.DeleteServer)
	group.PATCH("/servers/server/enabled", serversServer.SetServerEnabled)

	devicesServer := NewDevicesServer(nativeStub, iotStub, logger.WithField("integration.balena.service", "devices"))
	group.GET("/devices", devicesServer.ListDevicesInNamespace)
	group.PATCH("/devices/device/bind", devicesServer.BindDevice)
	group.PATCH("/devices/device/unbind", devicesServer.UnbindDevice)

	syncServer := NewSyncServer(nativeStub, iotStub, logger.WithField("integration.balena.service", "sync"))
	group.POST("/sync/now", syncServer.SyncNow)
	group.GET("/sync/log", syncServer.ListLog)

	toolsServer := NewToolsServer(nativeStub, iotStub, logger.WithField("integration.balena.service", "tools"))
	group.POST("/tools/verifyConnectionData", toolsServer.VerifyConnectionData)
}
