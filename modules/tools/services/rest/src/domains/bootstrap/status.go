package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"
)

const (
	STATUS_OK                        = "OK"
	STATUS_ROOT_USER_NOT_INITIALIZED = "ROOT_USER_NOT_INITIALIZED"
)

type StatusRouter struct {
	servicesHandler *services.ServicesConnectionHandler

	rootUserRouter *RootUserRouter
}

type statusRequest struct{}
type statusResponse struct {
	FullyBootstrapped bool   `json:"fullyBootstrapped"`
	Code              string `json:"code"`
}

func (r *StatusRouter) Status(ctx *gin.Context) {
	rootUserInited, err := r.rootUserRouter.isRootUserInited(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !rootUserInited {
		ctx.JSON(http.StatusOK, &statusResponse{FullyBootstrapped: true, Code: STATUS_ROOT_USER_NOT_INITIALIZED})
		return
	}

	// TODO: Set cache when everything initialized

	ctx.JSON(http.StatusOK, &statusResponse{FullyBootstrapped: true, Code: STATUS_OK})
}
