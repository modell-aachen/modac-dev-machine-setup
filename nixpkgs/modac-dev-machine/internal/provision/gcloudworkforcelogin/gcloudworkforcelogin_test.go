package gcloudworkforcelogin

import (
	"bytes"
	"testing"
)

func TestEmbeddedExecutable(t *testing.T) {
	if !bytes.HasPrefix(execScript, []byte("#!")) {
		t.Error("embedded executable is missing a shebang")
	}
	for _, want := range []string{
		"get-token",
		"ExecCredential",
		"--login",
		"--has-session",
		"--check",
		"--force-refresh",
		"--logout",
		"refresh_token",
	} {
		if !bytes.Contains(execScript, []byte(want)) {
			t.Errorf("embedded executable does not contain %q", want)
		}
	}
}
