# Cava Template Walkthrough

This guide walks through creating a custom template for [Cava](https://github.com/karlstav/cava), a console-based audio visualizer.

## Goal

Create a Peachy template that:
- Generates a color theme for Cava
- Uses gradient colors from your palette
- Reloads Cava automatically when colors change

## Step 1: Create the Template Folder

```bash
mkdir -p ~/.config/peachy/templates/cava
```

## Step 2: Create the Manifest

Create `~/.config/peachy/templates/cava/template.toml`:

```toml
name = "Cava"
description = "Audio visualizer with gradient colors"
version = "1.0"
condition = "cava"

[[files]]
template = "theme.ini"
destination = "~/.config/cava/peachy.ini"
```

Key points:
- `condition = "cava"` ensures this only runs if Cava is installed
- The template outputs to `~/.config/cava/peachy.ini`

## Step 3: Create the Template File

Create `~/.config/peachy/templates/cava/theme.ini`:

```ini
[color]
; Peachy generated theme
background = 'default'
foreground = '{magenta}'

; Gradient mode for colorful visualization
gradient = 1
gradient_count = 6

gradient_color_1 = '{blue}'
gradient_color_2 = '{cyan}'
gradient_color_3 = '{green}'
gradient_color_4 = '{yellow}'
gradient_color_5 = '{magenta}'
gradient_color_6 = '{red}'
```

This creates a rainbow gradient using your palette colors.

## Step 4: Create the Post-Apply Script

Create `~/.config/peachy/templates/cava/post-apply`:

```bash
#!/bin/bash
# Reload cava if running
pgrep -x cava && pkill -USR2 cava
exit 0
```

Make it executable:

```bash
chmod +x ~/.config/peachy/templates/cava/post-apply
```

## Step 5: Configure Cava

Edit your Cava config (`~/.config/cava/config`) to use the theme:

```ini
[color]
; Use Peachy theme
; Note: You need to source/include the theme file
; Cava doesn't support includes, so we reference the gradient colors directly

; Or simply copy the content when it's generated
```

Since Cava doesn't support includes, you have two options:

**Option A:** Symlink the entire config
```bash
ln -sf ~/.config/cava/peachy.ini ~/.config/cava/config
```

**Option B:** Update the post-apply script to merge configs:
```bash
#!/bin/bash
# Copy theme colors to config
config="$HOME/.config/cava/config"
theme="$HOME/.config/cava/peachy.ini"

if [[ -f "$config" ]] && [[ -f "$theme" ]]; then
    # Remove existing [color] section and append new one
    sed -i '/^\[color\]/,/^\[/d' "$config"
    cat "$theme" >> "$config"
fi

# Reload cava
pgrep -x cava && pkill -USR2 cava
exit 0
```

## Step 6: Test

Validate your template:

```bash
peachy templates validate
```

Apply a theme:

```bash
peachy apply mytheme
```

Check the output:

```bash
cat ~/.config/cava/peachy.ini
```

You should see your template with color variables replaced by actual hex values.

## Variations

### Monochrome Mode

For a single-color visualizer:

```ini
[color]
background = 'default'
foreground = '{blue}'
gradient = 0
```

### Two-Tone Gradient

```ini
[color]
gradient = 1
gradient_count = 2
gradient_color_1 = '{blue}'
gradient_color_2 = '{magenta}'
```

### Match Terminal Background

```ini
[color]
background = '{background}'
foreground = '{foreground}'
```

## Troubleshooting

**Template not processing:**
- Check if Cava is installed: `which cava`
- Run `peachy templates validate` to check for errors

**Colors not updating:**
- Ensure the post-apply script is executable
- Check if Cava is running: `pgrep cava`
- Try manually sending the reload signal: `pkill -USR2 cava`

**Symlink not created:**
- Check if the destination directory exists
- Look for errors in the apply output
