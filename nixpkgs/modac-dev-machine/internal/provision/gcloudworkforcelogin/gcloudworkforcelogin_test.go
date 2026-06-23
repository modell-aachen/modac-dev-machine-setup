package gcloudworkforcelogin

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildCredConfig(t *testing.T) {
	const execPath = "/home/dev/.local/bin/ariadne-gcp-workforce-login"

	data, err := buildCredConfig(execPath)
	if err != nil {
		t.Fatalf("buildCredConfig returned error: %v", err)
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("credential config is not valid JSON: %v", err)
	}

	if cfg["type"] != "external_account" {
		t.Errorf("type = %v, want external_account", cfg["type"])
	}
	if cfg["subject_token_type"] != "urn:ietf:params:oauth:token-type:id_token" {
		t.Errorf("subject_token_type = %v", cfg["subject_token_type"])
	}
	if cfg["workforce_pool_user_project"] != quotaProject {
		t.Errorf("workforce_pool_user_project = %v, want %s", cfg["workforce_pool_user_project"], quotaProject)
	}

	audience, _ := cfg["audience"].(string)
	if !strings.Contains(audience, workforcePool) || !strings.Contains(audience, workforceProvider) {
		t.Errorf("audience %q does not reference pool/provider", audience)
	}

	source, ok := cfg["credential_source"].(map[string]any)
	if !ok {
		t.Fatalf("credential_source missing or wrong type")
	}
	executable, ok := source["executable"].(map[string]any)
	if !ok {
		t.Fatalf("credential_source.executable missing or wrong type")
	}
	if executable["command"] != execPath {
		t.Errorf("command = %v, want %s", executable["command"], execPath)
	}
}

func TestEmbeddedExecutable(t *testing.T) {
	if !bytes.HasPrefix(execScript, []byte("#!")) {
		t.Error("embedded executable is missing a shebang")
	}
	for _, want := range []string{
		"--login",
		"--force-refresh",
		"--logout",
		"refresh_token",
		"urn:ietf:params:oauth:token-type:id_token",
	} {
		if !bytes.Contains(execScript, []byte(want)) {
			t.Errorf("embedded executable does not contain %q", want)
		}
	}
}
