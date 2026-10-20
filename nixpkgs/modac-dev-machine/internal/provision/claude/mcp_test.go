package claude

import (
	"reflect"
	"testing"
)

func grafanaServer(args ...string) map[string]any {
	list := make([]any, len(args))
	for i, a := range args {
		list[i] = a
	}
	return map[string]any{"command": "uvx", "args": list}
}

func TestRemoveObsoleteMcpArgsStripsDisableSiftFromGrafana(t *testing.T) {
	servers := map[string]any{
		"grafana": grafanaServer("mcp-grafana", "--disable-write", "--disable-sift", "--disable-admin"),
	}

	fixed := removeObsoleteMcpArgs(servers)

	if fixed != 1 {
		t.Errorf("fixed = %d, want 1", fixed)
	}
	want := grafanaServer("mcp-grafana", "--disable-write", "--disable-admin")
	if !reflect.DeepEqual(servers["grafana"], want) {
		t.Errorf("grafana = %#v, want %#v", servers["grafana"], want)
	}
}

func TestRemoveObsoleteMcpArgsLeavesCurrentConfigAlone(t *testing.T) {
	servers := map[string]any{
		"grafana": grafanaServer("mcp-grafana", "--disable-write"),
		"other":   grafanaServer("something", "--disable-sift"),
	}
	snapshot := map[string]any{
		"grafana": grafanaServer("mcp-grafana", "--disable-write"),
		"other":   grafanaServer("something", "--disable-sift"),
	}

	if fixed := removeObsoleteMcpArgs(servers); fixed != 0 {
		t.Errorf("fixed = %d, want 0", fixed)
	}
	if !reflect.DeepEqual(servers, snapshot) {
		t.Errorf("servers changed: %#v", servers)
	}
}

func TestMergeMcpServersAddsOnlyMissingServers(t *testing.T) {
	custom := grafanaServer("mcp-grafana", "--my-own-flag")
	servers := map[string]any{"grafana": custom}
	tmpl := map[string]any{
		"grafana": grafanaServer("mcp-grafana"),
		"new":     grafanaServer("new-server"),
	}

	if added := mergeMcpServers(servers, tmpl); added != 1 {
		t.Errorf("added = %d, want 1", added)
	}
	if !reflect.DeepEqual(servers["grafana"], custom) {
		t.Errorf("existing grafana entry was overwritten: %#v", servers["grafana"])
	}
	if _, ok := servers["new"]; !ok {
		t.Errorf("missing server was not added")
	}
}
