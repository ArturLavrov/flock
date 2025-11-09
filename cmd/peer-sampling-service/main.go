package main

import (
	gossipingAgent "flock/internal/agent"
	"flock/internal/models"
	peerSampingService "flock/internal/peer-sampling-service/http"
	"flock/internal/shortID"
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
	nodeName := generatePSSNodeName()
	bindPort := os.Args[2]

	agent := gossipingAgent.Start(nodeName, bindPort, getAllPeerSamplingServices())
	h := peerSampingService.Start(agent)

	log.Println("Peer sampling service started :8081")
	log.Fatal(http.ListenAndServe(":8081", h))
}

func generatePSSNodeName() string {
	return fmt.Sprintf("pss_%s", shortID.Generate())
}

func getAllPeerSamplingServices() []models.Peer {
	//TODO:(https://github.com/ArturLavrov/flock/issues/2)
	return []models.Peer{}
}
