package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/memberlist"
)

// Agent defines the behavior of a Flock agent.
// It keeps track of known nodes and provides sampling operations.
type Agent interface {
	GetAllNodes() []Node
	GetNodesByName(name string) []Node
}

type flockAgent struct {
	mlist   *memberlist.Memberlist
	events  chan memberlist.NodeEvent
	mu      sync.RWMutex
	nodes   map[string]Node // key = member.Name
	started time.Time
}

// Start creates and starts an agent node and returns it.
func Start(serviceName string, bindPort string, nodes []Node) Agent {
	cfg := memberlist.DefaultLocalConfig()
	cfg.Name = uuid.NewString()[:8]
	cfg.Delegate = nodeMetadataDelegate{serviceName: serviceName}
	cfg.BindPort = parsePort(bindPort)

	evCh := make(chan memberlist.NodeEvent, 16)
	cfg.Events = &memberlist.ChannelEventDelegate{Ch: evCh}

	list, err := memberlist.Create(cfg)
	if err != nil {
		log.Fatalf("Failed to create memberlist: %v", err)
	}

	a := &flockAgent{
		mlist:   list,
		events:  evCh,
		nodes:   make(map[string]Node),
		started: time.Now(),
	}

	if len(nodes) > 0 {
		n, err := list.Join([]string{nodes[0].getAddr()})
		if err != nil {
			log.Printf("Failed to join cluster: %v", err)
		} else {
			log.Printf("Successfully joined %d nodes!", n)
		}
	} else {
		log.Println("Started as a standalone node.")
	}

	a.refreshFromMemberlist()

	go a.handleEvents()

	return a
}

type nodeMetadataDelegate struct {
	serviceName string
}

func (d nodeMetadataDelegate) NodeMeta(limit int) []byte {
	meta := []byte(d.serviceName)
	if len(meta) > limit {
		return nil
	}
	return meta
}

func (nodeMetadataDelegate) NotifyMsg([]byte)                {}
func (nodeMetadataDelegate) GetBroadcasts(int, int) [][]byte { return nil }
func (nodeMetadataDelegate) LocalState(bool) []byte          { return nil }
func (nodeMetadataDelegate) MergeRemoteState([]byte, bool)   {}

// handleEvents consumes memberlist events and updates our in-memory nodes map.
func (a *flockAgent) handleEvents() {
	for ev := range a.events {
		switch ev.Event {
		case memberlist.NodeJoin:
			a.addOrUpdateNode(ev.Node)
			log.Printf("[agent] node joined: %s (%s)", ev.Node.Name, ev.Node.Address())
		case memberlist.NodeLeave:
			a.removeNode(ev.Node)
			log.Printf("[agent] node left: %s", ev.Node.Name)
		case memberlist.NodeUpdate:
			a.addOrUpdateNode(ev.Node)
			log.Printf("[agent] node updated: %s", ev.Node.Name)
		}
	}
}

// refreshFromMemberlist pulls current members from memberlist and syncs to memory.
func (a *flockAgent) refreshFromMemberlist() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, m := range a.mlist.Members() {
		a.nodes[m.Name] = memberToNode(m)
	}
}

func (a *flockAgent) addOrUpdateNode(m *memberlist.Node) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.nodes[m.Name] = memberToNode(m)
}

func (a *flockAgent) removeNode(m *memberlist.Node) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.nodes, m.Name)
}

// GetAllNodes returns all known nodes (copy).
func (a *flockAgent) GetAllNodes() []Node {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]Node, 0, len(a.nodes))
	for _, node := range a.nodes {
		result = append(result, node)
	}
	return result
}

// GetNodesByName returns nodes filtered by name.
func (a *flockAgent) GetNodesByName(name string) []Node {
	if name == "" {
		return a.GetAllNodes()
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]Node, 0)
	for _, node := range a.nodes {
		if node.ServiceName == name {
			result = append(result, node)
		}
	}
	return result
}

func parsePort(s string) int {
	var p int
	_, err := fmt.Sscanf(s, "%d", &p)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	return p
}

// memberToNode converts a memberlist node to our internal Node.
func memberToNode(m *memberlist.Node) Node {
	return Node{
		ID:            m.Name,
		IP:            m.Addr.String(),
		Port:          int(m.Port),
		ServiceName:   string(m.Meta),
		LastHeartbeat: time.Now().Unix(),
	}
}
