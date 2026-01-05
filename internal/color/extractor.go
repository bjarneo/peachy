package color

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"peachy/internal/shared"
)

// Extractor handles color extraction from images using ImageMagick
type Extractor struct{}

// NewExtractor creates a new color extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// CheckImageMagick verifies that ImageMagick is installed
func (e *Extractor) CheckImageMagick() error {
	// Try magick first (ImageMagick 7), then convert (ImageMagick 6)
	cmd := exec.Command("magick", "-version")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("convert", "-version")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("ImageMagick not found. Please install it: sudo pacman -S imagemagick")
		}
	}
	return nil
}

// colorWithCount holds a color and its pixel count for sorting by dominance
type colorWithCount struct {
	Color Color
	Count int
}

// ExtractColors extracts 32 dominant colors from an image using histogram
// Colors are returned sorted by frequency (most dominant first)
func (e *Extractor) ExtractColors(imagePath string) ([]Color, error) {
	if err := e.CheckImageMagick(); err != nil {
		return nil, err
	}

	// Use ImageMagick histogram to extract colors sorted by frequency
	// This matches Aether's approach: scale image, reduce to 32 colors, get histogram
	cmd := exec.Command("magick", imagePath,
		"-scale", "800x600>",
		"-colors", "32",
		"-depth", "8",
		"-quality", "85",
		"-format", "%c",
		"histogram:info:-")

	output, err := cmd.Output()
	if err != nil {
		// Fallback to convert command (ImageMagick 6)
		cmd = exec.Command("convert", imagePath,
			"-scale", "800x600>",
			"-colors", "32",
			"-depth", "8",
			"-quality", "85",
			"-format", "%c",
			"histogram:info:-")
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to extract colors: %w", err)
		}
	}

	colors, err := parseHistogramOutput(string(output))
	if err != nil {
		return nil, err
	}

	if len(colors) < 8 {
		return nil, fmt.Errorf("not enough colors extracted from image")
	}

	return colors, nil
}

// ExtractPalette extracts colors and assigns them to ANSI roles using the specified mode
func (e *Extractor) ExtractPalette(imagePath string, mode ExtractionMode, lightMode bool) (*Palette, error) {
	colors, err := e.ExtractColors(imagePath)
	if err != nil {
		return nil, err
	}

	var palette *Palette

	switch mode {
	case ModeMaterial:
		palette = GenerateMaterialPalette(colors, lightMode)
	case ModePastel:
		palette = GeneratePastelPalette(colors, lightMode)
	case ModeMonochromatic:
		palette = GenerateMonochromaticPalette(colors, lightMode)
	case ModeAnalogous:
		palette = GenerateAnalogousPalette(colors, lightMode)
	default:
		palette = GenerateNormalPalette(colors, lightMode)
	}

	// Apply brightness normalization for readability (matches Aether)
	NormalizeBrightness(palette)

	return palette, nil
}

// parseHistogramOutput parses ImageMagick histogram output
// Returns colors sorted by frequency (most dominant first)
func parseHistogramOutput(output string) ([]Color, error) {
	var colorData []colorWithCount

	// Match histogram lines: "  12345: (R,G,B) #RRGGBB srgb(...)"
	hexPattern := regexp.MustCompile(`^\s*(\d+):\s*\([^)]+\)\s*(#[0-9A-Fa-f]{6})`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		matches := hexPattern.FindStringSubmatch(line)
		if len(matches) >= 3 {
			count, _ := strconv.Atoi(matches[1])
			color, err := NewColorFromHex(matches[2])
			if err == nil {
				colorData = append(colorData, colorWithCount{Color: color, Count: count})
			}
		}
	}

	if len(colorData) == 0 {
		return nil, fmt.Errorf("no colors found in image")
	}

	// Sort by count descending (most dominant first)
	sort.Slice(colorData, func(i, j int) bool {
		return colorData[i].Count > colorData[j].Count
	})

	colors := make([]Color, len(colorData))
	for i, c := range colorData {
		colors[i] = c.Color
	}

	return colors, nil
}

// IsValidImage checks if the file is a valid image format
// This is a convenience wrapper around shared.IsValidImage for package compatibility
func IsValidImage(path string) bool {
	return shared.IsValidImage(path)
}
