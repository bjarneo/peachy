package ui

import (
	"path/filepath"
	"strings"

	"peachy/internal/cache"
	"peachy/internal/color"
	"peachy/internal/config"
	"peachy/internal/shared"
	"peachy/internal/ui/components"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// CacheScanCompleteMsg is sent when thumbnail cache scanning is complete
type CacheScanCompleteMsg struct {
	Count int
	Err   error
}

// ViewState represents which view/mode is active
type ViewState int

const (
	ViewMain ViewState = iota
	ViewEditor
	ViewFilePicker
	ViewHelp
	ViewSaveTheme
	ViewThemeBrowser
)

// Model is the main application model
type Model struct {
	// Components
	colorList   components.ColorList
	colorEditor components.ColorEditor
	filePicker  components.FilePicker
	preview     components.Preview
	help        components.Help

	// State
	palette         *color.Palette
	originalPalette *color.Palette
	imagePath       string
	configPath      string
	viewState       ViewState
	status          shared.Status

	// Extraction mode
	extractionMode color.ExtractionMode
	lightMode      bool

	// Services
	extractor      *color.Extractor
	thumbnailCache *cache.ThumbnailCache

	// Key bindings
	keys KeyMap

	// Window size
	width  int
	height int

	// Flags
	quitting     bool
	cacheScanned bool

	// Theme name input
	themeNameInput string

	// Theme browser
	themes      []string
	themeCursor int
}

// NewModel creates a new application model
func NewModel() Model {
	palette := color.NewPalette()
	thumbCache := cache.NewThumbnailCache()

	// Ensure directories exist
	_ = thumbCache.EnsureDirectories()

	// Create file picker with thumbnail cache
	filePicker := components.NewFilePicker()
	filePicker.SetThumbnailCache(thumbCache)

	// Create preview with thumbnail cache
	preview := components.NewPreview(palette)
	preview.SetThumbnailCache(thumbCache)

	return Model{
		colorList:      components.NewColorList(palette),
		colorEditor:    components.NewColorEditor(),
		filePicker:     filePicker,
		preview:        preview,
		help:           components.NewHelp(),
		palette:        palette,
		extractor:      color.NewExtractor(),
		thumbnailCache: thumbCache,
		extractionMode: color.ModeNormal,
		lightMode:      false,
		keys:           DefaultKeyMap(),
		viewState:      ViewMain,
		status:         shared.Info("Scanning wallpapers..."),
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	// Start scanning thumbnails in background
	return m.scanThumbnails()
}

// scanThumbnails starts background thumbnail scanning
func (m Model) scanThumbnails() tea.Cmd {
	return func() tea.Msg {
		if m.thumbnailCache == nil {
			return CacheScanCompleteMsg{Count: 0, Err: nil}
		}

		err := m.thumbnailCache.ScanAndCache()
		count, _ := m.thumbnailCache.GetCacheStats()

		return CacheScanCompleteMsg{Count: count, Err: err}
	}
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case CacheScanCompleteMsg:
		m.cacheScanned = true
		if msg.Err != nil {
			m.status = shared.Error("Cache scan failed: " + msg.Err.Error())
		} else if msg.Count > 0 {
			m.status = shared.Info("Ready. Press 'o' to open ~/Wallpapers")
		} else {
			m.status = shared.Info("Press 'o' to open an image, 'l' to load colors.toml")
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateComponentSizes()
		return m, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		// Handle based on current view
		switch m.viewState {
		case ViewHelp:
			return m.updateHelp(msg)
		case ViewFilePicker:
			return m.updateFilePicker(msg)
		case ViewEditor:
			return m.updateEditor(msg)
		case ViewSaveTheme:
			return m.updateSaveTheme(msg)
		case ViewThemeBrowser:
			return m.updateThemeBrowser(msg)
		default:
			return m.updateMain(msg)
		}
	}

	return m, nil
}

func (m Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		m.quitting = true
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		m.help.Show()
		m.viewState = ViewHelp
		return m, nil
	case key.Matches(msg, m.keys.Open):
		m.filePicker.Open()
		m.viewState = ViewFilePicker
		return m, nil
	case key.Matches(msg, m.keys.Extract):
		return m.extractColors()
	case key.Matches(msg, m.keys.Save):
		m.themeNameInput = ""
		m.viewState = ViewSaveTheme
		return m, nil
	case key.Matches(msg, m.keys.Load):
		return m.openThemeBrowser()
	case key.Matches(msg, m.keys.Reset):
		if m.originalPalette != nil {
			m.palette = m.originalPalette.Clone()
			m.colorList.SetPalette(m.palette)
			m.preview.SetPalette(m.palette)
			m.status = shared.Info("Reset to extracted colors")
		}
		return m, nil
	case msg.String() == "m":
		m.extractionMode = color.NextMode(m.extractionMode)
		m.status = shared.Info("Mode: " + color.ModeNames[m.extractionMode])
		if m.imagePath != "" {
			return m.extractColors()
		}
		return m, nil
	case msg.String() == "M":
		m.extractionMode = color.PrevMode(m.extractionMode)
		m.status = shared.Info("Mode: " + color.ModeNames[m.extractionMode])
		if m.imagePath != "" {
			return m.extractColors()
		}
		return m, nil
	case msg.String() == "t":
		m.lightMode = !m.lightMode
		if m.lightMode {
			m.status = shared.Info("Light mode enabled")
		} else {
			m.status = shared.Info("Dark mode enabled")
		}
		if m.imagePath != "" {
			return m.extractColors()
		}
		return m, nil
	case msg.String() == "a":
		return m.applyCurrentPalette()
	case key.Matches(msg, m.keys.Enter):
		m.colorEditor.Open(m.colorList.SelectedColor(), m.colorList.Cursor())
		m.viewState = ViewEditor
		return m, nil
	default:
		var cmd tea.Cmd
		m.colorList, cmd = m.colorList.Update(msg)
		return m, cmd
	}
}

func (m Model) updateEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Escape):
		m.colorEditor.Close()
		m.viewState = ViewMain
		return m, nil
	case key.Matches(msg, m.keys.Enter):
		if !m.colorEditor.Visible() {
			return m, nil
		}
		m.palette.SetColor(m.colorEditor.ColorIndex(), m.colorEditor.CurrentColor())
		m.colorList.SetPalette(m.palette)
		m.preview.SetPalette(m.palette)
		m.colorEditor.Close()
		m.viewState = ViewMain
		m.status = shared.Success("Color updated")
		return m, nil
	default:
		var cmd tea.Cmd
		m.colorEditor, cmd = m.colorEditor.Update(msg)
		m.palette.SetColor(m.colorEditor.ColorIndex(), m.colorEditor.CurrentColor())
		m.preview.SetPalette(m.palette)
		return m, cmd
	}
}

func (m Model) updateFilePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Escape):
		m.filePicker.Close()
		m.viewState = ViewMain
		return m, nil
	default:
		var cmd tea.Cmd
		m.filePicker, cmd = m.filePicker.Update(msg)

		if selected := m.filePicker.Selected(); selected != "" {
			m.imagePath = selected
			m.viewState = ViewMain
			return m.extractColors()
		}

		return m, cmd
	}
}

func (m Model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Escape), key.Matches(msg, m.keys.Help), key.Matches(msg, m.keys.Quit):
		m.help.Hide()
		m.viewState = ViewMain
		return m, nil
	}
	return m, nil
}

func (m Model) updateSaveTheme(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.viewState = ViewMain
		m.themeNameInput = ""
		return m, nil
	case "enter":
		if m.themeNameInput != "" {
			return m.saveTheme()
		}
		return m, nil
	case "backspace":
		if len(m.themeNameInput) > 0 {
			m.themeNameInput = m.themeNameInput[:len(m.themeNameInput)-1]
		}
		return m, nil
	default:
		// Add character to input (only allow valid filename chars)
		key := msg.String()
		if len(key) == 1 && isValidThemeChar(key[0]) {
			m.themeNameInput += key
		}
		return m, nil
	}
}

func isValidThemeChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_'
}

func (m Model) openThemeBrowser() (tea.Model, tea.Cmd) {
	themes, err := config.ListThemes()
	if err != nil {
		m.status = shared.Error("Error listing themes: " + err.Error())
		return m, nil
	}

	if len(themes) == 0 {
		m.status = shared.Info("No themes found. Press 's' to save a theme first.")
		return m, nil
	}

	m.themes = themes
	m.themeCursor = 0
	m.viewState = ViewThemeBrowser
	return m, nil
}

func (m Model) updateThemeBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.viewState = ViewMain
		return m, nil
	case "j", "down":
		if m.themeCursor < len(m.themes)-1 {
			m.themeCursor++
		}
		return m, nil
	case "k", "up":
		if m.themeCursor > 0 {
			m.themeCursor--
		}
		return m, nil
	case "enter":
		// Load the selected theme
		return m.loadSelectedTheme()
	case "a":
		// Apply the selected theme
		return m.applySelectedTheme()
	}
	return m, nil
}

func (m Model) loadSelectedTheme() (tea.Model, tea.Cmd) {
	if len(m.themes) == 0 {
		return m, nil
	}

	themeName := m.themes[m.themeCursor]
	palette, err := config.LoadTheme(themeName)
	if err != nil {
		m.status = shared.Error("Error loading theme: " + err.Error())
		m.viewState = ViewMain
		return m, nil
	}

	m.palette = palette
	m.colorList.SetPalette(m.palette)
	m.preview.SetPalette(m.palette)
	m.status = shared.Success("Loaded theme '" + themeName + "'")
	m.viewState = ViewMain
	return m, nil
}

func (m Model) applySelectedTheme() (tea.Model, tea.Cmd) {
	if len(m.themes) == 0 {
		return m, nil
	}

	themeName := m.themes[m.themeCursor]

	palette, err := config.LoadTheme(themeName)
	if err != nil {
		m.status = shared.Error("Error loading theme: " + err.Error())
		m.viewState = ViewMain
		return m, nil
	}

	if err := config.ApplyTheme(themeName); err != nil {
		m.status = shared.Error("Error applying theme: " + err.Error())
		m.viewState = ViewMain
		return m, nil
	}

	if err := config.ApplyThemeToSystem(palette, m.imagePath); err != nil {
		m.status = shared.Error("Error applying to system: " + err.Error())
		m.viewState = ViewMain
		return m, nil
	}

	m.palette = palette
	m.colorList.SetPalette(m.palette)
	m.preview.SetPalette(m.palette)

	if config.IsOmarchyInstalled() {
		m.status = shared.Success("Applied theme '" + themeName + "' to system")
	} else {
		m.status = shared.Success("Applied theme '" + themeName + "'")
	}
	m.viewState = ViewMain
	return m, nil
}

func (m Model) applyCurrentPalette() (tea.Model, tea.Cmd) {
	if err := config.ApplyThemeToSystem(m.palette, m.imagePath); err != nil {
		m.status = shared.Error("Error applying to system: " + err.Error())
		return m, nil
	}

	if config.IsOmarchyInstalled() {
		m.status = shared.Success("Applied theme to system")
	} else {
		m.status = shared.Success("Generated theme files to " + config.GetPeachyThemeDir())
	}
	return m, nil
}

func (m Model) extractColors() (tea.Model, tea.Cmd) {
	if m.imagePath == "" {
		m.status = shared.Info("No image selected. Press 'o' to open an image.")
		return m, nil
	}

	m.status = shared.Info("Extracting colors...")

	palette, err := m.extractor.ExtractPalette(m.imagePath, m.extractionMode, m.lightMode)
	if err != nil {
		m.status = shared.Error("Error: " + err.Error())
		return m, nil
	}

	m.palette = palette
	m.originalPalette = palette.Clone()
	m.colorList.SetPalette(m.palette)
	m.preview.SetPalette(m.palette)
	m.preview.SetImage(m.imagePath)

	modeName := color.ModeNames[m.extractionMode]
	m.status = shared.Success("Colors extracted (" + modeName + ") from " + filepath.Base(m.imagePath))

	return m, nil
}

func (m Model) saveTheme() (tea.Model, tea.Cmd) {
	name := m.themeNameInput

	if err := config.SaveTheme(name, m.palette); err != nil {
		m.status = shared.Error("Error saving theme: " + err.Error())
		m.viewState = ViewMain
		return m, nil
	}

	m.status = shared.Success("Theme '" + name + "' saved to " + config.GetThemePath(name))
	m.themeNameInput = ""
	m.viewState = ViewMain
	return m, nil
}

func (m *Model) updateComponentSizes() {
	// Guard against zero/invalid sizes before WindowSizeMsg is received
	if m.width < 40 || m.height < 20 {
		return
	}

	// Calculate component sizes - both panels same width
	panelWidth := (m.width - 6) / 2 // Account for gap and margins
	if panelWidth < 30 {
		panelWidth = 30
	}
	if panelWidth > 50 {
		panelWidth = 50
	}

	contentHeight := m.height - 6 // Account for header and help bar
	if contentHeight < 10 {
		contentHeight = 10
	}

	m.colorList.SetSize(panelWidth, contentHeight)
	m.preview.SetSize(panelWidth, contentHeight)
	m.colorEditor.SetSize(panelWidth, contentHeight)
	m.filePicker.SetSize(m.width-10, m.height-10)
	m.help.SetSize(60, 30)
}

// SetImagePath sets the initial image path
func (m *Model) SetImagePath(path string) {
	m.imagePath = path
}

// SetConfigPath sets the config file path
func (m *Model) SetConfigPath(path string) {
	m.configPath = path
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var sb strings.Builder

	// Render based on current view
	switch m.viewState {
	case ViewHelp:
		// Overlay help on main view
		sb.WriteString(m.renderMain())
		sb.WriteString("\n")
		// Center the help overlay
		sb.WriteString(m.help.View())
	case ViewFilePicker:
		sb.WriteString(m.filePicker.View())
	case ViewEditor:
		sb.WriteString(m.renderMainWithEditor())
	case ViewSaveTheme:
		sb.WriteString(m.renderMain())
		sb.WriteString("\n")
		sb.WriteString(m.renderSaveThemeDialog())
	case ViewThemeBrowser:
		sb.WriteString(m.renderThemeBrowser())
	default:
		sb.WriteString(m.renderMain())
	}

	return sb.String()
}

func (m Model) renderMain() string {
	return m.renderLayout(m.colorList.View(), m.preview.View())
}

func (m Model) renderMainWithEditor() string {
	return m.renderLayout(m.colorEditor.View(), m.preview.View())
}

func (m Model) renderSaveThemeDialog() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(HeaderStyle.Render("Save Theme"))
	sb.WriteString("\n\n")
	sb.WriteString("Theme name: ")
	sb.WriteString(m.themeNameInput)
	sb.WriteString("█")
	sb.WriteString("\n\n")
	sb.WriteString(HelpDescStyle.Render("Enter: save  Esc: cancel"))
	return sb.String()
}

func (m Model) renderThemeBrowser() string {
	var sb strings.Builder

	sb.WriteString(HeaderStyle.Render("Peachy - Theme Browser"))
	sb.WriteString("\n\n")

	if len(m.themes) == 0 {
		sb.WriteString("No themes found.\n")
	} else {
		for i, theme := range m.themes {
			if i == m.themeCursor {
				sb.WriteString("> ")
				sb.WriteString(HeaderStyle.Render(theme))
			} else {
				sb.WriteString("  ")
				sb.WriteString(theme)
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(HelpDescStyle.Render("j/k: navigate  Enter: load  a: apply  Esc: cancel"))
	return sb.String()
}

func (m Model) renderLayout(left, right string) string {
	// Build header
	headerTitle := HeaderStyle.Render("Peachy - Theme Creator")

	// Build subtitle line (filename + mode)
	var subtitleParts []string
	if m.imagePath != "" {
		subtitleParts = append(subtitleParts, filepath.Base(m.imagePath))
	}
	modeLabel := "[" + color.ModeNames[m.extractionMode] + "]"
	if m.lightMode {
		modeLabel += " (light)"
	}
	subtitleParts = append(subtitleParts, modeLabel)
	subtitleText := SubtitleStyle.Render(strings.Join(subtitleParts, "  "))

	// Main content - side by side (same as FilePicker)
	content := lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)

	// Status bar
	var statusContent string
	statusBarStyle := lipgloss.NewStyle().Padding(0, 1)

	switch m.status.Type {
	case shared.StatusError:
		errorIcon := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true).
			Render("● ")
		statusContent = statusBarStyle.
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("1")).
			Render(errorIcon + ErrorStyle.Render(m.status.Message))
	case shared.StatusSuccess:
		successIcon := lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true).
			Render("✓ ")
		statusContent = statusBarStyle.
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("2")).
			Render(successIcon + SuccessStyle.Render(m.status.Message))
	case shared.StatusInfo:
		infoIcon := lipgloss.NewStyle().
			Foreground(lipgloss.Color("6")).
			Bold(true).
			Render("● ")
		statusContent = statusBarStyle.
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("6")).
			Render(infoIcon + lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(m.status.Message))
	default:
		statusContent = statusBarStyle.Render(StatusStyle.Render(m.status.Message))
	}

	// Help bar
	helpBar := m.renderHelpBar()

	// Compose everything vertically using lipgloss
	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		headerTitle,
		subtitleText,
		"",
		content,
		"",
		statusContent,
		helpBar,
	)
}

func (m Model) renderHelpBar() string {
	var groups []string

	addKey := func(key, desc string) string {
		return HelpKeyStyle.Render(key) + HelpDescStyle.Render(":"+desc)
	}

	separator := HelpDescStyle.Render(" │ ")

	switch m.viewState {
	case ViewEditor:
		// Group: Adjust values
		adjustGroup := []string{addKey("h/l", "adjust"), addKey("j/k", "field"), addKey("#", "hex")}
		groups = append(groups, strings.Join(adjustGroup, " "))

		// Group: Actions
		actionGroup := []string{addKey("u", "reset"), addKey("Enter", "confirm"), addKey("Esc", "cancel")}
		groups = append(groups, strings.Join(actionGroup, " "))
	default:
		// Group: Navigation
		navGroup := []string{addKey("j/k", "nav"), addKey("Enter", "edit")}
		groups = append(groups, strings.Join(navGroup, " "))

		// Group: Files & Colors
		fileGroup := []string{addKey("o", "open"), addKey("e", "extract"), addKey("m", "mode"), addKey("t", "light/dark")}
		groups = append(groups, strings.Join(fileGroup, " "))

		// Group: Themes
		themeGroup := []string{addKey("s", "save"), addKey("l", "themes"), addKey("a", "apply")}
		groups = append(groups, strings.Join(themeGroup, " "))

		// Group: App
		appGroup := []string{addKey("?", "help"), addKey("q", "quit")}
		groups = append(groups, strings.Join(appGroup, " "))
	}

	return HelpBarStyle.Render(strings.Join(groups, separator))
}
