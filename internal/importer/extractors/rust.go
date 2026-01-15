package extractors

import (
	"regexp"
	"strings"
)

// RustExtractor extracts information from Rust source files.
type RustExtractor struct{}

// NewRustExtractor creates a new Rust extractor.
func NewRustExtractor() *RustExtractor {
	return &RustExtractor{}
}

// LanguageID returns the language identifier.
func (e *RustExtractor) LanguageID() string {
	return "rust"
}

// Extensions returns Rust file extensions.
func (e *RustExtractor) Extensions() []string {
	return []string{".rs"}
}

// IsTestFile returns true if the file is a Rust test file.
func (e *RustExtractor) IsTestFile(filename string) bool {
	return strings.Contains(filename, "/tests/") ||
		strings.HasSuffix(filename, "_test.rs")
}

// ExtractPackageDescription extracts the module doc comment.
func (e *RustExtractor) ExtractPackageDescription(content string) string {
	// Look for //! module doc comments at the start
	docRe := regexp.MustCompile(`(?m)^//!\s*(.+)$`)
	matches := docRe.FindAllStringSubmatch(content, 5)
	if len(matches) > 0 {
		var lines []string
		for _, m := range matches {
			lines = append(lines, m[1])
		}
		return strings.Join(lines, " ")
	}
	return ""
}

// ExtractTypes extracts type definitions from Rust source.
func (e *RustExtractor) ExtractTypes(content string, filename string) []ExtractedType {
	var types []ExtractedType
	lines := strings.Split(content, "\n")

	// Struct pattern: pub struct Name { ... }
	structRe := regexp.MustCompile(`(?:pub\s+)?struct\s+(\w+)(?:<[^>]+>)?\s*\{`)
	for _, match := range structRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		fields := e.parseStructFields(body)
		description := e.extractRustDocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Tuple struct pattern: pub struct Name(Type, Type);
	tupleStructRe := regexp.MustCompile(`(?:pub\s+)?struct\s+(\w+)\s*\(([^)]+)\)\s*;`)
	for _, match := range tupleStructRe.FindAllStringSubmatch(content, -1) {
		name := match[1]
		fieldsStr := match[2]

		var fields []ExtractedField
		parts := strings.Split(fieldsStr, ",")
		for i, part := range parts {
			part = strings.TrimSpace(part)
			part = strings.TrimPrefix(part, "pub ")
			if part != "" {
				fields = append(fields, ExtractedField{
					Name: string(rune('0' + i)),
					Type: e.MapTypeToSpec(part),
				})
			}
		}

		types = append(types, ExtractedType{
			Name:       name,
			Kind:       "struct",
			Fields:     fields,
			SourceFile: filename,
		})
	}

	// Enum pattern: pub enum Name { ... }
	enumRe := regexp.MustCompile(`(?:pub\s+)?enum\s+(\w+)(?:<[^>]+>)?\s*\{`)
	for _, match := range enumRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		variants := e.parseEnumVariants(body)
		description := e.extractRustDocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "enum",
			Description: description,
			Variants:    variants,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Type alias pattern: pub type Name = OtherType;
	aliasRe := regexp.MustCompile(`(?:pub\s+)?type\s+(\w+)(?:<[^>]+>)?\s*=\s*([^;]+);`)
	for _, match := range aliasRe.FindAllStringSubmatch(content, -1) {
		types = append(types, ExtractedType{
			Name:       match[1],
			Kind:       "alias",
			AliasOf:    strings.TrimSpace(match[2]),
			SourceFile: filename,
		})
	}

	return types
}

// parseStructFields parses fields from a Rust struct body.
func (e *RustExtractor) parseStructFields(body string) []ExtractedField {
	var fields []ExtractedField

	// Field pattern: pub name: Type, or name: Type,
	fieldRe := regexp.MustCompile(`(?:pub\s+)?(\w+)\s*:\s*([^,\n}]+)`)

	for _, match := range fieldRe.FindAllStringSubmatch(body, -1) {
		fieldName := match[1]
		fieldType := strings.TrimSpace(match[2])

		// Check for Option<T>
		optional := strings.HasPrefix(fieldType, "Option<")
		if optional {
			fieldType = strings.TrimPrefix(fieldType, "Option<")
			fieldType = strings.TrimSuffix(fieldType, ">")
		}

		fields = append(fields, ExtractedField{
			Name:     fieldName,
			Type:     e.MapTypeToSpec(fieldType),
			Optional: optional,
		})
	}

	return fields
}

// parseEnumVariants parses variants from a Rust enum body.
func (e *RustExtractor) parseEnumVariants(body string) []string {
	var variants []string

	// Variant pattern: Name or Name(fields) or Name { fields }
	variantRe := regexp.MustCompile(`(\w+)(?:\s*[({])?`)

	for _, match := range variantRe.FindAllStringSubmatch(body, -1) {
		variant := match[1]
		// Skip common keywords
		if variant == "pub" || variant == "struct" || variant == "fn" {
			continue
		}
		variants = append(variants, variant)
	}

	return variants
}

// extractRustDocAbove extracts /// doc comments above a line.
func (e *RustExtractor) extractRustDocAbove(lines []string, lineNum int) string {
	if lineNum <= 1 || lineNum > len(lines) {
		return ""
	}

	var docLines []string
	for i := lineNum - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "///") {
			doc := strings.TrimPrefix(line, "///")
			doc = strings.TrimSpace(doc)
			docLines = append([]string{doc}, docLines...)
		} else if strings.HasPrefix(line, "#[") {
			// Skip attributes
			continue
		} else if line == "" {
			continue
		} else {
			break
		}
	}

	return strings.Join(docLines, " ")
}

// ExtractFunctions extracts function definitions from Rust source.
func (e *RustExtractor) ExtractFunctions(content string, filename string) []ExtractedFunction {
	var functions []ExtractedFunction
	lines := strings.Split(content, "\n")

	// Function pattern: pub fn name(params) -> ReturnType {
	funcRe := regexp.MustCompile(`(?:pub\s+)?(?:async\s+)?fn\s+(\w+)(?:<[^>]+>)?\s*\(([^)]*)\)\s*(?:->\s*([^{]+))?\s*\{`)

	for _, match := range funcRe.FindAllStringSubmatchIndex(content, -1) {
		funcName := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Skip private functions (not pub) and test functions
		fullMatch := content[match[0]:match[1]]
		if !strings.Contains(fullMatch, "pub ") {
			continue
		}

		paramsStr := ""
		if match[4] != -1 && match[5] != -1 {
			paramsStr = content[match[4]:match[5]]
		}

		returnType := ""
		if match[6] != -1 && match[7] != -1 {
			returnType = strings.TrimSpace(content[match[6]:match[7]])
		}

		isAsync := strings.Contains(fullMatch, "async ")
		body := ExtractFunctionBody(content, match[0])
		description := e.extractRustDocAbove(lines, lineNum)
		logic := InferLogicFromBody(body)

		functions = append(functions, ExtractedFunction{
			Name:        funcName,
			Description: description,
			Parameters:  e.parseParams(paramsStr),
			Returns:     e.MapTypeToSpec(returnType),
			IsAsync:     isAsync,
			Logic:       logic,
			Body:        body,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	return functions
}

// parseParams parses Rust function parameters.
func (e *RustExtractor) parseParams(paramsStr string) []ExtractedParam {
	var params []ExtractedParam

	if strings.TrimSpace(paramsStr) == "" {
		return params
	}

	// Remove self parameter
	paramsStr = regexp.MustCompile(`&?(?:mut\s+)?self\s*,?\s*`).ReplaceAllString(paramsStr, "")

	parts := splitParams(paramsStr)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Pattern: name: Type or mut name: Type
		paramRe := regexp.MustCompile(`(?:mut\s+)?(\w+)\s*:\s*(.+)`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 2 {
			paramName := match[1]
			paramType := strings.TrimSpace(match[2])

			// Remove reference markers for cleaner display
			paramType = strings.TrimPrefix(paramType, "&")
			paramType = strings.TrimPrefix(paramType, "mut ")

			params = append(params, ExtractedParam{
				Name: paramName,
				Type: e.MapTypeToSpec(paramType),
			})
		}
	}

	return params
}

// ExtractTests extracts test cases from Rust test files.
func (e *RustExtractor) ExtractTests(content string, filename string) []ExtractedTest {
	var tests []ExtractedTest

	// Rust test pattern: #[test] fn test_name()
	testRe := regexp.MustCompile(`#\[test\]\s*(?:async\s+)?fn\s+(\w+)`)

	for _, match := range testRe.FindAllStringSubmatch(content, -1) {
		testFuncName := match[1]
		// Convert snake_case to readable name
		testName := strings.ReplaceAll(testFuncName, "_", " ")

		// Try to derive function being tested
		funcName := strings.TrimPrefix(testFuncName, "test_")

		tests = append(tests, ExtractedTest{
			Function:   funcName,
			Name:       testName,
			SourceFile: filename,
		})
	}

	return tests
}

// MapTypeToSpec maps Rust types to spec pseudo-types.
func (e *RustExtractor) MapTypeToSpec(rustType string) string {
	rustType = strings.TrimSpace(rustType)

	// Handle Vec types
	if strings.HasPrefix(rustType, "Vec<") {
		innerType := strings.TrimPrefix(rustType, "Vec<")
		innerType = strings.TrimSuffix(innerType, ">")
		return "List of " + e.MapTypeToSpec(innerType)
	}

	// Handle Option types
	if strings.HasPrefix(rustType, "Option<") {
		innerType := strings.TrimPrefix(rustType, "Option<")
		innerType = strings.TrimSuffix(innerType, ">")
		return "Optional " + e.MapTypeToSpec(innerType)
	}

	// Handle Result types
	if strings.HasPrefix(rustType, "Result<") {
		// Extract the Ok type
		innerType := strings.TrimPrefix(rustType, "Result<")
		commaIdx := strings.Index(innerType, ",")
		if commaIdx > 0 {
			innerType = innerType[:commaIdx]
		}
		return "Result of " + e.MapTypeToSpec(innerType)
	}

	// Handle HashMap types
	if strings.HasPrefix(rustType, "HashMap<") || strings.HasPrefix(rustType, "BTreeMap<") {
		return "Map"
	}

	// Handle reference types
	rustType = strings.TrimPrefix(rustType, "&")
	rustType = strings.TrimPrefix(rustType, "'static ")
	rustType = strings.TrimPrefix(rustType, "'_ ")

	// Direct type mappings
	typeMap := map[string]string{
		"String":  "Text",
		"str":     "Text",
		"&str":    "Text",
		"i8":      "Integer",
		"i16":     "Integer",
		"i32":     "Integer",
		"i64":     "Integer",
		"i128":    "Integer",
		"isize":   "Integer",
		"u8":      "Integer",
		"u16":     "Integer",
		"u32":     "Integer",
		"u64":     "Integer",
		"u128":    "Integer",
		"usize":   "Integer",
		"f32":     "Float",
		"f64":     "Float",
		"bool":    "Boolean",
		"()":      "Nothing",
		"":        "Nothing",
	}

	if mapped, ok := typeMap[rustType]; ok {
		return mapped
	}

	return rustType
}
