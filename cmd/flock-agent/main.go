package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	service "flock/internal/flock-agent"
	gossipAgent "flock/internal/models"
)

const agentVersion = "v0.1.0"

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <service_name> <bind_port> <agent_api_port>")
		os.Exit(1)
	}

	serviceName := os.Args[1]
	bindPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatalf("Invalid bind port: %v", err)
	}
	apiPort := fmt.Sprintf(":%s", os.Args[3])

	agentNode := gossipAgent.NewNode(serviceName, bindPort)
	nodes, discoveryErr := getAvailableNodes()
	eventCh := make(chan gossipAgent.MembershipEvent)
	startupReady := make(chan struct{})
	go func() {
		pending := make([]gossipAgent.MembershipEvent, 0)
		printEvent := func(event gossipAgent.MembershipEvent) {
			if err := gossipAgent.PrintMembershipEvent(event); err != nil {
				log.Printf("Failed to print membership event: %v", err)
			}
		}

		for {
			select {
			case event := <-eventCh:
				pending = append(pending, event)
			case <-startupReady:
				for _, event := range pending {
					printEvent(event)
				}
				for event := range eventCh {
					printEvent(event)
				}
				return
			}
		}
	}()

	agent, joinErr := gossipAgent.StartWithEvents(agentNode, nodes, func(event gossipAgent.MembershipEvent) {
		eventCh <- event
	})
	localNode := agent.GetLocalNode()
	if err := gossipAgent.PrintAgentInfo(gossipAgent.AgentInfo{
		Version:  agentVersion,
		NodeID:   localNode.ID,
		Address:  net.JoinHostPort(localNode.IP, strconv.Itoa(localNode.Port)),
		Protocol: "SWIM / Gossip",
		Status:   "ALIVE",
	}); err != nil {
		log.Printf("Failed to print agent startup information: %v", err)
	}
	close(startupReady)
	if joinErr != nil {
		log.Printf("Failed to join cluster: %v", joinErr)
	}
	if discoveryErr != nil {
		fmt.Printf("%v\n", discoveryErr)
	}

	h := service.HTTPHandler(agent)

	log.Printf("Agent REST running on %s\n", apiPort)
	log.Fatal(http.ListenAndServe(apiPort, h))
}

func getAvailableNodes() ([]gossipAgent.Node, error) {
	nodeLimit := 5

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:8081/api/v1/nodes?limit=%d", nodeLimit)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to call Peer Sampling Service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response status: %d", resp.StatusCode)
	}

	var result struct {
		Nodes []gossipAgent.Node `json:"nodes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode nodes: %v", err)
	}

	return result.Nodes, nil
}
