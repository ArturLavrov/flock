package seed

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "seed [peer-addresses...]",
		Short: "Start a Flock seed node",
		Args:  cobra.ArbitraryArgs,
	}
	command.Flags().Int("gossip-port", 4200, "Gossip communication port")
	command.Flags().Int("api-port", 8081, "Agent HTTP API port")
	command.Flags().BoolP("detach", "d", false, "Run seed node in background")

	return command
}
