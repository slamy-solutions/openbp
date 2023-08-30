package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var OAuthProviderTypes = []string{"github", "gitlab", "microsoft", "google", "discord", "facebook", "apple", "twitter", "oidc", "oidc2", "oidc3", "instagram"}
var OAuthProviderStringTypeToGRPC = map[string]oauth2.ProviderType{
	"github":    oauth2.ProviderType_GITHUB,
	"gitlab":    oauth2.ProviderType_GITLAB,
	"microsoft": oauth2.ProviderType_MICROSOFT,
	"google":    oauth2.ProviderType_GOOGLE,
	"discord":   oauth2.ProviderType_DISCORD,
	"facebook":  oauth2.ProviderType_FACEBOOK,
	"apple":     oauth2.ProviderType_APPLE,
	"twitter":   oauth2.ProviderType_TWITTER,
	"oidc":      oauth2.ProviderType_OIDC,
	"oidc2":     oauth2.ProviderType_OIDC2,
	"oidc3":     oauth2.ProviderType_OIDC3,
	"instagram": oauth2.ProviderType_INSTAGRAM,
}
var OAuthProviderGRPCTypeToString = map[oauth2.ProviderType]string{
	oauth2.ProviderType_GITHUB:    "github",
	oauth2.ProviderType_GITLAB:    "gitlab",
	oauth2.ProviderType_MICROSOFT: "microsoft",
	oauth2.ProviderType_GOOGLE:    "google",
	oauth2.ProviderType_DISCORD:   "discord",
	oauth2.ProviderType_FACEBOOK:  "facebook",
	oauth2.ProviderType_APPLE:     "apple",
	oauth2.ProviderType_TWITTER:   "twitter",
	oauth2.ProviderType_OIDC:      "oidc",
	oauth2.ProviderType_OIDC2:     "oidc2",
	oauth2.ProviderType_OIDC3:     "oidc3",
	oauth2.ProviderType_INSTAGRAM: "instagram",
}

type OAuthConfigRouter struct {
	nativeStub *native.NativeStub
}

func NewOAuthConfigRouter(nativeStub *native.NativeStub) *OAuthConfigRouter {
	return &OAuthConfigRouter{
		nativeStub: nativeStub,
	}
}

type oauthProviderConfig struct {
	Namespace    string `json:"namespace"`
	Enabled      bool   `json:"enabled"`
	Name         string `json:"name"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	AuthURL      string `json:"authURL"`
	TokenURL     string `json:"tokenURL"`
	UserApiUrl   string `json:"userApiURL"`
}
type oauthGetProvidersConfigsRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
}
type oauthGetProvidersConfigsResponse struct {
	Providers []oauthProviderConfig `json:"providers"`
}

func (r *OAuthConfigRouter) GetProvidersConfigs(ctx *gin.Context) {
	var requestData oauthGetProvidersConfigsRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.config.oauth.provider.*"},
			Actions:              []string{"native.iam.config.oauth.provider.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	// Get config
	response, err := r.nativeStub.Services.IAM.Authentication.OAuth.Config.ListProviderConfigs(ctx.Request.Context(), &oauth2.ListProviderConfigsRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Convert to response
	providers := make([]oauthProviderConfig, 0, len(response.Configs))
	for _, provider := range response.Configs {
		providers = append(providers, oauthProviderConfig{
			Namespace:    provider.Namespace,
			Enabled:      provider.Enabled,
			Name:         OAuthProviderGRPCTypeToString[provider.Type],
			ClientID:     provider.ClientId,
			ClientSecret: provider.ClientSecret,
			AuthURL:      provider.AuthUrl,
			TokenURL:     provider.TokenUrl,
			UserApiUrl:   provider.UserApiUrl,
		})
	}

	ctx.JSON(http.StatusOK, oauthGetProvidersConfigsResponse{
		Providers: providers,
	})
}

type oauthUpdateProviderConfigRequest struct {
	Namespace    string `json:"namespace"`
	Enabled      bool   `json:"enabled"`
	Name         string `json:"name" binding:"required,oneof=github gitlab microsoft google discord facebook apple twitter oidc oidc2 oidc3 instagram"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	AuthURL      string `json:"authURL"`
	TokenURL     string `json:"tokenURL"`
	UserApiUrl   string `json:"userApiURL"`
}
type oauthUpdateProviderConfigResponse struct{}

func (r *OAuthConfigRouter) UpdateProviderConfig(ctx *gin.Context) {
	var requestData oauthUpdateProviderConfigRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.config.oauth.provider." + requestData.Name},
			Actions:              []string{"native.iam.config.oauth.provider.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	// Update config
	if _, err := r.nativeStub.Services.IAM.Authentication.OAuth.Config.UpdateProviderConfig(ctx.Request.Context(), &oauth2.UpdateProviderConfigRequest{
		Config: &oauth2.ProviderConfig{
			Namespace:    requestData.Namespace,
			Enabled:      requestData.Enabled,
			Type:         OAuthProviderStringTypeToGRPC[requestData.Name],
			ClientId:     requestData.ClientID,
			ClientSecret: requestData.ClientSecret,
			AuthUrl:      requestData.AuthURL,
			TokenUrl:     requestData.TokenURL,
			UserApiUrl:   requestData.UserApiUrl,
		},
	}); err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.FailedPrecondition {
			ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "Probably the vault is sealed"})
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, oauthUpdateProviderConfigResponse{})
}
