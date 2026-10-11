package models

import (
	"net"
	"testing"

	"github.com/hashicorp/memberlist"
)

func TestMemberToNodeUsesMemberIDAndServiceName(t *testing.T) {
	member := &memberlist.Node{
		Name: "8d312a3f",
		Addr: net.ParseIP("127.0.0.1"),
		Port: 7946,
		Meta: []byte("pss"),
	}

	node := memberToNode(member)
	if node.ID != member.Name {
		t.Errorf("node ID = %q, want member ID %q", node.ID, member.Name)
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
