package models

import (
	"net"
	"strconv"
)

type Node struct {
	ID            string `json:"id"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	ServiceName   string `json:"service_name"`
	LastHeartbeat int64  `json:"last_heartbeat"`
}

func (n Node) getAddr() string {
	return net.JoinHostPort(n.IP, strconv.Itoa(n.Port))
}
