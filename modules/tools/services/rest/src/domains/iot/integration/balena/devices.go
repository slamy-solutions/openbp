package balena

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/integration/balena"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DevicesServer struct {
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewDevicesServer(nativeStub *native.NativeStub, iotStub *iot.IOTStub, logger *logrus.Entry) DevicesServer {
	return DevicesServer{
		nativeStub: nativeStub,
		iotStub:    iotStub,
		logger:     logger,
	}
}

type listDevicesInNamespaceRequest struct {
	Namespace     string `form:"namespace"`
	Skip          uint64 `form:"skip" binding:"gte=0"`
	Limit         uint64 `form:"limit" binding:"gt=0,lte=100"`
	BindingFilter string `form:"bindingFilter" binding:"regex=^(all|binded|unbinded)$"`
}
type listDevicesInNamespaceResponse struct {
	Devices    []formatedDevice `json:"devices"`
	TotalCount uint64           `json:"totalCount"`
}

func (s *DevicesServer) ListDevicesInNamespace(ctx *gin.Context) {
	var requestData listDevicesInNamespaceRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.namespace": requestData.Namespace,
		"list.limit":       requestData.Limit,
		"list.skip":        requestData.Skip,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.integration.balena.device"},
			Actions:              []string{"iot.core.integration.balena.device.list"},
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

	countBindingFilter := balena.CountDevicesInNamespaceRequest_ALL
	if requestData.BindingFilter == "binded" {
		countBindingFilter = balena.CountDevicesInNamespaceRequest_ONLY_BINDED
	} else if requestData.BindingFilter == "unbinded" {
		countBindingFilter = balena.CountDevicesInNamespaceRequest_ONLY_UNBINDED
	}

	countResponse, err := s.iotStub.Core.Integration.Balena.Device.CountInNamespace(ctx.Request.Context(), &balena.CountDevicesInNamespaceRequest{
		BalenaServersNamespace: requestData.Namespace,
		BindingFilter:          countBindingFilter,
	})
	if err != nil {
		err := errors.Join(errors.New("error while counting devices in namespace"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	listBindingFilter := balena.ListDevicesInNamespaceRequest_ALL
	if requestData.BindingFilter == "binded" {
		listBindingFilter = balena.ListDevicesInNamespaceRequest_ONLY_BINDED
	} else if requestData.BindingFilter == "unbinded" {
		listBindingFilter = balena.ListDevicesInNamespaceRequest_ONLY_UNBINDED
	}

	listStream, err := s.iotStub.Core.Integration.Balena.Device.ListInNamespace(ctx.Request.Context(), &balena.ListDevicesInNamespaceRequest{
		BalenaServersNamespace: requestData.Namespace,
		Skip:                   requestData.Skip,
		Limit:                  requestData.Limit,
		BindingFilter:          listBindingFilter,
	})
	if err != nil {
		err := errors.Join(errors.New("failed to open devices list stream"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	foundedDevices := make([]formatedDevice, 0, 1)
	for {
		listResponse, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err := errors.Join(errors.New("error while receiving next device from list"), err)
			logger.Error(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		foundedDevices = append(foundedDevices, formatedDeviceFromGRPC(listResponse.Device))
	}

	logger.Debug("Successfully listed balena devices in namespace")
	ctx.JSON(http.StatusOK, listDevicesInNamespaceResponse{
		Devices:    foundedDevices,
		TotalCount: countResponse.Count,
	})
}

type bindDeviceRequest struct {
	DeviceNamespace string `json:"deviceNamespace" binding:"lte=64"`
	DeviceUUID      string `json:"deviceUUID" binding:"lte=64"`

	BalenaDeviceUUID string `json:"balenaDeviceUUID" binding:"lte=64"`
}
type bindDeviceResponse struct{}

func (s *DevicesServer) BindDevice(ctx *gin.Context) {
	var requestData bindDeviceRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"device.namespace":  requestData.DeviceNamespace,
		"device.uuid":       requestData.DeviceUUID,
		"balena.deviceUUID": requestData.BalenaDeviceUUID,
	})

	getBalenaDeviceResponse, err := s.iotStub.Core.Integration.Balena.Device.Get(ctx.Request.Context(), &balena.GetDeviceRequest{
		Uuid: requestData.BalenaDeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Balena device not found"})
			return
		}

		err := errors.Join(errors.New("error while searching for balena device"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logger = logger.WithField("balena.deviceNamespace", getBalenaDeviceResponse.Device.BalenaServerNamespace)

	getIOTDeviceResponse, err := s.iotStub.Core.Device.Get(ctx.Request.Context(), &device.GetRequest{
		Namespace: requestData.DeviceNamespace,
		Uuid:      requestData.BalenaDeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "IoT device not found"})
			return
		}

		err := errors.Join(errors.New("error while searching for iot device"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getBalenaDeviceResponse.Device.BalenaServerNamespace,
			Resources:            []string{"iot.core.integration.balena.device"},
			Actions:              []string{"iot.core.integration.balena.device.bind"},
			NamespaceIndependent: false,
		},
		{
			Namespace:            requestData.DeviceNamespace,
			Resources:            []string{"iot.core.device." + getIOTDeviceResponse.Device.Name},
			Actions:              []string{"iot.core.integration.balena.device.bind"},
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

	_, err = s.iotStub.Core.Integration.Balena.Device.Bind(ctx.Request.Context(), &balena.BindDeviceRequest{
		DeviceNamespace:  requestData.DeviceNamespace,
		DeviceUUID:       requestData.DeviceUUID,
		BalenaDeviceUUID: requestData.BalenaDeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err := errors.Join(errors.New("error while binding device"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, bindDeviceResponse{})
}

type unbindDeviceRequest struct {
	BalenaDeviceUUID string `json:"balenaDeviceUUID" binding:"lte=64"`
}
type unbindDeviceResponse struct{}

func (s *DevicesServer) UnbindDevice(ctx *gin.Context) {
	var requestData unbindDeviceRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.deviceUUID": requestData.BalenaDeviceUUID,
	})

	getBalenaDeviceResponse, err := s.iotStub.Core.Integration.Balena.Device.Get(ctx.Request.Context(), &balena.GetDeviceRequest{
		Uuid: requestData.BalenaDeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Balena device not found"})
			return
		}

		err := errors.Join(errors.New("error while searching for balena device"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logger = logger.WithField("balena.deviceNamespace", getBalenaDeviceResponse.Device.BalenaServerNamespace)

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getBalenaDeviceResponse.Device.BalenaServerNamespace,
			Resources:            []string{"iot.core.integration.balena.device"},
			Actions:              []string{"iot.core.integration.balena.device.bind"},
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

	_, err = s.iotStub.Core.Integration.Balena.Device.UnBind(ctx.Request.Context(), &balena.UnBindDeviceRequest{
		BalenaDeviceUUID: requestData.BalenaDeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err := errors.Join(errors.New("error while unbinding device"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, unbindDeviceResponse{})
}
