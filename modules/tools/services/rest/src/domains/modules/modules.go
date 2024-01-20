package modules

import (
	"github.com/gin-gonic/gin"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	erp "github.com/slamy-solutions/openbp/modules/erp/libs/golang"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

func FillRouterGroup(group *gin.RouterGroup, systemStub *system.SystemStub, nativeStub *native.NativeStub, iotStub *iot.IOTStub, crmStub *crm.CRMStub, erpStub *erp.ERPStub) {
	rootUserRouter := &modulesRouter{
		systemStub: systemStub,
		nativeStub: nativeStub,
		iotStub:    iotStub,
		crmStub:    crmStub,
		erpStub:    erpStub,
	}

	group.GET("/status", rootUserRouter.GetLoadedModules)
}
