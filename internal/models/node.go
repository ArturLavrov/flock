package models

type Node struct {
	ID            string `json:"id"`
	Addr          string `json:"addr"`
	Tag           string `json:"tag"`
	LastHeartbeat int64  `json:"last_heartbeat"`
}
