package cluster

import (
	"strconv"
	"testing"
)

func TestNewNodeCreatesUUIDFormatID(t *testing.T) {
	node := NewNode("pss", 7946, false)

	if len(node.ID) != 8 {
		t.Fatalf("node ID length = %d, want 8", len(node.ID))
	}
	if _, err := strconv.ParseUint(node.ID, 16, 32); err != nil {
		t.Errorf("node ID %q is not an 8-character hexadecimal ID: %v", node.ID, err)
	}
	if node.ServiceName != "pss" || node.Port != 7946 {
		t.Errorf("node configuration = (%q, %d), want (%q, %d)", node.ServiceName, node.Port, "pss", 7946)
	}
	if node.isAnchor {
		t.Error("node isAnchor = true, want false")
	}
}

func TestNewNodeSetsAnchorFlag(t *testing.T) {
	node := NewNode("pss", 7946, true)
	if !node.isAnchor {
		t.Fatal("node isAnchor = false, want true")
	}
}

func TestNodeGetDisplayName(t *testing.T) {
	node := Node{ServiceName: "pss", ID: "8d312a3f"}

	if got, want := node.GetDisplayName(), "pss_8d312a3f"; got != want {
		t.Errorf("display name = %q, want %q", got, want)
	}
}

func TestNodeGetAddrAndPort(t *testing.T) {
	node := Node{IP: "127.0.0.1", Port: 7946}

	if got, want := node.getAddr(), "127.0.0.1:7946"; got != want {
		t.Errorf("getAddr() = %q, want %q", got, want)
	}
	if got, want := node.GetPort(), 7946; got != want {
		t.Errorf("getPort() = %d, want %d", got, want)
	}
}
