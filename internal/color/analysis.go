package color

import (
	"math"
	"sort"
)

// Constants for color analysis
const (
	// Saturation threshold below which a color is considered grayscale (0-100)
	MonochromeSaturationThreshold = 15.0

	// Percentage of low-saturation colors needed to classify image as monochrome (0-1)
	MonochromeImageThreshold = 0.7

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

	// Color scoring thresholds
	TooDarkThreshold   = 20.0
	TooBrightThreshold = 85.0

	// Score threshold above which a synthesized color is used instead of a poor match
	SynthesisScoreThreshold = 180.0

	// Minimum saturation for a color to be considered a valid ANSI match
	ANSIMinSaturationForMatch = 12.0

	// Monochrome palette settings
	MonochromeSaturation           = 5.0
	MonochromeAnsiSaturation       = 30.0
	MonochromeAnsiBrightSaturation = 40.0
	MonochromeTintStrength         = 0.15
	MonochromeColor8SatFactor      = 0.5

	// Subtle balanced palette saturation
	SubtlePaletteSaturation = 28.0
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

// HasLowColorDiversity detects if colors lack hue diversity using O(n) hue bucketing
func HasLowColorDiversity(colors []Color) bool {
	// Sample first 16 colors for faster analysis
	sampleSize := len(colors)
	if sampleSize > 16 {
		sampleSize = 16
	}

	// Filter to chromatic colors and bucket by hue
	const hueBucketSize = 30.0
	hueBuckets := make([]int, 12)
	chromaticCount := 0

	for i := 0; i < sampleSize; i++ {
		c := colors[i]
		if c.HSL.S >= MonochromeSaturationThreshold {
			chromaticCount++
			bucket := int(c.HSL.H/hueBucketSize) % 12
			hueBuckets[bucket]++
		}
	}

	// Need at least 3 chromatic colors to determine diversity
	if chromaticCount < 3 {
		return false
	}

	// Count occupied hue buckets
	occupiedBuckets := 0
	for _, count := range hueBuckets {
		if count > 0 {
			occupiedBuckets++
		}
	}

	// If colors span fewer than 3 hue buckets, they lack diversity
	return occupiedBuckets < 3
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
func FindBackgroundColor(colors []Color, lightMode bool) (Color, int) {
	// Search within top 12 dominant colors for better representation
	topCount := len(colors)
	if topCount > 12 {
		topCount = 12
	}
	topColors := colors[:topCount]

	bgIndex := 0
	bgLightness := 101.0
	if lightMode {
		bgLightness = -1.0
	}

	for i, c := range topColors {
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

	return topColors[bgIndex], bgIndex
}

// FindForegroundColor finds the best foreground color for a mode
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

// calculateColorScore calculates how well a color matches a target ANSI hue
// Balances hue accuracy, saturation preference, and lightness suitability
// Lower score = better match
func calculateColorScore(c Color, targetHue float64) float64 {
	// Hue accuracy - primary factor
	hueScore := CalculateHueDistance(c.HSL.H, targetHue) * 2.5

	// Saturation preference - strongly prefer chromatic colors
	var satScore float64
	if c.HSL.S < ANSIMinSaturationForMatch {
		satScore = 80
	} else if c.HSL.S < 20 {
		satScore = 40
	} else if c.HSL.S < 30 {
		satScore = 15
	} else {
		satScore = math.Max(0, (50-c.HSL.S)*0.3)
	}

	// Lightness suitability - prefer mid-range, penalize extremes
	var lightnessScore float64
	if c.HSL.L < TooDarkThreshold {
		lightnessScore = (TooDarkThreshold - c.HSL.L) * 2.5
	} else if c.HSL.L > TooBrightThreshold {
		lightnessScore = (c.HSL.L - TooBrightThreshold) * 2
	} else {
		lightnessScore = math.Abs(c.HSL.L-55) * 0.2
	}

	return hueScore + satScore + lightnessScore
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
		for i := range colors {
			if !usedIndices[i] {
				return i
			}
		}
		return 0
	}

	return bestIndex
}

// ansiAssignment represents a color pool index and its match quality score
type ansiAssignment struct {
	PoolIndex int
	Score     float64
}

// FindOptimalAnsiAssignment finds optimal ANSI color assignments using global greedy matching.
// Instead of assigning colors sequentially (red first, then green, etc.),
// this finds the globally best (ANSI slot, color) pair at each step,
// preventing earlier slots from stealing good matches from later ones.
func FindOptimalAnsiAssignment(colorPool []Color, usedIndices map[int]bool) [6]*ansiAssignment {
	type candidate struct {
		poolIndex int
		score     float64
	}

	// Pre-compute and sort scores for all (ANSI slot, color) pairs
	allScores := make([][]candidate, 6)
	for a, targetHue := range ANSIHues {
		candidates := make([]candidate, 0, len(colorPool))
		for i, c := range colorPool {
			if usedIndices[i] {
				continue
			}
			candidates = append(candidates, candidate{
				poolIndex: i,
				score:     calculateColorScore(c, targetHue),
			})
		}
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].score < candidates[j].score
		})
		allScores[a] = candidates
	}

	var assignments [6]*ansiAssignment
	assignedPool := make(map[int]bool)
	for k, v := range usedIndices {
		if v {
			assignedPool[k] = true
		}
	}

	// Iteratively assign the globally best pair
	for round := 0; round < 6; round++ {
		bestAnsi := -1
		bestPoolIndex := -1
		bestScore := math.MaxFloat64

		for a := 0; a < 6; a++ {
			if assignments[a] != nil {
				continue
			}
			// Find best unassigned candidate for this slot
			for _, cand := range allScores[a] {
				if assignedPool[cand.poolIndex] {
					continue
				}
				if cand.score < bestScore {
					bestScore = cand.score
					bestAnsi = a
					bestPoolIndex = cand.poolIndex
				}
				break // First unassigned is best (list is sorted)
			}
		}

		if bestAnsi == -1 {
			break
		}
		assignments[bestAnsi] = &ansiAssignment{PoolIndex: bestPoolIndex, Score: bestScore}
		assignedPool[bestPoolIndex] = true
	}

	return assignments
}

// SynthesizeAnsiColor creates an ANSI color when no good match exists in the image.
// Uses the average saturation and lightness of already-assigned colors
// to create a color that fits the palette's visual mood.
func SynthesizeAnsiColor(targetHue float64, existingColors []Color) Color {
	var totalS, totalL float64
	var count int

	for _, c := range existingColors {
		if c.HSL.S >= ANSIMinSaturationForMatch {
			totalS += c.HSL.S
			totalL += c.HSL.L
			count++
		}
	}

	// Fall back to reasonable defaults if no reference colors
	avgS := 50.0
	avgL := 55.0
	if count > 0 {
		avgS = totalS / float64(count)
		avgL = totalL / float64(count)
	}

	// Clamp to ensure the synthesized color is visually clear
	synS := math.Max(35, math.Min(75, avgS))
	synL := math.Max(40, math.Min(70, avgL))

	return NewColorFromHSL(targetHue, synS, synL)
}

// GenerateBrightVersion creates a lighter version of a color for bright ANSI slots.
// Scales the boost based on available headroom to avoid washing out bright colors.
func GenerateBrightVersion(c Color) Color {
	headroom := 90 - c.HSL.L
	boost := math.Max(5, math.Min(BrightColorLightnessBoost, headroom*0.6))
	newLightness := math.Min(90, c.HSL.L+boost)
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

// CalculateAverageHue calculates the average hue of chromatic colors using circular mean
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
func NormalizeBrightness(p *Palette) {
	bgLightness := p.Colors[0].HSL.L

	isVeryDarkBg := bgLightness < VeryDarkBackgroundThreshold
	isVeryLightBg := bgLightness > VeryLightBackgroundThreshold

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

	if index >= 1 && index <= 6 {
		p.Colors[index+8] = GenerateBrightVersion(p.Colors[index])
	}
}
