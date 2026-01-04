# Peachy Extraction Modes Implementation Plan

Based on Aether's `src/utils/imagemagick-color-extraction.js` and `src/services/color-harmony.js`.

## Extraction Modes

| Mode | Description | Default |
|------|-------------|---------|
| `normal` | Auto-detects image type and generates appropriate palette | Yes |
| `material` | Material Design backgrounds with image colors |  |
| `pastel` | Low saturation, high lightness (soft, muted) |  |
| `monochromatic` | Single hue with varying saturation/lightness |  |
| `analogous` | Adjacent hues on color wheel (±30 degrees) |  |

## Mode Details

### Normal (Default)

Auto-detects image characteristics:
1. **Monochrome detection**: If >70% of colors have saturation <15%, generate grayscale palette
2. **Low diversity detection**: If >60% of color pairs have similar hue (<30°) and lightness (<20%), generate subtle balanced palette
3. **Otherwise**: Generate chromatic palette from extracted colors

```go
func (e *Extractor) ExtractPalette(imagePath string, mode ExtractionMode) (*Palette, error)
```

### Material

Uses Material Design backgrounds with image colors:
- Light mode background: `#fafafa` (Grey 50)
- Dark mode background: `#121212`
- Light mode foreground: `#212121` (Grey 900)
- Dark mode foreground: `#ffffff`
- Color 8 (comments): `#757575` (light) / `#9e9e9e` (dark)
- ANSI colors 1-6: Best matches from image with refined saturation (min 35%) and lightness normalization

### Pastel

Soft, muted colors:
- Background: Very light/desaturated (`hsl(h, 10, 95)` light / `hsl(h, 15, 20)` dark)
- Foreground: Darker but pastel (`hsl(h, 25, 35)` light / `hsl(h, 20, 75)` dark)
- ANSI colors: Max saturation 35%, lightness ~50 (light) / ~70 (dark)

### Monochromatic

Single hue variations:
1. Find most frequent color with saturation >15%
2. Use its hue for all colors
3. Vary saturation: 40-55% for normal, 60-75% for bright
4. Vary lightness: ~45 (light mode) / ~55 (dark mode) base, ±5 per color

### Analogous

Adjacent hues (±30 degrees from base):
1. Find most saturated color as base hue
2. Generate offsets: `[-30, -20, -10, 10, 20, 30]` degrees
3. Apply to ANSI colors 1-6
4. Background/foreground use base hue with adjusted saturation/lightness

## Implementation

### 1. Add ExtractionMode type

```go
// internal/color/extractor.go

type ExtractionMode string

const (
    ModeNormal        ExtractionMode = "normal"
    ModeMaterial      ExtractionMode = "material"
    ModePastel        ExtractionMode = "pastel"
    ModeMonochromatic ExtractionMode = "monochromatic"
    ModeAnalogous     ExtractionMode = "analogous"
)

var AllModes = []ExtractionMode{
    ModeNormal,
    ModeMaterial,
    ModePastel,
    ModeMonochromatic,
    ModeAnalogous,
}
```

### 2. Add palette generators

```go
// internal/color/generators.go

func generateNormalPalette(colors []Color, lightMode bool) *Palette
func generateMaterialPalette(colors []Color, lightMode bool) *Palette
func generatePastelPalette(colors []Color, lightMode bool) *Palette
func generateMonochromaticPalette(colors []Color, lightMode bool) *Palette
func generateAnalogousPalette(colors []Color, lightMode bool) *Palette
```

### 3. Helper functions needed

```go
// internal/color/analysis.go

func isMonochromeImage(colors []Color) bool
func hasLowColorDiversity(colors []Color) bool
func findMostSaturatedColor(colors []Color) Color
func calculateHueDistance(hue1, hue2 float64) float64
func sortByLightness(colors []Color) []Color
func sortBySaturation(colors []Color) []Color
```

### 4. Update ExtractPalette signature

```go
func (e *Extractor) ExtractPalette(imagePath string, mode ExtractionMode, lightMode bool) (*Palette, error) {
    colors, err := e.ExtractColors(imagePath)
    if err != nil {
        return nil, err
    }

    switch mode {
    case ModeMaterial:
        return generateMaterialPalette(colors, lightMode), nil
    case ModePastel:
        return generatePastelPalette(colors, lightMode), nil
    case ModeMonochromatic:
        return generateMonochromaticPalette(colors, lightMode), nil
    case ModeAnalogous:
        return generateAnalogousPalette(colors, lightMode), nil
    default:
        return generateNormalPalette(colors, lightMode), nil
    }
}
```

### 5. UI integration

Add mode selector to main view:
- Key binding: `m` to cycle modes, or `M` to open mode picker
- Display current mode in status bar
- Re-extract when mode changes

```go
// internal/ui/model.go

type Model struct {
    // ...
    extractionMode color.ExtractionMode
    lightMode      bool
}
```

## Constants (from Aether)

```go
const (
    // Saturation threshold for grayscale detection
    MonochromeSaturationThreshold = 15.0

    // Percentage of low-saturation colors for monochrome image
    MonochromeImageThreshold = 0.7

    // Hue difference for similar colors
    SimilarHueRange = 30.0

    // Lightness difference for similar colors
    SimilarLightnessRange = 20.0

    // Low diversity threshold
    LowDiversityThreshold = 0.6

    // Minimum saturation for chromatic colors
    MinChromaticSaturation = 15.0

    // Bright color lightness boost
    BrightColorLightnessBoost = 18.0
)
```

## ANSI Color Hue Targets

```go
var ANSIHues = map[int]float64{
    1: 0,    // Red
    2: 120,  // Green
    3: 60,   // Yellow
    4: 240,  // Blue
    5: 300,  // Magenta
    6: 180,  // Cyan
}
```

## Files to create/modify

1. `internal/color/modes.go` - ExtractionMode type and constants
2. `internal/color/generators.go` - Palette generator functions
3. `internal/color/analysis.go` - Image analysis helpers
4. `internal/color/extractor.go` - Update ExtractPalette signature
5. `internal/ui/model.go` - Add mode state
6. `internal/ui/views.go` - Display current mode
7. `internal/ui/keys.go` - Add mode keybindings
