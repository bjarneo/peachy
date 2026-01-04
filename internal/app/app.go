package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"peachy/internal/ui"
)

// App represents the Peachy application
type App struct {
	imagePath  string
	configPath string
}

// New creates a new App instance
func New() *App {
	return &App{}
}

// WithImage sets the initial image path
func (a *App) WithImage(path string) *App {
	a.imagePath = path
	return a
}

// WithConfig sets the config file path
func (a *App) WithConfig(path string) *App {
	a.configPath = path
	return a
}

// Run starts the application
func (a *App) Run() error {
	// Create the model
	model := ui.NewModel()

	// Set initial paths if provided
	if a.imagePath != "" {
		model.SetImagePath(a.imagePath)
	}
	if a.configPath != "" {
		model.SetConfigPath(a.configPath)
	}

	// Create the program
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running application: %w", err)
	}

	return nil
}

// PrintVersion prints the version information
func PrintVersion() {
	fmt.Println("Peachy v0.1.0")
	fmt.Println("A TUI theme creator for extracting colors from images")
}

// PrintHelp prints the help message
func PrintHelp() {
	fmt.Println(`Peachy - TUI Theme Creator

Usage:
  peachy [flags] [image]

Flags:
  -a, --apply <name>    Apply a saved theme (writes to ~/.config/peachy/theme)
  -c, --config <path>   Path to colors.toml file
  -i, --image <path>    Path to image file
  -h, --help            Show this help message
  -v, --version         Show version information

Examples:
  peachy                         # Start with default palette
  peachy wallpaper.png           # Start and extract colors from image
  peachy -c ~/.config/peachy/colors.toml  # Load existing config
  peachy -i ~/Pictures/bg.jpg    # Specify image with flag
  peachy --apply mytheme         # Apply saved theme 'mytheme'

Keyboard Shortcuts:
  j/k, up/down  Navigate colors
  Enter         Edit selected color
  o             Open image file picker
  e             Extract colors from image
  m             Cycle extraction mode
  M             Cycle extraction mode (reverse)
  t             Toggle light/dark mode
  s             Save theme (prompts for name)
  l             Browse and load themes
  a             Apply current palette to system
  r             Reset to extracted colors
  ?             Show help
  q             Quit

Theme Browser:
  j/k           Navigate themes
  Enter         Load theme (populate colors in TUI)
  a             Apply theme (save to ~/.config/peachy/theme)
  Esc           Cancel

Extraction Modes:
  normal        Auto-detect image type (default)
  material      Material Design backgrounds with image colors
  pastel        Soft, muted colors
  monochromatic Single hue variations
  analogous     Adjacent hues on color wheel

Color Editor:
  h/l         Adjust value (small)
  H/L         Adjust value (large)
  j/k         Select field (Hue/Saturation/Lightness)
  #           Enter hex mode
  u           Reset to original
  Enter       Confirm changes
  Esc         Cancel changes

Theme Storage:
  Themes are saved to: ~/.config/peachy/themes/<name>.toml
  Active theme file:   ~/.config/peachy/current
  Generated configs:   ~/.config/peachy/generated/

Omarchy Integration:
  On Omarchy systems, applying a theme will:
  - Generate config files from embedded templates
  - Symlink to ~/.config/omarchy/themes/peachy
  - Run omarchy-theme-set to apply system-wide

Requirements:
  - ImageMagick (for color extraction)
`)
}

// CheckDependencies verifies required dependencies are installed
func CheckDependencies() error {
	// Check for ImageMagick
	_, err := os.Stat("/usr/bin/convert")
	if err != nil {
		// Try which
		_, err = os.Stat("/usr/local/bin/convert")
		if err != nil {
			return fmt.Errorf("ImageMagick not found. Please install it:\n  Arch: sudo pacman -S imagemagick\n  Ubuntu: sudo apt install imagemagick\n  macOS: brew install imagemagick")
		}
	}
	return nil
}
