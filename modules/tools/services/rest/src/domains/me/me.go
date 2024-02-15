package me

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/me/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/me/user"
)

func FillRouterGroup(logger *logrus.Entry, group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub, iotStub *iot.IOTStub) {

	authRouter := auth.NewAuthRouter(logger.WithField("domain.service", "auth"), systemStub, nativeStub)
	group.GET("/auth", authRouter.GetMyAuthInfo)
	group.PUT("/auth/password", authRouter.SetOrUpdatePassword)
	group.DELETE("/auth/password", authRouter.DeletePassword)
	group.POST("/auth/oauth/provider", authRouter.FinalizeOAuthRegistration)
	group.DELETE("/auth/oauth/provider", authRouter.ForgetOAuthProvider)

	userRouter := user.NewUserRouter(logger.WithField("domain.service", "user"), nativeStub)
	group.GET("/user", userRouter.GetMyUserInfo)
}
