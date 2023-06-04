package auth

import (
	"github.com/gin-gonic/gin"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func FillRouterGroup(group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub) {
	passwordRouter := &PasswordRouter{nativeStub: nativeStub}
	tokenRouter := &TokenRouter{nativeStub: nativeStub}

	// Login
	group.POST("/login/password", passwordRouter.Login)

	// Token
	group.POST("/token/refresh", tokenRouter.Refresh)
	group.POST("/token/validate", tokenRouter.Validate)
}
