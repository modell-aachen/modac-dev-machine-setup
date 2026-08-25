package dockerpackages

import "testing"

func TestDockerRepoLine(t *testing.T) {
	want := "deb [arch=amd64 signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu noble stable"
	if got := dockerRepoLine("amd64", "noble"); got != want {
		t.Errorf("dockerRepoLine() = %q, want %q", got, want)
	}
}

func TestNeedsArchPin(t *testing.T) {
	tests := []struct {
		name string
		repo string
		want bool
	}{
		{"generated entry", dockerRepoLine("amd64", "noble") + "\n", false},
		{"unpinned entry", "deb [signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu noble stable\n", true},
		{"hand-pinned entry with reordered options", "deb [signed-by=/etc/apt/keyrings/docker.asc arch=arm64] https://download.docker.com/linux/ubuntu noble stable\n", false},
		{"empty file", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := needsArchPin(tt.repo); got != tt.want {
				t.Errorf("needsArchPin() = %v, want %v", got, tt.want)
			}
		})
	}
}
