package models

import (
	"net"
	"testing"

	"github.com/hashicorp/memberlist"
)

func TestMemberToNodeUsesMemberIDAndServiceName(t *testing.T) {
	member := &memberlist.Node{
		Name: "pss_8d312a3f",
		Addr: net.ParseIP("127.0.0.1"),
		Port: 7946,
		Meta: []byte("pss"),
	}

	node := memberToNode(member)
	if node.ID != "8d312a3f" {
		t.Errorf("node ID = %q, want %q", node.ID, "8d312a3f")
	}
	if node.ServiceName != "pss" {
		t.Errorf("service name = %q, want %q", node.ServiceName, "pss")
	}
}

func TestNodeMetadataDelegateProvidesServiceName(t *testing.T) {
	delegate := nodeMetadataDelegate{serviceName: "pss"}

	if got := string(delegate.NodeMeta(512)); got != "pss" {
		t.Errorf("node metadata = %q, want %q", got, "pss")
	}
	if got := delegate.NodeMeta(2); got != nil {
		t.Errorf("oversized node metadata = %q, want nil", got)
	}
}

func TestMembershipEventsOnlyReportStateTransitions(t *testing.T) {
	var events []MembershipEvent
	agent := &flockAgent{
		nodes: make(map[string]Node),
		onMembershipEvent: func(event MembershipEvent) {
			events = append(events, event)
		},
	}
	member := &memberlist.Node{
		Name: "pss_node-b19c",
		Addr: net.ParseIP("127.0.0.1"),
		Port: 7002,
		Meta: []byte("pss"),
	}

	agent.handleEvent(memberlist.NodeEvent{Event: memberlist.NodeJoin, Node: member})
	agent.handleEvent(memberlist.NodeEvent{Event: memberlist.NodeJoin, Node: member})
	agent.handleEvent(memberlist.NodeEvent{Event: memberlist.NodeUpdate, Node: member})
	agent.handleEvent(memberlist.NodeEvent{Event: memberlist.NodeLeave, Node: member})
	agent.handleEvent(memberlist.NodeEvent{Event: memberlist.NodeLeave, Node: member})

	if len(events) != 2 {
		t.Fatalf("got %d membership events, want 2: %#v", len(events), events)
	}
	if events[0].Type != MembershipJoin || events[0].NodeID != "node-b19c" {
		t.Errorf("join event = %#v", events[0])
	}
	if events[1].Type != MembershipLeave || events[1].NodeID != "node-b19c" {
		t.Errorf("leave event = %#v", events[1])
	}
}
