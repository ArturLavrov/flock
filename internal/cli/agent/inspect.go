package agent

import "github.com/spf13/cobra"

func NewInspectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <agent-id>",
		Short: "Show agent information",
		Args:  cobra.ExactArgs(1),
	}
}
