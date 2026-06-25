package provision

import "testing"

func TestBinaryChanged(t *testing.T) {
	tests := []struct {
		name     string
		self     string
		resolved string
		want     bool
	}{
		{"same path means unchanged", "/nix/store/aaa/bin/machine", "/nix/store/aaa/bin/machine", false},
		{"different path means changed", "/nix/store/aaa/bin/machine", "/nix/store/bbb/bin/machine", true},
		{"unresolvable post-update binary counts as unchanged", "/nix/store/aaa/bin/machine", "", false},
		{"unknown current binary forces a reload", "", "/nix/store/bbb/bin/machine", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binaryChanged(tt.self, tt.resolved); got != tt.want {
				t.Errorf("binaryChanged(%q, %q) = %v, want %v", tt.self, tt.resolved, got, tt.want)
			}
		})
	}
}
