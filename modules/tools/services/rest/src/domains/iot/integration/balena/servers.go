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

type ServersServer struct {
	nativeStub *native.NativeStub
	iotStub    *iot.IOTStub

	logger *logrus.Entry
}

func NewServersServer(nativeStub *native.NativeStub, iotStub *iot.IOTStub, logger *logrus.Entry) ServersServer {
	return ServersServer{
		nativeStub: nativeStub,
		iotStub:    iotStub,
		logger:     logger,
	}
}

type createServerRequest struct {
	Namespace   string `json:"namespace" binding:"lte=64"`
	Name        string `json:"name" binding:"lte=64"`
	Description string `json:"description" binding:"lte=256"`
	BaseURL     string `json:"baseURL" binding:"lte=256"`
	AuthToken   string `json:"authToken" binding:"lte=1024"`
}
type createServerResponse struct {
	Server formatedServer `json:"server"`
}

func (s *ServersServer) CreateServer(ctx *gin.Context) {
	var requestData createServerRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.namespace": requestData.Namespace,
		"balena.baseURL":   requestData.BaseURL,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"iot.core.integration.balena.server." + requestData.Name},
			Actions:              []string{"iot.core.integration.balena.server.create"},
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

	createResponse, err := s.iotStub.Core.Integration.Balena.Server.Create(ctx.Request.Context(), &balena.CreateServerRequest{
		Namespace:   requestData.Namespace,
		Name:        requestData.Name,
		Description: requestData.Description,
		BaseURL:     requestData.BaseURL,
		AuthToken:   requestData.AuthToken,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "The vault is sealed."})
				return
			}
		}

		err := errors.Join(errors.New("error while trying to create balena server"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.WithField("balena.serverUUID", createResponse.Server.Uuid).Info("Successfully created (registered) new balena server")

	ctx.JSON(http.StatusOK, createServerResponse{
		Server: formatedServerFromGRPC(createResponse.Server),
	})
}

type listServersInNamespaceRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	Skip      uint64 `form:"skip" binding:"gte=0"`
	Limit     uint64 `form:"limit" binding:"gt=0,lte=100"`
}
type listServersInNamespaceResponse struct {
	Servers    []formatedServerWithMetadata `json:"servers"`
	TotalCount uint64                       `json:"totalCount"`
}

func (s *ServersServer) ListServersInNamespace(ctx *gin.Context) {
	var requestData listServersInNamespaceRequest
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
			Resources:            []string{"iot.core.integration.balena.server.*"},
			Actions:              []string{"iot.core.integration.balena.server.get"},
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

	countResponse, err := s.iotStub.Core.Integration.Balena.Server.CountInNamespace(ctx.Request.Context(), &balena.CountServersInNamespaceRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		err := errors.Join(errors.New("error while counting servers in namespace"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	listStream, err := s.iotStub.Core.Integration.Balena.Server.ListInNamespace(ctx.Request.Context(), &balena.ListServersInNamespaceRequest{
		Namespace: requestData.Namespace,
		Skip:      requestData.Skip,
		Limit:     requestData.Limit,
	})
	if err != nil {
		err := errors.Join(errors.New("failed to open servers list stream"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	foundedServer := make([]formatedServerWithMetadata, 0, 1)
	for {
		listResponse, err := listStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			err := errors.Join(errors.New("error while receiving next server from list"), err)
			logger.Error(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var lastSyncLog *formatedSyncLogEntry = nil
		lastLogResponse, err := s.iotStub.Core.Integration.Balena.Sync.GetLastSyncLog(ctx.Request.Context(), &balena.GetLastSyncLogRequest{
			ServerUUID: listResponse.Server.Uuid,
		})
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
				lastSyncLog = nil
			} else {
				err := errors.Join(errors.New("error while trying to get last sync log"), err)
				logger.Error(err.Error())
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		} else {
			l := formatedSyncLogEntryFromGRPC(lastLogResponse.Log)
			lastSyncLog = &l
		}

		foundedServer = append(foundedServer, formatedServerWithMetadata{
			Server:      formatedServerFromGRPC(listResponse.Server),
			LastSyncLog: lastSyncLog,
		})
	}

	logger.Debug("Successfully listed balena servers in namespace")
	ctx.JSON(http.StatusOK, listServersInNamespaceResponse{
		Servers:    foundedServer,
		TotalCount: countResponse.TotalCount,
	})
}

type updateServerRequest struct {
	UUID           string `json:"uuid" binding:"lte=64"`
	NewDescription string `json:"newDescription" binding:"lte=256"`
	NewBaseURL     string `json:"newBaseURL" binding:"lte=256"`
	NewAuthToken   string `json:"newAuthToken" binding:"lte=1024"`
}
type updateServerResponse struct {
	Server formatedServer `json:"server"`
}

func (s *ServersServer) UpdateServer(ctx *gin.Context) {
	var requestData updateServerRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.serverUUID": requestData.UUID,
		"balena.newBaseURL": requestData.NewBaseURL,
	})

	// Get server to check auth
	getResponse, err := s.iotStub.Core.Integration.Balena.Server.Get(ctx.Request.Context(), &balena.GetServerRequest{
		Uuid: requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Server not found"})
				return
			}
		}

		err := errors.Join(errors.New("failed to get server information"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger = s.logger.WithFields(logrus.Fields{
		"balena.serverNamespace": getResponse.Server.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getResponse.Server.Namespace,
			Resources:            []string{"iot.core.integration.balena.server." + getResponse.Server.Name},
			Actions:              []string{"iot.core.integration.balena.server.update"},
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

	updateResponse, err := s.iotStub.Core.Integration.Balena.Server.Update(ctx.Request.Context(), &balena.UpdateServerRequest{
		ServerUUID:     requestData.UUID,
		NewDescription: requestData.NewDescription,
		NewBaseURL:     requestData.NewBaseURL,
		NewAuthToken:   requestData.NewAuthToken,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Server not found"})
				return
			}
		}

		err := errors.Join(errors.New("failed to update server information"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully updated server")
	ctx.JSON(http.StatusOK, updateServerResponse{Server: formatedServerFromGRPC(updateResponse.Server)})
}

type deleteServerRequest struct {
	UUID string `form:"uuid" binding:"lte=64"`
}
type deleteServerResponse struct{}

func (s *ServersServer) DeleteServer(ctx *gin.Context) {
	var requestData deleteServerRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.serverUUID": requestData.UUID,
	})

	// Get server to check auth
	getResponse, err := s.iotStub.Core.Integration.Balena.Server.Get(ctx.Request.Context(), &balena.GetServerRequest{
		Uuid: requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Server not found"})
				return
			}
		}

		err := errors.Join(errors.New("failed to get server information"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger = s.logger.WithFields(logrus.Fields{
		"balena.serverNamespace": getResponse.Server.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getResponse.Server.Namespace,
			Resources:            []string{"iot.core.integration.balena.server." + getResponse.Server.Name},
			Actions:              []string{"iot.core.integration.balena.server.delete"},
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

	_, err = s.iotStub.Core.Integration.Balena.Server.Delete(ctx.Request.Context(), &balena.DeleteServerRequest{
		ServerUUID: requestData.UUID,
	})
	if err != nil {
		err := errors.Join(errors.New("failed to delete server"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully deleted server")
	ctx.JSON(http.StatusOK, deleteServerResponse{})
}

type setServerEnabledRequest struct {
	UUID    string `json:"uuid" binding:"lte=64"`
	Enabled bool   `json:"enabled"`
}
type setServerEnabledResponse struct {
	Server formatedServer `json:"server"`
}

func (s *ServersServer) SetServerEnabled(ctx *gin.Context) {
	var requestData setServerEnabledRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	logger := s.logger.WithFields(logrus.Fields{
		"balena.serverUUID":       requestData.UUID,
		"balena.serverNewEnabled": requestData.Enabled,
	})

	// Get server to check auth
	getResponse, err := s.iotStub.Core.Integration.Balena.Server.Get(ctx.Request.Context(), &balena.GetServerRequest{
		Uuid: requestData.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Server not found"})
				return
			}
		}

		err := errors.Join(errors.New("failed to get server information"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger = s.logger.WithFields(logrus.Fields{
		"balena.serverNamespace":  getResponse.Server.Namespace,
		"balena.serverOldEnabled": getResponse.Server.Enabled,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, s.nativeStub, []*auth.Scope{
		{
			Namespace:            getResponse.Server.Namespace,
			Resources:            []string{"iot.core.integration.balena.server." + getResponse.Server.Name},
			Actions:              []string{"iot.core.integration.balena.server.setEnabled"},
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

	setEnabledResponse, err := s.iotStub.Core.Integration.Balena.Server.SetEnabled(ctx.Request.Context(), &balena.SetServerEnabledRequest{
		ServerUUID: requestData.UUID,
		Enabled:    requestData.Enabled,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Server not found"})
				return
			}
		}

		err := errors.Join(errors.New("failed to change enabling state of the server"), err)
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Successfully changed server enabling state")
	ctx.JSON(http.StatusOK, setServerEnabledResponse{formatedServerFromGRPC(setEnabledResponse.Server)})
}
