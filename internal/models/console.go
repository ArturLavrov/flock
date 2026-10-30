package models

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/pterm/pterm"
)

type AgentInfo struct {
	Version  string
	NodeID   string
	Address  string
	Protocol string
	Status   string
}

var consoleOutputMu sync.Mutex

func PrintAgentInfo(info AgentInfo) error {
	return writeAgentInfo(os.Stdout, info, terminalColorsEnabled(), pterm.GetTerminalWidth())
}

func PrintMembershipEvent(event MembershipEvent) error {
	return writeMembershipEvent(os.Stdout, event, terminalColorsEnabled(), pterm.GetTerminalWidth())
}

func writeAgentInfo(w io.Writer, info AgentInfo, colors bool, terminalWidth int) error {
	consoleOutputMu.Lock()
	defer consoleOutputMu.Unlock()

	restoreColors := setColors(colors)
	defer restoreColors()

	_, err := io.WriteString(w, formatAgentInfo(info, terminalWidth))
	return err
}

func writeMembershipEvent(w io.Writer, event MembershipEvent, colors bool, terminalWidth int) error {
	consoleOutputMu.Lock()
	defer consoleOutputMu.Unlock()

	restoreColors := setColors(colors)
	defer restoreColors()

	_, err := io.WriteString(w, formatMembershipEvent(event, terminalWidth))
	return err
}

func setColors(enabled bool) func() {
	previous := pterm.PrintColor
	if enabled {
		pterm.EnableColor()
	} else {
		pterm.DisableColor()
	}

	return func() {
		if previous {
			pterm.EnableColor()
		} else {
			pterm.DisableColor()
		}
	}
}

func terminalColorsEnabled() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return false
	}

	_, _, err := pterm.GetTerminalSize()
	return err == nil
}

func formatAgentInfo(info AgentInfo, terminalWidth int) string {
	if terminalWidth <= 0 {
		terminalWidth = 80
	}
	if info.Status == "" {
		info.Status = "ALIVE"
	}

	const labelWidth = 12
	maxContentWidth := terminalWidth - 6
	if maxContentWidth < 8 {
		maxContentWidth = 8
	}
	valueWidth := maxContentWidth - labelWidth - 1
	if valueWidth < 1 {
		valueWidth = 1
	}

	fields := []struct {
		label string
		value string
		color pterm.Color
	}{
		{label: "Node ID", value: info.NodeID, color: pterm.FgWhite},
		{label: "Address", value: info.Address, color: pterm.FgWhite},
		{label: "Protocol", value: info.Protocol, color: pterm.FgWhite},
		{label: "Status", value: "● " + info.Status, color: pterm.FgGreen},
	}

	lines := make([]string, 0, len(fields)+2)
	heading := truncateRunes("NODE INFORMATION", maxContentWidth-1)
	lines = append(lines, " "+pterm.NewStyle(pterm.FgGray, pterm.Bold).Sprint(heading))
	lines = append(lines, "")
	for _, field := range fields {
		value := truncateRunes(field.value, valueWidth)
		label := pterm.NewStyle(pterm.FgGray).Sprint(fmt.Sprintf("%-*s", labelWidth, field.label))
		lines = append(lines, " "+label+pterm.NewStyle(field.color).Sprint(value))
	}

	box := pterm.BoxPrinter{
		BoxStyle:                pterm.NewStyle(pterm.FgGray),
		TextStyle:               pterm.NewStyle(),
		VerticalString:          "│",
		HorizontalString:        "─",
		TopLeftCornerString:     "╯",
		TopRightCornerString:    "╰",
		BottomLeftCornerString:  "╮",
		BottomRightCornerString: "╭",
		LeftPadding:             1,
		RightPadding:            1,
	}

	subtitle := "Distributed cluster membership"
	if info.Version != "" {
		subtitle += " · " + info.Version
	}
	subtitle = truncateRunes(subtitle, terminalWidth-2)
	readyText := truncateRunes("✓ Agent ready · Waiting for peers...", terminalWidth-2)
	ready := pterm.NewStyle(pterm.FgGreen).Sprint(readyText)

	var output strings.Builder
	output.WriteString("  ")
	output.WriteString(pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold).Sprint("◆ FLOCK"))
	output.WriteString("\n  ")
	output.WriteString(pterm.NewStyle(pterm.FgGray).Sprint(subtitle))
	output.WriteString("\n\n")
	output.WriteString(indentLines(box.Sprint(strings.Join(lines, "\n")), "  "))
	output.WriteString("\n\n  ")
	output.WriteString(ready)
	output.WriteString("\n\n  ")
	output.WriteString(pterm.NewStyle(pterm.FgGray, pterm.Bold).Sprint("CLUSTER EVENTS"))
	output.WriteString("\n")
	return output.String()
}

func indentLines(value, indentation string) string {
	lines := strings.Split(value, "\n")
	for i := range lines {
		lines[i] = indentation + lines[i]
	}
	return strings.Join(lines, "\n")
}

func truncateRunes(value string, maxWidth int) string {
	if maxWidth < 1 {
		return ""
	}
	if utf8.RuneCountInString(value) <= maxWidth {
		return value
	}
	if maxWidth == 1 {
		return "…"
	}
	runes := []rune(value)
	return string(runes[:maxWidth-1]) + "…"
}

func formatMembershipEvent(event MembershipEvent, terminalWidth int) string {
	timestamp := event.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}
	if terminalWidth <= 0 {
		terminalWidth = 80
	}

	label, color, message := eventDescription(event)
	message = truncateRunes(message, terminalWidth-22)
	line := fmt.Sprintf("%s  %-7s %s",
		pterm.NewStyle(pterm.FgGray).Sprint(timestamp.Format("15:04:05")),
		pterm.NewStyle(color).Sprint(label),
		message,
	)
	return "  " + line + "\n"
}

func eventDescription(event MembershipEvent) (string, pterm.Color, string) {
	node := event.NodeID
	if node == "" {
		node = "unknown node"
	}

	switch strings.ToUpper(event.Type) {
	case "JOIN":
		message := node + " joined"
		if event.Address != "" {
			message += " · " + event.Address
		}
		return "JOIN", pterm.FgGreen, message
	case "SUSPECT":
		return "SUSPECT", pterm.FgYellow, messageWithDetails(node, event.Details, "may be unreachable")
	case "ALIVE":
		return "ALIVE", pterm.FgGreen, messageWithDetails(node, event.Details, "recovered")
	case "DEAD":
		return "DEAD", pterm.FgRed, messageWithDetails(node, event.Details, "failed")
	case "LEAVE":
		message := node + " left"
		if event.Address != "" {
			message += " · " + event.Address
		}
		return "LEAVE", pterm.FgBlue, message
	default:
		return "UNKNOWN", pterm.FgGray, fmt.Sprintf("%s event type %q is unsupported", node, event.Type)
	}
}

func messageWithDetails(node, details, fallback string) string {
	if details == "" {
		details = fallback
	}
	return node + " " + details
}
