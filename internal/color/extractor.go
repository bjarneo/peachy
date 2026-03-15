package color

import (
	"fmt"
	"sort"

	"peachy/internal/shared"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

const (
	// imageScaleSize is the max dimension for fast processing (pixels)
	imageScaleSize = 300

	// minPixelsToSample is the minimum number of pixels needed for reliable extraction
	minPixelsToSample = 1000

	// maxPixelsToSample caps how many pixels we sample from scaled images
	maxPixelsToSample = 40000

	// dominantColorsToExtract is the target number of colors from median-cut
	dominantColorsToExtract = 48
)

// Extractor handles color extraction from images using native median-cut
type Extractor struct{}

// NewExtractor creates a new color extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// colorWithCount holds a color and its pixel count for sorting by dominance
type colorWithCount struct {
	Color Color
	Count int
}

// rgbPixel is a lightweight RGB value for pixel processing
type rgbPixel struct {
	r, g, b uint8
}

// colorBucket represents a group of pixels for median-cut quantization
type colorBucket struct {
	colors []rgbPixel
}

func (b *colorBucket) ranges() (rRange, gRange, bRange int) {
	if len(b.colors) == 0 {
		return 0, 0, 0
	}
	rMin, rMax := b.colors[0].r, b.colors[0].r
	gMin, gMax := b.colors[0].g, b.colors[0].g
	bMin, bMax := b.colors[0].b, b.colors[0].b

	for _, c := range b.colors[1:] {
		if c.r < rMin {
			rMin = c.r
		}
		if c.r > rMax {
			rMax = c.r
		}
		if c.g < gMin {
			gMin = c.g
		}
		if c.g > gMax {
			gMax = c.g
		}
		if c.b < bMin {
			bMin = c.b
		}
		if c.b > bMax {
			bMax = c.b
		}
	}
	return int(rMax) - int(rMin), int(gMax) - int(gMin), int(bMax) - int(bMin)
}

func (b *colorBucket) longestChannel() int {
	rr, gr, br := b.ranges()
	if rr >= gr && rr >= br {
		return 0
	}
	if gr >= br {
		return 1
	}
	return 2
}

func (b *colorBucket) split() (*colorBucket, *colorBucket) {
	ch := b.longestChannel()
	sort.Slice(b.colors, func(i, j int) bool {
		switch ch {
		case 0:
			return b.colors[i].r < b.colors[j].r
		case 1:
			return b.colors[i].g < b.colors[j].g
		default:
			return b.colors[i].b < b.colors[j].b
		}
	})
	mid := len(b.colors) / 2
	return &colorBucket{colors: b.colors[:mid]}, &colorBucket{colors: b.colors[mid:]}
}

func (b *colorBucket) averageColor() (rgbPixel, int) {
	if len(b.colors) == 0 {
		return rgbPixel{}, 0
	}
	var rSum, gSum, bSum int
	for _, c := range b.colors {
		rSum += int(c.r)
		gSum += int(c.g)
		bSum += int(c.b)
	}
	n := len(b.colors)
	return rgbPixel{
		r: uint8(rSum / n),
		g: uint8(gSum / n),
		b: uint8(bSum / n),
	}, n
}

func (b *colorBucket) volume() int {
	rr, gr, br := b.ranges()
	return rr * gr * br * len(b.colors)
}

// medianCut performs median-cut color quantization on pixel data
func medianCut(pixels []rgbPixel, numColors int) []colorWithCount {
	if len(pixels) == 0 {
		return nil
	}

	if len(pixels) <= numColors {
		seen := make(map[rgbPixel]bool)
		var result []colorWithCount
		for _, p := range pixels {
			if !seen[p] {
				seen[p] = true
				result = append(result, colorWithCount{
					Color: NewColorFromRGB(p.r, p.g, p.b),
					Count: 1,
				})
			}
		}
		return result
	}

	buckets := []*colorBucket{{colors: pixels}}

	for len(buckets) < numColors {
		// Find bucket with largest volume that can be split
		maxVolume := 0
		maxIdx := -1
		for i, b := range buckets {
			if len(b.colors) > 1 {
				v := b.volume()
				if v > maxVolume {
					maxVolume = v
					maxIdx = i
				}
			}
		}
		if maxIdx == -1 {
			break
		}

		left, right := buckets[maxIdx].split()
		newBuckets := make([]*colorBucket, 0, len(buckets)+1)
		newBuckets = append(newBuckets, buckets[:maxIdx]...)
		newBuckets = append(newBuckets, left, right)
		newBuckets = append(newBuckets, buckets[maxIdx+1:]...)
		buckets = newBuckets
	}

	var result []colorWithCount
	for _, b := range buckets {
		avg, count := b.averageColor()
		if count > 0 {
			result = append(result, colorWithCount{
				Color: NewColorFromRGB(avg.r, avg.g, avg.b),
				Count: count,
			})
		}
	}
	return result
}

// ExtractColors extracts dominant colors from an image using native median-cut
// Colors are returned sorted by frequency (most dominant first)
func (e *Extractor) ExtractColors(imagePath string) ([]Color, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}

	// Scale to imageScaleSize for fast processing (preserves aspect ratio)
	img = imaging.Fit(img, imageScaleSize, imageScaleSize, imaging.Lanczos)

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height
	sampleRate := totalPixels / maxPixelsToSample
	if sampleRate < 1 {
		sampleRate = 1
	}

	var pixels []rgbPixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y += sampleRate {
		for x := bounds.Min.X; x < bounds.Max.X; x += sampleRate {
			r, g, b, a := img.At(x, y).RGBA()
			// Skip transparent pixels
			if a < 0x8000 {
				continue
			}
			pixels = append(pixels, rgbPixel{
				r: uint8(r >> 8),
				g: uint8(g >> 8),
				b: uint8(b >> 8),
			})
		}
	}

	if len(pixels) < minPixelsToSample/10 {
		return nil, fmt.Errorf("not enough pixels to extract colors")
	}

	quantized := medianCut(pixels, dominantColorsToExtract)
	if len(quantized) == 0 {
		return nil, fmt.Errorf("no colors extracted from image")
	}

	// Sort by count (most dominant first)
	sort.Slice(quantized, func(i, j int) bool {
		return quantized[i].Count > quantized[j].Count
	})

	colors := make([]Color, len(quantized))
	for i, c := range quantized {
		colors[i] = c.Color
	}

	if len(colors) < 8 {
		return nil, fmt.Errorf("not enough colors extracted from image")
	}

	return colors, nil
}

// ExtractPalette extracts colors and assigns them to ANSI roles using the specified mode
func (e *Extractor) ExtractPalette(imagePath string, mode ExtractionMode, lightMode bool) (*Palette, error) {
	colors, err := e.ExtractColors(imagePath)
	if err != nil {
		return nil, err
	}

	var palette *Palette

	switch mode {
	case ModeMaterial:
		palette = GenerateMaterialPalette(colors, lightMode)
	case ModePastel:
		palette = GeneratePastelPalette(colors, lightMode)
	case ModeMonochromatic:
		palette = GenerateMonochromaticPalette(colors, lightMode)
	case ModeAnalogous:
		palette = GenerateAnalogousPalette(colors, lightMode)
	case ModeColorful:
		palette = GenerateColorfulPalette(colors, lightMode)
	case ModeMuted:
		palette = GenerateMutedPalette(colors, lightMode)
	case ModeBright:
		palette = GenerateBrightPalette(colors, lightMode)
	default:
		palette = GenerateNormalPalette(colors, lightMode)
	}

	// Apply brightness normalization for readability
	NormalizeBrightness(palette)

	return palette, nil
}

// IsValidImage checks if the file is a valid image format
func IsValidImage(path string) bool {
	return shared.IsValidImage(path)
}
