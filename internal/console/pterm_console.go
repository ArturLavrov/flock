package console

import "github.com/pterm/pterm"

type PtermConsole struct{}

func (PtermConsole) Print(message string) {
	pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo).Info(message)
}
