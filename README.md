<p align="center">
  <h1 align="center">🍑 Peachy</h1>
  <p align="center">
    <strong>A delightful TUI for creating terminal color themes from images</strong>
  </p>
  <p align="center">
    Extract beautiful 16-color palettes from your wallpapers, fine-tune every color, and export to <code>colors.toml</code>
  </p>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#usage">Usage</a> •
  <a href="#theme-management">Theme Management</a> •
  <a href="#extraction-modes">Extraction Modes</a>
</p>

---

## Why Peachy?

Ever found the perfect wallpaper but couldn't get your terminal colors to match? Peachy solves that. It extracts dominant colors from any image using ImageMagick, intelligently assigns them to ANSI color roles, and lets you tweak each one until it's *just right*.

**No more copy-pasting hex codes. No more guessing which shade of blue works best. Just point, extract, tune, and export.**

## Features

- **Smart Color Extraction** - Uses ImageMagick to extract 16 dominant colors and auto-assigns them to ANSI roles (black, red, green, yellow, blue, magenta, cyan, white + bright variants)

- **Multiple Extraction Modes** - Choose from Normal, Material, Pastel, Monochromatic, or Analogous extraction algorithms for different aesthetic styles

- **Image Preview in Terminal** - See your wallpapers right in the TUI using Unicode half-block rendering. Works in any terminal with true color support.

- **Thumbnail Caching** - Scans `~/Wallpapers` on startup and creates cached thumbnails for instant previews. No lag, no waiting.

- **Precision Color Tuning** - Adjust any color via HSL sliders or direct hex input. See changes in real-time.

- **colors.toml Export** - Save your palette in a clean, documented TOML format compatible with [Aether](https://github.com/bjarneo/aether) and other theming tools.

- **Vim-style Navigation** - Navigate with `j/k`, edit with `Enter`, quit with `q`. Feels like home.

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
# Launch Peachy
./peachy

# Or with an image
./peachy ~/Pictures/wallpaper.png

# Or load an existing config
./peachy -c ~/.config/peachy/colors.toml
```

That's it. Press `o` to open the file picker, select an image, and watch the magic happen.

## Usage

### Main View

```
+------------------------------------------------------------------+
| Peachy - Theme Creator                    wallpaper.png [Normal] |
+--------------------------------+---------------------------------+
|  COLORS                        |  PREVIEW                        |
| ---------------------------    | ------------------------------- |
|   0 #### #1A1A2E  black        |  Source: wallpaper.png          |
|   1 #### #F54242  red          |  ################################|
| > 2 #### #42F554  green        |  ################################|
|   3 #### #F5F542  yellow       |                                 |
|   4 #### #4287F5  blue         |  Normal: ################       |
|   5 #### #F542F5  magenta      |  Bright: ################       |
|   6 #### #42F5F5  cyan         |                                 |
|   7 #### #E0E0E0  white        |  Sample:                        |
|   ...                          |  +---------------------------+  |
|                                |  | Hello, World!             |  |
|                                |  | Error Success Warning     |  |
|                                |  +---------------------------+  |
+--------------------------------+---------------------------------+
| j/k:nav  Enter:edit  o:open  e:extract  m:mode  s:save  ?:help   |
+------------------------------------------------------------------+
```

### Keyboard Shortcuts

#### Main View

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Enter` | Edit selected color |
| `o` | Open file picker |
| `e` | Re-extract colors from image |
| `s` | Save to colors.toml |
| `S` | Save as named theme |
| `l` | Load from colors.toml |
| `r` | Reset to extracted colors |
| `m` | Cycle extraction mode |
| `M` | Cycle extraction mode (reverse) |
| `t` | Toggle light/dark mode |
| `?` | Show help |
| `q` | Quit |

#### File Picker

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate |
| `Enter` | Select file / Enter directory |
| `Backspace` | Go to parent directory |
| `/` | Start search |
| `w` | Jump to ~/Wallpapers |
| `~` | Jump to home directory |
| `p` | Toggle image preview |
| `Esc` | Cancel |

#### Search Mode (after pressing `/`)

| Key | Action |
|-----|--------|
| Type | Filter files by name |
| `Enter` | Select file |
| `Esc` | Cancel search |
| `Up/Down` | Navigate results |

#### Color Editor

| Key | Action |
|-----|--------|
| `j` / `k` | Select field (Hue/Saturation/Lightness) |
| `h` / `l` | Adjust value (small step) |
| `H` / `L` | Adjust value (large step) |
| `#` | Enter hex mode |
| `u` | Reset to original |
| `Enter` | Confirm changes |
| `Esc` | Cancel |

## Extraction Modes

Peachy supports multiple extraction algorithms for different aesthetic styles. Press `m` to cycle through modes.

| Mode | Description |
|------|-------------|
| **Normal** | Auto-detects image type. Generates grayscale palette for monochrome images, subtle balanced palette for low-diversity images, or chromatic palette for colorful images. (Default) |
| **Material** | Uses Material Design backgrounds (#fafafa light / #121212 dark) with refined image colors. Clean, professional aesthetic. |
| **Pastel** | Soft, muted colors with low saturation and high lightness. Great for easy-on-the-eyes themes. |
| **Monochromatic** | Single hue derived from the most frequent chromatic color. Creates cohesive, focused themes. |
| **Analogous** | Adjacent hues on the color wheel (plus/minus 30 degrees). Creates harmonious, visually pleasing palettes. |

### Light/Dark Mode

Press `t` to toggle between light and dark mode generation. This affects:
- Background/foreground color selection
- Lightness levels for all palette colors
- Comment color brightness

## Theme Management

Peachy supports saving and managing multiple named themes, making it easy to switch between different color schemes.

### Saving Themes

Press `S` in the TUI to save your current palette as a named theme. Themes are stored in:

```
~/.config/peachy/themes/<name>.toml
```

### Applying Themes

Use the `--apply` flag to set a theme as active:

```bash
# Apply a saved theme
peachy --apply mytheme

# This creates ~/.config/peachy/theme containing the theme name
```

### Theme Files

| Path | Description |
|------|-------------|
| `~/.config/peachy/themes/` | Directory containing saved themes |
| `~/.config/peachy/theme` | File containing the active theme name |

This makes it easy to integrate with other tools — just read `~/.config/peachy/theme` to get the active theme name, then load the corresponding `.toml` from the themes directory.

## Configuration

### File Locations

| Path | Description |
|------|-------------|
| `~/.config/peachy/colors.toml` | Default config file |
| `~/.config/peachy/themes/` | Saved named themes |
| `~/.config/peachy/theme` | Active theme name |
| `~/.cache/peachy/thumbnails/` | Cached image thumbnails |
| `~/Wallpapers/` | Default wallpaper directory |

### colors.toml Format

```toml
# Peachy Theme - colors.toml

# Accent and UI colors
accent = "#4287F5"
active_border_color = "#6BA3F7"
active_tab_background = "#4287F5"

# Cursor colors
cursor = "#E0E0E0"

# Primary colors
foreground = "#E0E0E0"
background = "#1A1A2E"

# Selection colors
selection_foreground = "#1A1A2E"
selection_background = "#4287F5"

# Normal colors (ANSI 0-7)
color0 = "#1A1A2E"    # black
color1 = "#F54242"    # red
color2 = "#42F554"    # green
color3 = "#F5F542"    # yellow
color4 = "#4287F5"    # blue
color5 = "#F542F5"    # magenta
color6 = "#42F5F5"    # cyan
color7 = "#E0E0E0"    # white

# Bright colors (ANSI 8-15)
color8 = "#4A4A5E"    # bright_black
color9 = "#F76B6B"    # bright_red
color10 = "#6BF76B"   # bright_green
color11 = "#F7F76B"   # bright_yellow
color12 = "#6BA3F7"   # bright_blue
color13 = "#F76BF7"   # bright_magenta
color14 = "#6BF7F7"   # bright_cyan
color15 = "#FFFFFF"   # bright_white
```

### CLI Options

```bash
peachy [flags] [image]

Flags:
  -a, --apply <name>    Apply a saved theme (writes to ~/.config/peachy/theme)
  -c, --config <path>   Path to colors.toml file
  -i, --image <path>    Path to image file
  -h, --help            Show help
  -v, --version         Show version
```

**Examples:**

```bash
# Launch the TUI
peachy

# Extract colors from an image
peachy ~/Pictures/wallpaper.png

# Load an existing config
peachy -c ~/.config/peachy/colors.toml

# Apply a saved theme
peachy --apply nord
```

## How It Works

### Color Extraction

Peachy uses ImageMagick to extract dominant colors:

```bash
convert image.png -resize 100x100 -colors 16 -unique-colors txt:-
```

Colors are then intelligently assigned to ANSI roles based on:
- **Luminance** — Darkest colors → black/background, lightest → white/foreground
- **Hue** — Colors are grouped by hue ranges (red, yellow, green, cyan, blue, magenta)
- **Saturation** — More saturated colors are preferred for the primary slots

### Thumbnail Caching

For fast image previews, Peachy generates 200x200 thumbnails:

```
~/Wallpapers/sunset.png (4K, 12MB)
        ↓
~/.cache/peachy/thumbnails/a1b2c3d4.png (200x200, 20KB)
        ↓
Terminal Preview (Unicode blocks)
```

Thumbnails are keyed by `MD5(path + mtime)` and automatically regenerate when source files change.

### Terminal Image Rendering

Images are rendered using Unicode half-block characters (`▀`):

- Each character = 2 vertical pixels
- Foreground color = top pixel
- Background color = bottom pixel
- Works in any terminal with true color (24-bit) support

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [go-toml](https://github.com/pelletier/go-toml) - TOML parsing
- [imaging](https://github.com/disintegration/imaging) - Image processing

## Related Projects

- [Aether](https://github.com/bjarneo/aether) - GTK theme creator (Peachy's big sibling)
- [pywal](https://github.com/dylanaraps/pywal) - Generate color schemes from images

## Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT

---

<p align="center">
  Made with 🍑 by <a href="https://x.com/iamdothash">@iamdothash</a>
</p>
