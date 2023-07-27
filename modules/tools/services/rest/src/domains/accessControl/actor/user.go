package actor

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	native "github.com/slamy-solutions/openbp/modules/native/libs/golang"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/actor/user"
	"github.com/slamy-solutions/openbp/modules/native/libs/golang/iam/auth"
	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/lib/authTools"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserRouter struct {
	nativeStub *native.NativeStub
}

func NewUserRouter(nativeStub *native.NativeStub) *UserRouter {
	return &UserRouter{
		nativeStub: nativeStub,
	}
}

type formatedUser struct {
	Namespace string    `json:"certificate"`
	UUID      string    `json:"uuid"`
	Login     string    `json:"login"`
	Identity  string    `json:"identity"`
	FullName  string    `json:"fullName"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Version   uint64    `json:"version"`
}

func NewFormatedUserFromGRPC(usr *user.User) *formatedUser {
	return &formatedUser{
		Namespace: usr.Namespace,
		UUID:      usr.Uuid,
		Login:     usr.Login,
		Identity:  usr.Identity,
		FullName:  usr.FullName,
		Avatar:    usr.Avatar,
		Email:     usr.Email,
		Created:   usr.Created.AsTime(),
		Updated:   usr.Updated.AsTime(),
		Version:   usr.Version,
	}
}

type listUsersRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	Skip      uint32 `form:"skip" binding:"gte=0"`
	Limit     uint32 `form:"limit" binding:"gte=0,lte=100"`
}
type listUsersResponse struct {
	Users      []formatedUser `json:"users"`
	TotalCount uint64         `json:"totalCount"`
}

func (r *UserRouter) ListUsers(ctx *gin.Context) {
	var requestData listUsersRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.actor.user"},
			Actions:              []string{"native.iam.actor.user.list"},
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

	cursor, err := r.nativeStub.Services.IAM.Actor.User.List(ctx.Request.Context(), &user.ListRequest{
		Namespace: requestData.Namespace,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	foundedUsers := make([]formatedUser, 0, requestData.Limit)
	for {
		userResponse, err := cursor.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		foundedUsers = append(foundedUsers, *NewFormatedUserFromGRPC(userResponse.User))
	}

	countResponse, err := r.nativeStub.Services.IAM.Actor.User.Count(ctx.Request.Context(), &user.CountRequest{
		Namespace: requestData.Namespace,
		UseCache:  true,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &listUsersResponse{Users: foundedUsers, TotalCount: countResponse.Count})
}

type createUserRequest struct {
	Namespace string `json:"namespace" binding:"lte=32"`
	Login     string `json:"login" binding:"lte=32,required,alphanum"`
	FullName  string `json:"fullName" binding:"lte=64"`
	Email     string `json:"email" binding:"lte=128,email"`
}
type createUserResponse struct {
	User formatedUser `json:"user"`
}

func (r *UserRouter) CreateUser(ctx *gin.Context) {
	var requestData createUserRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.actor.user"},
			Actions:              []string{"native.iam.actor.user.create"},
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

	createResponse, err := r.nativeStub.Services.IAM.Actor.User.Create(ctx.Request.Context(), &user.CreateRequest{
		Namespace: requestData.Namespace,
		Login:     requestData.Login,
		FullName:  requestData.FullName,
		Avatar:    "",
		Email:     "",
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.OK {
				ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "User with same login already exist"})
				return
			}
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &createUserResponse{User: *NewFormatedUserFromGRPC(createResponse.User)})
}

type deleteUserRequest struct {
	Namespace string `form:"namespace" binding:"lte=32"`
	UUID      string `form:"uuid" binding:"lte=32,required"`
}

func (r *UserRouter) DeleteUser(ctx *gin.Context) {
	var requestData deleteUserRequest
	if err := ctx.ShouldBind(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	// Check auth
	authData, err := authTools.CheckAuth(ctx, r.nativeStub, []*auth.Scope{
		{
			Namespace:            requestData.Namespace,
			Resources:            []string{"native.iam.actor.user." + requestData.UUID},
			Actions:              []string{"native.iam.actor.user.delete"},
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

	_, err = r.nativeStub.Services.IAM.Actor.User.Delete(ctx.Request.Context(), &user.DeleteRequest{
		Namespace: requestData.Namespace,
		Uuid:      requestData.UUID,
	})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
