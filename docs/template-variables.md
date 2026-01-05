# Template Variables

Reference for all color variables and modifiers available in Peachy templates.

## Color Variables

### Primary Colors

| Variable | Description |
|----------|-------------|
| `{background}` | Primary background color |
| `{foreground}` | Primary text/foreground color |

### ANSI Colors (0-7)

| Variable | Color | Typical Use |
|----------|-------|-------------|
| `{black}` | Black (ANSI 0) | Dark backgrounds |
| `{red}` | Red (ANSI 1) | Errors, deletions |
| `{green}` | Green (ANSI 2) | Success, additions |
| `{yellow}` | Yellow (ANSI 3) | Warnings, highlights |
| `{blue}` | Blue (ANSI 4) | Info, links, accents |
| `{magenta}` | Magenta (ANSI 5) | Keywords, special |
| `{cyan}` | Cyan (ANSI 6) | Strings, secondary |
| `{white}` | White (ANSI 7) | Light text |

### Bright ANSI Colors (8-15)

| Variable | Color | Typical Use |
|----------|-------|-------------|
| `{bright_black}` | Bright Black (ANSI 8) | Comments, muted text |
| `{bright_red}` | Bright Red (ANSI 9) | Critical errors |
| `{bright_green}` | Bright Green (ANSI 10) | Emphasized success |
| `{bright_yellow}` | Bright Yellow (ANSI 11) | Important warnings |
| `{bright_blue}` | Bright Blue (ANSI 12) | Active elements |
| `{bright_magenta}` | Bright Magenta (ANSI 13) | Focused items |
| `{bright_cyan}` | Bright Cyan (ANSI 14) | Important data |
| `{bright_white}` | Bright White (ANSI 15) | Bold text, headings |

### Numeric Aliases

You can also use numeric color references:

| Variable | Equivalent |
|----------|------------|
| `{color0}` | `{black}` |
| `{color1}` | `{red}` |
| `{color2}` | `{green}` |
| ... | ... |
| `{color15}` | `{bright_white}` |

## Modifiers

Modifiers transform color values for different formats.

### Standard Hex (Default)

No modifier needed. Returns `#RRGGBB` format.

```
{blue}           → #5294E2
{background}     → #1E1E2E
```

### Strip Hash (`.strip`)

Removes the `#` prefix. Returns `RRGGBB`.

```
{blue.strip}     → 5294E2
{red.strip}      → E06C75
```

Use for applications that don't want the hash prefix.

### RGB Decimal (`.rgb`)

Converts to comma-separated decimal values. Returns `R,G,B`.

```
{blue.rgb}       → 82,148,226
{red.rgb}        → 224,108,117
```

Use for applications expecting decimal RGB values.

### RGBA (`.rgba`)

Converts to CSS rgba() format with optional alpha.

**Default (alpha = 1.0):**
```
{blue.rgba}      → rgba(82,148,226,1.0)
```

**Custom alpha:**
```
{blue.rgba:0.5}  → rgba(82,148,226,0.5)
{black.rgba:0.8} → rgba(30,30,46,0.8)
```

### Hex with 0x Prefix (`.0x`)

Returns hex with `0x` prefix instead of `#`.

```
{blue.0x}        → 0x5294E2
{red.0x}         → 0xE06C75
```

Use for applications expecting C-style hex notation.

### Individual RGB Components (`.r`, `.g`, `.b`)

Extract individual color components as decimal values (0-255).

```
{blue.r}         → 82
{blue.g}         → 148
{blue.b}         → 226
```

Use when you need to manipulate colors component by component.

### Normalized RGB Components (`.rf`, `.gf`, `.bf`)

Extract individual color components as normalized floats (0.0-1.0).

```
{blue.rf}        → 0.321569
{blue.gf}        → 0.580392
{blue.bf}        → 0.886275
```

Use for macOS applications like iTerm2 that expect normalized color values.

### Yaru Icon Theme (`.yaru`)

Maps the color to the closest Ubuntu Yaru icon theme variant.

```
{blue.yaru}      → Yaru-blue
{red.yaru}       → Yaru-red
{green.yaru}     → Yaru-sage
```

## Examples

### Terminal Emulator (Kitty)

```
foreground {foreground}
background {background}
color0 {black}
color1 {red}
```

### CSS/GTK Application (Waybar)

```css
@define-color background {background};
@define-color foreground {foreground};
@define-color accent {blue};
@define-color overlay {black.rgba:0.8};
```

### Application with RGB Values

```ini
[colors]
background = {background.rgb}
accent_r = {blue.r}
accent_g = {blue.g}
accent_b = {blue.b}
```

### Hex Without Hash

```conf
base_color = {background.strip}
text_color = {foreground.strip}
```

### Semi-Transparent Elements

```css
.panel {
    background-color: {background.rgba:0.9};
}
.tooltip {
    background-color: {black.rgba:0.75};
}
```

### macOS iTerm2 Color Profile

```xml
<key>Red Component</key>
<real>{blue.rf}</real>
<key>Green Component</key>
<real>{blue.gf}</real>
<key>Blue Component</key>
<real>{blue.bf}</real>
```
