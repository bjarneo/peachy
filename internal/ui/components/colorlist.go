package components

import (
	"fmt"
	"strings"

	"peachy/internal/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ColorList displays and allows navigation of the 16-color palette
type ColorList struct {
	palette *color.Palette
	cursor  int
	width   int
	height  int
	focused bool
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

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12"))

	sb.WriteString(titleStyle.Render("PALETTE"))
	sb.WriteString("\n")

	divider := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(strings.Repeat("─", c.width-4))
	sb.WriteString(divider)
	sb.WriteString("\n")

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true)

	// Normal section
	sb.WriteString(sectionStyle.Render("  Normal"))
	sb.WriteString("\n")

	for i := 0; i < 16; i++ {
		if i == 8 {
			sb.WriteString("\n")
			sb.WriteString(sectionStyle.Render("  Bright"))
			sb.WriteString("\n")
		}

		col := c.palette.GetColor(i)
		roleName := c.palette.GetRoleName(i)

		// 6-wide swatch for visual prominence
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(col.Hex)).
			Render("      ")

		indexStr := fmt.Sprintf("%2d", i)
		hexStr := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(col.Hex)

		roleStr := lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Render(roleName)

		if i == c.cursor && c.focused {
			// Selected: accent left bar + subtle highlight
			indicator := lipgloss.NewStyle().
				Foreground(lipgloss.Color("13")).
				Bold(true).
				Render("▍")

			idx := lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Bold(true).
				Render(indexStr)

			hex := lipgloss.NewStyle().
				Foreground(lipgloss.Color("13")).
				Render(col.Hex)

			role := lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Bold(true).
				Render(roleName)

			line := fmt.Sprintf("%s %s %s %s  %s", indicator, idx, swatch, hex, role)
			sb.WriteString(line)
		} else if i == c.cursor {
			line := fmt.Sprintf("  %s %s %s  %s", indexStr, swatch, hexStr, roleStr)
			sb.WriteString(line)
		} else {
			dimIdx := lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")).
				Render(indexStr)
			line := fmt.Sprintf("  %s %s %s  %s", dimIdx, swatch, hexStr, roleStr)
			sb.WriteString(line)
		}

		sb.WriteString("\n")
	}

	borderColor := lipgloss.Color("8")
	if c.focused {
		borderColor = lipgloss.Color("5")
	}

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1).
		Width(c.width)

	return borderStyle.Render(sb.String())
}
