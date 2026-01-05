# Extraction Modes

Peachy supports multiple extraction algorithms for different aesthetic styles. Press `m` in the TUI to cycle through modes.

## Available Modes

| Mode | Description |
|------|-------------|
| **Normal** | Auto-detects image type. Generates grayscale palette for monochrome images, subtle balanced palette for low-diversity images, or chromatic palette for colorful images. (Default) |
| **Material** | Uses Material Design backgrounds (#fafafa light / #121212 dark) with refined image colors. Clean, professional aesthetic. |
| **Pastel** | Soft, muted colors with low saturation and high lightness. Great for easy-on-the-eyes themes. |
| **Monochromatic** | Single hue derived from the most frequent chromatic color. Creates cohesive, focused themes. |
| **Analogous** | Adjacent hues on the color wheel (plus/minus 30 degrees). Creates harmonious, visually pleasing palettes. |

## Light/Dark Mode

Press `t` to toggle between light and dark mode generation. This affects:

- Background/foreground color selection
- Lightness levels for all palette colors
- Comment color brightness
