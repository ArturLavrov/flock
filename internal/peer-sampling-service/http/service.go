package http

import (
	gossipAgent "flock/internal/flock-agent"
	"flock/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start(agent gossipAgent.Agent) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/v1/peers", func(c *gin.Context) {
		nodes := agent.GetAllNodes()

		var peers []models.Peer
		for _, n := range nodes {
			peers = append(peers, models.Peer{
				ID:          n.ID,
				Address:     n.Addr,
				ServiceName: n.Service,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"peers": peers,
		})
	})

	return r
}
