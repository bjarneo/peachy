package color

import "strings"

// ExtractionMode represents the palette extraction algorithm to use
type ExtractionMode string

const (
	ModeNormal        ExtractionMode = "normal"
	ModeMaterial      ExtractionMode = "material"
	ModePastel        ExtractionMode = "pastel"
	ModeMonochromatic ExtractionMode = "monochromatic"
	ModeAnalogous     ExtractionMode = "analogous"
	ModeColorful      ExtractionMode = "colorful"
	ModeMuted         ExtractionMode = "muted"
	ModeBright        ExtractionMode = "bright"
)

// AllModes contains all available extraction modes
var AllModes = []ExtractionMode{
	ModeNormal,
	ModeMaterial,
	ModePastel,
	ModeMonochromatic,
	ModeAnalogous,
	ModeColorful,
	ModeMuted,
	ModeBright,
}

// ModeNames provides human-readable names for each mode
var ModeNames = map[ExtractionMode]string{
	ModeNormal:        "Normal",
	ModeMaterial:      "Material",
	ModePastel:        "Pastel",
	ModeMonochromatic: "Monochromatic",
	ModeAnalogous:     "Analogous",
	ModeColorful:      "Colorful",
	ModeMuted:         "Muted",
	ModeBright:        "Bright",
}

// ModeDescriptions provides descriptions for each mode
var ModeDescriptions = map[ExtractionMode]string{
	ModeNormal:        "Auto-detects image type and generates appropriate palette",
	ModeMaterial:      "Material Design backgrounds with image colors",
	ModePastel:        "Soft, muted colors with low saturation",
	ModeMonochromatic: "Single hue with varying saturation/lightness",
	ModeAnalogous:     "Adjacent hues on color wheel (±30 degrees)",
	ModeColorful:      "Highly saturated, vibrant colors",
	ModeMuted:         "Desaturated, subdued professional colors",
	ModeBright:        "High lightness with punchy colors",
}

// NextMode returns the next mode in the cycle
func NextMode(current ExtractionMode) ExtractionMode {
	for i, mode := range AllModes {
		if mode == current {
			return AllModes[(i+1)%len(AllModes)]
		}
	}
	return ModeNormal
}

// PrevMode returns the previous mode in the cycle
func PrevMode(current ExtractionMode) ExtractionMode {
	for i, mode := range AllModes {
		if mode == current {
			if i == 0 {
				return AllModes[len(AllModes)-1]
			}
			return AllModes[i-1]
		}
	}
	return ModeNormal
}

// ParseMode parses a string into an ExtractionMode, returning ModeNormal as default
func ParseMode(s string) ExtractionMode {
	switch strings.ToLower(s) {
	case "monochromatic", "mono":
		return ModeMonochromatic
	case "analogous":
		return ModeAnalogous
	case "pastel":
		return ModePastel
	case "material":
		return ModeMaterial
	case "colorful":
		return ModeColorful
	case "muted":
		return ModeMuted
	case "bright":
		return ModeBright
	default:
		return ModeNormal
	}
}
