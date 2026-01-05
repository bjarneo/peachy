package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"peachy/internal/color"
)

const (
	// GeneratedTemplatesDir is where processed templates are stored before symlinking
	GeneratedTemplatesDir = "generated/templates"
)

// GetGeneratedTemplatesDir returns the directory for processed template files
func GetGeneratedTemplatesDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "peachy", GeneratedTemplatesDir)
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}

// ProcessAllTemplates discovers and processes all enabled templates
func ProcessAllTemplates(p *color.Palette) ([]*ProcessResult, error) {
	templates, err := DiscoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	if len(templates) == 0 {
		return nil, nil
	}

	variables := BuildVariables(p)
	var results []*ProcessResult

	for _, t := range templates {
		result := ProcessTemplate(t, variables)
		results = append(results, result)
	}

	return results, nil
}

// ProcessTemplate processes a single template
func ProcessTemplate(t *Template, variables map[string]string) *ProcessResult {
	result := &ProcessResult{
		Template: t,
	}

	// Check if disabled
	if !IsEnabled(t) {
		result.Skipped = true
		result.SkipReason = "disabled in manifest"
		return result
	}

	// Check condition
	if ok, reason := CheckCondition(t); !ok {
		result.Skipped = true
		result.SkipReason = reason
		return result
	}

	// Process each file in the manifest
	for _, file := range t.Manifest.Files {
		linkedPath, err := processTemplateFile(t, file, variables)
		if err != nil {
			result.Error = fmt.Errorf("failed to process %s: %w", file.Template, err)
			return result
		}
		result.FilesLinked = append(result.FilesLinked, linkedPath)
	}

	// Run post-apply script if it exists
	if err := runPostApply(t); err != nil {
		result.Error = fmt.Errorf("post-apply failed: %w", err)
		return result
	}

	result.Success = true
	return result
}

// processTemplateFile processes a single template file and creates symlink
func processTemplateFile(t *Template, file TemplateFile, variables map[string]string) (string, error) {
	// Read template file
	templatePath := filepath.Join(t.Path, file.Template)
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Process variables
	processed := ProcessContent(string(content), variables)

	// Determine output path in generated directory
	generatedDir := GetGeneratedTemplatesDir()
	outputDir := filepath.Join(generatedDir, t.Name)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	outputPath := filepath.Join(outputDir, file.Template)
	if err := os.WriteFile(outputPath, []byte(processed), 0644); err != nil {
		return "", fmt.Errorf("failed to write processed file: %w", err)
	}

	// Create symlink to destination
	destPath := expandPath(file.Destination)
	if err := createSymlink(outputPath, destPath); err != nil {
		return "", fmt.Errorf("failed to create symlink: %w", err)
	}

	return destPath, nil
}

// createSymlink creates a symlink from source to destination
func createSymlink(source, dest string) error {
	// Ensure parent directory exists
	parentDir := filepath.Dir(dest)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Remove existing file/symlink at destination
	if info, err := os.Lstat(dest); err == nil {
		if info.Mode()&os.ModeSymlink != 0 || info.Mode().IsRegular() {
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("failed to remove existing file: %w", err)
			}
		}
	}

	// Create symlink
	if err := os.Symlink(source, dest); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// runPostApply executes the post-apply script if it exists
func runPostApply(t *Template) error {
	scriptPath := GetPostApplyPath(t)
	if scriptPath == "" {
		return nil // No post-apply script
	}

	// Check if executable
	info, err := os.Stat(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to stat post-apply script: %w", err)
	}

	if info.Mode()&0111 == 0 {
		return fmt.Errorf("post-apply script is not executable")
	}

	// Execute script
	cmd := exec.Command(scriptPath)
	cmd.Dir = t.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script failed: %w\nOutput: %s", err, strings.TrimSpace(string(output)))
	}

	return nil
}

// ValidateTemplates validates all discovered templates without processing
func ValidateTemplates() ([]*ValidationResult, error) {
	templates, err := DiscoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	var results []*ValidationResult

	for _, t := range templates {
		result := validateTemplate(t)
		results = append(results, result)
	}

	return results, nil
}

// ValidationResult represents the result of validating a template
type ValidationResult struct {
	Template *Template
	Valid    bool
	Errors   []string
	Warnings []string
}

// validateTemplate validates a single template
func validateTemplate(t *Template) *ValidationResult {
	result := &ValidationResult{
		Template: t,
		Valid:    true,
	}

	// Check all template files exist
	for _, file := range t.Manifest.Files {
		templatePath := filepath.Join(t.Path, file.Template)
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			result.Errors = append(result.Errors, fmt.Sprintf("template file not found: %s", file.Template))
			result.Valid = false
		}

		// Check destination parent directory exists (warning only)
		destPath := expandPath(file.Destination)
		parentDir := filepath.Dir(destPath)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("destination directory does not exist: %s", parentDir))
		}
	}

	// Check post-apply script if present
	scriptPath := GetPostApplyPath(t)
	if scriptPath != "" {
		info, err := os.Stat(scriptPath)
		if err == nil && info.Mode()&0111 == 0 {
			result.Warnings = append(result.Warnings, "post-apply script is not executable")
		}
	}

	// Check condition
	if t.Manifest.Condition != "" {
		if _, err := exec.LookPath(t.Manifest.Condition); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("condition command not found: %s", t.Manifest.Condition))
		}
	}

	return result
}

// ListTemplates returns all discovered templates with their status
func ListTemplates() ([]*TemplateStatus, error) {
	templates, err := DiscoverTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to discover templates: %w", err)
	}

	var statuses []*TemplateStatus

	for _, t := range templates {
		status := &TemplateStatus{
			Template:       t,
			Enabled:        IsEnabled(t),
			ConditionMet:   true,
			ConditionError: "",
		}

		if ok, reason := CheckCondition(t); !ok {
			status.ConditionMet = false
			status.ConditionError = reason
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// TemplateStatus represents the status of a template
type TemplateStatus struct {
	Template       *Template
	Enabled        bool
	ConditionMet   bool
	ConditionError string
}
