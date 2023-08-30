package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/accessControl/config"
)

type OAuthRouter struct {
	nativeStub *native.NativeStub
}

type oauthLoginRequest struct {
	Namespace    string `json:"namespace" binding:""`
	Provider     string `json:"provider" binding:"required"`
	Code         string `json:"code" binding:"required"`
	CodeVerifier string `json:"codeVerifier" binding:""`
	RedirectURL  string `json:"redirectURL" binding:""`
}
type oauthLoginResponse struct {
	Login        string `json:"login"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (r *OAuthRouter) Login(ctx *gin.Context) {
	var requestData oauthLoginRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	tokenResponse, err := r.nativeStub.Services.IAM.Auth.CreateTokenWithOAuth2(ctx.Request.Context(), &auth.CreateTokenWithOAuth2Request{
		Namespace:    requestData.Namespace,
		Provider:     requestData.Provider,
		Code:         requestData.Code,
		CodeVerifier: requestData.CodeVerifier,
		RedirectURL:  requestData.RedirectURL,
		Metadata:     MetadataFromRequestContext(ctx).ToJSONString(),
		Scopes:       []*auth.Scope{},
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Cant create toke. Most probably, provider is invalid or not supported."})
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	switch tokenResponse.Status {
	case auth.CreateTokenWithOAuth2Response_OK:
		ctx.JSON(http.StatusOK, oauthLoginResponse{AccessToken: tokenResponse.AccessToken, RefreshToken: tokenResponse.RefreshToken})
	case auth.CreateTokenWithOAuth2Response_UNAUTHORIZED:
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Cant create toke. You cant access requested scopes."})
	case auth.CreateTokenWithOAuth2Response_UNAUTHENTICATED:
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Cant create token. Failed to authenticate."})
	case auth.CreateTokenWithOAuth2Response_IDENTITY_NOT_ACTIVE:
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Cant create token. Identity is not active."})
	default:
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("unsupported response status"))
	}
}

type oauthSupportedProvider struct {
	Name     string `json:"name"`
	ClientID string `json:"clientId"`
	AuthURL  string `json:"authUrl"`
}
type oauthGetSupportedProvidersRequest struct {
	Namespace string `form:"namespace" binding:""`
}
type oauthGetSupportedProvidersResponse struct {
	Providers []oauthSupportedProvider `json:"providers"`
}

func (r *OAuthRouter) GetSupportedProviders(ctx *gin.Context) {
	var requestData oauthGetSupportedProvidersRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	providersResponse, err := r.nativeStub.Services.IAM.Authentication.OAuth.Config.ListProviderConfigs(ctx.Request.Context(), &oauth2.ListProviderConfigsRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	providers := make([]oauthSupportedProvider, 0, len(providersResponse.Configs))
	for _, provider := range providersResponse.Configs {
		if provider.Enabled {
			providers = append(providers, oauthSupportedProvider{Name: config.OAuthProviderGRPCTypeToString[provider.Type], ClientID: provider.ClientId, AuthURL: provider.AuthUrl})
		}
	}

	ctx.JSON(http.StatusOK, oauthGetSupportedProvidersResponse{Providers: providers})
}
