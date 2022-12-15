package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"
)

type UserRouter struct {
	servicesHandler *services.ServicesConnectionHandler
}

type User struct {
	UUID  string `json:"uuid"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type listUsersRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	Skip      int    `json:"skip"`
	Limit     int    `json:"limit"`
}
type listUsersResponse struct {
	TotalCount   int    `json:"totalCount"`
	Users        []User `json:"users"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (r *UserRouter) List(ctx *gin.Context) {
	var requestData listUsersRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
}
