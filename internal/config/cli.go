package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"peachy/internal/color"
)

// PrintColoredPalette prints the palette with colored blocks for CLI output
func PrintColoredPalette(p *color.Palette) {
	fmt.Println("\nPalette:")

	// Normal colors with colored blocks
	fmt.Print("  Normal: ")
	for i := 0; i < 8; i++ {
		c := p.GetColor(i)
		fmt.Print(ColorBlock(c.Hex))
	}
	fmt.Println()

	// Bright colors with colored blocks
	fmt.Print("  Bright: ")
	for i := 8; i < 16; i++ {
		c := p.GetColor(i)
		fmt.Print(ColorBlock(c.Hex))
	}
	fmt.Println()

	// FG/BG
	fmt.Printf("\n  FG: %s %s\n", ColorBlock(p.Foreground.Hex), p.Foreground.Hex)
	fmt.Printf("  BG: %s %s\n", ColorBlock(p.Background.Hex), p.Background.Hex)
}

// ColorBlock returns a colored block string using 24-bit ANSI escape codes
func ColorBlock(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	var r, g, b int
	if _, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b); err != nil {
		return "    " // Return empty block on parse error
	}

	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm    \x1b[0m", r, g, b)
}

// ExportAllFormats exports all theme template files to the specified directory
func ExportAllFormats(p *color.Palette, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	templateFiles, err := GetEmbeddedTemplateFiles()
	if err != nil {
		return fmt.Errorf("getting templates: %w", err)
	}

	variables := BuildTemplateVariables(p)

	for _, filename := range templateFiles {
		content, err := ReadTemplateFile(filename)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", filename, err)
		}

		processed := ProcessTemplateContent(content, variables)

		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, []byte(processed), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}
	}

	return nil
}
