package components

import (
	"fmt"
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

	imageHeight := (p.height - 16) / 2
	if imageHeight < 3 {
		imageHeight = 3
	}

	var preview string
	var err error

	if p.thumbnailCache != nil {
		thumbPath, thumbErr := p.thumbnailCache.GetOrCreateThumbnail(p.imagePath)
		if thumbErr == nil && thumbPath != "" {
			preview, err = terminal.RenderImageFit(thumbPath, p.width-6, imageHeight)
		}
	}

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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12"))

	sb.WriteString(titleStyle.Render("PREVIEW"))
	sb.WriteString("\n")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(strings.Repeat("─", p.width-4))
	sb.WriteString(divider)
	sb.WriteString("\n")

	// Image preview
	if p.showImage && p.imagePreview != "" {
		infoStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
		sb.WriteString(infoStyle.Render(filepath.Base(p.imagePath)))
		sb.WriteString("\n")
		sb.WriteString(p.imagePreview)
	}

	// Paired color grid — normal and bright side by side
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true)

	sb.WriteString(labelStyle.Render("Colors"))
	sb.WriteString("\n")

	// Labels row
	names := []string{"BLK", "RED", "GRN", "YEL", "BLU", "MAG", "CYN", "WHT"}
	labelRow := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	for _, n := range names {
		sb.WriteString(labelRow.Render(fmt.Sprintf("%-4s", n)))
	}
	sb.WriteString("\n")

	// Normal row (0-7)
	for i := 0; i < 8; i++ {
		col := p.palette.GetColor(i)
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(col.Hex)).
			Render("    ")
		sb.WriteString(swatch)
	}
	sb.WriteString("\n")

	// Bright row (8-15)
	for i := 8; i < 16; i++ {
		col := p.palette.GetColor(i)
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(col.Hex)).
			Render("    ")
		sb.WriteString(swatch)
	}
	sb.WriteString("\n\n")

	// Syntax-highlighted code preview
	sb.WriteString(labelStyle.Render("Syntax Preview"))
	sb.WriteString("\n")

	sb.WriteString(p.renderCodePreview())
	sb.WriteString("\n")

	// FG/BG info
	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	fgSwatch := lipgloss.NewStyle().
		Background(lipgloss.Color(p.palette.Foreground.Hex)).
		Render("  ")
	bgSwatch := lipgloss.NewStyle().
		Background(lipgloss.Color(p.palette.Background.Hex)).
		Render("  ")

	sb.WriteString(infoStyle.Render("FG "))
	sb.WriteString(fgSwatch)
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(" " + p.palette.Foreground.Hex))
	sb.WriteString("   ")
	sb.WriteString(infoStyle.Render("BG "))
	sb.WriteString(bgSwatch)
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(" " + p.palette.Background.Hex))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1).
		Width(p.width)

	return borderStyle.Render(sb.String())
}

// renderCodePreview renders a fake code snippet with syntax highlighting
// using the palette colors for an instant visual of how the theme looks in an editor
func (p Preview) renderCodePreview() string {
	bg := p.palette.Background.Hex
	fg := p.palette.Foreground.Hex

	// Color helpers using actual palette hex values
	c := func(idx int, text string) string {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(p.palette.GetColor(idx).Hex)).
			Render(text)
	}

	// Build code lines
	var lines []string

	// Line 1: comment
	lines = append(lines, c(8, "// Extract colors from wallpaper"))

	// Line 2: func declaration
	lines = append(lines, c(4, "func")+" "+c(6, "extract")+"("+c(6, "path")+" "+c(5, "string")+") "+c(5, "Theme")+" {")

	// Line 3: variable assignment with string
	lines = append(lines, "    "+c(7, "img")+" := "+c(6, "openImage")+"("+c(7, "path")+")")

	// Line 4: variable with number
	lines = append(lines, "    "+c(7, "colors")+" := "+c(6, "quantize")+"("+c(7, "img")+", "+c(3, "48")+")")

	// Line 5: conditional
	lines = append(lines, "    "+c(4, "if")+" "+c(7, "len")+"("+c(7, "colors")+") "+c(1, ">")+" "+c(3, "0")+" {")

	// Line 6: string literal
	lines = append(lines, "        "+c(6, "fmt.Println")+"("+c(2, `"ready"`)+", "+c(11, "true")+")")

	// Line 7: close blocks
	lines = append(lines, "    }")

	// Line 8: return
	lines = append(lines, "    "+c(4, "return")+" "+c(5, "Theme")+"{"+c(1, "Name")+": "+c(2, `"peachy"`)+"}")

	// Line 9: close func
	lines = append(lines, "}")

	// Add line numbers and join
	lineNumStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(p.palette.GetColor(8).Hex)).
		Width(3).
		Align(lipgloss.Right)

	var codeLines []string
	for i, line := range lines {
		num := lineNumStyle.Render(fmt.Sprintf("%d", i+1))
		codeLines = append(codeLines, num+" "+line)
	}

	codeContent := strings.Join(codeLines, "\n")

	// Wrap in a background box using the palette's BG/FG
	codeBox := lipgloss.NewStyle().
		Background(lipgloss.Color(bg)).
		Foreground(lipgloss.Color(fg)).
		Padding(1).
		Width(p.width - 6)

	return codeBox.Render(codeContent)
}
