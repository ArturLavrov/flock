package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"flock/internal/agent"
)

func TestNodesEndpointReturnsWrappedList(t *testing.T) {
	current := agent.GossipAgent{}
	request := httptest.NewRequest(http.MethodGet, "/api/v1/nodes", nil)
	response := httptest.NewRecorder()

	HTTPHandler(&current).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var result struct {
		Nodes []json.RawMessage `json:"nodes"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Nodes == nil {
		t.Fatalf("response nodes = nil, want an empty array: %s", response.Body.String())
	}
	if len(result.Nodes) != 0 {
		t.Fatalf("response nodes = %s, want empty array", response.Body.String())
	}
}
