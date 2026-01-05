# VS Code Template

This template creates a VS Code theme extension that appears in the theme picker.

## Overview

The template creates a proper VS Code extension at `~/.vscode/extensions/peachy-theme/` with:
- Full UI theming (editor, sidebar, panels, tabs, etc.)
- Syntax highlighting with semantic token colors
- Terminal ANSI colors
- Git decoration colors

## Installation

Copy the template:

```bash
cp -r /path/to/peachy/examples/templates/vscode ~/.config/peachy/templates/
```

## How It Works

1. When you apply a Peachy theme, the template generates:
   - `~/.vscode/extensions/peachy-theme/package.json` - Extension manifest
   - `~/.vscode/extensions/peachy-theme/themes/peachy-color-theme.json` - Theme colors
2. VS Code detects the extension on next launch
3. Select "Peachy" from the theme picker

## Activating the Theme

After applying a Peachy theme:

1. Restart VS Code (or reload window with `Ctrl+Shift+P` → "Reload Window")
2. Open the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`)
3. Type "Color Theme" and select "Preferences: Color Theme"
4. Choose "Peachy" from the list

Or add to your `settings.json`:

```json
{
  "workbench.colorTheme": "Peachy"
}
```

## Files

### template.toml

```toml
name = "VS Code"
description = "Visual Studio Code theme extension"
version = "1.0"
condition = "code"

[[files]]
template = "package.json"
destination = "~/.vscode/extensions/peachy-theme/package.json"

[[files]]
template = "peachy-color-theme.json"
destination = "~/.vscode/extensions/peachy-theme/themes/peachy-color-theme.json"
```

### package.json

Extension manifest that registers the theme with VS Code.

### peachy-color-theme.json

Complete theme definition with:
- Editor colors
- UI element colors
- Syntax token colors
- Terminal colors

## What Gets Themed

### Editor
- Background and foreground
- Line numbers and cursor
- Selection and find highlights
- Bracket matching and highlights
- Indent guides and rulers

### UI Elements
- Activity bar and sidebar
- Tabs and title bar
- Status bar
- Panels and terminals
- Quick picker and menus
- Notifications

### Syntax Highlighting
- Comments, strings, numbers
- Keywords and operators
- Functions and methods
- Classes and types
- Variables and properties
- HTML/XML tags and attributes
- Markdown formatting
- JSON keys
- CSS properties and selectors

### Git Integration
- Added/modified/deleted decorations
- Diff editor colors
- Merge conflict highlighting

## Customization

Edit `peachy-color-theme.json` to customize colors. The template uses Peachy variables:

- `{foreground}`, `{background}` - Main colors
- `{black}` through `{white}` - Normal ANSI colors
- `{bright_black}` through `{bright_white}` - Bright ANSI colors

You can also add hex alpha suffixes for transparency:
```json
"editor.selectionBackground": "{blue}44"
```

## Troubleshooting

**Theme not appearing:**
- Restart VS Code after applying
- Check files exist in `~/.vscode/extensions/peachy-theme/`
- Verify JSON is valid: `jq . ~/.vscode/extensions/peachy-theme/themes/peachy-color-theme.json`

**Colors not updating:**
- VS Code caches themes - restart after re-applying
- Delete `~/.vscode/extensions/peachy-theme/` and re-apply

**Template not processing:**
- Verify VS Code is installed: `which code`
- Run `peachy templates validate`
