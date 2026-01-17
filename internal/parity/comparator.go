package parity

import (
	"fmt"
	"strings"

	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// Comparator performs semantic parity comparison
type Comparator struct {
	config     ComparisonConfig
	normalizer *Normalizer
}

// NewComparator creates a new semantic comparator
func NewComparator(config ComparisonConfig) *Comparator {
	return &Comparator{
		config:     config,
		normalizer: NewNormalizer(),
	}
}

// Compare compares source analysis against multiple generated analyses
func (c *Comparator) Compare(source *semantic.Analysis, generated map[string]*semantic.Analysis) *ParityResult {
	result := &ParityResult{
		ByLanguage: make(map[string]LanguageResult),
		Gaps:       []ParityGap{},
	}

	// Normalize source
	sourceFuncs := c.normalizer.NormalizeFunctions(source.Functions, c.config.IgnorePrivate)
	sourceTypes := c.normalizer.NormalizeTypes(source.Types, c.config.IgnorePrivate)

	// Compare each generated language
	var totalScore float64
	for lang, genAnalysis := range generated {
		langResult := c.compareLanguage(sourceFuncs, sourceTypes, genAnalysis)
		langResult.Language = treesitter.Language(lang)
		result.ByLanguage[lang] = langResult

		// Collect gaps
		for _, gap := range c.identifyGaps(source, genAnalysis, lang) {
			result.Gaps = append(result.Gaps, gap)
		}

		totalScore += langResult.OverallScore
	}

	// Calculate overall score
	if len(generated) > 0 {
		result.OverallScore = totalScore / float64(len(generated))
	}

	// Calculate dimension averages
	result.ByDimension = c.averageDimensions(result.ByLanguage)

	// Check convergence
	result.Converged = result.OverallScore >= c.config.Threshold

	return result
}

// compareLanguage compares source against a single language's generated code
func (c *Comparator) compareLanguage(sourceFuncs []NormalizedSignature, sourceTypes []NormalizedType, gen *semantic.Analysis) LanguageResult {
	result := LanguageResult{}

	// Normalize generated
	genFuncs := c.normalizer.NormalizeFunctions(gen.Functions, c.config.IgnorePrivate)
	genTypes := c.normalizer.NormalizeTypes(gen.Types, c.config.IgnorePrivate)

	// Calculate structural score
	result.ByDimension.Structural = c.calculateStructuralScore(sourceFuncs, genFuncs, sourceTypes, genTypes)

	// Calculate type parity score
	result.ByDimension.Type, result.MissingTypes, result.TypeErrors = c.calculateTypeScore(sourceTypes, genTypes)

	// Calculate behavioral (function) score
	result.ByDimension.Behavioral, result.MissingFuncs, result.SigErrors = c.calculateBehavioralScore(sourceFuncs, genFuncs)

	// Test parity - placeholder (would require test analysis)
	result.ByDimension.Test = 1.0 // Assume full if no test info

	// Idiomatic score - placeholder (would require style analysis)
	result.ByDimension.Idiomatic = 1.0

	// Calculate weighted overall score
	result.OverallScore = c.calculateWeightedScore(result.ByDimension)

	return result
}

// calculateStructuralScore compares overall code structure
func (c *Comparator) calculateStructuralScore(
	sourceFuncs, genFuncs []NormalizedSignature,
	sourceTypes, genTypes []NormalizedType,
) float64 {
	// Compare counts (rough structural similarity)
	funcScore := compareCount(len(sourceFuncs), len(genFuncs))
	typeScore := compareCount(len(sourceTypes), len(genTypes))

	return (funcScore + typeScore) / 2.0
}

// calculateTypeScore compares type definitions
func (c *Comparator) calculateTypeScore(sourceTypes, genTypes []NormalizedType) (float64, []string, []TypeMismatch) {
	var missingTypes []string
	var typeErrors []TypeMismatch

	// Build lookup for generated types
	genTypeMap := make(map[string]NormalizedType)
	for _, t := range genTypes {
		genTypeMap[t.Name] = t
	}

	matchedCount := 0
	for _, srcType := range sourceTypes {
		if genType, ok := genTypeMap[srcType.Name]; ok {
			matchedCount++

			// Check for differences
			match, diffs := TypeMatch(srcType, genType, c.config.StrictTypeMatching)
			if !match {
				typeErrors = append(typeErrors, TypeMismatch{
					TypeName: srcType.Name,
					SourceType: TypeInfo{
						Kind:       srcType.Kind,
						FieldCount: len(srcType.Fields),
						Methods:    srcType.Methods,
						Implements: srcType.Implements,
					},
					GenType: TypeInfo{
						Kind:       genType.Kind,
						FieldCount: len(genType.Fields),
						Methods:    genType.Methods,
						Implements: genType.Implements,
					},
					Differences: diffs,
				})
			}
		} else {
			missingTypes = append(missingTypes, srcType.Name)
		}
	}

	if len(sourceTypes) == 0 {
		return 1.0, missingTypes, typeErrors
	}

	// Score based on matches and errors
	matchScore := float64(matchedCount) / float64(len(sourceTypes))
	errorPenalty := float64(len(typeErrors)) / float64(len(sourceTypes)) * 0.5

	score := matchScore - errorPenalty
	if score < 0 {
		score = 0
	}

	return score, missingTypes, typeErrors
}

// calculateBehavioralScore compares function signatures
func (c *Comparator) calculateBehavioralScore(sourceFuncs, genFuncs []NormalizedSignature) (float64, []string, []SignatureMismatch) {
	var missingFuncs []string
	var sigErrors []SignatureMismatch

	// Build lookup for generated functions
	genFuncMap := make(map[string]NormalizedSignature)
	for _, f := range genFuncs {
		genFuncMap[f.Name] = f
	}

	matchedCount := 0
	for _, srcFunc := range sourceFuncs {
		if genFunc, ok := genFuncMap[srcFunc.Name]; ok {
			matchedCount++

			// Check for signature differences
			match, diffs := SignatureMatch(srcFunc, genFunc, c.config.StrictTypeMatching)
			if !match {
				sigErrors = append(sigErrors, SignatureMismatch{
					FuncName:    srcFunc.Name,
					SourceSig:   formatSignature(srcFunc),
					GenSig:      formatSignature(genFunc),
					Differences: diffs,
				})
			}
		} else {
			missingFuncs = append(missingFuncs, srcFunc.Name)
		}
	}

	if len(sourceFuncs) == 0 {
		return 1.0, missingFuncs, sigErrors
	}

	// Score based on matches and errors
	matchScore := float64(matchedCount) / float64(len(sourceFuncs))
	errorPenalty := float64(len(sigErrors)) / float64(len(sourceFuncs)) * 0.5

	score := matchScore - errorPenalty
	if score < 0 {
		score = 0
	}

	return score, missingFuncs, sigErrors
}

// calculateWeightedScore calculates the weighted overall score
func (c *Comparator) calculateWeightedScore(dims DimensionScores) float64 {
	return dims.Structural*c.config.Weights.Structural +
		dims.Type*c.config.Weights.Type +
		dims.Behavioral*c.config.Weights.Behavioral +
		dims.Test*c.config.Weights.Test +
		dims.Idiomatic*c.config.Weights.Idiomatic
}

// averageDimensions calculates average dimension scores across languages
func (c *Comparator) averageDimensions(byLang map[string]LanguageResult) DimensionScores {
	if len(byLang) == 0 {
		return DimensionScores{}
	}

	var sum DimensionScores
	for _, lr := range byLang {
		sum.Structural += lr.ByDimension.Structural
		sum.Type += lr.ByDimension.Type
		sum.Behavioral += lr.ByDimension.Behavioral
		sum.Test += lr.ByDimension.Test
		sum.Idiomatic += lr.ByDimension.Idiomatic
	}

	n := float64(len(byLang))
	return DimensionScores{
		Structural: sum.Structural / n,
		Type:       sum.Type / n,
		Behavioral: sum.Behavioral / n,
		Test:       sum.Test / n,
		Idiomatic:  sum.Idiomatic / n,
	}
}

// identifyGaps identifies specific parity gaps
func (c *Comparator) identifyGaps(source, gen *semantic.Analysis, lang string) []ParityGap {
	var gaps []ParityGap

	// Normalize for comparison
	sourceFuncs := c.normalizer.NormalizeFunctions(source.Functions, c.config.IgnorePrivate)
	genFuncs := c.normalizer.NormalizeFunctions(gen.Functions, c.config.IgnorePrivate)
	sourceTypes := c.normalizer.NormalizeTypes(source.Types, c.config.IgnorePrivate)
	genTypes := c.normalizer.NormalizeTypes(gen.Types, c.config.IgnorePrivate)

	// Build lookups
	genFuncMap := make(map[string]NormalizedSignature)
	for _, f := range genFuncs {
		genFuncMap[f.Name] = f
	}
	genTypeMap := make(map[string]NormalizedType)
	for _, t := range genTypes {
		genTypeMap[t.Name] = t
	}

	// Find missing or mismatched functions
	for i, srcFunc := range sourceFuncs {
		srcOriginal := source.Functions[i]

		if genFunc, ok := genFuncMap[srcFunc.Name]; ok {
			match, diffs := SignatureMatch(srcFunc, genFunc, c.config.StrictTypeMatching)
			if !match {
				gaps = append(gaps, ParityGap{
					Dimension: "behavioral",
					Severity:  "high",
					SourceItem: ItemReference{
						Type:      "function",
						Name:      srcOriginal.Name,
						Signature: srcOriginal.Signature,
						Location:  fmt.Sprintf("%s:%d", srcOriginal.Location.File, srcOriginal.Location.StartLine),
						Language:  string(source.Language),
					},
					GeneratedItem: &ItemReference{
						Type:      "function",
						Name:      srcFunc.Name,
						Signature: formatSignature(genFunc),
						Language:  lang,
					},
					Discrepancy:  strings.Join(diffs, "; "),
					SuggestedFix: fmt.Sprintf("Update function '%s' in %s to match source signature", srcFunc.Name, lang),
				})
			}
		} else {
			gaps = append(gaps, ParityGap{
				Dimension: "behavioral",
				Severity:  "high",
				SourceItem: ItemReference{
					Type:      "function",
					Name:      srcOriginal.Name,
					Signature: srcOriginal.Signature,
					Location:  fmt.Sprintf("%s:%d", srcOriginal.Location.File, srcOriginal.Location.StartLine),
					Language:  string(source.Language),
				},
				Discrepancy:  "function missing in generated code",
				SuggestedFix: fmt.Sprintf("Implement function '%s' in %s", srcOriginal.Name, lang),
			})
		}
	}

	// Find missing or mismatched types
	for i, srcType := range sourceTypes {
		srcOriginal := source.Types[i]

		if genType, ok := genTypeMap[srcType.Name]; ok {
			match, diffs := TypeMatch(srcType, genType, c.config.StrictTypeMatching)
			if !match {
				gaps = append(gaps, ParityGap{
					Dimension: "type",
					Severity:  "medium",
					SourceItem: ItemReference{
						Type:     "type",
						Name:     srcOriginal.Name,
						Location: fmt.Sprintf("%s:%d", srcOriginal.Location.File, srcOriginal.Location.StartLine),
						Language: string(source.Language),
					},
					GeneratedItem: &ItemReference{
						Type:     "type",
						Name:     srcType.Name,
						Language: lang,
					},
					Discrepancy:  strings.Join(diffs, "; "),
					SuggestedFix: fmt.Sprintf("Update type '%s' in %s to match source definition", srcType.Name, lang),
				})
			}
		} else {
			gaps = append(gaps, ParityGap{
				Dimension: "type",
				Severity:  "high",
				SourceItem: ItemReference{
					Type:     "type",
					Name:     srcOriginal.Name,
					Location: fmt.Sprintf("%s:%d", srcOriginal.Location.File, srcOriginal.Location.StartLine),
					Language: string(source.Language),
				},
				Discrepancy:  "type missing in generated code",
				SuggestedFix: fmt.Sprintf("Implement type '%s' in %s", srcOriginal.Name, lang),
			})
		}
	}

	return gaps
}

// Helper functions

func compareCount(source, gen int) float64 {
	if source == 0 && gen == 0 {
		return 1.0
	}
	if source == 0 {
		return 0.0
	}
	ratio := float64(gen) / float64(source)
	if ratio > 1.0 {
		ratio = 1.0 / ratio
	}
	return ratio
}

func formatSignature(sig NormalizedSignature) string {
	var params []string
	for _, p := range sig.Parameters {
		params = append(params, p.Name+": "+p.BaseType)
	}
	return fmt.Sprintf("%s(%s) -> %s", sig.Name, strings.Join(params, ", "), strings.Join(sig.Returns, ", "))
}

// GenerateFixInstructions generates detailed fix instructions from gaps
func GenerateFixInstructions(result *ParityResult, sourceLang string) string {
	if len(result.Gaps) == 0 {
		return fmt.Sprintf("All implementations have %.1f%% parity. No fixes needed.", result.OverallScore*100)
	}

	var sb strings.Builder
	sb.WriteString("# Parity Fix Instructions\n\n")
	sb.WriteString(fmt.Sprintf("**Overall Parity Score:** %.1f%%\n", result.OverallScore*100))
	sb.WriteString(fmt.Sprintf("**Reference Language:** %s\n\n", sourceLang))

	// Group gaps by language
	gapsByLang := make(map[string][]ParityGap)
	for _, gap := range result.Gaps {
		if gap.GeneratedItem != nil {
			gapsByLang[gap.GeneratedItem.Language] = append(gapsByLang[gap.GeneratedItem.Language], gap)
		} else {
			gapsByLang["unknown"] = append(gapsByLang["unknown"], gap)
		}
	}

	for lang, gaps := range gapsByLang {
		sb.WriteString(fmt.Sprintf("## %s (%d issues)\n\n", strings.ToUpper(lang), len(gaps)))

		// Group by severity
		high := filterBySeverity(gaps, "high")
		medium := filterBySeverity(gaps, "medium")
		low := filterBySeverity(gaps, "low")

		if len(high) > 0 {
			sb.WriteString("### Critical (High Severity)\n\n")
			for _, g := range high {
				writeGap(&sb, g)
			}
		}

		if len(medium) > 0 {
			sb.WriteString("### Important (Medium Severity)\n\n")
			for _, g := range medium {
				writeGap(&sb, g)
			}
		}

		if len(low) > 0 {
			sb.WriteString("### Minor (Low Severity)\n\n")
			for _, g := range low {
				writeGap(&sb, g)
			}
		}
	}

	sb.WriteString("\n---\n")
	sb.WriteString("After applying fixes, run parity analysis again to verify improvements.\n")

	return sb.String()
}

func filterBySeverity(gaps []ParityGap, severity string) []ParityGap {
	var result []ParityGap
	for _, g := range gaps {
		if g.Severity == severity {
			result = append(result, g)
		}
	}
	return result
}

func writeGap(sb *strings.Builder, g ParityGap) {
	sb.WriteString(fmt.Sprintf("- **%s** `%s`\n", g.Dimension, g.SourceItem.Name))
	sb.WriteString(fmt.Sprintf("  - Source: `%s` at `%s`\n", g.SourceItem.Signature, g.SourceItem.Location))
	sb.WriteString(fmt.Sprintf("  - Issue: %s\n", g.Discrepancy))
	sb.WriteString(fmt.Sprintf("  - Fix: %s\n\n", g.SuggestedFix))
}
