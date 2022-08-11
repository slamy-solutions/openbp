package routes

import (
	"encoding/json"
	"net/http"

	actorUserGRPC "github.com/slamy-solutions/open-erp/modules/native/services/api/src/grpc/native_actor_user"
	iamAuthGRPC "github.com/slamy-solutions/open-erp/modules/native/services/api/src/grpc/native_iam_auth"
	iamAuthenticationPasswordGRPC "github.com/slamy-solutions/open-erp/modules/native/services/api/src/grpc/native_iam_authentication_password"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func RegisterAuthRoutes(group *gin.RouterGroup, iamAuthClient iamAuthGRPC.IAMAuthServiceClient, actorUserClient actorUserGRPC.ActorUserServiceClient, iamAuthenticationPasswordClient iamAuthenticationPasswordGRPC.IAMAuthenticationPasswordServiceClient) {
	api := &authAPI{
		iamAuthClient:                   iamAuthClient,
		actorUserClient:                 actorUserClient,
		iamAuthenticationPasswordClient: iamAuthenticationPasswordClient,
	}
	group.POST("login", api.login)
	group.POST("register", api.register)
	group.POST("token/refresh", api.refreshToken)
	group.POST("token/validate", api.validateToken)
}

type authAPI struct {
	iamAuthClient                   iamAuthGRPC.IAMAuthServiceClient
	actorUserClient                 actorUserGRPC.ActorUserServiceClient
	iamAuthenticationPasswordClient iamAuthenticationPasswordGRPC.IAMAuthenticationPasswordServiceClient
}

type loginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type loginMetadata struct {
	IP        string `json:"ip"`
	UserAgent string `json:"userAgent"`
	Referer   string `json:"referer"`
}

func (api *authAPI) login(ctx *gin.Context) {
	var requestData loginRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	userResponse, err := api.actorUserClient.GetByLogin(ctx, &actorUserGRPC.GetByLoginRequest{
		Login:    requestData.Login,
		UseCache: false,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User with specified login doesnt exist or password doesnt match."})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	metadata := loginMetadata{
		IP:        ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
		Referer:   ctx.GetHeader("Referer"),
	}
	metadataString, err := json.Marshal(metadata)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	tokenResponse, err := api.iamAuthClient.CreateTokenWithPassword(ctx, &iamAuthGRPC.CreateTokenWithPasswordRequest{
		Namespace: "",
		Identity:  userResponse.User.Identity,
		Password:  requestData.Password,
		Metadata:  string(metadataString),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	switch tokenResponse.Status {
	case iamAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHENTICATED:
		{
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User with specified login doesnt exist or password doesnt match."})
			return
		}
	case iamAuthGRPC.CreateTokenWithPasswordResponse_UNAUTHORIZED:
		{
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Not enoth privileges to create token with specified scopes"})
			return
		}
	}

	ctx.JSON(http.StatusOK, loginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
	})
}

type registerRequest struct {
	Login    string `json:"login" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerResponse struct{}

func (api *authAPI) register(ctx *gin.Context) {
	var requestData registerRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Create new user
	userCreateResponse, err := api.actorUserClient.Create(ctx, &actorUserGRPC.CreateRequest{
		Login:    requestData.Login,
		FullName: requestData.FullName,
		Avatar:   "",
		Email:    requestData.Email,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok && e.Code() == codes.AlreadyExists {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "User with specified login already exists."})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	// Set up password for user
	_, err = api.iamAuthenticationPasswordClient.CreateOrUpdate(ctx, &iamAuthenticationPasswordGRPC.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  userCreateResponse.User.Identity,
		Password:  requestData.Password,
	})
	if err != nil {
		// Try to delete created user. Very low chance of error that will not affect system. Price of error handling is highter than impact.
		api.actorUserClient.Delete(ctx, &actorUserGRPC.DeleteRequest{Uuid: userCreateResponse.User.Uuid})
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, registerResponse{})
}

func (api *authAPI) refreshToken(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, registerResponse{})
}

func (api *authAPI) validateToken(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, registerResponse{})
}
