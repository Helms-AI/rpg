package extractors

import (
	"regexp"
	"strings"
)

// CSharpExtractor extracts information from C# source files.
type CSharpExtractor struct{}

// NewCSharpExtractor creates a new C# extractor.
func NewCSharpExtractor() *CSharpExtractor {
	return &CSharpExtractor{}
}

// LanguageID returns the language identifier.
func (e *CSharpExtractor) LanguageID() string {
	return "csharp"
}

// Extensions returns C# file extensions.
func (e *CSharpExtractor) Extensions() []string {
	return []string{".cs"}
}

// IsTestFile returns true if the file is a C# test file.
func (e *CSharpExtractor) IsTestFile(filename string) bool {
	return strings.HasSuffix(filename, "Tests.cs") ||
		strings.HasSuffix(filename, "Test.cs") ||
		strings.Contains(filename, ".Tests/") ||
		strings.Contains(filename, "/Tests/")
}

// ExtractPackageDescription extracts the class XML documentation.
func (e *CSharpExtractor) ExtractPackageDescription(content string) string {
	// Look for /// <summary> before class declaration
	summaryRe := regexp.MustCompile(`///\s*<summary>\s*\n?///?\s*([^<]+)`)
	if match := summaryRe.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// ExtractTypes extracts type definitions from C# source.
func (e *CSharpExtractor) ExtractTypes(content string, filename string) []ExtractedType {
	var types []ExtractedType
	lines := strings.Split(content, "\n")

	// Class pattern: public class Name { ... }
	classRe := regexp.MustCompile(`(?:public|internal)\s+(?:sealed\s+)?(?:partial\s+)?(?:abstract\s+)?class\s+(\w+)(?:<[^>]+>)?(?:\s*:\s*[\w,\s<>]+)?\s*\{`)
	for _, match := range classRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		fields := e.parseClassFields(body)
		description := e.extractXmlDocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Record pattern: public record Name(fields) or public record Name { ... }
	recordRe := regexp.MustCompile(`(?:public|internal)\s+(?:sealed\s+)?record\s+(?:struct\s+)?(\w+)(?:<[^>]+>)?\s*(?:\(([^)]*)\))?`)
	for _, match := range recordRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		var fields []ExtractedField
		if match[4] != -1 && match[5] != -1 {
			fieldsStr := content[match[4]:match[5]]
			fields = e.parseRecordFields(fieldsStr)
		}

		description := e.extractXmlDocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Struct pattern: public struct Name { ... }
	structRe := regexp.MustCompile(`(?:public|internal)\s+(?:readonly\s+)?struct\s+(\w+)(?:<[^>]+>)?\s*\{`)
	for _, match := range structRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		fields := e.parseClassFields(body)

		types = append(types, ExtractedType{
			Name:       name,
			Kind:       "struct",
			Fields:     fields,
			SourceFile: filename,
			LineNumber: lineNum,
		})
	}

	// Interface pattern: public interface IName { ... }
	interfaceRe := regexp.MustCompile(`(?:public|internal)\s+interface\s+(I?\w+)(?:<[^>]+>)?(?:\s*:\s*[\w,\s<>]+)?\s*\{`)
	for _, match := range interfaceRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])
		description := e.extractXmlDocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "interface",
			Description: description,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Enum pattern: public enum Name { ... }
	enumRe := regexp.MustCompile(`(?:public|internal)\s+enum\s+(\w+)\s*\{`)
	for _, match := range enumRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		variants := e.parseEnumVariants(body)

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

// parseClassFields parses fields/properties from a C# class body.
func (e *CSharpExtractor) parseClassFields(body string) []ExtractedField {
	var fields []ExtractedField

	// Property pattern: public Type Name { get; set; }
	propRe := regexp.MustCompile(`(?:public|internal)\s+(?:required\s+)?(\w+(?:<[^>]+>)?(?:\?)?)\s+(\w+)\s*\{`)

	for _, match := range propRe.FindAllStringSubmatch(body, -1) {
		propType := match[1]
		propName := match[2]

		// Check for nullable
		optional := strings.HasSuffix(propType, "?")
		propType = strings.TrimSuffix(propType, "?")

		fields = append(fields, ExtractedField{
			Name:     propName,
			Type:     e.MapTypeToSpec(propType),
			Optional: optional,
		})
	}

	return fields
}

// parseRecordFields parses fields from a record primary constructor.
func (e *CSharpExtractor) parseRecordFields(fieldsStr string) []ExtractedField {
	var fields []ExtractedField

	parts := splitParams(fieldsStr)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Pattern: Type Name or Type? Name
		paramRe := regexp.MustCompile(`(\w+(?:<[^>]+>)?(?:\?)?)\s+(\w+)`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 2 {
			propType := match[1]
			optional := strings.HasSuffix(propType, "?")
			propType = strings.TrimSuffix(propType, "?")

			fields = append(fields, ExtractedField{
				Name:     match[2],
				Type:     e.MapTypeToSpec(propType),
				Optional: optional,
			})
		}
	}

	return fields
}

// parseEnumVariants parses variants from a C# enum body.
func (e *CSharpExtractor) parseEnumVariants(body string) []string {
	var variants []string

	// Variant pattern: Name = value or just Name
	variantRe := regexp.MustCompile(`([A-Z]\w*)\s*(?:=|,|$)`)
	for _, match := range variantRe.FindAllStringSubmatch(body, -1) {
		variants = append(variants, match[1])
	}

	return variants
}

// extractXmlDocAbove extracts /// XML documentation above a line.
func (e *CSharpExtractor) extractXmlDocAbove(lines []string, lineNum int) string {
	if lineNum <= 1 || lineNum > len(lines) {
		return ""
	}

	var docLines []string
	inSummary := false

	for i := lineNum - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "///") {
			content := strings.TrimPrefix(line, "///")
			content = strings.TrimSpace(content)

			// Handle XML tags
			if strings.Contains(content, "</summary>") {
				inSummary = true
				content = strings.ReplaceAll(content, "</summary>", "")
			}
			if strings.Contains(content, "<summary>") {
				inSummary = false
				content = strings.ReplaceAll(content, "<summary>", "")
			}

			content = strings.TrimSpace(content)
			if content != "" && (inSummary || !strings.HasPrefix(content, "<")) {
				docLines = append([]string{content}, docLines...)
			}
		} else if strings.HasPrefix(line, "[") {
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

// ExtractFunctions extracts method definitions from C# source.
func (e *CSharpExtractor) ExtractFunctions(content string, filename string) []ExtractedFunction {
	var functions []ExtractedFunction
	lines := strings.Split(content, "\n")

	// Method pattern: public ReturnType MethodName(params) {
	methodRe := regexp.MustCompile(`(?:public|internal)\s+(?:static\s+)?(?:async\s+)?(?:virtual\s+)?(?:override\s+)?(\w+(?:<[^>]+>)?(?:\?)?)\s+(\w+)\s*\(([^)]*)\)\s*(?:where\s+[^{]+)?\s*\{`)

	for _, match := range methodRe.FindAllStringSubmatchIndex(content, -1) {
		returnType := content[match[2]:match[3]]
		methodName := content[match[4]:match[5]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Skip property-like methods
		if methodName == "get" || methodName == "set" {
			continue
		}

		paramsStr := ""
		if match[6] != -1 && match[7] != -1 {
			paramsStr = content[match[6]:match[7]]
		}

		fullMatch := content[match[0]:match[1]]
		isAsync := strings.Contains(fullMatch, " async ")
		body := ExtractFunctionBody(content, match[0])
		description := e.extractXmlDocAbove(lines, lineNum)
		logic := InferLogicFromBody(body)

		functions = append(functions, ExtractedFunction{
			Name:        methodName,
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

// parseParams parses C# method parameters.
func (e *CSharpExtractor) parseParams(paramsStr string) []ExtractedParam {
	var params []ExtractedParam

	if strings.TrimSpace(paramsStr) == "" {
		return params
	}

	parts := splitParams(paramsStr)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Remove attributes like [FromBody]
		part = regexp.MustCompile(`\[[^\]]+\]\s*`).ReplaceAllString(part, "")

		// Pattern: Type name or Type name = default
		paramRe := regexp.MustCompile(`(?:this\s+)?(?:params\s+)?(\w+(?:<[^>]+>)?(?:\?)?(?:\[\])?)\s+(\w+)(?:\s*=\s*(.+))?`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 2 {
			paramType := match[1]
			defaultVal := ""
			if len(match) > 3 && match[3] != "" {
				defaultVal = strings.TrimSpace(match[3])
			}

			params = append(params, ExtractedParam{
				Name:    match[2],
				Type:    e.MapTypeToSpec(paramType),
				Default: defaultVal,
			})
		}
	}

	return params
}

// ExtractTests extracts test cases from C# test files.
func (e *CSharpExtractor) ExtractTests(content string, filename string) []ExtractedTest {
	var tests []ExtractedTest

	// xUnit pattern: [Fact] or [Theory] public void TestName()
	testRe := regexp.MustCompile(`\[(?:Fact|Theory)[^\]]*\]\s*(?:public\s+)?(?:async\s+)?(?:Task\s+|void\s+)(\w+)\s*\(`)

	for _, match := range testRe.FindAllStringSubmatch(content, -1) {
		testFuncName := match[1]
		// Convert PascalCase to readable name
		testName := camelToSpaces(testFuncName)

		// Try to derive function being tested
		funcName := testFuncName
		funcName = strings.TrimPrefix(funcName, "Should")
		funcName = strings.TrimSuffix(funcName, "Test")

		tests = append(tests, ExtractedTest{
			Function:   funcName,
			Name:       testName,
			SourceFile: filename,
		})
	}

	// NUnit pattern: [Test] public void TestName()
	nunitTestRe := regexp.MustCompile(`\[Test(?:Case)?[^\]]*\]\s*(?:public\s+)?(?:async\s+)?(?:Task\s+|void\s+)(\w+)\s*\(`)
	for _, match := range nunitTestRe.FindAllStringSubmatch(content, -1) {
		testFuncName := match[1]
		testName := camelToSpaces(testFuncName)

		tests = append(tests, ExtractedTest{
			Function:   testFuncName,
			Name:       testName,
			SourceFile: filename,
		})
	}

	return tests
}

// MapTypeToSpec maps C# types to spec pseudo-types.
func (e *CSharpExtractor) MapTypeToSpec(csType string) string {
	csType = strings.TrimSpace(csType)

	// Handle nullable types
	if strings.HasSuffix(csType, "?") {
		innerType := strings.TrimSuffix(csType, "?")
		return "Optional " + e.MapTypeToSpec(innerType)
	}

	// Handle List types
	if strings.HasPrefix(csType, "List<") || strings.HasPrefix(csType, "IList<") ||
		strings.HasPrefix(csType, "IEnumerable<") || strings.HasPrefix(csType, "ICollection<") {
		innerType := regexp.MustCompile(`\w+<(.+)>`).FindStringSubmatch(csType)
		if len(innerType) > 1 {
			return "List of " + e.MapTypeToSpec(innerType[1])
		}
		return "List"
	}

	// Handle array types
	if strings.HasSuffix(csType, "[]") {
		innerType := strings.TrimSuffix(csType, "[]")
		return "List of " + e.MapTypeToSpec(innerType)
	}

	// Handle Dictionary types
	if strings.HasPrefix(csType, "Dictionary<") || strings.HasPrefix(csType, "IDictionary<") {
		return "Map"
	}

	// Handle Task types (async return)
	if strings.HasPrefix(csType, "Task<") {
		innerType := regexp.MustCompile(`Task<(.+)>`).FindStringSubmatch(csType)
		if len(innerType) > 1 {
			return e.MapTypeToSpec(innerType[1])
		}
	}
	if csType == "Task" {
		return "Nothing"
	}

	// Direct type mappings
	typeMap := map[string]string{
		"string":         "Text",
		"String":         "Text",
		"int":            "Integer",
		"Int32":          "Integer",
		"long":           "Integer",
		"Int64":          "Integer",
		"short":          "Integer",
		"Int16":          "Integer",
		"double":         "Float",
		"Double":         "Float",
		"float":          "Float",
		"Single":         "Float",
		"decimal":        "Float",
		"Decimal":        "Float",
		"bool":           "Boolean",
		"Boolean":        "Boolean",
		"void":           "Nothing",
		"Void":           "Nothing",
		"object":         "Any",
		"Object":         "Any",
		"DateTime":       "Timestamp",
		"DateTimeOffset": "Timestamp",
		"TimeSpan":       "Duration",
		"Guid":           "UUID",
		"":               "Nothing",
	}

	if mapped, ok := typeMap[csType]; ok {
		return mapped
	}

	return csType
}
