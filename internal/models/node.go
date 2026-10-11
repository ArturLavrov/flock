package models

type Node struct {
	ID            string
	Addr          string
	Service       string
	LastHeartbeat int64
}
