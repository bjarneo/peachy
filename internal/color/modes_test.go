package color

import "testing"

func TestNextMode(t *testing.T) {
	tests := []struct {
		current ExtractionMode
		want    ExtractionMode
	}{
		{ModeNormal, ModeMaterial},
		{ModeMaterial, ModePastel},
		{ModePastel, ModeMonochromatic},
		{ModeMonochromatic, ModeAnalogous},
		{ModeAnalogous, ModeColorful},
		{ModeColorful, ModeMuted},
		{ModeMuted, ModeBright},
		{ModeBright, ModeNormal},
	}

	for _, tt := range tests {
		t.Run(string(tt.current), func(t *testing.T) {
			got := NextMode(tt.current)
			if got != tt.want {
				t.Errorf("NextMode(%s) = %s, want %s", tt.current, got, tt.want)
			}
		})
	}
}

func TestPrevMode(t *testing.T) {
	tests := []struct {
		current ExtractionMode
		want    ExtractionMode
	}{
		{ModeNormal, ModeBright},
		{ModeMaterial, ModeNormal},
		{ModePastel, ModeMaterial},
		{ModeMonochromatic, ModePastel},
		{ModeAnalogous, ModeMonochromatic},
		{ModeColorful, ModeAnalogous},
		{ModeMuted, ModeColorful},
		{ModeBright, ModeMuted},
	}

	for _, tt := range tests {
		t.Run(string(tt.current), func(t *testing.T) {
			got := PrevMode(tt.current)
			if got != tt.want {
				t.Errorf("PrevMode(%s) = %s, want %s", tt.current, got, tt.want)
			}
		})
	}
}

func TestParseMode(t *testing.T) {
	tests := []struct {
		input string
		want  ExtractionMode
	}{
		{"normal", ModeNormal},
		{"Normal", ModeNormal},
		{"NORMAL", ModeNormal},
		{"material", ModeMaterial},
		{"Material", ModeMaterial},
		{"pastel", ModePastel},
		{"monochromatic", ModeMonochromatic},
		{"mono", ModeMonochromatic},
		{"Mono", ModeMonochromatic},
		{"analogous", ModeAnalogous},
		{"colorful", ModeColorful},
		{"Colorful", ModeColorful},
		{"muted", ModeMuted},
		{"Muted", ModeMuted},
		{"bright", ModeBright},
		{"Bright", ModeBright},
		{"invalid", ModeNormal},
		{"", ModeNormal},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseMode(tt.input)
			if got != tt.want {
				t.Errorf("ParseMode(%q) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestModeNames(t *testing.T) {
	for _, mode := range AllModes {
		name, ok := ModeNames[mode]
		if !ok {
			t.Errorf("Mode %s has no name", mode)
		}
		if name == "" {
			t.Errorf("Mode %s has empty name", mode)
		}
	}
}

func TestModeDescriptions(t *testing.T) {
	for _, mode := range AllModes {
		desc, ok := ModeDescriptions[mode]
		if !ok {
			t.Errorf("Mode %s has no description", mode)
		}
		if desc == "" {
			t.Errorf("Mode %s has empty description", mode)
		}
	}
}

func TestModeCycleComplete(t *testing.T) {
	start := ModeNormal
	current := start

	for i := 0; i < len(AllModes); i++ {
		current = NextMode(current)
	}

	if current != start {
		t.Errorf("After full cycle, got %s, want %s", current, start)
	}
}
