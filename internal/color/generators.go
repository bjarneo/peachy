package color

import "math"

// GenerateNormalPalette auto-detects image type and generates appropriate palette
func GenerateNormalPalette(colors []Color, lightMode bool) *Palette {
	if IsMonochromeImage(colors) {
		return generateGrayscalePalette(colors, lightMode)
	}

	if HasLowColorDiversity(colors) {
		return generateSubtleBalancedPalette(colors, lightMode)
	}

	return generateChromaticPalette(colors, lightMode)
}

// GenerateMaterialPalette creates a Material Design-inspired palette
func GenerateMaterialPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()
	usedIndices := make(map[int]bool)

	// Material Design backgrounds
	if lightMode {
		p.Colors[0], _ = NewColorFromHex("#fafafa") // Grey 50
		p.Colors[7], _ = NewColorFromHex("#212121") // Grey 900
	} else {
		p.Colors[0], _ = NewColorFromHex("#121212") // Dark background
		p.Colors[7], _ = NewColorFromHex("#ffffff") // Pure white
	}

	// Find best ANSI color matches from image colors
	for i, targetHue := range ANSIHues {
		matchIndex := FindBestColorMatch(targetHue, colors, usedIndices)
		matchedColor := colors[matchIndex]

		// Apply Material Design refinement
		refinedSaturation := math.Max(matchedColor.HSL.S, 35)
		var refinedLightness float64
		if lightMode {
			refinedLightness = math.Max(35, math.Min(60, matchedColor.HSL.L))
		} else {
			refinedLightness = math.Max(45, math.Min(70, matchedColor.HSL.L))
		}

		p.Colors[i+1] = NewColorFromHSL(matchedColor.HSL.H, refinedSaturation, refinedLightness)
		usedIndices[matchIndex] = true
	}

	// Color 8 (comment): Material Grey
	if lightMode {
		p.Colors[8], _ = NewColorFromHex("#757575")
	} else {
		p.Colors[8], _ = NewColorFromHex("#9e9e9e")
	}

	// Colors 9-14: Brighter versions of 1-6
	for i := 1; i <= 6; i++ {
		c := p.Colors[i]
		brightSaturation := math.Min(100, c.HSL.S+8)
		var brightLightness float64
		if lightMode {
			brightLightness = math.Max(30, c.HSL.L-8)
		} else {
			brightLightness = math.Min(75, c.HSL.L+8)
		}
		p.Colors[i+8] = NewColorFromHSL(c.HSL.H, brightSaturation, brightLightness)
	}

	// Color 15
	if lightMode {
		p.Colors[15], _ = NewColorFromHex("#000000")
	} else {
		p.Colors[15], _ = NewColorFromHex("#ffffff")
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// GeneratePastelPalette creates soft, muted colors from image hues
func GeneratePastelPalette(colors []Color, lightMode bool) *Palette {
	// Start with chromatic palette
	base := generateChromaticPalette(colors, lightMode)

	// Convert all colors to pastel
	for i := 0; i < 16; i++ {
		c := base.Colors[i]

		switch i {
		case 0: // Background
			if lightMode {
				base.Colors[i] = NewColorFromHSL(c.HSL.H, 10, 95)
			} else {
				base.Colors[i] = NewColorFromHSL(c.HSL.H, 15, 20)
			}
		case 7, 15: // Foreground
			if lightMode {
				base.Colors[i] = NewColorFromHSL(c.HSL.H, 25, 35)
			} else {
				base.Colors[i] = NewColorFromHSL(c.HSL.H, 20, 75)
			}
		case 8: // Comment gray
			if lightMode {
				base.Colors[i] = NewColorFromHSL(c.HSL.H, 15, 65)
			} else {
				base.Colors[i] = NewColorFromHSL(c.HSL.H, 12, 45)
			}
		default: // ANSI colors
			pastelSaturation := math.Min(35, c.HSL.S)
			var pastelLightness float64
			if lightMode {
				pastelLightness = 50
			} else {
				pastelLightness = 70
			}
			base.Colors[i] = NewColorFromHSL(c.HSL.H, pastelSaturation, pastelLightness)
		}
	}

	base.Background = base.Colors[0]
	base.Foreground = base.Colors[7]

	return base
}

// GenerateMonochromaticPalette creates single hue variations
func GenerateMonochromaticPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()

	// Find most frequent chromatic color
	baseColor := FindMostFrequentChromatic(colors)
	baseHue := baseColor.HSL.H

	sorted := SortByLightness(colors)
	darkest := sorted[0]
	lightest := sorted[len(sorted)-1]

	// Background and foreground
	if lightMode {
		p.Colors[0] = NewColorFromHSL(baseHue, 8, math.Max(85, lightest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 25, math.Min(30, darkest.HSL.L+10))
	} else {
		p.Colors[0] = NewColorFromHSL(baseHue, 15, math.Min(15, darkest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 10, math.Max(80, lightest.HSL.L-10))
	}

	// ANSI colors 1-6 with base hue and varying saturation/lightness
	saturationLevels := []float64{40, 50, 45, 55, 42, 48}
	var lightnessBase float64
	if lightMode {
		lightnessBase = 45
	} else {
		lightnessBase = 55
	}

	for i := 0; i < 6; i++ {
		lightness := lightnessBase + (float64(i)-2.5)*5
		p.Colors[i+1] = NewColorFromHSL(baseHue, saturationLevels[i], lightness)
	}

	// Color 8: Muted comment
	var color8Lightness float64
	if lightMode {
		color8Lightness = 40
	} else {
		color8Lightness = 65
	}
	p.Colors[8] = NewColorFromHSL(baseHue, 20, color8Lightness)

	// Colors 9-14: Brighter versions
	brightSaturationLevels := []float64{60, 70, 65, 75, 62, 68}
	for i := 0; i < 6; i++ {
		baseLightness := lightnessBase + (float64(i)-2.5)*5
		var adjustment float64
		if lightMode {
			adjustment = -8
		} else {
			adjustment = 8
		}
		lightness := math.Max(0, math.Min(100, baseLightness+adjustment))
		p.Colors[i+9] = NewColorFromHSL(baseHue, brightSaturationLevels[i], lightness)
	}

	// Color 15
	if lightMode {
		p.Colors[15] = NewColorFromHSL(baseHue, 30, math.Min(25, darkest.HSL.L+5))
	} else {
		p.Colors[15] = NewColorFromHSL(baseHue, 15, math.Max(85, lightest.HSL.L))
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// GenerateAnalogousPalette creates adjacent hues on color wheel (±30 degrees)
func GenerateAnalogousPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()

	// Find most saturated color as base
	chromatic := SortBySaturation(colors)
	var baseColor Color
	if len(chromatic) > 0 && chromatic[0].HSL.S > MonochromeSaturationThreshold {
		baseColor = chromatic[0]
	} else if len(colors) > 0 {
		baseColor = colors[0]
	} else {
		baseColor = NewColorFromHSL(200, 50, 50)
	}
	baseHue := baseColor.HSL.H

	sorted := SortByLightness(colors)
	darkest := sorted[0]
	lightest := sorted[len(sorted)-1]

	// Background and foreground
	if lightMode {
		p.Colors[0] = NewColorFromHSL(baseHue, 12, math.Max(90, lightest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 30, math.Min(25, darkest.HSL.L+10))
	} else {
		p.Colors[0] = NewColorFromHSL(baseHue, 18, math.Min(12, darkest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 15, math.Max(85, lightest.HSL.L-10))
	}

	// ANSI colors 1-6 with analogous hues (±30 degrees)
	analogousOffsets := []float64{-30, -20, -10, 10, 20, 30}
	saturationLevels := []float64{45, 50, 48, 52, 47, 50}
	var lightnessBase float64
	if lightMode {
		lightnessBase = 45
	} else {
		lightnessBase = 58
	}

	for i := 0; i < 6; i++ {
		hue := math.Mod(baseHue+analogousOffsets[i]+360, 360)
		var lightness float64
		if i%2 == 0 {
			lightness = lightnessBase - 3
		} else {
			lightness = lightnessBase + 3
		}
		p.Colors[i+1] = NewColorFromHSL(hue, saturationLevels[i], lightness)
	}

	// Color 8 (comment)
	if lightMode {
		p.Colors[8] = NewColorFromHSL(baseHue, 20, 55)
	} else {
		p.Colors[8] = NewColorFromHSL(baseHue, 15, 45)
	}

	// Colors 9-14: Brighter analogous
	for i := 0; i < 6; i++ {
		hue := math.Mod(baseHue+analogousOffsets[i]+360, 360)
		var lightness float64
		if lightMode {
			lightness = 38
		} else {
			lightness = 68
		}
		p.Colors[i+9] = NewColorFromHSL(hue, saturationLevels[i]+8, lightness)
	}

	// Color 15
	if lightMode {
		p.Colors[15] = NewColorFromHSL(baseHue, 20, 20)
	} else {
		p.Colors[15] = NewColorFromHSL(baseHue, 10, 95)
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// generateChromaticPalette generates a vibrant palette from diverse colors
// This matches Aether's generateChromaticPalette function
func generateChromaticPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()
	usedIndices := make(map[int]bool)

	// Find background (darkest or lightest)
	bgColor, bgIndex := FindBackgroundColor(colors, lightMode)
	usedIndices[bgIndex] = true

	// Find foreground (opposite of background)
	fgColor, fgIndex := FindForegroundColor(colors, lightMode, usedIndices)
	usedIndices[fgIndex] = true

	p.Colors[0] = bgColor
	p.Colors[7] = fgColor

	// Find best matches for ANSI colors 1-6
	for i, targetHue := range ANSIHues {
		matchIndex := FindBestColorMatch(targetHue, colors, usedIndices)
		p.Colors[i+1] = colors[matchIndex]
		usedIndices[matchIndex] = true
	}

	// Color 8 (bright black/gray) - brighter in dark mode, darker in light mode
	bg := p.Colors[0]
	var color8Lightness float64
	if IsDarkColor(bg) {
		color8Lightness = math.Min(100, bg.HSL.L+45)
	} else {
		color8Lightness = math.Max(0, bg.HSL.L-40)
	}
	p.Colors[8] = NewColorFromHSL(bg.HSL.H, bg.HSL.S*0.5, color8Lightness)

	// Colors 9-14: Bright versions of 1-6
	for i := 1; i <= 6; i++ {
		p.Colors[i+8] = GenerateBrightVersion(p.Colors[i])
	}

	// Color 15
	p.Colors[15] = GenerateBrightVersion(p.Colors[7])

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// generateGrayscalePalette generates a monochrome/grayscale palette
func generateGrayscalePalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()

	sorted := SortByLightness(colors)
	darkest := sorted[0]
	lightest := sorted[len(sorted)-1]
	baseHue := darkest.HSL.H

	// Background and foreground
	if lightMode {
		p.Colors[0] = lightest
		p.Colors[7] = darkest
	} else {
		p.Colors[0] = darkest
		p.Colors[7] = lightest
	}

	// Generate grayscale shades for ANSI colors 1-6
	const monochromeSaturation = 5.0

	if lightMode {
		startL := darkest.HSL.L + 10
		endL := math.Min(darkest.HSL.L+40, lightest.HSL.L-10)
		step := (endL - startL) / 5

		for i := 1; i <= 6; i++ {
			lightness := startL + float64(i-1)*step
			p.Colors[i] = NewColorFromHSL(baseHue, monochromeSaturation, lightness)
		}
	} else {
		startL := math.Max(darkest.HSL.L+30, lightest.HSL.L-40)
		endL := lightest.HSL.L - 10
		step := (endL - startL) / 5

		for i := 1; i <= 6; i++ {
			lightness := startL + float64(i-1)*step
			p.Colors[i] = NewColorFromHSL(baseHue, monochromeSaturation, lightness)
		}
	}

	// Color 8
	var color8Lightness float64
	if lightMode {
		color8Lightness = math.Max(0, darkest.HSL.L+5)
	} else {
		color8Lightness = math.Min(100, lightest.HSL.L-10)
	}
	p.Colors[8] = NewColorFromHSL(baseHue, monochromeSaturation*0.5, color8Lightness)

	// Colors 9-14
	for i := 1; i <= 6; i++ {
		c := p.Colors[i]
		var adjustment float64
		if lightMode {
			adjustment = -10
		} else {
			adjustment = 10
		}
		newL := math.Max(0, math.Min(100, c.HSL.L+adjustment))
		p.Colors[i+8] = NewColorFromHSL(baseHue, monochromeSaturation, newL)
	}

	// Color 15
	if lightMode {
		p.Colors[15] = NewColorFromHSL(baseHue, 2, math.Max(0, darkest.HSL.L-5))
	} else {
		p.Colors[15] = NewColorFromHSL(baseHue, 2, math.Min(100, lightest.HSL.L+5))
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// generateSubtleBalancedPalette generates a palette for low-diversity images
func generateSubtleBalancedPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()

	sorted := SortByLightness(colors)
	darkest := sorted[0]
	lightest := sorted[len(sorted)-1]

	// Calculate average hue from chromatic colors
	avgHue := CalculateAverageHue(colors)
	if avgHue == 0 {
		avgHue = darkest.HSL.H
	}

	// Background and foreground
	if lightMode {
		p.Colors[0] = lightest
		p.Colors[7] = darkest
	} else {
		p.Colors[0] = darkest
		p.Colors[7] = lightest
	}

	// ANSI colors 1-6 with balanced subtle saturation
	const subtleSaturation = 28.0

	for i := 0; i < 6; i++ {
		lightness := 50 + (float64(i)-2.5)*4 // Vary slightly (42-58%)
		p.Colors[i+1] = NewColorFromHSL(ANSIHues[i], subtleSaturation, lightness)
	}

	// Color 8
	var color8Lightness float64
	if lightMode {
		color8Lightness = math.Max(0, lightest.HSL.L-40)
	} else {
		color8Lightness = math.Min(100, darkest.HSL.L+45)
	}
	p.Colors[8] = NewColorFromHSL(avgHue, subtleSaturation*0.5, color8Lightness)

	// Colors 9-14
	brightSaturation := subtleSaturation + 8
	for i := 0; i < 6; i++ {
		baseLightness := 50 + (float64(i)-2.5)*4
		var adjustment float64
		if lightMode {
			adjustment = -8
		} else {
			adjustment = 8
		}
		lightness := math.Max(0, math.Min(100, baseLightness+adjustment))
		p.Colors[i+9] = NewColorFromHSL(ANSIHues[i], brightSaturation, lightness)
	}

	// Color 15
	if lightMode {
		p.Colors[15] = NewColorFromHSL(avgHue, subtleSaturation*0.3, math.Max(0, darkest.HSL.L-5))
	} else {
		p.Colors[15] = NewColorFromHSL(avgHue, subtleSaturation*0.3, math.Min(100, lightest.HSL.L+5))
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}
