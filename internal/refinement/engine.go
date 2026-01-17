// Package refinement provides the iterative refinement engine for achieving
// high parity between source, spec, and generated code.
package refinement

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
	"github.com/kon1790/rpg/internal/parity"
)

// Engine orchestrates the iterative refinement loop
type Engine struct {
	config       LoopConfig
	comparator   *parity.Comparator
	registry     *semantic.AnalyzerRegistry
	cache        *AnalysisCache
	convergence  *ConvergenceTracker
}

// NewEngine creates a new refinement engine
func NewEngine(config LoopConfig, registry *semantic.AnalyzerRegistry) *Engine {
	return &Engine{
		config:      config,
		comparator:  parity.NewComparator(config.ParityConfig),
		registry:    registry,
		cache:       NewAnalysisCache(),
		convergence: NewConvergenceTracker(config),
	}
}

// AnalyzeSource performs deep semantic analysis on source code
func (e *Engine) AnalyzeSource(ctx context.Context, sourcePath string, lang treesitter.Language) (*semantic.Analysis, error) {
	// Check cache first
	if e.cache.Source != nil {
		return e.cache.Source, nil
	}

	// Get analyzer for language
	analyzer, ok := e.registry.Get(lang)
	if !ok {
		return nil, fmt.Errorf("no analyzer for language %s", lang)
	}

	// Perform analysis
	analysis, err := analyzer.Analyze(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("analyzing source: %w", err)
	}

	// Cache result
	e.cache.Source = analysis
	return analysis, nil
}

// AnalyzeGenerated performs semantic analysis on generated code
func (e *Engine) AnalyzeGenerated(ctx context.Context, projectPath string, lang treesitter.Language) (*semantic.Analysis, error) {
	// Check cache
	if cached, ok := e.cache.Generated[string(lang)]; ok {
		return cached, nil
	}

	// Get analyzer
	analyzer, ok := e.registry.Get(lang)
	if !ok {
		return nil, fmt.Errorf("no analyzer for language %s", lang)
	}

	// Analyze
	analysis, err := analyzer.Analyze(projectPath)
	if err != nil {
		return nil, fmt.Errorf("analyzing generated %s: %w", lang, err)
	}

	// Cache
	e.cache.Generated[string(lang)] = analysis
	return analysis, nil
}

// CompareParity compares source against generated implementations
func (e *Engine) CompareParity(source *semantic.Analysis, generated map[string]*semantic.Analysis) *parity.ParityResult {
	return e.comparator.Compare(source, generated)
}

// Run executes the full refinement loop
func (e *Engine) Run(ctx context.Context, input *LoopInput) (*LoopResult, error) {
	result := &LoopResult{
		IterationHistory:  []Iteration{},
		GeneratedProjects: make(map[string]string),
	}

	// Detect source language if not provided
	sourceLang := treesitter.Language(input.SourceLanguage)
	if sourceLang == "" {
		sourceLang = treesitter.DetectLanguage(input.SourcePath)
		if sourceLang == "" {
			return nil, fmt.Errorf("could not detect source language for: %s", input.SourcePath)
		}
	}

	// Initial source analysis
	sourceAnalysis, err := e.AnalyzeSource(ctx, input.SourcePath, sourceLang)
	if err != nil {
		return nil, fmt.Errorf("initial source analysis: %w", err)
	}

	// Main refinement loop
	for i := 0; i < e.config.MaxIterations; i++ {
		iteration := Iteration{
			Number:      i + 1,
			Refinements: []Refinement{},
		}

		// Analyze generated code for each target language
		generatedAnalyses := make(map[string]*semantic.Analysis)
		for _, targetLang := range input.TargetLanguages {
			projectPath := e.getProjectPath(input.OutputDir, targetLang)
			result.GeneratedProjects[targetLang] = projectPath

			analysis, err := e.AnalyzeGenerated(ctx, projectPath, treesitter.Language(targetLang))
			if err != nil {
				// Project may not exist yet on first iteration
				continue
			}
			generatedAnalyses[targetLang] = analysis
		}

		// Calculate parity
		parityResult := e.CompareParity(sourceAnalysis, generatedAnalyses)
		iteration.ParityScore = parityResult.OverallScore
		iteration.ParityResult = parityResult

		// Calculate improvement from previous iteration
		if i > 0 {
			prevScore := result.IterationHistory[i-1].ParityScore
			iteration.ScoreImprovement = iteration.ParityScore - prevScore
		}

		// Determine refinement phase based on strategy
		iteration.Phase = e.determinePhase(i, parityResult)

		// Generate refinement instructions
		instructions := e.generateRefinementInstructions(parityResult, sourceAnalysis, iteration.Phase)
		iteration.RefinementsApplied = len(instructions.SpecRefinements) + e.countCodeRefinements(instructions)

		// Track convergence
		e.convergence.AddScore(iteration.ParityScore)
		result.IterationHistory = append(result.IterationHistory, iteration)

		// Check if converged
		if parityResult.OverallScore >= e.config.ConvergenceThreshold {
			result.Converged = true
			result.FinalScore = parityResult.OverallScore
			result.IterationsUsed = i + 1
			result.FinalSpec = input.SpecPath
			result.RefinementSummary = e.generateSummary(result)
			return result, nil
		}

		// Check if stuck
		metrics := e.convergence.GetMetrics()
		if metrics.IsStuck {
			result.UnresolvedGaps = parityResult.Gaps
			result.RefinementSummary = fmt.Sprintf("Refinement stalled after %d iterations: %s. Best score: %.2f%%",
				i+1, metrics.StuckReason, parityResult.OverallScore*100)
			break
		}

		// Clear generated cache for next iteration (source stays cached)
		e.cache.Generated = make(map[string]*semantic.Analysis)
	}

	// Did not converge within max iterations
	if len(result.IterationHistory) > 0 {
		lastIteration := result.IterationHistory[len(result.IterationHistory)-1]
		result.FinalScore = lastIteration.ParityScore
		result.UnresolvedGaps = lastIteration.ParityResult.Gaps
	}
	result.IterationsUsed = len(result.IterationHistory)
	result.FinalSpec = input.SpecPath
	if result.RefinementSummary == "" {
		result.RefinementSummary = e.generateSummary(result)
	}

	return result, nil
}

// determinePhase determines what to refine based on strategy and iteration
func (e *Engine) determinePhase(iteration int, parityResult *parity.ParityResult) string {
	switch e.config.RefinementStrategy {
	case StrategySpecFirst:
		if iteration < 2 {
			return "spec"
		}
		return "code"

	case StrategyCodeFirst:
		if iteration < 2 {
			return "code"
		}
		return "spec"

	case StrategyBalanced:
		if iteration%2 == 0 {
			return "spec"
		}
		return "code"

	case StrategyAdaptive:
		// Analyze gaps to determine what needs more work
		specGaps := 0
		codeGaps := 0
		for _, gap := range parityResult.Gaps {
			if gap.GeneratedItem == nil {
				specGaps++
			} else {
				codeGaps++
			}
		}
		if specGaps > codeGaps {
			return "spec"
		}
		return "code"

	default:
		return "both"
	}
}

// generateRefinementInstructions creates instructions for AI to apply refinements
func (e *Engine) generateRefinementInstructions(parityResult *parity.ParityResult, source *semantic.Analysis, phase string) *RefinementInstructions {
	instructions := &RefinementInstructions{
		Priority:        phase,
		SpecRefinements: []SpecChange{},
		CodeRefinements: make(map[string][]CodeChange),
	}

	if len(parityResult.Gaps) == 0 {
		instructions.Summary = "No refinements needed - parity achieved."
		return instructions
	}

	// Build summary
	var summaryParts []string
	summaryParts = append(summaryParts, fmt.Sprintf("Found %d parity gaps", len(parityResult.Gaps)))

	// Process gaps
	for _, gap := range parityResult.Gaps {
		switch gap.Dimension {
		case "type":
			if phase == "spec" || phase == "both" {
				instructions.SpecRefinements = append(instructions.SpecRefinements, SpecChange{
					Section:     "Types",
					Action:      e.determineAction(gap),
					ElementName: gap.SourceItem.Name,
					Description: gap.Discrepancy,
					Example:     gap.SuggestedFix,
				})
			}
			if (phase == "code" || phase == "both") && gap.GeneratedItem != nil {
				lang := gap.GeneratedItem.Language
				instructions.CodeRefinements[lang] = append(instructions.CodeRefinements[lang], CodeChange{
					Action:      e.determineAction(gap),
					ElementType: "type",
					ElementName: gap.SourceItem.Name,
					Description: gap.Discrepancy,
					SourceRef:   gap.SourceItem.Location,
				})
			}

		case "behavioral":
			if phase == "spec" || phase == "both" {
				instructions.SpecRefinements = append(instructions.SpecRefinements, SpecChange{
					Section:     "Functions",
					Action:      e.determineAction(gap),
					ElementName: gap.SourceItem.Name,
					Description: gap.Discrepancy,
					Example:     gap.SuggestedFix,
				})
			}
			if (phase == "code" || phase == "both") && gap.GeneratedItem != nil {
				lang := gap.GeneratedItem.Language
				instructions.CodeRefinements[lang] = append(instructions.CodeRefinements[lang], CodeChange{
					Action:      e.determineAction(gap),
					ElementType: "function",
					ElementName: gap.SourceItem.Name,
					Description: gap.Discrepancy,
					SourceRef:   gap.SourceItem.Location,
				})
			}

		case "structural":
			if phase == "spec" || phase == "both" {
				instructions.SpecRefinements = append(instructions.SpecRefinements, SpecChange{
					Section:     "Architecture",
					Action:      "modify",
					ElementName: gap.SourceItem.Name,
					Description: gap.Discrepancy,
				})
			}
		}
	}

	// Count refinements by target
	specCount := len(instructions.SpecRefinements)
	codeCount := 0
	for _, changes := range instructions.CodeRefinements {
		codeCount += len(changes)
	}
	summaryParts = append(summaryParts, fmt.Sprintf("%d spec changes, %d code changes needed", specCount, codeCount))

	instructions.Summary = strings.Join(summaryParts, ". ")
	return instructions
}

// determineAction determines what action to take based on gap
func (e *Engine) determineAction(gap parity.ParityGap) string {
	if strings.Contains(gap.Discrepancy, "missing") {
		return "add"
	}
	if strings.Contains(gap.Discrepancy, "mismatch") {
		return "modify"
	}
	return "modify"
}

// countCodeRefinements counts total code refinements across all languages
func (e *Engine) countCodeRefinements(instructions *RefinementInstructions) int {
	count := 0
	for _, changes := range instructions.CodeRefinements {
		count += len(changes)
	}
	return count
}

// getProjectPath returns the path for a generated project
func (e *Engine) getProjectPath(outputDir, lang string) string {
	return fmt.Sprintf("%s/%s", outputDir, lang)
}

// generateSummary generates a human-readable summary of the refinement process
func (e *Engine) generateSummary(result *LoopResult) string {
	var sb strings.Builder

	if result.Converged {
		sb.WriteString(fmt.Sprintf("✓ Converged after %d iteration(s) with %.1f%% parity.\n",
			result.IterationsUsed, result.FinalScore*100))
	} else {
		sb.WriteString(fmt.Sprintf("✗ Did not converge after %d iteration(s). Final score: %.1f%%\n",
			result.IterationsUsed, result.FinalScore*100))
	}

	// Score progression
	if len(result.IterationHistory) > 1 {
		sb.WriteString("\nScore progression: ")
		for i, iter := range result.IterationHistory {
			if i > 0 {
				sb.WriteString(" → ")
			}
			sb.WriteString(fmt.Sprintf("%.1f%%", iter.ParityScore*100))
		}
		sb.WriteString("\n")
	}

	// Unresolved gaps
	if len(result.UnresolvedGaps) > 0 {
		sb.WriteString(fmt.Sprintf("\n%d unresolved gap(s):\n", len(result.UnresolvedGaps)))
		for _, gap := range result.UnresolvedGaps {
			sb.WriteString(fmt.Sprintf("  - [%s] %s: %s\n", gap.Severity, gap.SourceItem.Name, gap.Discrepancy))
		}
	}

	return sb.String()
}

// GenerateRefinementPrompt generates a prompt for AI to apply refinements
func (e *Engine) GenerateRefinementPrompt(instructions *RefinementInstructions, sourceLang string) string {
	var sb strings.Builder

	sb.WriteString("# Refinement Instructions\n\n")
	sb.WriteString(instructions.Summary)
	sb.WriteString("\n\n")

	if instructions.Priority == "spec" || instructions.Priority == "both" {
		if len(instructions.SpecRefinements) > 0 {
			sb.WriteString("## Spec Changes\n\n")
			for _, change := range instructions.SpecRefinements {
				sb.WriteString(fmt.Sprintf("### %s: %s `%s`\n", change.Action, change.Section, change.ElementName))
				sb.WriteString(fmt.Sprintf("%s\n", change.Description))
				if change.Example != "" {
					sb.WriteString(fmt.Sprintf("\nSuggested fix: %s\n", change.Example))
				}
				sb.WriteString("\n")
			}
		}
	}

	if instructions.Priority == "code" || instructions.Priority == "both" {
		for lang, changes := range instructions.CodeRefinements {
			if len(changes) > 0 {
				sb.WriteString(fmt.Sprintf("## %s Code Changes\n\n", strings.ToUpper(lang)))
				for _, change := range changes {
					sb.WriteString(fmt.Sprintf("### %s %s `%s`\n", change.Action, change.ElementType, change.ElementName))
					sb.WriteString(fmt.Sprintf("%s\n", change.Description))
					if change.SourceRef != "" {
						sb.WriteString(fmt.Sprintf("Reference: %s\n", change.SourceRef))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String()
}

// ConvergenceTracker tracks convergence progress
type ConvergenceTracker struct {
	config LoopConfig
	scores []float64
}

// NewConvergenceTracker creates a new convergence tracker
func NewConvergenceTracker(config LoopConfig) *ConvergenceTracker {
	return &ConvergenceTracker{
		config: config,
		scores: []float64{},
	}
}

// AddScore adds a score to the tracker
func (ct *ConvergenceTracker) AddScore(score float64) {
	ct.scores = append(ct.scores, score)
}

// GetMetrics returns current convergence metrics
func (ct *ConvergenceTracker) GetMetrics() ConvergenceMetrics {
	metrics := ConvergenceMetrics{
		Scores: ct.scores,
	}

	if len(ct.scores) == 0 {
		return metrics
	}

	// Calculate trend (slope of recent scores)
	if len(ct.scores) >= 2 {
		recent := ct.scores[len(ct.scores)-min(3, len(ct.scores)):]
		metrics.Trend = recent[len(recent)-1] - recent[0]
	}

	// Check if stuck
	if len(ct.scores) >= ct.config.StuckWindow {
		windowScores := ct.scores[len(ct.scores)-ct.config.StuckWindow:]
		improvement := windowScores[len(windowScores)-1] - windowScores[0]

		if math.Abs(improvement) < ct.config.StuckThreshold {
			metrics.IsStuck = true
			if improvement <= 0 {
				metrics.StuckReason = "no improvement in recent iterations"
			} else {
				metrics.StuckReason = "improvement below threshold"
			}
		}
	}

	// Estimate iterations to converge
	if metrics.Trend > 0 && len(ct.scores) > 0 {
		currentScore := ct.scores[len(ct.scores)-1]
		remaining := ct.config.ConvergenceThreshold - currentScore
		if metrics.Trend > 0 {
			metrics.EstimatedIterationsToConverge = int(math.Ceil(remaining / metrics.Trend))
		}
	}

	return metrics
}

// Reset clears the tracker
func (ct *ConvergenceTracker) Reset() {
	ct.scores = []float64{}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
