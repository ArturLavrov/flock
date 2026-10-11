package main

import (
	peerSampingService "flock/internal/peer-sampling-service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <gossip_port>")
		os.Exit(1)
	}

	gossipPort := validateGossipPort(os.Args[1])
	h := peerSampingService.Start(gossipPort)

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


