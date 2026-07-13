package platform

import "testing"

func TestIsWSL(t *testing.T) {
	tests := []struct {
		name          string
		distroName    string
		kernelRelease string
		want          bool
	}{
		{"WSL_DISTRO_NAME set", "Ubuntu", "", true},
		{"microsoft kernel (WSL2)", "", "5.15.167.4-microsoft-standard-WSL2", true},
		{"microsoft kernel uppercase (WSL1)", "", "4.4.0-19041-Microsoft", true},
		{"both signals", "Ubuntu", "5.15.167.4-microsoft-standard-WSL2", true},
		{"plain Ubuntu kernel", "", "6.8.0-45-generic", false},
		{"no signals", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWSL(tt.distroName, tt.kernelRelease); got != tt.want {
				t.Errorf("isWSL(%q, %q) = %v, want %v", tt.distroName, tt.kernelRelease, got, tt.want)
			}
		})
	}
}
