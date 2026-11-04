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

func (PtermConsole) PrintBox(nodeInfo NodeInfo) {
	labelStyle := pterm.NewStyle(pterm.FgGray)
	valueStyle := pterm.NewStyle(pterm.FgWhite)
	statusStyle := pterm.NewStyle(pterm.FgGreen)
	rows := []string{
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Node ID"), valueStyle.Sprint(nodeInfo.NodeID)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Node Name"), valueStyle.Sprint(nodeInfo.NodeName)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Bootstrap"), valueStyle.Sprint(nodeInfo.IsBootstrapNode)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Gossip Port"), valueStyle.Sprint(nodeInfo.GossipPort)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "API Port"), valueStyle.Sprint(nodeInfo.ApiPort)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Protocol"), valueStyle.Sprint(nodeInfo.Protocol)),
		fmt.Sprintf("%s  %s", labelStyle.Sprintf("%-12s", "Status"), statusStyle.Sprint("● "+nodeInfo.Status)),
	}

	box := pterm.DefaultBox
	box.Title = "NODE INFORMATION"
	box.BoxStyle = pterm.NewStyle(pterm.FgGray)
	box.TextStyle = pterm.NewStyle(pterm.FgWhite)
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
