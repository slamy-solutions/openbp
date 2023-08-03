package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/iot/integration/balena"
)

func FillRouterGroup(logger *logrus.Entry, group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub, iotStub *iot.IOTStub) {
	balena.FillRouterGroup(logger.WithField("integartion.name", "balena"), group.Group("/balena"), systemStub, nativeStub, iotStub)
}
