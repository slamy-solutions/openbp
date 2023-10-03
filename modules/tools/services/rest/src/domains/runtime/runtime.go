package runtime

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	runtimeGRPC "github.com/slamy-solutions/openbp/modules/runtime/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func FillRouterGroup(logger *logrus.Entry, group *gin.RouterGroup, nativeStub *native.NativeStub, systemStub *system.SystemStub, runtimeStub *runtimeGRPC.RuntimeStub) {
	rpcRouter := &RPCRouter{
		systemStub:  systemStub,
		nativeStub:  nativeStub,
		logger:      logger.WithField("domain.service", "rpc"),
		runtimeStub: runtimeStub,
	}

	group.POST("/rpc/call", rpcRouter.Call)
}
