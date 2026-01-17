package parity

import (
	"testing"

	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
)

func TestNormalizeName(t *testing.T) {
	n := NewNormalizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"CreateUser", "create_user"},
		{"getUserByID", "get_user_by_i_d"},
		{"simple", "simple"},
		{"HTTPHandler", "h_t_t_p_handler"},
	}

	for _, tc := range tests {
		result := n.normalizeName(tc.input)
		if result != tc.expected {
			t.Errorf("normalizeName(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestNormalizeType(t *testing.T) {
	n := NewNormalizer()

	tests := []struct {
		input    string
		expected string
	}{
		{"int", "integer"},
		{"int64", "integer"},
		{"string", "string"},
		{"String", "string"},
		{"bool", "boolean"},
		{"float64", "float"},
		{"error", "error"},
		{"context.Context", "context"},
		{"CustomType", "custom_type"},
	}

	for _, tc := range tests {
		result := n.normalizeType(tc.input)
		if result != tc.expected {
			t.Errorf("normalizeType(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestParseType(t *testing.T) {
	n := NewNormalizer()

	tests := []struct {
		input     string
		baseType  string
		isPtr     bool
		isArray   bool
		isMap     bool
	}{
		{"int", "int", false, false, false},
		{"*string", "string", true, false, false},
		{"[]int", "int", false, true, false},
		{"map[string]int", "map", false, false, true},
		{"List<User>", "User", false, true, false},
		{"Optional<int>", "int", true, false, false},
	}

	for _, tc := range tests {
		baseType, isPtr, isArray, isMap := n.parseType(tc.input)
		if baseType != tc.baseType {
			t.Errorf("parseType(%s) baseType = %s, expected %s", tc.input, baseType, tc.baseType)
		}
		if isPtr != tc.isPtr {
			t.Errorf("parseType(%s) isPtr = %v, expected %v", tc.input, isPtr, tc.isPtr)
		}
		if isArray != tc.isArray {
			t.Errorf("parseType(%s) isArray = %v, expected %v", tc.input, isArray, tc.isArray)
		}
		if isMap != tc.isMap {
			t.Errorf("parseType(%s) isMap = %v, expected %v", tc.input, isMap, tc.isMap)
		}
	}
}

func TestSignatureMatch(t *testing.T) {
	sigA := NormalizedSignature{
		Name: "create_user",
		Parameters: []NormalizedParam{
			{Name: "ctx", BaseType: "context"},
			{Name: "name", BaseType: "string"},
		},
		Returns: []string{"user", "error"},
	}

	// Exact match
	sigB := NormalizedSignature{
		Name: "create_user",
		Parameters: []NormalizedParam{
			{Name: "ctx", BaseType: "context"},
			{Name: "name", BaseType: "string"},
		},
		Returns: []string{"user", "error"},
	}

	match, diffs := SignatureMatch(sigA, sigB, false)
	if !match {
		t.Errorf("Expected signatures to match, got diffs: %v", diffs)
	}

	// Different parameter type
	sigC := NormalizedSignature{
		Name: "create_user",
		Parameters: []NormalizedParam{
			{Name: "ctx", BaseType: "context"},
			{Name: "name", BaseType: "integer"}, // Different type
		},
		Returns: []string{"user", "error"},
	}

	match, diffs = SignatureMatch(sigA, sigC, false)
	if match {
		t.Error("Expected signatures to not match due to parameter type difference")
	}
	if len(diffs) == 0 {
		t.Error("Expected differences to be reported")
	}

	// Different return type
	sigD := NormalizedSignature{
		Name: "create_user",
		Parameters: []NormalizedParam{
			{Name: "ctx", BaseType: "context"},
			{Name: "name", BaseType: "string"},
		},
		Returns: []string{"user"}, // Missing error
	}

	match, diffs = SignatureMatch(sigA, sigD, false)
	if match {
		t.Error("Expected signatures to not match due to return type difference")
	}
}

func TestTypeMatch(t *testing.T) {
	typeA := NormalizedType{
		Name: "user",
		Kind: "struct",
		Fields: []NormalizedField{
			{Name: "id", BaseType: "integer"},
			{Name: "name", BaseType: "string"},
		},
	}

	// Exact match
	typeB := NormalizedType{
		Name: "user",
		Kind: "struct",
		Fields: []NormalizedField{
			{Name: "id", BaseType: "integer"},
			{Name: "name", BaseType: "string"},
		},
	}

	match, diffs := TypeMatch(typeA, typeB, false)
	if !match {
		t.Errorf("Expected types to match, got diffs: %v", diffs)
	}

	// Missing field
	typeC := NormalizedType{
		Name: "user",
		Kind: "struct",
		Fields: []NormalizedField{
			{Name: "id", BaseType: "integer"},
		},
	}

	match, diffs = TypeMatch(typeA, typeC, false)
	if match {
		t.Error("Expected types to not match due to field count difference")
	}
}

func TestComparator(t *testing.T) {
	config := DefaultConfig()
	comparator := NewComparator(config)

	// Create source analysis
	source := &semantic.Analysis{
		Language: treesitter.LanguageGo,
		Name:     "test-project",
		Types: []semantic.ResolvedType{
			{
				TypeDef: treesitter.TypeDef{
					Name:     "User",
					Kind:     treesitter.TypeKindStruct,
					IsPublic: true,
				},
				ResolvedFields: []semantic.ResolvedField{
					{Field: treesitter.Field{Name: "ID", Type: "int"}, ResolvedType: "int"},
					{Field: treesitter.Field{Name: "Name", Type: "string"}, ResolvedType: "string"},
				},
			},
		},
		Functions: []semantic.ResolvedFunction{
			{
				FunctionDef: treesitter.FunctionDef{
					Name:      "CreateUser",
					IsPublic:  true,
					Signature: "func CreateUser(name string) *User",
				},
				ResolvedParameters:  []semantic.ResolvedParameter{{Parameter: treesitter.Parameter{Name: "name", Type: "string"}, ResolvedType: "string"}},
				ResolvedReturnTypes: []string{"*User"},
			},
		},
	}

	// Create matching generated analysis
	generated := map[string]*semantic.Analysis{
		"python": {
			Language: treesitter.LanguagePython,
			Name:     "test-project",
			Types: []semantic.ResolvedType{
				{
					TypeDef: treesitter.TypeDef{
						Name:     "User",
						Kind:     treesitter.TypeKindClass,
						IsPublic: true,
					},
					ResolvedFields: []semantic.ResolvedField{
						{Field: treesitter.Field{Name: "id", Type: "int"}, ResolvedType: "int"},
						{Field: treesitter.Field{Name: "name", Type: "str"}, ResolvedType: "str"},
					},
				},
			},
			Functions: []semantic.ResolvedFunction{
				{
					FunctionDef: treesitter.FunctionDef{
						Name:      "create_user",
						IsPublic:  true,
						Signature: "def create_user(name: str) -> User",
					},
					ResolvedParameters:  []semantic.ResolvedParameter{{Parameter: treesitter.Parameter{Name: "name", Type: "str"}, ResolvedType: "str"}},
					ResolvedReturnTypes: []string{"User"},
				},
			},
		},
	}

	result := comparator.Compare(source, generated)

	// Should have reasonable parity
	if result.OverallScore < 0.5 {
		t.Errorf("Expected reasonable parity score, got %.2f", result.OverallScore)
	}

	// Should have python result
	if _, ok := result.ByLanguage["python"]; !ok {
		t.Error("Expected python result in ByLanguage")
	}
}

func TestConvertCase(t *testing.T) {
	tests := []struct {
		name       string
		targetCase string
		expected   string
	}{
		{"CreateUser", "snake_case", "create_user"},
		{"create_user", "PascalCase", "CreateUser"},
		{"create_user", "camelCase", "createUser"},
		{"CREATE_USER", "snake_case", "create_user"},
		{"CreateUser", "SCREAMING_CASE", "CREATE_USER"},
	}

	for _, tc := range tests {
		result := ConvertCase(tc.name, tc.targetCase)
		if result != tc.expected {
			t.Errorf("ConvertCase(%s, %s) = %s, expected %s", tc.name, tc.targetCase, result, tc.expected)
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	// Weights should sum to 1.0
	sum := config.Weights.Structural + config.Weights.Type + config.Weights.Behavioral + config.Weights.Test + config.Weights.Idiomatic
	if sum < 0.99 || sum > 1.01 {
		t.Errorf("Weights should sum to 1.0, got %.2f", sum)
	}

	// Threshold should be reasonable
	if config.Threshold < 0.8 || config.Threshold > 1.0 {
		t.Errorf("Threshold should be between 0.8 and 1.0, got %.2f", config.Threshold)
	}
}
