package console

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

type PtermConsole struct{}

func (console PtermConsole) Print(message string) {
	pterm.DefaultBasicText.Println(message)
}

func (PtermConsole) PrintFlockTitle() {
	pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold).Printfln("◆ FLOCK")
	pterm.NewStyle(pterm.FgGray).Println("Distributed cluster membership · v0.1.0")
}

func (PtermConsole) PrintAgentInfo(agentID, serviceName string, gossipPort, apiPort int, protocol, version, status string, isSeedNodeAgent bool) {
	labelStyle := pterm.NewStyle(pterm.FgGreen)
	valueStyle := pterm.NewStyle(pterm.FgWhite)
	boxStyle := pterm.NewStyle(pterm.FgGray)
	textStyle := pterm.NewStyle(pterm.FgWhite)
	if isSeedNodeAgent {
		valueStyle = pterm.NewStyle(pterm.FgGreen)
		boxStyle = pterm.NewStyle(pterm.FgGreen)
		textStyle = pterm.NewStyle(pterm.FgGreen)
	}
	statusStyle := pterm.NewStyle(pterm.FgGreen)
	rows := []string{
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Agent ID"), valueStyle.Sprint(agentID)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Service Name"), valueStyle.Sprint(serviceName)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Gossip Port"), valueStyle.Sprint(gossipPort)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "API Port"), valueStyle.Sprint(apiPort)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Protocol"), valueStyle.Sprint(protocol)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Version"), valueStyle.Sprint(version)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Status"), statusStyle.Sprint("● "+status)),
	}

	box := pterm.DefaultBox
	box.Title = "AGENT INFORMATION"
	box.BoxStyle = boxStyle
	box.TextStyle = textStyle
	box.LeftPadding = 1
	box.RightPadding = 1
	box.TopPadding = 1
	box.BottomPadding = 1
	box.Println(strings.Join(rows, "\n"))
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
