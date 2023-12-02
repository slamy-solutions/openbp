package crm

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	crm_onec_sync_grpc "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/onecsync"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type onecsyncRouter struct {
	nativeStub *native.NativeStub
	crmStub    *crm.CRMStub

	logger *logrus.Entry
}

type formatedOneCSyncLogEvent struct {
	UUID         string    `json:"uuid"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"errorMessage"`
	Timestamp    time.Time `json:"timestamp"`
	Log          string    `json:"log"`
}

type oneCSyncSyncNowRequest struct {
	Namespace string `json:"namespace"`
}
type oneCSyncSyncNowResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

func (r *onecsyncRouter) syncNow(ctx *gin.Context) {
	var request oneCSyncSyncNowRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": request.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            request.Namespace,
			Resources:            []string{"crm.onec"},
			Actions:              []string{"crm.onec.sync"},
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
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	syncResponse, err := r.crmStub.Core.OneCSync.Sync(ctx.Request.Context(), &crm_onec_sync_grpc.SyncRequest{
		Namespace: request.Namespace,
	})
	if err != nil {
		err := errors.New("failed to sync: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, oneCSyncSyncNowResponse{
		Success:      syncResponse.Ok,
		ErrorMessage: syncResponse.ErrorMessage,
	})
}

type oneCSyncGetLogRequest struct {
	Namespace string `form:"namespace"`
	Skip      int32  `form:"skip" binding:"gte=0"`
	Limit     int32  `form:"limit" binding:"gte=0,lte=100"`
}
type oneCSyncGetLogResponse struct {
	Events     []formatedOneCSyncLogEvent `json:"events"`
	TotalCount int32                      `json:"totalCount"`
}

func (r *onecsyncRouter) getLog(ctx *gin.Context) {
	var request oneCSyncGetLogRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": request.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            request.Namespace,
			Resources:            []string{"crm.onec.synclog"},
			Actions:              []string{"crm.onec.synclog.get"},
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
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	logResponse, err := r.crmStub.Core.OneCSync.GetLog(ctx.Request.Context(), &crm_onec_sync_grpc.GetLogRequest{
		Namespace: request.Namespace,
		Skip:      request.Skip,
		Limit:     request.Limit,
	})
	if err != nil {
		err := errors.New("failed to get log: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithError(500, err)
		return
	}

	formatedEvents := make([]formatedOneCSyncLogEvent, len(logResponse.Events))
	for i, event := range logResponse.Events {
		formatedEvents[i] = formatedOneCSyncLogEvent{
			UUID:         event.Uuid,
			Success:      event.Success,
			ErrorMessage: event.Error,
			Timestamp:    event.Timestamp.AsTime(),
			Log:          event.Log,
		}
	}

	ctx.JSON(200, oneCSyncGetLogResponse{
		Events:     formatedEvents,
		TotalCount: logResponse.TotalCount,
	})
}
