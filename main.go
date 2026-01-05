package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"peachy/internal/app"
	"peachy/internal/color"
	"peachy/internal/config"
	"peachy/internal/shared"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "peachy [image]",
	Short: "A TUI theme creator for extracting colors from images",
	Long: `Peachy is a terminal UI application that extracts color palettes from images
and generates terminal themes. It supports multiple extraction modes and
integrates with Omarchy for system-wide theming.

Examples:
  peachy                        Start with default palette
  peachy wallpaper.png          Open TUI with image
  peachy generate sunset.jpg    Generate theme from image
  peachy apply mytheme          Apply a saved theme
  peachy list                   List saved themes`,
	Args: cobra.MaximumNArgs(1),
	RunE: runMain,
}

var generateCmd = &cobra.Command{
	Use:   "generate <image>",
	Short: "Generate theme from image without TUI",
	Long: `Extract colors from an image and generate theme files.

Extraction Modes:
  normal        Balanced extraction (default)
  monochromatic Single-hue variations
  analogous     Adjacent colors on color wheel
  pastel        Soft, muted colors
  material      Material Design inspired

Examples:
  peachy generate wallpaper.jpg
  peachy generate wallpaper.jpg --light
  peachy generate wallpaper.jpg --mode pastel
  peachy generate --random --save random-theme
  peachy generate wallpaper.jpg --save mytheme --output ~/exports`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGenerate,
}

var applyCmd = &cobra.Command{
	Use:   "apply <theme>",
	Short: "Apply a saved theme",
	Long: `Apply a previously saved theme to the system.
This loads the theme and generates configuration files.
On Omarchy systems, it also applies the theme system-wide.`,
	Args: cobra.ExactArgs(1),
	RunE: runApply,
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List saved themes",
	Long:    `List all themes saved in ~/.config/peachy/themes/`,
	Args:    cobra.NoArgs,
	RunE:    runList,
}

var exportCmd = &cobra.Command{
	Use:   "export <theme> <output-dir>",
	Short: "Export theme to a folder",
	Long: `Export a saved theme to a folder with all config files.
Creates a complete theme folder suitable for Omarchy or manual use.

The folder will contain all template files:
  btop.theme, chromium.theme, colors.toml, icons.theme, neovim.lua, vscode.empty.json

Examples:
  peachy export mytheme ~/themes/mytheme
  peachy export mytheme ~/.config/omarchy/themes/mytheme`,
	Args: cobra.ExactArgs(2),
	RunE: runExport,
}

var deleteCmd = &cobra.Command{
	Use:     "delete <theme>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a saved theme",
	Long:    `Delete a theme from ~/.config/peachy/themes/`,
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var infoCmd = &cobra.Command{
	Use:   "info <theme>",
	Short: "Show theme details",
	Long:  `Display detailed color information for a saved theme.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runInfo,
}

var (
	// Root command flags
	flagConfig string

	// Generate command flags
	flagLightMode   bool
	flagExtractMode string
	flagNoApply     bool
	flagOutput      string
	flagSave        string
	flagRandom      bool
)

func init() {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("Peachy v{{.Version}}\n")

	// Root flags
	rootCmd.Flags().StringVarP(&flagConfig, "config", "c", "", "path to colors.toml file")

	// Generate flags
	generateCmd.Flags().BoolVarP(&flagLightMode, "light", "l", false, "generate light theme")
	generateCmd.Flags().StringVarP(&flagExtractMode, "mode", "m", "normal", "extraction mode (normal, monochromatic, analogous, pastel, material)")
	generateCmd.Flags().BoolVar(&flagNoApply, "no-apply", false, "generate files only, don't apply theme")
	generateCmd.Flags().StringVarP(&flagOutput, "output", "o", "", "output directory for exported files")
	generateCmd.Flags().StringVarP(&flagSave, "save", "s", "", "save theme with given name")
	generateCmd.Flags().BoolVarP(&flagRandom, "random", "r", false, "use random wallpaper from ~/Wallpapers")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(infoCmd)
}

func runMain(cmd *cobra.Command, args []string) error {
	var imagePath string
	if len(args) > 0 {
		imagePath = args[0]
	}

	application := app.New()

	if imagePath != "" {
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			return fmt.Errorf("image file not found: %s", imagePath)
		}
		application.WithImage(imagePath)
	}

	if flagConfig != "" {
		application.WithConfig(flagConfig)
	}

	return application.Run()
}

func runGenerate(cmd *cobra.Command, args []string) error {
	var imagePath string

	if flagRandom {
		// Find random wallpaper
		var err error
		imagePath, err = findRandomWallpaper()
		if err != nil {
			return err
		}
		fmt.Printf("Selected: %s\n\n", filepath.Base(imagePath))
	} else if len(args) > 0 {
		imagePath = args[0]
	} else {
		return fmt.Errorf("either provide an image path or use --random")
	}

	// Verify image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return fmt.Errorf("image file not found: %s", imagePath)
	}

	// Parse extraction mode
	mode := color.ParseMode(flagExtractMode)

	// Extract colors
	extractor := color.NewExtractor()
	palette, err := extractor.ExtractPalette(imagePath, mode, flagLightMode)
	if err != nil {
		return fmt.Errorf("extracting colors: %w", err)
	}

	fmt.Printf("Extracted colors from %s\n", filepath.Base(imagePath))
	fmt.Printf("Mode: %s", color.ModeNames[mode])
	if flagLightMode {
		fmt.Print(" (light)")
	}
	fmt.Println()

	// Print palette preview with colors
	config.PrintColoredPalette(palette)

	// Save theme if requested
	if flagSave != "" {
		if err := config.SaveTheme(flagSave, palette); err != nil {
			return fmt.Errorf("saving theme: %w", err)
		}
		fmt.Printf("\nSaved theme '%s'\n", flagSave)
	}

	// Export to output directory if specified
	if flagOutput != "" {
		if err := config.ExportAllFormats(palette, flagOutput); err != nil {
			return fmt.Errorf("exporting theme: %w", err)
		}
		fmt.Printf("Exported to %s\n", flagOutput)
	}

	// Apply unless --no-apply
	if !flagNoApply && flagOutput == "" {
		if err := config.ApplyThemeToSystem(palette, imagePath); err != nil {
			return fmt.Errorf("applying theme: %w", err)
		}

		if config.IsOmarchyInstalled() {
			fmt.Println("\nApplied theme to system")
		} else {
			fmt.Printf("\nGenerated theme files to %s\n", config.GetPeachyThemeDir())
		}
	}

	return nil
}

func runApply(cmd *cobra.Command, args []string) error {
	themeName := args[0]

	palette, err := config.LoadTheme(themeName)
	if err != nil {
		return fmt.Errorf("loading theme: %w", err)
	}

	if err := config.ApplyTheme(themeName); err != nil {
		return fmt.Errorf("applying theme: %w", err)
	}

	if err := config.ApplyThemeToSystem(palette, ""); err != nil {
		return fmt.Errorf("applying to system: %w", err)
	}

	if config.IsOmarchyInstalled() {
		fmt.Printf("Applied theme '%s' to system\n", themeName)
	} else {
		fmt.Printf("Applied theme '%s'\n", themeName)
	}

	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	themes, err := config.ListThemes()
	if err != nil {
		return fmt.Errorf("listing themes: %w", err)
	}

	if len(themes) == 0 {
		fmt.Println("No themes saved yet.")
		fmt.Println("Use 'peachy generate <image> --save <name>' to create one.")
		return nil
	}

	fmt.Println("Saved themes:")
	for _, theme := range themes {
		fmt.Printf("  %s\n", theme)
	}

	return nil
}

func runExport(cmd *cobra.Command, args []string) error {
	themeName := args[0]
	outputDir := config.ExpandPath(args[1])

	palette, err := config.LoadTheme(themeName)
	if err != nil {
		return fmt.Errorf("loading theme: %w", err)
	}

	if err := config.ExportAllFormats(palette, outputDir); err != nil {
		return err
	}

	fmt.Printf("Exported theme '%s' to %s\n", themeName, outputDir)
	return nil
}

func runDelete(cmd *cobra.Command, args []string) error {
	themeName := args[0]

	themePath := config.GetThemePath(themeName)
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		return fmt.Errorf("theme '%s' not found", themeName)
	}

	if err := os.Remove(themePath); err != nil {
		return fmt.Errorf("deleting theme: %w", err)
	}

	fmt.Printf("Deleted theme '%s'\n", themeName)
	return nil
}

func runInfo(cmd *cobra.Command, args []string) error {
	themeName := args[0]

	palette, err := config.LoadTheme(themeName)
	if err != nil {
		return fmt.Errorf("loading theme: %w", err)
	}

	fmt.Printf("Theme: %s\n", themeName)
	fmt.Printf("Path:  %s\n\n", config.GetThemePath(themeName))

	config.PrintColoredPalette(palette)

	return nil
}

// Helper functions

func findRandomWallpaper() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}

	wallpaperDir := filepath.Join(home, "Wallpapers")
	if _, err := os.Stat(wallpaperDir); os.IsNotExist(err) {
		return "", fmt.Errorf("wallpaper directory not found: %s", wallpaperDir)
	}

	var images []string
	err = filepath.Walk(wallpaperDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			return nil
		}
		if shared.IsValidImage(path) {
			images = append(images, path)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("scanning wallpapers: %w", err)
	}

	if len(images) == 0 {
		return "", fmt.Errorf("no images found in %s", wallpaperDir)
	}

	return images[rand.Intn(len(images))], nil
}
