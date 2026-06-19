package nixconf

import "testing"

func TestHasRequiredFeatures(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"empty file", "", false},
		{"exact line", "experimental-features = nix-command flakes\n", true},
		{"reordered features", "experimental-features = flakes nix-command\n", true},
		{"extra features", "experimental-features = nix-command flakes ca-derivations\n", true},
		{"no spaces around equals", "experimental-features=nix-command flakes\n", true},
		{"only nix-command", "experimental-features = nix-command\n", false},
		{"only flakes", "experimental-features = flakes\n", false},
		{"unrelated key", "max-jobs = 4\n", false},
		{"commented out", "# experimental-features = nix-command flakes\n", false},
		{"set among other lines", "max-jobs = 4\nexperimental-features = nix-command flakes\nsandbox = true\n", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasRequiredFeatures(tt.content); got != tt.want {
				t.Errorf("hasRequiredFeatures(%q) = %v, want %v", tt.content, got, tt.want)
			}
		})
	}
}
