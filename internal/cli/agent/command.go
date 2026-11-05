package agent

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "agent",
		Short: "Manage local Flock agents",
		Args:  cobra.NoArgs,
	}
	command.AddCommand(
		NewStartCommand(),
		NewListCommand(),
		NewStopCommand(),
		NewInspectCommand(),
		NewLogsCommand(),
	)

	return command
}
