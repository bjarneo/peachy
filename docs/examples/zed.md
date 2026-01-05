# Zed Editor Template

This template generates a theme for [Zed](https://zed.dev/), the high-performance code editor.

## Overview

The Zed template creates a complete theme file at `~/.config/zed/themes/peachy.json` with:
- Editor colors (background, foreground, syntax highlighting)
- UI elements (panels, tabs, status bar)
- Terminal colors
- Git decoration colors

## Installation

Copy the template to your Peachy templates directory:

```bash
cp -r /path/to/peachy/examples/templates/zed ~/.config/peachy/templates/
```

Or create it manually:

```bash
mkdir -p ~/.config/peachy/templates/zed
```

## Files

### template.toml

```toml
name = "Zed"
description = "Zed code editor theme"
version = "1.0"
condition = "zeditor"

[[files]]
template = "peachy.json"
destination = "~/.config/zed/themes/peachy.json"
```

The `condition = "zeditor"` ensures the template only processes when Zed is installed.

### peachy.json

The template generates a full Zed theme with:
- Dark appearance
- Semi-transparent backgrounds (using alpha suffixes like `90` for 90% opacity)
- Complete syntax highlighting for all language constructs
- Terminal ANSI colors matching your palette

## Activating the Theme

After applying a Peachy theme, activate it in Zed:

1. Open Zed
2. Press `Cmd+,` (macOS) or `Ctrl+,` (Linux) to open settings
3. Add or update the theme setting:

```json
{
  "theme": "Peachy"
}
```

Or use the command palette (`Cmd+Shift+P` / `Ctrl+Shift+P`) and search for "theme".

## Customization

### Adjusting Transparency

The template uses hex alpha values for transparency. Edit the template to adjust:

- `90` = 90% opacity (slightly transparent)
- `70` = 70% opacity
- `30` = 30% opacity (very transparent)

Example:
```json
"background": "{background}90",  // 90% opacity
"panel.background": "{black}90",
```

### Syntax Colors

Modify the syntax section in `peachy.json` to change how code is highlighted:

```json
"syntax": {
  "keyword": {
    "color": "{magenta}"  // Change to any palette color
  },
  "string": {
    "color": "{green}"
  }
}
```

## Troubleshooting

**Theme not appearing in Zed:**
- Ensure the file exists at `~/.config/zed/themes/peachy.json`
- Restart Zed after applying the theme
- Check the JSON is valid: `jq . ~/.config/zed/themes/peachy.json`

**Colors look wrong:**
- Verify your Peachy palette has been applied
- Check `peachy info <themename>` to see the colors

**Template not processing:**
- Verify Zed is installed: `which zeditor`
- Run `peachy templates validate`
