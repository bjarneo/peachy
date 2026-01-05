package color

import (
	"testing"
)

func TestNewColorFromHex(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantErr bool
	}{
		{"valid with hash", "#FF5500", 255, 85, 0, false},
		{"valid without hash", "FF5500", 255, 85, 0, false},
		{"black", "#000000", 0, 0, 0, false},
		{"white", "#FFFFFF", 255, 255, 255, false},
		{"lowercase", "#aabbcc", 170, 187, 204, false},
		{"invalid length", "#FFF", 0, 0, 0, true},
		{"invalid chars", "#GGGGGG", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewColorFromHex(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewColorFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if c.RGB.R != tt.wantR || c.RGB.G != tt.wantG || c.RGB.B != tt.wantB {
					t.Errorf("NewColorFromHex() = RGB(%d,%d,%d), want RGB(%d,%d,%d)",
						c.RGB.R, c.RGB.G, c.RGB.B, tt.wantR, tt.wantG, tt.wantB)
				}
			}
		})
	}
}

func TestNewColorFromRGB(t *testing.T) {
	c := NewColorFromRGB(255, 128, 64)

	if c.RGB.R != 255 || c.RGB.G != 128 || c.RGB.B != 64 {
		t.Errorf("NewColorFromRGB() RGB = (%d,%d,%d), want (255,128,64)",
			c.RGB.R, c.RGB.G, c.RGB.B)
	}

	if c.Hex != "#FF8040" {
		t.Errorf("NewColorFromRGB() Hex = %s, want #FF8040", c.Hex)
	}
}

func TestNewColorFromHSL(t *testing.T) {
	// Red
	c := NewColorFromHSL(0, 100, 50)
	if c.RGB.R != 255 || c.RGB.G != 0 || c.RGB.B != 0 {
		t.Errorf("NewColorFromHSL(0,100,50) = RGB(%d,%d,%d), want RGB(255,0,0)",
			c.RGB.R, c.RGB.G, c.RGB.B)
	}

	// Green
	c = NewColorFromHSL(120, 100, 50)
	if c.RGB.R != 0 || c.RGB.G != 255 || c.RGB.B != 0 {
		t.Errorf("NewColorFromHSL(120,100,50) = RGB(%d,%d,%d), want RGB(0,255,0)",
			c.RGB.R, c.RGB.G, c.RGB.B)
	}

	// Blue
	c = NewColorFromHSL(240, 100, 50)
	if c.RGB.R != 0 || c.RGB.G != 0 || c.RGB.B != 255 {
		t.Errorf("NewColorFromHSL(240,100,50) = RGB(%d,%d,%d), want RGB(0,0,255)",
			c.RGB.R, c.RGB.G, c.RGB.B)
	}
}

func TestRGBToHSLRoundtrip(t *testing.T) {
	tests := []struct {
		r, g, b uint8
	}{
		{255, 0, 0},     // Red
		{0, 255, 0},     // Green
		{0, 0, 255},     // Blue
		{255, 255, 0},   // Yellow
		{0, 255, 255},   // Cyan
		{255, 0, 255},   // Magenta
		{128, 128, 128}, // Gray
		{0, 0, 0},       // Black
		{255, 255, 255}, // White
		{100, 150, 200}, // Random
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			rgb := RGB{R: tt.r, G: tt.g, B: tt.b}
			hsl := RGBToHSL(rgb)
			rgbBack := HSLToRGB(hsl)

			// Allow small rounding errors
			if abs(int(rgb.R)-int(rgbBack.R)) > 1 ||
				abs(int(rgb.G)-int(rgbBack.G)) > 1 ||
				abs(int(rgb.B)-int(rgbBack.B)) > 1 {
				t.Errorf("Roundtrip failed: RGB(%d,%d,%d) -> HSL(%.1f,%.1f,%.1f) -> RGB(%d,%d,%d)",
					rgb.R, rgb.G, rgb.B, hsl.H, hsl.S, hsl.L, rgbBack.R, rgbBack.G, rgbBack.B)
			}
		})
	}
}

func TestRGBToHex(t *testing.T) {
	tests := []struct {
		rgb  RGB
		want string
	}{
		{RGB{255, 255, 255}, "#FFFFFF"},
		{RGB{0, 0, 0}, "#000000"},
		{RGB{255, 0, 0}, "#FF0000"},
		{RGB{0, 255, 0}, "#00FF00"},
		{RGB{0, 0, 255}, "#0000FF"},
		{RGB{170, 187, 204}, "#AABBCC"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := RGBToHex(tt.rgb)
			if got != tt.want {
				t.Errorf("RGBToHex(%v) = %s, want %s", tt.rgb, got, tt.want)
			}
		})
	}
}

func TestColorLuminance(t *testing.T) {
	// White should have high luminance
	white, _ := NewColorFromHex("#FFFFFF")
	if white.Luminance() < 0.9 {
		t.Errorf("White luminance = %f, want > 0.9", white.Luminance())
	}

	// Black should have low luminance
	black, _ := NewColorFromHex("#000000")
	if black.Luminance() > 0.1 {
		t.Errorf("Black luminance = %f, want < 0.1", black.Luminance())
	}

	// Gray should have medium luminance
	gray, _ := NewColorFromHex("#808080")
	if gray.Luminance() < 0.1 || gray.Luminance() > 0.5 {
		t.Errorf("Gray luminance = %f, want between 0.1 and 0.5", gray.Luminance())
	}
}

func TestColorWithHue(t *testing.T) {
	red, _ := NewColorFromHex("#FF0000")
	green := red.WithHue(120)

	if green.RGB.R > 10 || green.RGB.G < 240 || green.RGB.B > 10 {
		t.Errorf("WithHue(120) from red = RGB(%d,%d,%d), want approximately green",
			green.RGB.R, green.RGB.G, green.RGB.B)
	}
}

func TestColorWithSaturation(t *testing.T) {
	red, _ := NewColorFromHex("#FF0000")
	desaturated := red.WithSaturation(0)

	// Desaturated red should be gray
	if desaturated.RGB.R != desaturated.RGB.G || desaturated.RGB.G != desaturated.RGB.B {
		t.Errorf("WithSaturation(0) = RGB(%d,%d,%d), want gray",
			desaturated.RGB.R, desaturated.RGB.G, desaturated.RGB.B)
	}
}

func TestColorWithLightness(t *testing.T) {
	red, _ := NewColorFromHex("#FF0000")

	// Very dark
	dark := red.WithLightness(10)
	if dark.RGB.R > 60 {
		t.Errorf("WithLightness(10) R = %d, want < 60", dark.RGB.R)
	}

	// Very light
	light := red.WithLightness(90)
	if light.RGB.R < 200 {
		t.Errorf("WithLightness(90) R = %d, want > 200", light.RGB.R)
	}
}

func TestColorBrightenDarken(t *testing.T) {
	color, _ := NewColorFromHex("#808080")

	brighter := color.Brighten(20)
	if brighter.HSL.L <= color.HSL.L {
		t.Errorf("Brighten() L = %f, want > %f", brighter.HSL.L, color.HSL.L)
	}

	darker := color.Darken(20)
	if darker.HSL.L >= color.HSL.L {
		t.Errorf("Darken() L = %f, want < %f", darker.HSL.L, color.HSL.L)
	}
}

func TestColorString(t *testing.T) {
	c, _ := NewColorFromHex("#AABBCC")
	if c.String() != "#AABBCC" {
		t.Errorf("String() = %s, want #AABBCC", c.String())
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
