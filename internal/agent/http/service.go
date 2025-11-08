package http

import (
	"flock/internal/agent"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start(a agent.Agent) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/v1/nodes", func(c *gin.Context) {
		nodes := a.GetAllNodes()
		c.JSON(http.StatusOK, nodes)
	})

	return r
}
