package components

import (
	"path/filepath"
	"strings"

	"peachy/internal/cache"
	"peachy/internal/color"
	"peachy/internal/terminal"

	"github.com/charmbracelet/lipgloss"
)

// Preview displays a preview of the palette and optionally the source image
type Preview struct {
	palette        *color.Palette
	imagePath      string
	imagePreview   string
	width          int
	height         int
	showImage      bool
	thumbnailCache *cache.ThumbnailCache
}

// NewPreview creates a new preview component
func NewPreview(palette *color.Palette) Preview {
	return Preview{
		palette:   palette,
		width:     35,
		height:    20,
		showImage: true,
	}
}

// SetPalette updates the palette
func (p *Preview) SetPalette(pal *color.Palette) {
	p.palette = pal
}

// SetThumbnailCache sets the thumbnail cache
func (p *Preview) SetThumbnailCache(tc *cache.ThumbnailCache) {
	p.thumbnailCache = tc
}

// SetImage sets the source image for preview
func (p *Preview) SetImage(path string) {
	p.imagePath = path
	p.updateImagePreview()
}

// updateImagePreview renders the image preview
func (p *Preview) updateImagePreview() {
	if p.imagePath == "" {
		p.imagePreview = ""
		return
	}

	// Reserve space for the color swatches and sample (about 12 lines)
	// Use remaining space for image
	imageHeight := (p.height - 16) / 2
	if imageHeight < 3 {
		imageHeight = 3
	}

	// Always use cached thumbnail for efficiency
	var preview string
	var err error

	if p.thumbnailCache != nil {
		// Get or create thumbnail - always use the smaller cached version
		thumbPath, thumbErr := p.thumbnailCache.GetOrCreateThumbnail(p.imagePath)
		if thumbErr == nil && thumbPath != "" {
			// Use cached thumbnail (200x200 max) - much faster to render
			preview, err = terminal.RenderImageFit(thumbPath, p.width-6, imageHeight)
		}
	}

	// Only fall back to original if thumbnail cache is not available
	if preview == "" || err != nil {
		preview, err = terminal.RenderImageFit(p.imagePath, p.width-6, imageHeight)
	}

	if err != nil {
		p.imagePreview = ""
		return
	}
	p.imagePreview = preview
}

// SetSize sets the preview size
func (p *Preview) SetSize(w, h int) {
	p.width = w
	p.height = h
	p.updateImagePreview()
}

// View renders the preview
func (p Preview) View() string {
	var sb strings.Builder

	// Header - using ANSI colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")). // Bright Blue
		MarginBottom(1)
	sb.WriteString(titleStyle.Render("PREVIEW"))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", p.width-4))
	sb.WriteString("\n")

	// Image preview (if available)
	if p.showImage && p.imagePreview != "" {
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")) // Bright Black
		sb.WriteString(infoStyle.Render("Source: " + filepath.Base(p.imagePath)))
		sb.WriteString("\n")
		sb.WriteString(p.imagePreview)
	}

	// Color grid - Normal colors
	sb.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Bright Black
		Render("Normal:"))
	sb.WriteString("\n")

	// First row (0-7)
	for i := 0; i < 8; i++ {
		col := p.palette.GetColor(i)
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(col.Hex)).
			Render("    ")
		sb.WriteString(swatch)
	}
	sb.WriteString("\n\n")

	// Bright colors
	sb.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Bright Black
		Render("Bright:"))
	sb.WriteString("\n")

	// Second row (8-15)
	for i := 8; i < 16; i++ {
		col := p.palette.GetColor(i)
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(col.Hex)).
			Render("    ")
		sb.WriteString(swatch)
	}
	sb.WriteString("\n\n")

	// Sample text preview
	bgColor := p.palette.Background.Hex
	fgColor := p.palette.Foreground.Hex

	previewBg := lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor)).
		Padding(1).
		Width(p.width - 6)

	sb.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Bright Black
		Render("Sample:"))
	sb.WriteString("\n")

	// Build sample text with colors
	var sample strings.Builder
	sample.WriteString("Hello, World!\n")

	// Show some colored text
	sample.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(1).Hex)).
		Render("Error"))
	sample.WriteString(" ")
	sample.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(2).Hex)).
		Render("Success"))
	sample.WriteString(" ")
	sample.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(3).Hex)).
		Render("Warning"))
	sample.WriteString("\n")

	sample.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(4).Hex)).
		Render("Info"))
	sample.WriteString(" ")
	sample.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(5).Hex)).
		Render("Debug"))
	sample.WriteString(" ")
	sample.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(6).Hex)).
		Render("Trace"))

	sb.WriteString(previewBg.Render(sample.String()))
	sb.WriteString("\n\n")

	// Foreground/Background info - using ANSI colors
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")) // Bright Black

	fgSwatch := lipgloss.NewStyle().
		Background(lipgloss.Color(fgColor)).
		Render("  ")
	bgSwatch := lipgloss.NewStyle().
		Background(lipgloss.Color(bgColor)).
		Render("  ")

	sb.WriteString(infoStyle.Render("FG: "))
	sb.WriteString(fgSwatch)
	sb.WriteString(" " + fgColor)
	sb.WriteString("\n")
	sb.WriteString(infoStyle.Render("BG: "))
	sb.WriteString(bgSwatch)
	sb.WriteString(" " + bgColor)

	// Border - using ANSI colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")). // Bright Black
		Width(p.width).
		Height(p.height)

	return borderStyle.Render(sb.String())
}
