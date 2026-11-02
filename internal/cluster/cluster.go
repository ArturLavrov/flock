package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetAvailableNodes() ([]Node, error) {
	const nodeLimit = 5

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url := fmt.Sprintf("http://localhost:8081/api/v1/nodes?limit=%d", nodeLimit)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call peer sampling service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected peer sampling service response status: %d", resp.StatusCode)
	}

	var result struct {
		Nodes []Node `json:"nodes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode peer sampling service nodes: %w", err)
	}

	return result.Nodes, nil
}
