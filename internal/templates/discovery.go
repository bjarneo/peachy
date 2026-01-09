package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

const (
	TemplatesDir  = "templates"
	ManifestFile  = "template.toml"
	PostApplyFile = "post-apply"
)

// getConfigDir returns the peachy config directory (XDG compliant)
func getConfigDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "peachy")
}

// GetTemplatesDir returns the user's custom templates directory
func GetTemplatesDir() string {
	return filepath.Join(getConfigDir(), TemplatesDir)
}

// EnsureTemplatesDir creates the templates directory if it doesn't exist
func EnsureTemplatesDir() error {
	dir := GetTemplatesDir()
	return os.MkdirAll(dir, 0755)
}

// DiscoverTemplates finds all valid templates in the templates directory
func DiscoverTemplates() ([]*Template, error) {
	dir := GetTemplatesDir()

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil // No templates directory, return empty
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []*Template

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templatePath := filepath.Join(dir, entry.Name())
		manifestPath := filepath.Join(templatePath, ManifestFile)

		// Check if manifest exists
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue // Skip directories without manifest
		}

		template, err := LoadTemplate(templatePath)
		if err != nil {
			// Log error but continue with other templates
			continue
		}

		templates = append(templates, template)
	}

	return templates, nil
}

// LoadTemplate loads a template from a directory
func LoadTemplate(path string) (*Template, error) {
	manifestPath := filepath.Join(path, ManifestFile)

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest TemplateManifest
	if err := toml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Validate manifest
	if err := validateManifest(&manifest); err != nil {
		return nil, err
	}

	return &Template{
		Manifest: manifest,
		Path:     path,
		Name:     filepath.Base(path),
	}, nil
}

// validateManifest checks that required fields are present
func validateManifest(m *TemplateManifest) error {
	if m.Name == "" {
		return fmt.Errorf("manifest missing required field: name")
	}
	if len(m.Files) == 0 {
		return fmt.Errorf("manifest must define at least one file")
	}
	for i, f := range m.Files {
		if f.Template == "" {
			return fmt.Errorf("file %d missing required field: template", i)
		}
		if f.Destination == "" {
			return fmt.Errorf("file %d missing required field: destination", i)
		}
	}
	return nil
}

// CheckCondition verifies if a template's condition is met
func CheckCondition(t *Template) (bool, string) {
	if t.Manifest.Condition == "" {
		return true, ""
	}

	_, err := exec.LookPath(t.Manifest.Condition)
	if err != nil {
		return false, fmt.Sprintf("command '%s' not found in PATH", t.Manifest.Condition)
	}

	return true, ""
}

// IsEnabled returns whether a template is enabled
func IsEnabled(t *Template) bool {
	return !t.Manifest.Disabled
}

// GetPostApplyPath returns the path to the post-apply script if it exists
func GetPostApplyPath(t *Template) string {
	path := filepath.Join(t.Path, PostApplyFile)
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		return path
	}
	return ""
}
