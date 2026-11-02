package api

import (
	"net/http"

	"flock/internal/agent"

	"github.com/gin-gonic/gin"
)

func HTTPHandler(agent *agent.GossipAgent) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/api/v1/nodes", func(c *gin.Context) {
		c.JSON(http.StatusOK, agent.GetNodesByName(c.Query("name")))
	})

	return r
}
