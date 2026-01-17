// Package specparser provides utilities for parsing specification files.
package specparser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Parser parses markdown spec files into structured SpecAnalysis.
type Parser struct{}

// NewParser creates a new spec parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a spec file from the given path.
func (p *Parser) ParseFile(specPath string) (*SpecAnalysis, error) {
	content, err := os.ReadFile(specPath)
	if err != nil {
		return nil, err
	}
	return p.Parse(string(content), filepath.Base(specPath))
}

// Parse parses spec content into a SpecAnalysis.
func (p *Parser) Parse(content, name string) (*SpecAnalysis, error) {
	analysis := &SpecAnalysis{
		Name:          extractSpecName(content, name),
		Overview:      extractOverview(content),
		Types:         []SpecType{},
		Functions:     []SpecFunction{},
		Tests:         []SpecTest{},
		Dependencies:  []SpecDependency{},
		Configuration: []SpecConfig{},
	}

	// Extract sections based on markdown headers
	sections := splitIntoSections(content)

	for sectionName, sectionContent := range sections {
		sectionLower := strings.ToLower(sectionName)

		switch {
		case strings.Contains(sectionLower, "type") || strings.Contains(sectionLower, "data") || strings.Contains(sectionLower, "model"):
			analysis.Types = append(analysis.Types, parseTypes(sectionContent)...)
		case strings.Contains(sectionLower, "function") || strings.Contains(sectionLower, "api") || strings.Contains(sectionLower, "method"):
			analysis.Functions = append(analysis.Functions, parseFunctions(sectionContent)...)
		case strings.Contains(sectionLower, "test"):
			analysis.Tests = append(analysis.Tests, parseTests(sectionContent)...)
		case strings.Contains(sectionLower, "depend"):
			analysis.Dependencies = append(analysis.Dependencies, parseDependencies(sectionContent)...)
		case strings.Contains(sectionLower, "config") || strings.Contains(sectionLower, "environment"):
			analysis.Configuration = append(analysis.Configuration, parseConfiguration(sectionContent)...)
		}
	}

	analysis.CalculateTotals()
	return analysis, nil
}

// extractSpecName extracts the spec name from the content or falls back to filename.
func extractSpecName(content, filename string) string {
	// Try to get name from H1 header
	h1Pattern := regexp.MustCompile(`(?m)^#\s+(.+)$`)
	if matches := h1Pattern.FindStringSubmatch(content); len(matches) > 1 {
		name := strings.TrimSpace(matches[1])
		// Remove .spec.md or .md suffix if present
		name = strings.TrimSuffix(name, ".spec.md")
		name = strings.TrimSuffix(name, ".md")
		return name
	}

	// Fall back to filename
	name := strings.TrimSuffix(filename, ".spec.md")
	name = strings.TrimSuffix(name, ".md")
	return name
}

// extractOverview extracts the overview/description from the content.
func extractOverview(content string) string {
	// Look for Overview/Description section
	overviewPattern := regexp.MustCompile(`(?mi)^##\s+(overview|description|about)\s*$\n([\s\S]*?)(?:^##|\z)`)
	if matches := overviewPattern.FindStringSubmatch(content); len(matches) > 2 {
		return strings.TrimSpace(matches[2])
	}

	// Fall back to first paragraph after H1
	lines := strings.Split(content, "\n")
	inOverview := false
	var overviewLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			inOverview = true
			continue
		}
		if inOverview {
			if strings.HasPrefix(line, "##") {
				break
			}
			if strings.TrimSpace(line) != "" {
				overviewLines = append(overviewLines, line)
			}
			if len(overviewLines) > 0 && strings.TrimSpace(line) == "" {
				break // Stop at first empty line after content
			}
		}
	}

	return strings.TrimSpace(strings.Join(overviewLines, " "))
}

// splitIntoSections splits markdown content into sections by H2 headers.
func splitIntoSections(content string) map[string]string {
	sections := make(map[string]string)

	h2Pattern := regexp.MustCompile(`(?m)^##\s+(.+)$`)
	matches := h2Pattern.FindAllStringSubmatchIndex(content, -1)

	for i, match := range matches {
		_ = match[0] // headerStart - unused but kept for clarity
		headerEnd := match[1]
		nameStart := match[2]
		nameEnd := match[3]

		sectionName := content[nameStart:nameEnd]

		// Find end of this section
		var sectionEnd int
		if i+1 < len(matches) {
			sectionEnd = matches[i+1][0]
		} else {
			sectionEnd = len(content)
		}

		sectionContent := content[headerEnd:sectionEnd]
		sections[sectionName] = sectionContent
	}

	return sections
}

// parseTypes extracts type definitions from a section.
func parseTypes(content string) []SpecType {
	var types []SpecType

	// Pattern for H3/H4 type definitions
	typePattern := regexp.MustCompile(`(?mi)^###\s+(.+?)\s*(?:\((struct|interface|enum|class|type)\))?\s*$`)
	typeMatches := typePattern.FindAllStringSubmatchIndex(content, -1)

	for i, match := range typeMatches {
		nameStart := match[2]
		nameEnd := match[3]
		name := content[nameStart:nameEnd]

		kind := "struct"
		if match[4] != -1 && match[5] != -1 {
			kind = strings.ToLower(content[match[4]:match[5]])
		}

		// Find section content
		sectionStart := match[1]
		var sectionEnd int
		if i+1 < len(typeMatches) {
			sectionEnd = typeMatches[i+1][0]
		} else {
			sectionEnd = len(content)
		}
		sectionContent := content[sectionStart:sectionEnd]

		specType := SpecType{
			Name:        strings.TrimSpace(name),
			Kind:        kind,
			Description: extractDescription(sectionContent),
			Fields:      parseFields(sectionContent),
			Methods:     parseMethods(sectionContent),
			IsPublic:    isPublic(name),
		}

		types = append(types, specType)
	}

	// Also parse tables for type definitions
	tableTypes := parseTypesFromTables(content)
	types = append(types, tableTypes...)

	return types
}

// parseFields extracts fields from a type section.
func parseFields(content string) []SpecField {
	var fields []SpecField

	// Look for markdown table with fields
	tablePattern := regexp.MustCompile(`(?m)\|(.+)\|(.+)\|(.+)\|`)
	tableMatches := tablePattern.FindAllStringSubmatch(content, -1)

	headerFound := false
	for _, match := range tableMatches {
		if len(match) < 4 {
			continue
		}

		col1 := strings.TrimSpace(match[1])
		col2 := strings.TrimSpace(match[2])
		col3 := strings.TrimSpace(match[3])

		// Skip header and separator rows
		if strings.Contains(col1, "---") || strings.Contains(col2, "---") {
			headerFound = true
			continue
		}
		if strings.ToLower(col1) == "field" || strings.ToLower(col1) == "name" {
			headerFound = true
			continue
		}

		if !headerFound {
			continue
		}

		field := SpecField{
			Name:        col1,
			Type:        col2,
			Description: col3,
			Required:    !strings.Contains(strings.ToLower(col2), "optional"),
		}
		fields = append(fields, field)
	}

	// Also look for bullet list fields
	bulletPattern := regexp.MustCompile(`(?m)^[-*]\s+\x60?(\w+)\x60?\s*[:\-]\s*\x60?([^\x60\s]+)\x60?\s*[-:]?\s*(.*)$`)
	bulletMatches := bulletPattern.FindAllStringSubmatch(content, -1)

	for _, match := range bulletMatches {
		if len(match) < 4 {
			continue
		}

		field := SpecField{
			Name:        match[1],
			Type:        match[2],
			Description: strings.TrimSpace(match[3]),
			Required:    true,
		}
		fields = append(fields, field)
	}

	return fields
}

// parseMethods extracts method signatures from an interface section.
func parseMethods(content string) []string {
	var methods []string

	// Look for method signatures in code blocks or bullet lists
	methodPattern := regexp.MustCompile(`(?m)^[-*]\s+\x60([A-Z]\w+\([^)]*\)[^)\x60]*)\x60`)
	matches := methodPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			methods = append(methods, match[1])
		}
	}

	return methods
}

// parseTypesFromTables extracts type definitions from markdown tables.
func parseTypesFromTables(content string) []SpecType {
	var types []SpecType

	// Pattern for type tables
	tableStartPattern := regexp.MustCompile(`(?mi)^\|\s*(type|name)\s*\|`)
	if !tableStartPattern.MatchString(content) {
		return types
	}

	lines := strings.Split(content, "\n")
	inTable := false
	headerCols := []string{}

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
			cols := parseTableRow(line)

			if !inTable {
				// Check if this is a header row
				headerLower := strings.ToLower(strings.Join(cols, " "))
				if strings.Contains(headerLower, "type") || strings.Contains(headerLower, "name") {
					inTable = true
					headerCols = cols
					continue
				}
			}

			// Skip separator row
			if len(cols) > 0 && strings.Contains(cols[0], "-") {
				continue
			}

			if inTable && len(cols) >= 2 {
				specType := SpecType{
					Name:     cols[0],
					Kind:     "struct",
					IsPublic: isPublic(cols[0]),
				}

				// Try to extract kind from columns
				for j, col := range cols {
					if j < len(headerCols) {
						headerLower := strings.ToLower(headerCols[j])
						if strings.Contains(headerLower, "kind") || strings.Contains(headerLower, "type") {
							if col != specType.Name {
								specType.Kind = strings.ToLower(col)
							}
						} else if strings.Contains(headerLower, "description") {
							specType.Description = col
						}
					}
				}

				if specType.Name != "" && !strings.Contains(specType.Name, "-") {
					types = append(types, specType)
				}
			}
		} else if inTable && line == "" {
			inTable = false
			headerCols = []string{}
		}
	}

	return types
}

// parseTableRow parses a markdown table row into columns.
func parseTableRow(row string) []string {
	// Remove leading/trailing pipes
	row = strings.TrimPrefix(row, "|")
	row = strings.TrimSuffix(row, "|")

	cols := strings.Split(row, "|")
	for i := range cols {
		cols[i] = strings.TrimSpace(cols[i])
	}
	return cols
}

// parseFunctions extracts function definitions from a section.
func parseFunctions(content string) []SpecFunction {
	var functions []SpecFunction

	// Pattern for H3/H4 function definitions
	funcPattern := regexp.MustCompile("(?mi)^###\\s+`?([A-Za-z_][A-Za-z0-9_]*)`?\\s*(?:\\(([^)]*)?\\))?")
	funcMatches := funcPattern.FindAllStringSubmatchIndex(content, -1)

	for i, match := range funcMatches {
		nameStart := match[2]
		nameEnd := match[3]
		name := content[nameStart:nameEnd]

		// Find section content
		sectionStart := match[1]
		var sectionEnd int
		if i+1 < len(funcMatches) {
			sectionEnd = funcMatches[i+1][0]
		} else {
			sectionEnd = len(content)
		}
		sectionContent := content[sectionStart:sectionEnd]

		specFunc := SpecFunction{
			Name:        strings.TrimSpace(name),
			Description: extractDescription(sectionContent),
			Parameters:  parseParameters(sectionContent),
			Returns:     parseReturns(sectionContent),
			Logic:       extractLogic(sectionContent),
			Errors:      parseErrors(sectionContent),
			IsPublic:    isPublic(name),
		}

		// Check for async keyword
		if strings.Contains(strings.ToLower(sectionContent), "async") {
			specFunc.IsAsync = true
		}

		functions = append(functions, specFunc)
	}

	return functions
}

// parseParameters extracts function parameters.
func parseParameters(content string) []SpecParameter {
	var params []SpecParameter

	// Look for Parameters/Arguments section
	paramSection := regexp.MustCompile(`(?mi)(?:\*\*parameters?\*\*|parameters?:)\s*\n([\s\S]*?)(?:\n\*\*|\n###|\z)`)
	if matches := paramSection.FindStringSubmatch(content); len(matches) > 1 {
		content = matches[1]
	}

	// Parse bullet list parameters
	paramPattern := regexp.MustCompile(`(?m)^[-*]\s+\x60?(\w+)\x60?\s*[:\(]\s*\x60?([^\x60\)\s,]+)\x60?[\),]?\s*[-:]?\s*(.*)$`)
	paramMatches := paramPattern.FindAllStringSubmatch(content, -1)

	for _, match := range paramMatches {
		if len(match) < 4 {
			continue
		}

		param := SpecParameter{
			Name:        match[1],
			Type:        match[2],
			Description: strings.TrimSpace(match[3]),
			Required:    !strings.Contains(strings.ToLower(match[3]), "optional"),
		}
		params = append(params, param)
	}

	// Also check tables
	tableParams := parseParametersFromTable(content)
	params = append(params, tableParams...)

	return params
}

// parseParametersFromTable extracts parameters from a markdown table.
func parseParametersFromTable(content string) []SpecParameter {
	var params []SpecParameter

	tablePattern := regexp.MustCompile(`(?m)\|(.+)\|(.+)\|(.+)\|`)
	matches := tablePattern.FindAllStringSubmatch(content, -1)

	headerFound := false
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		col1 := strings.TrimSpace(match[1])
		col2 := strings.TrimSpace(match[2])
		col3 := strings.TrimSpace(match[3])

		// Skip separator
		if strings.Contains(col1, "-") {
			headerFound = true
			continue
		}

		// Skip header
		if strings.ToLower(col1) == "name" || strings.ToLower(col1) == "parameter" {
			headerFound = true
			continue
		}

		if headerFound {
			param := SpecParameter{
				Name:        col1,
				Type:        col2,
				Description: col3,
				Required:    true,
			}
			params = append(params, param)
		}
	}

	return params
}

// parseReturns extracts return type information.
func parseReturns(content string) []SpecReturn {
	var returns []SpecReturn

	// Look for Returns section
	returnPattern := regexp.MustCompile(`(?mi)(?:\*\*returns?\*\*|returns?:)\s*\x60?([^\x60\n]+)\x60?`)
	if matches := returnPattern.FindStringSubmatch(content); len(matches) > 1 {
		ret := SpecReturn{
			Type:        strings.TrimSpace(matches[1]),
			Description: "",
		}
		returns = append(returns, ret)
	}

	return returns
}

// extractLogic extracts the logic/implementation description.
func extractLogic(content string) string {
	// Look for Logic section
	logicPattern := regexp.MustCompile(`(?mi)(?:\*\*logic\*\*|logic:)\s*\n([\s\S]*?)(?:\n\*\*|\n###|\z)`)
	if matches := logicPattern.FindStringSubmatch(content); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

// parseErrors extracts error conditions.
func parseErrors(content string) []SpecError {
	var errors []SpecError

	// Look for Errors section
	errorPattern := regexp.MustCompile(`(?mi)(?:\*\*errors?\*\*|errors?:)\s*\n([\s\S]*?)(?:\n\*\*|\n###|\z)`)
	if matches := errorPattern.FindStringSubmatch(content); len(matches) > 1 {
		errorContent := matches[1]

		// Parse bullet list errors
		bulletPattern := regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
		bullets := bulletPattern.FindAllStringSubmatch(errorContent, -1)

		for _, match := range bullets {
			if len(match) > 1 {
				err := SpecError{
					Condition: strings.TrimSpace(match[1]),
				}
				errors = append(errors, err)
			}
		}
	}

	return errors
}

// parseTests extracts test case definitions.
func parseTests(content string) []SpecTest {
	var tests []SpecTest

	// Pattern for test case headers
	testPattern := regexp.MustCompile(`(?mi)^###\s+(.+)$`)
	testMatches := testPattern.FindAllStringSubmatchIndex(content, -1)

	for i, match := range testMatches {
		nameStart := match[2]
		nameEnd := match[3]
		name := content[nameStart:nameEnd]

		// Find section content
		sectionStart := match[1]
		var sectionEnd int
		if i+1 < len(testMatches) {
			sectionEnd = testMatches[i+1][0]
		} else {
			sectionEnd = len(content)
		}
		sectionContent := content[sectionStart:sectionEnd]

		test := SpecTest{
			Name:        strings.TrimSpace(name),
			Description: extractDescription(sectionContent),
			Given:       parseGiven(sectionContent),
			When:        extractWhen(sectionContent),
			Then:        parseThen(sectionContent),
		}

		tests = append(tests, test)
	}

	return tests
}

// parseGiven extracts test preconditions.
func parseGiven(content string) []SpecCondition {
	var conditions []SpecCondition

	givenPattern := regexp.MustCompile(`(?mi)(?:\*\*given\*\*|given:)\s*\n([\s\S]*?)(?:\n\*\*|\n###|\z)`)
	if matches := givenPattern.FindStringSubmatch(content); len(matches) > 1 {
		bulletPattern := regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
		bullets := bulletPattern.FindAllStringSubmatch(matches[1], -1)

		for _, match := range bullets {
			if len(match) > 1 {
				conditions = append(conditions, SpecCondition{
					Description: strings.TrimSpace(match[1]),
				})
			}
		}
	}

	return conditions
}

// extractWhen extracts the test action.
func extractWhen(content string) string {
	whenPattern := regexp.MustCompile(`(?mi)(?:\*\*when\*\*|when:)\s*(.+)$`)
	if matches := whenPattern.FindStringSubmatch(content); len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

// parseThen extracts test assertions.
func parseThen(content string) []SpecAssertion {
	var assertions []SpecAssertion

	thenPattern := regexp.MustCompile(`(?mi)(?:\*\*then\*\*|(?:expect|then):)\s*\n([\s\S]*?)(?:\n\*\*|\n###|\z)`)
	if matches := thenPattern.FindStringSubmatch(content); len(matches) > 1 {
		bulletPattern := regexp.MustCompile(`(?m)^[-*]\s+(.+)$`)
		bullets := bulletPattern.FindAllStringSubmatch(matches[1], -1)

		for _, match := range bullets {
			if len(match) > 1 {
				assertions = append(assertions, SpecAssertion{
					Description: strings.TrimSpace(match[1]),
				})
			}
		}
	}

	return assertions
}

// parseDependencies extracts dependency definitions.
func parseDependencies(content string) []SpecDependency {
	var deps []SpecDependency

	// Parse bullet list dependencies
	depPattern := regexp.MustCompile(`(?m)^[-*]\s+\x60?([^\x60\s]+)\x60?\s*[-:]?\s*(.*)$`)
	matches := depPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		dep := SpecDependency{
			Name:        match[1],
			Description: strings.TrimSpace(match[2]),
		}
		deps = append(deps, dep)
	}

	return deps
}

// parseConfiguration extracts configuration items.
func parseConfiguration(content string) []SpecConfig {
	var configs []SpecConfig

	// Parse table configs
	tablePattern := regexp.MustCompile(`(?m)\|(.+)\|(.+)\|(.+)\|`)
	matches := tablePattern.FindAllStringSubmatch(content, -1)

	headerFound := false
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		col1 := strings.TrimSpace(match[1])
		col2 := strings.TrimSpace(match[2])
		col3 := strings.TrimSpace(match[3])

		if strings.Contains(col1, "-") {
			headerFound = true
			continue
		}
		if strings.ToLower(col1) == "name" || strings.ToLower(col1) == "variable" {
			headerFound = true
			continue
		}

		if headerFound {
			config := SpecConfig{
				Name:        col1,
				Type:        col2,
				Description: col3,
			}
			configs = append(configs, config)
		}
	}

	// Also parse bullet list configs
	bulletPattern := regexp.MustCompile(`(?m)^[-*]\s+\x60?([A-Z_]+)\x60?\s*[:\-]\s*(.*)$`)
	bulletMatches := bulletPattern.FindAllStringSubmatch(content, -1)

	for _, match := range bulletMatches {
		if len(match) < 3 {
			continue
		}

		config := SpecConfig{
			Name:        match[1],
			Description: strings.TrimSpace(match[2]),
			Type:        "string",
		}
		configs = append(configs, config)
	}

	return configs
}

// extractDescription extracts a description from section content.
func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	var descLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines, headers, tables, code blocks
		if line == "" {
			if len(descLines) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "|") ||
			strings.HasPrefix(line, "```") || strings.HasPrefix(line, "**") {
			continue
		}
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			break
		}
		descLines = append(descLines, line)
	}

	return strings.Join(descLines, " ")
}

// isPublic determines if a name represents a public identifier.
func isPublic(name string) bool {
	if len(name) == 0 {
		return false
	}
	// In Go and many languages, exported names start with uppercase
	firstChar := rune(name[0])
	return firstChar >= 'A' && firstChar <= 'Z'
}
