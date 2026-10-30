package peersamplingservice

import (
	"log"
	"net/http"
	"strconv"

	gossitAgent "flock/internal/models"

	"github.com/gin-gonic/gin"
)

func Start(gossipPort string) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	nodes := getActiveNodes()
	port, err := strconv.Atoi(gossipPort)
	if err != nil {
		log.Fatalf("Invalid gossip port: %v", err)
	}
	agentNode := gossitAgent.NewNode("pss", port)
	agent := gossitAgent.Start(agentNode, nodes)

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
