// Package parser provides spec file reading functionality.
package parser

import (
	"strings"

	"github.com/kon1790/rpg/pkg/spec"
)

// Parse reads a spec file and returns it as raw content for AI interpretation.
func Parse(content string) (*spec.Spec, error) {
	return &spec.Spec{RawContent: content}, nil
}

// Validate checks if the content is valid (non-empty).
func Validate(content string) spec.ValidationResult {
	if strings.TrimSpace(content) == "" {
		return spec.ValidationResult{
			Valid: false,
			Errors: []spec.ValidationError{{
				Severity: "error",
				Code:     "EMPTY_CONTENT",
				Message:  "Spec file is empty",
			}},
		}
	}
	return spec.ValidationResult{Valid: true}
}
