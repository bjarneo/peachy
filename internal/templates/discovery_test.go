package templates

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateManifest(t *testing.T) {
	tests := []struct {
		name      string
		manifest  TemplateManifest
		wantError bool
	}{
		{
			name: "valid manifest",
			manifest: TemplateManifest{
				Name: "Test",
				Files: []TemplateFile{
					{Template: "test.conf", Destination: "~/.config/test"},
				},
			},
			wantError: false,
		},
		{
			name: "missing name",
			manifest: TemplateManifest{
				Files: []TemplateFile{
					{Template: "test.conf", Destination: "~/.config/test"},
				},
			},
			wantError: true,
		},
		{
			name: "no files",
			manifest: TemplateManifest{
				Name:  "Test",
				Files: []TemplateFile{},
			},
			wantError: true,
		},
		{
			name: "missing template in file",
			manifest: TemplateManifest{
				Name: "Test",
				Files: []TemplateFile{
					{Destination: "~/.config/test"},
				},
			},
			wantError: true,
		},
		{
			name: "missing destination in file",
			manifest: TemplateManifest{
				Name: "Test",
				Files: []TemplateFile{
					{Template: "test.conf"},
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateManifest(&tt.manifest)
			if (err != nil) != tt.wantError {
				t.Errorf("validateManifest() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		template *Template
		expected bool
	}{
		{
			name: "enabled by default",
			template: &Template{
				Manifest: TemplateManifest{Disabled: false},
			},
			expected: true,
		},
		{
			name: "explicitly disabled",
			template: &Template{
				Manifest: TemplateManifest{Disabled: true},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEnabled(tt.template)
			if result != tt.expected {
				t.Errorf("IsEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckCondition(t *testing.T) {
	tests := []struct {
		name     string
		template *Template
		expectOk bool
	}{
		{
			name: "no condition",
			template: &Template{
				Manifest: TemplateManifest{Condition: ""},
			},
			expectOk: true,
		},
		{
			name: "condition met (ls exists)",
			template: &Template{
				Manifest: TemplateManifest{Condition: "ls"},
			},
			expectOk: true,
		},
		{
			name: "condition not met",
			template: &Template{
				Manifest: TemplateManifest{Condition: "nonexistent-command-12345"},
			},
			expectOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, _ := CheckCondition(tt.template)
			if ok != tt.expectOk {
				t.Errorf("CheckCondition() = %v, want %v", ok, tt.expectOk)
			}
		})
	}
}

func TestLoadTemplate(t *testing.T) {
	// Create a temporary template directory
	tmpDir := t.TempDir()
	templateDir := filepath.Join(tmpDir, "testapp")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create manifest
	manifestContent := `name = "Test App"
description = "A test template"
version = "1.0"

[[files]]
template = "colors.conf"
destination = "~/.config/test/colors.conf"
`
	manifestPath := filepath.Join(templateDir, "template.toml")
	if err := os.WriteFile(manifestPath, []byte(manifestContent), 0644); err != nil {
		t.Fatalf("Failed to write manifest: %v", err)
	}

	// Load template
	template, err := LoadTemplate(templateDir)
	if err != nil {
		t.Fatalf("LoadTemplate() error = %v", err)
	}

	// Verify
	if template.Name != "testapp" {
		t.Errorf("template.Name = %q, want %q", template.Name, "testapp")
	}
	if template.Manifest.Name != "Test App" {
		t.Errorf("template.Manifest.Name = %q, want %q", template.Manifest.Name, "Test App")
	}
	if template.Manifest.Description != "A test template" {
		t.Errorf("template.Manifest.Description = %q, want %q", template.Manifest.Description, "A test template")
	}
	if len(template.Manifest.Files) != 1 {
		t.Errorf("len(template.Manifest.Files) = %d, want 1", len(template.Manifest.Files))
	}
}

func TestGetPostApplyPath(t *testing.T) {
	// Create a temporary template directory
	tmpDir := t.TempDir()
	templateDir := filepath.Join(tmpDir, "testapp")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	template := &Template{
		Path: templateDir,
	}

	// No post-apply script
	path := GetPostApplyPath(template)
	if path != "" {
		t.Errorf("GetPostApplyPath() = %q, want empty string", path)
	}

	// Create post-apply script
	postApplyPath := filepath.Join(templateDir, "post-apply")
	if err := os.WriteFile(postApplyPath, []byte("#!/bin/bash\nexit 0"), 0755); err != nil {
		t.Fatalf("Failed to write post-apply: %v", err)
	}

	// Now should find it
	path = GetPostApplyPath(template)
	if path != postApplyPath {
		t.Errorf("GetPostApplyPath() = %q, want %q", path, postApplyPath)
	}
}
