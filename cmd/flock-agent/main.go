package main

import (
	"encoding/json"
	service "flock/internal/flock-agent"
	gossipAgent "flock/internal/models"
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

	agent := gossipAgent.Start(nodeName, bindPort, getAvailableNodes())
	h := service.HTTPHandler(agent)

	log.Printf("Agent REST running on %s\n", apiPort)
	log.Fatal(http.ListenAndServe(apiPort, h))
}

func getAvailableNodes() []gossipAgent.Node {
	nodeLimit := 5

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:8081/v1/nodes?limit=%d", nodeLimit)

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
		Nodes []gossipAgent.Node `json:"nodes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("Failed to decode nodes: %v\n", err)
		return nil
	}

	return result.Nodes
}
