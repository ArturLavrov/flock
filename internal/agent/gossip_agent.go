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
	CurrentNode cluster.Node
	members     *memberlist.Memberlist
	mu          sync.RWMutex
	nodes       []cluster.Node
}

func New(current *cluster.Node, nodes []cluster.Node) *GossipAgent {
	if current == nil {
		log.Fatal("Current node cannot be nil")
	}
	if len(nodes) == 0 {
		log.Fatal("Cluster nodes cannot be empty")
	}

	return &GossipAgent{
		CurrentNode: *current,
	}
}

func (agent *GossipAgent) Start() error {
	config := memberlist.DefaultLocalConfig()
	config.Name = agent.CurrentNode.GetDisplayName()
	config.BindPort = agent.CurrentNode.GetPort()

	members, err := memberlist.Create(config)
	if err != nil {
		log.Fatal("Unable to establish memebership")
	}

	agent.members = members
	membership := agent.members.Members()
	nodes := make([]cluster.Node, 0, len(membership))
	for _, member := range membership {
		nodes = append(nodes, ToNode(member))
	}

	agent.nodes = nodes

	console := console.New()
	console.Print("Node successfully join cluster")
	return err
}

func (agent *GossipAgent) Stop() error {
	if agent.members == nil {
		return nil
	}
	return agent.members.Shutdown()
}

func (agent *GossipAgent) GetAllNodes() []cluster.Node {
	agent.mu.RLock()
	defer agent.mu.RUnlock()
	return append([]cluster.Node(nil), agent.nodes...)
}

func (agent *GossipAgent) GetNodesByName(name string) []cluster.Node {
	if name == "" {
		return agent.GetAllNodes()
	}

	agent.mu.RLock()
	defer agent.mu.RUnlock()

	nodes := make([]cluster.Node, 0)
	for _, node := range agent.nodes {
		if node.ServiceName == name {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (agent *GossipAgent) GetCurrentNode() cluster.Node {
	return agent.CurrentNode
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
