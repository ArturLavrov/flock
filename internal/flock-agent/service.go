package agent

import (
	"flock/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HTTPHandler(a models.Agent) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/v1/nodes", func(c *gin.Context) {
		c.JSON(http.StatusOK, a.GetAllNodes())
	})

	return r
}
