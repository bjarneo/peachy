package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sync"

	"peachy/internal/shared"

	"github.com/disintegration/imaging"
)

// ThumbnailCache manages cached thumbnail images for faster preview
type ThumbnailCache struct {
	cacheDir     string
	wallpaperDir string
	thumbnails   map[string]string // original path -> thumbnail path
	accessOrder  []string          // LRU tracking for cache eviction
	mu           sync.RWMutex
	errors       []error // Collected errors from concurrent operations
	errMu        sync.Mutex
}

// NewThumbnailCache creates a new thumbnail cache
func NewThumbnailCache() *ThumbnailCache {
	home, _ := os.UserHomeDir()

	cacheBase := os.Getenv("XDG_CACHE_HOME")
	if cacheBase == "" {
		cacheBase = filepath.Join(home, ".cache")
	}

	return &ThumbnailCache{
		cacheDir:     filepath.Join(cacheBase, "peachy", shared.CacheDir),
		wallpaperDir: filepath.Join(home, "Wallpapers"),
		thumbnails:   make(map[string]string),
		accessOrder:  make([]string, 0),
		errors:       make([]error, 0),
	}
}

// GetWallpaperDir returns the wallpapers directory path
func (c *ThumbnailCache) GetWallpaperDir() string {
	return c.wallpaperDir
}

// EnsureDirectories creates the wallpaper and cache directories if they don't exist
func (c *ThumbnailCache) EnsureDirectories() error {
	// Create wallpaper directory
	if err := os.MkdirAll(c.wallpaperDir, 0o755); err != nil {
		return fmt.Errorf("failed to create wallpaper directory: %w", err)
	}

	// Create cache directory
	if err := os.MkdirAll(c.cacheDir, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	return nil
}

// ScanAndCache scans the wallpaper directory and generates thumbnails
func (c *ThumbnailCache) ScanAndCache() error {
	if err := c.EnsureDirectories(); err != nil {
		return err
	}

	// Clear previous errors
	c.errMu.Lock()
	c.errors = make([]error, 0)
	c.errMu.Unlock()

	// Get list of image files
	entries, err := os.ReadDir(c.wallpaperDir)
	if err != nil {
		return fmt.Errorf("failed to read wallpaper directory: %w", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, shared.ConcurrentLimit)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !shared.IsValidImage(name) {
			continue
		}

		imagePath := filepath.Join(c.wallpaperDir, name)

		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if _, err := c.GetOrCreateThumbnail(path); err != nil {
				c.collectError(err)
			}
		}(imagePath)
	}

	wg.Wait()
	return c.firstError()
}

// collectError safely adds an error to the error collection
func (c *ThumbnailCache) collectError(err error) {
	c.errMu.Lock()
	defer c.errMu.Unlock()
	c.errors = append(c.errors, err)
}

// firstError returns the first collected error, if any
func (c *ThumbnailCache) firstError() error {
	c.errMu.Lock()
	defer c.errMu.Unlock()
	if len(c.errors) > 0 {
		return c.errors[0]
	}
	return nil
}

// Errors returns all collected errors from concurrent operations
func (c *ThumbnailCache) Errors() []error {
	c.errMu.Lock()
	defer c.errMu.Unlock()
	result := make([]error, len(c.errors))
	copy(result, c.errors)
	return result
}

// GetOrCreateThumbnail returns the thumbnail path, creating it if needed
func (c *ThumbnailCache) GetOrCreateThumbnail(imagePath string) (string, error) {
	// Check cache first
	c.mu.RLock()
	if thumb, ok := c.thumbnails[imagePath]; ok {
		c.mu.RUnlock()
		// Verify it still exists
		if _, err := os.Stat(thumb); err == nil {
			c.updateAccessOrder(imagePath)
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
		c.addToAccessOrder(imagePath)
		c.mu.Unlock()
		return thumbPath, nil
	}

	// Evict old entries if cache is full
	c.evictIfNeeded()

	// Generate thumbnail
	if err := c.generateThumbnail(imagePath, thumbPath); err != nil {
		return "", err
	}

	c.mu.Lock()
	c.thumbnails[imagePath] = thumbPath
	c.addToAccessOrder(imagePath)
	c.mu.Unlock()

	return thumbPath, nil
}

// updateAccessOrder moves an item to the end of the access order (most recently used)
func (c *ThumbnailCache) updateAccessOrder(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove from current position
	for i, p := range c.accessOrder {
		if p == path {
			c.accessOrder = append(c.accessOrder[:i], c.accessOrder[i+1:]...)
			break
		}
	}
	// Add to end
	c.accessOrder = append(c.accessOrder, path)
}

// addToAccessOrder adds an item to the access order (must hold write lock)
func (c *ThumbnailCache) addToAccessOrder(path string) {
	c.accessOrder = append(c.accessOrder, path)
}

// evictIfNeeded removes the oldest cache entries if the cache is full
func (c *ThumbnailCache) evictIfNeeded() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for len(c.thumbnails) >= shared.MaxCacheSize && len(c.accessOrder) > 0 {
		// Remove oldest entry
		oldest := c.accessOrder[0]
		c.accessOrder = c.accessOrder[1:]

		if thumbPath, ok := c.thumbnails[oldest]; ok {
			delete(c.thumbnails, oldest)
			// Optionally remove the file from disk
			_ = os.Remove(thumbPath)
		}
	}
}

// generateThumbnail creates a thumbnail image
func (c *ThumbnailCache) generateThumbnail(srcPath, dstPath string) error {
	// Open source image
	src, err := imaging.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open image: %w", err)
	}

	// Resize to thumbnail size using Box filter (fast, good quality for downscaling)
	thumb := imaging.Fit(src, shared.ThumbnailSize, shared.ThumbnailSize, imaging.Box)

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

	return os.MkdirAll(c.cacheDir, 0o755)
}

// generateCacheKey creates a unique key for the cache based on path and mtime
func generateCacheKey(path string, mtime int64) string {
	data := fmt.Sprintf("%s:%d", path, mtime)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
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

	_ = filepath.Walk(c.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		size += info.Size()
		return nil
	})

	return
}
