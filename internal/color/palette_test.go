package color

import "testing"

func TestNewPalette(t *testing.T) {
	p := NewPalette()

	// Should have 16 colors
	for i := 0; i < 16; i++ {
		if p.Colors[i].Hex == "" {
			t.Errorf("Color %d has empty hex", i)
		}
	}

	// Background should be color 0
	if p.Background.Hex != p.Colors[0].Hex {
		t.Errorf("Background = %s, want %s", p.Background.Hex, p.Colors[0].Hex)
	}

	// Foreground should be color 7
	if p.Foreground.Hex != p.Colors[7].Hex {
		t.Errorf("Foreground = %s, want %s", p.Foreground.Hex, p.Colors[7].Hex)
	}
}

func TestPaletteGetColor(t *testing.T) {
	p := NewPalette()

	// Valid index
	c := p.GetColor(5)
	if c.Hex == "" {
		t.Error("GetColor(5) returned empty color")
	}

	// Invalid index - should return empty color
	c = p.GetColor(-1)
	if c.Hex != "" {
		t.Errorf("GetColor(-1) = %s, want empty", c.Hex)
	}

	c = p.GetColor(16)
	if c.Hex != "" {
		t.Errorf("GetColor(16) = %s, want empty", c.Hex)
	}
}

func TestPaletteSetColor(t *testing.T) {
	p := NewPalette()
	newColor, _ := NewColorFromHex("#123456")

	// Set color at index 5
	p.SetColor(5, newColor)
	if p.Colors[5].Hex != "#123456" {
		t.Errorf("SetColor(5) = %s, want #123456", p.Colors[5].Hex)
	}

	// Setting color 0 should update background
	p.SetColor(0, newColor)
	if p.Background.Hex != "#123456" {
		t.Errorf("Background after SetColor(0) = %s, want #123456", p.Background.Hex)
	}

	// Setting color 7 should update foreground
	p.SetColor(7, newColor)
	if p.Foreground.Hex != "#123456" {
		t.Errorf("Foreground after SetColor(7) = %s, want #123456", p.Foreground.Hex)
	}

	// Invalid index should be ignored
	p.SetColor(-1, newColor) // Should not panic
	p.SetColor(16, newColor) // Should not panic
}

func TestPaletteGetRoleName(t *testing.T) {
	p := NewPalette()

	tests := []struct {
		index int
		want  string
	}{
		{0, "black"},
		{1, "red"},
		{2, "green"},
		{3, "yellow"},
		{4, "blue"},
		{5, "magenta"},
		{6, "cyan"},
		{7, "white"},
		{8, "bright_black"},
		{15, "bright_white"},
		{-1, ""},
		{16, ""},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := p.GetRoleName(tt.index)
			if got != tt.want {
				t.Errorf("GetRoleName(%d) = %q, want %q", tt.index, got, tt.want)
			}
		})
	}
}

func TestPaletteClone(t *testing.T) {
	p := NewPalette()
	clone := p.Clone()

	// Clone should have same colors
	for i := 0; i < 16; i++ {
		if clone.Colors[i].Hex != p.Colors[i].Hex {
			t.Errorf("Clone color %d = %s, want %s", i, clone.Colors[i].Hex, p.Colors[i].Hex)
		}
	}

	// Modifying clone should not affect original
	newColor, _ := NewColorFromHex("#AABBCC")
	clone.SetColor(5, newColor)
	if p.Colors[5].Hex == "#AABBCC" {
		t.Error("Modifying clone affected original palette")
	}
}

func TestAssignColorsToRoles(t *testing.T) {
	// Create a diverse set of colors
	colors := []Color{}
	for i := 0; i < 20; i++ {
		h := float64(i * 18) // Spread across hues
		c := NewColorFromHSL(h, 70, 50)
		colors = append(colors, c)
	}

	p := AssignColorsToRoles(colors)

	// Should have all 16 colors
	for i := 0; i < 16; i++ {
		if p.Colors[i].Hex == "" {
			t.Errorf("AssignColorsToRoles: color %d is empty", i)
		}
	}

	// Black (0) should be darkest
	// White (7) should be lightest
	if p.Colors[0].Luminance() > p.Colors[7].Luminance() {
		t.Error("Black should be darker than white")
	}
}

func TestAssignColorsToRolesWithFewColors(t *testing.T) {
	// Test with fewer than 16 colors
	colors := []Color{
		NewColorFromHSL(0, 70, 50),   // red
		NewColorFromHSL(120, 70, 50), // green
		NewColorFromHSL(240, 70, 50), // blue
	}

	p := AssignColorsToRoles(colors)

	// Should still have 16 colors (padded)
	for i := 0; i < 16; i++ {
		if p.Colors[i].Hex == "" {
			t.Errorf("AssignColorsToRoles with few colors: color %d is empty", i)
		}
	}
}

func TestPadColors(t *testing.T) {
	colors := []Color{
		NewColorFromHSL(0, 70, 50),
		NewColorFromHSL(120, 70, 50),
	}

	padded := padColors(colors, 10)
	if len(padded) < 10 {
		t.Errorf("padColors length = %d, want at least 10", len(padded))
	}

	// Original colors should be preserved
	if padded[0].Hex != colors[0].Hex {
		t.Error("padColors did not preserve original colors")
	}
}

func TestRoleNames(t *testing.T) {
	// Verify all 16 roles have names
	for i := 0; i < 16; i++ {
		role := ColorRole(i)
		name, ok := RoleNames[role]
		if !ok {
			t.Errorf("Role %d has no name", i)
		}
		if name == "" {
			t.Errorf("Role %d has empty name", i)
		}
	}
}
