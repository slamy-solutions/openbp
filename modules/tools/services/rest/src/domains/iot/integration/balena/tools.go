package balena

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type ToolsServer struct {
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewToolsServer(nativeStub *native.NativeStub, iotStub *iot.IOTStub, logger *logrus.Entry) ToolsServer {
	return ToolsServer{
		nativeStub: nativeStub,
		iotStub:    iotStub,
		logger:     logger,
	}
}

type verifyConnectionDataRequest struct {
	URL      string `json:"url" binding:"required,lte=256"`
	APIToken string `json:"apiToken" binding:"required,lte=1024"`
}
type verifyConnectionDataResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (s *ToolsServer) VerifyConnectionData(ctx *gin.Context) {
	var requestData verifyConnectionDataRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.url": requestData.URL,
	})

	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{})
	if err != nil {
		err := errors.Join(errors.New("failed to check auth"), err)
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	verifyResponse, err := s.iotStub.Core.Integration.Balena.Tools.VerifyConnectionData(ctx.Request.Context(), &balena.VerifyConnectionDataRequest{
		Url:   requestData.URL,
		Token: requestData.APIToken,
	})
	if err != nil {
		err := errors.Join(errors.New("failed to verify connection data"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var status string = ""
	switch verifyResponse.Status {
	case balena.VerifyConnectionDataResponse_OK:
		status = "OK"
	case balena.VerifyConnectionDataResponse_BAD_URL:
		status = "BAD_URL"
	case balena.VerifyConnectionDataResponse_SERVER_UNAVAILABLE:
		status = "SERVER_UNAVAILABLE"
	case balena.VerifyConnectionDataResponse_SERVER_BAD_RESPONSE:
		status = "SERVER_BAD_RESPONSE"
	default:
		{
			err := errors.New("unknown status")
			logger.Error(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	logger.WithFields(logrus.Fields{
		"connection.status":  verifyResponse.Status,
		"connection.message": verifyResponse.Message,
	}).Debug("Balena connection verified")
	ctx.JSON(http.StatusOK, verifyConnectionDataResponse{
		Status:  status,
		Message: verifyResponse.Message,
	})
}
