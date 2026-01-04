package main

import (
	"fmt"
	"os"
	"strings"

	"peachy/internal/app"
	"peachy/internal/config"
)

func main() {
	// Parse command line arguments
	var imagePath, configPath, applyTheme string
	var showHelp, showVersion bool

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "-h", "--help":
			showHelp = true
		case "-v", "--version":
			showVersion = true
		case "-c", "--config":
			if i+1 < len(args) {
				i++
				configPath = args[i]
			} else {
				fmt.Fprintln(os.Stderr, "Error: --config requires a path argument")
				os.Exit(1)
			}
		case "-i", "--image":
			if i+1 < len(args) {
				i++
				imagePath = args[i]
			} else {
				fmt.Fprintln(os.Stderr, "Error: --image requires a path argument")
				os.Exit(1)
			}
		case "-a", "--apply":
			if i+1 < len(args) {
				i++
				applyTheme = args[i]
			} else {
				fmt.Fprintln(os.Stderr, "Error: --apply requires a theme name argument")
				os.Exit(1)
			}
		default:
			// Check if it looks like an image file
			if !strings.HasPrefix(arg, "-") {
				lower := strings.ToLower(arg)
				if strings.HasSuffix(lower, ".png") ||
					strings.HasSuffix(lower, ".jpg") ||
					strings.HasSuffix(lower, ".jpeg") ||
					strings.HasSuffix(lower, ".webp") ||
					strings.HasSuffix(lower, ".gif") ||
					strings.HasSuffix(lower, ".bmp") {
					imagePath = arg
				} else {
					fmt.Fprintf(os.Stderr, "Unknown argument: %s\n", arg)
					fmt.Fprintln(os.Stderr, "Use --help for usage information")
					os.Exit(1)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
				fmt.Fprintln(os.Stderr, "Use --help for usage information")
				os.Exit(1)
			}
		}
	}

	// Handle help and version
	if showHelp {
		app.PrintHelp()
		return
	}

	if showVersion {
		app.PrintVersion()
		return
	}

	// Handle apply theme
	if applyTheme != "" {
		if err := config.ApplyTheme(applyTheme); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Applied theme '%s'\n", applyTheme)
		return
	}

	// Create and configure the app
	application := app.New()

	if imagePath != "" {
		// Verify image exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Image file not found: %s\n", imagePath)
			os.Exit(1)
		}
		application.WithImage(imagePath)
	}

	if configPath != "" {
		application.WithConfig(configPath)
	}

	// Run the application
	if err := application.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
