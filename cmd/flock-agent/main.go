package main

import (
	"encoding/json"
	gossipingAgent "flock/internal/flock-agent"
	"flock/internal/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <node_name> <bind_port> <agent_api_port>")
		os.Exit(1)
	}

	nodeName := os.Args[1]
	bindPort := os.Args[2]
	apiPort := fmt.Sprintf(":%s", os.Args[3])

	agent := models.Start(nodeName, bindPort, getAvailablePeers())
	h := gossipingAgent.HTTPHandler(agent)

	log.Printf("Agent REST running on %s\n", apiPort)
	log.Fatal(http.ListenAndServe(apiPort, h))
}

func getAvailablePeers() []models.Peer {
	peerLimit := 5

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:8081/v1/peers?limit=%d", peerLimit)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Failed to call Peer Sampling Service: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected response status: %d\n", resp.StatusCode)
		return nil
	}

	var result struct {
		Peers []models.Peer `json:"peers"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to decode peers: %v\n", err)
		return nil
	}

	return result.Peers
}
