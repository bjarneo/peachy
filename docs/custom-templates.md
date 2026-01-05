# Custom Templates

Generate configuration files for any application using your Peachy color palette.

> **Cross-Platform Support**: Custom templates work on both **Linux** and **macOS**. Templates are stored in `~/.config/peachy/templates/` on both platforms.

## Quick Start

1. Create a template folder:
```bash
mkdir -p ~/.config/peachy/templates/myapp
```

2. Add a manifest file (`template.toml`):
```toml
name = "My App"
description = "Theme for my application"

[[files]]
template = "colors.conf"
destination = "~/.config/myapp/colors.conf"
```

3. Create your template file (`colors.conf`):
```
foreground = {foreground}
background = {background}
accent = {blue}
```

4. Apply a theme and your template will be processed automatically:
```bash
peachy apply mytheme
```

## How It Works

When you apply a theme (via TUI or CLI), Peachy:

1. Discovers templates in `~/.config/peachy/templates/`
2. Checks conditions (if specified)
3. Processes template files, replacing color variables
4. Writes processed files to `~/.config/peachy/generated/templates/`
5. Creates symlinks from generated files to destinations
6. Runs post-apply scripts (if present)

## Directory Structure

```
~/.config/peachy/templates/
├── cava/
│   ├── template.toml     # Manifest (required)
│   ├── theme.ini         # Template file
│   └── post-apply        # Post-apply script (optional)
├── kitty/
│   ├── template.toml
│   └── colors.conf
└── waybar/
    ├── template.toml
    ├── colors.css
    └── post-apply
```

## Manifest Format

The `template.toml` file defines your template:

```toml
name = "Application Name"        # Required: Display name
description = "What this does"   # Optional: Description
version = "1.0"                  # Optional: Version string
condition = "app-binary"         # Optional: Command that must exist
disabled = false                 # Optional: Skip this template

[[files]]                        # At least one file required
template = "colors.conf"         # Source template filename
destination = "~/.config/app/colors.conf"  # Target path
```

### Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | Display name for the template |
| `description` | No | Brief description |
| `version` | No | Template version |
| `condition` | No | Command that must exist in PATH for template to be processed |
| `disabled` | No | Set to `true` to skip this template |
| `[[files]]` | Yes | One or more file mappings |
| `files.template` | Yes | Source template filename in the template folder |
| `files.destination` | Yes | Target path for symlink (supports `~`) |

### Multiple Files

A single template can generate multiple config files:

```toml
name = "Polybar"

[[files]]
template = "colors.ini"
destination = "~/.config/polybar/colors.ini"

[[files]]
template = "modules.ini"
destination = "~/.config/polybar/peachy-modules.ini"
```

### Conditional Templates

Only process if an application is installed:

```toml
name = "Cava"
condition = "cava"  # Only runs if 'cava' command exists

[[files]]
template = "theme.ini"
destination = "~/.config/cava/peachy.ini"
```

## Post-Apply Scripts

Add a `post-apply` script (no extension) to run commands after template processing:

```bash
#!/usr/bin/env bash
# Reload the application
pkill -USR2 myapp
```

Requirements:
- Must be executable (`chmod +x post-apply`)
- Script runs from the template directory
- Non-zero exit code is logged but doesn't stop other templates
- Use `#!/usr/bin/env bash` for cross-platform compatibility

## Platform-Specific Templates

Some applications are platform-specific (e.g., Hyprland on Linux, iTerm2 on macOS). Use the `condition` field to ensure templates only run when the application is available:

```toml
name = "Hyprland"
condition = "hyprctl"  # Only processes if hyprctl exists

[[files]]
template = "colors.conf"
destination = "~/.config/hypr/colors.conf"
```

This way, Linux-specific templates won't cause errors on macOS and vice versa.

## CLI Commands

```bash
# List all custom templates and their status
peachy templates list

# Validate templates (check for errors)
peachy templates validate

# Apply templates with a specific theme
peachy templates apply --theme mytheme

# Dry run - show what would happen
peachy templates apply --dry-run

# Create the templates directory
peachy templates init
```

## Installing Examples

Copy example templates from the Peachy repository:

```bash
# Copy all examples
cp -r /path/to/peachy/examples/templates/* ~/.config/peachy/templates/

# Or copy specific ones
cp -r /path/to/peachy/examples/templates/kitty ~/.config/peachy/templates/
```

## Next Steps

- [Template Variables](template-variables.md) - All available color variables and modifiers
- [Cava Example](examples/cava.md) - Audio visualizer template walkthrough
- [Zed Example](examples/zed.md) - Zed editor theme
- [VS Code Example](examples/vscode.md) - VS Code live theming via settings.json
