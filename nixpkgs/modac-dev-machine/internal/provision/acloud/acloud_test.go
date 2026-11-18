package acloud

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/modell-aachen/machine/internal/output"
)

const (
	latestVersion = "1.2.0"
	newBinary     = "#!/bin/sh\necho acloud version 1.2.0\n"
	oldBinary     = "old acloud"
	testToken     = "secret-token"
)

type nexusFake struct {
	files map[string]string
	codes map[string]int

	mu       sync.Mutex
	requests []*http.Request
}

func newNexusFake(t *testing.T, archive []byte, checksums string) *nexusFake {
	t.Helper()
	return &nexusFake{
		files: map[string]string{
			"/VERSION":          latestVersion + "\n",
			"/checksums.txt":    checksums,
			"/" + archiveName(): string(archive),
		},
		codes: map[string]int{},
	}
}

func (f *nexusFake) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	f.requests = append(f.requests, r.Clone(r.Context()))
	f.mu.Unlock()
	if code, ok := f.codes[r.URL.Path]; ok {
		w.WriteHeader(code)
		return
	}
	body, ok := f.files[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	_, _ = w.Write([]byte(body))
}

func (f *nexusFake) requestedPaths() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	paths := make([]string, len(f.requests))
	for i, r := range f.requests {
		paths[i] = r.URL.Path
	}
	return paths
}

func buildArchive(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, body := range entries {
		header := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(body)), Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("write tar header: %v", err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatalf("write tar body: %v", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

func checksumLine(archive []byte) string {
	sum := sha256.Sum256(archive)
	return fmt.Sprintf("%s  %s\n", hex.EncodeToString(sum[:]), archiveName())
}

func validChecksums(archive []byte) string {
	return "deadbeef  acloud_other_arch.tar.gz\n" + checksumLine(archive)
}

func fixedVersion(version string) func(string) string {
	return func(string) string { return version }
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func dirEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names
}

func TestInstall(t *testing.T) {
	archive := buildArchive(t, map[string]string{"README.md": "readme", binaryName: newBinary})
	tamperedChecksums := strings.Replace(checksumLine(archive), checksumLine(archive)[:8], "00000000", 1)

	tests := []struct {
		name             string
		existingBinary   string
		installedVersion string
		token            string
		checksums        string
		codes            map[string]int
		wantVersion      string
		wantErr          []string
		wantBinary       string
		wantNoArchive    bool
		wantAuth         bool
		versionFile      string
	}{
		{
			name:        "fresh install puts the binary into the bin directory",
			checksums:   validChecksums(archive),
			wantVersion: latestVersion,
			wantBinary:  newBinary,
		},
		{
			name:             "already latest skips the download",
			existingBinary:   oldBinary,
			installedVersion: latestVersion,
			checksums:        validChecksums(archive),
			wantVersion:      latestVersion,
			wantBinary:       oldBinary,
			wantNoArchive:    true,
		},
		{
			name:             "older installed version is replaced",
			existingBinary:   oldBinary,
			installedVersion: "1.1.0",
			checksums:        validChecksums(archive),
			wantVersion:      latestVersion,
			wantBinary:       newBinary,
		},
		{
			name:             "checksum mismatch aborts and keeps the old binary",
			existingBinary:   oldBinary,
			installedVersion: "1.1.0",
			checksums:        tamperedChecksums,
			wantErr:          []string{"checksum mismatch", checksumLine(archive)[:64], "00000000"},
			wantBinary:       oldBinary,
		},
		{
			name:      "missing checksum entry aborts",
			checksums: "deadbeef  acloud_other_arch.tar.gz\n",
			wantErr:   []string{"no checksum", archiveName()},
		},
		{
			name:        "unpublished release is skipped without error",
			checksums:   validChecksums(archive),
			codes:       map[string]int{"/VERSION": http.StatusNotFound},
			wantVersion: "",
		},
		{
			name:      "unauthorized names the token variable",
			checksums: validChecksums(archive),
			codes:     map[string]int{"/VERSION": http.StatusUnauthorized},
			wantErr:   []string{tokenEnv},
		},
		{
			name:      "forbidden archive names the token variable",
			checksums: validChecksums(archive),
			codes:     map[string]int{"/" + archiveName(): http.StatusForbidden},
			wantErr:   []string{tokenEnv},
		},
		{
			name:        "empty VERSION file fails",
			checksums:   validChecksums(archive),
			versionFile: " \n",
			wantErr:     []string{"empty"},
		},
		{
			name:      "server error fails with the status",
			checksums: validChecksums(archive),
			codes:     map[string]int{"/checksums.txt": http.StatusInternalServerError},
			wantErr:   []string{"500"},
		},
		{
			name:        "token is sent as basic auth",
			token:       testToken,
			checksums:   validChecksums(archive),
			wantVersion: latestVersion,
			wantBinary:  newBinary,
			wantAuth:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := newNexusFake(t, archive, tt.checksums)
			if tt.versionFile != "" {
				fake.files["/VERSION"] = tt.versionFile
			}
			for path, code := range tt.codes {
				fake.codes[path] = code
			}
			server := httptest.NewServer(fake)
			defer server.Close()

			binDir := filepath.Join(t.TempDir(), "bin")
			binPath := filepath.Join(binDir, binaryName)
			if tt.existingBinary != "" {
				if err := os.MkdirAll(binDir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(binPath, []byte(tt.existingBinary), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			inst := installer{
				client:           server.Client(),
				baseURL:          server.URL + "/",
				token:            tt.token,
				binDir:           binDir,
				installedVersion: fixedVersion(tt.installedVersion),
			}
			version, err := inst.install(&output.Context{})

			if len(tt.wantErr) > 0 {
				if err == nil {
					t.Fatalf("install() error = nil, want error containing %v", tt.wantErr)
				}
				for _, want := range tt.wantErr {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("install() error = %q, want it to contain %q", err, want)
					}
				}
			} else if err != nil {
				t.Fatalf("install() error = %v", err)
			}

			if version != tt.wantVersion {
				t.Errorf("install() version = %q, want %q", version, tt.wantVersion)
			}

			if tt.wantBinary != "" {
				if got := readFile(t, binPath); got != tt.wantBinary {
					t.Errorf("binary content = %q, want %q", got, tt.wantBinary)
				}
				info, err := os.Stat(binPath)
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode().Perm()&0o111 == 0 {
					t.Errorf("binary mode = %v, want executable", info.Mode())
				}
			}

			if names, err := os.ReadDir(binDir); err == nil && len(names) > 0 {
				if got := dirEntries(t, binDir); len(got) != 1 || got[0] != binaryName {
					t.Errorf("bin directory = %v, want only %q", got, binaryName)
				}
			}

			for _, path := range fake.requestedPaths() {
				if tt.wantNoArchive && path != "/VERSION" {
					t.Errorf("requested %s, want only /VERSION", path)
				}
			}

			for _, r := range fake.requests {
				user, pass, ok := r.BasicAuth()
				if tt.wantAuth && (!ok || user != nexusUser || pass != testToken) {
					t.Errorf("%s auth = (%q, %q, %v), want (%q, %q)", r.URL.Path, user, pass, ok, nexusUser, testToken)
				}
				if !tt.wantAuth && ok {
					t.Errorf("%s sent basic auth without a token", r.URL.Path)
				}
			}
		})
	}
}

func TestInstallRejectsArchiveWithoutBinary(t *testing.T) {
	archive := buildArchive(t, map[string]string{"README.md": "readme"})
	fake := newNexusFake(t, archive, checksumLine(archive))
	server := httptest.NewServer(fake)
	defer server.Close()
	binDir := t.TempDir()

	inst := installer{
		client:           server.Client(),
		baseURL:          server.URL + "/",
		binDir:           binDir,
		installedVersion: fixedVersion(""),
	}
	if _, err := inst.install(&output.Context{}); err == nil || !strings.Contains(err.Error(), binaryName) {
		t.Fatalf("install() error = %v, want error naming %q", err, binaryName)
	}
	if got := dirEntries(t, binDir); len(got) != 0 {
		t.Errorf("bin directory = %v, want empty", got)
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{name: "release output", output: "acloud version 1.2.0 (abc123, 2026-11-18)\n", want: "1.2.0"},
		{name: "only first line counts", output: "acloud version 1.2.0\nacloud version 9.9.9\n", want: "1.2.0"},
		{name: "unexpected output", output: "something else\n", want: ""},
		{name: "empty output", output: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseVersion(tt.output); got != tt.want {
				t.Errorf("parseVersion(%q) = %q, want %q", tt.output, got, tt.want)
			}
		})
	}
}

func TestCompletionVersion(t *testing.T) {
	if got := completionVersion("1.2.0"); got != "1_2_0" {
		t.Errorf("completionVersion(%q) = %q, want %q", "1.2.0", got, "1_2_0")
	}
}
