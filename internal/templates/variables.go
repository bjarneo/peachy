package templates

import (
	"fmt"
	"regexp"
	"strings"

	"peachy/internal/color"
)

// BuildVariables creates the variable map from a palette
func BuildVariables(p *color.Palette) map[string]string {
	vars := make(map[string]string)

	// Primary colors
	vars["background"] = p.Background.Hex
	vars["foreground"] = p.Foreground.Hex

	// Semantic color names
	colorNames := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"bright_black", "bright_red", "bright_green", "bright_yellow",
		"bright_blue", "bright_magenta", "bright_cyan", "bright_white",
	}

	for i, name := range colorNames {
		if i < len(p.Colors) {
			vars[name] = p.Colors[i].Hex
		}
	}

	// Also add color0-15 aliases
	for i := 0; i < 16 && i < len(p.Colors); i++ {
		vars[fmt.Sprintf("color%d", i)] = p.Colors[i].Hex
	}

	return vars
}

// ProcessContent processes template content and replaces variables
func ProcessContent(content string, variables map[string]string) string {
	result := content

	for key, value := range variables {
		// Replace {key}
		result = strings.ReplaceAll(result, "{"+key+"}", value)

		// Replace {key.strip} (removes # from hex colors)
		strippedValue := strings.TrimPrefix(value, "#")
		result = strings.ReplaceAll(result, "{"+key+".strip}", strippedValue)

		// Replace {key.rgb} (converts hex to decimal RGB: r,g,b)
		if strings.HasPrefix(value, "#") {
			rgbValue := hexToDecimalRGB(value)
			result = strings.ReplaceAll(result, "{"+key+".rgb}", rgbValue)
		}

		// Replace {key.0x} (hex with 0x prefix instead of #)
		if strings.HasPrefix(value, "#") {
			oxValue := "0x" + strings.TrimPrefix(value, "#")
			result = strings.ReplaceAll(result, "{"+key+".0x}", oxValue)
		}

		// Replace {key.r}, {key.g}, {key.b} (individual RGB components as integers 0-255)
		if strings.HasPrefix(value, "#") {
			r, g, b := hexToRGBComponents(value)
			result = strings.ReplaceAll(result, "{"+key+".r}", fmt.Sprintf("%d", r))
			result = strings.ReplaceAll(result, "{"+key+".g}", fmt.Sprintf("%d", g))
			result = strings.ReplaceAll(result, "{"+key+".b}", fmt.Sprintf("%d", b))

			// Replace {key.rf}, {key.gf}, {key.bf} (normalized RGB as floats 0.0-1.0)
			// Useful for iTerm2 and other macOS apps that expect normalized values
			result = strings.ReplaceAll(result, "{"+key+".rf}", fmt.Sprintf("%.6f", float64(r)/255.0))
			result = strings.ReplaceAll(result, "{"+key+".gf}", fmt.Sprintf("%.6f", float64(g)/255.0))
			result = strings.ReplaceAll(result, "{"+key+".bf}", fmt.Sprintf("%.6f", float64(b)/255.0))
		}

		// Replace {key.rgba:alpha} patterns
		rgbaRegex := regexp.MustCompile(`\{` + regexp.QuoteMeta(key) + `\.rgba(?::(\d*\.?\d+))?\}`)
		result = rgbaRegex.ReplaceAllStringFunc(result, func(match string) string {
			alphaMatch := rgbaRegex.FindStringSubmatch(match)
			alpha := "1.0"
			if len(alphaMatch) > 1 && alphaMatch[1] != "" {
				alpha = alphaMatch[1]
			}
			return hexToRgbaString(value, alpha)
		})

		// Replace {key.yaru} (maps color to Yaru icon theme variant)
		if strings.HasPrefix(value, "#") {
			yaruValue := hexToYaruIconTheme(value)
			result = strings.ReplaceAll(result, "{"+key+".yaru}", yaruValue)
		}
	}

	return result
}

// hexToDecimalRGB converts #RRGGBB to "R,G,B" decimal format
func hexToDecimalRGB(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return hex
	}

	var r, g, b int
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("%d,%d,%d", r, g, b)
}

// hexToRGBComponents extracts individual R, G, B values from hex
func hexToRGBComponents(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}

	var r, g, b int
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return r, g, b
}

// hexToRgbaString converts #RRGGBB to "rgba(R,G,B,A)" CSS format
func hexToRgbaString(hex string, alpha string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return hex
	}

	var r, g, b int
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return fmt.Sprintf("rgba(%d,%d,%d,%s)", r, g, b, alpha)
}

// hexToYaruIconTheme maps a hex color to a Yaru icon theme variant based on hue
func hexToYaruIconTheme(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "Yaru"
	}

	c, err := color.NewColorFromHex(hex)
	if err != nil {
		return "Yaru"
	}
	hue := c.HSL.H

	switch {
	case hue >= 345 || hue < 15:
		return "Yaru-red"
	case hue >= 15 && hue < 30:
		return "Yaru-wartybrown"
	case hue >= 30 && hue < 60:
		return "Yaru-yellow"
	case hue >= 60 && hue < 90:
		return "Yaru-olive"
	case hue >= 90 && hue < 165:
		return "Yaru-sage"
	case hue >= 165 && hue < 195:
		return "Yaru-prussiangreen"
	case hue >= 195 && hue < 255:
		return "Yaru-blue"
	case hue >= 255 && hue < 285:
		return "Yaru-purple"
	default:
		return "Yaru-magenta"
	}
}
