package iot

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	iot "github.com/slamy-solutions/openbp/modules/iot/libs/golang"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/device"
	"github.com/slamy-solutions/openbp/modules/iot/libs/golang/core/fleet"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FleetRouter struct {
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewFleetRouter(logger *logrus.Entry, nativeStub *native.NativeStub, iotStub *iot.IOTStub) *FleetRouter {
	return &FleetRouter{
		logger:     logger,
		nativeStub: nativeStub,
		iotStub:    iotStub,
	}
}

type FleetCreateRequest struct {
	Namespace   string `json:"namespace" binding:"lte=64"`
	Name        string `json:"name" binding:"required,lte=64"`
	Description string `json:"description" binding:"required,lte=256"`
}
type FleetCreateResponse struct {
	Fleet formatedFleet `json:"fleet"`
}

func (r *FleetRouter) Create(ctx *gin.Context) {
	var requestData FleetCreateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.name":      requestData.Name,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + requestData.Name},
			Actions:              []string{"iot.core.fleet.create"},
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

	createResponse, err := r.iotStub.Core.Fleet.Create(ctx.Request.Context(), &fleet.CreateRequest{
		Namespace:   requestData.Namespace,
		Name:        requestData.Name,
		Description: requestData.Description,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.AlreadyExists {
			logger.Debug("Failed to create fleet. Fleet with same name already exist")
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Fleet with same name already exist"})
			return
		}

		err = errors.New("failed to create fleet: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.WithFields(logrus.Fields{
		"fleet.uuid": createResponse.Fleet.Uuid,
	}).Info("Successfully created fleet")
	ctx.JSON(http.StatusOK, FleetCreateResponse{Fleet: formatedFleetFromGRPC(createResponse.Fleet)})
}

type FleetGetRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}
type FleetGetResponse struct {
	Fleet formatedFleet `json:"fleet"`
}

func (r *FleetRouter) Get(ctx *gin.Context) {
	var requestData FleetGetRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.uuid":      requestData.UUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + requestData.UUID},
			Actions:              []string{"iot.core.fleet.get"},
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

	getResponse, err := r.iotStub.Core.Fleet.Get(ctx.Request.Context(), &fleet.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			logger.Debug("Failed to get fleet. Fleet not found")
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to get fleet: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Debug("Successfully got fleet")
	ctx.JSON(http.StatusOK, FleetGetResponse{Fleet: formatedFleetFromGRPC(getResponse.Fleet)})
}

type FleetListRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	Skip      uint64 `form:"skip" binding:"gte=0"`
	Limit     uint64 `form:"limit" binding:"gt=0,lte=100"`
}
type FleetListResponse struct {
	Fleets     []formatedFleet `json:"fleets"`
	TotalCount uint64          `json:"totalCount"`
}

func (r *FleetRouter) List(ctx *gin.Context) {
	var requestData FleetListRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"params.skip":     requestData.Skip,
		"params.limit":    requestData.Limit,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet.*"},
			Actions:              []string{"iot.core.fleet.get"},
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

	listStream, err := r.iotStub.Core.Fleet.List(ctx.Request.Context(), &fleet.ListRequest{
		Namespace: requestData.Namespace,
		Skip:      requestData.Skip,
		Limit:     requestData.Limit,
	})
	if err != nil {
		err = errors.New("failed to open fleet list stream: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	foundedFleets := make([]formatedFleet, 0, 10)
	for {
		response, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err = errors.New("error while getting fleets from the listing stream: " + err.Error())
			logger.Error(err.Error())

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		foundedFleets = append(foundedFleets, formatedFleetFromGRPC(response.Fleet))
	}

	countResponse, err := r.iotStub.Core.Fleet.Count(ctx.Request.Context(), &fleet.CountRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		err = errors.New("failed to count fleets in namespace: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Debug("Successfully listed fleets")
	ctx.JSON(http.StatusOK, FleetListResponse{Fleets: foundedFleets, TotalCount: countResponse.Count})
}

type FleetUpdateRequest struct {
	Namespace      string `json:"namespace" binding:"lte=64"`
	UUID           string `json:"uuid" binding:"required,lte=32"`
	NewDescription string `json:"newDescription" binding:"required,lte=256"`
}
type FleetUpdateResponse struct {
	Fleet formatedFleet `json:"fleet"`
}

func (r *FleetRouter) Update(ctx *gin.Context) {
	var requestData FleetUpdateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.uuid":      requestData.UUID,
	})

	// Get the fleet name to check auth
	getFleetResponse, err := r.iotStub.Core.Fleet.Get(ctx.Request.Context(), &fleet.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to update fleet: failed to get fleet information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + getFleetResponse.Fleet.Name},
			Actions:              []string{"iot.core.fleet.update"},
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

	updateResponse, err := r.iotStub.Core.Fleet.Update(ctx.Request.Context(), &fleet.UpdateRequest{
		Namespace:   requestData.Namespace,
		Uuid:        requestData.UUID,
		Description: requestData.NewDescription,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			logger.Debug("failed to update fleet: fleet not found")
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to update fleet: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully updated fleet")
	ctx.JSON(http.StatusOK, FleetUpdateResponse{Fleet: formatedFleetFromGRPC(updateResponse.Fleet)})
}

type FleetDeleteRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	UUID      string `form:"uuid" binding:"required,lte=32"`
}

func (r *FleetRouter) Delete(ctx *gin.Context) {
	var requestData FleetDeleteRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.uuid":      requestData.UUID,
	})

	// Get the fleet name to check auth
	getFleetResponse, err := r.iotStub.Core.Fleet.Get(ctx.Request.Context(), &fleet.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to delete fleet: failed to get fleet information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + getFleetResponse.Fleet.Name},
			Actions:              []string{"iot.core.fleet.delete"},
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

	_, err = r.iotStub.Core.Fleet.Delete(ctx.Request.Context(), &fleet.DeleteRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		err = errors.New("failed to delete fleet: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully deleted fleet")
	ctx.JSON(http.StatusOK, gin.H{})
}

type FleetListDevicesRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	UUID      string `form:"uuid" binding:"required,lte=32"`
	Skip      uint64 `form:"skip" binding:"gte=0"`
	Limit     uint64 `form:"limit" binding:"gt=0,lte=100"`
}
type FleetListDevicesResponse struct {
	Devices    []formatedDevice `json:"devices"`
	TotalCount uint64           `json:"totalCount"`
}

func (r *FleetRouter) ListDevices(ctx *gin.Context) {
	var requestData FleetListDevicesRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.uuid":      requestData.UUID,
		"params.skip":     requestData.Skip,
		"params.limit":    requestData.Limit,
	})

	// Get the fleet name to check auth
	getFleetResponse, err := r.iotStub.Core.Fleet.Get(ctx.Request.Context(), &fleet.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to list fleet devices: failed to get fleet information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + getFleetResponse.Fleet.Name},
			Actions:              []string{"iot.core.fleet.listDevices"},
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

	listStream, err := r.iotStub.Core.Fleet.ListDevices(ctx.Request.Context(), &fleet.ListDevicesRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
		Skip:      requestData.Skip,
		Limit:     requestData.Limit,
	})
	if err != nil {
		err = errors.New("failed to open fleet devices list stream: " + err.Error())
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

			err = errors.New("error while getting fleets from the listing stream: " + err.Error())
			logger.Error(err.Error())

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		foundedDevices = append(foundedDevices, formatedDeviceFromGRPC(response.Device.Device))
	}

	countResponse, err := r.iotStub.Core.Fleet.CountDevices(ctx.Request.Context(), &fleet.CountDevicesRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		err = errors.New("failed to count fleet devices: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Debug("Successfully listed fleet devices")
	ctx.JSON(http.StatusOK, FleetListDevicesResponse{Devices: foundedDevices, TotalCount: countResponse.Count})
}

type FleetAddDeviceRequest struct {
	Namespace  string `json:"namespace" binding:"lte=64"`
	FleetUUID  string `json:"fleetUUID" binding:"required,lte=32"`
	DeviceUUID string `json:"deviceUUID" binding:"required,lte=32"`
}
type FleetAddDeviceResponse struct{}

func (r *FleetRouter) AddDevice(ctx *gin.Context) {
	var requestData FleetAddDeviceRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.uuid":      requestData.FleetUUID,
		"device.uuid":     requestData.DeviceUUID,
	})

	// Get the fleet name to check auth
	getFleetResponse, err := r.iotStub.Core.Fleet.Get(ctx.Request.Context(), &fleet.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.FleetUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to add device to the fleet: failed to get fleet information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get the device name to check auth
	getDeviceResponse, err := r.iotStub.Core.Device.Get(ctx.Request.Context(), &device.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.DeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err = errors.New("failed to add device to the fleet: failed to get device information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + getFleetResponse.Fleet.Name, "iot.core.device." + getDeviceResponse.Device.Name},
			Actions:              []string{"iot.core.fleet.addDevice"},
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

	_, err = r.iotStub.Core.Fleet.AddDevice(ctx.Request.Context(), &fleet.AddDeviceRequest{
		Namespace:  requestData.Namespace,
		FleetUUID:  requestData.FleetUUID,
		DeviceUUID: requestData.DeviceUUID,
	})
	if err != nil {
		err = errors.New("failed to add device to the fleet: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully added device to the fleet")
	ctx.JSON(http.StatusOK, gin.H{})
}

type FleetRemoveDeviceRequest struct {
	Namespace  string `form:"namespace" binding:"lte=64"`
	FleetUUID  string `form:"fleetUUID" binding:"required,lte=32"`
	DeviceUUID string `form:"deviceUUID" binding:"required,lte=32"`
}
type FleetRemoveDeviceResponse struct{}

func (r *FleetRouter) RemoveDevice(ctx *gin.Context) {
	var requestData FleetRemoveDeviceRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"fleet.namespace": requestData.Namespace,
		"fleet.uuid":      requestData.FleetUUID,
		"device.uuid":     requestData.DeviceUUID,
	})

	// Get the fleet name to check auth
	getFleetResponse, err := r.iotStub.Core.Fleet.Get(ctx.Request.Context(), &fleet.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.FleetUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Fleet not found"})
			return
		}

		err = errors.New("failed to remove device from the fleet: failed to get fleet information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get the device name to check auth
	getDeviceResponse, err := r.iotStub.Core.Device.Get(ctx.Request.Context(), &device.GetRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.DeviceUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		err = errors.New("failed to remove device from the fleet: failed to get device information: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.fleet." + getFleetResponse.Fleet.Name, "iot.core.device." + getDeviceResponse.Device.Name},
			Actions:              []string{"iot.core.fleet.removeDevice"},
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

	_, err = r.iotStub.Core.Fleet.AddDevice(ctx.Request.Context(), &fleet.AddDeviceRequest{
		Namespace:  requestData.Namespace,
		FleetUUID:  requestData.FleetUUID,
		DeviceUUID: requestData.DeviceUUID,
	})
	if err != nil {
		err = errors.New("failed to remove device from the fleet: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully removed device from the fleet")
	ctx.JSON(http.StatusOK, gin.H{})
}
