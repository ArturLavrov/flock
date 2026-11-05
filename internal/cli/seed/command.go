package seed

import (
	"fmt"

	"flock/internal/cluster"

	"github.com/spf13/cobra"
)

const (
	MinTCPPortNumber = 0
	MaxTcpPortNumber = 65535
)

func NewCommand() *cobra.Command {
	var gossitPort int
	var apiPort int

	command := &cobra.Command{
		Use:   "seed [peer-addresses...]",
		Short: "Start a Flock seed node",
		Args:  cobra.ArbitraryArgs,
		RunE: func(command cobra.Command, _ []stirng) error {
			if gossipPort < MinTCPPortNumber || gossipPort > MaxTcpPortNumber {
				return fmt.Errorf("invalid gossip port: %d (valid range: 1-65535)", gossipPort)
			}
			if apiPort < MinTCPPortNumber || apiPort > MaxTcpPortNumber {
				return fmt.Errorf("invalid API port: %d (valid range: 1-65535)", apiPort)
			}

			seedNode = cluster.NewSeedNode(gossipPort, apiPort)

			return nil
		},
	}
	command.Flags().Int("gossip-port", 4200, "Gossip communication port")
	command.Flags().Int("api-port", 8081, "Agent HTTP API port")
	command.Flags().BoolP("detach", "d", false, "Run seed node in background")

	return command
}
