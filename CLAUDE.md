# Peachy

Terminal UI (TUI) application for extracting color palettes from images. Built with Go using the Charmbracelet stack.

## Tech Stack

- **Bubbletea** - TUI framework (Elm architecture)
- **Lipgloss** - Styling and layout
- **Bubbles** - Reusable components

## Project Structure

```
internal/
├── ui/
│   ├── model.go        # Main app model, Update/View, state machine
│   ├── styles.go       # ANSI color definitions and lipgloss styles
│   ├── keys.go         # Key bindings
│   ├── views.go        # View utilities
│   └── components/
│       ├── colorlist.go    # 16-color palette list (Normal/Bright sections)
│       ├── coloreditor.go  # HSL sliders + hex input modal
│       ├── filepicker.go   # File browser with thumbnails
│       ├── preview.go      # Image + color swatch preview
│       ├── imagepreview.go # Unicode block image rendering
│       └── help.go         # Help overlay
├── color/
│   ├── extractor.go    # Color extraction from images
│   ├── palette.go      # 16-color palette management
│   ├── modes.go        # Extraction modes (Normal, Vibrant, Pastel, etc.)
│   └── convert.go      # Color space conversions
├── config/
│   ├── paths.go        # XDG paths, theme directories
│   ├── toml.go         # Theme file I/O
│   └── templates.go    # System theme generation
├── cache/
│   └── thumbnails.go   # Wallpaper thumbnail caching
└── terminal/
    └── blocks.go       # Unicode block rendering
```

## View States

The app uses a state machine with these views:
- `ViewMain` - Split pane: color list + preview
- `ViewEditor` - Color editing modal (HSL/hex)
- `ViewFilePicker` - File browser
- `ViewHelp` - Keyboard shortcuts overlay
- `ViewSaveTheme` - Theme name input
- `ViewThemeBrowser` - Load/apply saved themes

## Key Patterns

- All UI uses ANSI colors 0-15 for theming consistency
- Vim-style navigation (j/k/h/l/g/G)
- Components implement `Init()`, `Update()`, `View()` pattern
- Live preview updates during color editing

## Build & Run

```bash
go build
./peachy
```

## Key Bindings (Main View)

- `j/k` - Navigate colors
- `Enter` - Edit selected color
- `o` - Open image file picker
- `e` - Extract colors from image
- `m/M` - Cycle extraction mode
- `t` - Toggle light/dark mode
- `s` - Save theme
- `l` - Load theme
- `a` - Apply theme to system
- `?` - Help
- `q` - Quit
