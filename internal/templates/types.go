package templates

// TemplateManifest represents a template.toml file
type TemplateManifest struct {
	Name        string         `toml:"name"`
	Description string         `toml:"description"`
	Version     string         `toml:"version"`
	Condition   string         `toml:"condition"` // Optional: command that must exist
	Disabled    bool           `toml:"disabled"`  // Optional: skip this template
	Files       []TemplateFile `toml:"files"`
}

// TemplateFile represents a single file mapping in the manifest
type TemplateFile struct {
	Template    string `toml:"template"`    // Source template filename
	Destination string `toml:"destination"` // Target path (supports ~)
}

// Template represents a discovered template with its manifest and path
type Template struct {
	Manifest TemplateManifest
	Path     string // Directory path containing the template
	Name     string // Directory name (e.g., "cava")
}

// ProcessResult represents the result of processing a template
type ProcessResult struct {
	Template    *Template
	Success     bool
	FilesLinked []string
	Error       error
	Skipped     bool
	SkipReason  string
}
