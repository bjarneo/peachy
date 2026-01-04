package color

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RGB represents a color in RGB color space
type RGB struct {
	R, G, B uint8
}

// HSL represents a color in HSL color space
type HSL struct {
	H float64 // Hue: 0-360
	S float64 // Saturation: 0-100
	L float64 // Lightness: 0-100
}

// Color represents a color with multiple representations
type Color struct {
	Hex string
	RGB RGB
	HSL HSL
}

// NewColorFromHex creates a Color from a hex string (with or without #)
func NewColorFromHex(hex string) (Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return Color{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	rgb := RGB{R: uint8(r), G: uint8(g), B: uint8(b)}
	hsl := RGBToHSL(rgb)

	return Color{
		Hex: "#" + strings.ToUpper(hex),
		RGB: rgb,
		HSL: hsl,
	}, nil
}

// NewColorFromRGB creates a Color from RGB values
func NewColorFromRGB(r, g, b uint8) Color {
	rgb := RGB{R: r, G: g, B: b}
	hsl := RGBToHSL(rgb)
	return Color{
		Hex: RGBToHex(rgb),
		RGB: rgb,
		HSL: hsl,
	}
}

// NewColorFromHSL creates a Color from HSL values
func NewColorFromHSL(h, s, l float64) Color {
	hsl := HSL{H: h, S: s, L: l}
	rgb := HSLToRGB(hsl)
	return Color{
		Hex: RGBToHex(rgb),
		RGB: rgb,
		HSL: hsl,
	}
}

// RGBToHex converts RGB to hex string
func RGBToHex(rgb RGB) string {
	return fmt.Sprintf("#%02X%02X%02X", rgb.R, rgb.G, rgb.B)
}

// RGBToHSL converts RGB to HSL
func RGBToHSL(rgb RGB) HSL {
	r := float64(rgb.R) / 255.0
	g := float64(rgb.G) / 255.0
	b := float64(rgb.B) / 255.0

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	l := (max + min) / 2.0

	var h, s float64

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return HSL{
		H: h * 360,
		S: s * 100,
		L: l * 100,
	}
}

// HSLToRGB converts HSL to RGB
func HSLToRGB(hsl HSL) RGB {
	h := hsl.H / 360.0
	s := hsl.S / 100.0
	l := hsl.L / 100.0

	var r, g, b float64

	if s == 0 {
		r = l
		g = l
		b = l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q
		r = hueToRGBComponent(p, q, h+1.0/3.0)
		g = hueToRGBComponent(p, q, h)
		b = hueToRGBComponent(p, q, h-1.0/3.0)
	}

	return RGB{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
	}
}

// hueToRGBComponent is a helper function for HSL to RGB conversion
func hueToRGBComponent(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

// Luminance calculates the relative luminance of a color (0-1)
func (c Color) Luminance() float64 {
	r := float64(c.RGB.R) / 255.0
	g := float64(c.RGB.G) / 255.0
	b := float64(c.RGB.B) / 255.0

	// Apply gamma correction
	if r <= 0.03928 {
		r = r / 12.92
	} else {
		r = math.Pow((r+0.055)/1.055, 2.4)
	}
	if g <= 0.03928 {
		g = g / 12.92
	} else {
		g = math.Pow((g+0.055)/1.055, 2.4)
	}
	if b <= 0.03928 {
		b = b / 12.92
	} else {
		b = math.Pow((b+0.055)/1.055, 2.4)
	}

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// WithHue returns a new color with adjusted hue
func (c Color) WithHue(h float64) Color {
	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	return NewColorFromHSL(h, c.HSL.S, c.HSL.L)
}

// WithSaturation returns a new color with adjusted saturation
func (c Color) WithSaturation(s float64) Color {
	s = math.Max(0, math.Min(100, s))
	return NewColorFromHSL(c.HSL.H, s, c.HSL.L)
}

// WithLightness returns a new color with adjusted lightness
func (c Color) WithLightness(l float64) Color {
	l = math.Max(0, math.Min(100, l))
	return NewColorFromHSL(c.HSL.H, c.HSL.S, l)
}

// Brighten returns a brighter version of the color
func (c Color) Brighten(amount float64) Color {
	return c.WithLightness(c.HSL.L + amount)
}

// Darken returns a darker version of the color
func (c Color) Darken(amount float64) Color {
	return c.WithLightness(c.HSL.L - amount)
}

// Saturate returns a more saturated version of the color
func (c Color) Saturate(amount float64) Color {
	return c.WithSaturation(c.HSL.S + amount)
}

// Desaturate returns a less saturated version of the color
func (c Color) Desaturate(amount float64) Color {
	return c.WithSaturation(c.HSL.S - amount)
}

// String returns the hex representation
func (c Color) String() string {
	return c.Hex
}
