package agent

import (
	"fmt"
	"net/http"

	gossip "flock/internal/agent"
	"flock/internal/api"
	"flock/internal/cluster"
	"flock/internal/console"

	"github.com/spf13/cobra"
)

const (
	MinTCPPortNumber = 0
	MaxTcpPortNumber = 65535
)

func NewStartCommand() *cobra.Command {
	var serviceName string
	var gossipPort int
	var apiPort int
	var isSeedNode bool

	command := &cobra.Command{
		Use:   "start",
		Short: "Start a new Flock agent",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if serviceName == "" {
				return fmt.Errorf("service name cannot be empty")
			}
			if gossipPort < MinTCPPortNumber || gossipPort > MaxTcpPortNumber {
				return fmt.Errorf("invalid gossip port: %d (valid range: 1-65535)", gossipPort)
			}
			if apiPort < MinTCPPortNumber || apiPort > MaxTcpPortNumber {
				return fmt.Errorf("invalid API port: %d (valid range: 1-65535)", apiPort)
			}

			currentNode := cluster.NewNode(serviceName, gossipPort, false)
			var err error
			var nodes []cluster.Node
			nodes, err = cluster.GetAvailableNodes()
			if err != nil {
				return fmt.Errorf("unable to discover cluster nodes: %w", err)
			}

			agent := gossip.New(&currentNode, nodes, apiPort)
			if err := agent.Start(); err != nil {
				return fmt.Errorf("unable to start agent: %w", err)
			}

			apiAddress := fmt.Sprintf(":%d", apiPort)
			output := console.New()
			output.Log(console.LogLevelInfo, fmt.Sprintf("Agent HTTP API running on %s", apiAddress))

			if err := http.ListenAndServe(apiAddress, api.HTTPHandler(agent)); err != nil {
				return fmt.Errorf("agent HTTP API stopped: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&serviceName, "service-name", "", "Name of the service")
	command.Flags().IntVar(&gossipPort, "gossip-port", 0, "Gossip port")
	command.Flags().IntVar(&apiPort, "api-port", 0, "API port")
	command.Flags().BoolVar(&isSeedNode, "seed-node", false, "Start this agent as a seed node")
	command.Flags().BoolP("detach", "d", false, "Run the agent in the background")
	if err := command.MarkFlagRequired("service-name"); err != nil {
		panic(err)
	}

	return command
}
