package main

import (
	gossipingAgent "flock/internal/agent"
	"flock/internal/models"
	peerSampingService "flock/internal/peer-sampling-service/http"
	"flock/internal/shortID"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 1 {
		//TODO: think about proper cli design.
		//http addres and gossip port should be specified.
		fmt.Println("Usage: go run main.go [pss_node_host: pss_node_port]")
		os.Exit(1)
	}

	var pss_node []*models.Peer
	pss_nodes := parsePssNodesList(os.Args[1:])
	root_node := pss_node[0];
	
	agent := gossipingAgent.Start(root_node.ID, root_node.Port, pss_nodes)
	h := peerSampingService.Start(agent)

	log.Println("Peer sampling service started :8081")
	log.Fatal(http.ListenAndServe(":8081", h))
}

func parsePssNodesList(addresses []string) []*models.Peer {
	var pss_node_ID = shortID.Generate();
	var generatePssNodeName = func () string {return fmt.Sprintf("pss_%s", pss_node_ID)}
	var pss_list []*models.Peer
	for _, addr := range addresses {
		host, port, err := net.SplitHostPort(addr)
			if err != nil {
				log.Fatalf("invalid address %q: %v\n", addr, err)
			}
			validateGossipPort(port)
			pss_node := &models.Peer{
				//TODO: clarify a difference between ID and ServiceName
				ID:pss_node_ID,
				Address:addr,
				ServiceName:generatePssNodeName(),
				Host: host,
				Port: port,
			}
			pss_list = append(pss_list, pss_node)
	}

	if len(pss_list) == 1 {
		return []*models.Peer{}
	}
		
	return pss_list
}

func validateGossipPort(portStr string) string {
	const(
		MIN_PORT_VALUE = 1
		MAX_PORT_VALUE = 2
	)
	
	if portStr == "" {
		log.Fatal("Port cannot be empty")
	}

	portValue, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %s", portStr)
	}

	if portValue < MIN_PORT_VALUE || portValue > MAX_PORT_VALUE {
		log.Fatalf("Port out of range: %d (valid range: 1–65535)", portValue)
	}

	return portStr
}



