# CLI Reference

Peachy provides a full CLI for headless theme generation and management.

## Commands

```bash
peachy [image]              # Launch TUI (optionally with image)
peachy generate <image>     # Generate theme from image
peachy apply <theme>        # Apply a saved theme
peachy list                 # List saved themes
peachy info <theme>         # Show theme color details
peachy export <theme> <dir> # Export theme to folder
peachy delete <theme>       # Delete a saved theme
```

## Generate Command

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

### Generate Flags

| Flag | Description |
|------|-------------|
| `-m, --mode` | Extraction mode: normal, monochromatic, analogous, pastel, material |
| `-l, --light` | Generate light theme instead of dark |
| `-s, --save` | Save theme with given name |
| `-o, --output` | Output directory for exported files |
| `-r, --random` | Use random wallpaper from ~/Wallpapers |
| `--no-apply` | Generate files only, don't apply theme |

## Export Command

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

## Other Commands

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

## Root Flags

```bash
peachy -c ~/.config/custom/colors.toml    # Load custom config
peachy --version                          # Show version
peachy --help                             # Show help
```

## Template Commands

```bash
peachy templates list       # List installed templates
peachy templates validate   # Check templates for errors
peachy templates apply      # Apply templates manually
peachy templates init       # Create templates directory
```
