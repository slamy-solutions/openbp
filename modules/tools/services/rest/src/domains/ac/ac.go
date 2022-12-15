package ac

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/slamy-solutions/openbp/modules/tools/services/rest/src/services"
)

func FillRouterGroup(group *gin.RouterGroup, servicesHandler *services.ServicesConnectionHandler) {
	group.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
