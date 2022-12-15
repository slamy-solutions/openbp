package bootstrap

import (
	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"
)

const (
	_ENV_ALLOW_ROOT_USER_INIT       = "ALLOW_ROOT_USER_INIT"
	_ENV_DEFAULT_ROOT_USER_LOGIN    = "DEFAULT_ROOT_USER_LOGIN"
	_ENV_DEFAULT_ROOT_USER_PASSWORD = "DEFAULT_ROOT_USER_PASSWORD"
)

func FillRouterGroup(group *gin.RouterGroup, servicesHandler *services.ServicesConnectionHandler) {
	rootUserRouter := &RootUserRouter{servicesHandler: servicesHandler, allowInit: true}
	statusRouter := &StatusRouter{servicesHandler: servicesHandler, rootUserRouter: rootUserRouter}

	group.GET("/status", statusRouter.Status)
	group.POST("/rootUser", rootUserRouter.InitRootUser)
}
