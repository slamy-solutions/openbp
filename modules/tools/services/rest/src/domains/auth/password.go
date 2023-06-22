package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/models"
)

type PasswordRouter struct {
	nativeStub *native.NativeStub
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
	userGetResponse, err := r.nativeStub.Services.ActorUser.GetByLogin(ctx.Request.Context(), &user.GetByLoginRequest{
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
	tokenCreationResponse, err := r.nativeStub.Services.IAM.Auth.CreateTokenWithPassword(ctx.Request.Context(), &auth.CreateTokenWithPasswordRequest{
		Namespace: "",
		Identity:  userGetResponse.User.Identity,
		Password:  requestData.Password,
		Metadata:  MetadataFromRequestContext(ctx).ToJSONString(),
		Scopes:    []*auth.Scope{},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	switch tokenCreationResponse.Status {
	case auth.CreateTokenWithPasswordResponse_OK:
		ctx.JSON(http.StatusOK, passwordLoginResponse{AccessToken: tokenCreationResponse.AccessToken, RefreshToken: tokenCreationResponse.RefreshToken})
	case auth.CreateTokenWithPasswordResponse_IDENTITY_NOT_ACTIVE:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordIdentityNotActive))
	case auth.CreateTokenWithPasswordResponse_CREDENTIALS_INVALID:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordCredentialsInvalid))
	case auth.CreateTokenWithPasswordResponse_UNAUTHORIZED:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordNotEnoughtPrivileges))
	default:
		ctx.JSON(http.StatusUnauthorized, models.NewAPIError(models.ErrorAuthPasswordUnauthorizedUnknown))
	}
}

/*type passwordRestRequest struct {
	Login string `json:"login" binding:"required"`
}
type passwordResetResponse struct{}

func (r *PasswordRouter) Reset(ctx *gin.Context) {
	var requestData passwordLoginRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Try to find user with this login
	userGetResponse, err := r.servicesHandler.Native.ActorUser.GetByLogin(ctx.Request.Context(), &user.GetByLoginRequest{
		Login:    requestData.Login,
		UseCache: true,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			// User not found. Silently returning success
			if st.Code() == codes.NotFound {
				ctx.JSON(http.StatusOK, gin.H{})
				return
			}
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// TODO:
	// If user email is not set or it is not verified - silently return success.
	if userGetResponse.User.Email == "" {
		ctx.JSON(http.StatusOK, gin.H{})
		return
	}

	ctx.AbortWithError(http.StatusNotImplemented, err)
}*/
