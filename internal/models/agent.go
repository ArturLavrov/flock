package models

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
)

// Agent defines the behavior of a Flock agent.
// It keeps track of known nodes and provides sampling operations.
type Agent interface {
	// GetAllNodes returns all known nodes (from memberlist)
	GetAllNodes() []Node
	GetAllNodesByService(service string) []Node
}

type flockAgent struct {
	mlist   *memberlist.Memberlist
	events  chan memberlist.NodeEvent
	mu      sync.RWMutex
	nodes   map[string]Node // key = member.Name
	started time.Time
}

// Start creates and starts an agent node and returns it.
func Start(nodeName string, bindPort string, peers []Peer) Agent {
	cfg := memberlist.DefaultLocalConfig()
	cfg.Name = nodeName
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

	if len(peers) > 0 {
		//TODO: improve this logic later
		// we join by address (host:port)
		n, err := list.Join([]string{peers[0].Address})
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
	for _, n := range a.nodes {
		result = append(result, n)
	}
	return result
}

// GetAllNodesByService returns nodes filtered by service name (optional helper).
func (a *flockAgent) GetAllNodesByService(service string) []Node {
	if service == "" {
		return a.GetAllNodes()
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]Node, 0)
	for _, n := range a.nodes {
		if n.Service == service {
			result = append(result, n)
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
		Addr:          fmt.Sprintf("%s:%d", m.Addr.To4().String(), m.Port),
		Service:       m.Name,
		LastHeartbeat: time.Now().Unix(),
	}
}
