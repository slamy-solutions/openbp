package bootstrap

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/slamy-solutions/openbp/modules/native/libs/golang/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/authentication/password"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/identity"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/policy"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/keyvaluestorage"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"
)

const (
	ROOT_USER_CREATED_KEY = "TOOLS_REST_BOOTSTRAP_ROOT_CREATED"
)

type RootUserRouter struct {
	servicesHandler *services.ServicesConnectionHandler

	allowInit bool
}

type initRootUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type initRootUserResponse struct{}

func (r *RootUserRouter) InitRootUser(ctx *gin.Context) {
	var requestData initRootUserRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	userExist, err := r.isRootUserInited(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if userExist {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"message": "Root user was already initialized"})
		return
	}

	// Create user
	userResponse, err := r.servicesHandler.Native.ActorUser.Create(ctx.Request.Context(), &user.CreateRequest{
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
	_, err = r.servicesHandler.Native.IAMAuthenticationPassword.CreateOrUpdate(ctx.Request.Context(), &password.CreateOrUpdateRequest{
		Namespace: "",
		Identity:  userResponse.User.Identity,
		Password:  requestData.Password,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Create policy for root with access to entire system
	policyResponse, err := r.servicesHandler.Native.IAMPolicy.Create(ctx.Request.Context(), &policy.CreatePolicyRequest{
		Namespace: "",
		Name:      "root",
		Resources: []string{"*"},
		Actions:   []string{"*"},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Assign policy to the user identity
	_, err = r.servicesHandler.Native.IAMIdentity.AddPolicy(ctx.Request.Context(), &identity.AddPolicyRequest{
		IdentityNamespace: "",
		IdentityUUID:      userResponse.User.Identity,
		PolicyNamespace:   "",
		PolicyUUID:        policyResponse.Policy.Uuid,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Set ROOT_USER_CREATED_KEY to identify that root user was created
	_, err = r.servicesHandler.Native.KeyValueStorage.Set(ctx.Request.Context(), &keyvaluestorage.SetRequest{
		Namespace: "",
		Key:       ROOT_USER_CREATED_KEY,
		Value:     []byte{},
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Root user successfully initialized"})
}

func (r *RootUserRouter) isRootUserInited(ctx context.Context) (bool, error) {
	if !r.allowInit {
		return true, nil
	}

	// Check if root user doesnt exist
	// TODO: change to 'Exist' instead of 'Get'
	_, err := r.servicesHandler.Native.KeyValueStorage.Get(ctx, &keyvaluestorage.GetRequest{
		Namespace: "",
		Key:       ROOT_USER_CREATED_KEY,
		UseCache:  true,
	})
	if err != nil {
		if st, ok := status.FromError(err); !ok || st.Code() != codes.NotFound {
			return false, err
		}
	} else {
		// Root user was already created
		return true, nil
	}

	//TODO: check list of users. If there are other users => return that already exist

	return false, nil
}
