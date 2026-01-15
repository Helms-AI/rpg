package languages

import (
	"fmt"
	"strings"
)

// Registry holds all registered language adapters.
type Registry struct {
	adapters map[string]LanguageAdapter
}

// NewRegistry creates a new language registry with all built-in adapters.
func NewRegistry() *Registry {
	r := &Registry{
		adapters: make(map[string]LanguageAdapter),
	}

	// Register all built-in language adapters
	r.Register(NewGoAdapter())
	r.Register(NewRustAdapter())
	r.Register(NewJavaAdapter())
	r.Register(NewPythonAdapter())
	r.Register(NewTypeScriptAdapter())
	r.Register(NewCSharpAdapter())

	return r
}

// Register adds a language adapter to the registry.
func (r *Registry) Register(adapter LanguageAdapter) {
	lang := adapter.GetLanguage()
	r.adapters[strings.ToLower(lang.ID)] = adapter
}

// Get returns a language adapter by ID.
func (r *Registry) Get(languageID string) (LanguageAdapter, error) {
	adapter, ok := r.adapters[strings.ToLower(languageID)]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", languageID)
	}
	return adapter, nil
}

// List returns all registered languages.
func (r *Registry) List() []Language {
	languages := make([]Language, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		languages = append(languages, adapter.GetLanguage())
	}
	return languages
}

// IDs returns all registered language IDs.
func (r *Registry) IDs() []string {
	ids := make([]string, 0, len(r.adapters))
	for id := range r.adapters {
		ids = append(ids, id)
	}
	return ids
}
