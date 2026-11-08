package cluster

import (
	"fmt"
	"net"
	"strconv"

	"github.com/google/uuid"
)

const (
	SeedNodeName = "flock-seed-node"
)

type Node struct {
	ID            string `json:"id"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	ServiceName   string `json:"service_name"`
	LastHeartbeat int64  `json:"last_heartbeat"`
	IsSeedNode    bool   `json:"is_seed_node"`
}

func NewNode(serviceName string, port int) Node {
	if serviceName == SeedNodeName {
		fmt.Errorf("% is a reserved name only for seed not and cannot be used!", SeedNodeName)
	}
	return Node{
		ID:          uuid.NewString()[:8],
		ServiceName: serviceName,
		Port:        port,
		IsSeedNode:  false,
	}
}

func NewSeedNode(gossitPort int) Node {
	return Node{
		ID:          uuid.NewString()[:8],
		ServiceName: SeedNodeName,
		Port:        gossitPort,
		IsSeedNode:  true,
	}
}

func (n *Node) getAddr() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(n.Port))
}

func (n *Node) GetDisplayName() string {
	return n.ServiceName + "_" + n.ID
}

func (n *Node) GetPort() int {
	return n.Port
}

func (n *Node) IsSeed() bool {
	return n.ServiceName == SeedNodeName
}
