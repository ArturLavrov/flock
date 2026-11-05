package cli

import (
	"github.com/spf13/cobra"

	"flock/internal/cli/agent"
)

const helpTemplate = `{{with .Short}}{{.}}

{{end}}Usage:
  {{.UseLine}}
{{if .HasSubCommands}}
Commands:
{{range .Commands}}{{if and (not .Hidden) (ne .Name "help")}}  {{rpad .Name .NamePadding}} {{.Short}}
{{end}}{{end}}{{end}}{{if .LocalFlags.HasAvailableFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
{{end}}`

func NewRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "flock",
		Short: "Decentralized service discovery",
	}
	rootCommand.SetHelpTemplate(helpTemplate)
	rootCommand.AddCommand(agent.NewCommand())

	return rootCommand
}
