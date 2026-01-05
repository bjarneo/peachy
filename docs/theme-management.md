# Theme Management

Peachy supports saving and managing multiple named themes, making it easy to switch between different color schemes.

## Saving Themes

Press `s` in the TUI to save your current palette as a named theme. You'll be prompted to enter a name. Themes are stored in:

```
~/.config/peachy/themes/<name>.toml
```

## Browsing Themes

Press `l` to open the theme browser. From there you can:

- **Navigate** with `j/k` or arrow keys
- **Load** a theme with `Enter` — this populates the colors in the TUI for editing
- **Apply** a theme with `a` — this saves it as the active theme and writes `colors.toml`

## Applying Themes

You can also apply a theme from the command line:

```bash
# Apply a saved theme
peachy --apply mytheme

# This creates ~/.config/peachy/theme containing the theme name
# and saves the theme colors to ~/.config/peachy/colors.toml
```

## Theme Files

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
