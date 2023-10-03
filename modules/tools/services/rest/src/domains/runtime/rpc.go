package runtime

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	runtime "github.com/slamy-solutions/openbp/modules/runtime/libs/golang"
	rpc "github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/rpc"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type RPCRouter struct {
	systemStub  *system.SystemStub
	nativeStub  *native.NativeStub
	runtimeStub *runtime.RuntimeStub

	logger *logrus.Entry
}

type rpcCallRequest struct {
	Namespace   string      `json:"namespace"`
	RuntimeName string      `json:"runtimeName" binding:"required"`
	MethodName  string      `json:"methodName" binding:"required"`
	Data        interface{} `json:"data"`
	Timeout     uint32      `json:"timeout" binding:"required,gt=0"`
}
type rpcCallResponse struct {
	Error        interface{} `json:"error"`
	ErrorMessage string      `json:"errorMessage"`
	Data         interface{} `json:"data"`
}

func (r *RPCRouter) Call(ctx *gin.Context) {
	var requestData rpcCallRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"method.namespace":   requestData.Namespace,
		"method.runtimeName": requestData.RuntimeName,
		"method.name":        requestData.MethodName,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"runtime.manager.rpc.method." + requestData.RuntimeName + "." + requestData.MethodName},
			Actions:              []string{"runtime.manager.rpc.call"},
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

	payload, err := json.Marshal(requestData.Data)
	if err != nil {
		err := errors.New("failed to marshal payload: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Call method
	response, err := r.runtimeStub.Manager.RPC.Call(ctx.Request.Context(), &rpc.CallRequest{
		Namespace:   requestData.Namespace,
		RuntimeName: requestData.RuntimeName,
		MethodName:  requestData.MethodName,
		Payload:     string(payload),
		Timeout:     requestData.Timeout,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.DeadlineExceeded {
			ctx.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{"message": "rpc timeout"})
			return
		}

		err := errors.New("failed to call method: " + err.Error())
		logger.Error(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var responseData interface{}
	if response.Response != "" {
		err = json.Unmarshal([]byte(response.Response), &responseData)
		if err != nil {
			err := errors.New("failed to unmarshal response: " + err.Error())
			logger.Error(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	var errorData interface{}
	if response.Error != "" {
		err = json.Unmarshal([]byte(response.Error), &errorData)
		if err != nil {
			err := errors.New("failed to unmarshal error: " + err.Error())
			logger.Error(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, rpcCallResponse{
		Error:        errorData,
		ErrorMessage: response.ErrorMessage,
		Data:         responseData,
	})
}
