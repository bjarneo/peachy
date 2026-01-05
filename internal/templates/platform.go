package templates

import (
	"os"
	"path/filepath"
	"runtime"
)

// Platform constants
const (
	PlatformLinux  = "linux"
	PlatformDarwin = "darwin" // macOS
)

// GetPlatform returns the current operating system
func GetPlatform() string {
	return runtime.GOOS
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return runtime.GOOS == PlatformLinux
}

// IsMacOS returns true if running on macOS
func IsMacOS() bool {
	return runtime.GOOS == PlatformDarwin
}

// GetDefaultConfigDir returns the platform-appropriate config directory
// On Linux: ~/.config (XDG_CONFIG_HOME)
// On macOS: ~/.config (for consistency, though ~/Library/Application Support is also valid)
func GetDefaultConfigDir() string {
	// Check XDG_CONFIG_HOME first (works on both platforms)
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		return configHome
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Use ~/.config for consistency across platforms
	// This is the standard on Linux and commonly used on macOS for CLI tools
	return filepath.Join(home, ".config")
}

// GetDefaultCacheDir returns the platform-appropriate cache directory
// On Linux: ~/.cache (XDG_CACHE_HOME)
// On macOS: ~/Library/Caches or ~/.cache
func GetDefaultCacheDir() string {
	// Check XDG_CACHE_HOME first
	if cacheHome := os.Getenv("XDG_CACHE_HOME"); cacheHome != "" {
		return cacheHome
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Use ~/.cache for consistency across platforms
	return filepath.Join(home, ".cache")
}

// IsSupportedPlatform returns true if the current platform is supported
func IsSupportedPlatform() bool {
	return IsLinux() || IsMacOS()
}

// GetPlatformName returns a human-readable platform name
func GetPlatformName() string {
	switch runtime.GOOS {
	case PlatformLinux:
		return "Linux"
	case PlatformDarwin:
		return "macOS"
	default:
		return runtime.GOOS
	}
}
