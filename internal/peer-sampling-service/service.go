package peersamplingservice

import (
	"net/http"

	gossitAgent "flock/internal/models"

	"github.com/gin-gonic/gin"
)

func Start(gossipPort string) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	nodes := getActiveNodes()
	agent := gossitAgent.Start("pss", gossipPort, nodes)

	r.GET("/api/v1/nodes", func(c *gin.Context) {
		nodes := agent.GetAllNodes()
		c.JSON(http.StatusOK, gin.H{
			"nodes": nodes,
		})
	})

	return r
}

func getActiveNodes() []gossitAgent.Node {
	var nodes []gossitAgent.Node
	return nodes
}
