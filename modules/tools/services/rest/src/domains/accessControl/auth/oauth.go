package auth

import (
	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
)

type OAuthRouter struct {
	nativeStub *native.NativeStub
}

func NewOAuthRouter(nativeStub *native.NativeStub) *OAuthRouter {
	return &OAuthRouter{
		nativeStub: nativeStub,
	}
}

type oauthListRegisteredProvidersRequest struct {
}

func (r *OAuthRouter) ListRegisteredProviders(ctx *gin.Context) {

}
