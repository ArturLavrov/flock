package agent

import "testing"

func TestStartCommandDoesNotExposeSeedNodeFlag(t *testing.T) {
	command := NewStartCommand()

	if flag := command.Flags().Lookup("seed-node"); flag != nil {
		t.Fatalf("seed-node flag is still registered: %v", flag)
	}
}
