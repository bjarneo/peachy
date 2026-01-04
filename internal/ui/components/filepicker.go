package components

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"peachy/internal/cache"
	"peachy/internal/color"
	"peachy/internal/terminal"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FileEntry represents a file or directory
type FileEntry struct {
	Name  string
	Path  string
	IsDir bool
}

// FilePicker allows browsing and selecting image files
type FilePicker struct {
	currentDir     string
	entries        []FileEntry
	allEntries     []FileEntry // Unfiltered entries for search
	cursor         int
	selected       string
	err            error
	width          int
	height         int
	visible        bool
	offset         int // For scrolling
	maxVisible     int
	showPreview    bool
	previewWidth   int
	previewHeight  int
	previewCache   string // Cached preview render
	previewPath    string // Path of cached preview
	thumbnailCache *cache.ThumbnailCache

	// Search mode
	searchMode  bool
	searchQuery string
}

// NewFilePicker creates a new file picker with ~/Wallpapers as default
func NewFilePicker() FilePicker {
	thumbCache := cache.NewThumbnailCache()

	// Ensure wallpaper directory exists
	_ = thumbCache.EnsureDirectories()

	fp := FilePicker{
		currentDir:     thumbCache.GetWallpaperDir(),
		width:          80,
		height:         25,
		visible:        false,
		maxVisible:     15,
		showPreview:    true,
		previewWidth:   35,
		previewHeight:  12,
		thumbnailCache: thumbCache,
	}
	fp.loadDir()
	return fp
}

// SetThumbnailCache sets the thumbnail cache
func (f *FilePicker) SetThumbnailCache(tc *cache.ThumbnailCache) {
	f.thumbnailCache = tc
	// Update default directory to wallpaper dir
	f.currentDir = tc.GetWallpaperDir()
}

// Open shows the file picker
func (f *FilePicker) Open() {
	f.visible = true
	f.selected = ""
	f.previewCache = ""
	f.previewPath = ""
	f.loadDir()
	f.updatePreview()
}

// Close hides the file picker
func (f *FilePicker) Close() {
	f.visible = false
	f.previewCache = ""
	f.previewPath = ""
}

// Visible returns whether the picker is visible
func (f FilePicker) Visible() bool {
	return f.visible
}

// Selected returns the selected file path (empty if none)
func (f FilePicker) Selected() string {
	return f.selected
}

// CurrentDir returns the current directory
func (f FilePicker) CurrentDir() string {
	return f.currentDir
}

// SetDir changes to a directory
func (f *FilePicker) SetDir(dir string) {
	f.currentDir = dir
	f.cursor = 0
	f.offset = 0
	f.previewCache = ""
	f.previewPath = ""
	f.loadDir()
	f.updatePreview()
}

// SetSize sets the picker size
func (f *FilePicker) SetSize(w, h int) {
	f.width = w
	f.height = h
	f.maxVisible = h - 10
	if f.maxVisible < 5 {
		f.maxVisible = 5
	}
	// Preview takes right third of the space
	f.previewWidth = w / 3
	if f.previewWidth < 20 {
		f.previewWidth = 20
	}
	f.previewHeight = h - 8
	if f.previewHeight < 5 {
		f.previewHeight = 5
	}
}

func (f *FilePicker) loadDir() {
	f.entries = nil
	f.allEntries = nil
	f.err = nil
	f.searchMode = false
	f.searchQuery = ""

	// Add parent directory entry
	if f.currentDir != "/" {
		f.entries = append(f.entries, FileEntry{
			Name:  "..",
			Path:  filepath.Dir(f.currentDir),
			IsDir: true,
		})
	}

	// Read directory
	files, err := os.ReadDir(f.currentDir)
	if err != nil {
		f.err = err
		f.allEntries = f.entries
		return
	}

	// Separate dirs and files
	var dirs, images []FileEntry

	for _, file := range files {
		// Skip hidden files
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		entry := FileEntry{
			Name:  file.Name(),
			Path:  filepath.Join(f.currentDir, file.Name()),
			IsDir: file.IsDir(),
		}

		if file.IsDir() {
			dirs = append(dirs, entry)
		} else if color.IsValidImage(file.Name()) {
			images = append(images, entry)
		}
	}

	// Sort alphabetically
	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(images, func(i, j int) bool {
		return strings.ToLower(images[i].Name) < strings.ToLower(images[j].Name)
	})

	// Dirs first, then images
	f.entries = append(f.entries, dirs...)
	f.entries = append(f.entries, images...)

	// Keep a copy for search filtering
	f.allEntries = make([]FileEntry, len(f.entries))
	copy(f.allEntries, f.entries)
}

// updatePreview updates the preview for the current selection
func (f *FilePicker) updatePreview() {
	if !f.showPreview || f.cursor >= len(f.entries) {
		f.previewCache = ""
		f.previewPath = ""
		return
	}

	entry := f.entries[f.cursor]

	// Only preview image files
	if entry.IsDir {
		f.previewCache = ""
		f.previewPath = ""
		return
	}

	// Check if already cached
	if entry.Path == f.previewPath && f.previewCache != "" {
		return
	}

	// Always use cached thumbnail for efficiency
	var preview string
	var err error

	if f.thumbnailCache != nil {
		// Get or create thumbnail - this ensures we always use the smaller cached version
		thumbPath, thumbErr := f.thumbnailCache.GetOrCreateThumbnail(entry.Path)
		if thumbErr == nil && thumbPath != "" {
			// Use cached thumbnail (200x200 max) - much faster to render
			preview, err = terminal.RenderImageFit(thumbPath, f.previewWidth-4, f.previewHeight-4)
		}
	}

	// Only fall back to original if thumbnail cache is not available
	if preview == "" || err != nil {
		preview, err = terminal.RenderImageFit(entry.Path, f.previewWidth-4, f.previewHeight-4)
	}

	if err != nil {
		f.previewCache = ""
		f.previewPath = ""
		return
	}

	f.previewCache = preview
	f.previewPath = entry.Path
}

// Init implements tea.Model
func (f FilePicker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (f FilePicker) Update(msg tea.Msg) (FilePicker, tea.Cmd) {
	if !f.visible {
		return f, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle search mode input
		if f.searchMode {
			return f.handleSearchInput(msg)
		}

		prevCursor := f.cursor

		switch msg.String() {
		case "up", "k":
			if f.cursor > 0 {
				f.cursor--
				f.adjustScroll()
			}
		case "down", "j":
			if f.cursor < len(f.entries)-1 {
				f.cursor++
				f.adjustScroll()
			}
		case "enter":
			if f.cursor < len(f.entries) {
				entry := f.entries[f.cursor]
				if entry.IsDir {
					f.SetDir(entry.Path)
				} else {
					f.selected = entry.Path
					f.visible = false
				}
			}
		case "backspace":
			// Go to parent directory
			parent := filepath.Dir(f.currentDir)
			if parent != f.currentDir {
				f.SetDir(parent)
			}
		case "home", "g":
			f.cursor = 0
			f.offset = 0
		case "end", "G":
			f.cursor = len(f.entries) - 1
			f.adjustScroll()
		case "~":
			// Go to home directory
			home, _ := os.UserHomeDir()
			f.SetDir(home)
		case "w":
			// Go to Wallpapers directory
			if f.thumbnailCache != nil {
				f.SetDir(f.thumbnailCache.GetWallpaperDir())
			}
		case "p":
			// Toggle preview
			f.showPreview = !f.showPreview
			if f.showPreview {
				f.updatePreview()
			}
		case "/":
			// Enter search mode
			f.searchMode = true
			f.searchQuery = ""
		}

		// Update preview if cursor changed
		if f.cursor != prevCursor {
			f.updatePreview()
		}
	}

	return f, nil
}

// handleSearchInput handles keyboard input in search mode
func (f FilePicker) handleSearchInput(msg tea.KeyMsg) (FilePicker, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Exit search mode and restore all entries
		f.searchMode = false
		f.searchQuery = ""
		f.entries = f.allEntries
		f.cursor = 0
		f.offset = 0
		f.updatePreview()
	case "enter":
		// Exit search mode but keep filtered results
		f.searchMode = false
		if f.cursor < len(f.entries) {
			entry := f.entries[f.cursor]
			if entry.IsDir {
				f.SetDir(entry.Path)
			} else {
				f.selected = entry.Path
				f.visible = false
			}
		}
	case "backspace":
		if len(f.searchQuery) > 0 {
			f.searchQuery = f.searchQuery[:len(f.searchQuery)-1]
			f.filterEntries()
		}
	case "up", "ctrl+p":
		if f.cursor > 0 {
			f.cursor--
			f.adjustScroll()
			f.updatePreview()
		}
	case "down", "ctrl+n":
		if f.cursor < len(f.entries)-1 {
			f.cursor++
			f.adjustScroll()
			f.updatePreview()
		}
	default:
		// Add character to search query
		key := msg.String()
		if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
			f.searchQuery += key
			f.filterEntries()
		}
	}

	return f, nil
}

// filterEntries filters entries based on search query
func (f *FilePicker) filterEntries() {
	if f.searchQuery == "" {
		f.entries = f.allEntries
		f.cursor = 0
		f.offset = 0
		f.updatePreview()
		return
	}

	query := strings.ToLower(f.searchQuery)
	var filtered []FileEntry

	for _, entry := range f.allEntries {
		if strings.Contains(strings.ToLower(entry.Name), query) {
			filtered = append(filtered, entry)
		}
	}

	f.entries = filtered
	f.cursor = 0
	f.offset = 0
	f.updatePreview()
}

func (f *FilePicker) adjustScroll() {
	// Keep cursor visible
	if f.cursor < f.offset {
		f.offset = f.cursor
	}
	if f.cursor >= f.offset+f.maxVisible {
		f.offset = f.cursor - f.maxVisible + 1
	}
}

// View renders the file picker
func (f FilePicker) View() string {
	if !f.visible {
		return ""
	}

	// Render file list
	fileList := f.renderFileList()

	// Render preview if enabled
	var preview string
	if f.showPreview {
		preview = f.renderPreview()
	}

	// Combine side by side
	var content string
	if f.showPreview && preview != "" {
		content = lipgloss.JoinHorizontal(lipgloss.Top, fileList, "  ", preview)
	} else {
		content = fileList
	}

	return content
}

func (f FilePicker) renderFileList() string {
	var sb strings.Builder

	listWidth := f.width - f.previewWidth - 10
	if !f.showPreview {
		listWidth = f.width - 6
	}
	if listWidth < 30 {
		listWidth = 30
	}

	// Title - using ANSI colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("13")) // Bright Magenta
	sb.WriteString(titleStyle.Render("Select Image"))
	sb.WriteString("\n\n")

	// Current path
	pathStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // Bright Black (gray)
		Width(listWidth - 4)
	truncatedPath := f.currentDir
	if len(truncatedPath) > listWidth-8 {
		truncatedPath = "..." + truncatedPath[len(truncatedPath)-(listWidth-11):]
	}
	sb.WriteString(pathStyle.Render(truncatedPath))
	sb.WriteString("\n")

	// Search bar (show when in search mode or has query)
	if f.searchMode || f.searchQuery != "" {
		searchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")). // Magenta
			Bold(true)
		inputStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")). // Bright White
			Background(lipgloss.Color("8")).  // Bright Black
			Padding(0, 1)

		searchLabel := searchStyle.Render("Search: ")
		searchInput := f.searchQuery
		if f.searchMode {
			searchInput += "▌" // Cursor
		}
		sb.WriteString(searchLabel)
		sb.WriteString(inputStyle.Render(searchInput))
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("─", listWidth-4))
	sb.WriteString("\n")

	// Error message
	if f.err != nil {
		sb.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")). // Red
			Render("Error: " + f.err.Error()))
		sb.WriteString("\n")
	}

	// No results message
	if len(f.entries) == 0 && f.searchQuery != "" {
		noResultStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Bright Black
			Italic(true)
		sb.WriteString(noResultStyle.Render("No matching files"))
		sb.WriteString("\n")
	}

	// File list - using ANSI colors
	dirStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")) // Bright Blue
	fileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // White
	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("4")).  // Blue
		Foreground(lipgloss.Color("15")). // Bright White
		Width(listWidth - 6)
	matchStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")). // Yellow
		Bold(true)

	// Render visible entries
	end := f.offset + f.maxVisible
	if end > len(f.entries) {
		end = len(f.entries)
	}

	for i := f.offset; i < end; i++ {
		entry := f.entries[i]

		var icon string
		var style lipgloss.Style

		if entry.IsDir {
			icon = "📁 "
			style = dirStyle
		} else {
			icon = "🖼  "
			style = fileStyle
		}

		name := entry.Name
		maxNameLen := listWidth - 12
		if maxNameLen < 10 {
			maxNameLen = 10
		}
		if len(name) > maxNameLen {
			name = name[:maxNameLen-3] + "..."
		}

		// Highlight matching part if searching
		var displayName string
		if f.searchQuery != "" {
			displayName = highlightMatch(name, f.searchQuery, style, matchStyle)
		} else {
			displayName = style.Render(name)
		}

		if i == f.cursor {
			// For selected item, render without highlight styling
			line := "> " + icon + name
			line = selectedStyle.Render(line)
			sb.WriteString(line)
		} else {
			sb.WriteString("  " + icon)
			sb.WriteString(displayName)
		}
		sb.WriteString("\n")
	}

	// Pad to consistent height
	for i := end - f.offset; i < f.maxVisible; i++ {
		sb.WriteString("\n")
	}

	// Scroll indicator
	if len(f.entries) > f.maxVisible {
		scrollInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Bright Black
			Render("[▴/▾ scroll]")
		sb.WriteString(scrollInfo)
		sb.WriteString("\n")
	} else {
		sb.WriteString("\n")
	}

	// Help
	sb.WriteString("\n")
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")) // Bright Black
	if f.searchMode {
		sb.WriteString(helpStyle.Render("Type to search • Enter: select • Esc: cancel"))
	} else {
		sb.WriteString(helpStyle.Render("Enter: select • /: search • w: wallpapers • p: preview"))
	}

	// Border - using ANSI colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("5")). // Magenta
		Padding(1).
		Width(listWidth)

	return borderStyle.Render(sb.String())
}

// highlightMatch highlights the matching part of a string
func highlightMatch(text, query string, normalStyle, highlightStyle lipgloss.Style) string {
	lowerText := strings.ToLower(text)
	lowerQuery := strings.ToLower(query)

	idx := strings.Index(lowerText, lowerQuery)
	if idx == -1 {
		return normalStyle.Render(text)
	}

	before := text[:idx]
	match := text[idx : idx+len(query)]
	after := text[idx+len(query):]

	return normalStyle.Render(before) + highlightStyle.Render(match) + normalStyle.Render(after)
}

func (f FilePicker) renderPreview() string {
	var sb strings.Builder

	// Title - using ANSI colors
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")) // Bright Blue
	sb.WriteString(titleStyle.Render("PREVIEW"))
	sb.WriteString("\n\n")

	if f.previewCache != "" {
		// Show image info
		if f.cursor < len(f.entries) {
			entry := f.entries[f.cursor]
			infoStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("8")) // Bright Black

			filename := filepath.Base(entry.Path)
			if len(filename) > f.previewWidth-6 {
				filename = filename[:f.previewWidth-9] + "..."
			}

			// Get dimensions
			w, h, err := terminal.GetImageDimensions(entry.Path)
			if err == nil {
				sb.WriteString(infoStyle.Render(filename))
				sb.WriteString("\n")
				sb.WriteString(infoStyle.Render(strings.Repeat("─", f.previewWidth-4)))
				sb.WriteString("\n")
				sb.WriteString(infoStyle.Render(lipgloss.NewStyle().
					Foreground(lipgloss.Color("8")).
					Render("Size: ") + lipgloss.NewStyle().
					Foreground(lipgloss.Color("7")).
					Render(strings.TrimSpace(strings.Join([]string{intToString(w), "x", intToString(h)}, "")))))
				sb.WriteString("\n\n")
			}
		}

		// Show preview
		sb.WriteString(f.previewCache)
	} else if f.cursor < len(f.entries) && f.entries[f.cursor].IsDir {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Bright Black
			Italic(true)
		sb.WriteString(emptyStyle.Render("Directory selected"))
	} else {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // Bright Black
			Italic(true)
		sb.WriteString(emptyStyle.Render("No preview"))
	}

	// Border - using ANSI colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")). // Bright Black
		Padding(1).
		Width(f.previewWidth).
		Height(f.previewHeight)

	return borderStyle.Render(sb.String())
}

// intToString converts an integer to string without importing strconv
func intToString(i int) string {
	if i == 0 {
		return "0"
	}
	var result []byte
	negative := i < 0
	if negative {
		i = -i
	}
	for i > 0 {
		result = append([]byte{byte('0' + i%10)}, result...)
		i /= 10
	}
	if negative {
		result = append([]byte{'-'}, result...)
	}
	return string(result)
}
