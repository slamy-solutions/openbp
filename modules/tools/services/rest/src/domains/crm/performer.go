package crm

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	"github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/performer"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type performerRouter struct {
	crmStub    *crm.CRMStub
	nativeStub *native.NativeStub

	logger *logrus.Entry
}

type formatedPerformer struct {
	Namespace      string `json:"namespace"`
	UUID           string `json:"uuid"`
	DepartmentUUID string `json:"departmentUUID"`
	UserUUID       string `json:"userUUID"`

	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type performerGetAllRequest struct {
	Namespace string `form:"namespace"`
}
type performerGetAllResponse struct {
	Performers []formatedPerformer `json:"performers"`
}

func (r *performerRouter) GetAll(ctx *gin.Context) {
	var request performerGetAllRequest
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
			Resources:            []string{"crm.performer"},
			Actions:              []string{"crm.performer.get"},
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
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	// Get performers
	response, err := r.crmStub.Core.Performer.List(ctx, &performer.ListPerformersRequest{
		Namespace: request.Namespace,
		UseCache:  true,
	})
	if err != nil {
		logger.WithError(err).Error("Failed to get performers")
		ctx.AbortWithError(500, err)
		return
	}

	// Format performers
	performers := make([]formatedPerformer, len(response.Performers))
	for i, performer := range response.Performers {
		performers[i] = formatedPerformer{
			Namespace:      performer.Namespace,
			UUID:           performer.UUID,
			DepartmentUUID: performer.DepartmentUUID,
			UserUUID:       performer.UserUUID,

			Name:      performer.Name,
			AvatarURL: performer.AvatarUrl,
		}
	}

	// Return performers
	ctx.JSON(200, performerGetAllResponse{
		Performers: performers,
	})
}

type createPerformerRequest struct {
	Namespace      string `json:"namespace"`
	DepartmentUUID string `json:"departmentUUID"`
	UserUUID       string `json:"userUUID"`
}
type createPerformerResponse struct {
	Performer formatedPerformer `json:"performer"`
}

func (r *performerRouter) Create(ctx *gin.Context) {
	var request createPerformerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace":      request.Namespace,
		"departmentUUID": request.DepartmentUUID,
		"userUUID":       request.UserUUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            request.Namespace,
			Resources:            []string{"crm.performer"},
			Actions:              []string{"crm.performer.create"},
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
	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	// Create performer
	response, err := r.crmStub.Core.Performer.Create(ctx, &performer.CreatePerformerRequest{
		Namespace:      request.Namespace,
		DepartmentUUID: request.DepartmentUUID,
		UserUUID:       request.UserUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.JSON(400, gin.H{"error": st.Message()})
				return
			}

			if st.Code() == codes.AlreadyExists {
				ctx.JSON(409, gin.H{"error": "Performer already exist."})
				return
			}

			if st.Code() == codes.NotFound {
				ctx.JSON(404, gin.H{"error": "User or department not found."})
				return
			}
		}

		logger.WithError(err).Error("Failed to create performer")
		ctx.AbortWithError(500, err)
		return
	}

	// Format performer
	performer := formatedPerformer{
		Namespace:      response.Performer.Namespace,
		UUID:           response.Performer.UUID,
		DepartmentUUID: response.Performer.DepartmentUUID,
		UserUUID:       response.Performer.UserUUID,

		Name:      response.Performer.Name,
		AvatarURL: response.Performer.AvatarUrl,
	}

	// Return performer
	ctx.JSON(200, createPerformerResponse{
		Performer: performer,
	})
}

type deletePerformerRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid"`
}
type deletePerformerResponse struct{}

func (r *performerRouter) Delete(ctx *gin.Context) {
	var request deletePerformerRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": request.Namespace,
		"UUID":      request.UUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            request.Namespace,
			Resources:            []string{"crm.performer"},
			Actions:              []string{"crm.performer.delete"},
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

	logger = authTools.FillLoggerWithAuthMetadata(logger, authData)

	// Delete performer
	_, err = r.crmStub.Core.Performer.Delete(ctx, &performer.DeletePerformerRequest{
		Namespace: request.Namespace,
		UUID:      request.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.JSON(400, gin.H{"error": st.Message()})
				return
			}

			if st.Code() == codes.NotFound {
				ctx.JSON(404, gin.H{"error": "Performer not found."})
				return
			}
		}

		logger.WithError(err).Error("Failed to delete performer")
		ctx.AbortWithError(500, err)
		return
	}

	// Return performer
	ctx.JSON(200, deletePerformerResponse{})
}
