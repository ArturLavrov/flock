package agent

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"flock/internal/cluster"
)

func TestNewStoresCurrentNode(t *testing.T) {
	current := cluster.NewNode("agent", 7946, true)
	nodes := []cluster.Node{cluster.NewNode("peer", 7947, false)}

	got := New(&current, nodes)

	if got.CurrentNode != current {
		t.Errorf("current node = %#v, want %#v", got.CurrentNode, current)
	}
}

func TestNewRejectsInvalidInput(t *testing.T) {
	if testCase := os.Getenv("FLOCK_AGENT_INVALID_INPUT"); testCase != "" {
		switch testCase {
		case "nil-current":
			New(nil, []cluster.Node{{}})
		case "empty-nodes":
			New(&cluster.Node{}, nil)
		}
		return
	}

	tests := []struct {
		name        string
		wantMessage string
	}{
		{name: "nil-current", wantMessage: "Current node cannot be nil"},
		{name: "empty-nodes", wantMessage: "Cluster nodes cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestNewRejectsInvalidInput$")
			cmd.Env = append(os.Environ(), "FLOCK_AGENT_INVALID_INPUT="+tt.name)
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("New unexpectedly returned without terminating")
			}
			if !strings.Contains(string(output), tt.wantMessage) {
				t.Errorf("output %q does not contain %q", output, tt.wantMessage)
			}
		})
	}
}

func TestStartAndStopReturnNil(t *testing.T) {
	current := cluster.NewNode("agent", 7946, false)
	a := New(&current, []cluster.Node{cluster.NewNode("peer", 7947, false)})

	if err := a.Start(); err != nil {
		t.Errorf("Start() error = %v, want nil", err)
	}
	if len(a.nodes) != 1 || a.nodes[0].ID != current.ID || a.nodes[0].ServiceName != current.ServiceName {
		t.Errorf("nodes stored after Start() = %#v, want current node %#v", a.nodes, current)
	}
	if err := a.Stop(); err != nil {
		t.Errorf("Stop() error = %v, want nil", err)
	}
}

func TestGetNodesByName(t *testing.T) {
	agent := &GossipAgent{
		nodes: []cluster.Node{
			{ID: "node-1", ServiceName: "worker"},
			{ID: "node-2", ServiceName: "worker"},
			{ID: "node-3", ServiceName: "database"},
		},
	}

	tests := []struct {
		testName  string
		queryName string
		want      []string
	}{
		{testName: "matching service", queryName: "worker", want: []string{"node-1", "node-2"}},
		{testName: "empty name returns all", queryName: "", want: []string{"node-1", "node-2", "node-3"}},
		{testName: "missing service", queryName: "missing", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := agent.GetNodesByName(tt.queryName)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d nodes, want %d: %#v", len(got), len(tt.want), got)
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("node %d ID = %q, want %q", i, got[i].ID, id)
				}
			}
		})
	}
}
