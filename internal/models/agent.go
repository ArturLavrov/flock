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
	GetAllPeers() []Peer
	GetAllPeersByName(name string) []Peer
}

type flockAgent struct {
	mlist   *memberlist.Memberlist
	events  chan memberlist.NodeEvent
	mu      sync.RWMutex
	peers   map[string]Peer // key = member.Name
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
		peers:   make(map[string]Peer),
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
		a.peers[m.Name] = memberToPeer(m)
	}
}

func (a *flockAgent) addOrUpdateNode(m *memberlist.Node) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.peers[m.Name] = memberToPeer(m)
}

func (a *flockAgent) removeNode(m *memberlist.Node) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.peers, m.Name)
}

// GetAllPeers returns all known peers (copy).
func (a *flockAgent) GetAllPeers() []Peer {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]Peer, 0, len(a.peers))
	for _, peer := range a.peers {
		result = append(result, peer)
	}
	return result
}

// GetAllPeersByName returns peers filtered by name.
func (a *flockAgent) GetAllPeersByName(name string) []Peer {
	if name == "" {
		return a.GetAllPeers()
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make([]Peer, 0)
	for _, peer := range a.peers {
		if peer.ServiceName == name {
			result = append(result, peer)
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

// memberToPeer converts a memberlist node to our internal Peer.
func memberToPeer(m *memberlist.Node) Peer {
	return Peer{
		ID:            m.Name,
		Address:       fmt.Sprintf("%s:%d", m.Addr.To4().String(), m.Port),
		ServiceName:   m.Name,
		LastHeartbeat: time.Now().Unix(),
	}
}
