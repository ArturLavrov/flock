package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	agent "flock/internal/agent"
	api "flock/internal/api"
	cluster "flock/internal/cluster"
	"flock/internal/console"
)

func main() {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: go run main.go <service_name> <gossip_port> <api_port> [is_anchor_node]")
		os.Exit(1)
	}

	serviceName := getServiceName(os.Args[1])
	gossipPort := getGossipPort(os.Args[2])
	apiPort := getApiPort(os.Args[3])
	discoveredNodes := []cluster.Node{}
	var err error
	isAnchorNode := false
	if len(os.Args) == 5 {
		isAnchorNode = getIsAnchorNode(os.Args[4])
	}
	currentNode := cluster.NewNode(serviceName, gossipPort, isAnchorNode)

	if !isAnchorNode {
		discoveredNodes, err = cluster.GetAvailableNodes()
		if err != nil {
			log.Fatalf("Unable to discover cluster nodes: %v", err)
		}
	}

	agent := agent.New(&currentNode, discoveredNodes)
	if err := agent.Start(); err != nil {
		log.Fatalf("Unable to start agent: %v", err)
	}

	api.HTTPHandler(agent)
	console.New().Print(fmt.Sprintf("Agent HTTP API running on %s", apiPort))
}

func getGossipPort(gossipPortArg string) int {
	port, err := strconv.Atoi(gossipPortArg)
	if err != nil || port < 1 || port > 65535 {
		log.Fatalf("Invalid gossip port: %q (valid range: 1-65535)", gossipPortArg)
	}
	return port
}

func getApiPort(apiPortArg string) string {
	port, err := strconv.Atoi(apiPortArg)
	if err != nil || port < 1 || port > 65535 {
		log.Fatalf("Invalid API port: %q (valid range: 1-65535)", apiPortArg)
	}
	return ":" + strconv.Itoa(port)
}

func getIsAnchorNode(isAnchorNodeArg string) bool {
	isAnchorNode, err := strconv.ParseBool(isAnchorNodeArg)
	if err != nil {
		log.Fatalf("Invalid anchor node value: %q (must be true or false)", isAnchorNodeArg)
	}
	return isAnchorNode
}

func getServiceName(serviceNameArg string) string {
	if serviceNameArg == "" {
		log.Fatal("Service name cannot be empty")
	}
	return serviceNameArg
}
