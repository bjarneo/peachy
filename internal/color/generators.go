package color

import "math"

// GenerateNormalPalette auto-detects image type and generates appropriate palette
func GenerateNormalPalette(colors []Color, lightMode bool) *Palette {
	if IsMonochromeImage(colors) {
		return generateMonochromePalette(colors, lightMode)
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
		p.Colors[0], _ = NewColorFromHex("#fafafa")
		p.Colors[7], _ = NewColorFromHex("#212121")
	} else {
		p.Colors[0], _ = NewColorFromHex("#121212")
		p.Colors[7], _ = NewColorFromHex("#ffffff")
	}

	// Find best ANSI color matches from image colors
	for i, targetHue := range ANSIHues {
		matchIndex := FindBestColorMatch(targetHue, colors, usedIndices)
		matchedColor := colors[matchIndex]

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
	return transformChromaticPalette(colors, lightMode, func(i int, c Color, light bool) Color {
		switch i {
		case 0:
			if light {
				return NewColorFromHSL(c.HSL.H, 10, 95)
			}
			return NewColorFromHSL(c.HSL.H, 15, 20)
		case 7, 15:
			if light {
				return NewColorFromHSL(c.HSL.H, 25, 35)
			}
			return NewColorFromHSL(c.HSL.H, 20, 75)
		case 8:
			if light {
				return NewColorFromHSL(c.HSL.H, 15, 65)
			}
			return NewColorFromHSL(c.HSL.H, 12, 45)
		default:
			pastelSat := math.Min(35, c.HSL.S)
			if light {
				return NewColorFromHSL(c.HSL.H, pastelSat, 50)
			}
			return NewColorFromHSL(c.HSL.H, pastelSat, 70)
		}
	})
}

// GenerateColorfulPalette creates a highly saturated, vibrant palette
func GenerateColorfulPalette(colors []Color, lightMode bool) *Palette {
	return transformChromaticPalette(colors, lightMode, func(i int, c Color, light bool) Color {
		switch i {
		case 0:
			if light {
				return NewColorFromHSL(c.HSL.H, 8, 98)
			}
			return NewColorFromHSL(c.HSL.H, 12, 8)
		case 7, 15:
			if light {
				return NewColorFromHSL(c.HSL.H, 15, 10)
			}
			return NewColorFromHSL(c.HSL.H, 10, 95)
		case 8:
			if light {
				return NewColorFromHSL(c.HSL.H, 20, 50)
			}
			return NewColorFromHSL(c.HSL.H, 15, 55)
		default:
			sat := math.Max(75, math.Min(95, c.HSL.S+30))
			if light {
				return NewColorFromHSL(c.HSL.H, sat, math.Max(35, math.Min(55, c.HSL.L)))
			}
			return NewColorFromHSL(c.HSL.H, sat, math.Max(55, math.Min(70, c.HSL.L)))
		}
	})
}

// GenerateMutedPalette creates a desaturated, subdued palette
func GenerateMutedPalette(colors []Color, lightMode bool) *Palette {
	return transformChromaticPalette(colors, lightMode, func(i int, c Color, light bool) Color {
		switch i {
		case 0:
			if light {
				return NewColorFromHSL(c.HSL.H, 5, 95)
			}
			return NewColorFromHSL(c.HSL.H, 8, 15)
		case 7, 15:
			if light {
				return NewColorFromHSL(c.HSL.H, 10, 20)
			}
			return NewColorFromHSL(c.HSL.H, 8, 85)
		case 8:
			if light {
				return NewColorFromHSL(c.HSL.H, 8, 60)
			}
			return NewColorFromHSL(c.HSL.H, 6, 50)
		default:
			sat := math.Max(15, math.Min(35, c.HSL.S*0.5))
			if light {
				return NewColorFromHSL(c.HSL.H, sat, math.Max(40, math.Min(60, c.HSL.L)))
			}
			return NewColorFromHSL(c.HSL.H, sat, math.Max(50, math.Min(65, c.HSL.L)))
		}
	})
}

// GenerateBrightPalette creates a high-lightness palette with punchy colors
func GenerateBrightPalette(colors []Color, lightMode bool) *Palette {
	return transformChromaticPalette(colors, lightMode, func(i int, c Color, light bool) Color {
		switch i {
		case 0:
			if light {
				return NewColorFromHSL(c.HSL.H, 6, 98)
			}
			return NewColorFromHSL(c.HSL.H, 10, 6)
		case 7, 15:
			if light {
				return NewColorFromHSL(c.HSL.H, 12, 15)
			}
			return NewColorFromHSL(c.HSL.H, 8, 98)
		case 8:
			if light {
				return NewColorFromHSL(c.HSL.H, 15, 55)
			}
			return NewColorFromHSL(c.HSL.H, 12, 65)
		default:
			sat := math.Max(45, math.Min(70, c.HSL.S))
			if light {
				return NewColorFromHSL(c.HSL.H, sat, math.Max(45, math.Min(65, c.HSL.L+10)))
			}
			return NewColorFromHSL(c.HSL.H, sat, math.Max(65, math.Min(80, c.HSL.L+15)))
		}
	})
}

// GenerateMonochromaticPalette creates single hue variations
func GenerateMonochromaticPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()

	baseColor := FindMostFrequentChromatic(colors)
	baseHue := baseColor.HSL.H

	sorted := SortByLightness(colors)
	darkest := sorted[0]
	lightest := sorted[len(sorted)-1]

	if lightMode {
		p.Colors[0] = NewColorFromHSL(baseHue, 8, math.Max(85, lightest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 25, math.Min(30, darkest.HSL.L+10))
	} else {
		p.Colors[0] = NewColorFromHSL(baseHue, 15, math.Min(15, darkest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 10, math.Max(80, lightest.HSL.L-10))
	}

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

	if lightMode {
		p.Colors[8] = NewColorFromHSL(baseHue, 20, 40)
	} else {
		p.Colors[8] = NewColorFromHSL(baseHue, 20, 65)
	}

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

	if lightMode {
		p.Colors[0] = NewColorFromHSL(baseHue, 12, math.Max(90, lightest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 30, math.Min(25, darkest.HSL.L+10))
	} else {
		p.Colors[0] = NewColorFromHSL(baseHue, 18, math.Min(12, darkest.HSL.L))
		p.Colors[7] = NewColorFromHSL(baseHue, 15, math.Max(85, lightest.HSL.L-10))
	}

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

	if lightMode {
		p.Colors[8] = NewColorFromHSL(baseHue, 20, 55)
	} else {
		p.Colors[8] = NewColorFromHSL(baseHue, 15, 45)
	}

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
// Uses optimal global ANSI assignment and synthesizes missing hues
func generateChromaticPalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()
	usedIndices := make(map[int]bool)

	// Find background from top dominant colors
	bgColor, bgIndex := FindBackgroundColor(colors, lightMode)
	usedIndices[bgIndex] = true

	// Find foreground (opposite of background)
	fgColor, fgIndex := FindForegroundColor(colors, lightMode, usedIndices)
	usedIndices[fgIndex] = true

	p.Colors[0] = bgColor
	p.Colors[7] = fgColor

	// Use global optimal assignment for ANSI colors 1-6
	assignments := FindOptimalAnsiAssignment(colors, usedIndices)

	// Apply assignments, collecting matched colors for synthesis reference
	var matchedColors []Color
	for i := 0; i < 6; i++ {
		assignment := assignments[i]
		if assignment != nil && assignment.Score < SynthesisScoreThreshold {
			c := colors[assignment.PoolIndex]
			if c.HSL.S >= ANSIMinSaturationForMatch {
				p.Colors[i+1] = c
				matchedColors = append(matchedColors, c)
				usedIndices[assignment.PoolIndex] = true
				continue
			}
		}
		// Mark for synthesis
		p.Colors[i+1] = Color{}
	}

	// Synthesize any missing ANSI colors to match the palette's mood
	for i := 0; i < 6; i++ {
		if p.Colors[i+1].Hex == "" {
			p.Colors[i+1] = SynthesizeAnsiColor(ANSIHues[i], matchedColors)
			matchedColors = append(matchedColors, p.Colors[i+1])
		}
	}

	// Generate color8 (bright black/gray)
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

// detectMonochromeTint detects the dominant tint hue from mostly-gray colors using circular mean
func detectMonochromeTint(colors []Color) (hue float64, hasTint bool) {
	var sinSum, cosSum float64
	var count int

	for _, c := range colors {
		if c.HSL.S > 3 {
			rad := c.HSL.H * math.Pi / 180
			sinSum += math.Sin(rad)
			cosSum += math.Cos(rad)
			count++
		}
	}

	if count == 0 {
		return 0, false
	}

	avgHue := math.Atan2(sinSum/float64(count), cosSum/float64(count)) * 180 / math.Pi
	avgHue = math.Mod(avgHue+360, 360)
	return avgHue, true
}

// applyTint applies tint influence to an ANSI hue based on the image's dominant tone
func applyTint(ansiHue, tintHue float64, hasTint bool) float64 {
	if !hasTint {
		return ansiHue
	}
	hueDiff := math.Mod(tintHue-ansiHue+540, 360) - 180
	return math.Mod(ansiHue+hueDiff*MonochromeTintStrength+360, 360)
}

// generateMonochromePalette generates a monochrome palette with distinguishable hue-tinted colors.
// Uses proper ANSI hue targets at subdued saturation so colors remain
// functional for syntax highlighting while matching the monochrome mood.
func generateMonochromePalette(colors []Color, lightMode bool) *Palette {
	p := NewPalette()

	sorted := SortByLightness(colors)
	darkest := sorted[0]
	lightest := sorted[len(sorted)-1]
	tintHue, hasTint := detectMonochromeTint(colors)

	// Background and foreground from actual image extremes
	if lightMode {
		p.Colors[0] = lightest
		p.Colors[7] = darkest
	} else {
		p.Colors[0] = darkest
		p.Colors[7] = lightest
	}

	// ANSI colors 1-6: proper hues with subdued saturation, tinted toward image tone
	var lightnessBase float64
	if lightMode {
		lightnessBase = 45
	} else {
		lightnessBase = 60
	}
	for i := 0; i < 6; i++ {
		hue := applyTint(ANSIHues[i], tintHue, hasTint)
		lightness := lightnessBase + (float64(i)-2.5)*4
		p.Colors[i+1] = NewColorFromHSL(hue, MonochromeAnsiSaturation, lightness)
	}

	// Color 8: neutral gray for comments
	var color8Lightness float64
	if lightMode {
		color8Lightness = math.Max(0, lightest.HSL.L-35)
	} else {
		color8Lightness = math.Min(100, darkest.HSL.L+40)
	}
	p.Colors[8] = NewColorFromHSL(tintHue, MonochromeSaturation*MonochromeColor8SatFactor, color8Lightness)

	// Colors 9-14: brighter, slightly more saturated versions of 1-6
	for i := 0; i < 6; i++ {
		hue := applyTint(ANSIHues[i], tintHue, hasTint)
		baseLightness := lightnessBase + (float64(i)-2.5)*4
		var adjustment float64
		if lightMode {
			adjustment = -6
		} else {
			adjustment = 6
		}
		lightness := math.Max(0, math.Min(100, baseLightness+adjustment))
		p.Colors[i+9] = NewColorFromHSL(hue, MonochromeAnsiBrightSaturation, lightness)
	}

	// Color 15: near-white or near-black from image
	if lightMode {
		p.Colors[15] = NewColorFromHSL(tintHue, 5, math.Max(0, darkest.HSL.L-5))
	} else {
		p.Colors[15] = NewColorFromHSL(tintHue, 5, math.Min(100, lightest.HSL.L+5))
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

	avgHue := CalculateAverageHue(colors)
	if avgHue == 0 {
		avgHue = darkest.HSL.H
	}

	if lightMode {
		p.Colors[0] = lightest
		p.Colors[7] = darkest
	} else {
		p.Colors[0] = darkest
		p.Colors[7] = lightest
	}

	for i := 0; i < 6; i++ {
		lightness := 50 + (float64(i)-2.5)*4
		p.Colors[i+1] = NewColorFromHSL(ANSIHues[i], SubtlePaletteSaturation, lightness)
	}

	var color8Lightness float64
	if lightMode {
		color8Lightness = math.Max(0, lightest.HSL.L-40)
	} else {
		color8Lightness = math.Min(100, darkest.HSL.L+45)
	}
	p.Colors[8] = NewColorFromHSL(avgHue, SubtlePaletteSaturation*0.5, color8Lightness)

	brightSaturation := SubtlePaletteSaturation + 8
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

	if lightMode {
		p.Colors[15] = NewColorFromHSL(avgHue, SubtlePaletteSaturation*0.3, math.Max(0, darkest.HSL.L-5))
	} else {
		p.Colors[15] = NewColorFromHSL(avgHue, SubtlePaletteSaturation*0.3, math.Min(100, lightest.HSL.L+5))
	}

	p.Background = p.Colors[0]
	p.Foreground = p.Colors[7]

	return p
}

// transformChromaticPalette builds a chromatic palette then transforms each color
func transformChromaticPalette(colors []Color, lightMode bool, transform func(index int, c Color, light bool) Color) *Palette {
	base := generateChromaticPalette(colors, lightMode)

	for i := 0; i < 16; i++ {
		base.Colors[i] = transform(i, base.Colors[i], lightMode)
	}

	base.Background = base.Colors[0]
	base.Foreground = base.Colors[7]

	return base
}
