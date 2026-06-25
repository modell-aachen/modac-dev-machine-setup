package setupenvs

import "testing"

func TestLastNonEmptyLine(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "single line",
			in:   "boom",
			want: "boom",
		},
		{
			name: "trailing blank lines are ignored",
			in:   "first\nlast error\n\n  \n",
			want: "last error",
		},
		{
			name: "op error on last line",
			in:   "[ERROR] 2026/06/25 \"GitHub Plugin Marketplace Token\" isn't an item in the \"Entwicklung\" vault",
			want: "[ERROR] 2026/06/25 \"GitHub Plugin Marketplace Token\" isn't an item in the \"Entwicklung\" vault",
		},
		{
			name: "empty input falls back",
			in:   "   \n\n",
			want: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lastNonEmptyLine(tt.in); got != tt.want {
				t.Errorf("lastNonEmptyLine() = %q, want %q", got, tt.want)
			}
		})
	}
}
