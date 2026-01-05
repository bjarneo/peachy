package shared

// Color indices for the 16-color ANSI palette
const (
	ColorBlack         = 0
	ColorRed           = 1
	ColorGreen         = 2
	ColorYellow        = 3
	ColorBlue          = 4
	ColorMagenta       = 5
	ColorCyan          = 6
	ColorWhite         = 7
	ColorBrightBlack   = 8
	ColorBrightRed     = 9
	ColorBrightGreen   = 10
	ColorBrightYellow  = 11
	ColorBrightBlue    = 12
	ColorBrightMagenta = 13
	ColorBrightCyan    = 14
	ColorBrightWhite   = 15
)

// PaletteSize is the number of colors in a terminal palette
const PaletteSize = 16

// NormalColorCount is the number of normal (non-bright) colors
const NormalColorCount = 8

// UI dimensions
const (
	MinPanelWidth   = 30
	MaxPanelWidth   = 50
	MinPanelHeight  = 10
	MinWindowWidth  = 40
	MinWindowHeight = 20
)

// Cache settings
const (
	ThumbnailSize   = 200 // Max dimension for thumbnails
	MaxCacheSize    = 500 // Maximum number of cached thumbnails
	CacheDir        = "thumbnails"
	ConcurrentLimit = 4 // Max concurrent thumbnail generation
)

// ValidImageExtensions contains all supported image file extensions
var ValidImageExtensions = []string{".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp", ".tiff"}
