package config

import (
	"fmt"
	"os"
	"strings"

	"peachy/internal/color"

	toml "github.com/pelletier/go-toml/v2"
)

// ColorsConfig represents the colors.toml file structure
type ColorsConfig struct {
	// Accent and UI colors
	Accent              string `toml:"accent"`
	ActiveBorderColor   string `toml:"active_border_color"`
	ActiveTabBackground string `toml:"active_tab_background"`

	// Cursor colors
	Cursor string `toml:"cursor"`

	// Primary colors
	Foreground string `toml:"foreground"`
	Background string `toml:"background"`

	// Selection colors
	SelectionForeground string `toml:"selection_foreground"`
	SelectionBackground string `toml:"selection_background"`

	// Normal colors (ANSI 0-7)
	Color0 string `toml:"color0"`
	Color1 string `toml:"color1"`
	Color2 string `toml:"color2"`
	Color3 string `toml:"color3"`
	Color4 string `toml:"color4"`
	Color5 string `toml:"color5"`
	Color6 string `toml:"color6"`
	Color7 string `toml:"color7"`

	// Bright colors (ANSI 8-15)
	Color8  string `toml:"color8"`
	Color9  string `toml:"color9"`
	Color10 string `toml:"color10"`
	Color11 string `toml:"color11"`
	Color12 string `toml:"color12"`
	Color13 string `toml:"color13"`
	Color14 string `toml:"color14"`
	Color15 string `toml:"color15"`
}

// LoadConfig loads a colors.toml file and returns a palette
func LoadConfig(path string) (*color.Palette, error) {
	path = ExpandPath(path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg ColorsConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return configToPalette(&cfg)
}

// SaveConfig saves a palette to a colors.toml file
func SaveConfig(path string, p *color.Palette) error {
	path = ExpandPath(path)

	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	cfg := paletteToConfig(p)

	// Format with comments for readability
	content := formatConfigWithComments(cfg)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// configToPalette converts a ColorsConfig to a Palette
func configToPalette(cfg *ColorsConfig) (*color.Palette, error) {
	p := color.NewPalette()

	colorFields := []string{
		cfg.Color0, cfg.Color1, cfg.Color2, cfg.Color3,
		cfg.Color4, cfg.Color5, cfg.Color6, cfg.Color7,
		cfg.Color8, cfg.Color9, cfg.Color10, cfg.Color11,
		cfg.Color12, cfg.Color13, cfg.Color14, cfg.Color15,
	}

	for i, hex := range colorFields {
		if hex == "" {
			continue
		}
		c, err := color.NewColorFromHex(hex)
		if err != nil {
			return nil, fmt.Errorf("invalid color at index %d: %w", i, err)
		}
		p.Colors[i] = c
	}

	// Set background and foreground
	if cfg.Background != "" {
		c, err := color.NewColorFromHex(cfg.Background)
		if err == nil {
			p.Background = c
		}
	}
	if cfg.Foreground != "" {
		c, err := color.NewColorFromHex(cfg.Foreground)
		if err == nil {
			p.Foreground = c
		}
	}

	return p, nil
}

// paletteToConfig converts a Palette to a ColorsConfig
func paletteToConfig(p *color.Palette) *ColorsConfig {
	return &ColorsConfig{
		// Accent colors default to blue
		Accent:              p.Colors[4].Hex,
		ActiveBorderColor:   p.Colors[12].Hex,
		ActiveTabBackground: p.Colors[4].Hex,

		// Cursor
		Cursor: p.Foreground.Hex,

		// Primary
		Foreground: p.Foreground.Hex,
		Background: p.Background.Hex,

		// Selection
		SelectionForeground: p.Background.Hex,
		SelectionBackground: p.Colors[4].Hex,

		// Normal colors
		Color0: p.Colors[0].Hex,
		Color1: p.Colors[1].Hex,
		Color2: p.Colors[2].Hex,
		Color3: p.Colors[3].Hex,
		Color4: p.Colors[4].Hex,
		Color5: p.Colors[5].Hex,
		Color6: p.Colors[6].Hex,
		Color7: p.Colors[7].Hex,

		// Bright colors
		Color8:  p.Colors[8].Hex,
		Color9:  p.Colors[9].Hex,
		Color10: p.Colors[10].Hex,
		Color11: p.Colors[11].Hex,
		Color12: p.Colors[12].Hex,
		Color13: p.Colors[13].Hex,
		Color14: p.Colors[14].Hex,
		Color15: p.Colors[15].Hex,
	}
}

// formatConfigWithComments formats the config with section comments
func formatConfigWithComments(cfg *ColorsConfig) string {
	return fmt.Sprintf(`# Peachy Theme - colors.toml

# Accent and UI colors
accent = "%s"
active_border_color = "%s"
active_tab_background = "%s"

# Cursor colors
cursor = "%s"

# Primary colors
foreground = "%s"
background = "%s"

# Selection colors
selection_foreground = "%s"
selection_background = "%s"

# Normal colors (ANSI 0-7)
color0 = "%s"
color1 = "%s"
color2 = "%s"
color3 = "%s"
color4 = "%s"
color5 = "%s"
color6 = "%s"
color7 = "%s"

# Bright colors (ANSI 8-15)
color8 = "%s"
color9 = "%s"
color10 = "%s"
color11 = "%s"
color12 = "%s"
color13 = "%s"
color14 = "%s"
color15 = "%s"
`,
		cfg.Accent, cfg.ActiveBorderColor, cfg.ActiveTabBackground,
		cfg.Cursor,
		cfg.Foreground, cfg.Background,
		cfg.SelectionForeground, cfg.SelectionBackground,
		cfg.Color0, cfg.Color1, cfg.Color2, cfg.Color3,
		cfg.Color4, cfg.Color5, cfg.Color6, cfg.Color7,
		cfg.Color8, cfg.Color9, cfg.Color10, cfg.Color11,
		cfg.Color12, cfg.Color13, cfg.Color14, cfg.Color15,
	)
}

// ConfigExists checks if the default config file exists
func ConfigExists() bool {
	path := GetDefaultConfigPath()
	_, err := os.Stat(path)
	return err == nil
}

// SaveTheme saves a palette to the themes directory with the given name
func SaveTheme(name string, p *color.Palette) error {
	if err := EnsureThemesDir(); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	path := GetThemePath(name)
	cfg := paletteToConfig(p)
	content := formatConfigWithComments(cfg)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write theme: %w", err)
	}

	return nil
}

// LoadTheme loads a theme by name from the themes directory
func LoadTheme(name string) (*color.Palette, error) {
	path := GetThemePath(name)
	return LoadConfig(path)
}

// ApplyTheme sets a theme as the active theme by writing its name to the theme file
func ApplyTheme(name string) error {
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Verify the theme exists
	themePath := GetThemePath(name)
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return fmt.Errorf("theme '%s' not found", name)
	}

	path := GetActiveThemePath()
	if err := os.WriteFile(path, []byte(name+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}

	return nil
}

// GetActiveTheme returns the name of the currently active theme
func GetActiveTheme() (string, error) {
	path := GetActiveThemePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// ListThemes returns a list of available theme names
func ListThemes() ([]string, error) {
	dir := GetThemesDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var themes []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".toml") {
			name := strings.TrimSuffix(entry.Name(), ".toml")
			themes = append(themes, name)
		}
	}
	return themes, nil
}

// ExportTheme exports a palette to a custom directory
func ExportTheme(p *color.Palette, imagePath, outputDir string) error {
	outputDir = ExpandPath(outputDir)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Save colors.toml
	colorsPath := outputDir + "/colors.toml"
	cfg := paletteToConfig(p)
	content := formatConfigWithComments(cfg)
	if err := os.WriteFile(colorsPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write colors.toml: %w", err)
	}

	return nil
}
