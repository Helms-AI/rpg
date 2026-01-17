package semantic

// DefaultRegistry creates a new analyzer registry with all analyzers registered
func DefaultRegistry() *AnalyzerRegistry {
	registry := NewAnalyzerRegistry()

	// Register all available analyzers
	registry.Register(NewGoAnalyzer())
	registry.Register(NewTypeScriptAnalyzer())
	registry.Register(NewPythonAnalyzer())
	registry.Register(NewJavaAnalyzer())
	registry.Register(NewRustAnalyzer())
	registry.Register(NewCSharpAnalyzer())

	return registry
}
