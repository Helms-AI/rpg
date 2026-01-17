// Package refinement provides the iterative refinement engine for achieving
// high parity between source, spec, and generated code.
package refinement

import (
	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
	"github.com/kon1790/rpg/internal/parity"
)

// LoopConfig configures the refinement loop
type LoopConfig struct {
	// ConvergenceThreshold is the minimum parity score to consider converged
	ConvergenceThreshold float64 `json:"convergenceThreshold"`

	// MaxIterations is the maximum number of refinement iterations
	MaxIterations int `json:"maxIterations"`

	// StuckThreshold is the minimum improvement required per iteration
	StuckThreshold float64 `json:"stuckThreshold"`

	// StuckWindow is the number of iterations to consider for stuck detection
	StuckWindow int `json:"stuckWindow"`

	// RefinementStrategy determines how refinements are applied
	RefinementStrategy Strategy `json:"refinementStrategy"`

	// ParityConfig configures the parity comparison
	ParityConfig parity.ComparisonConfig `json:"parityConfig"`
}

// Strategy represents the refinement strategy
type Strategy string

const (
	// StrategySpecFirst refines the spec before code
	StrategySpecFirst Strategy = "spec-first"

	// StrategyCodeFirst refines code before spec
	StrategyCodeFirst Strategy = "code-first"

	// StrategyBalanced alternates between spec and code refinement
	StrategyBalanced Strategy = "balanced"

	// StrategyAdaptive chooses strategy based on gap analysis
	StrategyAdaptive Strategy = "adaptive"
)

// DefaultLoopConfig returns sensible defaults
func DefaultLoopConfig() LoopConfig {
	return LoopConfig{
		ConvergenceThreshold: 0.95,
		MaxIterations:        5,
		StuckThreshold:       0.02,
		StuckWindow:          3,
		RefinementStrategy:   StrategyBalanced,
		ParityConfig:         parity.DefaultConfig(),
	}
}

// LoopInput contains input for the refinement loop
type LoopInput struct {
	// SourcePath is the path to the source code to analyze
	SourcePath string `json:"sourcePath"`

	// SourceLanguage is the source language (auto-detected if empty)
	SourceLanguage string `json:"sourceLanguage,omitempty"`

	// TargetLanguages are the languages to generate
	TargetLanguages []string `json:"targetLanguages"`

	// SpecPath is the path to the spec file (generated if not exists)
	SpecPath string `json:"specPath,omitempty"`

	// OutputDir is the base directory for generated code
	OutputDir string `json:"outputDir"`
}

// LoopResult contains the result of the refinement loop
type LoopResult struct {
	// Converged indicates if the parity threshold was met
	Converged bool `json:"converged"`

	// FinalScore is the final parity score
	FinalScore float64 `json:"finalScore"`

	// IterationsUsed is the number of iterations performed
	IterationsUsed int `json:"iterationsUsed"`

	// IterationHistory contains details for each iteration
	IterationHistory []Iteration `json:"iterationHistory"`

	// FinalSpec is the path to the final spec
	FinalSpec string `json:"finalSpec,omitempty"`

	// GeneratedProjects maps language to project path
	GeneratedProjects map[string]string `json:"generatedProjects,omitempty"`

	// UnresolvedGaps contains gaps that couldn't be resolved
	UnresolvedGaps []parity.ParityGap `json:"unresolvedGaps,omitempty"`

	// RefinementSummary provides human-readable summary
	RefinementSummary string `json:"refinementSummary"`
}

// Iteration contains details for a single refinement iteration
type Iteration struct {
	// Number is the iteration number (1-based)
	Number int `json:"number"`

	// ParityScore is the parity score after this iteration
	ParityScore float64 `json:"parityScore"`

	// ParityResult contains full parity analysis
	ParityResult *parity.ParityResult `json:"parityResult,omitempty"`

	// RefinementsApplied is the number of refinements applied
	RefinementsApplied int `json:"refinementsApplied"`

	// Refinements contains the refinements applied
	Refinements []Refinement `json:"refinements,omitempty"`

	// Phase indicates what was refined (spec, code, or both)
	Phase string `json:"phase"`

	// ScoreImprovement is the improvement from previous iteration
	ScoreImprovement float64 `json:"scoreImprovement"`
}

// Refinement represents a single refinement action
type Refinement struct {
	// Type is the refinement type (add_function, fix_signature, add_type, etc.)
	Type string `json:"type"`

	// Target is what was refined (language or "spec")
	Target string `json:"target"`

	// ElementName is the name of the element refined
	ElementName string `json:"elementName"`

	// Before is the state before refinement
	Before string `json:"before,omitempty"`

	// After is the state after refinement
	After string `json:"after,omitempty"`

	// Gap is the gap this refinement addresses
	Gap *parity.ParityGap `json:"gap,omitempty"`
}

// AnalysisCache caches semantic analyses to avoid re-parsing
type AnalysisCache struct {
	Source    *semantic.Analysis
	Generated map[string]*semantic.Analysis
}

// NewAnalysisCache creates a new cache
func NewAnalysisCache() *AnalysisCache {
	return &AnalysisCache{
		Generated: make(map[string]*semantic.Analysis),
	}
}

// RefinementInstructions contains instructions for AI to apply refinements
type RefinementInstructions struct {
	// Summary provides overview of needed refinements
	Summary string `json:"summary"`

	// SpecRefinements are changes needed to the spec
	SpecRefinements []SpecChange `json:"specRefinements,omitempty"`

	// CodeRefinements are changes needed to generated code
	CodeRefinements map[string][]CodeChange `json:"codeRefinements,omitempty"`

	// Priority indicates which to do first
	Priority string `json:"priority"` // "spec" or "code"
}

// SpecChange represents a change to the spec
type SpecChange struct {
	Section     string `json:"section"`     // Types, Functions, etc.
	Action      string `json:"action"`      // add, modify, remove
	ElementName string `json:"elementName"` // Name of the element
	Description string `json:"description"` // What to change
	Example     string `json:"example"`     // Example of the change
}

// CodeChange represents a change to generated code
type CodeChange struct {
	File        string `json:"file"`
	Action      string `json:"action"`      // add, modify, remove
	ElementType string `json:"elementType"` // function, type, field
	ElementName string `json:"elementName"`
	Description string `json:"description"`
	SourceRef   string `json:"sourceRef"` // Reference to source implementation
}

// ConvergenceMetrics tracks convergence progress
type ConvergenceMetrics struct {
	// Scores tracks parity scores over iterations
	Scores []float64 `json:"scores"`

	// IsStuck indicates if progress has stalled
	IsStuck bool `json:"isStuck"`

	// StuckReason explains why progress stalled
	StuckReason string `json:"stuckReason,omitempty"`

	// Trend is the recent score trend (positive = improving)
	Trend float64 `json:"trend"`

	// EstimatedIterationsToConverge estimates remaining iterations
	EstimatedIterationsToConverge int `json:"estimatedIterationsToConverge"`
}

// Languages provides language-specific information
type Languages struct {
	available map[treesitter.Language]bool
}

// NewLanguages creates a new Languages tracker
func NewLanguages() *Languages {
	return &Languages{
		available: make(map[treesitter.Language]bool),
	}
}

// SetAvailable marks a language as available
func (l *Languages) SetAvailable(lang treesitter.Language, available bool) {
	l.available[lang] = available
}

// IsAvailable checks if a language is available
func (l *Languages) IsAvailable(lang treesitter.Language) bool {
	return l.available[lang]
}
