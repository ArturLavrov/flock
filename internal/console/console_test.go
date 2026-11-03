package console

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pterm/pterm"
)

func TestPtermConsoleLogsMessageAtInfoLevel(t *testing.T) {
	var output bytes.Buffer
	previousWriter := pterm.DefaultLogger.Writer
	previousLevel := pterm.DefaultLogger.Level
	t.Cleanup(func() {
		pterm.DefaultLogger.Writer = previousWriter
		pterm.DefaultLogger.Level = previousLevel
	})

	pterm.DefaultLogger.Writer = &output
	pterm.DefaultLogger.Level = pterm.LogLevelInfo

	(PtermConsole{}).Log(LogLevelInfo, "console test message")

	if !strings.Contains(output.String(), "console test message") {
		t.Fatalf("logger output %q does not contain the message", output.String())
	}
}

func TestPtermConsoleLogUsesSelectedLevel(t *testing.T) {
	var output bytes.Buffer
	previousWriter := pterm.DefaultLogger.Writer
	previousLevel := pterm.DefaultLogger.Level
	t.Cleanup(func() {
		pterm.DefaultLogger.Writer = previousWriter
		pterm.DefaultLogger.Level = previousLevel
	})

	pterm.DefaultLogger.Writer = &output
	pterm.DefaultLogger.Level = pterm.LogLevelWarn

	console := PtermConsole{}
	console.Log(LogLevelDebug, "hidden debug message")
	console.Log(LogLevelInfo, "hidden info message")
	console.Log(LogLevelWarn, "visible warning message")
	console.Log(LogLevelError, "visible error message")

	logOutput := output.String()
	if strings.Contains(logOutput, "hidden") {
		t.Errorf("output includes messages below the configured level: %q", logOutput)
	}
	if !strings.Contains(logOutput, "visible warning message") || !strings.Contains(logOutput, "visible error message") {
		t.Errorf("output %q is missing warning or error messages", logOutput)
	}
}
