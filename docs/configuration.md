# Configuration

## File Locations

| Path | Description |
|------|-------------|
| `~/.config/peachy/colors.toml` | Default config file |
| `~/.config/peachy/themes/` | Saved named themes |
| `~/.config/peachy/current` | Active theme name |
| `~/.config/peachy/generated/` | Generated config files |
| `~/.cache/peachy/thumbnails/` | Cached image thumbnails |
| `~/Wallpapers/` | Default wallpaper directory |

## colors.toml Format

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
