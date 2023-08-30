package accesscontrol

import (
	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"

	actorAPI "github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/accessControl/actor"
	authAPI "github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/accessControl/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/accessControl/config"
)

func FillRouterGroup(group *gin.RouterGroup, nativeStub *native.NativeStub) {
	identityRouter := &IdentityRouter{nativeStub: nativeStub}
	group.GET("/iam/identity/list", identityRouter.List)
	group.GET("/iam/identity", identityRouter.Get)
	group.POST("/iam/identity", identityRouter.Create)
	group.DELETE("/iam/identity", identityRouter.Delete)
	group.PATCH("/iam/identity", identityRouter.Update)
	group.PATCH("/iam/identity/active", identityRouter.SetActive)
	group.PATCH("/iam/identity/addPolicy", identityRouter.AddPolicy)
	group.PATCH("/iam/identity/removePolicy", identityRouter.RemovePolicy)
	group.PATCH("/iam/identity/addRole", identityRouter.AddRole)
	group.PATCH("/iam/identity/removeRole", identityRouter.RemoveRole)

	policyRouter := &PolicyRouter{nativeStub: nativeStub}
	group.POST("/iam/policy", policyRouter.Create)
	group.GET("/iam/policy/list", policyRouter.List)
	group.GET("/iam/policy", policyRouter.Get)
	group.PATCH("/iam/policy", policyRouter.Update)
	group.DELETE("/iam/policy", policyRouter.Delete)

	roleRouter := &RoleRouter{nativeStub: nativeStub}
	group.POST("/iam/role", roleRouter.Create)
	group.GET("/iam/role/list", roleRouter.List)
	group.GET("/iam/role", roleRouter.Get)
	group.DELETE("/iam/role", roleRouter.Delete)
	group.PATCH("/iam/role", roleRouter.Update)
	group.PATCH("/iam/role/addPolicy", roleRouter.AddPolicy)
	group.PATCH("/iam/role/removePolicy", roleRouter.RemovePolicy)

	authenticationPasswordRouter := authAPI.NewPasswordRouter(nativeStub)
	group.GET("/iam/auth/password", authenticationPasswordRouter.GetStatus)
	group.DELETE("/iam/auth/password", authenticationPasswordRouter.Disable)
	group.PUT("/iam/auth/password", authenticationPasswordRouter.SetOrUpdate)

	authenticationCertificateRouter := authAPI.NewCertificateRouter(nativeStub)
	group.GET("/iam/auth/certificate/listForIdentity", authenticationCertificateRouter.ListCertificatesForIdentity)
	group.POST("/iam/auth/certificate", authenticationCertificateRouter.RegisterKeyAndGenerateCertificate)
	group.PATCH("/iam/auth/certificate/disable", authenticationCertificateRouter.Disable)
	group.DELETE("/iam/auth/certificate", authenticationCertificateRouter.Delete)

	actorUserRouter := actorAPI.NewUserRouter(nativeStub)
	group.GET("/iam/actor/user", actorUserRouter.ListUsers)
	group.POST("/iam/actor/user", actorUserRouter.CreateUser)
	group.DELETE("/iam/actor/user", actorUserRouter.DeleteUser)

	configOAuthRouter := config.NewOAuthConfigRouter(nativeStub)
	group.GET("/config/oauth/provider/list", configOAuthRouter.GetProvidersConfigs)
	group.PATCH("/config/oauth/provider", configOAuthRouter.UpdateProviderConfig)
}
