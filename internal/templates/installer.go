package templates

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

const (
	RepoBaseURL = "https://raw.githubusercontent.com/bjarneo/peachy/main/examples/templates"
)

// AvailableTemplate represents a template that can be installed
type AvailableTemplate struct {
	Name        string
	Description string
	Files       []string // Files to download (template.toml is always included)
	LinuxOnly   bool
	MacOSOnly   bool
}

// AvailableTemplates lists all templates that can be installed from the repository
var AvailableTemplates = []AvailableTemplate{
	{Name: "alacritty", Description: "GPU-accelerated terminal emulator", Files: []string{"colors.toml", "post-apply"}},
	{Name: "btop", Description: "Resource monitor with colorful interface", Files: []string{"peachy.theme"}},
	{Name: "cava", Description: "Console audio visualizer", Files: []string{"theme.ini", "post-apply"}},
	{Name: "dunst", Description: "Lightweight notification daemon", Files: []string{"colors.conf", "post-apply"}},
	{Name: "foot", Description: "Fast Wayland terminal emulator", Files: []string{"colors.ini"}, LinuxOnly: true},
	{Name: "ghostty", Description: "Zig-based GPU-accelerated terminal", Files: []string{"colors.conf"}},
	{Name: "gtk", Description: "GTK3/GTK4 application theming", Files: []string{"colors.css"}},
	{Name: "hyprland", Description: "Dynamic tiling Wayland compositor", Files: []string{"colors.conf", "post-apply"}, LinuxOnly: true},
	{Name: "hyprlock", Description: "Screen locker for Hyprland", Files: []string{"colors.conf"}, LinuxOnly: true},
	{Name: "iterm2", Description: "macOS terminal emulator", Files: []string{"peachy.itermcolors", "post-apply"}, MacOSOnly: true},
	{Name: "kitty", Description: "GPU-based terminal emulator", Files: []string{"colors.conf", "post-apply"}},
	{Name: "mako", Description: "Wayland notification daemon", Files: []string{"config", "post-apply"}, LinuxOnly: true},
	{Name: "neovim", Description: "Neovim colorscheme (aether.nvim)", Files: []string{"aether.lua"}},
	{Name: "rofi", Description: "Application launcher", Files: []string{"colors.rasi"}},
	{Name: "swayosd", Description: "On-screen display for Sway/Hyprland", Files: []string{"style.css"}, LinuxOnly: true},
	{Name: "vencord", Description: "Discord client mod theme", Files: []string{"peachy.theme.css"}},
	{Name: "vscode", Description: "VS Code live theme via settings.json", Files: []string{"package.json", "peachy-color-theme.json"}},
	{Name: "walker", Description: "Wayland application launcher", Files: []string{"colors.css"}, LinuxOnly: true},
	{Name: "warp", Description: "Modern terminal with AI features", Files: []string{"peachy.yaml"}},
	{Name: "waybar", Description: "Wayland status bar", Files: []string{"colors.css", "post-apply"}, LinuxOnly: true},
	{Name: "wofi", Description: "Wayland application launcher", Files: []string{"colors.css"}, LinuxOnly: true},
	{Name: "zed", Description: "Zed code editor theme", Files: []string{"peachy.json"}},
}

// GetAvailableTemplates returns templates available for the current platform
func GetAvailableTemplates() []AvailableTemplate {
	var result []AvailableTemplate
	isLinux := runtime.GOOS == "linux"
	isMacOS := runtime.GOOS == "darwin"

	for _, t := range AvailableTemplates {
		// Skip platform-specific templates
		if t.LinuxOnly && !isLinux {
			continue
		}
		if t.MacOSOnly && !isMacOS {
			continue
		}
		result = append(result, t)
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// GetTemplateByName finds a template by name
func GetTemplateByName(name string) *AvailableTemplate {
	for _, t := range AvailableTemplates {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// IsTemplateInstalled checks if a template is already installed
func IsTemplateInstalled(name string) bool {
	path := filepath.Join(GetTemplatesDir(), name)
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// InstallTemplate downloads and installs a template from the repository
func InstallTemplate(name string) error {
	t := GetTemplateByName(name)
	if t == nil {
		return fmt.Errorf("unknown template: %s", name)
	}

	// Check platform compatibility
	isLinux := runtime.GOOS == "linux"
	isMacOS := runtime.GOOS == "darwin"

	if t.LinuxOnly && !isLinux {
		return fmt.Errorf("template '%s' is only available on Linux", name)
	}
	if t.MacOSOnly && !isMacOS {
		return fmt.Errorf("template '%s' is only available on macOS", name)
	}

	// Ensure templates directory exists
	if err := EnsureTemplatesDir(); err != nil {
		return fmt.Errorf("creating templates directory: %w", err)
	}

	destDir := filepath.Join(GetTemplatesDir(), name)

	// Create template directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("creating template directory: %w", err)
	}

	// Always download template.toml
	files := append([]string{"template.toml"}, t.Files...)

	for _, file := range files {
		url := fmt.Sprintf("%s/%s/%s", RepoBaseURL, name, file)
		destPath := filepath.Join(destDir, file)

		if err := downloadFile(url, destPath); err != nil {
			// post-apply is optional, don't fail if it doesn't exist
			if file == "post-apply" {
				continue
			}
			// Clean up on failure
			_ = os.RemoveAll(destDir)
			return fmt.Errorf("downloading %s: %w", file, err)
		}

		// Make post-apply executable
		if file == "post-apply" {
			_ = os.Chmod(destPath, 0755)
		}
	}

	return nil
}

// UninstallTemplate removes an installed template
func UninstallTemplate(name string) error {
	path := filepath.Join(GetTemplatesDir(), name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("template '%s' is not installed", name)
	}
	return os.RemoveAll(path)
}

// downloadFile downloads a file from a URL to the specified path
func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	return err
}

// InstalledTemplateInfo contains information about an installed template
type InstalledTemplateInfo struct {
	Name        string
	Description string
	Installed   bool
}

// ListTemplateStatus returns installation status for all available templates
func ListTemplateStatus() []InstalledTemplateInfo {
	available := GetAvailableTemplates()
	result := make([]InstalledTemplateInfo, len(available))

	for i, t := range available {
		result[i] = InstalledTemplateInfo{
			Name:        t.Name,
			Description: t.Description,
			Installed:   IsTemplateInstalled(t.Name),
		}
	}

	return result
}
