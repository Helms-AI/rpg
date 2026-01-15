package extractors

import (
	"regexp"
	"strings"
)

// JavaExtractor extracts information from Java source files.
type JavaExtractor struct{}

// NewJavaExtractor creates a new Java extractor.
func NewJavaExtractor() *JavaExtractor {
	return &JavaExtractor{}
}

// LanguageID returns the language identifier.
func (e *JavaExtractor) LanguageID() string {
	return "java"
}

// Extensions returns Java file extensions.
func (e *JavaExtractor) Extensions() []string {
	return []string{".java"}
}

// IsTestFile returns true if the file is a Java test file.
func (e *JavaExtractor) IsTestFile(filename string) bool {
	return strings.HasSuffix(filename, "Test.java") ||
		strings.HasSuffix(filename, "Tests.java") ||
		strings.Contains(filename, "/test/")
}

// ExtractPackageDescription extracts the class Javadoc.
func (e *JavaExtractor) ExtractPackageDescription(content string) string {
	// Look for class-level Javadoc
	javadocRe := regexp.MustCompile(`/\*\*\s*((?:[^*]|\*(?!/))*)\*/\s*(?:public\s+)?(?:final\s+)?(?:abstract\s+)?class`)
	if match := javadocRe.FindStringSubmatch(content); len(match) > 1 {
		return cleanJavadoc(match[1])
	}
	return ""
}

// cleanJavadoc removes Javadoc formatting.
func cleanJavadoc(doc string) string {
	// Remove leading * from each line
	lines := strings.Split(doc, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		// Skip @param, @return, @throws lines
		if strings.HasPrefix(line, "@") {
			continue
		}
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, " ")
}

// ExtractTypes extracts type definitions from Java source.
func (e *JavaExtractor) ExtractTypes(content string, filename string) []ExtractedType {
	var types []ExtractedType
	lines := strings.Split(content, "\n")

	// Class pattern: public class Name { ... }
	classRe := regexp.MustCompile(`(?:public\s+)?(?:final\s+)?(?:abstract\s+)?class\s+(\w+)(?:<[^>]+>)?(?:\s+extends\s+\w+)?(?:\s+implements\s+[\w,\s]+)?\s*\{`)
	for _, match := range classRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		body := ExtractFunctionBody(content, match[0])
		fields := e.parseClassFields(body)
		description := e.extractJavadocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Record pattern: public record Name(fields) { ... }
	recordRe := regexp.MustCompile(`(?:public\s+)?record\s+(\w+)\s*\(([^)]*)\)`)
	for _, match := range recordRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		fieldsStr := ""
		if match[4] != -1 && match[5] != -1 {
			fieldsStr = content[match[4]:match[5]]
		}
		fields := e.parseRecordFields(fieldsStr)
		description := e.extractJavadocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "struct",
			Description: description,
			Fields:      fields,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Interface pattern: public interface Name { ... }
	interfaceRe := regexp.MustCompile(`(?:public\s+)?interface\s+(\w+)(?:<[^>]+>)?(?:\s+extends\s+[\w,\s<>]+)?\s*\{`)
	for _, match := range interfaceRe.FindAllStringSubmatchIndex(content, -1) {
		name := content[match[2]:match[3]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])
		description := e.extractJavadocAbove(lines, lineNum)

		types = append(types, ExtractedType{
			Name:        name,
			Kind:        "interface",
			Description: description,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	// Enum pattern: public enum Name { ... }
	enumRe := regexp.MustCompile(`(?:public\s+)?enum\s+(\w+)\s*\{`)
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

// parseClassFields parses fields from a Java class body.
func (e *JavaExtractor) parseClassFields(body string) []ExtractedField {
	var fields []ExtractedField

	// Field pattern: private Type name;
	fieldRe := regexp.MustCompile(`(?:private|protected|public)?\s*(?:final\s+)?(\w+(?:<[^>]+>)?)\s+(\w+)\s*[;=]`)

	for _, match := range fieldRe.FindAllStringSubmatch(body, -1) {
		fieldType := match[1]
		fieldName := match[2]

		// Skip static fields and constants
		if strings.Contains(match[0], "static") {
			continue
		}

		fields = append(fields, ExtractedField{
			Name: fieldName,
			Type: e.MapTypeToSpec(fieldType),
		})
	}

	return fields
}

// parseRecordFields parses fields from a record definition.
func (e *JavaExtractor) parseRecordFields(fieldsStr string) []ExtractedField {
	var fields []ExtractedField

	parts := splitParams(fieldsStr)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Pattern: Type name
		paramRe := regexp.MustCompile(`(\w+(?:<[^>]+>)?)\s+(\w+)`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 2 {
			fields = append(fields, ExtractedField{
				Name: match[2],
				Type: e.MapTypeToSpec(match[1]),
			})
		}
	}

	return fields
}

// parseEnumVariants parses variants from an enum body.
func (e *JavaExtractor) parseEnumVariants(body string) []string {
	var variants []string

	// Enum variants are uppercase identifiers at the start
	// Stop at first method or semicolon
	semicolonIdx := strings.Index(body, ";")
	if semicolonIdx != -1 {
		body = body[:semicolonIdx]
	}

	variantRe := regexp.MustCompile(`([A-Z][A-Z0-9_]*)(?:\([^)]*\))?`)
	for _, match := range variantRe.FindAllStringSubmatch(body, -1) {
		variants = append(variants, match[1])
	}

	return variants
}

// extractJavadocAbove extracts Javadoc comment above a line.
func (e *JavaExtractor) extractJavadocAbove(lines []string, lineNum int) string {
	if lineNum <= 1 || lineNum > len(lines) {
		return ""
	}

	// Look for /** ... */ pattern above
	var javadocLines []string
	inJavadoc := false

	for i := lineNum - 2; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])

		if strings.HasSuffix(line, "*/") {
			inJavadoc = true
			line = strings.TrimSuffix(line, "*/")
		}

		if inJavadoc {
			if strings.HasPrefix(line, "/**") {
				line = strings.TrimPrefix(line, "/**")
				javadocLines = append([]string{line}, javadocLines...)
				break
			}
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "@") && line != "" {
				javadocLines = append([]string{line}, javadocLines...)
			}
		} else if line != "" && !strings.HasPrefix(line, "@") {
			break
		}
	}

	return strings.Join(javadocLines, " ")
}

// ExtractFunctions extracts method definitions from Java source.
func (e *JavaExtractor) ExtractFunctions(content string, filename string) []ExtractedFunction {
	var functions []ExtractedFunction
	lines := strings.Split(content, "\n")

	// Method pattern: public ReturnType methodName(params) {
	methodRe := regexp.MustCompile(`(?:public|protected)\s+(?:static\s+)?(?:final\s+)?(?:synchronized\s+)?(\w+(?:<[^>]+>)?)\s+(\w+)\s*\(([^)]*)\)\s*(?:throws\s+[\w,\s]+)?\s*\{`)

	for _, match := range methodRe.FindAllStringSubmatchIndex(content, -1) {
		returnType := content[match[2]:match[3]]
		methodName := content[match[4]:match[5]]
		lineNum := FindLineNumber(content, content[match[0]:match[1]])

		// Skip getters/setters and constructor-like methods
		if strings.HasPrefix(methodName, "get") || strings.HasPrefix(methodName, "set") ||
			strings.HasPrefix(methodName, "is") {
			continue
		}

		paramsStr := ""
		if match[6] != -1 && match[7] != -1 {
			paramsStr = content[match[6]:match[7]]
		}

		body := ExtractFunctionBody(content, match[0])
		description := e.extractJavadocAbove(lines, lineNum)
		logic := InferLogicFromBody(body)

		functions = append(functions, ExtractedFunction{
			Name:        methodName,
			Description: description,
			Parameters:  e.parseParams(paramsStr),
			Returns:     e.MapTypeToSpec(returnType),
			Logic:       logic,
			Body:        body,
			SourceFile:  filename,
			LineNumber:  lineNum,
		})
	}

	return functions
}

// parseParams parses Java method parameters.
func (e *JavaExtractor) parseParams(paramsStr string) []ExtractedParam {
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

		// Remove annotations like @NonNull
		part = regexp.MustCompile(`@\w+\s*`).ReplaceAllString(part, "")

		// Pattern: Type name or final Type name
		paramRe := regexp.MustCompile(`(?:final\s+)?(\w+(?:<[^>]+>)?)\s+(\w+)`)
		if match := paramRe.FindStringSubmatch(part); len(match) > 2 {
			params = append(params, ExtractedParam{
				Name: match[2],
				Type: e.MapTypeToSpec(match[1]),
			})
		}
	}

	return params
}

// ExtractTests extracts test cases from Java test files.
func (e *JavaExtractor) ExtractTests(content string, filename string) []ExtractedTest {
	var tests []ExtractedTest

	// JUnit pattern: @Test void testName() or @Test public void testName()
	testRe := regexp.MustCompile(`@Test\s*(?:public\s+)?(?:void|[\w<>]+)\s+(\w+)\s*\(`)

	for _, match := range testRe.FindAllStringSubmatch(content, -1) {
		testFuncName := match[1]
		// Convert camelCase to readable name
		testName := camelToSpaces(testFuncName)

		// Try to determine function being tested
		funcName := strings.TrimPrefix(testFuncName, "test")
		funcName = strings.TrimPrefix(funcName, "should")

		tests = append(tests, ExtractedTest{
			Function:   funcName,
			Name:       testName,
			SourceFile: filename,
		})
	}

	return tests
}

// camelToSpaces converts camelCase to spaces.
func camelToSpaces(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// MapTypeToSpec maps Java types to spec pseudo-types.
func (e *JavaExtractor) MapTypeToSpec(javaType string) string {
	javaType = strings.TrimSpace(javaType)

	// Handle generic List types
	if strings.HasPrefix(javaType, "List<") || strings.HasPrefix(javaType, "ArrayList<") {
		innerType := regexp.MustCompile(`\w+<(.+)>`).FindStringSubmatch(javaType)
		if len(innerType) > 1 {
			return "List of " + e.MapTypeToSpec(innerType[1])
		}
		return "List"
	}

	// Handle Optional types
	if strings.HasPrefix(javaType, "Optional<") {
		innerType := regexp.MustCompile(`Optional<(.+)>`).FindStringSubmatch(javaType)
		if len(innerType) > 1 {
			return "Optional " + e.MapTypeToSpec(innerType[1])
		}
	}

	// Handle Map types
	if strings.HasPrefix(javaType, "Map<") || strings.HasPrefix(javaType, "HashMap<") {
		return "Map"
	}

	// Direct type mappings
	typeMap := map[string]string{
		"String":     "Text",
		"int":        "Integer",
		"Integer":    "Integer",
		"long":       "Integer",
		"Long":       "Integer",
		"double":     "Float",
		"Double":     "Float",
		"float":      "Float",
		"Float":      "Float",
		"boolean":    "Boolean",
		"Boolean":    "Boolean",
		"void":       "Nothing",
		"Void":       "Nothing",
		"Object":     "Any",
		"Instant":    "Timestamp",
		"LocalDate":  "Timestamp",
		"LocalDateTime": "Timestamp",
		"UUID":       "UUID",
		"":           "Nothing",
	}

	if mapped, ok := typeMap[javaType]; ok {
		return mapped
	}

	return javaType
}
