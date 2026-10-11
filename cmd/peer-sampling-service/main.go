package main

import (
	gossipingAgent "flock/internal/flock-agent"
	"flock/internal/models"
	peerSampingService "flock/internal/peer-sampling-service/http"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/google/uuid"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <gossip_port>")
		os.Exit(1)
	}

	gossipPort := validateGossipPort(os.Args[1])
	pssNodeName := generatePssNodeName()
	peers := getPeers()

	agent := gossipingAgent.Start(pssNodeName, gossipPort, peers)
	h := peerSampingService.Start(agent)

	log.Println("Peer sampling service started :8081")
	log.Fatal(http.ListenAndServe(":8081", h))
}

func validateGossipPort(portStr string) string {
	const (
		MIN_PORT_VALUE = 1
		MAX_PORT_VALUE = 65535
	)

	if portStr == "" {
		log.Fatal("Port cannot be empty")
	}

	portValue, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s", portStr)
	}

	if portValue < MIN_PORT_VALUE || portValue > MAX_PORT_VALUE {
		log.Fatalf("Port out of range: %d (valid range: 1-65535)", portValue)
	}

	return portStr
}

func generatePssNodeName() string {
	var shortId = strings.Split(uuid.New().String(), "-")[0]
	return fmt.Sprintf("pss_%s", shortId)
}

func getPeers() []models.Peer {
	var peers []models.Peer
	return peers
}


