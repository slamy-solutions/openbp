package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/oauth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/models"
)

type PasswordRouter struct {
	servicesHandler *services.ServicesConnectionHandler
}

type passwordLoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type passwordLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (r *PasswordRouter) Login(ctx *gin.Context) {
	var requestData passwordLoginRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Try to find user with this login
	userGetResponse, err := r.servicesHandler.Native.ActorUser.GetByLogin(ctx, &user.GetByLoginRequest{
		Login:    requestData.Login,
		UseCache: false,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordCredentialsInvalid))
				return
			}
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Try to verify password and create authorization token for user
	tokenCreationResponse, err := r.servicesHandler.Native.IAMOAuth.CreateTokenWithPassword(ctx, &oauth.CreateTokenWithPasswordRequest{
		Namespace: "",
		Identity:  userGetResponse.User.Identity,
		Password:  requestData.Password,
		Metadata:  MetadataFromRequestContext(ctx).ToJSONString(),
		Scopes:    []*oauth.Scope{},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	switch tokenCreationResponse.Status {
	case oauth.CreateTokenWithPasswordResponse_OK:
		ctx.JSON(http.StatusOK, passwordLoginResponse{AccessToken: tokenCreationResponse.AccessToken, RefreshToken: tokenCreationResponse.RefreshToken})
		break
	case oauth.CreateTokenWithPasswordResponse_IDENTITY_NOT_ACTIVE:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordIdentityNotActive))
		break
	case oauth.CreateTokenWithPasswordResponse_CREDENTIALS_INVALID:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordCredentialsInvalid))
		break
	case oauth.CreateTokenWithPasswordResponse_UNAUTHORIZED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordNotEnoughtPrivileges))
		break
	default:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordUnauthorizedUnknown))
	}
}