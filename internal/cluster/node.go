package cluster

import (
	"net"
	"strconv"

	"github.com/google/uuid"
)

type Node struct {
	ID            string `json:"id"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	ServiceName   string `json:"service_name"`
	LastHeartbeat int64  `json:"last_heartbeat"`
	isAnchor      bool   `json:"is_anchor"`
}

func NewNode(serviceName string, port int, isAnchor bool) Node {
	return Node{
		ID:          uuid.NewString()[:8],
		ServiceName: serviceName,
		Port:        port,
		isAnchor:    isAnchor,
	}
}

func (n Node) getAddr() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(n.Port))
}

func (n Node) GetDisplayName() string {
	return n.ServiceName + "_" + n.ID
}

func (n Node) GetPort() int {
	return n.Port
}
