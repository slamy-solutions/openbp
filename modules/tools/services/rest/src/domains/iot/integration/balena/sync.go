package balena

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SyncServer struct {
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewSyncServer(nativeStub *native.NativeStub, iotStub *iot.IOTStub, logger *logrus.Entry) SyncServer {
	return SyncServer{
		nativeStub: nativeStub,
		iotStub:    iotStub,
		logger:     logger,
	}
}

type syncNowRequest struct {
	ServerUUID string `json:"serverUUID" binding:"lte=64"`
}
type syncNowResponse struct{}

func (s *SyncServer) SyncNow(ctx *gin.Context) {
	var requestData syncNowRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.serverUUID": requestData.ServerUUID,
	})

	getServerResponse, err := s.iotStub.Core.Integration.Balena.Server.Get(ctx.Request.Context(), &balena.GetServerRequest{
		Uuid: requestData.ServerUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "balena server not found"})
			return
		}

		err := errors.Join(errors.New("failed to get balena server"), err)
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getServerResponse.Server.Namespace,
			Resources:            []string{"iot.core.integration.balena.server." + getServerResponse.Server.Name},
			Actions:              []string{"iot.core.integration.balena.sync.run"},
			NamespaceIndependent: false,
		},
	})
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

	_, err = s.iotStub.Core.Integration.Balena.Sync.SyncNow(ctx.Request.Context(), &balena.SyncNowRequest{
		BalenaServerUUID: requestData.ServerUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "balena server not found"})
			return
		}

		err := errors.Join(errors.New("failed to sync balena server"), err)
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("balena server manually synced")
	ctx.JSON(http.StatusOK, syncNowResponse{})
}

type listSyncLogRequest struct {
	ServerUUID string `form:"serverUUID" binding:"lte=64"`
	Skip       int64  `form:"skip" binding:"gte=0"`
	Limit      int64  `form:"limit" binding:"gt=0,lte=100"`
}
type listSyncLogResponse struct {
	LogEntries []formatedSyncLogEntry `json:"logEntries"`
	TotalCobnt int64                  `json:"totalCount"`
}

func (s *SyncServer) ListLog(ctx *gin.Context) {
	var requestData listSyncLogRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.serverUUID": requestData.ServerUUID,
	})

	getServerResponse, err := s.iotStub.Core.Integration.Balena.Server.Get(ctx.Request.Context(), &balena.GetServerRequest{
		Uuid: requestData.ServerUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "balena server not found"})
			return
		}

		err := errors.Join(errors.New("failed to get balena server"), err)
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getServerResponse.Server.Namespace,
			Resources:            []string{"iot.core.integration.balena.server." + getServerResponse.Server.Name},
			Actions:              []string{"iot.core.integration.balena.sync.log.list"},
			NamespaceIndependent: false,
		},
	})
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

	countLogResponse, err := s.iotStub.Core.Integration.Balena.Sync.CountLog(ctx.Request.Context(), &balena.CountSyncLogRequest{
		ServerUUID: requestData.ServerUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "balena server not found"})
			return
		}

		err := errors.Join(errors.New("failed to count balena server sync log"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	listLogResponse, err := s.iotStub.Core.Integration.Balena.Sync.ListLog(ctx.Request.Context(), &balena.ListSyncLogRequest{
		ServerUUID: requestData.ServerUUID,
		Skip:       uint64(requestData.Skip),
		Limit:      uint64(requestData.Limit),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "balena server not found"})
			return
		}

		err := errors.Join(errors.New("failed to list balena server sync log"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	formatedLogEntries := make([]formatedSyncLogEntry, 0, requestData.Limit)
	for {
		logResponse, err := listLogResponse.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err := errors.Join(errors.New("error while receiving balena sync log"), err)
			logger.Error(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		formatedLogEntries = append(formatedLogEntries, formatedSyncLogEntryFromGRPC(logResponse.Log))
	}

	logger.Debug("balena sync log listed")
	ctx.JSON(http.StatusOK, listSyncLogResponse{
		LogEntries: formatedLogEntries,
		TotalCobnt: int64(countLogResponse.TotalCount),
	})
}
