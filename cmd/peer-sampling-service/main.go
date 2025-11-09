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
	"strconv"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("Usage: go run main.go <gossip_port>")
		os.Exit(1)
	}

	bindPort := validateGossipPort(os.Args[1])
	nodeName := generatePSSNodeName()

	agent := gossipingAgent.Start(nodeName, bindPort, getAllPeerSamplingServices())
	h := peerSampingService.Start(agent)

	log.Println("Peer sampling service started :8081")
	log.Fatal(http.ListenAndServe(":8081", h))
}

func validateGossipPort(portStr string) int {
	if portStr == "" {
		log.Fatal("Port cannot be empty")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s", portStr)
	}

	if port < 1 || port > 65535 {
		log.Fatalf("Port out of range: %d (valid range: 1–65535)", port)
	}

	return port
}

func generatePSSNodeName() string {
	return fmt.Sprintf("pss_%s", shortID.Generate())
}

func getAllPeerSamplingServices() []models.Peer {
	//TODO:(https://github.com/ArturLavrov/flock/issues/2)
	return []models.Peer{}
}
