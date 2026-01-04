package config

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"peachy/internal/color"
)

//go:embed templates/*
var templatesFS embed.FS

const (
	OmarchyThemesDir = ".config/omarchy/themes"
	OmarchyThemeName = "peachy"
)

// GetPeachyThemeDir returns the path to peachy theme output directory
func GetPeachyThemeDir() string {
	return filepath.Join(GetConfigDir(), GeneratedDir)
}

// GetOmarchyThemeDir returns the path to omarchy themes directory
func GetOmarchyThemeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, OmarchyThemesDir, OmarchyThemeName)
}

// IsOmarchyInstalled checks if omarchy-theme-set exists
func IsOmarchyInstalled() bool {
	_, err := exec.LookPath("omarchy-theme-set")
	return err == nil
}

// GetReferenceThemeFiles returns the list of files in a reference theme directory
func GetReferenceThemeFiles(themePath string) ([]string, error) {
	// Follow symlinks to get the actual directory
	realPath, err := filepath.EvalSymlinks(themePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve theme path: %w", err)
	}

	entries, err := os.ReadDir(realPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// GetEmbeddedTemplateFiles returns the list of embedded template files
func GetEmbeddedTemplateFiles() ([]string, error) {
	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded templates: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

// BuildTemplateVariables creates the variable map from a palette
func BuildTemplateVariables(p *color.Palette) map[string]string {
	vars := make(map[string]string)

	// Primary colors
	vars["background"] = p.Background.Hex
	vars["foreground"] = p.Foreground.Hex

	// Semantic color names
	colorNames := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"bright_black", "bright_red", "bright_green", "bright_yellow",
		"bright_blue", "bright_magenta", "bright_cyan", "bright_white",
	}

	for i, name := range colorNames {
		if i < len(p.Colors) {
			vars[name] = p.Colors[i].Hex
		}
	}

	// Also add color0-15 aliases for backwards compatibility
	for i := 0; i < 16 && i < len(p.Colors); i++ {
		vars[fmt.Sprintf("color%d", i)] = p.Colors[i].Hex
	}

	return vars
}

// ProcessTemplateContent processes template content and replaces variables
func ProcessTemplateContent(content string, variables map[string]string) string {
	result := content

	// Replace {key} patterns
	for key, value := range variables {
		// Replace {key}
		result = strings.ReplaceAll(result, "{"+key+"}", value)

		// Replace {key.strip} (removes # from hex colors)
		strippedValue := strings.TrimPrefix(value, "#")
		result = strings.ReplaceAll(result, "{"+key+".strip}", strippedValue)

		// Replace {key.rgb} (converts hex to decimal RGB: r,g,b)
		if strings.HasPrefix(value, "#") {
			rgbValue := hexToDecimalRGB(value)
			result = strings.ReplaceAll(result, "{"+key+".rgb}", rgbValue)
		}

		// Replace {key.rgba:alpha} patterns
		rgbaRegex := regexp.MustCompile(`\{` + regexp.QuoteMeta(key) + `\.rgba(?::(\d*\.?\d+))?\}`)
		result = rgbaRegex.ReplaceAllStringFunc(result, func(match string) string {
			alphaMatch := rgbaRegex.FindStringSubmatch(match)
			alpha := "1.0"
			if len(alphaMatch) > 1 && alphaMatch[1] != "" {
				alpha = alphaMatch[1]
			}
			return hexToRgbaString(value, alpha)
		})

		// Replace {key.yaru} (maps color to Yaru icon theme variant)
		if strings.HasPrefix(value, "#") {
			yaruValue := hexToYaruIconTheme(value)
			result = strings.ReplaceAll(result, "{"+key+".yaru}", yaruValue)
		}
	}

	return result
}

// hexToDecimalRGB converts #RRGGBB to "R,G,B" decimal format
func hexToDecimalRGB(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return hex
	}

	var r, g, b int
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d,%d,%d", r, g, b)
}

// hexToRgbaString converts #RRGGBB to "rgba(R,G,B,A)" CSS format
func hexToRgbaString(hex string, alpha string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return hex
	}

	var r, g, b int
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("rgba(%d,%d,%d,%s)", r, g, b, alpha)
}

// hexToYaruIconTheme maps a hex color to a Yaru icon theme variant based on hue
func hexToYaruIconTheme(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "Yaru"
	}

	var r, g, b int
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

	// Convert to HSL to get hue
	c, err := color.NewColorFromHex(hex)
	if err != nil {
		return "Yaru"
	}
	hue := c.HSL.H

	// Map hue ranges to Yaru icon theme variants
	switch {
	case hue >= 345 || hue < 15:
		return "Yaru-red"
	case hue >= 15 && hue < 30:
		return "Yaru-wartybrown"
	case hue >= 30 && hue < 60:
		return "Yaru-yellow"
	case hue >= 60 && hue < 90:
		return "Yaru-olive"
	case hue >= 90 && hue < 165:
		return "Yaru-sage"
	case hue >= 165 && hue < 195:
		return "Yaru-prussiangreen"
	case hue >= 195 && hue < 255:
		return "Yaru-blue"
	case hue >= 255 && hue < 285:
		return "Yaru-purple"
	default:
		return "Yaru-magenta"
	}
}

// GenerateThemeFiles processes embedded templates and generates theme files
func GenerateThemeFiles(p *color.Palette, referenceFiles []string, wallpaperPath string) error {
	outputDir := GetPeachyThemeDir()

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create theme directory: %w", err)
	}

	// Create backgrounds subdirectory
	bgDir := filepath.Join(outputDir, "backgrounds")
	if err := os.MkdirAll(bgDir, 0o755); err != nil {
		return fmt.Errorf("failed to create backgrounds directory: %w", err)
	}

	// Copy wallpaper to backgrounds folder if provided
	if wallpaperPath != "" {
		if err := copyWallpaper(wallpaperPath, bgDir); err != nil {
			return fmt.Errorf("failed to copy wallpaper: %w", err)
		}
	}

	variables := BuildTemplateVariables(p)

	// Get list of embedded templates
	embeddedFiles, err := GetEmbeddedTemplateFiles()
	if err != nil {
		return fmt.Errorf("failed to get embedded templates: %w", err)
	}

	// Create a set of reference files for quick lookup
	refSet := make(map[string]bool)
	for _, f := range referenceFiles {
		refSet[f] = true
	}

	// Process each embedded template that matches a reference file
	for _, fileName := range embeddedFiles {
		// Skip if not in reference files (unless reference is empty, then process all)
		if len(referenceFiles) > 0 && !refSet[fileName] {
			continue
		}

		// Skip special files
		if fileName == "backgrounds" || fileName == "preview.png" {
			continue
		}

		// Read embedded template
		content, err := templatesFS.ReadFile("templates/" + fileName)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", fileName, err)
		}

		// Process the template
		processed := ProcessTemplateContent(string(content), variables)

		// Write output file
		outputPath := filepath.Join(outputDir, fileName)
		if err := os.WriteFile(outputPath, []byte(processed), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", fileName, err)
		}
	}

	return nil
}

// copyWallpaper copies the wallpaper file to the backgrounds directory
func copyWallpaper(src, destDir string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read wallpaper: %w", err)
	}

	// Get filename from source path
	filename := filepath.Base(src)
	destPath := filepath.Join(destDir, filename)

	// Write to destination
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write wallpaper: %w", err)
	}

	return nil
}

// CreateOmarchySymlink creates symlink from omarchy themes to peachy theme
func CreateOmarchySymlink() error {
	if !IsOmarchyInstalled() {
		return nil // Not an omarchy system, skip silently
	}

	peachyThemeDir := GetPeachyThemeDir()
	omarchyThemeDir := GetOmarchyThemeDir()

	// Ensure parent directory exists
	parentDir := filepath.Dir(omarchyThemeDir)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return fmt.Errorf("failed to create omarchy themes directory: %w", err)
	}

	// Remove existing symlink or directory
	if info, err := os.Lstat(omarchyThemeDir); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			_ = os.Remove(omarchyThemeDir)
		} else if info.IsDir() {
			_ = os.RemoveAll(omarchyThemeDir)
		}
	}

	// Create symlink
	if err := os.Symlink(peachyThemeDir, omarchyThemeDir); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// GetCacheDir returns the cache directory path
func GetCacheDir() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, AppName)
}

// RunOmarchyThemeSet runs omarchy-theme-set to apply the theme
func RunOmarchyThemeSet() error {
	if !IsOmarchyInstalled() {
		return nil // Not an omarchy system, skip silently
	}

	// Ensure cache directory exists for log file
	cacheDir := GetCacheDir()
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Open log file
	logPath := filepath.Join(cacheDir, "apply.log")
	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer func() { _ = logFile.Close() }()

	cmd := exec.Command("omarchy-theme-set", OmarchyThemeName)
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run omarchy-theme-set (see %s): %w", logPath, err)
	}

	return nil
}

// ApplyThemeToSystem generates theme files, creates symlinks, and applies the theme
func ApplyThemeToSystem(p *color.Palette, wallpaperPath string) error {
	// Generate all embedded template files
	if err := GenerateThemeFiles(p, nil, wallpaperPath); err != nil {
		return fmt.Errorf("failed to generate theme files: %w", err)
	}

	// Create omarchy symlink if on omarchy system
	if err := CreateOmarchySymlink(); err != nil {
		return fmt.Errorf("failed to create omarchy symlink: %w", err)
	}

	// Run omarchy-theme-set if on omarchy system
	if err := RunOmarchyThemeSet(); err != nil {
		return fmt.Errorf("failed to apply omarchy theme: %w", err)
	}

	return nil
}
