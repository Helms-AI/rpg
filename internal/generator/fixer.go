package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/languages"
	"github.com/kon1790/rpg/internal/parity"
	"github.com/kon1790/rpg/internal/specparser"
)

// Fixer handles fixing parity gaps in generated code.
type Fixer struct {
	generator *Generator
}

// NewFixer creates a new gap fixer.
func NewFixer(registry *languages.Registry) *Fixer {
	return &Fixer{
		generator: NewGenerator(registry),
	}
}

// FixGaps attempts to fix parity gaps by generating missing code.
func (f *Fixer) FixGaps(gaps []parity.ParityGap, spec *specparser.SpecAnalysis, language, outputDir string) ([]string, int, error) {
	if len(gaps) == 0 {
		return nil, 0, nil
	}

	adapter, err := f.generator.registry.Get(language)
	if err != nil {
		return nil, 0, fmt.Errorf("unsupported language %s: %w", language, err)
	}

	lang := adapter.GetLanguage()
	var modifiedFiles []string
	fixed := 0

	// Group gaps by dimension for organized fixing
	typeGaps := filterGapsByDimension(gaps, "type")
	behavioralGaps := filterGapsByDimension(gaps, "behavioral")
	testGaps := filterGapsByDimension(gaps, "test")

	// Fix type gaps
	if len(typeGaps) > 0 {
		files, count := f.fixTypeGaps(typeGaps, spec, lang, outputDir)
		modifiedFiles = append(modifiedFiles, files...)
		fixed += count
	}

	// Fix behavioral (function) gaps
	if len(behavioralGaps) > 0 {
		files, count := f.fixBehavioralGaps(behavioralGaps, spec, lang, outputDir)
		modifiedFiles = append(modifiedFiles, files...)
		fixed += count
	}

	// Fix test gaps
	if len(testGaps) > 0 {
		files, count := f.fixTestGaps(testGaps, spec, lang, outputDir)
		modifiedFiles = append(modifiedFiles, files...)
		fixed += count
	}

	return modifiedFiles, fixed, nil
}

// fixTypeGaps fixes missing or incorrect type definitions.
func (f *Fixer) fixTypeGaps(gaps []parity.ParityGap, spec *specparser.SpecAnalysis, lang languages.Language, outputDir string) ([]string, int) {
	var modifiedFiles []string
	fixed := 0

	// Collect types that need to be added or fixed
	var typesToFix []specparser.SpecType
	for _, gap := range gaps {
		if gap.SourceItem.Type != "type" {
			continue
		}

		// Find the type in spec
		specType := findTypeInSpec(spec, gap.SourceItem.Name)
		if specType != nil {
			typesToFix = append(typesToFix, *specType)
		}
	}

	if len(typesToFix) == 0 {
		return modifiedFiles, fixed
	}

	// Determine target file path
	var targetPath string
	switch lang.ID {
	case "go":
		targetPath = filepath.Join(outputDir, "types.go")
	case "typescript":
		targetPath = filepath.Join(outputDir, "src", "types.ts")
	case "python":
		targetPath = filepath.Join(outputDir, "src", "types.py")
	case "java":
		targetPath = filepath.Join(outputDir, "src", "main", "java", toPackageName(spec.Name), "Types.java")
	case "rust":
		targetPath = filepath.Join(outputDir, "src", "types.rs")
	case "csharp":
		targetPath = filepath.Join(outputDir, "src", "Types.cs")
	default:
		return modifiedFiles, fixed
	}

	// Read existing file content
	existingContent, _ := os.ReadFile(targetPath)

	// Generate new type definitions
	var newTypes strings.Builder
	for _, t := range typesToFix {
		// Check if type already exists in file
		if strings.Contains(string(existingContent), t.Name) {
			continue
		}

		typeCode := f.generator.generateType(t, lang)
		newTypes.WriteString(typeCode)
		newTypes.WriteString("\n")
		fixed++
	}

	if newTypes.Len() == 0 {
		return modifiedFiles, fixed
	}

	// Append to file or create new
	var finalContent string
	if len(existingContent) > 0 {
		// Find appropriate place to insert (before closing brace for C#)
		content := string(existingContent)
		if lang.ID == "csharp" && strings.Contains(content, "}") {
			// Insert before last closing brace
			lastBrace := strings.LastIndex(content, "}")
			finalContent = content[:lastBrace] + newTypes.String() + content[lastBrace:]
		} else {
			// Append at end
			finalContent = content + "\n" + newTypes.String()
		}
	} else {
		// Create new file with proper header
		var header strings.Builder
		switch lang.ID {
		case "go":
			header.WriteString(fmt.Sprintf("package %s\n\n", toPackageName(spec.Name)))
		case "typescript":
			header.WriteString("// Type definitions\n\n")
		case "python":
			header.WriteString("from dataclasses import dataclass\nfrom typing import Optional, List, Dict, Any\n\n")
		case "java":
			header.WriteString(fmt.Sprintf("package %s;\n\n", toPackageName(spec.Name)))
		case "rust":
			header.WriteString("use serde::{Deserialize, Serialize};\n\n")
		case "csharp":
			header.WriteString(fmt.Sprintf("namespace %s\n{\n", toPascalCase(spec.Name)))
		}
		finalContent = header.String() + newTypes.String()
		if lang.ID == "csharp" {
			finalContent += "}\n"
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return modifiedFiles, 0
	}

	// Write file
	if err := os.WriteFile(targetPath, []byte(finalContent), 0644); err != nil {
		return modifiedFiles, 0
	}

	modifiedFiles = append(modifiedFiles, targetPath)
	return modifiedFiles, fixed
}

// fixBehavioralGaps fixes missing or incorrect functions.
func (f *Fixer) fixBehavioralGaps(gaps []parity.ParityGap, spec *specparser.SpecAnalysis, lang languages.Language, outputDir string) ([]string, int) {
	var modifiedFiles []string
	fixed := 0

	// Collect functions that need to be added or fixed
	var funcsToFix []specparser.SpecFunction
	for _, gap := range gaps {
		if gap.SourceItem.Type != "function" {
			continue
		}

		// Find the function in spec
		specFunc := findFunctionInSpec(spec, gap.SourceItem.Name)
		if specFunc != nil {
			funcsToFix = append(funcsToFix, *specFunc)
		}
	}

	if len(funcsToFix) == 0 {
		return modifiedFiles, fixed
	}

	// Determine target file path
	var targetPath string
	switch lang.ID {
	case "go":
		targetPath = filepath.Join(outputDir, "service.go")
	case "typescript":
		targetPath = filepath.Join(outputDir, "src", "service.ts")
	case "python":
		targetPath = filepath.Join(outputDir, "src", "service.py")
	case "java":
		targetPath = filepath.Join(outputDir, "src", "main", "java", toPackageName(spec.Name), "Service.java")
	case "rust":
		targetPath = filepath.Join(outputDir, "src", "service.rs")
	case "csharp":
		targetPath = filepath.Join(outputDir, "src", "Service.cs")
	default:
		return modifiedFiles, fixed
	}

	// Read existing file content
	existingContent, _ := os.ReadFile(targetPath)

	// Generate new function definitions
	var newFuncs strings.Builder
	for _, fn := range funcsToFix {
		// Check if function already exists in file
		if strings.Contains(string(existingContent), fn.Name) {
			continue
		}

		funcCode := f.generator.generateFunction(fn, lang)
		newFuncs.WriteString(funcCode)
		newFuncs.WriteString("\n")
		fixed++
	}

	if newFuncs.Len() == 0 {
		return modifiedFiles, fixed
	}

	// Append to file or create new
	var finalContent string
	if len(existingContent) > 0 {
		content := string(existingContent)
		// For languages with class/namespace wrappers, insert before closing brace
		if (lang.ID == "java" || lang.ID == "csharp") && strings.Contains(content, "}") {
			// Find the appropriate place to insert (before last closing brace)
			lastBrace := strings.LastIndex(content, "}")
			if lang.ID == "csharp" {
				// C# has nested braces
				secondLastBrace := strings.LastIndex(content[:lastBrace], "}")
				if secondLastBrace > 0 {
					lastBrace = secondLastBrace
				}
			}
			finalContent = content[:lastBrace] + newFuncs.String() + content[lastBrace:]
		} else {
			finalContent = content + "\n" + newFuncs.String()
		}
	} else {
		// Create new file with proper header
		var header strings.Builder
		switch lang.ID {
		case "go":
			header.WriteString(fmt.Sprintf("package %s\n\n", toPackageName(spec.Name)))
		case "typescript":
			header.WriteString("// Functions\n\n")
		case "python":
			header.WriteString("from typing import Optional, Any\n\n")
		case "java":
			header.WriteString(fmt.Sprintf("package %s;\n\npublic class Service {\n", toPackageName(spec.Name)))
		case "rust":
			header.WriteString("use crate::types::*;\n\n")
		case "csharp":
			header.WriteString(fmt.Sprintf("namespace %s\n{\n    public static class Service\n    {\n", toPascalCase(spec.Name)))
		}
		finalContent = header.String() + newFuncs.String()
		if lang.ID == "java" {
			finalContent += "}\n"
		} else if lang.ID == "csharp" {
			finalContent += "    }\n}\n"
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return modifiedFiles, 0
	}

	// Write file
	if err := os.WriteFile(targetPath, []byte(finalContent), 0644); err != nil {
		return modifiedFiles, 0
	}

	modifiedFiles = append(modifiedFiles, targetPath)
	return modifiedFiles, fixed
}

// fixTestGaps fixes missing tests.
func (f *Fixer) fixTestGaps(gaps []parity.ParityGap, spec *specparser.SpecAnalysis, lang languages.Language, outputDir string) ([]string, int) {
	var modifiedFiles []string
	fixed := 0

	// Collect tests that need to be added
	var testsToFix []specparser.SpecTest
	for _, gap := range gaps {
		// Find the test in spec
		specTest := findTestInSpec(spec, gap.SourceItem.Name)
		if specTest != nil {
			testsToFix = append(testsToFix, *specTest)
		}
	}

	if len(testsToFix) == 0 {
		return modifiedFiles, fixed
	}

	// Determine target file path
	var targetPath string
	switch lang.ID {
	case "go":
		targetPath = filepath.Join(outputDir, "service_test.go")
	case "typescript":
		targetPath = filepath.Join(outputDir, "src", "service.test.ts")
	case "python":
		targetPath = filepath.Join(outputDir, "tests", "test_service.py")
	case "java":
		targetPath = filepath.Join(outputDir, "src", "test", "java", toPackageName(spec.Name), "Tests.java")
	case "rust":
		targetPath = filepath.Join(outputDir, "src", "tests.rs")
	case "csharp":
		targetPath = filepath.Join(outputDir, "tests", "Tests.cs")
	default:
		return modifiedFiles, fixed
	}

	// Read existing file content
	existingContent, _ := os.ReadFile(targetPath)

	// Generate new test definitions
	var newTests strings.Builder
	for _, t := range testsToFix {
		testName := toPascalCase(strings.ReplaceAll(t.Name, " ", "_"))
		// Check if test already exists in file
		if strings.Contains(string(existingContent), testName) {
			continue
		}

		testCode := f.generator.generateTest(t, lang)
		newTests.WriteString(testCode)
		newTests.WriteString("\n")
		fixed++
	}

	if newTests.Len() == 0 {
		return modifiedFiles, fixed
	}

	// Append to file or create new
	var finalContent string
	if len(existingContent) > 0 {
		content := string(existingContent)
		// For languages with class/module wrappers, insert before closing brace
		if (lang.ID == "java" || lang.ID == "csharp" || lang.ID == "rust") && strings.Contains(content, "}") {
			lastBrace := strings.LastIndex(content, "}")
			finalContent = content[:lastBrace] + newTests.String() + content[lastBrace:]
		} else {
			finalContent = content + "\n" + newTests.String()
		}
	} else {
		// Create new file with proper header
		var header strings.Builder
		switch lang.ID {
		case "go":
			header.WriteString(fmt.Sprintf("package %s\n\nimport \"testing\"\n\n", toPackageName(spec.Name)))
		case "typescript":
			header.WriteString("import { describe, it, expect } from 'vitest';\n\n")
		case "python":
			header.WriteString("import pytest\n\n")
		case "java":
			header.WriteString(fmt.Sprintf("package %s;\n\nimport org.junit.jupiter.api.Test;\nimport static org.junit.jupiter.api.Assertions.*;\n\npublic class Tests {\n", toPackageName(spec.Name)))
		case "rust":
			header.WriteString("#[cfg(test)]\nmod tests {\n    use super::*;\n\n")
		case "csharp":
			header.WriteString(fmt.Sprintf("namespace %s.Tests\n{\n    using Xunit;\n\n    public class Tests\n    {\n", toPascalCase(spec.Name)))
		}
		finalContent = header.String() + newTests.String()
		switch lang.ID {
		case "java":
			finalContent += "}\n"
		case "rust":
			finalContent += "}\n"
		case "csharp":
			finalContent += "    }\n}\n"
		}
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return modifiedFiles, 0
	}

	// Write file
	if err := os.WriteFile(targetPath, []byte(finalContent), 0644); err != nil {
		return modifiedFiles, 0
	}

	modifiedFiles = append(modifiedFiles, targetPath)
	return modifiedFiles, fixed
}

// GenerateFixInstructions generates human-readable fix instructions for gaps.
func GenerateFixInstructions(gaps []parity.ParityGap, spec *specparser.SpecAnalysis, language string) string {
	if len(gaps) == 0 {
		return "All code has been generated successfully. No fixes needed."
	}

	var sb strings.Builder
	sb.WriteString("# Code Generation Fix Instructions\n\n")
	sb.WriteString(fmt.Sprintf("**Target Language:** %s\n", language))
	sb.WriteString(fmt.Sprintf("**Remaining Gaps:** %d\n\n", len(gaps)))

	// Group by severity
	high := filterGapsBySeverity(gaps, "high")
	medium := filterGapsBySeverity(gaps, "medium")
	low := filterGapsBySeverity(gaps, "low")

	if len(high) > 0 {
		sb.WriteString("## Critical (High Severity)\n\n")
		for _, g := range high {
			writeGapInstruction(&sb, g, spec)
		}
	}

	if len(medium) > 0 {
		sb.WriteString("## Important (Medium Severity)\n\n")
		for _, g := range medium {
			writeGapInstruction(&sb, g, spec)
		}
	}

	if len(low) > 0 {
		sb.WriteString("## Minor (Low Severity)\n\n")
		for _, g := range low {
			writeGapInstruction(&sb, g, spec)
		}
	}

	sb.WriteString("\n---\n")
	sb.WriteString("Generate the missing code to achieve full parity with the specification.\n")

	return sb.String()
}

func writeGapInstruction(sb *strings.Builder, gap parity.ParityGap, spec *specparser.SpecAnalysis) {
	sb.WriteString(fmt.Sprintf("### %s: `%s`\n\n", gap.Dimension, gap.SourceItem.Name))
	sb.WriteString(fmt.Sprintf("- **Issue:** %s\n", gap.Discrepancy))
	sb.WriteString(fmt.Sprintf("- **Fix:** %s\n", gap.SuggestedFix))

	// Add spec reference
	switch gap.SourceItem.Type {
	case "type":
		if t := findTypeInSpec(spec, gap.SourceItem.Name); t != nil {
			sb.WriteString("- **Spec Definition:**\n")
			sb.WriteString(fmt.Sprintf("  - Kind: %s\n", t.Kind))
			if len(t.Fields) > 0 {
				sb.WriteString("  - Fields:\n")
				for _, f := range t.Fields {
					sb.WriteString(fmt.Sprintf("    - `%s: %s`\n", f.Name, f.Type))
				}
			}
		}
	case "function":
		if fn := findFunctionInSpec(spec, gap.SourceItem.Name); fn != nil {
			sb.WriteString("- **Spec Definition:**\n")
			if len(fn.Parameters) > 0 {
				sb.WriteString("  - Parameters:\n")
				for _, p := range fn.Parameters {
					sb.WriteString(fmt.Sprintf("    - `%s: %s`\n", p.Name, p.Type))
				}
			}
			if len(fn.Returns) > 0 {
				sb.WriteString(fmt.Sprintf("  - Returns: `%s`\n", fn.Returns[0].Type))
			}
			if fn.Logic != "" {
				sb.WriteString(fmt.Sprintf("  - Logic: %s\n", fn.Logic))
			}
		}
	}

	sb.WriteString("\n")
}

// Helper functions

func filterGapsByDimension(gaps []parity.ParityGap, dimension string) []parity.ParityGap {
	var result []parity.ParityGap
	for _, g := range gaps {
		if g.Dimension == dimension {
			result = append(result, g)
		}
	}
	return result
}

func filterGapsBySeverity(gaps []parity.ParityGap, severity string) []parity.ParityGap {
	var result []parity.ParityGap
	for _, g := range gaps {
		if g.Severity == severity {
			result = append(result, g)
		}
	}
	return result
}

func findTypeInSpec(spec *specparser.SpecAnalysis, name string) *specparser.SpecType {
	for i, t := range spec.Types {
		if t.Name == name {
			return &spec.Types[i]
		}
	}
	return nil
}

func findFunctionInSpec(spec *specparser.SpecAnalysis, name string) *specparser.SpecFunction {
	for i, fn := range spec.Functions {
		if fn.Name == name {
			return &spec.Functions[i]
		}
	}
	return nil
}

func findTestInSpec(spec *specparser.SpecAnalysis, name string) *specparser.SpecTest {
	for i, t := range spec.Tests {
		if t.Name == name || strings.Contains(t.Name, name) {
			return &spec.Tests[i]
		}
	}
	return nil
}
