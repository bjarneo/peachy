package components

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"peachy/internal/terminal"
)

// ImagePreview displays an image preview using Unicode blocks
type ImagePreview struct {
	imagePath    string
	rendered     string
	width        int
	height       int
	err          error
	showInfo     bool
	imgWidth     int
	imgHeight    int
}

// NewImagePreview creates a new image preview component
func NewImagePreview() ImagePreview {
	return ImagePreview{
		width:    40,
		height:   15,
		showInfo: true,
	}
}

// SetSize sets the preview size
func (p *ImagePreview) SetSize(w, h int) {
	// Account for borders and padding
	p.width = w - 4
	p.height = h - 6

	if p.width < 10 {
		p.width = 10
	}
	if p.height < 5 {
		p.height = 5
	}

	// Re-render if we have an image
	if p.imagePath != "" {
		p.render()
	}
}

// SetImage sets the image to preview
func (p *ImagePreview) SetImage(path string) {
	p.imagePath = path
	p.err = nil
	p.rendered = ""

	if path == "" {
		return
	}

	// Get image dimensions
	w, h, err := terminal.GetImageDimensions(path)
	if err != nil {
		p.err = err
		return
	}
	p.imgWidth = w
	p.imgHeight = h

	p.render()
}

// render renders the image to Unicode blocks
func (p *ImagePreview) render() {
	if p.imagePath == "" {
		return
	}

	rendered, err := terminal.RenderImageFit(p.imagePath, p.width, p.height)
	if err != nil {
		p.err = err
		p.rendered = ""
		return
	}

	p.rendered = rendered
	p.err = nil
}

// Clear clears the preview
func (p *ImagePreview) Clear() {
	p.imagePath = ""
	p.rendered = ""
	p.err = nil
	p.imgWidth = 0
	p.imgHeight = 0
}

// HasImage returns whether an image is loaded
func (p ImagePreview) HasImage() bool {
	return p.imagePath != "" && p.rendered != ""
}

// View renders the preview
func (p ImagePreview) View() string {
	var sb strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A78BFA")).
		MarginBottom(1)
	sb.WriteString(titleStyle.Render("IMAGE PREVIEW"))
	sb.WriteString("\n")

	if p.imagePath == "" {
		// No image
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")).
			Italic(true)
		sb.WriteString("\n")
		sb.WriteString(emptyStyle.Render("No image selected"))
		sb.WriteString("\n")
		sb.WriteString(emptyStyle.Render("Press 'o' to open file picker"))
	} else if p.err != nil {
		// Error
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444"))
		sb.WriteString("\n")
		sb.WriteString(errorStyle.Render("Error: " + p.err.Error()))
	} else {
		// Show image info
		if p.showInfo {
			infoStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6B7280"))
			filename := filepath.Base(p.imagePath)
			if len(filename) > p.width-5 {
				filename = filename[:p.width-8] + "..."
			}
			sb.WriteString(infoStyle.Render(fmt.Sprintf("%s (%dx%d)", filename, p.imgWidth, p.imgHeight)))
			sb.WriteString("\n")
		}

		// Render the image
		sb.WriteString(p.rendered)
	}

	// Border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#374151")).
		Width(p.width + 4).
		Height(p.height + 4)

	return borderStyle.Render(sb.String())
}

// ViewCompact renders a compact version (for file picker)
func (p ImagePreview) ViewCompact() string {
	if p.imagePath == "" || p.err != nil {
		return ""
	}

	var sb strings.Builder

	// Just the image, no border
	sb.WriteString(p.rendered)

	return sb.String()
}
