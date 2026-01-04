package color

import (
	"math"
	"sort"
)

// Constants for color analysis (matching Aether)
const (
	// Saturation threshold below which a color is considered grayscale (0-100)
	MonochromeSaturationThreshold = 15.0

	// Percentage of low-saturation colors needed to classify image as monochrome (0-1)
	MonochromeImageThreshold = 0.7

	// Hue difference (degrees) within which colors are considered similar (0-360)
	SimilarHueRange = 30.0

	// Lightness difference (%) within which colors are considered similar (0-100)
	SimilarLightnessRange = 20.0

	// Low diversity threshold - percentage of similar color pairs (0-1)
	LowDiversityThreshold = 0.6

	// Minimum saturation for chromatic colors (0-100)
	MinChromaticSaturation = 15.0

	// Lightness increase for bright ANSI colors (9-14)
	BrightColorLightnessBoost = 18.0

	// Saturation multiplier for bright ANSI colors
	BrightColorSaturationBoost = 1.1

	// Brightness normalization thresholds
	VeryDarkBackgroundThreshold  = 20.0
	VeryLightBackgroundThreshold = 80.0
	MinLightnessOnDarkBg         = 55.0
	MaxLightnessOnLightBg        = 45.0
	AbsoluteMinLightness         = 25.0
	OutlierLightnessThreshold    = 25.0
	BrightThemeThreshold         = 50.0
	DarkColorThreshold           = 50.0
)

// ANSI color hue targets (in degrees, 0-360)
var ANSIHues = []float64{
	0,   // Red (color 1)
	120, // Green (color 2)
	60,  // Yellow (color 3)
	240, // Blue (color 4)
	300, // Magenta (color 5)
	180, // Cyan (color 6)
}

// IsMonochromeImage detects if the extracted colors are mostly monochrome/grayscale
func IsMonochromeImage(colors []Color) bool {
	lowSaturationCount := 0

	for _, c := range colors {
		if c.HSL.S < MonochromeSaturationThreshold {
			lowSaturationCount++
		}
	}

	return float64(lowSaturationCount)/float64(len(colors)) > MonochromeImageThreshold
}

// HasLowColorDiversity detects if colors are too similar to each other
func HasLowColorDiversity(colors []Color) bool {
	similarCount := 0
	totalComparisons := 0

	for i := 0; i < len(colors); i++ {
		for j := i + 1; j < len(colors); j++ {
			c1 := colors[i]
			c2 := colors[j]

			// Skip grayscale colors
			if c1.HSL.S < MonochromeSaturationThreshold || c2.HSL.S < MonochromeSaturationThreshold {
				continue
			}

			totalComparisons++

			hueDiff := CalculateHueDistance(c1.HSL.H, c2.HSL.H)
			lightnessDiff := math.Abs(c1.HSL.L - c2.HSL.L)

			if hueDiff < SimilarHueRange && lightnessDiff < SimilarLightnessRange {
				similarCount++
			}
		}
	}

	if totalComparisons == 0 {
		return false
	}

	return float64(similarCount)/float64(totalComparisons) > LowDiversityThreshold
}

// CalculateHueDistance calculates circular hue distance between two hues
func CalculateHueDistance(hue1, hue2 float64) float64 {
	diff := math.Abs(hue1 - hue2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// IsDarkColor determines if a color is considered "dark" based on lightness
func IsDarkColor(c Color) bool {
	return c.HSL.L < DarkColorThreshold
}

// FindMostSaturatedColor finds the color with highest saturation
func FindMostSaturatedColor(colors []Color) Color {
	if len(colors) == 0 {
		return NewColorFromHSL(0, 50, 50)
	}

	best := colors[0]
	for _, c := range colors[1:] {
		if c.HSL.S > best.HSL.S {
			best = c
		}
	}
	return best
}

// FindMostFrequentChromatic finds the first color with saturation > threshold
// Colors are assumed to be sorted by frequency (most dominant first)
func FindMostFrequentChromatic(colors []Color) Color {
	for _, c := range colors {
		if c.HSL.S > MonochromeSaturationThreshold {
			return c
		}
	}
	// Fallback to first color
	if len(colors) > 0 {
		return colors[0]
	}
	return NewColorFromHSL(0, 50, 50)
}

// SortByLightness returns colors sorted by lightness (darkest first)
func SortByLightness(colors []Color) []Color {
	sorted := make([]Color, len(colors))
	copy(sorted, colors)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].HSL.L < sorted[j].HSL.L
	})
	return sorted
}

// SortBySaturation returns colors sorted by saturation (most saturated first)
func SortBySaturation(colors []Color) []Color {
	sorted := make([]Color, len(colors))
	copy(sorted, colors)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].HSL.S > sorted[j].HSL.S
	})
	return sorted
}

// FindBackgroundColor finds the best background color for a mode
// Returns the color and its original index
func FindBackgroundColor(colors []Color, lightMode bool) (Color, int) {
	bgIndex := 0
	bgLightness := 101.0
	if lightMode {
		bgLightness = -1.0
	}

	for i, c := range colors {
		if lightMode {
			if c.HSL.L > bgLightness {
				bgLightness = c.HSL.L
				bgIndex = i
			}
		} else {
			if c.HSL.L < bgLightness {
				bgLightness = c.HSL.L
				bgIndex = i
			}
		}
	}

	return colors[bgIndex], bgIndex
}

// FindForegroundColor finds the best foreground color for a mode
// Returns the color and its original index
func FindForegroundColor(colors []Color, lightMode bool, usedIndices map[int]bool) (Color, int) {
	fgIndex := 0
	fgLightness := -1.0
	if lightMode {
		fgLightness = 101.0
	}

	for i, c := range colors {
		if usedIndices[i] {
			continue
		}

		if lightMode {
			if c.HSL.L < fgLightness {
				fgLightness = c.HSL.L
				fgIndex = i
			}
		} else {
			if c.HSL.L > fgLightness {
				fgLightness = c.HSL.L
				fgIndex = i
			}
		}
	}

	return colors[fgIndex], fgIndex
}

// FindBestColorMatch finds the best color for a target ANSI hue
func FindBestColorMatch(targetHue float64, colors []Color, usedIndices map[int]bool) int {
	bestIndex := -1
	bestScore := math.MaxFloat64

	for i, c := range colors {
		if usedIndices[i] {
			continue
		}

		score := calculateColorScore(c, targetHue)
		if score < bestScore {
			bestScore = score
			bestIndex = i
		}
	}

	if bestIndex == -1 {
		// Find first unused
		for i := range colors {
			if !usedIndices[i] {
				return i
			}
		}
		return 0
	}

	return bestIndex
}

// calculateColorScore calculates how well a color matches a target hue
// Lower score is better (matches Aether's implementation)
func calculateColorScore(c Color, targetHue float64) float64 {
	// Hue difference is primary factor - multiply by 3 to prioritize it
	hueDiff := CalculateHueDistance(c.HSL.H, targetHue) * 3

	// Heavily penalize extremely desaturated colors
	saturationPenalty := 0.0
	if c.HSL.S < MinChromaticSaturation {
		saturationPenalty = 50
	}

	// Reward higher saturation: prefer stronger colors when hues are similar
	// Invert saturation (100 - s) so lower score = better (more saturated)
	saturationReward := (100 - c.HSL.S) / 2

	// Very minimal lightness penalties - just avoid extremes
	lightnessPenalty := 0.0
	if c.HSL.L < 20 || c.HSL.L > 85 {
		lightnessPenalty = 10
	}

	return hueDiff + saturationPenalty + saturationReward + lightnessPenalty
}

// GenerateBrightVersion creates a lighter version of a color for bright ANSI slots
func GenerateBrightVersion(c Color) Color {
	newLightness := math.Min(100, c.HSL.L+BrightColorLightnessBoost)
	newSaturation := math.Min(100, c.HSL.S*BrightColorSaturationBoost)
	return NewColorFromHSL(c.HSL.H, newSaturation, newLightness)
}

// GetChromaticColors filters colors to only those with saturation above threshold
func GetChromaticColors(colors []Color) []Color {
	var result []Color
	for _, c := range colors {
		if c.HSL.S > MonochromeSaturationThreshold {
			result = append(result, c)
		}
	}
	return result
}

// CalculateAverageHue calculates the average hue of chromatic colors
func CalculateAverageHue(colors []Color) float64 {
	chromatic := GetChromaticColors(colors)
	if len(chromatic) == 0 {
		return 0
	}

	var sumSin, sumCos float64
	for _, c := range chromatic {
		rad := c.HSL.H * math.Pi / 180
		sumSin += math.Sin(rad)
		sumCos += math.Cos(rad)
	}

	avgRad := math.Atan2(sumSin/float64(len(chromatic)), sumCos/float64(len(chromatic)))
	avgHue := avgRad * 180 / math.Pi
	if avgHue < 0 {
		avgHue += 360
	}
	return avgHue
}

// NormalizeBrightness adjusts ANSI colors to ensure readability
// This matches Aether's normalizeBrightness function
func NormalizeBrightness(p *Palette) {
	bgLightness := p.Colors[0].HSL.L

	isVeryDarkBg := bgLightness < VeryDarkBackgroundThreshold
	isVeryLightBg := bgLightness > VeryLightBackgroundThreshold

	// Analyze colors 1-7
	type colorInfo struct {
		index      int
		lightness  float64
		hue        float64
		saturation float64
	}

	ansiColors := make([]colorInfo, 7)
	for i := 1; i <= 7; i++ {
		c := p.Colors[i]
		ansiColors[i-1] = colorInfo{
			index:      i,
			lightness:  c.HSL.L,
			hue:        c.HSL.H,
			saturation: c.HSL.S,
		}
	}

	avgLightness := 0.0
	for _, c := range ansiColors {
		avgLightness += c.lightness
	}
	avgLightness /= float64(len(ansiColors))

	isBrightTheme := avgLightness > BrightThemeThreshold

	if isVeryDarkBg {
		for _, ci := range ansiColors {
			adjustColorForDarkBackground(p, ci.index, ci.lightness)
		}
		return
	}

	if isVeryLightBg {
		for _, ci := range ansiColors {
			adjustColorForLightBackground(p, ci.index, ci.lightness)
		}
		return
	}

	// Normal background - apply outlier detection
	for _, ci := range ansiColors {
		if math.Abs(ci.lightness-avgLightness) > OutlierLightnessThreshold {
			adjustOutlierColor(p, ci.index, ci.lightness, avgLightness, isBrightTheme)
		}
	}
}

func adjustColorForDarkBackground(p *Palette, index int, lightness float64) {
	if lightness >= MinLightnessOnDarkBg {
		return
	}

	adjustedLightness := MinLightnessOnDarkBg + float64(index)*3
	c := p.Colors[index]
	p.Colors[index] = NewColorFromHSL(c.HSL.H, c.HSL.S, adjustedLightness)

	// Regenerate bright version if applicable
	if index >= 1 && index <= 6 {
		p.Colors[index+8] = GenerateBrightVersion(p.Colors[index])
	}
}

func adjustColorForLightBackground(p *Palette, index int, lightness float64) {
	if lightness <= MaxLightnessOnLightBg {
		return
	}

	adjustedLightness := math.Max(AbsoluteMinLightness, MaxLightnessOnLightBg-float64(index)*2)
	c := p.Colors[index]
	p.Colors[index] = NewColorFromHSL(c.HSL.H, c.HSL.S, adjustedLightness)

	// Regenerate bright version if applicable
	if index >= 1 && index <= 6 {
		p.Colors[index+8] = GenerateBrightVersion(p.Colors[index])
	}
}

func adjustOutlierColor(p *Palette, index int, lightness, avgLightness float64, isBrightTheme bool) {
	isDarkOutlierInBrightTheme := isBrightTheme && lightness < avgLightness-OutlierLightnessThreshold
	isBrightOutlierInDarkTheme := !isBrightTheme && lightness > avgLightness+OutlierLightnessThreshold

	if !isDarkOutlierInBrightTheme && !isBrightOutlierInDarkTheme {
		return
	}

	var adjustedLightness float64
	if isDarkOutlierInBrightTheme {
		adjustedLightness = avgLightness - 10
	} else {
		adjustedLightness = avgLightness + 10
	}

	c := p.Colors[index]
	p.Colors[index] = NewColorFromHSL(c.HSL.H, c.HSL.S, adjustedLightness)

	// Regenerate bright version if applicable
	if index >= 1 && index <= 6 {
		p.Colors[index+8] = GenerateBrightVersion(p.Colors[index])
	}
}
