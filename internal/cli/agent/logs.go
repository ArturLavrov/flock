package agent

import "github.com/spf13/cobra"

func NewLogsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logs <agent-id>",
		Short: "Show agent logs",
		Args:  cobra.ExactArgs(1),
	}
}
