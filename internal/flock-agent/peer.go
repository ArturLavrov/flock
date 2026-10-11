package agent

type Peer struct {
	ID          string `json:"id"`
	Address     string `json:"address"`
	ServiceName string `json:"service_name"`
	Host        string
	Port        string
}
