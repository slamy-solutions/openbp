package namespace

import (
	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

func FillRouterGroup(group *gin.RouterGroup, nativeStub *native.NativeStub) {
	listRouter := &ListRouter{nativeStub: nativeStub}

	group.GET("/list", listRouter.List)
	group.POST("/list/namespace", listRouter.Create)
	group.DELETE("/list/namespace", listRouter.Delete)
}
