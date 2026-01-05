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
  <a href="#cli-commands">CLI</a> •
  <a href="#tui-usage">TUI Usage</a> •
  <a href="#extraction-modes">Extraction Modes</a> •
  <a href="#theme-management">Themes</a> •
  <a href="#custom-templates">Templates</a>
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
# Launch the TUI
peachy

# Open TUI with an image
peachy ~/Pictures/wallpaper.png

# Generate theme from command line (no TUI)
peachy generate wallpaper.jpg

# Generate from random wallpaper and save
peachy generate --random --save mytheme

# Apply a saved theme
peachy apply mytheme
```

That's it. Press `o` to open the file picker, select an image, and watch the magic happen.

## CLI Commands

Peachy provides a full CLI for headless theme generation and management.

```bash
peachy [image]              # Launch TUI (optionally with image)
peachy generate <image>     # Generate theme from image
peachy apply <theme>        # Apply a saved theme
peachy list                 # List saved themes
peachy info <theme>         # Show theme color details
peachy export <theme> <dir> # Export theme to folder
peachy delete <theme>       # Delete a saved theme
```

### Generate Command

Generate themes from the command line without opening the TUI.

```bash
peachy generate wallpaper.jpg                    # Generate and apply
peachy generate wallpaper.jpg --save mytheme     # Save as named theme
peachy generate wallpaper.jpg --mode pastel      # Use pastel extraction
peachy generate wallpaper.jpg --light            # Generate light theme
peachy generate --random --save random-theme     # Random wallpaper from ~/Wallpapers
peachy generate wallpaper.jpg --no-apply         # Generate files only
peachy generate wallpaper.jpg -o ~/exports       # Export to custom directory
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-m, --mode` | Extraction mode: normal, monochromatic, analogous, pastel, material |
| `-l, --light` | Generate light theme instead of dark |
| `-s, --save` | Save theme with given name |
| `-o, --output` | Output directory for exported files |
| `-r, --random` | Use random wallpaper from ~/Wallpapers |
| `--no-apply` | Generate files only, don't apply theme |

### Export Command

Export a saved theme to a folder with all config files. Creates a complete theme folder suitable for Omarchy or manual use.

```bash
peachy export mytheme ~/themes/mytheme
peachy export mytheme ~/.config/omarchy/themes/mytheme
```

The exported folder contains:
- `btop.theme` - btop system monitor
- `chromium.theme` - Chromium browser
- `colors.toml` - Universal color config
- `icons.theme` - Desktop icon theme
- `neovim.lua` - Neovim colorscheme
- `vscode.empty.json` - VS Code theme

### Other Commands

```bash
# List all saved themes
peachy list

# Show detailed color info for a theme
peachy info mytheme

# Apply a saved theme to the system
peachy apply mytheme

# Delete a theme
peachy delete mytheme
```

### Root Flags

```bash
peachy -c ~/.config/custom/colors.toml    # Load custom config
peachy --version                          # Show version
peachy --help                             # Show help
```

## TUI Usage

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
| `s` | Save theme (prompts for name) |
| `l` | Browse and load themes |
| `a` | Apply current palette to system |
| `r` | Reset to extracted colors |
| `m` | Cycle extraction mode |
| `M` | Cycle extraction mode (reverse) |
| `t` | Toggle light/dark mode |
| `?` | Show help |
| `q` | Quit |

#### Theme Browser

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate themes |
| `Enter` | Load theme (populate colors) |
| `a` | Apply theme (save as active) |
| `Esc` | Cancel |

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

Press `s` in the TUI to save your current palette as a named theme. You'll be prompted to enter a name. Themes are stored in:

```
~/.config/peachy/themes/<name>.toml
```

### Browsing Themes

Press `l` to open the theme browser. From there you can:

- **Navigate** with `j/k` or arrow keys
- **Load** a theme with `Enter` — this populates the colors in the TUI for editing
- **Apply** a theme with `a` — this saves it as the active theme and writes `colors.toml`

### Applying Themes

You can also apply a theme from the command line:

```bash
# Apply a saved theme
peachy --apply mytheme

# This creates ~/.config/peachy/theme containing the theme name
# and saves the theme colors to ~/.config/peachy/colors.toml
```

### Theme Files

| Path | Description |
|------|-------------|
| `~/.config/peachy/themes/` | Directory containing saved themes |
| `~/.config/peachy/current` | File containing the active theme name |
| `~/.config/peachy/generated/` | Generated config files (ghostty, kitty, etc.) |

This makes it easy to integrate with other tools — just read `~/.config/peachy/current` to get the active theme name, or use the generated configs in the `generated/` directory.

## Omarchy Integration

On [Omarchy](https://omarchy.app) systems, Peachy provides seamless integration:

- **Template Processing**: Generates configs for ghostty, kitty, alacritty, hyprland, waybar, wofi, btop, neovim, and more
- **Automatic Symlink**: Creates `~/.config/omarchy/themes/peachy` → `~/.config/peachy/generated/`
- **System Apply**: Runs `omarchy-theme-set peachy` to apply the theme system-wide

When you press `a` (apply) in the theme browser or use `--apply`, Peachy will:
1. Generate all config files from embedded templates
2. Create the omarchy symlink
3. Run `omarchy-theme-set` to activate the theme

On non-Omarchy systems, themes are still saved and can be manually copied to your app configs.

## Configuration

### File Locations

| Path | Description |
|------|-------------|
| `~/.config/peachy/colors.toml` | Default config file |
| `~/.config/peachy/themes/` | Saved named themes |
| `~/.config/peachy/current` | Active theme name |
| `~/.config/peachy/generated/` | Generated config files |
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

## Custom Templates

Peachy supports custom templates to automatically generate config files for any application when you apply a theme. Templates work on both **Linux** and **macOS**.

### Quick Install

Install templates with the interactive installer:

```bash
# Run the installer
curl -fsSL https://raw.githubusercontent.com/bjarneo/peachy/main/docs/install-templates.sh | bash

# Or download and run locally
curl -O https://raw.githubusercontent.com/bjarneo/peachy/main/docs/install-templates.sh
chmod +x install-templates.sh
./install-templates.sh
```

### Available Templates

| Template | Description | Platform |
|----------|-------------|----------|
| alacritty | GPU-accelerated terminal | All |
| btop | Resource monitor | All |
| cava | Console audio visualizer | All |
| dunst | Notification daemon | All |
| foot | Fast Wayland terminal | Linux |
| ghostty | Zig-based terminal | All |
| gtk | GTK3/GTK4 theming | All |
| hyprland | Wayland compositor | Linux |
| hyprlock | Hyprland screen locker | Linux |
| iterm2 | macOS terminal | macOS |
| kitty | GPU terminal emulator | All |
| mako | Wayland notifications | Linux |
| neovim | aether.nvim colorscheme | All |
| rofi | Application launcher | All |
| swayosd | On-screen display | Linux |
| vencord | Discord theme | All |
| walker | Wayland launcher | Linux |
| warp | AI terminal | All |
| waybar | Wayland status bar | Linux |
| wofi | Wayland launcher | Linux |

### Manual Installation

Copy templates from the repo to your config:

```bash
# Clone and copy specific templates
git clone https://github.com/bjarneo/peachy.git
cp -r peachy/examples/templates/kitty ~/.config/peachy/templates/
cp -r peachy/examples/templates/alacritty ~/.config/peachy/templates/
```

### Template Commands

```bash
peachy templates list       # List installed templates
peachy templates validate   # Check templates for errors
peachy templates apply      # Apply templates manually
peachy templates init       # Create templates directory
```

### Creating Custom Templates

Create your own templates for any application. See [Custom Templates Documentation](docs/custom-templates.md) for details.

```bash
mkdir -p ~/.config/peachy/templates/myapp
```

Create `template.toml`:
```toml
name = "My App"
description = "Theme for my application"
condition = "myapp"  # Only run if myapp exists

[[files]]
template = "colors.conf"
destination = "~/.config/myapp/colors.conf"
```

Create `colors.conf` with color variables:
```
background = {background}
foreground = {foreground}
accent = {blue}
```

## Related Projects

- [Aether](https://github.com/bjarneo/aether) - GTK theme creator (Peachy's big sibling)

## License

MIT

---

<p align="center">
  Made with 🍑 by <a href="https://x.com/iamdothash">@iamdothash</a>
</p>
