package agent

import (
	"log"
	"strings"
	"sync"
	"time"

	"flock/internal/cluster"
	"flock/internal/console"

	"github.com/hashicorp/memberlist"
)

type GossipAgent struct {
	Current cluster.Node
	members *memberlist.Memberlist
	mu      sync.RWMutex
	nodes   []cluster.Node
}

func New(current *cluster.Node, nodes []cluster.Node) *GossipAgent {
	if current == nil {
		log.Fatal("Current node cannot be nil")
	}
	if len(nodes) == 0 {
		log.Fatal("Cluster nodes cannot be empty")
	}

	return &GossipAgent{
		Current: *current,
	}
}

func (a *GossipAgent) Start() error {
	config := memberlist.DefaultLocalConfig()
	config.Name = a.Current.GetDisplayName()
	config.BindPort = a.Current.GetPort()

	members, err := memberlist.Create(config)
	if err != nil {
		log.Fatal("Unable to establish memebership")
	}

	a.members = members
	membership := a.members.Members()
	nodes := make([]cluster.Node, 0, len(membership))
	for _, member := range membership {
		nodes = append(nodes, ToNode(member))
	}

	a.mu.Lock()
	a.nodes = nodes
	a.mu.Unlock()

	console.New().Print("Node successfully join cluster")
	return err
}

func (a *GossipAgent) Stop() error {
	if a.members == nil {
		return nil
	}
	return a.members.Shutdown()
}

func (a *GossipAgent) GetAllNodes() []cluster.Node {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return append([]cluster.Node(nil), a.nodes...)
}

func (a *GossipAgent) GetNodesByName(name string) []cluster.Node {
	return nil
}

func (a *GossipAgent) GetCurrentNode() cluster.Node {
	return a.Current
}

func ToNode(member *memberlist.Node) cluster.Node {
	nameSeparator := strings.LastIndex(member.Name, "_")
	node := cluster.Node{
		IP:            member.Addr.String(),
		Port:          int(member.Port),
		LastHeartbeat: time.Now().Unix(),
	}
	if nameSeparator >= 0 {
		node.ServiceName = member.Name[:nameSeparator]
		node.ID = member.Name[nameSeparator+1:]
	} else {
		node.ServiceName = member.Name
	}
	if len(member.Meta) > 0 {
		node.ServiceName = string(member.Meta)
	}
	return node
}
