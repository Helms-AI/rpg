package extractors

import (
	"regexp"
	"strings"
)

// GoExtractor extracts information from Go source files.
type GoExtractor struct{}

// NewGoExtractor creates a new Go extractor.
func NewGoExtractor() *GoExtractor {
	return &GoExtractor{}
}

// LanguageID returns the language identifier.
func (e *GoExtractor) LanguageID() string {
	return "go"
}

// Extensions returns Go file extensions.
func (e *GoExtractor) Extensions() []string {
	return []string{".go"}
}

// IsTestFile returns true if the file is a Go test file.
func (e *GoExtractor) IsTestFile(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

// ExtractPackageDescription extracts the package comment.
func (e *GoExtractor) ExtractPackageDescription(content string) string {
	// Look for package comment: // Package name ...
	packageCommentRe := regexp.MustCompile(`(?m)^//\s*Package\s+\w+\s+(.+)$`)
	if match := packageCommentRe.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	// Look for block comment before package declaration
	blockCommentRe := regexp.MustCompile(`(?s)/\*(.+?)\*/\s*package`)
	if match := blockCommentRe.FindStringSubmatch(content); len(match) > 1 {
		comment := strings.TrimSpace(match[1])
		// Remove leading asterisks from each line
		lines := strings.Split(comment, "\n")
		var cleaned []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
			if line != "" {
				cleaned = append(cleaned, line)
			}
		}
		return strings.Join(cleaned, " ")
	}

	return ""
}

// ExtractTypes extracts type definitions from Go source.
func (e *GoExtractor) ExtractTypes(content string, filename string) []ExtractedType {
	var types []ExtractedType
	lines := strings.Split(content, "\n")

	// Struct pattern: type Name struct { ... }
	structRe := regexp.MustCompile(`(?m)^type\s+(\w+)\s+struct\s*\{`)
	for _, match := range structRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Extract fields from struct body
		body := ExtractFunctionBody(content, match[0])
		fields := e.parseStructFields(body)

		// Extract comment above
		description := ExtractCommentAbove(lines, lineNum, "//", "/*", "*/")

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Interface pattern: type Name interface { ... }
	interfaceRe := regexp.MustCompile(`(?m)^type\s+(\w+)\s+interface\s*\{`)
	for _, match := range interfaceRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])
		description := ExtractCommentAbove(lines, lineNum, "//", "/*", "*/")

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "interface",
			Description: description,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Type alias pattern: type Name = OtherType
	aliasRe := regexp.MustCompile(`(?m)^type\s+(\w+)\s*=\s*(\S+)`)
	for _, match := range aliasRe.FindAllStringSubmatch(content, -1) {
		types = append(types, ExtractedType{
			Name:       match[1],
			Kind:       "alias",
			AliasOf:    match[2],
			SourceFile: filename,
		})
	}

	// Enum-like pattern: type Name int/string with const block
	// Look for: type Status string followed by const ( Status1 Status = "value" ... )
	enumTypeRe := regexp.MustCompile(`type\s+(\w+)\s+(int|string|uint|int64)`)
	for _, match := range enumTypeRe.FindAllStringSubmatch(content, -1) {
		typeName := match[1]
		// Look for const block with this type
		constBlockRe := regexp.MustCompile(`const\s*\(\s*([^)]+)\)`)
		for _, constMatch := range constBlockRe.FindAllStringSubmatch(content, -1) {
			constBody := constMatch[1]
			if strings.Contains(constBody, typeName) {
				variants := e.parseEnumVariants(constBody, typeName)
				if len(variants) > 0 {
					types = append(types, ExtractedType{
						Name:       typeName,
						Kind:       "enum",
						Variants:   variants,
						SourceFile: filename,
					})
				}
			}
		}
	}

	return types
}

// parseStructFields parses fields from a struct body.
func (e *GoExtractor) parseStructFields(body string) []ExtractedField {
	var fields []ExtractedField

	fieldRe := regexp.MustCompile(`(\w+)\s+(\S+)(?:\s+` + "`" + `[^` + "`" + `]+` + "`" + `)?`)
	lines := strings.Split(body, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if match := fieldRe.FindStringSubmatch(line); len(match) > 2 {
			fieldName := match[1]
			fieldType := match[2]

			// Check for optional (pointer types)
			optional := strings.HasPrefix(fieldType, "*")
			fieldType = strings.TrimPrefix(fieldType, "*")

			fields = append(fields, ExtractedField{
				Name:     fieldName,
				Type:     e.MapTypeToSpec(fieldType),
				Optional: optional,
			})
		}
	}

	return fields
}

// parseEnumVariants parses enum variants from a const block.
func (e *GoExtractor) parseEnumVariants(constBody string, typeName string) []string {
	var variants []string

	// Pattern: VariantName TypeName = "value" or VariantName TypeName = iota
	variantRe := regexp.MustCompile(`(\w+)\s+` + typeName)
	for _, match := range variantRe.FindAllStringSubmatch(constBody, -1) {
		variants = append(variants, match[1])
	}

	return variants
}

// ExtractFunctions extracts function definitions from Go source.
func (e *GoExtractor) ExtractFunctions(content string, filename string) []ExtractedFunction {
	var functions []ExtractedFunction
	lines := strings.Split(content, "\n")

	// Function pattern: func Name(params) returnType {
	// Also handles: func (r *Receiver) Name(params) returnType {
	funcRe := regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?(\w+)\s*\(([^)]*)\)\s*(?:\(([^)]+)\)|(\w+(?:\.\w+)?))?\s*\{`)

	for _, match := range funcRe.FindAllStringSubmatchIndex(content, -1) {
		funcName := content[match[2]:match[3]]

		// Skip test functions and internal functions
		if strings.HasPrefix(funcName, "Test") || strings.HasPrefix(funcName, "Benchmark") {
			continue
		}

		// Skip unexported functions (start with lowercase)
		if funcName[0] >= 'a' && funcName[0] <= 'z' {
			continue
		}

		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Extract parameters
		paramsStr := ""
		if match[4] != -1 && match[5] != -1 {
			paramsStr = content[match[4]:match[5]]
		}
		params := e.parseParams(paramsStr)

		// Extract return type
		returnType := ""
		if match[6] != -1 && match[7] != -1 {
			returnType = content[match[6]:match[7]]
		} else if match[8] != -1 && match[9] != -1 {
			returnType = content[match[8]:match[9]]
		}

		// Extract function body
		body := ExtractFunctionBody(content, match[0])

		// Get comment above
		description := ExtractCommentAbove(lines, lineNum, "//", "/*", "*/")

		// Infer logic from body
		logic := InferLogicFromBody(body)

		functions = append(functions, ExtractedFunction{
			Name:        funcName,
			Description: description,
			Parameters:  params,
			Returns:     e.MapTypeToSpec(returnType),
			Logic:       logic,
			Body:        body,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	return functions
}

// parseParams parses function parameters.
func (e *GoExtractor) parseParams(paramsStr string) []ExtractedParam {
	var params []ExtractedParam

	if strings.TrimSpace(paramsStr) == "" {
		return params
	}

	// Split by comma, but handle variadic params
	parts := strings.Split(paramsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Handle variadic: name ...type
		isVariadic := strings.Contains(part, "...")
		part = strings.Replace(part, "...", "", 1)

		// Split into name and type
		fields := strings.Fields(part)
		if len(fields) >= 2 {
			paramName := fields[0]
			paramType := fields[len(fields)-1]

			// Handle default for variadic (optional)
			defaultVal := ""
			if isVariadic {
				defaultVal = "optional"
			}

			params = append(params, ExtractedParam{
				Name:    paramName,
				Type:    e.MapTypeToSpec(paramType),
				Default: defaultVal,
			})
		} else if len(fields) == 1 {
			// Type only (previous param shares type)
			params = append(params, ExtractedParam{
				Type: e.MapTypeToSpec(fields[0]),
			})
		}
	}

	return params
}

// ExtractTests extracts test cases from Go test files.
func (e *GoExtractor) ExtractTests(content string, filename string) []ExtractedTest {
	var tests []ExtractedTest

	// Look for table-driven tests
	// Pattern: func TestXxx(t *testing.T) { tests := []struct { ... }
	testFuncRe := regexp.MustCompile(`func\s+(Test\w+)\s*\(\s*\w+\s+\*testing\.T\s*\)\s*\{`)

	for _, match := range testFuncRe.FindAllStringSubmatchIndex(content, -1) {
		testFuncName := content[match[2]:match[3]]
		// Derive function being tested from test name
		funcName := strings.TrimPrefix(testFuncName, "Test")

		// Extract test function body
		body := ExtractFunctionBody(content, match[0])

		// Look for table-driven test cases
		tableCases := e.parseTableDrivenTests(body, funcName)
		tests = append(tests, tableCases...)

		// If no table cases found, create a single test case
		if len(tableCases) == 0 {
			tests = append(tests, ExtractedTest{
				Function:   funcName,
				Name:       "test " + funcName,
				SourceFile: filename,
			})
		}
	}

	return tests
}

// parseTableDrivenTests parses table-driven test cases from a test function body.
func (e *GoExtractor) parseTableDrivenTests(body string, funcName string) []ExtractedTest {
	var tests []ExtractedTest

	// Look for test case definitions in table
	// Pattern: { name: "...", input: ..., want: ... }
	testCaseRe := regexp.MustCompile(`\{\s*name:\s*"([^"]+)"[^}]*?\}`)

	for _, match := range testCaseRe.FindAllStringSubmatch(body, -1) {
		testName := match[1]
		caseBody := match[0]

		// Try to extract given/expect values
		given := e.extractTestValue(caseBody, []string{"text", "input", "given", "args"})
		expect := e.extractTestValue(caseBody, []string{"want", "expect", "expected", "result"})

		tests = append(tests, ExtractedTest{
			Function: funcName,
			Name:     testName,
			Given:    given,
			Expect:   expect,
		})
	}

	return tests
}

// extractTestValue extracts a test value from a test case body.
func (e *GoExtractor) extractTestValue(caseBody string, fieldNames []string) interface{} {
	for _, name := range fieldNames {
		// Try string value
		re := regexp.MustCompile(name + `:\s*"([^"]*)"`)
		if match := re.FindStringSubmatch(caseBody); len(match) > 1 {
			return match[1]
		}

		// Try numeric value
		re = regexp.MustCompile(name + `:\s*(\d+)`)
		if match := re.FindStringSubmatch(caseBody); len(match) > 1 {
			return match[1]
		}

		// Try boolean value
		re = regexp.MustCompile(name + `:\s*(true|false)`)
		if match := re.FindStringSubmatch(caseBody); len(match) > 1 {
			return match[1]
		}
	}
	return nil
}

// MapTypeToSpec maps Go types to spec pseudo-types.
func (e *GoExtractor) MapTypeToSpec(goType string) string {
	goType = strings.TrimSpace(goType)

	// Handle slice/array types
	if strings.HasPrefix(goType, "[]") {
		innerType := e.MapTypeToSpec(strings.TrimPrefix(goType, "[]"))
		return "List of " + innerType
	}

	// Handle pointer types as optional
	if strings.HasPrefix(goType, "*") {
		innerType := e.MapTypeToSpec(strings.TrimPrefix(goType, "*"))
		return "Optional " + innerType
	}

	// Handle map types
	if strings.HasPrefix(goType, "map[") {
		return "Map"
	}

	// Direct type mappings
	typeMap := map[string]string{
		"string":        "Text",
		"int":           "Integer",
		"int8":          "Integer",
		"int16":         "Integer",
		"int32":         "Integer",
		"int64":         "Integer",
		"uint":          "Integer",
		"uint8":         "Integer",
		"uint16":        "Integer",
		"uint32":        "Integer",
		"uint64":        "Integer",
		"float32":       "Float",
		"float64":       "Float",
		"bool":          "Boolean",
		"time.Time":     "Timestamp",
		"time.Duration": "Duration",
		"error":         "Error",
		"interface{}":   "Any",
		"any":           "Any",
		"":              "Nothing",
	}

	if mapped, ok := typeMap[goType]; ok {
		return mapped
	}

	// Return as-is for custom types
	return goType
}
