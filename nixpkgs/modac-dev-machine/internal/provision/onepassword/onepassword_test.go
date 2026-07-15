package onepassword

import (
	"maps"
	"testing"
)

func TestAptSourceLine(t *testing.T) {
	tests := []struct {
		name string
		arch string
		want string
	}{
		{
			name: "amd64",
			arch: "amd64",
			want: "deb [arch=amd64 signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/amd64 stable main",
		},
		{
			name: "arm64",
			arch: "arm64",
			want: "deb [arch=arm64 signed-by=/usr/share/keyrings/1password-archive-keyring.gpg] https://downloads.1password.com/linux/debian/arm64 stable main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aptSourceLine(tt.arch); got != tt.want {
				t.Errorf("aptSourceLine(%q) = %q, want %q", tt.arch, got, tt.want)
			}
		})
	}
}

func TestParseSessionExports(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   map[string]string
	}{
		{
			name:   "single export line",
			output: `export OP_SESSION_my="abc123"` + "\n",
			want:   map[string]string{"OP_SESSION_my": "abc123"},
		},
		{
			name: "export line with trailing comment output",
			output: `export OP_SESSION_modellaachen="tok-en_value"` + "\n" +
				`# This command is meant to be used with your shell's eval function.` + "\n",
			want: map[string]string{"OP_SESSION_modellaachen": "tok-en_value"},
		},
		{
			name:   "no export lines",
			output: "some unrelated output\n",
			want:   map[string]string{},
		},
		{
			name:   "empty output",
			output: "",
			want:   map[string]string{},
		},
		{
			name:   "ignores non-session exports",
			output: `export PATH="/usr/bin"` + "\n" + `export OP_SESSION_x="y"` + "\n",
			want:   map[string]string{"OP_SESSION_x": "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSessionExports(tt.output)
			if !maps.Equal(got, tt.want) {
				t.Errorf("parseSessionExports(%q) = %v, want %v", tt.output, got, tt.want)
			}
		})
	}
}
