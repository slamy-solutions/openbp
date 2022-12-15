package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/models"
)

type TokenRouter struct {
	servicesHandler *services.ServicesConnectionHandler
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
type refreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func (r *TokenRouter) Refresh(ctx *gin.Context) {
	var requestData refreshTokenRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	refreshResponse, err := r.servicesHandler.Native.IAMOAuth.RefreshToken(ctx.Request.Context(), &oauth.RefreshTokenRequest{
		RefreshToken: requestData.RefreshToken,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	switch refreshResponse.Status {
	case oauth.RefreshTokenResponse_OK:
		ctx.JSON(http.StatusOK, refreshTokenResponse{AccessToken: refreshResponse.AccessToken})
		break
	case oauth.RefreshTokenResponse_TOKEN_INVALID:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenInvalid))
		break
	case oauth.RefreshTokenResponse_TOKEN_NOT_FOUND:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenNotFound))
		break
	case oauth.RefreshTokenResponse_TOKEN_DISABLED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenDisabled))
		break
	case oauth.RefreshTokenResponse_TOKEN_EXPIRED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenExpired))
		break
	case oauth.RefreshTokenResponse_TOKEN_IS_NOT_REFRESH_TOKEN:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenIsNotRefreshToken))
		break
	case oauth.RefreshTokenResponse_IDENTITY_NOT_FOUND:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshIdentityNotFound))
		break
	case oauth.RefreshTokenResponse_IDENTITY_NOT_ACTIVE:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshIdentityNotActive))
		break
	case oauth.RefreshTokenResponse_IDENTITY_UNAUTHENTICATED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshIdentityUnauthenticated))
		break
	default:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshUnknown))
	}
}

type validateTokenRequest struct {
	AccessToken string `json:"accessToken" binding:"required"`
}
type validateTokenResponse struct {
	Valid bool `json:"valid"`

	Message string `json:"message"`
}

func (r *TokenRouter) Validate(ctx *gin.Context) {
	var requestData validateTokenRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	checkResponse, err := r.servicesHandler.Native.IAMOAuth.CheckAccess(ctx.Request.Context(), &oauth.CheckAccessRequest{
		AccessToken: requestData.AccessToken,
		Scopes:      []*oauth.Scope{},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	valid := checkResponse.Status == oauth.CheckAccessResponse_OK

	ctx.JSON(http.StatusOK, validateTokenResponse{Valid: valid, Message: checkResponse.Message})
}
