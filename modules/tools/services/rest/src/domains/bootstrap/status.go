package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

type StatusRouter struct {
	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

type statusResponse struct {
	VaultSealed             bool `json:"vaultSealed"`
	RootUserCreated         bool `json:"rootUserCreated"`
	RootUserCreationBlocked bool `json:"rootUserCreationBlocked"`
}

func (r *StatusRouter) Status(ctx *gin.Context) {
	response := &statusResponse{
		VaultSealed:             false,
		RootUserCreated:         false,
		RootUserCreationBlocked: false,
	}

	vaultSealed, err := isVaultSealed(ctx.Request.Context(), r.systemStub)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response.VaultSealed = vaultSealed
	if vaultSealed {
		// Return immediatelly, because if vault is sealed other services may not work
		ctx.JSON(http.StatusOK, response)
		return
	}

	rootUserCreated, err := isRootUserCreated(ctx.Request.Context(), r.nativeStub)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response.RootUserCreated = rootUserCreated

	response.RootUserCreationBlocked = isRootUserCreationBlocked()

	ctx.JSON(http.StatusOK, response)
}
