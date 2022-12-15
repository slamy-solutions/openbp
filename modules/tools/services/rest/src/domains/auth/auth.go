package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"
)

func FillRouterGroup(group *gin.RouterGroup, servicesHandler *services.ServicesConnectionHandler) {
	passwordRouter := &PasswordRouter{servicesHandler: servicesHandler}
	tokenRouter := &TokenRouter{servicesHandler: servicesHandler}

	// Login
	group.POST("/login/password", passwordRouter.Login)

	// Token
	group.POST("/token/refresh", tokenRouter.Refresh)
	group.POST("/token/validate", tokenRouter.Validate)
}
