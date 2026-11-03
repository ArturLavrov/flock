package agent

import (
	"net"
	"testing"
	"time"

	"flock/internal/cluster"
)

func TestStartJoinsBootstrapNodesAndPropagatesJoin(t *testing.T) {
	firstNode := cluster.NewNode("first", availablePort(t), true)
	first := New(&firstNode, nil)
	if err := first.Start(); err != nil {
		t.Fatalf("start first agent: %v", err)
	}
	t.Cleanup(func() {
		if err := first.Stop(); err != nil {
			t.Errorf("stop first agent: %v", err)
		}
	})

	bootstrapNodes := first.GetAllNodes()
	if len(bootstrapNodes) != 1 {
		t.Fatalf("first agent nodes = %#v, want its local node", bootstrapNodes)
	}

	secondNode := cluster.NewNode("second", availablePort(t), false)
	second := New(&secondNode, bootstrapNodes)
	if err := second.Start(); err != nil {
		t.Fatalf("start second agent: %v", err)
	}
	t.Cleanup(func() {
		if err := second.Stop(); err != nil {
			t.Errorf("stop second agent: %v", err)
		}
	})

	waitForNode(t, first, secondNode.ID)
	waitForNode(t, second, firstNode.ID)
}

func availablePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate port: %v", err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}

func waitForNode(t *testing.T, agent *GossipAgent, nodeID string) {
	t.Helper()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		for _, node := range agent.GetAllNodes() {
			if node.ID == nodeID {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("node %q was not propagated; current membership: %#v", nodeID, agent.GetAllNodes())
}
