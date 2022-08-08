package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAuthGroup(g *gin.RouterGroup) {
	g.GET("ping", ping)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

type loginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func login(c *gin.Context) {
	var requestData loginRequest
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

}

func register(c *gin.Context) {

}

func refreshToken(c *gin.Context) {

}

func validateToken(c *gin.Context) {

}
