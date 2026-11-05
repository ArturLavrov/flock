package agent

import "github.com/spf13/cobra"

func NewStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <agent-id>",
		Short: "Stop a running agent",
		Args:  cobra.ExactArgs(1),
	}
}
