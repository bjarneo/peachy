package components

import (
	"fmt"
	"strings"

	"peachy/internal/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditorField represents which field is being edited
type EditorField int

const (
	FieldHue EditorField = iota
	FieldSaturation
	FieldLightness
	FieldHex
)

// ColorEditor allows editing a single color via HSL or hex
type ColorEditor struct {
	original    color.Color
	current     color.Color
	colorIndex  int
	activeField EditorField
	hexInput    string
	hexMode     bool
	width       int
	height      int
	visible     bool
}

// NewColorEditor creates a new color editor
func NewColorEditor() ColorEditor {
	return ColorEditor{
		activeField: FieldHue,
		width:       40,
		height:      15,
		visible:     false,
	}
}

// Open opens the editor with a color to edit
func (e *ColorEditor) Open(col color.Color, index int) {
	e.original = col
	e.current = col
	e.colorIndex = index
	e.activeField = FieldHue
	e.hexInput = col.Hex
	e.hexMode = false
	e.visible = true
}

// Close closes the editor without saving
func (e *ColorEditor) Close() {
	e.visible = false
}

// Visible returns whether the editor is visible
func (e ColorEditor) Visible() bool {
	return e.visible
}

// CurrentColor returns the current edited color
func (e ColorEditor) CurrentColor() color.Color {
	return e.current
}

// OriginalColor returns the original color before editing
func (e ColorEditor) OriginalColor() color.Color {
	return e.original
}

// ColorIndex returns the index of the color being edited
func (e ColorEditor) ColorIndex() int {
	return e.colorIndex
}

// Reset resets the color to the original
func (e *ColorEditor) Reset() {
	e.current = e.original
	e.hexInput = e.original.Hex
}

// SetSize sets the editor size
func (e *ColorEditor) SetSize(w, h int) {
	e.width = w
	e.height = h
}

// Init implements tea.Model
func (e ColorEditor) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (e ColorEditor) Update(msg tea.Msg) (ColorEditor, tea.Cmd) {
	if !e.visible {
		return e, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if e.hexMode {
			return e.updateHexMode(msg)
		}
		return e.updateSliderMode(msg)
	}

	return e, nil
}

func (e ColorEditor) updateSliderMode(msg tea.KeyMsg) (ColorEditor, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if e.activeField > FieldHue {
			e.activeField--
		}
	case "down", "j":
		if e.activeField < FieldLightness {
			e.activeField++
		}
	case "left", "h":
		e.adjustValue(-1)
	case "right", "l":
		e.adjustValue(1)
	case "H":
		e.adjustValue(-10)
	case "L":
		e.adjustValue(10)
	case "#":
		e.hexMode = true
		e.hexInput = e.current.Hex
	case "u":
		e.Reset()
	}

	return e, nil
}

func (e ColorEditor) updateHexMode(msg tea.KeyMsg) (ColorEditor, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Try to parse hex input
		if col, err := color.NewColorFromHex(e.hexInput); err == nil {
			e.current = col
		}
		e.hexMode = false
	case "esc":
		e.hexInput = e.current.Hex
		e.hexMode = false
	case "backspace":
		if len(e.hexInput) > 0 {
			e.hexInput = e.hexInput[:len(e.hexInput)-1]
		}
	default:
		// Only allow hex characters
		char := msg.String()
		if len(char) == 1 && strings.ContainsAny(char, "0123456789abcdefABCDEF#") {
			if len(e.hexInput) < 7 {
				e.hexInput += strings.ToUpper(char)
			}
		}
	}

	return e, nil
}

func (e *ColorEditor) adjustValue(delta float64) {
	switch e.activeField {
	case FieldHue:
		newHue := e.current.HSL.H + delta*5
		e.current = e.current.WithHue(newHue)
	case FieldSaturation:
		newSat := e.current.HSL.S + delta*2
		e.current = e.current.WithSaturation(newSat)
	case FieldLightness:
		newLight := e.current.HSL.L + delta*2
		e.current = e.current.WithLightness(newLight)
	}
	e.hexInput = e.current.Hex
}

// View implements tea.Model
func (e ColorEditor) View() string {
	if !e.visible {
		return ""
	}

	var sb strings.Builder

	// Title - using ANSI colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")) // Bright Magenta

	roleName := color.RoleNames[color.ColorRole(e.colorIndex)]
	sb.WriteString(titleStyle.Render(fmt.Sprintf("Edit Color %d (%s)", e.colorIndex, roleName)))
	sb.WriteString("\n\n")

	// Color preview
	previewStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(e.current.Hex)).
		Width(e.width - 6).
		Height(2).
		Align(lipgloss.Center)
	sb.WriteString(previewStyle.Render(" "))
	sb.WriteString("\n\n")

	// HSL sliders - using ANSI colors
	labelStyle := lipgloss.NewStyle().
		Width(12).
		Foreground(lipgloss.Color("8")) // Bright Black

	activeStyle := lipgloss.NewStyle().
		Width(12).
		Foreground(lipgloss.Color("13")). // Bright Magenta
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Width(8).
		Foreground(lipgloss.Color("7")) // White

	// Hue
	hueLabel := labelStyle
	if e.activeField == FieldHue && !e.hexMode {
		hueLabel = activeStyle
	}
	sb.WriteString(hueLabel.Render("Hue"))
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%.0f°", e.current.HSL.H)))
	sb.WriteString(e.renderSlider(e.current.HSL.H, 360, e.activeField == FieldHue && !e.hexMode))
	sb.WriteString("\n")

	// Saturation
	satLabel := labelStyle
	if e.activeField == FieldSaturation && !e.hexMode {
		satLabel = activeStyle
	}
	sb.WriteString(satLabel.Render("Saturation"))
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%.0f%%", e.current.HSL.S)))
	sb.WriteString(e.renderSlider(e.current.HSL.S, 100, e.activeField == FieldSaturation && !e.hexMode))
	sb.WriteString("\n")

	// Lightness
	lightLabel := labelStyle
	if e.activeField == FieldLightness && !e.hexMode {
		lightLabel = activeStyle
	}
	sb.WriteString(lightLabel.Render("Lightness"))
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%.0f%%", e.current.HSL.L)))
	sb.WriteString(e.renderSlider(e.current.HSL.L, 100, e.activeField == FieldLightness && !e.hexMode))
	sb.WriteString("\n\n")

	// Hex input
	hexLabel := labelStyle
	if e.hexMode {
		hexLabel = activeStyle
	}
	sb.WriteString(hexLabel.Render("Hex"))
	if e.hexMode {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")). // Bright Magenta
			Render(e.hexInput + "█"))
	} else {
		sb.WriteString(valueStyle.Render(e.current.Hex))
	}
	sb.WriteString("\n\n")

	// Help
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")) // Bright Black

	if e.hexMode {
		sb.WriteString(helpStyle.Render("Type hex value • Enter: confirm • Esc: cancel"))
	} else {
		sb.WriteString(helpStyle.Render("h/l: adjust • H/L: big adjust • #: hex input • u: reset"))
	}

	// Border - using ANSI colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5")). // Magenta
		Padding(1).
		Width(e.width)

	return borderStyle.Render(sb.String())
}

func (e ColorEditor) renderSlider(value, max float64, active bool) string {
	width := 20
	filled := int((value / max) * float64(width))
	if filled > width {
		filled = width
	}

	// Using ANSI colors
	track := lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Bright Black
	fill := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))  // Magenta
	if active {
		fill = lipgloss.NewStyle().Foreground(lipgloss.Color("13")) // Bright Magenta
	}

	return fill.Render(strings.Repeat("█", filled)) +
		track.Render(strings.Repeat("░", width-filled))
}
