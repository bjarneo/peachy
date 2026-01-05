package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Retro cassette futurism rainbow colors
var retroColors = []string{
	"#FF0055", // Hot pink/red
	"#FF6600", // Orange
	"#FFCC00", // Yellow
	"#00FF66", // Green
	"#00CCFF", // Cyan
	"#6633FF", // Purple
	"#FF00CC", // Magenta
}

// ASCII art logo
var logoLines = []string{
	`  ___  ___  ___  ___| |_  _ _`,
	` | . || -_||.'||  _||   || | |`,
	` |  _||___||__,||___||_|_||_  |`,
	` |_|                     |___|`,
}

// SplashModel is the model for the splash screen
type SplashModel struct {
	frame    int
	maxFrame int
	width    int
	height   int
}

// NewSplashModel creates a new splash screen model
func NewSplashModel() SplashModel {
	return SplashModel{
		frame:    0,
		maxFrame: 24, // ~2 seconds at 12fps
	}
}

// splashTickMsg is sent on each animation frame
type splashTickMsg time.Time

// splashTick sends tick messages for splash animation
func splashTick() tea.Cmd {
	return tea.Tick(time.Millisecond*83, func(t time.Time) tea.Msg {
		return splashTickMsg(t)
	})
}

// View renders the splash screen
func (m SplashModel) View() string {
	var b strings.Builder

	// Calculate stripe animation offset
	offset := m.frame % len(retroColors)

	// Build the animated rainbow stripes above logo
	stripeWidth := 50
	if m.width > 0 && m.width < stripeWidth {
		stripeWidth = m.width
	}

	// Render 3 stripe lines that scroll
	for row := 0; row < 3; row++ {
		var stripeLine strings.Builder
		for col := 0; col < stripeWidth; col++ {
			// Create moving diagonal stripe effect
			colorIdx := (col/3 + offset + row) % len(retroColors)
			color := retroColors[colorIdx]
			style := lipgloss.NewStyle().Background(lipgloss.Color(color))
			stripeLine.WriteString(style.Render(" "))
		}
		b.WriteString(stripeLine.String())
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Render logo with animated color cycling
	for i, line := range logoLines {
		// Each line gets a color from the rainbow, cycling with animation
		colorIdx := (i + offset) % len(retroColors)
		color := retroColors[colorIdx]
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Bold(true)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Render 3 stripe lines below logo
	for row := 0; row < 3; row++ {
		var stripeLine strings.Builder
		for col := 0; col < stripeWidth; col++ {
			// Reverse direction for bottom stripes
			colorIdx := (stripeWidth/3 - col/3 + offset + row) % len(retroColors)
			if colorIdx < 0 {
				colorIdx += len(retroColors)
			}
			color := retroColors[colorIdx]
			style := lipgloss.NewStyle().Background(lipgloss.Color(color))
			stripeLine.WriteString(style.Render(" "))
		}
		b.WriteString(stripeLine.String())
		b.WriteString("\n")
	}

	// Tagline
	tagline := lipgloss.NewStyle().
		Foreground(lipgloss.Color(retroColors[(offset+3)%len(retroColors)])).
		Italic(true).
		Render("Terminal themes from images")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("%s\n", tagline))

	// Center everything
	content := b.String()

	if m.width > 0 && m.height > 0 {
		contentLines := strings.Count(content, "\n") + 1
		topPadding := (m.height - contentLines) / 2
		if topPadding < 0 {
			topPadding = 0
		}

		centered := lipgloss.NewStyle().
			Width(m.width).
			Align(lipgloss.Center).
			Render(content)

		return strings.Repeat("\n", topPadding) + centered
	}

	return content
}
