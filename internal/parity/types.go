// Package parity provides semantic parity analysis between source and generated code.
package parity

import (
	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// ParityResult contains the full parity analysis results
type ParityResult struct {
	// OverallScore is the weighted parity score (0.0-1.0)
	OverallScore float64 `json:"overallScore"`

	// Converged indicates if parity threshold was met
	Converged bool `json:"converged"`

	// ByDimension contains scores per comparison dimension
	ByDimension DimensionScores `json:"byDimension"`

	// ByLanguage contains per-language analysis
	ByLanguage map[string]LanguageResult `json:"byLanguage"`

	// Gaps contains detailed gap information
	Gaps []ParityGap `json:"gaps"`

	// ConvergenceTrend tracks score progression over iterations
	ConvergenceTrend []float64 `json:"convergenceTrend,omitempty"`
}

// DimensionScores contains scores for each comparison dimension
type DimensionScores struct {
	Structural float64 `json:"structural"` // AST shape similarity
	Type       float64 `json:"type"`       // Type definitions match
	Behavioral float64 `json:"behavioral"` // Function signatures match
	Test       float64 `json:"test"`       // Test cases equivalent
	Idiomatic  float64 `json:"idiomatic"`  // Language conventions followed
}

// LanguageResult contains parity analysis for a single language
type LanguageResult struct {
	Language     treesitter.Language `json:"language"`
	OverallScore float64             `json:"overallScore"`
	ByDimension  DimensionScores     `json:"byDimension"`
	MissingTypes []string            `json:"missingTypes,omitempty"`
	MissingFuncs []string            `json:"missingFuncs,omitempty"`
	TypeErrors   []TypeMismatch      `json:"typeErrors,omitempty"`
	SigErrors    []SignatureMismatch `json:"sigErrors,omitempty"`
}

// ParityGap represents a specific parity issue
type ParityGap struct {
	// Dimension is the parity dimension (structural, type, behavioral, etc.)
	Dimension string `json:"dimension"`

	// Severity is the gap severity (high, medium, low)
	Severity string `json:"severity"`

	// SourceItem describes the source element
	SourceItem ItemReference `json:"sourceItem"`

	// GeneratedItem describes the generated element (if present)
	GeneratedItem *ItemReference `json:"generatedItem,omitempty"`

	// Discrepancy describes what's different
	Discrepancy string `json:"discrepancy"`

	// SuggestedFix provides fix guidance
	SuggestedFix string `json:"suggestedFix"`
}

// ItemReference identifies a code element
type ItemReference struct {
	Type      string `json:"type"`      // "function", "type", "field", etc.
	Name      string `json:"name"`      // Element name
	Signature string `json:"signature"` // Full signature if applicable
	Location  string `json:"location"`  // File:line
	Language  string `json:"language"`  // Language ID
}

// TypeMismatch describes a type definition mismatch
type TypeMismatch struct {
	TypeName    string        `json:"typeName"`
	SourceType  TypeInfo      `json:"sourceType"`
	GenType     TypeInfo      `json:"genType"`
	Differences []string      `json:"differences"`
}

// TypeInfo summarizes a type definition
type TypeInfo struct {
	Kind       string   `json:"kind"`
	FieldCount int      `json:"fieldCount"`
	Methods    []string `json:"methods"`
	Implements []string `json:"implements,omitempty"`
}

// SignatureMismatch describes a function signature mismatch
type SignatureMismatch struct {
	FuncName     string   `json:"funcName"`
	SourceSig    string   `json:"sourceSig"`
	GenSig       string   `json:"genSig"`
	Differences  []string `json:"differences"`
}

// ComparisonConfig configures the parity comparison
type ComparisonConfig struct {
	// Weights for each dimension (should sum to 1.0)
	Weights DimensionWeights `json:"weights"`

	// Threshold is the minimum score to consider converged
	Threshold float64 `json:"threshold"`

	// StrictTypeMatching requires exact type matches
	StrictTypeMatching bool `json:"strictTypeMatching"`

	// IgnorePrivate ignores non-exported elements
	IgnorePrivate bool `json:"ignorePrivate"`
}

// DimensionWeights contains weights for each dimension
type DimensionWeights struct {
	Structural float64 `json:"structural"`
	Type       float64 `json:"type"`
	Behavioral float64 `json:"behavioral"`
	Test       float64 `json:"test"`
	Idiomatic  float64 `json:"idiomatic"`
}

// DefaultConfig returns the default comparison configuration
func DefaultConfig() ComparisonConfig {
	return ComparisonConfig{
		Weights: DimensionWeights{
			Structural: 0.20,
			Type:       0.25,
			Behavioral: 0.35,
			Test:       0.15,
			Idiomatic:  0.05,
		},
		Threshold:          0.95,
		StrictTypeMatching: false,
		IgnorePrivate:      true,
	}
}

// NormalizedSignature represents a cross-language normalized function signature
type NormalizedSignature struct {
	Name       string              `json:"name"`
	Parameters []NormalizedParam   `json:"parameters"`
	Returns    []string            `json:"returns"`
	IsAsync    bool                `json:"isAsync"`
	IsPublic   bool                `json:"isPublic"`
	Complexity int                 `json:"complexity"`
}

// NormalizedParam represents a normalized parameter
type NormalizedParam struct {
	Name     string `json:"name"`
	BaseType string `json:"baseType"` // Normalized to common vocabulary
	IsPtr    bool   `json:"isPtr"`
	IsArray  bool   `json:"isArray"`
	IsMap    bool   `json:"isMap"`
}

// NormalizedType represents a cross-language normalized type
type NormalizedType struct {
	Name       string            `json:"name"`
	Kind       string            `json:"kind"` // struct, interface, enum, alias
	Fields     []NormalizedField `json:"fields"`
	Methods    []string          `json:"methods"`
	Implements []string          `json:"implements"`
	IsPublic   bool              `json:"isPublic"`
}

// NormalizedField represents a normalized field
type NormalizedField struct {
	Name     string `json:"name"`
	BaseType string `json:"baseType"`
	IsPtr    bool   `json:"isPtr"`
	IsArray  bool   `json:"isArray"`
	IsMap    bool   `json:"isMap"`
}

// AnalysisSet contains analyses for comparison
type AnalysisSet struct {
	Source    *semantic.Analysis            `json:"source"`
	Generated map[string]*semantic.Analysis `json:"generated"` // language -> analysis
}
