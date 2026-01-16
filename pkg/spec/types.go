// Package spec defines the data structures for specifications.
package spec

// Spec represents a specification file - pure content for AI interpretation.
type Spec struct {
	RawContent string `json:"rawContent"` // Full markdown content
}

// ValidationError represents an error found during spec validation.
type ValidationError struct {
	Severity string `json:"severity"` // "error", "warning", "info"
	Code     string `json:"code"`
	Message  string `json:"message"`
}

// ValidationResult contains the results of spec validation.
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}
