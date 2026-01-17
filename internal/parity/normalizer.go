package parity

import (
	"regexp"
	"strings"

	"github.com/kon1790/rpg/internal/importer/semantic"
	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// Normalizer provides cross-language normalization for comparison
type Normalizer struct {
	// Type vocabulary mapping: language-specific types -> normalized types
	typeVocab map[string]string
}

// NewNormalizer creates a new normalizer
func NewNormalizer() *Normalizer {
	return &Normalizer{
		typeVocab: buildTypeVocabulary(),
	}
}

// buildTypeVocabulary builds the common type vocabulary
func buildTypeVocabulary() map[string]string {
	return map[string]string{
		// Integers
		"int":     "integer",
		"int8":    "integer",
		"int16":   "integer",
		"int32":   "integer",
		"int64":   "integer",
		"uint":    "integer",
		"uint8":   "integer",
		"uint16":  "integer",
		"uint32":  "integer",
		"uint64":  "integer",
		"i8":      "integer",
		"i16":     "integer",
		"i32":     "integer",
		"i64":     "integer",
		"u8":      "integer",
		"u16":     "integer",
		"u32":     "integer",
		"u64":     "integer",
		"isize":   "integer",
		"usize":   "integer",
		"Integer": "integer",
		"Long":    "integer",
		"Short":   "integer",
		"number":  "number",

		// Floats
		"float":   "float",
		"float32": "float",
		"float64": "float",
		"f32":     "float",
		"f64":     "float",
		"Float":   "float",
		"Double":  "float",
		"double":  "float",

		// Strings
		"string":       "string",
		"String":       "string",
		"&str":         "string",
		"str":          "string",
		"StringBuilder": "string",

		// Booleans
		"bool":    "boolean",
		"boolean": "boolean",
		"Boolean": "boolean",

		// Void/Unit
		"void":   "void",
		"()":     "void",
		"None":   "void",
		"null":   "null",
		"nil":    "null",

		// Byte
		"byte":  "byte",
		"Byte":  "byte",
		"[]byte": "bytes",
		"bytes": "bytes",

		// Any/Object
		"any":         "any",
		"interface{}": "any",
		"Object":      "any",
		"object":      "any",
		"dynamic":     "any",

		// Error
		"error":     "error",
		"Error":     "error",
		"Exception": "error",
		"Result":    "result",

		// Context
		"context.Context": "context",
		"Context":         "context",
		"CancellationToken": "context",
	}
}

// NormalizeFunctions normalizes functions for cross-language comparison
func (n *Normalizer) NormalizeFunctions(funcs []semantic.ResolvedFunction, ignorePrivate bool) []NormalizedSignature {
	var normalized []NormalizedSignature

	for _, fn := range funcs {
		if ignorePrivate && !fn.IsPublic {
			continue
		}

		sig := NormalizedSignature{
			Name:       n.normalizeName(fn.Name),
			IsAsync:    fn.IsAsync,
			IsPublic:   fn.IsPublic,
			Complexity: fn.Complexity,
		}

		// Normalize parameters
		for _, p := range fn.ResolvedParameters {
			sig.Parameters = append(sig.Parameters, n.normalizeParam(p))
		}

		// Normalize return types
		for _, rt := range fn.ResolvedReturnTypes {
			sig.Returns = append(sig.Returns, n.normalizeType(rt))
		}

		normalized = append(normalized, sig)
	}

	return normalized
}

// NormalizeTypes normalizes types for cross-language comparison
func (n *Normalizer) NormalizeTypes(types []semantic.ResolvedType, ignorePrivate bool) []NormalizedType {
	var normalized []NormalizedType

	for _, t := range types {
		if ignorePrivate && !t.IsPublic {
			continue
		}

		nt := NormalizedType{
			Name:       n.normalizeName(t.Name),
			Kind:       string(t.Kind),
			Methods:    t.Methods,
			Implements: t.ImplementsInterfaces,
			IsPublic:   t.IsPublic,
		}

		// Normalize fields
		for _, f := range t.ResolvedFields {
			nt.Fields = append(nt.Fields, n.normalizeField(f))
		}

		normalized = append(normalized, nt)
	}

	return normalized
}

// normalizeParam normalizes a parameter
func (n *Normalizer) normalizeParam(p semantic.ResolvedParameter) NormalizedParam {
	baseType, isPtr, isArray, isMap := n.parseType(p.ResolvedType)

	return NormalizedParam{
		Name:     n.normalizeName(p.Name),
		BaseType: n.normalizeType(baseType),
		IsPtr:    isPtr,
		IsArray:  isArray,
		IsMap:    isMap,
	}
}

// normalizeField normalizes a field
func (n *Normalizer) normalizeField(f semantic.ResolvedField) NormalizedField {
	return NormalizedField{
		Name:     n.normalizeName(f.Name),
		BaseType: n.normalizeType(f.ResolvedType),
		IsPtr:    f.IsPointer,
		IsArray:  f.IsSlice,
		IsMap:    f.IsMap,
	}
}

// normalizeName normalizes an identifier name to snake_case
func (n *Normalizer) normalizeName(name string) string {
	// Convert camelCase and PascalCase to snake_case
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// normalizeType normalizes a type name to common vocabulary
func (n *Normalizer) normalizeType(typeName string) string {
	// Strip modifiers and get base type
	base, _, _, _ := n.parseType(typeName)

	// Look up in vocabulary
	if normalized, ok := n.typeVocab[base]; ok {
		return normalized
	}

	// For custom types, just normalize the name
	return n.normalizeName(base)
}

// parseType parses a type string to extract modifiers
func (n *Normalizer) parseType(typeName string) (baseType string, isPtr bool, isArray bool, isMap bool) {
	typeName = strings.TrimSpace(typeName)

	// Check for pointer/reference
	if strings.HasPrefix(typeName, "*") || strings.HasPrefix(typeName, "&") {
		isPtr = true
		typeName = typeName[1:]
	}

	// Check for array/slice
	if strings.HasPrefix(typeName, "[]") {
		isArray = true
		typeName = typeName[2:]
	}
	if strings.HasPrefix(typeName, "List<") || strings.HasPrefix(typeName, "Vec<") ||
		strings.HasSuffix(typeName, "[]") || strings.HasPrefix(typeName, "Array<") {
		isArray = true
		// Extract inner type
		typeName = extractGenericType(typeName)
	}

	// Check for map/dictionary
	if strings.HasPrefix(typeName, "map[") || strings.HasPrefix(typeName, "Map<") ||
		strings.HasPrefix(typeName, "HashMap<") || strings.HasPrefix(typeName, "Dictionary<") ||
		strings.HasPrefix(typeName, "dict[") {
		isMap = true
		// For maps, we just note it's a map
		baseType = "map"
		return
	}

	// Check for Optional/Nullable
	if strings.HasPrefix(typeName, "Optional<") || strings.HasPrefix(typeName, "Option<") ||
		strings.HasSuffix(typeName, "?") {
		isPtr = true // Treat optional as nullable
		typeName = extractGenericType(typeName)
		if strings.HasSuffix(typeName, "?") {
			typeName = typeName[:len(typeName)-1]
		}
	}

	baseType = typeName
	return
}

// extractGenericType extracts the type from a generic like List<T>
func extractGenericType(typeName string) string {
	// Simple extraction: find content between < and >
	start := strings.Index(typeName, "<")
	end := strings.LastIndex(typeName, ">")
	if start != -1 && end > start {
		inner := typeName[start+1 : end]
		// If there's a comma (like Dict<K,V>), take first part or simplify
		if comma := strings.Index(inner, ","); comma != -1 {
			return strings.TrimSpace(inner[:comma])
		}
		return strings.TrimSpace(inner)
	}

	// Handle array suffix like int[]
	if strings.HasSuffix(typeName, "[]") {
		return typeName[:len(typeName)-2]
	}

	return typeName
}

// SignatureMatch checks if two normalized signatures match
func SignatureMatch(a, b NormalizedSignature, strict bool) (bool, []string) {
	var diffs []string

	// Name must match
	if a.Name != b.Name {
		diffs = append(diffs, "name mismatch")
		return false, diffs
	}

	// Parameter count must match
	if len(a.Parameters) != len(b.Parameters) {
		diffs = append(diffs, "parameter count mismatch")
	}

	// Compare parameters
	minParams := len(a.Parameters)
	if len(b.Parameters) < minParams {
		minParams = len(b.Parameters)
	}
	for i := 0; i < minParams; i++ {
		if a.Parameters[i].BaseType != b.Parameters[i].BaseType {
			diffs = append(diffs, "parameter type mismatch: "+a.Parameters[i].Name)
		}
		if strict && a.Parameters[i].IsPtr != b.Parameters[i].IsPtr {
			diffs = append(diffs, "parameter pointer mismatch: "+a.Parameters[i].Name)
		}
	}

	// Return type must match
	if len(a.Returns) != len(b.Returns) {
		diffs = append(diffs, "return type count mismatch")
	} else {
		for i := range a.Returns {
			if a.Returns[i] != b.Returns[i] {
				diffs = append(diffs, "return type mismatch")
			}
		}
	}

	// Async should match (optional check)
	if strict && a.IsAsync != b.IsAsync {
		diffs = append(diffs, "async mismatch")
	}

	return len(diffs) == 0, diffs
}

// TypeMatch checks if two normalized types match
func TypeMatch(a, b NormalizedType, strict bool) (bool, []string) {
	var diffs []string

	// Name must match
	if a.Name != b.Name {
		diffs = append(diffs, "name mismatch")
		return false, diffs
	}

	// Kind should match
	if a.Kind != b.Kind {
		diffs = append(diffs, "kind mismatch: "+a.Kind+" vs "+b.Kind)
	}

	// Field count should be similar
	if len(a.Fields) != len(b.Fields) {
		diffs = append(diffs, "field count mismatch")
	}

	// Compare fields by name
	aFields := make(map[string]NormalizedField)
	for _, f := range a.Fields {
		aFields[f.Name] = f
	}

	for _, bf := range b.Fields {
		if af, ok := aFields[bf.Name]; ok {
			if af.BaseType != bf.BaseType {
				diffs = append(diffs, "field type mismatch: "+bf.Name)
			}
		} else {
			diffs = append(diffs, "missing field: "+bf.Name)
		}
	}

	// Method count
	if len(a.Methods) != len(b.Methods) && strict {
		diffs = append(diffs, "method count mismatch")
	}

	return len(diffs) == 0, diffs
}

// LanguageNaming provides language-specific naming patterns
type LanguageNaming struct {
	// Case convention: camelCase, PascalCase, snake_case, SCREAMING_CASE
	FunctionCase string
	TypeCase     string
	FieldCase    string
	ConstCase    string
}

// GetLanguageNaming returns naming conventions for a language
func GetLanguageNaming(lang treesitter.Language) LanguageNaming {
	switch lang {
	case treesitter.LanguageGo:
		return LanguageNaming{
			FunctionCase: "PascalCase",
			TypeCase:     "PascalCase",
			FieldCase:    "PascalCase",
			ConstCase:    "PascalCase",
		}
	case treesitter.LanguageRust:
		return LanguageNaming{
			FunctionCase: "snake_case",
			TypeCase:     "PascalCase",
			FieldCase:    "snake_case",
			ConstCase:    "SCREAMING_CASE",
		}
	case treesitter.LanguagePython:
		return LanguageNaming{
			FunctionCase: "snake_case",
			TypeCase:     "PascalCase",
			FieldCase:    "snake_case",
			ConstCase:    "SCREAMING_CASE",
		}
	case treesitter.LanguageJava, treesitter.LanguageCSharp:
		return LanguageNaming{
			FunctionCase: "camelCase",
			TypeCase:     "PascalCase",
			FieldCase:    "camelCase",
			ConstCase:    "SCREAMING_CASE",
		}
	case treesitter.LanguageTypeScript:
		return LanguageNaming{
			FunctionCase: "camelCase",
			TypeCase:     "PascalCase",
			FieldCase:    "camelCase",
			ConstCase:    "SCREAMING_CASE",
		}
	default:
		return LanguageNaming{
			FunctionCase: "camelCase",
			TypeCase:     "PascalCase",
			FieldCase:    "camelCase",
			ConstCase:    "SCREAMING_CASE",
		}
	}
}

// ConvertCase converts a name to the specified case
func ConvertCase(name string, targetCase string) string {
	// First normalize to parts
	parts := splitIntoParts(name)

	switch targetCase {
	case "camelCase":
		return toCamelCase(parts)
	case "PascalCase":
		return toPascalCase(parts)
	case "snake_case":
		return toSnakeCase(parts)
	case "SCREAMING_CASE":
		return toScreamingCase(parts)
	default:
		return name
	}
}

// splitIntoParts splits an identifier into word parts
func splitIntoParts(name string) []string {
	// Handle snake_case
	if strings.Contains(name, "_") {
		return strings.Split(strings.ToLower(name), "_")
	}

	// Handle camelCase/PascalCase
	var parts []string
	re := regexp.MustCompile(`[A-Z][a-z]*|[a-z]+|[0-9]+`)
	matches := re.FindAllString(name, -1)
	for _, m := range matches {
		parts = append(parts, strings.ToLower(m))
	}
	return parts
}

func toCamelCase(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, p := range parts[1:] {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

func toPascalCase(parts []string) string {
	var result string
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

func toSnakeCase(parts []string) string {
	return strings.Join(parts, "_")
}

func toScreamingCase(parts []string) string {
	return strings.ToUpper(strings.Join(parts, "_"))
}
