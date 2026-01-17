// Package generator provides code generation from spec analysis.
package generator

import (
	"time"

	"github.com/kon1790/rpg/internal/parity"
	"github.com/kon1790/rpg/internal/specparser"
)

// GenerateSourceFromSpecInput contains the input parameters for source generation.
type GenerateSourceFromSpecInput struct {
	// SpecPath is the path to the specification file
	SpecPath string `json:"specPath" jsonschema:"required"`

	// Language is the target language for code generation
	Language string `json:"language" jsonschema:"required"`

	// OutputDir is the optional output directory (defaults to outputDir/specName/language)
	OutputDir string `json:"outputDir,omitempty"`

	// MaxIterations limits the parity loop iterations (default: 10)
	MaxIterations int `json:"maxIterations,omitempty"`

	// ParityTarget is the target parity percentage (default: 100.0)
	ParityTarget float64 `json:"parityTarget,omitempty"`
}

// GenerateSourceFromSpecOutput contains the results of source generation.
type GenerateSourceFromSpecOutput struct {
	// Success indicates if target parity was achieved
	Success bool `json:"success"`

	// Iterations is the number of parity loop iterations performed
	Iterations int `json:"iterations"`

	// FinalParity is the final parity percentage achieved
	FinalParity float64 `json:"finalParity"`

	// SpecAnalysis contains the parsed spec information
	SpecAnalysis SpecAnalysisSummary `json:"specAnalysis"`

	// GeneratedFiles lists all files that were generated
	GeneratedFiles []GeneratedFile `json:"generatedFiles"`

	// OutputDir is the directory where files were written
	OutputDir string `json:"outputDir"`

	// ParityReport contains the final parity analysis
	ParityReport ParityReportSummary `json:"parityReport"`

	// History contains the iteration history
	History []IterationResult `json:"history"`

	// GenerationPrompt provides AI instructions for completing generation
	GenerationPrompt string `json:"generationPrompt,omitempty"`
}

// SpecAnalysisSummary provides a summary of the parsed spec.
type SpecAnalysisSummary struct {
	// Name is the spec name
	Name string `json:"name"`

	// Overview is the spec description
	Overview string `json:"overview"`

	// TypeCount is the number of types defined
	TypeCount int `json:"typeCount"`

	// FunctionCount is the number of functions defined
	FunctionCount int `json:"functionCount"`

	// TestCount is the number of tests defined
	TestCount int `json:"testCount"`

	// DependencyCount is the number of dependencies
	DependencyCount int `json:"dependencyCount"`

	// ConfigCount is the number of configuration items
	ConfigCount int `json:"configCount"`

	// TotalItems is the sum of all items
	TotalItems int `json:"totalItems"`

	// Types contains the type definitions
	Types []specparser.SpecType `json:"types,omitempty"`

	// Functions contains the function definitions
	Functions []specparser.SpecFunction `json:"functions,omitempty"`

	// Tests contains the test definitions
	Tests []specparser.SpecTest `json:"tests,omitempty"`
}

// GeneratedFile represents a single generated file.
type GeneratedFile struct {
	// Path is the relative path within the output directory
	Path string `json:"path"`

	// Content is the generated file content
	Content string `json:"content"`

	// Size is the content size in bytes
	Size int `json:"size"`

	// Category is the type of file (type, function, test, config)
	Category string `json:"category"`

	// Elements lists the spec elements this file implements
	Elements []string `json:"elements,omitempty"`
}

// ParityReportSummary provides a summary of parity analysis.
type ParityReportSummary struct {
	// OverallScore is the overall parity score (0.0-1.0)
	OverallScore float64 `json:"overallScore"`

	// Converged indicates if parity target was reached
	Converged bool `json:"converged"`

	// ByDimension contains scores per comparison dimension
	ByDimension DimensionScores `json:"byDimension"`

	// GapCount is the number of remaining gaps
	GapCount int `json:"gapCount"`

	// Gaps contains detailed gap information
	Gaps []ParityGapSummary `json:"gaps,omitempty"`
}

// DimensionScores contains scores for each parity dimension.
type DimensionScores struct {
	Structural float64 `json:"structural"`
	Type       float64 `json:"type"`
	Behavioral float64 `json:"behavioral"`
	Test       float64 `json:"test"`
	Idiomatic  float64 `json:"idiomatic"`
}

// ParityGapSummary provides a summary of a parity gap.
type ParityGapSummary struct {
	// Dimension is the parity dimension (type, behavioral, etc.)
	Dimension string `json:"dimension"`

	// Severity is the gap severity (high, medium, low)
	Severity string `json:"severity"`

	// SourceItem is what's expected from the spec
	SourceItem string `json:"sourceItem"`

	// Discrepancy describes what's missing or different
	Discrepancy string `json:"discrepancy"`

	// SuggestedFix provides guidance on how to fix
	SuggestedFix string `json:"suggestedFix"`
}

// IterationResult contains details of a single parity loop iteration.
type IterationResult struct {
	// Iteration is the iteration number (1-based)
	Iteration int `json:"iteration"`

	// ParityScore is the parity score after this iteration
	ParityScore float64 `json:"parityScore"`

	// GapsFixed is the number of gaps fixed in this iteration
	GapsFixed int `json:"gapsFixed"`

	// GapsRemaining is the number of gaps still remaining
	GapsRemaining int `json:"gapsRemaining"`

	// FilesModified lists files that were modified
	FilesModified []string `json:"filesModified,omitempty"`

	// Duration is how long this iteration took
	Duration time.Duration `json:"duration"`
}

// GenerationContext contains all context needed for code generation.
type GenerationContext struct {
	// Spec is the parsed specification
	Spec *specparser.SpecAnalysis

	// Language is the target language ID
	Language string

	// LanguageConventions contains language-specific conventions
	LanguageConventions string

	// PromptTemplate is the language-specific generation prompt
	PromptTemplate string

	// OutputDir is the target output directory
	OutputDir string

	// ProjectStructure is the recommended file structure
	ProjectStructure []ProjectFile
}

// ProjectFile represents a file in the project structure.
type ProjectFile struct {
	Path        string `json:"path"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// GapFixContext contains context for fixing parity gaps.
type GapFixContext struct {
	// Gap is the gap to fix
	Gap parity.ParityGap

	// Spec is the specification
	Spec *specparser.SpecAnalysis

	// Language is the target language
	Language string

	// OutputDir is the output directory
	OutputDir string
}
