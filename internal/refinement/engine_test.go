package refinement

import (
	"testing"

	"github.com/kon1790/rpg/internal/parity"
)

func TestDefaultLoopConfig(t *testing.T) {
	config := DefaultLoopConfig()

	if config.ConvergenceThreshold < 0.9 || config.ConvergenceThreshold > 1.0 {
		t.Errorf("ConvergenceThreshold should be between 0.9 and 1.0, got %.2f", config.ConvergenceThreshold)
	}

	if config.MaxIterations < 1 {
		t.Errorf("MaxIterations should be at least 1, got %d", config.MaxIterations)
	}

	if config.StuckWindow < 2 {
		t.Errorf("StuckWindow should be at least 2, got %d", config.StuckWindow)
	}
}

func TestConvergenceTracker(t *testing.T) {
	config := LoopConfig{
		ConvergenceThreshold: 0.95,
		StuckThreshold:       0.02,
		StuckWindow:          3,
	}

	tracker := NewConvergenceTracker(config)

	// Add improving scores
	tracker.AddScore(0.70)
	tracker.AddScore(0.80)
	tracker.AddScore(0.85)

	metrics := tracker.GetMetrics()

	if len(metrics.Scores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(metrics.Scores))
	}

	if metrics.Trend <= 0 {
		t.Errorf("Expected positive trend, got %.4f", metrics.Trend)
	}

	if metrics.IsStuck {
		t.Error("Should not be stuck with improving scores")
	}
}

func TestConvergenceTrackerStuck(t *testing.T) {
	config := LoopConfig{
		ConvergenceThreshold: 0.95,
		StuckThreshold:       0.02,
		StuckWindow:          3,
	}

	tracker := NewConvergenceTracker(config)

	// Add flat scores (stuck)
	tracker.AddScore(0.80)
	tracker.AddScore(0.80)
	tracker.AddScore(0.81)

	metrics := tracker.GetMetrics()

	if !metrics.IsStuck {
		t.Error("Should be stuck with flat scores")
	}

	if metrics.StuckReason == "" {
		t.Error("Should have a stuck reason")
	}
}

func TestNewEngine(t *testing.T) {
	config := DefaultLoopConfig()
	engine := NewEngine(config, nil)

	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}

	if engine.comparator == nil {
		t.Error("Engine should have a comparator")
	}

	if engine.cache == nil {
		t.Error("Engine should have a cache")
	}

	if engine.convergence == nil {
		t.Error("Engine should have a convergence tracker")
	}
}

func TestEngineDeterminePhase(t *testing.T) {
	parityResult := &parity.ParityResult{
		Gaps: []parity.ParityGap{
			{Dimension: "type", SourceItem: parity.ItemReference{Name: "User"}},
			{Dimension: "behavioral", GeneratedItem: &parity.ItemReference{Language: "python"}},
		},
	}

	tests := []struct {
		strategy  Strategy
		iteration int
		expected  string
	}{
		{StrategySpecFirst, 0, "spec"},
		{StrategySpecFirst, 1, "spec"},
		{StrategySpecFirst, 2, "code"},
		{StrategyCodeFirst, 0, "code"},
		{StrategyCodeFirst, 2, "spec"},
		{StrategyBalanced, 0, "spec"},
		{StrategyBalanced, 1, "code"},
		{StrategyBalanced, 2, "spec"},
	}

	for _, tc := range tests {
		config := DefaultLoopConfig()
		config.RefinementStrategy = tc.strategy
		engine := NewEngine(config, nil)

		result := engine.determinePhase(tc.iteration, parityResult)
		if result != tc.expected {
			t.Errorf("Strategy %s, iteration %d: expected %s, got %s",
				tc.strategy, tc.iteration, tc.expected, result)
		}
	}
}

func TestEngineGenerateRefinementInstructions(t *testing.T) {
	config := DefaultLoopConfig()
	engine := NewEngine(config, nil)

	// Create a parity result with gaps
	parityResult := &parity.ParityResult{
		OverallScore: 0.75,
		Gaps: []parity.ParityGap{
			{
				Dimension: "type",
				Severity:  "high",
				SourceItem: parity.ItemReference{
					Type:     "type",
					Name:     "User",
					Location: "user.go:10",
				},
				Discrepancy:  "type missing in generated code",
				SuggestedFix: "Add User type definition",
			},
			{
				Dimension: "behavioral",
				Severity:  "high",
				SourceItem: parity.ItemReference{
					Type:      "function",
					Name:      "CreateUser",
					Signature: "func CreateUser(name string) *User",
					Location:  "user.go:25",
				},
				GeneratedItem: &parity.ItemReference{
					Type:     "function",
					Name:     "create_user",
					Language: "python",
				},
				Discrepancy:  "parameter mismatch",
				SuggestedFix: "Add missing parameter",
			},
		},
	}

	// Test spec phase
	instructions := engine.generateRefinementInstructions(parityResult, nil, "spec")

	if len(instructions.SpecRefinements) == 0 {
		t.Error("Expected spec refinements for spec phase")
	}

	if instructions.Priority != "spec" {
		t.Errorf("Expected priority 'spec', got '%s'", instructions.Priority)
	}

	// Test code phase
	instructions = engine.generateRefinementInstructions(parityResult, nil, "code")

	if len(instructions.CodeRefinements) == 0 {
		t.Error("Expected code refinements for code phase")
	}

	if instructions.Priority != "code" {
		t.Errorf("Expected priority 'code', got '%s'", instructions.Priority)
	}

	// Test both phase
	instructions = engine.generateRefinementInstructions(parityResult, nil, "both")

	if len(instructions.SpecRefinements) == 0 {
		t.Error("Expected spec refinements for both phase")
	}

	if len(instructions.CodeRefinements) == 0 {
		t.Error("Expected code refinements for both phase")
	}
}

func TestEngineGenerateSummary(t *testing.T) {
	config := DefaultLoopConfig()
	engine := NewEngine(config, nil)

	// Test converged result
	result := &LoopResult{
		Converged:      true,
		FinalScore:     0.96,
		IterationsUsed: 3,
		IterationHistory: []Iteration{
			{Number: 1, ParityScore: 0.75},
			{Number: 2, ParityScore: 0.85},
			{Number: 3, ParityScore: 0.96},
		},
	}

	summary := engine.generateSummary(result)

	if summary == "" {
		t.Error("Summary should not be empty")
	}

	if !contains(summary, "Converged") {
		t.Error("Summary should mention convergence")
	}

	if !contains(summary, "96") {
		t.Error("Summary should include final score")
	}

	// Test non-converged result
	result.Converged = false
	result.FinalScore = 0.80
	result.UnresolvedGaps = []parity.ParityGap{
		{Severity: "high", SourceItem: parity.ItemReference{Name: "TestFunc"}, Discrepancy: "missing"},
	}

	summary = engine.generateSummary(result)

	if !contains(summary, "Did not converge") {
		t.Error("Summary should mention non-convergence")
	}

	if !contains(summary, "unresolved") {
		t.Error("Summary should mention unresolved gaps")
	}
}

func TestEngineGenerateRefinementPrompt(t *testing.T) {
	config := DefaultLoopConfig()
	engine := NewEngine(config, nil)

	instructions := &RefinementInstructions{
		Summary:  "Found 2 gaps",
		Priority: "both",
		SpecRefinements: []SpecChange{
			{
				Section:     "Types",
				Action:      "add",
				ElementName: "User",
				Description: "Add User type",
			},
		},
		CodeRefinements: map[string][]CodeChange{
			"python": {
				{
					Action:      "modify",
					ElementType: "function",
					ElementName: "create_user",
					Description: "Fix parameter types",
				},
			},
		},
	}

	prompt := engine.GenerateRefinementPrompt(instructions, "go")

	if prompt == "" {
		t.Error("Prompt should not be empty")
	}

	if !contains(prompt, "Spec Changes") {
		t.Error("Prompt should include spec changes section")
	}

	if !contains(prompt, "PYTHON") {
		t.Error("Prompt should include python code changes")
	}

	if !contains(prompt, "User") {
		t.Error("Prompt should mention User type")
	}
}

func TestAnalysisCache(t *testing.T) {
	cache := NewAnalysisCache()

	if cache.Source != nil {
		t.Error("Initial source should be nil")
	}

	if len(cache.Generated) != 0 {
		t.Error("Initial generated map should be empty")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
