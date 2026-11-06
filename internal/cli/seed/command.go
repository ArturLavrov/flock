package seed

import (
	"fmt"
	"net/http"

	gossipAgent "flock/internal/agent"
	"flock/internal/api"
	"flock/internal/cluster"
	"flock/internal/console"

	"github.com/spf13/cobra"
)

const (
	MinTCPPortNumber = 0
	MaxTcpPortNumber = 65535
)

func NewCommand() *cobra.Command {
	var gossipPort int
	var apiPort int

	command := &cobra.Command{
		Use:   "seed [peer-addresses...]",
		Short: "Start a Flock seed node",
		Args:  cobra.ArbitraryArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if gossipPort < MinTCPPortNumber || gossipPort > MaxTcpPortNumber {
				return fmt.Errorf("invalid gossip port: %d (valid range: 1-65535)", gossipPort)
			}
			if apiPort < MinTCPPortNumber || apiPort > MaxTcpPortNumber {
				return fmt.Errorf("invalid API port: %d (valid range: 1-65535)", apiPort)
			}

			seedNode := cluster.NewSeedNode(gossipPort)
			clusterNodes := []cluster.Node{}
			agent := gossipAgent.New(&seedNode, clusterNodes, apiPort)
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
	command.Flags().Int("gossip-port", 4200, "Gossip communication port")
	command.Flags().Int("api-port", 8081, "Agent HTTP API port")
	command.Flags().BoolP("detach", "d", false, "Run seed node in background")

	return command
}
