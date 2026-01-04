package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// ANSI Colors - these use the terminal's color palette (0-15)
// The UI will automatically adapt when the user changes their terminal theme
var (
	// Basic ANSI colors (0-7)
	ANSIBlack   = lipgloss.Color("0")
	ANSIRed     = lipgloss.Color("1")
	ANSIGreen   = lipgloss.Color("2")
	ANSIYellow  = lipgloss.Color("3")
	ANSIBlue    = lipgloss.Color("4")
	ANSIMagenta = lipgloss.Color("5")
	ANSICyan    = lipgloss.Color("6")
	ANSIWhite   = lipgloss.Color("7")

	// Bright ANSI colors (8-15)
	ANSIBrightBlack   = lipgloss.Color("8")
	ANSIBrightRed     = lipgloss.Color("9")
	ANSIBrightGreen   = lipgloss.Color("10")
	ANSIBrightYellow  = lipgloss.Color("11")
	ANSIBrightBlue    = lipgloss.Color("12")
	ANSIBrightMagenta = lipgloss.Color("13")
	ANSIBrightCyan    = lipgloss.Color("14")
	ANSIBrightWhite   = lipgloss.Color("15")
)

// Semantic color mappings using ANSI colors
var (
	ColorPrimary   = ANSIMagenta       // Main accent color
	ColorSecondary = ANSIBrightBlue    // Secondary accent
	ColorAccent    = ANSIBrightMagenta // Highlights
	ColorMuted     = ANSIBrightBlack   // Comments/muted text
	ColorBorder    = ANSIBrightBlack   // Borders
	ColorBg        = ANSIBlack         // Background
	ColorFg        = ANSIWhite         // Foreground
	ColorSelected  = ANSIBlue          // Selected items
	ColorError     = ANSIRed           // Error messages
	ColorSuccess   = ANSIGreen         // Success messages
	ColorWarning   = ANSIYellow        // Warnings
	ColorInfo      = ANSICyan          // Info messages
)

// Styles for various UI elements
var (
	// App styles
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ANSIBrightWhite).
			Background(ColorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Panel styles
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1)

	PanelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			MarginBottom(1)

	// List styles
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	ListSelectedStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ANSIBrightWhite).
				Background(ColorSelected)

	// Color swatch style (placeholder, will be styled dynamically)
	SwatchStyle = lipgloss.NewStyle().
			Width(4).
			Height(1)

	// Help bar styles
	HelpBarStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)

	HelpKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Editor styles
	EditorLabelStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Width(12)

	EditorValueStyle = lipgloss.NewStyle().
				Foreground(ColorFg)

	EditorActiveStyle = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true)

	SliderTrackStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	SliderFillStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary)

	// Status styles
	StatusStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	// Preview styles
	PreviewHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				MarginBottom(1)

	// File picker styles
	FilePickerDirStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSecondary)

	FilePickerFileStyle = lipgloss.NewStyle().
				Foreground(ColorFg)

	FilePickerSelectedStyle = lipgloss.NewStyle().
				Background(ColorSelected).
				Foreground(ANSIBrightWhite)

	// Modal/overlay styles
	OverlayStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)
)

// ColorSwatch creates a colored block for displaying a color
func ColorSwatch(hexColor string) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(hexColor)).
		Width(4).
		Render("    ")
}

// ColorSwatchSmall creates a smaller colored block
func ColorSwatchSmall(hexColor string) string {
	return lipgloss.NewStyle().
		Background(lipgloss.Color(hexColor)).
		Width(2).
		Render("  ")
}
