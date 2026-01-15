package extractors

import (
	"regexp"
	"strings"
)

// TypeScriptExtractor extracts information from TypeScript source files.
type TypeScriptExtractor struct{}

// NewTypeScriptExtractor creates a new TypeScript extractor.
func NewTypeScriptExtractor() *TypeScriptExtractor {
	return &TypeScriptExtractor{}
}

// LanguageID returns the language identifier.
func (e *TypeScriptExtractor) LanguageID() string {
	return "typescript"
}

// Extensions returns TypeScript file extensions.
func (e *TypeScriptExtractor) Extensions() []string {
	return []string{".ts", ".tsx"}
}

// IsTestFile returns true if the file is a TypeScript test file.
func (e *TypeScriptExtractor) IsTestFile(filename string) bool {
	return strings.Contains(filename, ".test.") ||
		strings.Contains(filename, ".spec.") ||
		strings.Contains(filename, "__tests__")
}

// ExtractPackageDescription extracts the module description from JSDoc.
func (e *TypeScriptExtractor) ExtractPackageDescription(content string) string {
	// Look for file-level JSDoc: /** @module ... */ or /** @fileoverview ... */
	jsdocRe := regexp.MustCompile(`/\*\*\s*\n?\s*\*?\s*(?:@(?:module|fileoverview|file)\s+)?([^@*]+)`)
	if match := jsdocRe.FindStringSubmatch(content); len(match) > 1 {
		desc := strings.TrimSpace(match[1])
		desc = strings.ReplaceAll(desc, "\n", " ")
		desc = regexp.MustCompile(`\s+`).ReplaceAllString(desc, " ")
		return desc
	}
	return ""
}

// ExtractTypes extracts type definitions from TypeScript source.
func (e *TypeScriptExtractor) ExtractTypes(content string, filename string) []ExtractedType {
	var types []ExtractedType
	lines := strings.Split(content, "\n")

	// Interface pattern: interface Name { ... }
	interfaceRe := regexp.MustCompile(`(?:export\s+)?interface\s+(\w+)(?:\s+extends\s+[\w,\s]+)?\s*\{`)
	for _, match := range interfaceRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		fields := e.parseInterfaceFields(body)
		description := ExtractCommentAbove(lines, lineNum, "//", "/**", "*/")

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Type alias pattern: type Name = { ... } or type Name = OtherType
	typeAliasRe := regexp.MustCompile(`(?:export\s+)?type\s+(\w+)(?:<[^>]+>)?\s*=\s*`)
	for _, match := range typeAliasRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])
		description := ExtractCommentAbove(lines, lineNum, "//", "/**", "*/")

		// Check if it's an object type or union type
		afterEquals := content[match[1]:]
		if strings.HasPrefix(strings.TrimSpace(afterEquals), "{") {
			// Object type alias
			body := ExtractFunctionBody(content, match[0])
			fields := e.parseInterfaceFields(body)
			types = append(types, ExtractedType{
				Name:        name,
				Kind:        "struct",
				Description: description,
				Fields:      fields,
				SourceFile:  filename,
				LineNumber:  lineNum,
			})
		} else if strings.Contains(afterEquals[:min(100, len(afterEquals))], "|") {
			// Union type (enum-like)
			variants := e.parseUnionVariants(afterEquals)
			types = append(types, ExtractedType{
				Name:        name,
				Kind:        "enum",
				Description: description,
				Variants:    variants,
				SourceFile:  filename,
				LineNumber:  lineNum,
			})
		} else {
			// Simple type alias
			aliasRe := regexp.MustCompile(`^([^;\n{]+)`)
			if aliasMatch := aliasRe.FindStringSubmatch(afterEquals); len(aliasMatch) > 1 {
				types = append(types, ExtractedType{
					Name:        name,
					Kind:        "alias",
					AliasOf:     strings.TrimSpace(aliasMatch[1]),
					Description: description,
					SourceFile:  filename,
					LineNumber:  lineNum,
				})
			}
		}
	}

	// Enum pattern: enum Name { ... }
	enumRe := regexp.MustCompile(`(?:export\s+)?enum\s+(\w+)\s*\{`)
	for _, match := range enumRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		variants := e.parseEnumVariants(body)
		description := ExtractCommentAbove(lines, lineNum, "//", "/**", "*/")

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "enum",
			Description: description,
			Variants:    variants,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	return types
}

// parseInterfaceFields parses fields from an interface body.
func (e *TypeScriptExtractor) parseInterfaceFields(body string) []ExtractedField {
	var fields []ExtractedField

	// Field pattern: name?: Type; or name: Type;
	fieldRe := regexp.MustCompile(`(\w+)(\?)?\s*:\s*([^;,\n]+)`)

	for _, match := range fieldRe.FindAllStringSubmatch(body, -1) {
		fieldName := match[1]
		optional := match[2] == "?"
		fieldType := strings.TrimSpace(match[3])

		fields = append(fields, ExtractedField{
			Name:     fieldName,
			Type:     e.MapTypeToSpec(fieldType),
			Optional: optional,
		})
	}

	return fields
}

// parseUnionVariants parses variants from a union type.
func (e *TypeScriptExtractor) parseUnionVariants(content string) []string {
	var variants []string

	// Extract until semicolon or newline
	endIdx := strings.IndexAny(content, ";\n")
	if endIdx == -1 {
		endIdx = min(200, len(content))
	}
	unionStr := content[:endIdx]

	// Split by |
	parts := strings.Split(unionStr, "|")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		// Remove quotes for string literals
		part = strings.Trim(part, `"'`)
		if part != "" {
			variants = append(variants, part)
		}
	}

	return variants
}

// parseEnumVariants parses variants from an enum body.
func (e *TypeScriptExtractor) parseEnumVariants(body string) []string {
	var variants []string

	// Variant pattern: Name = value or just Name
	variantRe := regexp.MustCompile(`(\w+)(?:\s*=)?`)
	for _, match := range variantRe.FindAllStringSubmatch(body, -1) {
		variants = append(variants, match[1])
	}

	return variants
}

// ExtractFunctions extracts function definitions from TypeScript source.
func (e *TypeScriptExtractor) ExtractFunctions(content string, filename string) []ExtractedFunction {
	var functions []ExtractedFunction
	lines := strings.Split(content, "\n")

	// Function declaration pattern: function name(params): ReturnType {
	funcDeclRe := regexp.MustCompile(`(?:export\s+)?(?:async\s+)?function\s+(\w+)(?:<[^>]+>)?\s*\(([^)]*)\)\s*(?::\s*([^{]+))?\s*\{`)

	for _, match := range funcDeclRe.FindAllStringSubmatchIndex(content, -1) {
		funcName := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		paramsStr := ""
		if match[4] != -1 && match[5] != -1 {
			paramsStr = content[match[4]:match[5]]
		}

		returnType := ""
		if match[6] != -1 && match[7] != -1 {
			returnType = strings.TrimSpace(content[match[6]:match[7]])
		}

		isAsync := strings.Contains(content[match[0]:match[1]], "async ")
		body := ExtractFunctionBody(content, match[0])
		description := ExtractCommentAbove(lines, lineNum, "//", "/**", "*/")
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

	// Arrow function pattern: const name = (params): ReturnType => {
	arrowFuncRe := regexp.MustCompile(`(?:export\s+)?const\s+(\w+)\s*=\s*(?:async\s+)?\(([^)]*)\)\s*(?::\s*([^=]+))?\s*=>`)

	for _, match := range arrowFuncRe.FindAllStringSubmatchIndex(content, -1) {
		funcName := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		paramsStr := ""
		if match[4] != -1 && match[5] != -1 {
			paramsStr = content[match[4]:match[5]]
		}

		returnType := ""
		if match[6] != -1 && match[7] != -1 {
			returnType = strings.TrimSpace(content[match[6]:match[7]])
		}

		isAsync := strings.Contains(content[match[0]:match[1]], "async ")
		body := ExtractFunctionBody(content, match[0])
		description := ExtractCommentAbove(lines, lineNum, "//", "/**", "*/")
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

// parseParams parses TypeScript function parameters.
func (e *TypeScriptExtractor) parseParams(paramsStr string) []ExtractedParam {
	var params []ExtractedParam

	if strings.TrimSpace(paramsStr) == "" {
		return params
	}

	// Split by comma (careful with generic types)
	parts := splitParams(paramsStr)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Pattern: name?: Type = default or name: Type
		paramRe := regexp.MustCompile(`(\w+)(\?)?\s*:\s*([^=]+)(?:\s*=\s*(.+))?`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 1 {
			paramName := match[1]
			optional := match[2] == "?"
			paramType := strings.TrimSpace(match[3])
			defaultVal := ""
			if len(match) > 4 && match[4] != "" {
				defaultVal = strings.TrimSpace(match[4])
			}

			// Optional params have default
			if optional && defaultVal == "" {
				defaultVal = "optional"
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

// ExtractTests extracts test cases from TypeScript test files.
func (e *TypeScriptExtractor) ExtractTests(content string, filename string) []ExtractedTest {
	var tests []ExtractedTest

	// Jest/Vitest patterns: it("name", ...) or test("name", ...)
	testRe := regexp.MustCompile(`(?:it|test)\s*\(\s*['"]([^'"]+)['"]`)

	// Find describe blocks to get function being tested
	describeRe := regexp.MustCompile(`describe\s*\(\s*['"]([^'"]+)['"]`)
	currentFunction := ""

	// Simple approach: find all describe and test/it blocks
	describeMatches := describeRe.FindAllStringSubmatch(content, -1)
	if len(describeMatches) > 0 {
		currentFunction = describeMatches[0][1]
	}

	for _, match := range testRe.FindAllStringSubmatch(content, -1) {
		testName := match[1]

		tests = append(tests, ExtractedTest{
			Function:   currentFunction,
			Name:       testName,
			SourceFile: filename,
		})
	}

	return tests
}

// MapTypeToSpec maps TypeScript types to spec pseudo-types.
func (e *TypeScriptExtractor) MapTypeToSpec(tsType string) string {
	tsType = strings.TrimSpace(tsType)

	// Handle array types
	if strings.HasSuffix(tsType, "[]") {
		innerType := e.MapTypeToSpec(strings.TrimSuffix(tsType, "[]"))
		return "List of " + innerType
	}
	if strings.HasPrefix(tsType, "Array<") {
		innerType := strings.TrimPrefix(tsType, "Array<")
		innerType = strings.TrimSuffix(innerType, ">")
		return "List of " + e.MapTypeToSpec(innerType)
	}

	// Handle optional/nullable types
	if strings.HasSuffix(tsType, "| null") || strings.HasSuffix(tsType, "| undefined") {
		innerType := strings.TrimSuffix(tsType, "| null")
		innerType = strings.TrimSuffix(innerType, "| undefined")
		return "Optional " + e.MapTypeToSpec(strings.TrimSpace(innerType))
	}

	// Handle Promise types
	if strings.HasPrefix(tsType, "Promise<") {
		innerType := strings.TrimPrefix(tsType, "Promise<")
		innerType = strings.TrimSuffix(innerType, ">")
		return e.MapTypeToSpec(innerType)
	}

	// Direct type mappings
	typeMap := map[string]string{
		"string":    "Text",
		"number":    "Number",
		"boolean":   "Boolean",
		"Date":      "Timestamp",
		"void":      "Nothing",
		"undefined": "Nothing",
		"null":      "Nothing",
		"unknown":   "Any",
		"any":       "Any",
		"never":     "Nothing",
		"object":    "Any",
		"":          "Nothing",
	}

	if mapped, ok := typeMap[tsType]; ok {
		return mapped
	}

	return tsType
}

// splitParams splits parameters handling nested generics.
func splitParams(params string) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, ch := range params {
		if ch == '<' || ch == '(' || ch == '[' || ch == '{' {
			depth++
		} else if ch == '>' || ch == ')' || ch == ']' || ch == '}' {
			depth--
		}

		if ch == ',' && depth == 0 {
			result = append(result, current.String())
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
