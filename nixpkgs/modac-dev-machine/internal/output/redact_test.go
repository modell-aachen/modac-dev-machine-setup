package output

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	classicToken     = "ghp_FRqsUNHbJCWa0123456789abcdefghijKLMN"
	oauthToken       = "gho_0123456789abcdefghijKLMNOPQRSTUVwxyz"
	fineGrainedToken = "github_pat_11ABCDEFG0123456789_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123"
)

func TestRedactSecrets(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no secret", "copying path '/nix/store/abc-gh-2.105.0'\n", "copying path '/nix/store/abc-gh-2.105.0'\n"},
		{"classic token in nix access-tokens option", "--option access-tokens 'github.com=" + classicToken + "' profile install\n", "--option access-tokens 'github.com=ghp_<REDACTED>' profile install\n"},
		{"oauth token", "token " + oauthToken + "\n", "token gho_<REDACTED>\n"},
		{"fine-grained token", "GITHUB_TOKEN=" + fineGrainedToken + "\n", "GITHUB_TOKEN=github_pat_<REDACTED>\n"},
		{"several tokens on one line", classicToken + " " + oauthToken, "ghp_<REDACTED> gho_<REDACTED>"},
		{"already censored by devbox", "github.com=ghp_<CENSORED>", "github.com=ghp_<CENSORED>"},
		{"short ghp_ prefix is no token", "ghp_short", "ghp_short"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(redactSecrets([]byte(tt.in))); got != tt.want {
				t.Errorf("redactSecrets(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRedactingWriterTokenSplitAcrossWrites(t *testing.T) {
	var buf bytes.Buffer
	w := &redactingWriter{dst: &buf}

	line := "access-tokens 'github.com=" + classicToken + "'\n"
	split := strings.Index(line, classicToken) + 10
	for _, chunk := range []string{line[:split], line[split:]} {
		if _, err := w.Write([]byte(chunk)); err != nil {
			t.Fatalf("Write: %v", err)
		}
	}

	if want := "access-tokens 'github.com=ghp_<REDACTED>'\n"; buf.String() != want {
		t.Errorf("log = %q, want %q", buf.String(), want)
	}
}

func TestRedactingWriterFlushWritesPendingLine(t *testing.T) {
	var buf bytes.Buffer
	w := &redactingWriter{dst: &buf}

	if _, err := w.Write([]byte("no trailing newline " + classicToken)); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("incomplete line written before flush: %q", buf.String())
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	if want := "no trailing newline ghp_<REDACTED>"; buf.String() != want {
		t.Errorf("log = %q, want %q", buf.String(), want)
	}
}

func TestRedactingWriterFlushesOversizedLine(t *testing.T) {
	var buf bytes.Buffer
	w := &redactingWriter{dst: &buf}

	if _, err := w.Write(bytes.Repeat([]byte("x"), maxPendingBytes+1)); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if buf.Len() != maxPendingBytes+1 {
		t.Errorf("wrote %d bytes, want %d", buf.Len(), maxPendingBytes+1)
	}
}

func TestRedactingWriterTokenAtForcedFlush(t *testing.T) {
	filler := strings.Repeat("x", maxPendingBytes-10) + " "
	tests := []struct {
		name   string
		writes []string
	}{
		{"token split across writes at the limit", []string{filler, classicToken[:15], classicToken[15:] + "\n"}},
		// 236 trailing spaces put the token across a cut 256 bytes before the end.
		{"token in the middle of one oversized write", []string{filler + classicToken + strings.Repeat(" ", 236) + "\n"}},
		{"token followed by more text without newline", []string{filler + classicToken + " more", " text\n"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := &redactingWriter{dst: &buf}
			for _, chunk := range tt.writes {
				if _, err := w.Write([]byte(chunk)); err != nil {
					t.Fatalf("Write: %v", err)
				}
			}
			if err := w.Flush(); err != nil {
				t.Fatalf("Flush: %v", err)
			}

			if strings.Contains(buf.String(), classicToken[15:]) {
				t.Error("log contains the tail of a token split at the forced flush")
			}
			if !strings.Contains(buf.String(), " ghp_<REDACTED>") {
				t.Error("log misses the redacted token")
			}
		})
	}
}

func TestLogFileContainsNoToken(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LEAKED_TOKEN", classicToken)

	out, err := New()
	if err != nil {
		t.Fatalf("New(): %v", err)
	}

	if err := out.RunCommand("sh", "-c", "printf 'stdout %s\\n' \"$LEAKED_TOKEN\"; printf 'stderr %s' \"$LEAKED_TOKEN\" >&2; exit 1"); err == nil {
		t.Fatal("RunCommand() error = nil, want exit status 1")
	}
	if err := out.RunCommand("echo", classicToken); err != nil {
		t.Fatalf("RunCommand(echo): %v", err)
	}
	if _, err := out.MultiWriter(nil).Write([]byte(oauthToken + "\n")); err != nil {
		t.Fatalf("MultiWriter write: %v", err)
	}
	if err := out.Close(); err != nil {
		t.Fatalf("Close(): %v", err)
	}

	log, err := os.ReadFile(out.LogPath())
	if err != nil {
		t.Fatalf("reading log: %v", err)
	}
	for _, token := range []string{classicToken, oauthToken} {
		if bytes.Contains(log, []byte(token)) {
			t.Errorf("log contains token %s…:\n%s", token[:8], log)
		}
	}
	for _, want := range []string{"stdout ghp_<REDACTED>", "stderr ghp_<REDACTED>", "$ " + mustLookPath(t, "echo") + " ghp_<REDACTED>", "gho_<REDACTED>"} {
		if !bytes.Contains(log, []byte(want)) {
			t.Errorf("log misses %q:\n%s", want, log)
		}
	}
}

func mustLookPath(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Fatalf("LookPath(%s): %v", name, err)
	}
	return path
}
