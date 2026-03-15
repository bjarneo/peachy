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

	// Title with role name
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13"))

	roleName := color.RoleNames[color.ColorRole(e.colorIndex)]
	sb.WriteString(titleStyle.Render(fmt.Sprintf("Edit: %s", roleName)))
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(fmt.Sprintf("  (color %d)", e.colorIndex)))
	sb.WriteString("\n\n")

	// Side-by-side original vs current preview
	previewWidth := (e.width - 10) / 2
	if previewWidth < 10 {
		previewWidth = 10
	}

	origBox := lipgloss.NewStyle().
		Background(lipgloss.Color(e.original.Hex)).
		Width(previewWidth).
		Height(2).
		Align(lipgloss.Center)
	currBox := lipgloss.NewStyle().
		Background(lipgloss.Color(e.current.Hex)).
		Width(previewWidth).
		Height(2).
		Align(lipgloss.Center)

	origLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Width(previewWidth).
		Align(lipgloss.Center).
		Render("original")
	currLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("13")).
		Bold(true).
		Width(previewWidth).
		Align(lipgloss.Center).
		Render("current")

	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		origBox.Render(" "),
		" ",
		currBox.Render(" "),
	))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		origLabel,
		" ",
		currLabel,
	))
	sb.WriteString("\n\n")

	// HSL sliders
	labelStyle := lipgloss.NewStyle().
		Width(5).
		Foreground(lipgloss.Color("8"))

	activeLabel := lipgloss.NewStyle().
		Width(5).
		Foreground(lipgloss.Color("13")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Width(7).
		Foreground(lipgloss.Color("7")).
		Align(lipgloss.Right)

	type sliderInfo struct {
		label string
		value float64
		max   float64
		field EditorField
	}

	sliders := []sliderInfo{
		{"H", e.current.HSL.H, 360, FieldHue},
		{"S", e.current.HSL.S, 100, FieldSaturation},
		{"L", e.current.HSL.L, 100, FieldLightness},
	}

	for _, s := range sliders {
		active := e.activeField == s.field && !e.hexMode

		lbl := labelStyle
		if active {
			lbl = activeLabel
		}

		var valStr string
		if s.field == FieldHue {
			valStr = fmt.Sprintf("%.0f°", s.value)
		} else {
			valStr = fmt.Sprintf("%.0f%%", s.value)
		}

		sb.WriteString(lbl.Render(s.label))
		sb.WriteString(valueStyle.Render(valStr))
		sb.WriteString(" ")
		sb.WriteString(e.renderSlider(s.value, s.max, active))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	// Hex input
	hexLabel := labelStyle
	if e.hexMode {
		hexLabel = activeLabel
	}
	sb.WriteString(hexLabel.Render("Hex"))
	if e.hexMode {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Bold(true).
			Render(e.hexInput + "█"))
	} else {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Render("  " + e.current.Hex))
	}

	// RGB values
	sb.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(fmt.Sprintf("  R:%d G:%d B:%d", e.current.RGB.R, e.current.RGB.G, e.current.RGB.B)))

	sb.WriteString("\n\n")

	// Contextual help
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	if e.hexMode {
		sb.WriteString(helpStyle.Render("Type hex value"))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("  "))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("Enter"))
		sb.WriteString(helpStyle.Render(" confirm  "))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("Esc"))
		sb.WriteString(helpStyle.Render(" cancel"))
	} else {
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("h/l"))
		sb.WriteString(helpStyle.Render(" adjust  "))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("H/L"))
		sb.WriteString(helpStyle.Render(" big  "))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("#"))
		sb.WriteString(helpStyle.Render(" hex  "))
		sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true).Render("u"))
		sb.WriteString(helpStyle.Render(" reset"))
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5")).
		Padding(1).
		Width(e.width)

	return borderStyle.Render(sb.String())
}

func (e ColorEditor) renderSlider(value, max float64, active bool) string {
	width := e.width - 20
	if width < 12 {
		width = 12
	}
	if width > 25 {
		width = 25
	}

	filled := int((value / max) * float64(width))
	if filled > width {
		filled = width
	}

	track := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	fill := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	if active {
		fill = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	}

	// Use half-block for the cursor position to show exact value
	var slider string
	if filled > 0 {
		slider = fill.Render(strings.Repeat("━", filled))
	}
	if filled < width {
		slider += track.Render(strings.Repeat("─", width-filled))
	}

	return slider
}
