package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	agent "flock/internal/agent"
	api "flock/internal/api"
	cluster "flock/internal/cluster"
	"flock/internal/console"
)

func main() {
	if len(os.Args) < 4 || len(os.Args) > 5 {
		fmt.Println("Usage: go run main.go <service_name> <gossip_port> <api_port> [is_seed_node]")
		os.Exit(1)
	}

	flk_console := console.New()
	serviceName := getServiceName(os.Args[1])
	gossipPort := getGossipPort(os.Args[2])
	apiPort := getAPIPort(os.Args[3])
	isSeedNode := false
	if len(os.Args) == 5 {
		isSeedNode = getIsSeedNode(os.Args[4])
	}
	currentNode := cluster.NewNode(serviceName, gossipPort, isSeedNode)

	discoveredNodes := []cluster.Node{}
	var err error
	// TODO: Handle seed discovery when multiple seed nodes are configured.
	if !isSeedNode {
		discoveredNodes, err = cluster.GetAvailableNodes()
		if err != nil {
			log.Fatalf("Unable to discover cluster nodes: %v", err)
		}
	}
	gossipAgent := agent.New(&currentNode, discoveredNodes, apiPort)
	if err := gossipAgent.Start(); err != nil {
		log.Fatalf("Unable to start agent: %v", err)
	}

	apiAddress := fmt.Sprintf(":%d", apiPort)
	flk_console.Log(console.LogLevelInfo, fmt.Sprintf("Agent HTTP API running on %s", apiAddress))
	if err := http.ListenAndServe(apiAddress, api.HTTPHandler(gossipAgent)); err != nil {
		log.Fatalf("Agent HTTP API stopped: %v", err)
	}
}

func getGossipPort(gossipPortArg string) int {
	port, err := strconv.Atoi(gossipPortArg)
	if err != nil || port < 1 || port > 65535 {
		log.Fatalf("Invalid gossip port: %q (valid range: 1-65535)", gossipPortArg)
	}
	return port
}

func getAPIPort(apiPortArg string) int {
	port, err := strconv.Atoi(apiPortArg)
	if err != nil || port < 1 || port > 65535 {
		log.Fatalf("Invalid API port: %q (valid range: 1-65535)", apiPortArg)
	}
	return port
}

func getIsSeedNode(isSeedNodeArg string) bool {
	isSeedNode, err := strconv.ParseBool(isSeedNodeArg)
	if err != nil {
		log.Fatalf("Invalid seed node value: %q (must be true or false)", isSeedNodeArg)
	}
	return isSeedNode
}

func getServiceName(serviceNameArg string) string {
	if serviceNameArg == "" {
		log.Fatal("Service name cannot be empty")
	}
	return serviceNameArg
}
