# VS Code Template

This template applies your Peachy color palette directly to VS Code's `settings.json` for live theming without needing a separate extension.

## Overview

Unlike traditional VS Code themes, this template modifies `settings.json` directly using `workbench.colorCustomizations` and `editor.tokenColorCustomizations`. This means:

- **Live updates** - Colors change immediately without restart
- **No extension needed** - Works with vanilla VS Code
- **Preserves settings** - Only touches color customizations, all other settings remain

## Requirements

- VS Code (`code` command available)
- `jq` - JSON processor for safe settings manipulation

Install jq:
```bash
# Arch
sudo pacman -S jq

# Ubuntu/Debian
sudo apt install jq

# macOS
brew install jq
```

## Installation

Copy the template:

```bash
cp -r /path/to/peachy/examples/templates/vscode ~/.config/peachy/templates/
```

Make the post-apply script executable:

```bash
chmod +x ~/.config/peachy/templates/vscode/post-apply
```

## How It Works

1. When you apply a Peachy theme, the template generates `apply-theme.sh` with your colors
2. The `post-apply` script runs `apply-theme.sh`
3. The script uses `jq` to safely merge color settings into `settings.json`
4. VS Code detects the settings change and updates colors live

## Files

### template.toml

```toml
name = "VS Code"
description = "Visual Studio Code live theme via settings.json"
version = "1.0"
condition = "code"

[[files]]
template = "apply-theme.sh"
destination = "~/.config/peachy/scripts/vscode-apply.sh"
```

### apply-theme.sh

A bash script that:
- Validates existing `settings.json`
- Creates backup if JSON is malformed
- Generates comprehensive UI color customizations
- Generates syntax token customizations with TextMate rules
- Merges colors into settings using `jq`

### post-apply

Runs the apply script after template processing:

```bash
#!/bin/bash
bash ~/.config/peachy/scripts/vscode-apply.sh
```

## What Gets Themed

### UI Elements
- Editor background and foreground
- Activity bar, sidebar, panels
- Tabs and title bar
- Status bar (with semantic colors for debugging/remote)
- Input fields, buttons, dropdowns
- Lists and trees
- Notifications and badges
- Scrollbars and minimap
- Breadcrumbs and menus

### Syntax Highlighting
- Comments (italic)
- Strings, numbers, booleans
- Keywords and operators
- Functions and methods
- Classes and types
- Variables and properties
- Tags and attributes (HTML/XML)
- Regex and escape sequences

### Git Integration
- Added/modified/deleted decorations
- Diff editor colors
- Merge conflict highlighting

### Terminal
- Full ANSI color support (16 colors)
- Cursor and selection colors

## Manual Application

If you want to apply colors without using the template system:

```bash
# Generate the script with current theme
peachy templates apply --theme mytheme

# Run manually
bash ~/.config/peachy/scripts/vscode-apply.sh
```

## Resetting Colors

To remove Peachy customizations and return to your VS Code theme:

```bash
# Remove color customizations from settings.json
jq 'del(.["workbench.colorCustomizations"]) | del(.["editor.tokenColorCustomizations"])' \
  ~/.config/Code/User/settings.json > /tmp/settings.json && \
  mv /tmp/settings.json ~/.config/Code/User/settings.json
```

## Troubleshooting

**"jq required" error:**
- Install jq (see Requirements above)

**Settings not updating:**
- Check VS Code is looking at the right settings file
- On some systems it may be `~/.config/Code - OSS/User/settings.json`
- Modify the `settings_file` path in `apply-theme.sh` if needed

**Colors partially applied:**
- Run `peachy templates validate` to check for template errors
- Check `~/.config/peachy/scripts/vscode-apply.sh` was generated

**Backup created warning:**
- Your `settings.json` had invalid JSON
- Check `settings.json.backup` if you need to recover settings

## Customization

Edit `apply-theme.sh` to customize which colors are applied. The script uses shell variables that reference Peachy's template variables with `.strip` modifier (removes `#` prefix):

```bash
black="{black.strip}"
red="{red.strip}"
# etc.
```

You can modify the `ui_colors` and `syntax_colors` heredocs to add, remove, or change color mappings.
