package agent

import (
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"

	"flock/internal/cluster"
	"flock/internal/console"

	"github.com/hashicorp/memberlist"
)

func TestNewStoresCurrentNode(t *testing.T) {
	current := cluster.NewNode("agent", 7946, true)
	nodes := []cluster.Node{cluster.NewNode("peer", 7947, false)}

	got := New(&current, nodes, 8081)

	if got.CurrentNode != current {
		t.Errorf("current node = %#v, want %#v", got.CurrentNode, current)
	}
}

func TestAgentProperties(t *testing.T) {
	current := cluster.NewNode("payments", 4200, true)
	a := New(&current, nil, 8081)
	anotherAgent := New(&current, nil, 8081)

	if len(a.AgentID) != 8 {
		t.Errorf("AgentID = %q, want an 8-character ID", a.AgentID)
	}
	if a.AgentID == current.ID {
		t.Errorf("AgentID = %q, should be distinct from the node ID", a.AgentID)
	}
	if a.AgentID == anotherAgent.AgentID {
		t.Errorf("two agents unexpectedly share ID %q", a.AgentID)
	}
	if a.AgentID == "" {
		t.Error("AgentID is empty")
	}
	if a.ServiceName != current.ServiceName {
		t.Errorf("ServiceName = %q, want %q", a.ServiceName, current.ServiceName)
	}
	if a.GossipPort != current.Port {
		t.Errorf("GossipPort = %d, want %d", a.GossipPort, current.Port)
	}
	if a.ApiPort != 8081 {
		t.Errorf("ApiPort = %d, want 8081", a.ApiPort)
	}
	if a.Protocol != "Gossip" || a.Status != "Starting" {
		t.Errorf("Protocol and status = %q/%q, want Gossip/Starting", a.Protocol, a.Status)
	}
}

func TestNewRejectsInvalidInput(t *testing.T) {
	if testCase := os.Getenv("FLOCK_AGENT_INVALID_INPUT"); testCase != "" {
		switch testCase {
		case "nil-current":
			New(nil, []cluster.Node{{}}, 8081)
		}
		return
	}

	tests := []struct {
		name        string
		wantMessage string
	}{
		{name: "nil-current", wantMessage: "Current node cannot be nil"},
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
	a := New(&current, nil, 8081)

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

type recordingConsole struct {
	messages []string
	levels   []console.LogLevel
}

func (output *recordingConsole) Print(message string) {
	output.messages = append(output.messages, message)
}

func (*recordingConsole) PrintFlockTitle() {}

func (*recordingConsole) PrintAgentInfo(string, string, int, int, string, string, string) {}

func (output *recordingConsole) Log(level console.LogLevel, text string) {
	output.messages = append(output.messages, text)
	output.levels = append(output.levels, level)
}

func TestMembershipEventHandlerTracksAndReportsJoinAndLeave(t *testing.T) {
	agent := &GossipAgent{}
	output := &recordingConsole{}
	handler := membershipEventHandler{agent: agent, console: output}
	member := &memberlist.Node{
		Name: "worker_node1234",
		Addr: net.ParseIP("127.0.0.1"),
		Port: 7946,
		Meta: []byte("worker"),
	}

	handler.NotifyJoin(member)
	nodes := agent.GetAllNodes()
	if len(nodes) != 1 || nodes[0].ID != "node1234" || nodes[0].ServiceName != "worker" {
		t.Fatalf("nodes after join = %#v", nodes)
	}
	if len(output.messages) != 1 || output.messages[0] != "Node joined: worker_node1234" {
		t.Errorf("join messages = %#v", output.messages)
	}
	if len(output.levels) != 1 || output.levels[0] != console.LogLevelInfo {
		t.Errorf("join log levels = %#v, want Info", output.levels)
	}

	handler.NotifyLeave(member)
	if nodes := agent.GetAllNodes(); len(nodes) != 0 {
		t.Errorf("nodes after leave = %#v, want none", nodes)
	}
	if len(output.messages) != 2 || output.messages[1] != "Node leave: worker_node1234" {
		t.Errorf("leave messages = %#v", output.messages)
	}
	if len(output.levels) != 2 || output.levels[1] != console.LogLevelWarn {
		t.Errorf("leave log levels = %#v, want Warn", output.levels)
	}
}
