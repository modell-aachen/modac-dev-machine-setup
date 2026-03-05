package claude

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/util"
)

// Run sets up Claude configuration
func Run(out *output.Context, plat platform.Platform) error {
	_ = plat
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeDir := filepath.Join(homeDir, ".claude")

	// Create .claude directory
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	if err := setupClaudeMd(out, claudeDir); err != nil {
		return err
	}

	if err := setupMcpConfig(out, claudeDir); err != nil {
		return err
	}

	return nil
}

func setupClaudeMd(out *output.Context, claudeDir string) error {
	claudeMdPath := filepath.Join(claudeDir, "CLAUDE.md")

	if util.FileExists(claudeMdPath) {
		out.Skipped("CLAUDE.md already exists")
		return nil
	}

	templatesDir, err := util.GetTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to find templates directory: %w", err)
	}

	out.Step("Creating CLAUDE.md from template")
	templatePath := filepath.Join(templatesDir, "team-claude.md")
	if err := copyFile(templatePath, claudeMdPath); err != nil {
		return fmt.Errorf("failed to copy template: %w", err)
	}

	if err := os.Chmod(claudeMdPath, 0664); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	return nil
}

func setupMcpConfig(out *output.Context, claudeDir string) error {
	templatesDir, err := util.GetTemplatesDir()
	if err != nil {
		return fmt.Errorf("failed to find templates directory: %w", err)
	}

	templatePath := filepath.Join(templatesDir, "claude-mcp.json")
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read MCP template: %w", err)
	}

	var templateConfig map[string]any
	if err := json.Unmarshal(templateData, &templateConfig); err != nil {
		return fmt.Errorf("failed to parse MCP template: %w", err)
	}

	templateServers, _ := templateConfig["mcpServers"].(map[string]any)
	if len(templateServers) == 0 {
		out.Skipped("No MCP servers in template")
		return nil
	}

	mcpPath := filepath.Join(claudeDir, ".mcp.json")

	var existing map[string]any
	if util.FileExists(mcpPath) {
		data, err := os.ReadFile(mcpPath)
		if err != nil {
			return fmt.Errorf("failed to read existing MCP config: %w", err)
		}
		if err := json.Unmarshal(data, &existing); err != nil {
			return fmt.Errorf("failed to parse existing MCP config: %w", err)
		}
	} else {
		existing = map[string]any{}
	}

	existingServers, ok := existing["mcpServers"].(map[string]any)
	if !ok {
		existingServers = map[string]any{}
	}

	added := 0
	for name, config := range templateServers {
		if _, exists := existingServers[name]; !exists {
			existingServers[name] = config
			added++
		}
	}

	if added == 0 {
		out.Skipped("MCP servers already configured")
		return nil
	}

	existing["mcpServers"] = existingServers

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal MCP config: %w", err)
	}

	out.Step(fmt.Sprintf("Adding %d MCP server(s) to .mcp.json", added))
	if err := os.WriteFile(mcpPath, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("failed to write MCP config: %w", err)
	}

	return nil
}


func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
