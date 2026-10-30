package models

import (
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/memberlist"
)

// Agent defines the behavior of a Flock agent.
// It keeps track of known nodes and provides sampling operations.
type Agent interface {
	GetAllNodes() []Node
	GetNodesByName(name string) []Node
	GetLocalNode() Node
}

type MembershipEvent struct {
	Type      string
	NodeID    string
	Address   string
	Details   string
	Timestamp time.Time
}

const (
	MembershipJoin  = "JOIN"
	MembershipLeave = "LEAVE"
)

type flockAgent struct {
	mlist             *memberlist.Memberlist
	events            chan memberlist.NodeEvent
	mu                sync.RWMutex
	nodes             map[string]Node // key = member.Name
	started           time.Time
	onMembershipEvent func(MembershipEvent)
}

// Start creates and starts an agent node and returns it.
func Start(node Node, nodes []Node) Agent {
	agent, err := start(node, nodes, nil)
	if err != nil {
		log.Printf("Failed to join cluster: %v", err)
	}
	return agent
}

// StartWithEvents starts an agent and reports membership transitions.
func StartWithEvents(node Node, nodes []Node, onMembershipEvent func(MembershipEvent)) (Agent, error) {
	return start(node, nodes, onMembershipEvent)
}

func start(node Node, nodes []Node, onMembershipEvent func(MembershipEvent)) (Agent, error) {
	cfg := memberlist.DefaultLocalConfig()
	cfg.Name = node.getDisplayName()
	cfg.Delegate = nodeMetadataDelegate{serviceName: node.ServiceName}
	cfg.BindPort = node.getPort()
	if onMembershipEvent != nil {
		cfg.Logger = log.New(io.Discard, "", 0)
	}

	evCh := make(chan memberlist.NodeEvent, 16)
	cfg.Events = &memberlist.ChannelEventDelegate{Ch: evCh}

	list, err := memberlist.Create(cfg)
	if err != nil {
		log.Fatalf("Failed to create memberlist: %v", err)
	}

	a := &flockAgent{
		mlist:             list,
		events:            evCh,
		nodes:             make(map[string]Node),
		started:           time.Now(),
		onMembershipEvent: onMembershipEvent,
	}

	a.refreshFromMemberlist()
	go a.handleEvents()

	var startupErr error
	if len(nodes) > 0 {
		n, err := list.Join([]string{nodes[0].getAddr()})
		if err != nil {
			if onMembershipEvent == nil {
				log.Printf("Failed to join cluster: %v", err)
			} else {
				startupErr = err
			}
		} else if onMembershipEvent == nil {
			log.Printf("Successfully joined %d nodes!", n)
		}
	} else if onMembershipEvent == nil {
		log.Println("Started as a standalone node.")
	}

	return a, startupErr
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
		a.handleEvent(ev)
	}
}

func (a *flockAgent) handleEvent(ev memberlist.NodeEvent) {
	node := memberToNode(ev.Node)

	switch ev.Event {
	case memberlist.NodeJoin:
		isNew := a.addOrUpdateNode(ev.Node)
		if a.onMembershipEvent == nil {
			log.Printf("[agent] node joined: %s (%s)", ev.Node.Name, ev.Node.Address())
		}
		if isNew {
			a.publishMembershipEvent(MembershipJoin, node)
		}
	case memberlist.NodeLeave:
		wasKnown := a.removeNode(ev.Node)
		if a.onMembershipEvent == nil {
			log.Printf("[agent] node left: %s", ev.Node.Name)
		}
		if wasKnown {
			a.publishMembershipEvent(MembershipLeave, node)
		}
	case memberlist.NodeUpdate:
		a.addOrUpdateNode(ev.Node)
		if a.onMembershipEvent == nil {
			log.Printf("[agent] node updated: %s", ev.Node.Name)
		}
	}
}

func (a *flockAgent) publishMembershipEvent(eventType string, node Node) {
	if a.onMembershipEvent == nil {
		return
	}
	a.onMembershipEvent(MembershipEvent{
		Type:      eventType,
		NodeID:    node.ID,
		Address:   node.getAddr(),
		Timestamp: time.Now(),
	})
}

// refreshFromMemberlist pulls current members from memberlist and syncs to memory.
func (a *flockAgent) refreshFromMemberlist() {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, m := range a.mlist.Members() {
		a.nodes[m.Name] = memberToNode(m)
	}
}

func (a *flockAgent) addOrUpdateNode(m *memberlist.Node) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, exists := a.nodes[m.Name]
	a.nodes[m.Name] = memberToNode(m)
	return !exists
}

func (a *flockAgent) removeNode(m *memberlist.Node) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, exists := a.nodes[m.Name]
	delete(a.nodes, m.Name)
	return exists
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

func (a *flockAgent) GetLocalNode() Node {
	return memberToNode(a.mlist.LocalNode())
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

// memberToNode converts a memberlist node to our internal Node.
func memberToNode(m *memberlist.Node) Node {
	id := strings.TrimPrefix(m.Name, string(m.Meta)+"_")
	return Node{
		ID:            id,
		IP:            m.Addr.String(),
		Port:          int(m.Port),
		ServiceName:   string(m.Meta),
		LastHeartbeat: time.Now().Unix(),
	}
}
