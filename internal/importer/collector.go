package importer

import (
	"os"
	"path/filepath"
	"strings"
)

// FileCategory represents the type of file collected
type FileCategory string

const (
	CategorySource FileCategory = "source"
	CategoryTest   FileCategory = "test"
	CategoryAPI    FileCategory = "api_spec"
	CategoryConfig FileCategory = "config"
	CategoryDoc    FileCategory = "doc"
)

// FileContent holds the path and content of a collected file
type FileContent struct {
	Path     string       `json:"path"`
	Category FileCategory `json:"category"`
	Content  string       `json:"content"`
}

// ProjectFiles contains all collected files from a project
type ProjectFiles struct {
	Name        string        `json:"name"`
	Language    string        `json:"language"`
	RootPath    string        `json:"rootPath"`
	SourceFiles []FileContent `json:"sourceFiles"`
	TestFiles   []FileContent `json:"testFiles"`
	APISpecs    []FileContent `json:"apiSpecs"`
	ConfigFiles []FileContent `json:"configFiles"`
	DocFiles    []FileContent `json:"docFiles"`
}

// Source file extensions by language
var sourceExtensions = map[string][]string{
	".go":    {"go"},
	".ts":    {"typescript"},
	".tsx":   {"typescript"},
	".js":    {"javascript"},
	".jsx":   {"javascript"},
	".py":    {"python"},
	".java":  {"java"},
	".rs":    {"rust"},
	".cs":    {"csharp"},
	".kt":    {"kotlin"},
	".swift": {"swift"},
	".rb":    {"ruby"},
	".php":   {"php"},
	".scala": {"scala"},
	".c":     {"c"},
	".cpp":   {"cpp"},
	".h":     {"c"},
	".hpp":   {"cpp"},
}

// Test file patterns
var testPatterns = []string{
	"_test.go",
	".test.ts",
	".test.tsx",
	".test.js",
	".test.jsx",
	".spec.ts",
	".spec.tsx",
	".spec.js",
	".spec.jsx",
	"test_",
	"_test.py",
	"Test.java",
	"Tests.java",
	"_test.rs",
	"Tests.cs",
	"Test.cs",
	"_test.rb",
	"_spec.rb",
}

// API specification files
var apiSpecFiles = []string{
	"openapi.yaml",
	"openapi.yml",
	"openapi.json",
	"swagger.yaml",
	"swagger.yml",
	"swagger.json",
	"asyncapi.yaml",
	"asyncapi.yml",
	"asyncapi.json",
	"api.yaml",
	"api.yml",
	"api.json",
	"spec.yaml",
	"spec.yml",
	"spec.json",
}

// Configuration files
var configFiles = []string{
	"go.mod",
	"go.sum",
	"package.json",
	"package-lock.json",
	"tsconfig.json",
	"requirements.txt",
	"pyproject.toml",
	"setup.py",
	"Cargo.toml",
	"Cargo.lock",
	"pom.xml",
	"build.gradle",
	"build.gradle.kts",
	"settings.gradle",
	"settings.gradle.kts",
	".csproj",
	".sln",
	".fsproj",
	".env",
	".env.example",
	".env.sample",
	"config.yaml",
	"config.yml",
	"config.json",
	"application.properties",
	"application.yml",
	"application.yaml",
	"appsettings.json",
	"Makefile",
	"Dockerfile",
	"docker-compose.yml",
	"docker-compose.yaml",
}

// Documentation files
var docFiles = []string{
	"README.md",
	"README",
	"readme.md",
	"CONTRIBUTING.md",
	"API.md",
	"CHANGELOG.md",
	"ARCHITECTURE.md",
	"docs/",
}

// Directories to skip during collection
var skipDirs = []string{
	".git",
	"node_modules",
	"vendor",
	"target",
	"build",
	"dist",
	"out",
	".idea",
	".vscode",
	"__pycache__",
	".pytest_cache",
	".mypy_cache",
	".next",
	".nuxt",
	"coverage",
	".gradle",
	".mvn",
	"bin",
	"obj",
}

// CollectProjectFiles scans a directory and collects all relevant files with their contents
func CollectProjectFiles(dir string) (*ProjectFiles, error) {
	project := &ProjectFiles{
		Name:        filepath.Base(dir),
		RootPath:    dir,
		SourceFiles: []FileContent{},
		TestFiles:   []FileContent{},
		APISpecs:    []FileContent{},
		ConfigFiles: []FileContent{},
		DocFiles:    []FileContent{},
	}

	// Track language detection
	langCounts := make(map[string]int)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Get relative path
		relPath, _ := filepath.Rel(dir, path)

		// Skip hidden directories and known skip directories
		if info.IsDir() {
			baseName := filepath.Base(path)
			if strings.HasPrefix(baseName, ".") && baseName != "." {
				return filepath.SkipDir
			}
			for _, skip := range skipDirs {
				if baseName == skip {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Skip hidden files
		baseName := filepath.Base(path)
		if strings.HasPrefix(baseName, ".") && !strings.HasPrefix(baseName, ".env") {
			return nil
		}

		// Skip very large files (> 1MB)
		if info.Size() > 1024*1024 {
			return nil
		}

		// Categorize and collect the file
		category := CategorizeFile(relPath)
		if category == "" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		fc := FileContent{
			Path:     relPath,
			Category: category,
			Content:  string(content),
		}

		switch category {
		case CategorySource:
			project.SourceFiles = append(project.SourceFiles, fc)
			// Track language
			ext := filepath.Ext(path)
			if langs, ok := sourceExtensions[ext]; ok {
				for _, lang := range langs {
					langCounts[lang]++
				}
			}
		case CategoryTest:
			project.TestFiles = append(project.TestFiles, fc)
		case CategoryAPI:
			project.APISpecs = append(project.APISpecs, fc)
		case CategoryConfig:
			project.ConfigFiles = append(project.ConfigFiles, fc)
		case CategoryDoc:
			project.DocFiles = append(project.DocFiles, fc)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Determine primary language
	maxCount := 0
	for lang, count := range langCounts {
		if count > maxCount {
			maxCount = count
			project.Language = lang
		}
	}

	return project, nil
}

// CategorizeFile determines the category of a file based on its path and name
func CategorizeFile(path string) FileCategory {
	baseName := filepath.Base(path)
	lowerPath := strings.ToLower(path)
	lowerName := strings.ToLower(baseName)
	ext := filepath.Ext(path)

	// Check for API specs first (exact filename match)
	for _, apiFile := range apiSpecFiles {
		if lowerName == apiFile || strings.HasSuffix(lowerPath, "/"+apiFile) {
			return CategoryAPI
		}
	}

	// Check for documentation files
	for _, docFile := range docFiles {
		if strings.HasSuffix(docFile, "/") {
			// Directory pattern
			if strings.Contains(lowerPath, strings.TrimSuffix(docFile, "/")) {
				if ext == ".md" || ext == ".txt" || ext == ".rst" {
					return CategoryDoc
				}
			}
		} else if lowerName == strings.ToLower(docFile) {
			return CategoryDoc
		}
	}

	// Check for config files
	for _, cfgFile := range configFiles {
		if strings.HasSuffix(cfgFile, ".csproj") || strings.HasSuffix(cfgFile, ".sln") || strings.HasSuffix(cfgFile, ".fsproj") {
			if strings.HasSuffix(lowerName, strings.ToLower(cfgFile)) {
				return CategoryConfig
			}
		} else if lowerName == strings.ToLower(cfgFile) {
			return CategoryConfig
		}
	}

	// Check for test files
	for _, testPat := range testPatterns {
		if strings.HasPrefix(testPat, "_") || strings.HasPrefix(testPat, ".") {
			// Suffix pattern
			if strings.HasSuffix(lowerName, strings.ToLower(testPat)) {
				return CategoryTest
			}
		} else if strings.HasSuffix(testPat, "_") {
			// Prefix pattern
			if strings.HasPrefix(lowerName, strings.ToLower(testPat)) {
				return CategoryTest
			}
		} else {
			// Contains pattern (like Test.java)
			if strings.Contains(baseName, testPat) {
				return CategoryTest
			}
		}
	}

	// Check for source files
	if _, ok := sourceExtensions[ext]; ok {
		return CategorySource
	}

	return ""
}

// IsRelevantFile checks if a file should be collected
func IsRelevantFile(path string) bool {
	return CategorizeFile(path) != ""
}

// GetTotalFileCount returns the total number of files collected
func (p *ProjectFiles) GetTotalFileCount() int {
	return len(p.SourceFiles) + len(p.TestFiles) + len(p.APISpecs) + len(p.ConfigFiles) + len(p.DocFiles)
}

// GetTotalContentSize returns the total size of all collected content
func (p *ProjectFiles) GetTotalContentSize() int {
	total := 0
	for _, f := range p.SourceFiles {
		total += len(f.Content)
	}
	for _, f := range p.TestFiles {
		total += len(f.Content)
	}
	for _, f := range p.APISpecs {
		total += len(f.Content)
	}
	for _, f := range p.ConfigFiles {
		total += len(f.Content)
	}
	for _, f := range p.DocFiles {
		total += len(f.Content)
	}
	return total
}
