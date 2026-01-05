package templates

import (
	"testing"
)

func TestProcessContent(t *testing.T) {
	variables := map[string]string{
		"background": "#1E1E2E",
		"foreground": "#CDD6F4",
		"blue":       "#89B4FA",
		"red":        "#F38BA8",
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple replacement",
			input:    "color = {blue}",
			expected: "color = #89B4FA",
		},
		{
			name:     "strip modifier",
			input:    "color = {blue.strip}",
			expected: "color = 89B4FA",
		},
		{
			name:     "rgb modifier",
			input:    "color = {blue.rgb}",
			expected: "color = 137,180,250",
		},
		{
			name:     "rgba modifier default",
			input:    "color = {blue.rgba}",
			expected: "color = rgba(137,180,250,1.0)",
		},
		{
			name:     "rgba modifier with alpha",
			input:    "color = {blue.rgba:0.5}",
			expected: "color = rgba(137,180,250,0.5)",
		},
		{
			name:     "0x modifier",
			input:    "color = {blue.0x}",
			expected: "color = 0x89B4FA",
		},
		{
			name:     "r component",
			input:    "r = {blue.r}",
			expected: "r = 137",
		},
		{
			name:     "g component",
			input:    "g = {blue.g}",
			expected: "g = 180",
		},
		{
			name:     "b component",
			input:    "b = {blue.b}",
			expected: "b = 250",
		},
		{
			name:     "rf component (normalized)",
			input:    "r = {blue.rf}",
			expected: "r = 0.537255",
		},
		{
			name:     "gf component (normalized)",
			input:    "g = {blue.gf}",
			expected: "g = 0.705882",
		},
		{
			name:     "bf component (normalized)",
			input:    "b = {blue.bf}",
			expected: "b = 0.980392",
		},
		{
			name:     "multiple replacements",
			input:    "bg = {background}\nfg = {foreground}",
			expected: "bg = #1E1E2E\nfg = #CDD6F4",
		},
		{
			name:     "no replacement needed",
			input:    "static text",
			expected: "static text",
		},
		{
			name:     "unknown variable preserved",
			input:    "color = {unknown}",
			expected: "color = {unknown}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessContent(tt.input, variables)
			if result != tt.expected {
				t.Errorf("ProcessContent() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHexToDecimalRGB(t *testing.T) {
	tests := []struct {
		hex      string
		expected string
	}{
		{"#89B4FA", "137,180,250"},
		{"89B4FA", "137,180,250"},
		{"#000000", "0,0,0"},
		{"#FFFFFF", "255,255,255"},
		{"#FF0000", "255,0,0"},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			result := hexToDecimalRGB(tt.hex)
			if result != tt.expected {
				t.Errorf("hexToDecimalRGB(%q) = %q, want %q", tt.hex, result, tt.expected)
			}
		})
	}
}

func TestHexToRGBComponents(t *testing.T) {
	tests := []struct {
		hex       string
		expectedR int
		expectedG int
		expectedB int
	}{
		{"#89B4FA", 137, 180, 250},
		{"#000000", 0, 0, 0},
		{"#FFFFFF", 255, 255, 255},
		{"#FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#0000FF", 0, 0, 255},
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			r, g, b := hexToRGBComponents(tt.hex)
			if r != tt.expectedR || g != tt.expectedG || b != tt.expectedB {
				t.Errorf("hexToRGBComponents(%q) = (%d,%d,%d), want (%d,%d,%d)",
					tt.hex, r, g, b, tt.expectedR, tt.expectedG, tt.expectedB)
			}
		})
	}
}
