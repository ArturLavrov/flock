package peersamplingservice

import (
	"fmt"
	"net/http"
	"strings"

	gossitAgent "flock/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Start(gossipPort string) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	pssNodeName := generatePssNodeName()
	nodes := getActiveNodes()
	agent := gossitAgent.Start(pssNodeName, gossipPort, nodes)

	r.GET("/v1/nodes", func(c *gin.Context) {
		nodes := agent.GetAllNodes()
		c.JSON(http.StatusOK, gin.H{
			"nodes": nodes,
		})
	})

	return r
}

func generatePssNodeName() string {
	shortID := strings.Split(uuid.New().String(), "-")[0]
	return fmt.Sprintf("pss_%s", shortID)
}

func getActiveNodes() []models.Node {
	var nodes []models.Node
	return nodes
}
