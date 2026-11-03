package api

import (
	"net/http"

	"flock/internal/agent"
	"flock/internal/cluster"

	"github.com/gin-gonic/gin"
)

func HTTPHandler(agent *agent.GossipAgent) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/api/v1/nodes", func(c *gin.Context) {
		nodes := agent.GetNodesByName(c.Query("name"))
		if nodes == nil {
			nodes = []cluster.Node{}
		}
		c.JSON(http.StatusOK, gin.H{"nodes": nodes})
	})

	return r
}
