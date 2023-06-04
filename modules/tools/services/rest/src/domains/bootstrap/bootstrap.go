package bootstrap

import (
	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func FillRouterGroup(group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub) {
	rootUserRouter := &RootUserRouter{systemStub: systemStub, nativeStub: nativeStub}
	statusRouter := &StatusRouter{systemStub: systemStub, nativeStub: nativeStub}

	group.GET("/status", statusRouter.Status)
	group.POST("/rootUser", rootUserRouter.InitRootUser)
}
