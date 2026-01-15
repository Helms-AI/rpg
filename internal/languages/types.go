// Package languages defines language adapters and conventions.
package languages

// Language represents a supported target language with its conventions.
type Language struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	FileExtension    string            `json:"fileExtension"`
	Conventions      Conventions       `json:"conventions"`
	Idioms           []string          `json:"idioms"`
	ProjectStructure ProjectStructure  `json:"projectStructure"`
	ErrorPatterns    ErrorPatterns     `json:"errorPatterns"`
	Dependencies     DependencyInfo    `json:"dependencies"`
}

// Conventions defines naming and style conventions for a language.
type Conventions struct {
	Naming        NamingConventions `json:"naming"`
	ErrorHandling string            `json:"errorHandling"`
	FileNaming    string            `json:"fileNaming"`
	Imports       string            `json:"imports"`
	DocStyle      string            `json:"docStyle"`
}

// NamingConventions defines naming patterns for different identifiers.
type NamingConventions struct {
	Functions  string `json:"functions"`
	Variables  string `json:"variables"`
	Constants  string `json:"constants"`
	Types      string `json:"types"`
	Packages   string `json:"packages"`
	Private    string `json:"private,omitempty"`
}

// ProjectStructure defines the typical project layout for a language.
type ProjectStructure struct {
	SourceDir     string   `json:"sourceDir"`
	TestDir       string   `json:"testDir,omitempty"`
	TestSuffix    string   `json:"testSuffix"`
	PackageFile   string   `json:"packageFile"`
	EntryPoint    string   `json:"entryPoint,omitempty"`
	CommonDirs    []string `json:"commonDirs,omitempty"`
}

// ErrorPatterns defines how errors are handled in the language.
type ErrorPatterns struct {
	Style       string `json:"style"` // "exceptions", "result", "tuple", "optional"
	CustomError string `json:"customError,omitempty"`
	WrapError   string `json:"wrapError,omitempty"`
}

// DependencyInfo provides package manager information.
type DependencyInfo struct {
	Manager       string `json:"manager"`
	InstallCmd    string `json:"installCmd"`
	AddCmd        string `json:"addCmd"`
	LockFile      string `json:"lockFile,omitempty"`
}

// LanguageAdapter provides language-specific behavior.
type LanguageAdapter interface {
	// GetLanguage returns the language configuration.
	GetLanguage() Language

	// GetPromptContext returns language-specific prompt instructions.
	GetPromptContext() string

	// GetProjectStructure returns the recommended project structure for a spec.
	GetProjectStructure(specName string, hasTests bool) []ProjectFile

	// MapType maps a pseudo-code type to the language's equivalent.
	MapType(pseudoType string) string
}

// ProjectFile represents a file in a project structure.
type ProjectFile struct {
	Path        string `json:"path"`
	Purpose     string `json:"purpose"`
	Description string `json:"description,omitempty"`
}
