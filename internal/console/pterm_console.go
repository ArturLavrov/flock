package console

import (
	"fmt"

	"github.com/pterm/pterm"
)

type PtermConsole struct{}

func (console PtermConsole) Print(message string) {
	pterm.DefaultBasicText.Println(message)
}

func (PtermConsole) Log(level LogLevel, text string) {
	logger := pterm.DefaultLogger
	switch level {
	case LogLevelDebug:
		logger.Debug(text)
	case LogLevelInfo:
		logger.Info(text)
	case LogLevelWarn:
		logger.Warn(text)
	case LogLevelError:
		logger.Error(text)
	default:
		logger.Error(fmt.Sprintf("Invalid log level %d: %s", level, text))
	}
}
