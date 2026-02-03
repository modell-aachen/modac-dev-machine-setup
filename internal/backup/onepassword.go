package backup

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type OpClient struct {
	vault string
}

type OpItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func NewOpClient(vault string) *OpClient {
	return &OpClient{vault: vault}
}

func (c *OpClient) CheckInstalled() error {
	if _, err := exec.LookPath("op"); err != nil {
		return fmt.Errorf("1Password CLI 'op' is required but not installed")
	}
	return nil
}

func (c *OpClient) CheckSignedIn() error {
	cmd := exec.Command("op", "whoami")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("you must be signed in to 1Password CLI (run: eval \"$(op signin)\")")
	}
	return nil
}

func (c *OpClient) ListItems() ([]OpItem, error) {
	cmd := exec.Command("op", "item", "list", "--vault", c.vault, "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("op item list failed: %w", err)
	}

	var items []OpItem
	if err := json.Unmarshal(output, &items); err != nil {
		return nil, fmt.Errorf("failed to parse op item list: %w", err)
	}

	return items, nil
}

func (c *OpClient) FindItemByTitle(title string) (string, error) {
	items, err := c.ListItems()
	if err != nil {
		return "", err
	}

	for _, item := range items {
		if item.Title == title {
			return item.ID, nil
		}
	}

	return "", nil
}

func (c *OpClient) DeleteItem(itemID string) error {
	cmd := exec.Command("op", "item", "delete", itemID, "--vault", c.vault)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("op item delete failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (c *OpClient) CreateDocument(filePath, title string) error {
	cmd := exec.Command("op", "document", "create", filePath, "--title", title, "--vault", c.vault)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("op document create failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (c *OpClient) GetDocument(itemID, outPath string) error {
	cmd := exec.Command("op", "document", "get", itemID, "--vault", c.vault, "--out-file", outPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if error is because item is not a document
		if strings.Contains(string(output), "not a document") {
			// Try alternative method: get item and save to file
			return c.getDocumentAlternative(itemID, outPath)
		}
		return fmt.Errorf("op document get failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (c *OpClient) getDocumentAlternative(itemID, outPath string) error {
	cmd := exec.Command("op", "document", "get", itemID, "--vault", c.vault)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("op document get (alternative) failed: %w", err)
	}

	return writeFile(outPath, output)
}
