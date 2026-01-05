<p align="center">
  <h1 align="center">🍑 Peachy</h1>
  <p align="center">
    <strong>A delightful TUI for creating terminal color themes from images</strong>
  </p>
  <p align="center">
    Extract beautiful 16-color palettes from your wallpapers, fine-tune every color, and export to <code>colors.toml</code>
  </p>
</p>

---

## Why Peachy?

Ever found the perfect wallpaper but couldn't get your terminal colors to match? Peachy extracts dominant colors from any image using ImageMagick, intelligently assigns them to ANSI color roles, and lets you tweak each one until it's *just right*.

## Features

- **Smart Color Extraction** - Uses ImageMagick to extract 16 dominant colors and auto-assigns them to ANSI roles
- **Multiple Extraction Modes** - Normal, Material, Pastel, Monochromatic, or Analogous algorithms
- **Image Preview in Terminal** - See wallpapers using Unicode half-block rendering
- **Precision Color Tuning** - Adjust colors via HSL sliders or direct hex input
- **Vim-style Navigation** - Navigate with `j/k`, edit with `Enter`, quit with `q`

## Installation

### Prerequisites

- **Go 1.21+**
- **ImageMagick** (for color extraction)

```bash
# Arch Linux
sudo pacman -S imagemagick

# Ubuntu/Debian
sudo apt install imagemagick

# macOS
brew install imagemagick
```

### Build from Source

```bash
git clone https://github.com/bjarneo/peachy.git
cd peachy
go build -o peachy .

# Optional: Install to PATH
sudo mv peachy /usr/local/bin/
```

## Quick Start

```bash
# Launch the TUI
peachy

# Open TUI with an image
peachy ~/Pictures/wallpaper.png

# Generate theme from command line (no TUI)
peachy generate wallpaper.jpg

# Apply a saved theme
peachy apply mytheme
```

Press `o` to open the file picker, select an image, and watch the magic happen.

## Basic Usage

```bash
peachy [image]              # Launch TUI (optionally with image)
peachy generate <image>     # Generate theme from image
peachy apply <theme>        # Apply a saved theme
peachy list                 # List saved themes
peachy export <theme> <dir> # Export theme to folder
```

### Key TUI Shortcuts

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Edit selected color |
| `o` | Open file picker |
| `s` | Save theme |
| `a` | Apply palette to system |
| `m` | Cycle extraction mode |
| `t` | Toggle light/dark mode |
| `q` | Quit |

## Documentation

- [CLI Reference](docs/cli-reference.md) - Full command documentation
- [TUI Usage](docs/tui-usage.md) - Complete keyboard shortcuts
- [Extraction Modes](docs/extraction-modes.md) - Color extraction algorithms
- [Theme Management](docs/theme-management.md) - Saving and applying themes
- [Configuration](docs/configuration.md) - File locations and colors.toml format
- [Custom Templates](docs/custom-templates.md) - Create templates for any app
- [Template Variables](docs/template-variables.md) - Color variables and modifiers

## Custom Templates

Generate config files for any application when applying themes:

```bash
# Install templates interactively
curl -fsSL https://raw.githubusercontent.com/bjarneo/peachy/main/docs/install-templates.sh | bash

# List available templates
peachy templates list
```

Supports: alacritty, btop, cava, dunst, foot, ghostty, gtk, hyprland, iterm2, kitty, mako, neovim, rofi, waybar, wofi, and more.

## Related Projects

- [Aether](https://github.com/bjarneo/aether) - GTK theme creator (Peachy's big sibling)

## License

MIT

---

<p align="center">
  Made with 🍑 by <a href="https://x.com/iamdothash">@iamdothash</a>
</p>
