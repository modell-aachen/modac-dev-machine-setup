package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfileDefaultsToDevWhenFileMissing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	profile, err := LoadProfile()
	if err != nil {
		t.Fatalf("LoadProfile() error = %v", err)
	}
	if profile != ProfileDev {
		t.Errorf("LoadProfile() = %q, want %q", profile, ProfileDev)
	}
}

func TestSaveProfileRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
	}{
		{"dev profile", ProfileDev},
		{"service profile", ProfileService},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOME", t.TempDir())

			if err := SaveProfile(tt.profile); err != nil {
				t.Fatalf("SaveProfile(%q) error = %v", tt.profile, err)
			}

			got, err := LoadProfile()
			if err != nil {
				t.Fatalf("LoadProfile() error = %v", err)
			}
			if got != tt.profile {
				t.Errorf("LoadProfile() = %q, want %q", got, tt.profile)
			}
		})
	}
}

func TestLoadProfileTrimsWhitespace(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := filepath.Join(home, ".machine")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "profile"), []byte("service\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := LoadProfile()
	if err != nil {
		t.Fatalf("LoadProfile() error = %v", err)
	}
	if got != ProfileService {
		t.Errorf("LoadProfile() = %q, want %q", got, ProfileService)
	}
}

func TestLoadProfileRejectsUnknownContent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	dir := filepath.Join(home, ".machine")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "profile"), []byte("laptop"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := LoadProfile(); err == nil {
		t.Error("LoadProfile() expected error for unknown profile content, got nil")
	}
}

func TestParseProfile(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Profile
		wantErr bool
	}{
		{"dev", "dev", ProfileDev, false},
		{"service", "service", ProfileService, false},
		{"unknown value", "laptop", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProfile(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseProfile(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseProfile(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
