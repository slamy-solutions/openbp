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

	department "github.com/slamy-solutions/openbp/modules/crm/libs/golang/core/department"
)

type departmentRouter struct {
	crmStub    *crm.CRMStub
	nativeStub *native.NativeStub

	logger *logrus.Entry
}

type formatedDepartment struct {
	Namespace string `json:"namespace"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
}

func formatedDepartmentFromGRPC(department *department.Department) formatedDepartment {
	return formatedDepartment{
		Namespace: department.Namespace,
		UUID:      department.Uuid,
		Name:      department.Name,
	}
}

type createDepartmentRequest struct {
	Namespace string `json:"namespace" binding:"lte=64"`
	Name      string `json:"name" binding:"required,gt=0,lte=128"`
}
type createDepartmentResponse struct {
	Department formatedDepartment `json:"department"`
}

func (r *departmentRouter) createDepartment(ctx *gin.Context) {
	var req createDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": req.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.department"},
			Actions:              []string{"crm.department.create"},
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

	response, err := r.crmStub.Core.Department.Create(ctx.Request.Context(), &department.CreateDepartmentRequest{
		Namespace: req.Namespace,
		Name:      req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.WithField("department.uuid", response.Department.Uuid).WithField("department.name", req.Name).Info("Department created")
	ctx.JSON(http.StatusOK, createDepartmentResponse{
		Department: formatedDepartmentFromGRPC(response.Department),
	})
}

type getDepartmentsRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
}
type getDepartmentsResponse struct {
	Departments []formatedDepartment `json:"departments"`
}

func (r *departmentRouter) getDepartments(ctx *gin.Context) {
	var req getDepartmentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger := r.logger.WithFields(logrus.Fields{
		"namespace": req.Namespace,
	})

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            req.Namespace,
			Resources:            []string{"crm.department"},
			Actions:              []string{"crm.department.get"},
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

	response, err := r.crmStub.Core.Department.GetAll(ctx.Request.Context(), &department.GetAllDepartmentsRequest{
		Namespace: req.Namespace,
		UseCache:  true,
	})
	if err != nil {
		logger.Error("Failed to get departments: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	formatedDepartments := make([]formatedDepartment, len(response.Departments))
	for i, department := range response.Departments {
		formatedDepartments[i] = formatedDepartmentFromGRPC(department)
	}

	ctx.JSON(http.StatusOK, getDepartmentsResponse{
		Departments: formatedDepartments,
	})
}

type getDepartmentRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}
type getDepartmentResponse struct {
	Department formatedDepartment `json:"department"`
}

func (r *departmentRouter) getDepartment(ctx *gin.Context) {
	var req getDepartmentRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			Resources:            []string{"crm.department"},
			Actions:              []string{"crm.department.get"},
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

	response, err := r.crmStub.Core.Department.Get(ctx.Request.Context(), &department.GetDepartmentRequest{
		Namespace: req.Namespace,
		Uuid:      req.UUID,
		UseCache:  true,
	})
	if err != nil {
		logger.Error("Failed to get department: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, getDepartmentResponse{
		Department: formatedDepartmentFromGRPC(response.Department),
	})
}

type updateDepartmentRequest struct {
	Namespace string `json:"namespace" binding:"lte=64"`
	UUID      string `json:"uuid" binding:"required,lte=64"`
	Name      string `json:"name" binding:"required,gt=0,lte=128"`
}
type updateDepartmentResponse struct {
	Department formatedDepartment `json:"department"`
}

func (r *departmentRouter) updateDepartment(ctx *gin.Context) {
	var req updateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			Resources:            []string{"crm.department"},
			Actions:              []string{"crm.department.update"},
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

	response, err := r.crmStub.Core.Department.Update(ctx.Request.Context(), &department.UpdateDepartmentRequest{
		Namespace: req.Namespace,
		Uuid:      req.UUID,
		Name:      req.Name,
	})
	if err != nil {
		logger.Error("Failed to update department: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updateDepartmentResponse{
		Department: formatedDepartmentFromGRPC(response.Department),
	})
}

type deleteDepartmentRequest struct {
	Namespace string `form:"namespace" binding:"lte=64"`
	UUID      string `form:"uuid" binding:"required,lte=64"`
}
type deleteDepartmentResponse struct {
	Department formatedDepartment `json:"department"`
}

func (r *departmentRouter) deleteDepartment(ctx *gin.Context) {
	var req deleteDepartmentRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			Resources:            []string{"crm.department"},
			Actions:              []string{"crm.department.delete"},
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

	response, err := r.crmStub.Core.Department.Delete(ctx.Request.Context(), &department.DeleteDepartmentRequest{
		Namespace: req.Namespace,
		Uuid:      req.UUID,
	})
	if err != nil {
		logger.Error("Failed to delete department: " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deleteDepartmentResponse{
		Department: formatedDepartmentFromGRPC(response.Department),
	})
}
