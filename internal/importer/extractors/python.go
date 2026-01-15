package extractors

import (
	"regexp"
	"strings"
)

// PythonExtractor extracts information from Python source files.
type PythonExtractor struct{}

// NewPythonExtractor creates a new Python extractor.
func NewPythonExtractor() *PythonExtractor {
	return &PythonExtractor{}
}

// LanguageID returns the language identifier.
func (e *PythonExtractor) LanguageID() string {
	return "python"
}

// Extensions returns Python file extensions.
func (e *PythonExtractor) Extensions() []string {
	return []string{".py"}
}

// IsTestFile returns true if the file is a Python test file.
func (e *PythonExtractor) IsTestFile(filename string) bool {
	return strings.HasPrefix(filename, "test_") ||
		strings.HasSuffix(filename, "_test.py") ||
		strings.Contains(filename, "/tests/")
}

// ExtractPackageDescription extracts the module docstring.
func (e *PythonExtractor) ExtractPackageDescription(content string) string {
	// Look for module docstring at the beginning
	docstringRe := regexp.MustCompile(`^(?:#!/[^\n]*\n)?(?:#[^\n]*\n)*\s*(?:"""([^"]*)"""|'''([^']*)''')`)
	if match := docstringRe.FindStringSubmatch(content); len(match) > 1 {
		if match[1] != "" {
			return strings.TrimSpace(match[1])
		}
		if match[2] != "" {
			return strings.TrimSpace(match[2])
		}
	}
	return ""
}

// ExtractTypes extracts type definitions from Python source.
func (e *PythonExtractor) ExtractTypes(content string, filename string) []ExtractedType {
	var types []ExtractedType
	lines := strings.Split(content, "\n")

	// Dataclass pattern: @dataclass class Name:
	dataclassRe := regexp.MustCompile(`@dataclass(?:\([^)]*\))?\s*\nclass\s+(\w+)(?:\([^)]*\))?:`)
	for _, match := range dataclassRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Extract fields from class body
		classBody := e.extractPythonClassBody(content, match[1])
		fields := e.parseDataclassFields(classBody)
		description := e.extractDocstring(lines, lineNum+1) // Docstring is after class line

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Regular class pattern: class Name:
	classRe := regexp.MustCompile(`(?m)^class\s+(\w+)(?:\([^)]*\))?:`)
	for _, match := range classRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Skip if already captured as dataclass
		alreadyCaptured := false
		for _, t := range types {
			if t.Name == name {
				alreadyCaptured = true
				break
			}
		}
		if alreadyCaptured {
			continue
		}

		description := e.extractDocstring(lines, lineNum+1)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// TypedDict pattern: class Name(TypedDict):
	typedDictRe := regexp.MustCompile(`class\s+(\w+)\s*\(\s*TypedDict\s*\):`)
	for _, match := range typedDictRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		classBody := e.extractPythonClassBody(content, match[1])
		fields := e.parseTypedDictFields(classBody)

		types = append(types, ExtractedType{
			Name:       name,
			Kind:       "struct",
			Fields:     fields,
			SourceFile: filename,
			LineNumber: lineNum,
		})
	}

	// Enum pattern: class Name(Enum):
	enumRe := regexp.MustCompile(`class\s+(\w+)\s*\(\s*(?:Enum|IntEnum|StrEnum)\s*\):`)
	for _, match := range enumRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		classBody := e.extractPythonClassBody(content, match[1])
		variants := e.parseEnumVariants(classBody)

		types = append(types, ExtractedType{
			Name:       name,
			Kind:       "enum",
			Variants:   variants,
			SourceFile: filename,
			LineNumber: lineNum,
		})
	}

	return types
}

// extractPythonClassBody extracts the body of a Python class.
func (e *PythonExtractor) extractPythonClassBody(content string, startIdx int) string {
	lines := strings.Split(content[startIdx:], "\n")
	if len(lines) == 0 {
		return ""
	}

	var body strings.Builder
	baseIndent := -1

	for i, line := range lines {
		if i == 0 {
			continue // Skip the class line
		}

		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			body.WriteString("\n")
			continue
		}

		indent := len(line) - len(trimmed)
		if baseIndent == -1 {
			baseIndent = indent
		}

		// If we encounter a line with less or equal indent (and not empty), we've left the class
		if indent <= 0 && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			break
		}

		body.WriteString(line)
		body.WriteString("\n")
	}

	return body.String()
}

// parseDataclassFields parses fields from a dataclass body.
func (e *PythonExtractor) parseDataclassFields(body string) []ExtractedField {
	var fields []ExtractedField

	// Field pattern: name: Type = default or name: Type
	fieldRe := regexp.MustCompile(`(\w+)\s*:\s*([^=\n]+)(?:\s*=\s*(.+))?`)

	for _, match := range fieldRe.FindAllStringSubmatch(body, -1) {
		fieldName := match[1]
		fieldType := strings.TrimSpace(match[2])

		// Check for Optional
		optional := strings.HasPrefix(fieldType, "Optional[")
		if optional {
			fieldType = strings.TrimPrefix(fieldType, "Optional[")
			fieldType = strings.TrimSuffix(fieldType, "]")
		}

		defaultVal := ""
		if len(match) > 3 && match[3] != "" {
			defaultVal = strings.TrimSpace(match[3])
		}

		fields = append(fields, ExtractedField{
			Name:     fieldName,
			Type:     e.MapTypeToSpec(fieldType),
			Optional: optional,
			Default:  defaultVal,
		})
	}

	return fields
}

// parseTypedDictFields parses fields from a TypedDict body.
func (e *PythonExtractor) parseTypedDictFields(body string) []ExtractedField {
	return e.parseDataclassFields(body) // Same pattern
}

// parseEnumVariants parses variants from an Enum class body.
func (e *PythonExtractor) parseEnumVariants(body string) []string {
	var variants []string

	// Variant pattern: NAME = value
	variantRe := regexp.MustCompile(`([A-Z_][A-Z0-9_]*)\s*=`)
	for _, match := range variantRe.FindAllStringSubmatch(body, -1) {
		variants = append(variants, match[1])
	}

	return variants
}

// extractDocstring extracts a docstring after a given line.
func (e *PythonExtractor) extractDocstring(lines []string, lineNum int) string {
	if lineNum < 1 || lineNum > len(lines) {
		return ""
	}

	// Look for docstring on the next few lines
	for i := lineNum - 1; i < min(lineNum+2, len(lines)); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, `"""`) || strings.HasPrefix(line, `'''`) {
			quote := line[:3]
			if strings.HasSuffix(line, quote) && len(line) > 6 {
				// Single line docstring
				return strings.Trim(line, `"'`)
			}
			// Multi-line docstring
			var docLines []string
			docLines = append(docLines, strings.TrimPrefix(line, quote))
			for j := i + 1; j < len(lines); j++ {
				if strings.Contains(lines[j], quote) {
					docLines = append(docLines, strings.TrimSuffix(strings.TrimSpace(lines[j]), quote))
					break
				}
				docLines = append(docLines, strings.TrimSpace(lines[j]))
			}
			return strings.Join(docLines, " ")
		}
	}

	return ""
}

// ExtractFunctions extracts function definitions from Python source.
func (e *PythonExtractor) ExtractFunctions(content string, filename string) []ExtractedFunction {
	var functions []ExtractedFunction
	lines := strings.Split(content, "\n")

	// Function pattern: def name(params) -> ReturnType:
	funcRe := regexp.MustCompile(`(?m)^(?:async\s+)?def\s+(\w+)\s*\(([^)]*)\)\s*(?:->\s*([^:]+))?:`)

	for _, match := range funcRe.FindAllStringSubmatchIndex(content, -1) {
		funcName := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Skip private functions and test functions
		if strings.HasPrefix(funcName, "_") || strings.HasPrefix(funcName, "test_") {
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

		isAsync := strings.Contains(content[match[0]:match[1]], "async ")

		// Extract function body
		body := e.extractPythonFunctionBody(content, match[1])
		description := e.extractDocstring(lines, lineNum+1)
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

// extractPythonFunctionBody extracts the body of a Python function.
func (e *PythonExtractor) extractPythonFunctionBody(content string, startIdx int) string {
	return e.extractPythonClassBody(content, startIdx) // Same indentation logic
}

// parseParams parses Python function parameters.
func (e *PythonExtractor) parseParams(paramsStr string) []ExtractedParam {
	var params []ExtractedParam

	if strings.TrimSpace(paramsStr) == "" {
		return params
	}

	// Remove self/cls parameter
	paramsStr = regexp.MustCompile(`^\s*(?:self|cls)\s*,?\s*`).ReplaceAllString(paramsStr, "")

	parts := splitParams(paramsStr)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "self" || part == "cls" {
			continue
		}

		// Skip *args and **kwargs
		if strings.HasPrefix(part, "*") {
			continue
		}

		// Pattern: name: Type = default or name = default or just name
		paramRe := regexp.MustCompile(`(\w+)(?:\s*:\s*([^=]+))?(?:\s*=\s*(.+))?`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 1 {
			paramName := match[1]
			paramType := ""
			if len(match) > 2 {
				paramType = strings.TrimSpace(match[2])
			}
			defaultVal := ""
			if len(match) > 3 && match[3] != "" {
				defaultVal = strings.TrimSpace(match[3])
			}

			params = append(params, ExtractedParam{
				Name:    paramName,
				Type:    e.MapTypeToSpec(paramType),
				Default: defaultVal,
			})
		}
	}

	return params
}

// ExtractTests extracts test cases from Python test files.
func (e *PythonExtractor) ExtractTests(content string, filename string) []ExtractedTest {
	var tests []ExtractedTest

	// pytest pattern: def test_name():
	testRe := regexp.MustCompile(`def\s+(test_\w+)\s*\(`)

	for _, match := range testRe.FindAllStringSubmatch(content, -1) {
		testFuncName := match[1]
		// Derive function being tested
		funcName := strings.TrimPrefix(testFuncName, "test_")
		// Convert snake_case to readable name
		testName := strings.ReplaceAll(funcName, "_", " ")

		tests = append(tests, ExtractedTest{
			Function:   funcName,
			Name:       testName,
			SourceFile: filename,
		})
	}

	return tests
}

// MapTypeToSpec maps Python types to spec pseudo-types.
func (e *PythonExtractor) MapTypeToSpec(pyType string) string {
	pyType = strings.TrimSpace(pyType)

	// Handle List types
	if strings.HasPrefix(pyType, "List[") || strings.HasPrefix(pyType, "list[") {
		innerType := strings.TrimPrefix(pyType, "List[")
		innerType = strings.TrimPrefix(innerType, "list[")
		innerType = strings.TrimSuffix(innerType, "]")
		return "List of " + e.MapTypeToSpec(innerType)
	}

	// Handle Optional types
	if strings.HasPrefix(pyType, "Optional[") {
		innerType := strings.TrimPrefix(pyType, "Optional[")
		innerType = strings.TrimSuffix(innerType, "]")
		return "Optional " + e.MapTypeToSpec(innerType)
	}

	// Handle Dict types
	if strings.HasPrefix(pyType, "Dict[") || strings.HasPrefix(pyType, "dict[") {
		return "Map"
	}

	// Direct type mappings
	typeMap := map[string]string{
		"str":      "Text",
		"int":      "Integer",
		"float":    "Float",
		"bool":     "Boolean",
		"datetime": "Timestamp",
		"date":     "Timestamp",
		"None":     "Nothing",
		"Any":      "Any",
		"bytes":    "Bytes",
		"":         "Any",
	}

	if mapped, ok := typeMap[pyType]; ok {
		return mapped
	}

	return pyType
}
