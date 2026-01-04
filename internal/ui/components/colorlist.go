package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"peachy/internal/color"
)

// ColorList displays and allows navigation of the 16-color palette
type ColorList struct {
	palette  *color.Palette
	cursor   int
	width    int
	height   int
	focused  bool
}

// NewColorList creates a new color list component
func NewColorList(palette *color.Palette) ColorList {
	return ColorList{
		palette: palette,
		cursor:  0,
		width:   40,
		height:  20,
		focused: true,
	}
}

// SetPalette updates the palette
func (c *ColorList) SetPalette(p *color.Palette) {
	c.palette = p
}

// GetPalette returns the current palette
func (c *ColorList) GetPalette() *color.Palette {
	return c.palette
}

// Cursor returns the current cursor position
func (c ColorList) Cursor() int {
	return c.cursor
}

// SetCursor sets the cursor position
func (c *ColorList) SetCursor(pos int) {
	if pos >= 0 && pos < 16 {
		c.cursor = pos
	}
}

// SelectedColor returns the currently selected color
func (c ColorList) SelectedColor() color.Color {
	return c.palette.GetColor(c.cursor)
}

// SetColor updates a color in the palette
func (c *ColorList) SetColor(index int, col color.Color) {
	c.palette.SetColor(index, col)
}

// Focus sets the focus state
func (c *ColorList) Focus() {
	c.focused = true
}

// Blur removes focus
func (c *ColorList) Blur() {
	c.focused = false
}

// Focused returns whether the list is focused
func (c ColorList) Focused() bool {
	return c.focused
}

// SetSize updates the component size
func (c *ColorList) SetSize(w, h int) {
	c.width = w
	c.height = h
}

// Init implements tea.Model
func (c ColorList) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (c ColorList) Update(msg tea.Msg) (ColorList, tea.Cmd) {
	if !c.focused {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if c.cursor > 0 {
				c.cursor--
			}
		case "down", "j":
			if c.cursor < 15 {
				c.cursor++
			}
		case "home", "g":
			c.cursor = 0
		case "end", "G":
			c.cursor = 15
		}
	}

	return c, nil
}

// View implements tea.Model
func (c ColorList) View() string {
	var sb strings.Builder

	// Header - using ANSI colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")). // Bright Blue
		MarginBottom(1)
	sb.WriteString(titleStyle.Render("COLORS"))
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", c.width-4))
	sb.WriteString("\n")

	// Color list
	for i := 0; i < 16; i++ {
		col := c.palette.GetColor(i)
		roleName := c.palette.GetRoleName(i)

		// Create swatch
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(col.Hex)).
			Render("    ")

		// Format line
		indexStr := fmt.Sprintf("%2d", i)
		hexStr := col.Hex

		// Style based on selection - using ANSI colors
		var line string
		if i == c.cursor && c.focused {
			line = fmt.Sprintf("> %s %s %s  %s", indexStr, swatch, hexStr, roleName)
			line = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")). // Bright White
				Background(lipgloss.Color("4")).  // Blue
				Width(c.width - 4).
				Render(line)
		} else {
			prefix := "  "
			if i == c.cursor {
				prefix = "> "
			}
			line = fmt.Sprintf("%s%s %s %s  %s", prefix, indexStr, swatch, hexStr, roleName)
		}

		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Border style - using ANSI colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")). // Bright Black
		Width(c.width).
		Height(c.height)

	if c.focused {
		borderStyle = borderStyle.BorderForeground(lipgloss.Color("5")) // Magenta
	}

	return borderStyle.Render(sb.String())
}
