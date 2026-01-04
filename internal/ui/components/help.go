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

	// Title - using ANSI colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")). // Bright Magenta
		MarginBottom(1)
	sb.WriteString(titleStyle.Render("Peachy - Keyboard Shortcuts"))
	sb.WriteString("\n\n")

	// Section style - using ANSI colors
	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")) // Bright Blue

	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3")). // Yellow
		Width(15)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // White

	// Navigation section
	sb.WriteString(sectionStyle.Render("Navigation"))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 40))
	sb.WriteString("\n")

	shortcuts := []struct {
		key  string
		desc string
	}{
		{"j / Down", "Move down"},
		{"k / Up", "Move up"},
		{"g / Home", "Go to first"},
		{"G / End", "Go to last"},
		{"Enter", "Edit selected color"},
	}

	for _, s := range shortcuts {
		sb.WriteString(keyStyle.Render(s.key))
		sb.WriteString(descStyle.Render(s.desc))
		sb.WriteString("\n")
	}

	// Actions section
	sb.WriteString("\n")
	sb.WriteString(sectionStyle.Render("Actions"))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 40))
	sb.WriteString("\n")

	actions := []struct {
		key  string
		desc string
	}{
		{"o", "Open image file picker"},
		{"e", "Extract colors from image"},
		{"m", "Cycle extraction mode"},
		{"t", "Toggle light/dark mode"},
		{"s", "Save colors.toml"},
		{"l", "Load colors.toml"},
		{"r", "Reset to extracted colors"},
		{"?", "Toggle this help"},
		{"q / Ctrl+C", "Quit"},
	}

	for _, a := range actions {
		sb.WriteString(keyStyle.Render(a.key))
		sb.WriteString(descStyle.Render(a.desc))
		sb.WriteString("\n")
	}

	// Color editor section
	sb.WriteString("\n")
	sb.WriteString(sectionStyle.Render("Color Editor"))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 40))
	sb.WriteString("\n")

	editorKeys := []struct {
		key  string
		desc string
	}{
		{"h / l", "Adjust value (small)"},
		{"H / L", "Adjust value (large)"},
		{"j / k", "Select field"},
		{"#", "Enter hex mode"},
		{"u", "Reset to original"},
		{"Enter", "Confirm changes"},
		{"Esc", "Cancel changes"},
	}

	for _, e := range editorKeys {
		sb.WriteString(keyStyle.Render(e.key))
		sb.WriteString(descStyle.Render(e.desc))
		sb.WriteString("\n")
	}

	// Footer - using ANSI colors
	sb.WriteString("\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Bright Black
		Italic(true)
	sb.WriteString(footerStyle.Render("Press ? or Esc to close"))

	// Border - using ANSI colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("5")). // Magenta
		Padding(1, 2).
		Width(h.width)

	return borderStyle.Render(sb.String())
}
