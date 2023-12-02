package crm

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func FillRouterGroup(logger *logrus.Entry, group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub, crmStub *crm.CRMStub) {
	settingsRouter := settingsRouter{nativeStub: nativeStub, crmStub: crmStub, logger: logger.WithField("domain.service", "settings")}
	group.GET("/settings", settingsRouter.GetSettings)
	group.PATCH("/settings", settingsRouter.SetSettings)
	group.POST("/settings/onec/connection", settingsRouter.checkOneCConnection)

	clientRoter := clientRouter{crmStub: crmStub, nativeStub: nativeStub}
	group.POST("/clients/client", clientRoter.CreateClient)
	group.GET("/clients", clientRoter.GetAllClients)
	group.GET("/clients/client", clientRoter.GetClient)
	group.PATCH("/clients/client", clientRoter.UpdateClient)
	group.DELETE("/clients/client", clientRoter.DeleteClient)

	group.POST("/clients/client/contacts/contact", clientRoter.AddContactPerson)
	group.GET("/clients/client/contacts", clientRoter.GetContactPersonsForClient)
	group.PATCH("/clients/client/contacts/contact", clientRoter.UpdateContactPerson)
	group.DELETE("/clients/client/contacts/contact", clientRoter.DeleteContactPerson)

	onecSyncRouter := onecsyncRouter{nativeStub: nativeStub, crmStub: crmStub, logger: logger.WithField("domain.service", "onec_sync")}
	group.POST("/onec/sync/now", onecSyncRouter.syncNow)
	group.GET("/onec/sync/log", onecSyncRouter.getLog)

	projectRouter := projectRouter{crmStub: crmStub, nativeStub: nativeStub, logger: logger.WithField("domain.service", "project")}
	group.POST("/projects/project", projectRouter.createProject)
	group.GET("/projects", projectRouter.getProjects)
	group.GET("/projects/project", projectRouter.getProject)
	group.PATCH("/projects/project", projectRouter.updateProject)
	group.DELETE("/projects/project", projectRouter.deleteProject)

	departmentRouter := departmentRouter{crmStub: crmStub, nativeStub: nativeStub, logger: logger.WithField("domain.service", "department")}
	group.POST("/departments/department", departmentRouter.createDepartment)
	group.GET("/departments", departmentRouter.getDepartments)
	group.GET("/departments/department", departmentRouter.getDepartment)
	group.PATCH("/departments/department", departmentRouter.updateDepartment)
	group.DELETE("/departments/department", departmentRouter.deleteDepartment)
}
