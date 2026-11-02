package console

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pterm/pterm"
)

func TestPtermConsolePrintsMessageAtInfoLevel(t *testing.T) {
	var output bytes.Buffer
	previousWriter := pterm.DefaultLogger.Writer
	previousLevel := pterm.DefaultLogger.Level
	t.Cleanup(func() {
		pterm.DefaultLogger.Writer = previousWriter
		pterm.DefaultLogger.Level = previousLevel
	})

	pterm.DefaultLogger.Writer = &output
	pterm.DefaultLogger.Level = pterm.LogLevelInfo

	(PtermConsole{}).Print("console test message")

	if !strings.Contains(output.String(), "console test message") {
		t.Fatalf("logger output %q does not contain the message", output.String())
	}
}
