package onepassword

import (
	"maps"
	"testing"
)

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
