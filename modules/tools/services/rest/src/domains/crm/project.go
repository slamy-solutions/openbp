package crm

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	crm "github.com/slamy-solutions/openbp/modules/crm/libs/golang"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	project "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/project"
)

type projectRouter struct {
	crmStub    *crm.CRMStub
	nativeStub *native.NativeStub

	logger *logrus.Entry
}

type formatedProject struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`

	ClientUUID     string `json:"clientUUID"`
	ContactUUID    string `json:"contactUUID"`
	DepartmentUUID string `json:"departmentUUID"`

	NotRelevant bool `json:"notRelevant"`
}

func formatedProjectFromGRPC(p *project.Project) formatedProject {
	return formatedProject{
		Namespace: p.Namespace,
		UUID:      p.Uuid,
		Name:      p.Name,

		ClientUUID:     p.ClientUUID,
		ContactUUID:    p.ContactUUID,
		DepartmentUUID: p.DepartmentUUID,

		NotRelevant: p.NotRelevant,
	}
}

type createProjectRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name" binding:"required,lt=255,gt=0"`

	ClientUUID     string `json:"clientUUID" binding:"required,lt=64,gt=0"`
	ContactUUID    string `json:"contactUUID" binding:"required,lt=64,gt=0"`
	DepartmentUUID string `json:"departmentUUID" binding:"required,lt=64,gt=0"`

	NotRelevant bool `json:"notRelevant"`
}
type createProjectResponse struct {
	Project formatedProject `json:"project"`
}

func (r *projectRouter) createProject(ctx *gin.Context) {
	var req createProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": req.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.project"},
			Actions:              []string{"crm.project.create"},
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

	p, err := r.crmStub.Core.Project.Create(ctx.Request.Context(), &project.CreateProjectRequest{
		Namespace:      req.Namespace,
		Name:           req.Name,
		ClientUUID:     req.ClientUUID,
		ContactUUID:    req.ContactUUID,
		DepartmentUUID: req.DepartmentUUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.InvalidArgument {
				ctx.JSON(400, gin.H{"error": "ClientUUID, ContactUUID or DepartmentUUID has invalid format"})
				return
			}
		}

		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.WithField("projectUUID", p.Project.Uuid).Info("Project created")
	ctx.JSON(200, createProjectResponse{
		Project: formatedProjectFromGRPC(p.Project),
	})
}

type getProjectsRequest struct {
	Namespace      string `form:"namespace"`
	ClientUUID     string `form:"clientUUID" binding:"required,lt=64,gt=0"`
	DepartmentUUID string `form:"departmentUUID" binding:"required,lt=64,gt=0"`
}
type getProjectsResponse struct {
	Projects []formatedProject `json:"projects"`
}

func (r *projectRouter) getProjects(ctx *gin.Context) {
	var req getProjectsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace":      req.Namespace,
		"clientUUID":     req.ClientUUID,
		"departmentUUID": req.DepartmentUUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.project"},
			Actions:              []string{"crm.project.get"},
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

	projects, err := r.crmStub.Core.Project.GetAll(ctx.Request.Context(), &project.GetAllProjectsRequest{
		Namespace:      req.Namespace,
		ClientUUID:     req.ClientUUID,
		DepartmentUUID: req.DepartmentUUID,
		UseCache:       true,
	})
	if err != nil {
		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	formatedProjects := make([]formatedProject, 0, len(projects.Projects))
	for _, p := range projects.Projects {
		formatedProjects = append(formatedProjects, formatedProjectFromGRPC(p))
	}

	ctx.JSON(200, getProjectsResponse{
		Projects: formatedProjects,
	})
}

type getProjectRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid" binding:"required,lt=64,gt=0"`
}
type getProjectResponse struct {
	Project formatedProject `json:"project"`
}

func (r *projectRouter) getProject(ctx *gin.Context) {
	var req getProjectRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": req.Namespace,
		"uuid":      req.UUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.project"},
			Actions:              []string{"crm.project.get"},
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

	p, err := r.crmStub.Core.Project.Get(ctx.Request.Context(), &project.GetProjectRequest{
		Namespace: req.Namespace,
		Uuid:      req.UUID,
		UseCache:  true,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.JSON(404, gin.H{"error": "project not found"})
				return
			}
		}

		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(200, getProjectResponse{
		Project: formatedProjectFromGRPC(p.Project),
	})
}

type updateProjectRequest struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid" binding:"required,lt=64,gt=0"`
	Name      string `json:"name" binding:"required,lt=255,gt=0"`

	ClientUUID     string `json:"clientUUID" binding:"required,lt=64,gt=0"`
	ContactUUID    string `json:"contactUUID" binding:"required,lt=64,gt=0"`
	DepartmentUUID string `json:"departmentUUID" binding:"required,lt=64,gt=0"`

	NotRelevant bool `json:"notRelevant"`
}
type updateProjectResponse struct {
	Project formatedProject `json:"project"`
}

func (r *projectRouter) updateProject(ctx *gin.Context) {
	var req updateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": req.Namespace,
		"uuid":      req.UUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.project"},
			Actions:              []string{"crm.project.update"},
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

	p, err := r.crmStub.Core.Project.Update(ctx.Request.Context(), &project.UpdateProjectRequest{
		Namespace:      req.Namespace,
		Uuid:           req.UUID,
		Name:           req.Name,
		ClientUUID:     req.ClientUUID,
		ContactUUID:    req.ContactUUID,
		DepartmentUUID: req.DepartmentUUID,
		NotRelevant:    req.NotRelevant,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.JSON(404, gin.H{"error": "project not found"})
				return
			}
			if st.Code() == codes.InvalidArgument {
				ctx.JSON(400, gin.H{"error": "ClientUUID, ContactUUID or DepartmentUUID has invalid format"})
				return
			}
		}

		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.WithField("projectUUID", p.Project.Uuid).Info("Project updated")
	ctx.JSON(200, updateProjectResponse{
		Project: formatedProjectFromGRPC(p.Project),
	})
}

type deleteProjectRequest struct {
	Namespace string `form:"namespace"`
	UUID      string `form:"uuid" binding:"required,lt=64,gt=0"`
}
type deleteProjectResponse struct{}

func (r *projectRouter) deleteProject(ctx *gin.Context) {
	var req deleteProjectRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": req.Namespace,
		"uuid":      req.UUID,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.project"},
			Actions:              []string{"crm.project.delete"},
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

	_, err = r.crmStub.Core.Project.Delete(ctx.Request.Context(), &project.DeleteProjectRequest{
		Namespace: req.Namespace,
		Uuid:      req.UUID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				ctx.JSON(404, gin.H{"error": "project not found"})
				return
			}
		}

		logger.Error(err.Error())
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	logger.Info("Project deleted")
	ctx.JSON(200, deleteProjectResponse{})
}
