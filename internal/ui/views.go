package ui

// This file contains additional view-related utilities.
// Main view logic is in model.go View() method.

import (
	"github.com/charmbracelet/lipgloss"
)

// CenterOverlay centers content as an overlay
func CenterOverlay(bg, overlay string, width, height int) string {
	// Simple overlay - just return the overlay for now
	// A more sophisticated implementation would layer them
	return overlay
}

// RenderColorBlock renders a color block with label
func RenderColorBlock(hexColor string, label string, width int) string {
	block := lipgloss.NewStyle().
		Background(lipgloss.Color(hexColor)).
		Width(width).
		Height(2).
		Align(lipgloss.Center).
		Render(label)

	return block
}

// TruncateWithEllipsis truncates a string with ellipsis if too long
func TruncateWithEllipsis(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
