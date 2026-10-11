package peersamplingservice

import (
	"fmt"
	"net/http"
	"strings"

	"flock/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Start(gossipPort string) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	pssNodeName := generatePssNodeName()
	peers := getPeers()
	agent := models.Start(pssNodeName, gossipPort, peers)

	r.GET("/v1/peers", func(c *gin.Context) {
		peers := agent.GetAllPeers()
		c.JSON(http.StatusOK, gin.H{
			"peers": peers,
		})
	})

	return r
}

func generatePssNodeName() string {
	shortID := strings.Split(uuid.New().String(), "-")[0]
	return fmt.Sprintf("pss_%s", shortID)
}

func getPeers() []models.Peer {
	var peers []models.Peer
	return peers
}
