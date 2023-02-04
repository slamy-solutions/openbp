package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
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

	refreshResponse, err := r.servicesHandler.Native.IAMAuth.RefreshToken(ctx.Request.Context(), &auth.RefreshTokenRequest{
		RefreshToken: requestData.RefreshToken,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	switch refreshResponse.Status {
	case auth.RefreshTokenResponse_OK:
		ctx.JSON(http.StatusOK, refreshTokenResponse{AccessToken: refreshResponse.AccessToken})
	case auth.RefreshTokenResponse_TOKEN_INVALID:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenInvalid))
	case auth.RefreshTokenResponse_TOKEN_NOT_FOUND:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenNotFound))
	case auth.RefreshTokenResponse_TOKEN_DISABLED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenDisabled))
	case auth.RefreshTokenResponse_TOKEN_EXPIRED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenExpired))
	case auth.RefreshTokenResponse_TOKEN_IS_NOT_REFRESH_TOKEN:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshTokenIsNotRefreshToken))
	case auth.RefreshTokenResponse_IDENTITY_NOT_FOUND:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshIdentityNotFound))
	case auth.RefreshTokenResponse_IDENTITY_NOT_ACTIVE:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshIdentityNotActive))
	case auth.RefreshTokenResponse_IDENTITY_UNAUTHENTICATED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthTokenRefreshIdentityUnauthenticated))
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

	checkResponse, err := r.servicesHandler.Native.IAMAuth.CheckAccessWithToken(ctx.Request.Context(), &auth.CheckAccessWithTokenRequest{
		AccessToken: requestData.AccessToken,
		Scopes:      []*auth.Scope{},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	valid := checkResponse.Status == auth.CheckAccessWithTokenResponse_OK

	ctx.JSON(http.StatusOK, validateTokenResponse{Valid: valid, Message: checkResponse.Message})
}
