package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/oauth2"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"

	oauthTools "github.com/slamy-solutions/openbp/modules/tools/services/rest/src/domains/accessControl/config"
)

type AuthRouter struct {
	logger *logrus.Entry

	systemStub *system.SystemStub
	nativeStub *native.NativeStub
}

func NewAuthRouter(logger *logrus.Entry, systemStub *system.SystemStub, nativeStub *native.NativeStub) *AuthRouter {
	return &AuthRouter{
		logger: logger,

		systemStub: systemStub,
		nativeStub: nativeStub,
	}
}

type myPasswordInfo struct {
	Enabled bool `json:"enabled"`
}
type myAvailableOAuthProvider struct {
	Name     string `json:"name"`
	ClientId string `json:"clientId"`
	AuthURL  string `json:"authUrl"`
}
type myOAuthProvider struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}
type myOAuthInfo struct {
	AvailableProviders  []myAvailableOAuthProvider `json:"availableProviders"`
	ConfiguredProviders []myOAuthProvider          `json:"configuredProviders"`
}

// type getMyAuthInfoRequest struct{}
type getMyAuthInfoResponse struct {
	Password myPasswordInfo `json:"password"`

	OAuth myOAuthInfo `json:"oauth"`
}

func (r *AuthRouter) GetMyAuthInfo(ctx *gin.Context) {
	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"me"},
			Actions:              []string{"me.auth.password.get", "me.auth.oauth.get"},
			NamespaceIndependent: true,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		r.logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger := authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	passwordExistResponse, err := r.nativeStub.Services.IAM.Authentication.Password.Exists(ctx.Request.Context(), &password.ExistsRequest{
		Namespace: authData.Namespace,
		Identity:  authData.IdentityUUID,
	})
	if err != nil {
		err := errors.New("failed to check password existence: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	oauthConfigResponse, err := r.nativeStub.Services.IAM.Authentication.OAuth.Config.ListProviderConfigs(ctx.Request.Context(), &oauth2.ListProviderConfigsRequest{
		Namespace: authData.Namespace,
	})
	if err != nil {
		err := errors.New("failed to get oauth config: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	oauthDataResponse, err := r.nativeStub.Services.IAM.Authentication.OAuth.OAuth2.GetRegisteredIdentityProviders(ctx.Request.Context(), &oauth2.GetRegisteredIdentityProvidersRequest{
		Namespace: authData.Namespace,
		Identity:  authData.IdentityUUID,
	})
	if err != nil {
		err := errors.New("failed to get oauth data: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var configuredProviders []myOAuthProvider
	for _, provider := range oauthDataResponse.Providers {
		configuredProviders = append(configuredProviders, myOAuthProvider{
			Name:   oauthTools.OAuthProviderGRPCTypeToString[provider.Provider],
			UserId: provider.UserDetails.Id,
		})
	}

	var availableProviders []myAvailableOAuthProvider
	for _, provider := range oauthConfigResponse.Configs {
		if provider.Enabled {
			availableProviders = append(availableProviders, myAvailableOAuthProvider{
				Name:     oauthTools.OAuthProviderGRPCTypeToString[provider.Type],
				ClientId: provider.ClientId,
				AuthURL:  provider.AuthUrl,
			})
		}
	}

	ctx.JSON(http.StatusOK, getMyAuthInfoResponse{
		Password: myPasswordInfo{
			Enabled: passwordExistResponse.Exists,
		},
		OAuth: myOAuthInfo{
			AvailableProviders:  availableProviders,
			ConfiguredProviders: configuredProviders,
		},
	})
}

type setPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (r *AuthRouter) SetOrUpdatePassword(ctx *gin.Context) {
	var requestData setPasswordRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"me"},
			Actions:              []string{"me.auth.password.set", "me.auth.password.update"},
			NamespaceIndependent: true,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		r.logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger := authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Set or update password
	_, err = r.nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx.Request.Context(), &password.CreateOrUpdateRequest{
		Namespace: authData.Namespace,
		Identity:  authData.IdentityUUID,
		Password:  requestData.Password,
	})
	if err != nil {
		err := errors.New("failed to set or update password: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *AuthRouter) DeletePassword(ctx *gin.Context) {
	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"me"},
			Actions:              []string{"me.auth.password.delete"},
			NamespaceIndependent: true,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		r.logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger := authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Delete password
	_, err = r.nativeStub.Services.IAM.Authentication.Password.Delete(ctx.Request.Context(), &password.DeleteRequest{
		Namespace: authData.Namespace,
		Identity:  authData.IdentityUUID,
	})
	if err != nil {
		err := errors.New("failed to delete password: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

type oauthFinalizeRegistrationRequest struct {
	Provider string `json:"provider" binding:"required,oneof=github gitlab microsoft google discord facebook apple twitter oidc oidc2 oidc3 instagram"`
	Code     string `json:"code" binding:"required"`
}

type oauthFinalizeRegistrationResponse struct {
	Status string `json:"status"`
}

func (r *AuthRouter) FinalizeOAuthRegistration(ctx *gin.Context) {
	var requestData oauthFinalizeRegistrationRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"me"},
			Actions:              []string{"me.auth.oauth.register"},
			NamespaceIndependent: true,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		r.logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger := authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	registrationResponse, err := r.nativeStub.Services.IAM.Authentication.OAuth.OAuth2.RegisterProviderForIdentity(ctx.Request.Context(), &oauth2.RegisterProviderForIdentityRequest{
		Namespace:    authData.Namespace,
		Identity:     authData.IdentityUUID,
		Code:         requestData.Code,
		Provider:     oauthTools.OAuthProviderStringTypeToGRPC[requestData.Provider],
		CodeVerifier: "",
		RedirectUrl:  "",
	})
	if err != nil {
		err := errors.New("failed to register provider for identity: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	switch registrationResponse.Status {
	case oauth2.RegisterProviderForIdentityResponse_OK:
		ctx.JSON(http.StatusOK, oauthFinalizeRegistrationResponse{Status: "OK"})
		return
	case oauth2.RegisterProviderForIdentityResponse_ALREADY_REGISTERED:
		ctx.JSON(http.StatusOK, oauthFinalizeRegistrationResponse{Status: "ALREADY_REGISTERED"})
		return
	case oauth2.RegisterProviderForIdentityResponse_PROVIDER_DISABLED:
		ctx.JSON(http.StatusOK, oauthFinalizeRegistrationResponse{Status: "PROVIDER_DISABLED"})
		return
	case oauth2.RegisterProviderForIdentityResponse_ERROR_WHILE_RETRIEVING_AUTH_TOKEN:
		ctx.JSON(http.StatusOK, oauthFinalizeRegistrationResponse{Status: "ERROR_WHILE_RETRIEVING_AUTH_TOKEN"})
		return
	case oauth2.RegisterProviderForIdentityResponse_ERROR_WHILE_FETCHING_USER_DETAILS:
		ctx.JSON(http.StatusOK, oauthFinalizeRegistrationResponse{Status: "ERROR_WHILE_FETCHING_USER_DETAILS"})
		return
	default:
		ctx.AbortWithError(http.StatusInternalServerError, errors.New("unknown status: "+registrationResponse.Status.String()))
	}
}

type oauthForgetProviderRequest struct {
	Provider string `json:"provider" binding:"required,oneof=github gitlab microsoft google discord facebook apple twitter oidc oidc2 oidc3 instagram"`
}

func (r *AuthRouter) ForgetOAuthProvider(ctx *gin.Context) {
	var requestData oauthForgetProviderRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"me"},
			Actions:              []string{"me.auth.oauth.unregister"},
			NamespaceIndependent: true,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		r.logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger := authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	_, err = r.nativeStub.Services.IAM.Authentication.OAuth.OAuth2.ForgetIdentityProvider(ctx.Request.Context(), &oauth2.ForgetIdentityProviderRequest{
		Namespace: authData.Namespace,
		Identity:  authData.IdentityUUID,
		Provider:  oauthTools.OAuthProviderStringTypeToGRPC[requestData.Provider],
	})
	if err != nil {
		err := errors.New("failed to unregister provider for identity: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
