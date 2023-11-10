package crm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	settings "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/settings"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type settingsRouter struct {
	nativeStub *native.NativeStub
	crmStub    *crm.CRMStub

	logger *logrus.Entry
}

type formatedSettings struct {
	BackendType string  `json:"backendType" validation:"required,oneof=NATIVE ONE_C"`
	RemoteURL   *string `json:"backendURL"`
	Token       *string `json:"token"`
}

func (s *formatedSettings) ToGRPC(namespace string) *settings.Settings {
	var oneBackend *settings.OneCBackendSettings
	if s.BackendType == "ONE_C" {
		oneBackend = &settings.OneCBackendSettings{
			RemoteURL: *s.RemoteURL,
			Token:     *s.Token,
		}
	}

	return &settings.Settings{
		Namespace:   namespace,
		BackendType: settings.BackendType(settings.BackendType_value[s.BackendType]),
		Backend:     &settings.Settings_OneC{OneC: oneBackend},
	}
}
func formatedSettingsFromGRPC(s *settings.Settings) formatedSettings {
	remoteURL := ""
	token := ""

	if s.GetOneC() != nil {
		remoteURL = s.GetOneC().RemoteURL
		token = s.GetOneC().Token
	}

	return formatedSettings{
		BackendType: s.BackendType.String(),
		RemoteURL:   &remoteURL,
		Token:       &token,
	}
}

type getSettingsRequest struct {
	Namespace string `form:"namespace"`
}
type getSettingsResponse struct {
	Settings formatedSettings `json:"settings"`
}

func (r *settingsRouter) GetSettings(ctx *gin.Context) {
	var requestData getSettingsRequest
	if err := ctx.ShouldBindQuery(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithField("namespace", requestData.Namespace)

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.settings"},
			Actions:              []string{"crm.settings.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	s, err := r.crmStub.Core.Settings.GetSettings(ctx.Request.Context(), &settings.GetSettingsRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		err := errors.New("failed to get settings: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, getSettingsResponse{
		Settings: formatedSettingsFromGRPC(s.Settings),
	})
}

type setSettingsRequest struct {
	Namespace string           `json:"namespace"`
	Settings  formatedSettings `json:"settings" binding:"required"`
}

func (r *settingsRouter) SetSettings(ctx *gin.Context) {
	var requestData setSettingsRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithField("namespace", requestData.Namespace)

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"crm.settings"},
			Actions:              []string{"crm.settings.set"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	_, err = r.crmStub.Core.Settings.SetSettings(ctx.Request.Context(), &settings.SetSettingsRequest{
		Settings: requestData.Settings.ToGRPC(requestData.Namespace),
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Settings where updated")
	ctx.JSON(http.StatusOK, gin.H{})
}

type checkOneCConnectionRequest struct {
	RemoteURL string `json:"backendURL" binding:"required"`
	Token     string `json:"token" binding:"required"`
}
type checkOneCConnectionResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (r *settingsRouter) checkOneCConnection(ctx *gin.Context) {
	var requestData checkOneCConnectionRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{})
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
	r.logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	pingMessage := "Hello 1C from OpenBP"

	// Make request to server
	var pingRequestData struct {
		Ping string `json:"ping"`
	}
	pingRequestData.Ping = pingMessage
	pingRequestDataBytes, err := json.Marshal(pingRequestData)
	if err != nil {
		err := errors.New("failed to marshal ping request data: " + err.Error())
		r.logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	pingResponse, err := http.Post(requestData.RemoteURL+"/ping/"+requestData.Token, "application/json", bytes.NewBuffer(pingRequestDataBytes))
	if err != nil {
		ctx.JSON(http.StatusOK, checkOneCConnectionResponse{
			Success:    false,
			StatusCode: 0,
			Message:    strings.ReplaceAll(err.Error(), requestData.Token, "xxxx-xxxx-xxxx-xxxx"),
		})
		return
	}

	pingResponseBodyBytes, err := io.ReadAll(pingResponse.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, checkOneCConnectionResponse{
			Success:    false,
			StatusCode: 0,
			Message:    err.Error(),
		})
		return
	}

	if pingResponse.StatusCode != 200 {
		ctx.JSON(http.StatusOK, checkOneCConnectionResponse{
			Success:    false,
			StatusCode: pingResponse.StatusCode,
			Message:    strings.ReplaceAll(string(pingResponseBodyBytes), requestData.Token, "xxxx-xxxx-xxxx-xxxx"),
		})
		return
	}

	var pingResponseData struct {
		Pong string `json:"pong"`
	}
	err = json.Unmarshal(pingResponseBodyBytes, &pingResponseData)
	if err != nil {
		ctx.JSON(http.StatusOK, checkOneCConnectionResponse{
			Success:    false,
			StatusCode: 0,
			Message:    "Failed to unmarshall JSON response: " + err.Error(),
		})
		return
	}

	if pingResponseData.Pong != pingMessage {
		ctx.JSON(http.StatusOK, checkOneCConnectionResponse{
			Success:    false,
			StatusCode: 0,
			Message:    "Invalid response. Ping and pong messages are not equal",
		})
		return
	}

	ctx.JSON(http.StatusOK, checkOneCConnectionResponse{
		Success:    true,
		StatusCode: 200,
		Message:    "Connection is OK",
	})
}
