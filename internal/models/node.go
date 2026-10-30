package models

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
}

func NewNode(serviceName string, port int) Node {
	return Node{
		ID:          uuid.NewString()[:8],
		ServiceName: serviceName,
		Port:        port,
	}
}

func (n Node) getAddr() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(n.Port))
}

func (n Node) getDisplayName() string {
	return n.ServiceName + "_" + n.ID
}

func (n Node) getPort() int {
	return n.Port
}
