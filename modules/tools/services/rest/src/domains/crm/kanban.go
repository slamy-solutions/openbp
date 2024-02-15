package crm

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/kanban"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type kanbanRouter struct {
	crmStub    *crm.CRMStub
	nativeStub *native.NativeStub

	logger *logrus.Entry
}

type kanbanCreateStageRequest struct {
	Namespace      string `json:"namespace"`
	Name           string `json:"name"`
	DepartmentUUID string `json:"departmentUUID"`
}
type kanbanCreateStageResponse struct {
	Stage *kanban.TicketStage `json:"stage"`
}

func (r *kanbanRouter) CreateStage(ctx *gin.Context) {
	var request kanbanCreateStageRequest
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
			Resources:            []string{"crm.kanban.stage"},
			Actions:              []string{"crm.kanban.stage.create"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Create stage
	response, err := r.crmStub.Core.Kanban.CreateStage(ctx, &kanban.CreateStageRequest{
		Namespace:        request.Namespace,
		Name:             request.Name,
		DepartmentUUID:   request.DepartmentUUID,
		ArrangementIndex: 0,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to create stage")
		ctx.AbortWithError(500, err)
		return
	}

	// Return stage
	ctx.JSON(200, kanbanCreateStageResponse{Stage: response.Stage})
}

type kanbanGetStagesRequest struct {
	Namespace      string `form:"namespace"`
	DepartmentUUID string `form:"departmentUUID"`
}

type kanbanGetStagesResponse struct {
	Stages []*kanban.TicketStage `json:"stages"`
}

func (r *kanbanRouter) GetStages(ctx *gin.Context) {
	var request kanbanGetStagesRequest
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
			Resources:            []string{"crm.kanban.stage"},
			Actions:              []string{"crm.kanban.stage.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Get stages
	response, err := r.crmStub.Core.Kanban.GetStages(ctx, &kanban.GetStagesRequest{
		Namespace:      request.Namespace,
		DepartmentUUID: request.DepartmentUUID,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to get stages")
		ctx.AbortWithError(500, err)
		return
	}

	responseStages := response.Stages
	if responseStages == nil {
		responseStages = make([]*kanban.TicketStage, 0)
	}

	// Return stages
	ctx.JSON(200, kanbanGetStagesResponse{Stages: responseStages})
}

type kanbanUpdateStageRequest struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`

	Name string `json:"name"`
}
type kanbanUpdateStageResponse struct {
	Stage *kanban.TicketStage `json:"stage"`
}

func (r *kanbanRouter) UpdateStage(ctx *gin.Context) {
	var request kanbanUpdateStageRequest
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
			Resources:            []string{"crm.kanban.stage"},
			Actions:              []string{"crm.kanban.stage.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Update stage
	response, err := r.crmStub.Core.Kanban.UpdateStage(ctx, &kanban.UpdateStageRequest{
		Namespace: request.Namespace,
		UUID:      request.UUID,
		Name:      request.Name,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to update stage")
		ctx.AbortWithError(500, err)
		return
	}

	// Return stage
	ctx.JSON(200, kanbanUpdateStageResponse{Stage: response.Stage})
}

type kanbanDeleteStageRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid"`
}
type kanbanDeleteStageResponse struct{}

func (r *kanbanRouter) DeleteStage(ctx *gin.Context) {
	var request kanbanDeleteStageRequest
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
			Resources:            []string{"crm.kanban.stage"},
			Actions:              []string{"crm.kanban.stage.delete"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Delete stage
	_, err = r.crmStub.Core.Kanban.DeleteStage(ctx, &kanban.DeleteStageRequest{
		Namespace: request.Namespace,
		UUID:      request.UUID,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to delete stage")
		ctx.AbortWithError(500, err)
		return
	}

	// Return empty response
	ctx.JSON(200, kanbanDeleteStageResponse{})
}

type kanbanSwapStagesPriorityRequest struct {
	Namespace string `form:"namespace"`
	UUID1     string `form:"uuid1" binding:"required"`
	UUID2     string `form:"uuid2" binding:"required"`
}
type kanbanSwapStagesPriorityResponse struct{}

func (r *kanbanRouter) SwapStagesPriority(ctx *gin.Context) {
	var request kanbanSwapStagesPriorityRequest
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
			Resources:            []string{"crm.kanban.stage"},
			Actions:              []string{"crm.kanban.stage.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Delete stage
	_, err = r.crmStub.Core.Kanban.SwapStagesOrder(ctx, &kanban.SwapStagesOrderRequest{
		Namespace:  request.Namespace,
		StageUUID1: request.UUID1,
		StageUUID2: request.UUID2,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(400, gin.H{"message": "Invalid stage UUIDs"})
				return
			}

			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Stage not found"})
				return
			}
		}

		logger.WithError(err).Error("Failed to swap stage priorities stage")
		ctx.AbortWithError(500, err)
		return
	}

	// Return empty response
	ctx.JSON(200, kanbanSwapStagesPriorityResponse{})
}

type kanbanCreateTicketRequest struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Description string `json:"description"`

	DepartmentUUID    string `json:"departmentUUID"`
	PerformerUUID     string `json:"performerUUID"`
	ClientUUID        string `json:"clientUUID"`
	ContactPersonUUID string `json:"contactPersonUUID"`
	ProjectUUID       string `json:"projectUUID"`

	TrackingStoryPointsPlan uint32 `json:"trackingStoryPointsPlan"`
}
type kanbanCreateTicketResponse struct {
	Ticket *kanban.Ticket `json:"ticket"`
}

func (r *kanbanRouter) CreateTicket(ctx *gin.Context) {
	var request kanbanCreateTicketRequest
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
			Resources:            []string{"crm.kanban.ticket"},
			Actions:              []string{"crm.kanban.ticket.create"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	response, err := r.crmStub.Core.Kanban.CreateTicket(ctx.Request.Context(), &kanban.CreateTicketRequest{
		Namespace:               request.Namespace,
		Name:                    request.Name,
		Description:             request.Description,
		Files:                   []string{},
		Priority:                0,
		ClientUUID:              request.ClientUUID,
		ContactPersonUUID:       request.ContactPersonUUID,
		ProjectUUID:             request.ProjectUUID,
		DepartmentUUID:          request.DepartmentUUID,
		PerformerUUID:           request.PerformerUUID,
		TrackingStoryPointsPlan: request.TrackingStoryPointsPlan,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to create ticket")
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, kanbanCreateTicketResponse{Ticket: response.Ticket})
}

type kanbanDeleteTicketRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid"`
}
type kanbanDeleteTicketResponse struct{}

func (r *kanbanRouter) DeleteTicket(ctx *gin.Context) {
	var request kanbanDeleteTicketRequest
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
			Resources:            []string{"crm.kanban.ticket"},
			Actions:              []string{"crm.kanban.ticket.delete"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	_, err = r.crmStub.Core.Kanban.DeleteTicket(ctx.Request.Context(), &kanban.DeleteTicketRequest{
		Namespace: request.Namespace,
		UUID:      request.UUID,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to delete ticket")
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, kanbanDeleteTicketResponse{})
}

type kanbanGetTicketsRequest struct {
	Namespace      string  `form:"namespace"`
	DepartmentUUID *string `form:"departmentUUID"`
	PerformerUUID  *string `form:"performerUUID"`
}
type kanbanGetTicketsResponse struct {
	Tickets []*kanban.Ticket `json:"tickets"`
}

func (r *kanbanRouter) GetTickets(ctx *gin.Context) {
	var request kanbanGetTicketsRequest
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
			Resources:            []string{"crm.kanban.ticket"},
			Actions:              []string{"crm.kanban.ticket.get"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	var departmentUUID *string = request.DepartmentUUID
	if departmentUUID != nil && *departmentUUID == "" {
		departmentUUID = nil
	}
	var performerUUID *string = request.PerformerUUID
	if performerUUID != nil && *performerUUID == "" {
		performerUUID = nil
	}

	response, err := r.crmStub.Core.Kanban.GetTickets(ctx.Request.Context(), &kanban.GetTicketsRequest{
		Namespace:      request.Namespace,
		DepartmentUUID: departmentUUID,
		PerformerUUID:  performerUUID,
		UseCache:       true,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.AbortWithStatusJSON(400, gin.H{"message": "Invalid departmentUUID or performerUUID"})
				return
			}
		}

		logger.WithError(err).Error("Failed get tickets")
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, kanbanGetTicketsResponse{Tickets: response.Tickets})
}

type kanbanCloseTicketRequest struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
}
type kanbanCloseTicketResponse struct{}

func (r *kanbanRouter) CloseTicket(ctx *gin.Context) {
	var request kanbanCloseTicketRequest
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
			Resources:            []string{"crm.kanban.ticket"},
			Actions:              []string{"crm.kanban.ticket.close"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	_, err = r.crmStub.Core.Kanban.CloseTicket(ctx.Request.Context(), &kanban.CloseTicketRequest{
		Namespace: request.Namespace,
		UUID:      request.UUID,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to close ticket")
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, kanbanCloseTicketResponse{})
}

type kanbanUpdateTicketBasicInfoRequest struct {
	Namespace   string   `json:"namespace"`
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
}
type kanbanUpdateTicketBasicInfoResponse struct {
	Ticket *kanban.Ticket `json:"ticket"`
}

func (r *kanbanRouter) UpdateTicketBasicInfo(ctx *gin.Context) {
	var request kanbanUpdateTicketBasicInfoRequest
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
			Resources:            []string{"crm.kanban.ticket"},
			Actions:              []string{"crm.kanban.ticket.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to check auth")
		ctx.AbortWithError(500, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger = authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	response, err := r.crmStub.Core.Kanban.UpdateTicketBasicInfo(ctx.Request.Context(), &kanban.UpdateTicketBasicInfoRequest{
		Namespace:   request.Namespace,
		UUID:        request.UUID,
		Name:        request.Name,
		Description: request.Description,
		Files:       request.Files,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.AbortWithStatusJSON(404, gin.H{"message": "Ticket not found"})
				return
			}
		}

		logger.WithError(err).Error("Failed to update ticket basic info")
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, kanbanUpdateTicketBasicInfoResponse{
		Ticket: response.Ticket,
	})
}
