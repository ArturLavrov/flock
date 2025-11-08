package main

import (
	gossipingAgent "flock/internal/agent"
	"flock/internal/models"
	peerSampingService "flock/internal/peer-sampling-service/http"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <node_name>, <bind_port>")
		os.Exit(1)
	}
	nodeName := os.Args[1]
	bindPort := os.Args[2]

	agent := gossipingAgent.Start(nodeName, bindPort, getAllPeerSamplingServices())
	h := peerSampingService.Start(agent)

	log.Println("Peer sampling service started :8081")
	log.Fatal(http.ListenAndServe(":8081", h))
}

func getAllPeerSamplingServices() []models.Peer {
	return []models.Peer{}
}
