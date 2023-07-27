package iot

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceRouter struct {
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewDeviceRouter(logger *logrus.Entry, nativeStub *native.NativeStub, iotStub *iot.IOTStub) *DeviceRouter {
	return &DeviceRouter{
		logger:     logger,
		nativeStub: nativeStub,
		iotStub:    iotStub,
	}
}

type DeviceCreateRequest struct {
	Namespace   string `json:"namespace" binding:"required,lte=64"`
	Name        string `json:"name" binding:"required,lte=64"`
	Description string `json:"description" binding:"required,lte=256"`
}
type DeviceCreateResponse struct {
	Device formatedDevice `json:"device"`
}

func (r *DeviceRouter) Create(ctx *gin.Context) {
	var requestData DeviceCreateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"device.namespace": requestData.Namespace,
		"device.name":      requestData.Name,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.device." + requestData.Name},
			Actions:              []string{"iot.core.device.create"},
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

	createResponse, err := r.iotStub.Core.Device.Create(ctx.Request.Context(), &device.CreateRequest{
		Namespace:   requestData.Namespace,
		Name:        requestData.Name,
		Description: requestData.Description,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
			logger.Debug("Failed to create device. Device with same name already exist")
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Device with same name already exist"})
			return
		}

		err = errors.New("failed to create device: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.WithFields(logrus.Fields{
		"device.uuid":     createResponse.Device.Uuid,
		"device.identity": createResponse.Device.Identity,
	}).Info("Successfully created device")
	ctx.JSON(http.StatusOK, DeviceCreateResponse{Device: formatedDeviceFromGRPC(createResponse.Device)})
}

type DeviceListRequest struct {
	Namespace string `form:"namespace" binding:"required,lte=64"`
	Skip      uint64 `form:"skip" binding:"gte=0"`
	Limit     uint64 `form:"limit" binding:"gt=0,lte=100"`
}
type DeviceListResponse struct {
	Devices    []formatedDevice `json:"devices"`
	TotalCount uint64           `json:"totalCount"`
}

func (r *DeviceRouter) List(ctx *gin.Context) {
	var requestData DeviceListRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"list.namespace": requestData.Namespace,
		"list.skip":      requestData.Skip,
		"list.limit":     requestData.Limit,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.device.*"},
			Actions:              []string{"iot.core.device.get"},
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

	listStream, err := r.iotStub.Core.Device.List(ctx.Request.Context(), &device.ListRequest{
		Namespace: requestData.Namespace,
		Skip:      requestData.Skip,
		Limit:     requestData.Limit,
	})
	if err != nil {
		err = errors.New("failed to open device list stream: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	foundedDevices := make([]formatedDevice, 0, 10)
	for {
		response, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.New("error while getting devices from the listing stream: " + err.Error())
			logger.Error(err.Error())

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		foundedDevices = append(foundedDevices, formatedDeviceFromGRPC(response.Device))
	}

	countResponse, err := r.iotStub.Core.Device.Count(ctx.Request.Context(), &device.CountRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		err = errors.New("failed to count devices in namespace: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Debug("Successfully listed devices")
	ctx.JSON(http.StatusOK, DeviceListResponse{Devices: foundedDevices, TotalCount: countResponse.Count})
}

type DeviceUpdateRequest struct {
	Namespace      string `json:"namespace" binding:"required,lte=64"`
	UUID           string `json:"uuid" binding:"required,lte=32"`
	NewDescription string `json:"newDescription" binding:"required,lte=256"`
}
type DeviceUpdateResponse struct {
	Device formatedDevice `json:"device"`
}

func (r *DeviceRouter) Update(ctx *gin.Context) {
	var requestData DeviceUpdateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"device.namespace": requestData.Namespace,
		"device.uuid":      requestData.UUID,
	})

	// Get the device name to check auth
	getDeviceResponse, err := r.iotStub.Core.Device.Get(ctx.Request.Context(), &device.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err = errors.New("failed to update device: failed to get device information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.device." + getDeviceResponse.Device.Name},
			Actions:              []string{"iot.core.device.update"},
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

	updateResponse, err := r.iotStub.Core.Device.Update(ctx.Request.Context(), &device.UpdateRequest{
		Namespace:   requestData.Namespace,
		Uuid:        requestData.UUID,
		Description: requestData.NewDescription,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			logger.Debug("failed to update device: device not found")
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err = errors.New("failed to update device: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully updated device")
	ctx.JSON(http.StatusOK, DeviceUpdateResponse{Device: formatedDeviceFromGRPC(updateResponse.Device)})
}

type DeviceDeleteRequest struct {
	Namespace string `form:"namespace" binding:"required,lte=64"`
	UUID      string `form:"uuid" binding:"required,lte=32"`
}

func (r *DeviceRouter) Delete(ctx *gin.Context) {
	var requestData DeviceDeleteRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"device.namespace": requestData.Namespace,
		"device.uuid":      requestData.UUID,
	})

	// Get the device name to check auth
	getDeviceResponse, err := r.iotStub.Core.Device.Get(ctx.Request.Context(), &device.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err = errors.New("failed to delete device: failed to get device information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.device." + getDeviceResponse.Device.Name},
			Actions:              []string{"iot.core.device.delete"},
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

	_, err = r.iotStub.Core.Device.Delete(ctx.Request.Context(), &device.DeleteRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		err = errors.New("failed to delete device: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully deleted device")
	ctx.JSON(http.StatusOK, gin.H{})
}
