package accesscontrol

import (
	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

func FillRouterGroup(group *gin.RouterGroup, nativeStub *native.NativeStub) {
	identityRouter := &IdentityRouter{nativeStub: nativeStub}
	group.GET("/iam/identity/list", identityRouter.List)
	group.POST("/iam/identity", identityRouter.Create)
	group.DELETE("/iam/identity", identityRouter.Delete)

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
	group.PATCH("/iam/role/addPolicy", roleRouter.AddPolicy)
	group.PATCH("/iam/role/removePolicy", roleRouter.RemovePolicy)
}
