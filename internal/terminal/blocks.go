package terminal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

// renderCache caches rendered Unicode output to avoid re-processing
var (
	renderCache   = make(map[string]string)
	renderCacheMu sync.RWMutex
)

// Unicode block characters for image rendering
// Each character represents 2 vertical pixels
const (
	BlockFull      = '█' // Full block - both pixels same color
	BlockUpperHalf = '▀' // Upper half - top pixel colored, bottom is background
	BlockLowerHalf = '▄' // Lower half - bottom pixel colored, top is background
	BlockEmpty     = ' ' // Empty - both pixels are background
)

// RenderOptions configures image rendering
type RenderOptions struct {
	Width  int  // Target width in characters
	Height int  // Target height in characters (each char = 2 pixels)
	Dither bool // Whether to apply dithering (not implemented yet)
}

// DefaultRenderOptions returns sensible defaults
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		Width:  40,
		Height: 20,
		Dither: false,
	}
}

// RenderImage renders an image to Unicode block characters with ANSI colors
func RenderImage(imagePath string, opts RenderOptions) (string, error) {
	// Load the image
	img, err := loadImage(imagePath)
	if err != nil {
		return "", err
	}

	return RenderImageData(img, opts), nil
}

// RenderImageData renders an image.Image to Unicode blocks
func RenderImageData(img image.Image, opts RenderOptions) string {
	// Resize image to fit target dimensions
	// Height is doubled because each character represents 2 vertical pixels
	targetWidth := opts.Width
	targetHeight := opts.Height * 2

	// Use NearestNeighbor for speed - quality doesn't matter for terminal preview
	resized := imaging.Fit(img, targetWidth, targetHeight, imaging.NearestNeighbor)

	bounds := resized.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	var sb strings.Builder

	// Process 2 rows at a time (each character = 2 vertical pixels)
	for y := 0; y < height; y += 2 {
		for x := 0; x < width; x++ {
			// Get top pixel
			topColor := resized.At(bounds.Min.X+x, bounds.Min.Y+y)

			// Get bottom pixel (may not exist if odd height)
			var bottomColor color.Color
			if y+1 < height {
				bottomColor = resized.At(bounds.Min.X+x, bounds.Min.Y+y+1)
			} else {
				bottomColor = topColor
			}

			// Convert to RGB
			topR, topG, topB := colorToRGB(topColor)
			bottomR, bottomG, bottomB := colorToRGB(bottomColor)

			// Render the character with appropriate colors
			sb.WriteString(renderBlockChar(topR, topG, topB, bottomR, bottomG, bottomB))
		}
		sb.WriteString("\x1b[0m\n") // Reset colors at end of line
	}

	return sb.String()
}

// loadImage loads an image from a file path
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// colorToRGB converts a color.Color to RGB values (0-255)
func colorToRGB(c color.Color) (r, g, b uint8) {
	rr, gg, bb, _ := c.RGBA()
	return uint8(rr >> 8), uint8(gg >> 8), uint8(bb >> 8)
}

// renderBlockChar renders a single block character with ANSI true color
// Top pixel uses foreground color, bottom pixel uses background color
func renderBlockChar(topR, topG, topB, bottomR, bottomG, bottomB uint8) string {
	// Use upper half block (▀) with:
	// - Foreground = top pixel color
	// - Background = bottom pixel color
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm%c",
		topR, topG, topB,
		bottomR, bottomG, bottomB,
		BlockUpperHalf)
}

// RenderImageFit renders an image, fitting it within maxWidth x maxHeight
// while maintaining aspect ratio. Results are cached for performance.
func RenderImageFit(imagePath string, maxWidth, maxHeight int) (string, error) {
	// Generate cache key from path + dimensions + mtime
	info, err := os.Stat(imagePath)
	if err != nil {
		return "", err
	}
	cacheKey := generateRenderCacheKey(imagePath, maxWidth, maxHeight, info.ModTime().Unix())

	// Check cache first
	renderCacheMu.RLock()
	if cached, ok := renderCache[cacheKey]; ok {
		renderCacheMu.RUnlock()
		return cached, nil
	}
	renderCacheMu.RUnlock()

	// Load and render
	img, err := loadImage(imagePath)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// Calculate aspect ratio and fit within bounds
	// Remember: each character is 2 pixels tall
	charAspect := 2.0 // Characters are roughly twice as tall as wide

	// Calculate scaling to maintain aspect ratio
	scaleW := float64(maxWidth) / float64(imgWidth)
	scaleH := float64(maxHeight*2) / float64(imgHeight) // *2 because each char = 2 pixels

	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	targetWidth := int(float64(imgWidth) * scale)
	targetHeight := int(float64(imgHeight) * scale / charAspect)

	if targetWidth < 1 {
		targetWidth = 1
	}
	if targetHeight < 1 {
		targetHeight = 1
	}

	opts := RenderOptions{
		Width:  targetWidth,
		Height: targetHeight,
	}

	result := RenderImageData(img, opts)

	// Cache the result
	renderCacheMu.Lock()
	renderCache[cacheKey] = result
	renderCacheMu.Unlock()

	return result, nil
}

// generateRenderCacheKey creates a unique cache key for rendered output
func generateRenderCacheKey(path string, width, height int, mtime int64) string {
	data := fmt.Sprintf("%s:%d:%d:%d", path, width, height, mtime)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// ClearRenderCache clears the render cache
func ClearRenderCache() {
	renderCacheMu.Lock()
	renderCache = make(map[string]string)
	renderCacheMu.Unlock()
}

// GetImageDimensions returns the dimensions of an image file
func GetImageDimensions(imagePath string) (width, height int, err error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer func() { _ = file.Close() }()

	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return config.Width, config.Height, nil
}
