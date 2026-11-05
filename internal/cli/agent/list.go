package agent

import "github.com/spf13/cobra"

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List local Flock agents",
		Args:    cobra.NoArgs,
	}
}
