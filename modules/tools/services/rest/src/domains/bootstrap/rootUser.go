package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"

	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/role"
	system "github.com/slamy-solutions/openbp/modules/system/libs/golang"
)

const (
	ROOT_USER_CREATED_KEY = "TOOLS_REST_BOOTSTRAP_ROOT_CREATED"
)

type RootUserRouter struct {
	nativeStub *native.NativeStub
	systemStub *system.SystemStub
}

type initRootUserRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (r *RootUserRouter) InitRootUser(ctx *gin.Context) {
	var requestData initRootUserRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	vaultSealed, err := isVaultSealed(ctx.Request.Context(), r.systemStub)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if vaultSealed {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "The vault is sealed."})
		return
	}

	blocked := isRootUserCreationBlocked()
	if blocked {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Root user creation is blocked."})
		return
	}

	userExist, err := isRootUserCreated(ctx.Request.Context(), r.nativeStub)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if userExist {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Root user was already initialized"})
		return
	}

	// Create user
	userResponse, err := r.nativeStub.Services.IAM.Actor.User.Create(ctx.Request.Context(), &user.CreateRequest{
		Login:    requestData.Login,
		FullName: "Root",
		Avatar:   "",
		Email:    "",
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Assign password to user
	_, err = r.nativeStub.Services.IAM.Authentication.Password.CreateOrUpdate(ctx.Request.Context(), &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  userResponse.User.Identity,
		Password:  requestData.Password,
	})
	if err != nil {
		r.nativeStub.Services.IAM.Actor.User.Delete(ctx.Request.Context(), &user.DeleteRequest{Namespace: "", Uuid: userResponse.User.Uuid})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Get Root Role
	roleResponse, err := r.nativeStub.Services.IAM.Role.GetBuiltInRole(ctx.Request.Context(), &role.GetBuiltInRoleRequest{
		Namespace: "",
		Type:      role.BuiltInRoleType_GLOBAL_ROOT,
	})
	if err != nil {
		r.nativeStub.Services.IAM.Authentication.Password.Delete(ctx.Request.Context(), &password.DeleteRequest{Namespace: "", Identity: userResponse.User.Identity})
		r.nativeStub.Services.IAM.Actor.User.Delete(ctx.Request.Context(), &user.DeleteRequest{Namespace: "", Uuid: userResponse.User.Uuid})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Assign role to the user identity
	_, err = r.nativeStub.Services.IAM.Identity.AddRole(ctx.Request.Context(), &identity.AddRoleRequest{
		IdentityNamespace: "",
		IdentityUUID:      userResponse.User.Identity,
		RoleNamespace:     "",
		RoleUUID:          roleResponse.Role.Uuid,
	})
	if err != nil {
		r.nativeStub.Services.IAM.Authentication.Password.Delete(ctx, &password.DeleteRequest{Namespace: "", Identity: userResponse.User.Identity})
		r.nativeStub.Services.IAM.Actor.User.Delete(ctx, &user.DeleteRequest{Namespace: "", Uuid: userResponse.User.Uuid})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Root user successfully initialized"})
}
