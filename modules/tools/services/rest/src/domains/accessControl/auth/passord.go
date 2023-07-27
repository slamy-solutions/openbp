package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PasswordRouter struct {
	nativeStub *native.NativeStub
}

func NewPasswordRouter(nativeStub *native.NativeStub) *PasswordRouter {
	return &PasswordRouter{
		nativeStub: nativeStub,
	}
}

type getStatusRequest struct {
	Namespace    string `form:"namespace" binding:"lte=32"`
	IdentityUUID string `form:"identityUUID" binding:"lte=64"`
}

type getStatusResponse struct {
	Seted bool `json:"seted"`
}

func (r *PasswordRouter) GetStatus(ctx *gin.Context) {
	var requestData getStatusRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.authentication.password.*"},
			Actions:              []string{"native.iam.authentication.password.status"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	existsResponse, err := r.nativeStub.Services.IAM.Authentication.Password.Exists(ctx.Request.Context(), &password.ExistsRequest{
		Namespace: requestData.Namespace,
		Identity:  requestData.IdentityUUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &getStatusResponse{Seted: existsResponse.Exists})
}

type disableRequest struct {
	Namespace    string `form:"namespace" binding:"lte=32"`
	IdentityUUID string `form:"identityUUID" binding:"lte=64"`
}

func (r *PasswordRouter) Disable(ctx *gin.Context) {
	var requestData disableRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.authentication.password.*"},
			Actions:              []string{"native.iam.authentication.password.disable"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	_, err = r.nativeStub.Services.IAM.Authentication.Password.Delete(ctx.Request.Context(), &password.DeleteRequest{
		Namespace: requestData.Namespace,
		Identity:  requestData.IdentityUUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type updateRequest struct {
	Namespace    string `json:"namespace" binding:"lte=32"`
	IdentityUUID string `json:"identityUUID" binding:"lte=64"`
	NewPassword  string `json:"newPassword" binding:"lte=64"`
}

func (r *PasswordRouter) SetOrUpdate(ctx *gin.Context) {
	var requestData updateRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.authentication.password.*"},
			Actions:              []string{"native.iam.authentication.password.update"},
			NamespaceIndependent: false,
		},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}

	_, err = r.nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx.Request.Context(), &password.CreateOrUpdateRequest{
		Namespace: requestData.Namespace,
		Identity:  requestData.IdentityUUID,
		Password:  requestData.NewPassword,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.FailedPrecondition {
				ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "The vault is sealed or namespace doesnt exist."})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
