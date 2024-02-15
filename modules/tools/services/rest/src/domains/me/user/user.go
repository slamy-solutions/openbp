package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
)

type UserRouter struct {
	logger *logrus.Entry

	nativeStub *native.NativeStub
}

func NewUserRouter(logger *logrus.Entry, nativeStub *native.NativeStub) *UserRouter {
	return &UserRouter{
		logger: logger,

		nativeStub: nativeStub,
	}
}

// type getMyAuthInfoRequest struct{}
type GetMyUserInfoResponse struct {
	User *user.User `json:"user"`
}

func (r *UserRouter) GetMyUserInfo(ctx *gin.Context) {
	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            "",
			Resources:            []string{"me"},
			Actions:              []string{"me.user.get"},
			NamespaceIndependent: true,
		},
	})
	if err != nil {
		err := errors.New("failed to check auth: " + err.Error())
		r.logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if !authData.AccessGranted {
		ctx.AbortWithStatusJSON(authData.StatusCode, gin.H{"message": authData.ErrorMessage})
		return
	}
	logger := authTools.FillLoggerWithAuthMetadata(r.logger, authData)

	// Get user info
	userInfo, err := r.nativeStub.Services.IAM.Actor.User.GetByIdentity(ctx.Request.Context(), &user.GetByIdentityRequest{
		Namespace: authData.Namespace,
		Identity:  authData.IdentityUUID,
		UseCache:  true,
	})
	if err != nil {
		err := errors.New("failed to get user info: " + err.Error())
		logger.Error(err.Error())

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, GetMyUserInfoResponse{
		User: userInfo.User,
	})
}
