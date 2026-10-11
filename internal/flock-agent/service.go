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
		peers := a.GetAllPeers()
		nodes := make([]models.Node, 0, len(peers))
		for _, peer := range peers {
			nodes = append(nodes, models.Node{
				ID:            peer.ID,
				Addr:          peer.Address,
				Service:       peer.ServiceName,
				LastHeartbeat: peer.LastHeartbeat,
			})
		}
		c.JSON(http.StatusOK, nodes)
	})

	return r
}
