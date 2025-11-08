package models

type Peer struct {
	ID          string `json:"id"`
	Address     string `json:"address"`
	ServiceName string `json:"service_name"`
}
