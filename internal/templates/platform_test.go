package templates

import (
	"runtime"
	"testing"
)

func TestGetPlatform(t *testing.T) {
	platform := GetPlatform()
	if platform != runtime.GOOS {
		t.Errorf("GetPlatform() = %q, want %q", platform, runtime.GOOS)
	}
}

func TestIsSupportedPlatform(t *testing.T) {
	supported := IsSupportedPlatform()
	// Should be true on Linux and macOS
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if !supported {
			t.Errorf("IsSupportedPlatform() = false, want true on %s", runtime.GOOS)
		}
	}
}

func TestGetPlatformName(t *testing.T) {
	name := GetPlatformName()
	switch runtime.GOOS {
	case "linux":
		if name != "Linux" {
			t.Errorf("GetPlatformName() = %q, want %q", name, "Linux")
		}
	case "darwin":
		if name != "macOS" {
			t.Errorf("GetPlatformName() = %q, want %q", name, "macOS")
		}
	}
}

func TestGetDefaultConfigDir(t *testing.T) {
	dir := GetDefaultConfigDir()
	if dir == "" {
		t.Error("GetDefaultConfigDir() returned empty string")
	}
}

func TestGetDefaultCacheDir(t *testing.T) {
	dir := GetDefaultCacheDir()
	if dir == "" {
		t.Error("GetDefaultCacheDir() returned empty string")
	}
}
