package claude

import (
	"reflect"
	"testing"
)

func template() map[string]any {
	return map[string]any{
		"extraKnownMarketplaces": map[string]any{
			"claude-skills": map[string]any{
				"source":     map[string]any{"source": "github", "repo": "modell-aachen/claude-skills"},
				"autoUpdate": true,
			},
		},
		"enabledPlugins": map[string]any{
			"dev-workflow@claude-skills": true,
		},
	}
}

func TestMergePluginSettings(t *testing.T) {
	tests := []struct {
		name      string
		settings  map[string]any
		wantAdded int
		wantKeep  map[string]bool // enabledPlugins keys that must remain enabled
	}{
		{
			name:      "empty settings gets marketplace and plugin",
			settings:  map[string]any{},
			wantAdded: 2,
			wantKeep:  map[string]bool{"dev-workflow@claude-skills": true},
		},
		{
			name: "preserves unrelated enabled plugins and marketplaces",
			settings: map[string]any{
				"enabledPlugins": map[string]any{
					"typescript-lsp@claude-plugins-official": true,
				},
			},
			wantAdded: 2,
			wantKeep: map[string]bool{
				"typescript-lsp@claude-plugins-official": true,
				"dev-workflow@claude-skills":             true,
			},
		},
		{
			name: "marketplace present but plugin missing adds only the plugin",
			settings: map[string]any{
				"extraKnownMarketplaces": map[string]any{
					"claude-skills": map[string]any{"source": map[string]any{}},
				},
			},
			wantAdded: 1,
			wantKeep:  map[string]bool{"dev-workflow@claude-skills": true},
		},
		{
			name: "fully configured is a no-op",
			settings: map[string]any{
				"extraKnownMarketplaces": map[string]any{
					"claude-skills": map[string]any{"source": map[string]any{}},
				},
				"enabledPlugins": map[string]any{
					"dev-workflow@claude-skills": true,
				},
			},
			wantAdded: 0,
			wantKeep:  map[string]bool{"dev-workflow@claude-skills": true},
		},
		{
			name: "does not overwrite a user-disabled plugin",
			settings: map[string]any{
				"enabledPlugins": map[string]any{
					"dev-workflow@claude-skills": false,
				},
			},
			wantAdded: 1, // only the marketplace is added
			wantKeep:  map[string]bool{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			added := mergePluginSettings(tt.settings, template())

			if added != tt.wantAdded {
				t.Errorf("added = %d, want %d", added, tt.wantAdded)
			}

			enabled, _ := tt.settings["enabledPlugins"].(map[string]any)
			for plugin, want := range tt.wantKeep {
				if got, _ := enabled[plugin].(bool); got != want {
					t.Errorf("enabledPlugins[%q] = %v, want %v", plugin, enabled[plugin], want)
				}
			}

			// The user override case must be left untouched.
			if tt.name == "does not overwrite a user-disabled plugin" {
				if got := enabled["dev-workflow@claude-skills"]; got != false {
					t.Errorf("user-disabled plugin was overwritten: got %v", got)
				}
			}

			// The marketplace must always be registered afterwards.
			markets, _ := tt.settings["extraKnownMarketplaces"].(map[string]any)
			if _, ok := markets["claude-skills"]; !ok {
				t.Errorf("claude-skills marketplace missing after merge")
			}
		})
	}
}

func TestMergePluginSettingsIsIdempotent(t *testing.T) {
	settings := map[string]any{}

	first := mergePluginSettings(settings, template())
	snapshot := deepCopy(settings)

	second := mergePluginSettings(settings, template())

	if first == 0 {
		t.Fatalf("first merge added nothing")
	}
	if second != 0 {
		t.Errorf("second merge added %d entries, want 0", second)
	}
	if !reflect.DeepEqual(settings, snapshot) {
		t.Errorf("settings changed on the second merge:\n got %#v\nwant %#v", settings, snapshot)
	}
}

func deepCopy(m map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			out[k] = deepCopy(nested)
			continue
		}
		out[k] = v
	}
	return out
}
