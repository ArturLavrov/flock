package models

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

func TestAgentInfoFormattingWithoutColors(t *testing.T) {
	var output bytes.Buffer
	err := writeAgentInfo(&output, AgentInfo{
		Version:  "v0.1.0",
		NodeID:   "node-a83f",
		Address:  "127.0.0.1:7001",
		Protocol: "SWIM / Gossip",
		Status:   "ALIVE",
	}, false, 80)
	if err != nil {
		t.Fatal(err)
	}

	got := output.String()
	for _, want := range []string{
		"◆ FLOCK",
		"Distributed cluster membership · v0.1.0",
		"╭", "╮", "╰", "╯",
		"NODE INFORMATION",
		"node-a83f",
		"127.0.0.1:7001",
		"SWIM / Gossip",
		"● ALIVE",
		"Agent ready · Waiting for peers...",
		"CLUSTER EVENTS",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("startup output missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "\x1b[") {
		t.Errorf("color-disabled output contains ANSI escape codes: %q", got)
	}
	lines := strings.Split(got, "\n")
	if !strings.HasPrefix(lines[3], "  ╭") || !strings.HasSuffix(lines[3], "╮") {
		t.Errorf("box top border has incorrect corners: %q", lines[3])
	}
	bottom := -1
	for i, line := range lines {
		if strings.HasPrefix(line, "  ╰") {
			bottom = i
			break
		}
	}
	if bottom < 0 {
		t.Fatalf("box bottom border not found:\n%s", got)
	}
	if !strings.HasSuffix(lines[bottom], "╯") {
		t.Errorf("box bottom border has incorrect corners: %q", lines[bottom])
	}
	if !strings.Contains(got, "Node ID     node-a83f") ||
		!strings.Contains(got, "Address     127.0.0.1:7001") {
		t.Errorf("node values are not aligned in two columns:\n%s", got)
	}
}

func TestAgentInfoFitsNarrowTerminal(t *testing.T) {
	const width = 36
	var formatted bytes.Buffer
	if err := writeAgentInfo(&formatted, AgentInfo{
		Version:  "v0.1.0",
		NodeID:   "node-with-a-long-identifier",
		Address:  "2001:db8::1:7001",
		Protocol: "SWIM / Gossip",
		Status:   "ALIVE",
	}, false, width); err != nil {
		t.Fatal(err)
	}

	for _, line := range strings.Split(strings.TrimSuffix(formatted.String(), "\n"), "\n") {
		if got := utf8.RuneCountInString(line); got > width {
			t.Errorf("startup line has width %d, exceeds terminal width %d: %q", got, width, line)
		}
	}
}

func TestMembershipEventFormatting(t *testing.T) {
	tests := []struct {
		name  string
		event MembershipEvent
		want  []string
	}{
		{
			name: "join",
			event: MembershipEvent{
				Type: "JOIN", NodeID: "node-b19c", Address: "127.0.0.1:7002",
				Timestamp: time.Date(2026, 10, 28, 21, 24, 5, 0, time.UTC),
			},
			want: []string{"21:24:05", "JOIN", "node-b19c joined · 127.0.0.1:7002"},
		},
		{
			name:  "suspect",
			event: MembershipEvent{Type: "SUSPECT", NodeID: "node-b19c", Details: "probe timeout"},
			want:  []string{"SUSPECT", "node-b19c probe timeout"},
		},
		{
			name:  "alive",
			event: MembershipEvent{Type: "ALIVE", NodeID: "node-b19c"},
			want:  []string{"ALIVE", "node-b19c recovered"},
		},
		{
			name:  "dead",
			event: MembershipEvent{Type: "DEAD", NodeID: "node-c52d"},
			want:  []string{"DEAD", "node-c52d failed"},
		},
		{
			name:  "leave",
			event: MembershipEvent{Type: "LEAVE", NodeID: "node-b19c"},
			want:  []string{"LEAVE", "node-b19c left"},
		},
		{
			name:  "unknown",
			event: MembershipEvent{Type: "REJOIN", NodeID: "node-a"},
			want:  []string{"UNKNOWN", `node-a event type "REJOIN" is unsupported`},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			if err := writeMembershipEvent(&output, test.event, false, 80); err != nil {
				t.Fatal(err)
			}
			for _, want := range test.want {
				if !strings.Contains(output.String(), want) {
					t.Errorf("event output missing %q: %q", want, output.String())
				}
			}
			if strings.Contains(output.String(), "\x1b[") {
				t.Errorf("color-disabled event output contains ANSI codes: %q", output.String())
			}
		})
	}
}

func TestConcurrentMembershipOutputDoesNotInterleave(t *testing.T) {
	const eventCount = 100
	var output bytes.Buffer
	var wait sync.WaitGroup

	for i := 0; i < eventCount; i++ {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			event := MembershipEvent{
				Type:      "JOIN",
				NodeID:    fmt.Sprintf("node-%03d", i),
				Timestamp: time.Date(2026, 10, 28, 21, 24, 5, 0, time.UTC),
			}
			if err := writeMembershipEvent(&output, event, false, 80); err != nil {
				t.Errorf("write event: %v", err)
			}
		}(i)
	}
	wait.Wait()

	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != eventCount {
		t.Fatalf("got %d event lines, want %d", len(lines), eventCount)
	}
	for _, line := range lines {
		if !strings.Contains(line, "21:24:05  JOIN    node-") || !strings.HasSuffix(line, " joined") {
			t.Errorf("malformed or interleaved event line: %q", line)
		}
	}
}
