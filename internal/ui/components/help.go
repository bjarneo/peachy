package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Help displays a help overlay
type Help struct {
	visible bool
	width   int
	height  int
}

// NewHelp creates a new help component
func NewHelp() Help {
	return Help{
		visible: false,
		width:   60,
		height:  25,
	}
}

// Toggle toggles the help visibility
func (h *Help) Toggle() {
	h.visible = !h.visible
}

// Show shows the help
func (h *Help) Show() {
	h.visible = true
}

// Hide hides the help
func (h *Help) Hide() {
	h.visible = false
}

// Visible returns whether help is visible
func (h Help) Visible() bool {
	return h.visible
}

// SetSize sets the help overlay size
func (h *Help) SetSize(w, ht int) {
	h.width = w
	h.height = ht
}

// View renders the help overlay
func (h Help) View() string {
	if !h.visible {
		return ""
	}

	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))
	sb.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12"))

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")).
		Width(14)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7"))

	divStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	renderSection := func(title string, items []struct{ key, desc string }) {
		sb.WriteString(sectionStyle.Render(title))
		sb.WriteString("\n")
		sb.WriteString(divStyle.Render(strings.Repeat("─", 38)))
		sb.WriteString("\n")
		for _, item := range items {
			sb.WriteString(keyStyle.Render(item.key))
			sb.WriteString(descStyle.Render(item.desc))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	renderSection("Navigation", []struct{ key, desc string }{
		{"j / k", "Move down / up"},
		{"g / G", "First / last"},
		{"Enter", "Edit selected color"},
	})

	renderSection("Actions", []struct{ key, desc string }{
		{"o", "Open image file picker"},
		{"e", "Extract colors from image"},
		{"m / M", "Cycle extraction mode"},
		{"t", "Toggle light/dark mode"},
		{"s", "Save theme"},
		{"l", "Load/browse themes"},
		{"a", "Apply theme to system"},
		{"r", "Reset to extracted colors"},
		{"q", "Quit"},
	})

	renderSection("Color Editor", []struct{ key, desc string }{
		{"h / l", "Adjust value (small)"},
		{"H / L", "Adjust value (large)"},
		{"j / k", "Select field"},
		{"#", "Enter hex mode"},
		{"u", "Reset to original"},
		{"Enter", "Confirm changes"},
		{"Esc", "Cancel"},
	})

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true)
	sb.WriteString(footerStyle.Render("Press ? or Esc to close"))

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("5")).
		Padding(1, 2).
		Width(h.width)

	return borderStyle.Render(sb.String())
}
