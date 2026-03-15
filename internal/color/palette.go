package color

import (
	"sort"
)

// ColorRole represents the semantic role of a color in the ANSI palette
type ColorRole int

const (
	RoleBlack ColorRole = iota
	RoleRed
	RoleGreen
	RoleYellow
	RoleBlue
	RoleMagenta
	RoleCyan
	RoleWhite
	RoleBrightBlack
	RoleBrightRed
	RoleBrightGreen
	RoleBrightYellow
	RoleBrightBlue
	RoleBrightMagenta
	RoleBrightCyan
	RoleBrightWhite
)

// RoleNames maps color roles to human-readable names
var RoleNames = map[ColorRole]string{
	RoleBlack:         "black",
	RoleRed:           "red",
	RoleGreen:         "green",
	RoleYellow:        "yellow",
	RoleBlue:          "blue",
	RoleMagenta:       "magenta",
	RoleCyan:          "cyan",
	RoleWhite:         "white",
	RoleBrightBlack:   "bright_black",
	RoleBrightRed:     "bright_red",
	RoleBrightGreen:   "bright_green",
	RoleBrightYellow:  "bright_yellow",
	RoleBrightBlue:    "bright_blue",
	RoleBrightMagenta: "bright_magenta",
	RoleBrightCyan:    "bright_cyan",
	RoleBrightWhite:   "bright_white",
}

// HueRanges defines the hue ranges for each color category
// Hue is in degrees (0-360)
var HueRanges = map[ColorRole][2]float64{
	RoleRed:     {345, 15}, // Wraps around 0
	RoleYellow:  {45, 75},
	RoleGreen:   {75, 165},
	RoleCyan:    {165, 195},
	RoleBlue:    {195, 255},
	RoleMagenta: {255, 345},
}

// Palette represents a 16-color ANSI palette
type Palette struct {
	Background Color
	Foreground Color
	Colors     [16]Color
}

// NewPalette creates an empty palette with default dark colors
func NewPalette() *Palette {
	p := &Palette{}

	// Default dark theme
	defaults := []string{
		"#1a1a2e", // 0: black
		"#f54242", // 1: red
		"#42f554", // 2: green
		"#f5f542", // 3: yellow
		"#4287f5", // 4: blue
		"#f542f5", // 5: magenta
		"#42f5f5", // 6: cyan
		"#e0e0e0", // 7: white
		"#4a4a5e", // 8: bright black
		"#f76b6b", // 9: bright red
		"#6bf76b", // 10: bright green
		"#f7f76b", // 11: bright yellow
		"#6ba3f7", // 12: bright blue
		"#f76bf7", // 13: bright magenta
		"#6bf7f7", // 14: bright cyan
		"#ffffff", // 15: bright white
	}

	for i, hex := range defaults {
		c, _ := NewColorFromHex(hex)
		p.Colors[i] = c
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// GetColor returns the color at the given index
func (p *Palette) GetColor(index int) Color {
	if index < 0 || index >= 16 {
		return Color{}
	}
	return p.Colors[index]
}

// SetColor sets the color at the given index
func (p *Palette) SetColor(index int, c Color) {
	if index >= 0 && index < 16 {
		p.Colors[index] = c

		// Update background/foreground if needed
		switch index {
		case 0:
			p.Background = c
		case 7:
			p.Foreground = c
		}
	}
}

// GetRoleName returns the role name for a color index
func (p *Palette) GetRoleName(index int) string {
	if index < 0 || index >= 16 {
		return ""
	}
	return RoleNames[ColorRole(index)]
}

// AssignColorsToRoles takes a slice of colors and assigns them to ANSI roles
// based on their hue and luminance
func AssignColorsToRoles(colors []Color) *Palette {
	if len(colors) < 16 {
		// Pad with generated colors if we don't have enough
		colors = padColors(colors, 16)
	}

	p := NewPalette()

	// Sort by luminance
	sorted := make([]Color, len(colors))
	copy(sorted, colors)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Luminance() < sorted[j].Luminance()
	})

	// Assign darkest to black (0), lightest to white (7)
	p.Colors[0] = sorted[0]
	p.Colors[7] = sorted[len(sorted)-1]

	// Assign second darkest to bright black (8), second lightest to bright white (15)
	if len(sorted) > 1 {
		p.Colors[8] = sorted[1]
		p.Colors[15] = sorted[len(sorted)-2]
	}

	// Remaining colors to assign based on hue
	remaining := make([]Color, 0)
	for i := 2; i < len(sorted)-2; i++ {
		remaining = append(remaining, sorted[i])
	}

	// Assign colors by hue category
	hueAssignments := []struct {
		normal int
		bright int
	}{
		{1, 9},  // red
		{2, 10}, // green
		{3, 11}, // yellow
		{4, 12}, // blue
		{5, 13}, // magenta
		{6, 14}, // cyan
	}

	for _, assignment := range hueAssignments {
		roleNormal := ColorRole(assignment.normal)
		normalColor, brightColor := findColorsForRole(remaining, roleNormal)
		p.Colors[assignment.normal] = normalColor
		p.Colors[assignment.bright] = brightColor
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// findColorsForRole finds the best normal and bright color for a given hue role
func findColorsForRole(colors []Color, role ColorRole) (normal, bright Color) {
	hueRange, ok := HueRanges[role]
	if !ok {
		// Fall back to default colors if role not in HueRanges
		normal, _ = NewColorFromHex("#888888")
		bright, _ = NewColorFromHex("#aaaaaa")
		return
	}

	var candidates []Color

	for _, c := range colors {
		h := c.HSL.H
		// Handle wrap-around for red
		if hueRange[0] > hueRange[1] {
			if h >= hueRange[0] || h <= hueRange[1] {
				candidates = append(candidates, c)
			}
		} else {
			if h >= hueRange[0] && h <= hueRange[1] {
				candidates = append(candidates, c)
			}
		}
	}

	if len(candidates) == 0 {
		// Generate a color in the target hue
		targetHue := (hueRange[0] + hueRange[1]) / 2
		if hueRange[0] > hueRange[1] {
			targetHue = hueRange[0] - 15
		}
		normal = NewColorFromHSL(targetHue, 70, 50)
		bright = NewColorFromHSL(targetHue, 70, 70)
		return
	}

	// Sort by saturation (most saturated first)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].HSL.S > candidates[j].HSL.S
	})

	normal = candidates[0]
	if len(candidates) > 1 {
		bright = candidates[1]
	} else {
		bright = normal.Brighten(20)
	}

	// Ensure bright is actually brighter
	if bright.HSL.L <= normal.HSL.L {
		bright = normal.Brighten(15)
	}

	return
}

// padColors extends the color slice to the target length by generating variations
func padColors(colors []Color, target int) []Color {
	if len(colors) >= target {
		return colors
	}

	result := make([]Color, len(colors))
	copy(result, colors)

	for len(result) < target {
		// Generate variations of existing colors
		for _, c := range colors {
			if len(result) >= target {
				break
			}
			// Add a brighter version
			result = append(result, c.Brighten(20))
		}
		for _, c := range colors {
			if len(result) >= target {
				break
			}
			// Add a different hue variation
			result = append(result, c.WithHue(c.HSL.H+30))
		}
	}

	return result
}

// Clone creates a copy of the palette
func (p *Palette) Clone() *Palette {
	clone := &Palette{
		Background: p.Background,
		Foreground: p.Foreground,
	}
	for i := 0; i < 16; i++ {
		clone.Colors[i] = p.Colors[i]
	}
	return clone
}
