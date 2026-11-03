package agent

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"flock/internal/cluster"
	"flock/internal/console"

	"github.com/hashicorp/memberlist"
)

type GossipAgent struct {
	CurrentNode    cluster.Node
	members        *memberlist.Memberlist
	mu             sync.RWMutex
	nodes          []cluster.Node
	bootstrapNodes []cluster.Node
}

func New(current *cluster.Node, nodes []cluster.Node) *GossipAgent {
	if current == nil {
		log.Fatal("Current node cannot be nil")
	}

	return &GossipAgent{
		CurrentNode:    *current,
		bootstrapNodes: append([]cluster.Node(nil), nodes...),
	}
}

func (agent *GossipAgent) Start() error {
	config := memberlist.DefaultLocalConfig()
	config.Name = agent.CurrentNode.GetDisplayName()
	config.BindPort = agent.CurrentNode.GetPort()
	config.Events = membershipEventHandler{agent: agent, console: console.New()}

	members, err := memberlist.Create(config)
	if err != nil {
		return fmt.Errorf("unable to establish membership: %w", err)
	}

	agent.members = members
	for _, member := range members.Members() {
		agent.upsertNode(ToNode(member))
	}

	addresses := make([]string, 0, len(agent.bootstrapNodes))
	for _, node := range agent.bootstrapNodes {
		if node.IP == "" || node.Port == 0 {
			_ = members.Shutdown()
			agent.members = nil
			return fmt.Errorf("invalid bootstrap node address for %s", node.GetDisplayName())
		}
		addresses = append(addresses, net.JoinHostPort(node.IP, strconv.Itoa(node.Port)))
	}
	if len(addresses) > 0 {
		if _, err := members.Join(addresses); err != nil {
			_ = members.Shutdown()
			agent.members = nil
			return fmt.Errorf("unable to join discovered nodes: %w", err)
		}
	}

	for _, member := range members.Members() {
		agent.upsertNode(ToNode(member))
	}

	output := console.New()
	output.Log(console.LogLevelInfo, "Node successfully join cluster")
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

type membershipEventHandler struct {
	agent   *GossipAgent
	console console.Console
}

func (handler membershipEventHandler) NotifyJoin(member *memberlist.Node) {
	handler.agent.upsertNode(ToNode(member))
	handler.console.Log(console.LogLevelInfo, fmt.Sprintf("Node joined: %s", member.Name))
}

func (handler membershipEventHandler) NotifyUpdate(member *memberlist.Node) {
	handler.agent.upsertNode(ToNode(member))
}

func (handler membershipEventHandler) NotifyLeave(member *memberlist.Node) {
	node := ToNode(member)

	handler.agent.mu.Lock()
	for i, existing := range handler.agent.nodes {
		if existing.ID == node.ID {
			handler.agent.nodes = append(handler.agent.nodes[:i], handler.agent.nodes[i+1:]...)
			break
		}
	}
	handler.agent.mu.Unlock()

	handler.console.Log(console.LogLevelWarn, fmt.Sprintf("Node leave: %s", member.Name))
}

func (agent *GossipAgent) upsertNode(node cluster.Node) {
	agent.mu.Lock()
	defer agent.mu.Unlock()

	for i, existing := range agent.nodes {
		if existing.ID == node.ID {
			agent.nodes[i] = node
			return
		}
	}
	agent.nodes = append(agent.nodes, node)
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
