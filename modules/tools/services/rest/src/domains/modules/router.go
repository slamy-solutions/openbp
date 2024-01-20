package modules

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	crm_settings_grpc "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/settings"
	erp "github.com/slamy-solutions/openbp/modules/erp/libs/golang"
	erp_catalog_grpc "github.com/slamy-solutions/openbp/modules/erp/libs/golang/core/catalog"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const MODULES_CACHE_TIMEOUT = time.Second * 10
const MODULES_CACHE_KEY = "tools_rest_general_modules_list"

type modulesRouter struct {
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub
	crmStub    *crm.CRMStub
	erpStub    *erp.ERPStub
}

type loadedModulesResponse struct {
	System bool `json:"system"`
	Native bool `json:"native"`
	IOT    bool `json:"iot"`
	CRM    bool `json:"crm"`
	ERP    bool `json:"erp"`
	Tools  bool `json:"tools"`
}

func (r *modulesRouter) GetLoadedModules(ctx *gin.Context) {
	// Get data from cache it possible
	cacheData, _ := r.systemStub.Cache.Get(ctx.Request.Context(), MODULES_CACHE_KEY)
	if cacheData != nil {
		ctx.Status(http.StatusOK)
		ctx.Header("Content-Type", "application/json")
		ctx.Writer.Write(cacheData)
		return
	}

	response := &loadedModulesResponse{
		System: true,
		Native: true,
		Tools:  true,
	}

	checkWaiter := sync.WaitGroup{}

	checkWaiter.Add(1)
	go func() {
		defer checkWaiter.Done()
		_, err := r.iotStub.Core.Fleet.Count(ctx.Request.Context(), &fleet.CountRequest{Namespace: ""})
		response.IOT = err == nil
	}()

	checkWaiter.Add(1)
	go func() {
		defer checkWaiter.Done()
		_, err := r.crmStub.Core.Settings.GetSettings(ctx.Request.Context(), &crm_settings_grpc.GetSettingsRequest{Namespace: "", UseCache: true})
		response.CRM = err == nil
	}()

	checkWaiter.Add(1)
	go func() {
		defer checkWaiter.Done()
		_, err := r.erpStub.Core.Catalog.Catalog.GetAll(ctx.Request.Context(), &erp_catalog_grpc.GetAllCatalogsRequest{
			Namespace: "", UseCache: true,
		})
		response.ERP = err == nil
	}()

	checkWaiter.Wait()

	//Add data to the cache
	cacheData, _ = json.Marshal(response)
	r.systemStub.Cache.Set(ctx.Request.Context(), MODULES_CACHE_KEY, cacheData, MODULES_CACHE_TIMEOUT)

	ctx.JSON(http.StatusOK, response)
}
