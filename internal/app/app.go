package app

import (
	"fmt"
	"os"

	"peachy/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
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
