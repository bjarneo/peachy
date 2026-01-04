package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

const (
	ThumbnailSize = 200 // Max dimension for thumbnails
	CacheDir      = "thumbnails"
)

// ThumbnailCache manages cached thumbnail images for faster preview
type ThumbnailCache struct {
	cacheDir    string
	wallpaperDir string
	thumbnails  map[string]string // original path -> thumbnail path
	mu          sync.RWMutex
}

// NewThumbnailCache creates a new thumbnail cache
func NewThumbnailCache() *ThumbnailCache {
	home, _ := os.UserHomeDir()

	cacheBase := os.Getenv("XDG_CACHE_HOME")
	if cacheBase == "" {
		cacheBase = filepath.Join(home, ".cache")
	}

	return &ThumbnailCache{
		cacheDir:     filepath.Join(cacheBase, "peachy", CacheDir),
		wallpaperDir: filepath.Join(home, "Wallpapers"),
		thumbnails:   make(map[string]string),
	}
}

// GetWallpaperDir returns the wallpapers directory path
func (c *ThumbnailCache) GetWallpaperDir() string {
	return c.wallpaperDir
}

// EnsureDirectories creates the wallpaper and cache directories if they don't exist
func (c *ThumbnailCache) EnsureDirectories() error {
	// Create wallpaper directory
	if err := os.MkdirAll(c.wallpaperDir, 0755); err != nil {
		return fmt.Errorf("failed to create wallpaper directory: %w", err)
	}

	// Create cache directory
	if err := os.MkdirAll(c.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	return nil
}

// ScanAndCache scans the wallpaper directory and generates thumbnails
func (c *ThumbnailCache) ScanAndCache() error {
	if err := c.EnsureDirectories(); err != nil {
		return err
	}

	// Get list of image files
	entries, err := os.ReadDir(c.wallpaperDir)
	if err != nil {
		return fmt.Errorf("failed to read wallpaper directory: %w", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // Limit concurrent thumbnail generation

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !isValidImage(name) {
			continue
		}

		imagePath := filepath.Join(c.wallpaperDir, name)

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			c.GetOrCreateThumbnail(path)
		}(imagePath)
	}

	wg.Wait()
	return nil
}

// GetOrCreateThumbnail returns the thumbnail path, creating it if needed
func (c *ThumbnailCache) GetOrCreateThumbnail(imagePath string) (string, error) {
	// Check cache first
	c.mu.RLock()
	if thumb, ok := c.thumbnails[imagePath]; ok {
		c.mu.RUnlock()
		// Verify it still exists
		if _, err := os.Stat(thumb); err == nil {
			return thumb, nil
		}
	}
	c.mu.RUnlock()

	// Generate cache key from path and modification time
	info, err := os.Stat(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat image: %w", err)
	}

	cacheKey := generateCacheKey(imagePath, info.ModTime().Unix())
	thumbPath := filepath.Join(c.cacheDir, cacheKey+".png")

	// Check if thumbnail already exists on disk
	if _, err := os.Stat(thumbPath); err == nil {
		c.mu.Lock()
		c.thumbnails[imagePath] = thumbPath
		c.mu.Unlock()
		return thumbPath, nil
	}

	// Generate thumbnail
	if err := c.generateThumbnail(imagePath, thumbPath); err != nil {
		return "", err
	}

	c.mu.Lock()
	c.thumbnails[imagePath] = thumbPath
	c.mu.Unlock()

	return thumbPath, nil
}

// generateThumbnail creates a thumbnail image
func (c *ThumbnailCache) generateThumbnail(srcPath, dstPath string) error {
	// Open source image
	src, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// Resize to thumbnail size using Box filter (fast, good quality for downscaling)
	thumb := imaging.Fit(src, ThumbnailSize, ThumbnailSize, imaging.Box)

	// Save as JPEG for smaller size and faster encoding
	if err := imaging.Save(thumb, dstPath); err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	return nil
}

// GetThumbnailPath returns the cached thumbnail path if it exists
func (c *ThumbnailCache) GetThumbnailPath(imagePath string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if thumb, ok := c.thumbnails[imagePath]; ok {
		return thumb
	}
	return ""
}

// HasThumbnail checks if a thumbnail exists for the given image
func (c *ThumbnailCache) HasThumbnail(imagePath string) bool {
	return c.GetThumbnailPath(imagePath) != ""
}

// LoadThumbnail loads a cached thumbnail image
func (c *ThumbnailCache) LoadThumbnail(imagePath string) (image.Image, error) {
	thumbPath, err := c.GetOrCreateThumbnail(imagePath)
	if err != nil {
		return nil, err
	}

	return imaging.Open(thumbPath)
}

// ClearCache removes all cached thumbnails
func (c *ThumbnailCache) ClearCache() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.RemoveAll(c.cacheDir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	c.thumbnails = make(map[string]string)

	return os.MkdirAll(c.cacheDir, 0755)
}

// generateCacheKey creates a unique key for the cache based on path and mtime
func generateCacheKey(path string, mtime int64) string {
	data := fmt.Sprintf("%s:%d", path, mtime)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// isValidImage checks if a file is a valid image format
func isValidImage(name string) bool {
	lower := strings.ToLower(name)
	validExts := []string{".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp"}
	for _, ext := range validExts {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// GetCacheDir returns the cache directory path
func (c *ThumbnailCache) GetCacheDir() string {
	return c.cacheDir
}

// GetCacheStats returns cache statistics
func (c *ThumbnailCache) GetCacheStats() (count int, size int64) {
	c.mu.RLock()
	count = len(c.thumbnails)
	c.mu.RUnlock()

	filepath.Walk(c.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})

	return
}
