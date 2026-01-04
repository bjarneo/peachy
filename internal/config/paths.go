package config

import (
	"os"
	"path/filepath"
)

const (
	AppName    = "peachy"
	ConfigFile = "colors.toml"
	ThemesDir  = "themes"
	ThemeFile  = "theme"
)

// GetConfigDir returns the configuration directory path
// Respects XDG_CONFIG_HOME if set
func GetConfigDir() string {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, AppName)
}

// GetDefaultConfigPath returns the default colors.toml path
func GetDefaultConfigPath() string {
	return filepath.Join(GetConfigDir(), ConfigFile)
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	dir := GetConfigDir()
	return os.MkdirAll(dir, 0755)
}

// GetThemesDir returns the themes directory path
func GetThemesDir() string {
	return filepath.Join(GetConfigDir(), ThemesDir)
}

// EnsureThemesDir creates the themes directory if it doesn't exist
func EnsureThemesDir() error {
	dir := GetThemesDir()
	return os.MkdirAll(dir, 0755)
}

// GetThemePath returns the path for a named theme
func GetThemePath(name string) string {
	return filepath.Join(GetThemesDir(), name+".toml")
}

// GetActiveThemePath returns the path to the active theme file
func GetActiveThemePath() string {
	return filepath.Join(GetConfigDir(), ThemeFile)
}

// ExpandPath expands ~ to home directory
func ExpandPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
