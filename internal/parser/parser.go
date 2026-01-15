// Package parser provides markdown spec file parsing functionality.
package parser

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/kon1790/rpg/pkg/spec"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Parse parses a markdown spec file and returns a structured Spec.
func Parse(content string) (*spec.Spec, error) {
	source := []byte(content)
	md := goldmark.New()
	reader := text.NewReader(source)
	doc := md.Parser().Parse(reader)

	p := &specParser{
		source:   source,
		spec:     &spec.Spec{},
		sections: make(map[string]string),
	}

	// Extract sections from the document
	p.extractSections(doc)

	// Parse each section
	p.parseTitle()
	p.parseMeta()
	p.parseTargetLanguages()
	// Fallback: if AST-based extraction failed, parse directly from source
	if len(p.spec.TargetLanguages) == 0 {
		p.parseTargetLanguagesFromSource()
	}
	p.parseDependencies()
	p.parseTypes()
	p.parseFunctions()
	p.parseTests()

	return p.spec, nil
}

// Validate validates a spec without full parsing.
func Validate(content string) spec.ValidationResult {
	result := spec.ValidationResult{Valid: true}
	var errors []spec.ValidationError

	// Check for required sections
	if !strings.Contains(content, "# ") {
		errors = append(errors, spec.ValidationError{
			Severity: "error",
			Code:     "MISSING_TITLE",
			Message:  "Spec must have a title (# Title)",
		})
	}

	if !strings.Contains(content, "## Target Languages") && !strings.Contains(content, "## Target languages") {
		errors = append(errors, spec.ValidationError{
			Severity: "error",
			Code:     "MISSING_TARGET_LANGUAGES",
			Message:  "Spec must have a '## Target Languages' section",
		})
	}

	if !strings.Contains(content, "## Functions") && !strings.Contains(content, "## functions") {
		errors = append(errors, spec.ValidationError{
			Severity: "error",
			Code:     "MISSING_FUNCTIONS",
			Message:  "Spec must have a '## Functions' section",
		})
	}

	if len(errors) > 0 {
		result.Valid = false
		result.Errors = errors
	}

	return result
}

// specParser is an internal parser state.
type specParser struct {
	source   []byte
	spec     *spec.Spec
	sections map[string]string
}

// extractSections walks the AST and extracts content by section.
func (p *specParser) extractSections(doc ast.Node) {
	var currentSection string
	var currentContent bytes.Buffer

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if heading, ok := n.(*ast.Heading); ok {
			// Save previous section
			if currentSection != "" {
				p.sections[currentSection] = currentContent.String()
			}

			// Start new section
			currentSection = strings.ToLower(string(heading.Text(p.source)))
			currentContent.Reset()
			return ast.WalkSkipChildren, nil
		}

		// Collect text content
		if textNode, ok := n.(*ast.Text); ok {
			currentContent.Write(textNode.Segment.Value(p.source))
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return
	}

	// Save last section
	if currentSection != "" {
		p.sections[currentSection] = currentContent.String()
	}
}

// parseTitle extracts the spec title from the first H1.
func (p *specParser) parseTitle() {
	// Find title in source (first # heading)
	lines := strings.Split(string(p.source), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			p.spec.Name = strings.TrimPrefix(line, "# ")
			p.spec.Name = strings.TrimSpace(p.spec.Name)
			break
		}
	}
}

// parseMeta extracts metadata from the Meta section.
func (p *specParser) parseMeta() {
	section := p.getSection("meta")
	if section == "" {
		return
	}

	// Parse key-value pairs
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			line = strings.TrimPrefix(line, "- ")
		}
		if strings.HasPrefix(line, "**") {
			// Parse **key**: value format
			re := regexp.MustCompile(`\*\*(\w+)\*\*:\s*(.+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) == 3 {
				key := strings.ToLower(matches[1])
				value := strings.TrimSpace(matches[2])
				switch key {
				case "version":
					p.spec.Version = value
				case "author":
					p.spec.Author = value
				case "license":
					p.spec.License = value
				case "description":
					p.spec.Description = value
				}
			}
		} else if strings.Contains(line, ":") {
			// Parse key: value format
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.ToLower(strings.TrimSpace(parts[0]))
				value := strings.TrimSpace(parts[1])
				switch key {
				case "version":
					p.spec.Version = value
				case "author":
					p.spec.Author = value
				case "license":
					p.spec.License = value
				case "description":
					p.spec.Description = value
				}
			}
		}
	}
}

// parseTargetLanguages extracts the target languages list.
func (p *specParser) parseTargetLanguages() {
	section := p.getSection("target languages")
	if section == "" {
		return
	}

	// Parse list items from the section
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			lang := strings.TrimPrefix(line, "- ")
			lang = strings.TrimSpace(lang)
			// Remove any extra info after the language name
			if idx := strings.Index(lang, "\n"); idx != -1 {
				lang = lang[:idx]
			}
			if lang != "" {
				p.spec.TargetLanguages = append(p.spec.TargetLanguages, lang)
			}
		}
	}
}

// parseTargetLanguagesFromSource parses target languages directly from source as fallback.
func (p *specParser) parseTargetLanguagesFromSource() {
	source := string(p.source)
	lowerSource := strings.ToLower(source)

	// Find Target Languages section
	idx := strings.Index(lowerSource, "## target languages")
	if idx == -1 {
		return
	}

	// Find end of section
	endIdx := len(source)
	nextSection := strings.Index(lowerSource[idx+19:], "\n## ")
	if nextSection != -1 {
		endIdx = idx + 19 + nextSection
	}

	section := source[idx+19 : endIdx]
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			lang := strings.TrimPrefix(line, "- ")
			lang = strings.TrimSpace(lang)
			if lang != "" {
				p.spec.TargetLanguages = append(p.spec.TargetLanguages, lang)
			}
		}
	}
}

// parseDependencies extracts the dependencies list.
func (p *specParser) parseDependencies() {
	section := p.getSection("dependencies")
	if section == "" {
		return
	}

	lines := strings.Split(section, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			dep := strings.TrimPrefix(line, "- ")
			dep = strings.TrimSpace(dep)
			if dep != "" {
				p.spec.Dependencies = append(p.spec.Dependencies, dep)
			}
		}
	}
}

// parseTypes extracts type definitions.
func (p *specParser) parseTypes() {
	// Find the Types section in raw source and parse it
	source := string(p.source)
	typesStart := strings.Index(source, "## Types")
	if typesStart == -1 {
		return
	}

	// Find the end of the Types section (next ## heading or EOF)
	typesEnd := len(source)
	nextSection := strings.Index(source[typesStart+8:], "\n## ")
	if nextSection != -1 {
		typesEnd = typesStart + 8 + nextSection
	}

	typesContent := source[typesStart:typesEnd]

	// Parse each ### heading as a type
	typeBlocks := strings.Split(typesContent, "\n### ")
	for i, block := range typeBlocks {
		if i == 0 {
			continue // Skip the section header
		}

		lines := strings.Split(block, "\n")
		if len(lines) == 0 {
			continue
		}

		typeDef := spec.TypeDef{
			Name: strings.TrimSpace(lines[0]),
		}

		// Determine type kind and parse fields
		blockContent := strings.Join(lines[1:], "\n")
		if strings.Contains(blockContent, "is one of:") {
			typeDef.Kind = "enum"
			// Parse variants
			inVariants := false
			for _, line := range lines[1:] {
				line = strings.TrimSpace(line)
				if line == "is one of:" {
					inVariants = true
					continue
				}
				if inVariants && strings.HasPrefix(line, "- ") {
					variant := strings.TrimPrefix(line, "- ")
					typeDef.Variants = append(typeDef.Variants, spec.Field{
						Name: strings.TrimSpace(variant),
					})
				}
			}
		} else if strings.Contains(blockContent, "contains:") {
			typeDef.Kind = "struct"
			// Parse fields
			inFields := false
			for _, line := range lines[1:] {
				line = strings.TrimSpace(line)
				if line == "contains:" {
					inFields = true
					continue
				}
				if inFields && strings.HasPrefix(line, "- ") {
					fieldLine := strings.TrimPrefix(line, "- ")
					field := parseField(fieldLine)
					typeDef.Fields = append(typeDef.Fields, field)
				}
			}
		}

		if typeDef.Name != "" {
			p.spec.Types = append(p.spec.Types, typeDef)
		}
	}
}

// parseFunctions extracts function definitions.
func (p *specParser) parseFunctions() {
	// Find the Functions section in raw source
	source := string(p.source)
	funcStart := strings.Index(source, "## Functions")
	if funcStart == -1 {
		return
	}

	// Find the end of the Functions section
	funcEnd := len(source)
	nextSection := strings.Index(source[funcStart+12:], "\n## ")
	if nextSection != -1 {
		funcEnd = funcStart + 12 + nextSection
	}

	funcContent := source[funcStart:funcEnd]

	// Parse each ### heading as a function
	funcBlocks := strings.Split(funcContent, "\n### ")
	for i, block := range funcBlocks {
		if i == 0 {
			continue // Skip the section header
		}

		fn := parseFunction(block)
		if fn.Name != "" {
			p.spec.Functions = append(p.spec.Functions, fn)
		}
	}
}

// parseFunction parses a single function block.
func parseFunction(block string) spec.Function {
	lines := strings.Split(block, "\n")
	if len(lines) == 0 {
		return spec.Function{}
	}

	fn := spec.Function{}

	// Parse name and modifiers from first line
	nameLine := strings.TrimSpace(lines[0])
	if strings.Contains(nameLine, "[async]") {
		fn.Async = true
		nameLine = strings.Replace(nameLine, "[async]", "", 1)
	}
	if strings.Contains(nameLine, "[pure]") {
		fn.Pure = true
		nameLine = strings.Replace(nameLine, "[pure]", "", 1)
	}
	fn.Name = strings.TrimSpace(nameLine)

	// Parse the rest of the block
	var currentSection string
	var logicLines []string
	inCodeBlock := false

	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)

		// Track code blocks
		if strings.HasPrefix(trimmed, "```") {
			if inCodeBlock {
				inCodeBlock = false
				fn.Logic = strings.Join(logicLines, "\n")
				logicLines = nil
			} else {
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			logicLines = append(logicLines, line)
			continue
		}

		// Parse sections
		if strings.HasPrefix(trimmed, "**accepts:**") {
			currentSection = "accepts"
			continue
		}
		if strings.HasPrefix(trimmed, "**returns:**") {
			currentSection = "returns"
			returnType := strings.TrimPrefix(trimmed, "**returns:**")
			returnType = strings.TrimSpace(returnType)
			if returnType != "" {
				fn.Returns = &spec.Return{Type: returnType}
			}
			continue
		}
		if strings.HasPrefix(trimmed, "**logic:**") {
			currentSection = "logic"
			continue
		}
		if strings.HasPrefix(trimmed, "**errors:**") {
			currentSection = "errors"
			continue
		}

		// Parse description (text before any section)
		if currentSection == "" && trimmed != "" && !strings.HasPrefix(trimmed, "- ") {
			if fn.Description != "" {
				fn.Description += " "
			}
			fn.Description += trimmed
			continue
		}

		// Parse section content
		if strings.HasPrefix(trimmed, "- ") {
			content := strings.TrimPrefix(trimmed, "- ")
			switch currentSection {
			case "accepts":
				param := parseParam(content)
				fn.Accepts = append(fn.Accepts, param)
			case "errors":
				fn.Errors = append(fn.Errors, content)
			}
		}
	}

	return fn
}

// parseParam parses a parameter definition.
func parseParam(line string) spec.Param {
	param := spec.Param{}

	// Format: name: Type (description) (defaults to "value")
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		param.Name = strings.TrimSpace(line)
		return param
	}

	param.Name = strings.TrimSpace(parts[0])
	rest := strings.TrimSpace(parts[1])

	// Extract default value
	if idx := strings.Index(rest, "(defaults to"); idx != -1 {
		defaultPart := rest[idx:]
		rest = strings.TrimSpace(rest[:idx])
		// Extract the default value
		re := regexp.MustCompile(`\(defaults to "?([^"]+)"?\)`)
		matches := re.FindStringSubmatch(defaultPart)
		if len(matches) == 2 {
			param.Default = matches[1]
		}
	}

	// Extract description in parentheses
	if idx := strings.Index(rest, "("); idx != -1 {
		param.Type = strings.TrimSpace(rest[:idx])
		descEnd := strings.Index(rest, ")")
		if descEnd > idx {
			param.Description = rest[idx+1 : descEnd]
		}
	} else {
		param.Type = rest
	}

	return param
}

// parseField parses a field definition.
func parseField(line string) spec.Field {
	field := spec.Field{}

	// Format: name: Type (description)
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		field.Name = strings.TrimSpace(line)
		return field
	}

	field.Name = strings.TrimSpace(parts[0])
	rest := strings.TrimSpace(parts[1])

	// Check for optional
	if strings.HasPrefix(rest, "Optional ") {
		field.Optional = true
		rest = strings.TrimPrefix(rest, "Optional ")
	}

	// Extract description in parentheses
	if idx := strings.Index(rest, "("); idx != -1 {
		field.Type = strings.TrimSpace(rest[:idx])
		descEnd := strings.Index(rest, ")")
		if descEnd > idx {
			field.Description = rest[idx+1 : descEnd]
		}
	} else {
		field.Type = rest
	}

	return field
}

// parseTests extracts test definitions.
func (p *specParser) parseTests() {
	// Find the Tests section in raw source
	source := string(p.source)
	testsStart := strings.Index(source, "## Tests")
	if testsStart == -1 {
		return
	}

	// Find the end of the Tests section
	testsEnd := len(source)
	nextSection := strings.Index(source[testsStart+8:], "\n## ")
	if nextSection != -1 {
		testsEnd = testsStart + 8 + nextSection
	}

	testsContent := source[testsStart:testsEnd]

	// Parse tests grouped by function
	funcBlocks := strings.Split(testsContent, "\n### ")
	for i, block := range funcBlocks {
		if i == 0 {
			continue // Skip the section header
		}

		lines := strings.Split(block, "\n")
		if len(lines) == 0 {
			continue
		}

		funcName := strings.TrimSpace(lines[0])

		// Parse individual test cases (#### headings)
		testCases := strings.Split(block, "\n#### ")
		for j, testBlock := range testCases {
			if j == 0 {
				continue // Skip the function name part
			}

			test := parseTestCase(testBlock, funcName)
			if test.Name != "" {
				p.spec.Tests = append(p.spec.Tests, test)
			}
		}
	}
}

// parseTestCase parses a single test case.
func parseTestCase(block string, funcName string) spec.TestCase {
	lines := strings.Split(block, "\n")
	if len(lines) == 0 {
		return spec.TestCase{}
	}

	test := spec.TestCase{
		Function: funcName,
	}

	// First line is the test name (after "test: ")
	nameLine := strings.TrimSpace(lines[0])
	if strings.HasPrefix(nameLine, "test: ") {
		test.Name = strings.TrimPrefix(nameLine, "test: ")
	} else {
		test.Name = nameLine
	}

	// Parse given and expect
	for _, line := range lines[1:] {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "given:") {
			value := strings.TrimPrefix(trimmed, "given:")
			value = strings.TrimSpace(value)
			if value != "" {
				test.Given = value
			}
		}
		if strings.HasPrefix(trimmed, "expect:") {
			value := strings.TrimPrefix(trimmed, "expect:")
			value = strings.TrimSpace(value)
			if value != "" {
				test.Expect = value
			}
		}
	}

	return test
}

// getSection returns the content of a section by name (case-insensitive).
func (p *specParser) getSection(name string) string {
	// Direct lookup
	if content, ok := p.sections[name]; ok {
		return content
	}

	// Fallback: search in source
	source := string(p.source)
	lowerSource := strings.ToLower(source)
	searchKey := "## " + name

	idx := strings.Index(lowerSource, searchKey)
	if idx == -1 {
		return ""
	}

	// Find end of section
	endIdx := len(source)
	nextSection := strings.Index(lowerSource[idx+len(searchKey):], "\n## ")
	if nextSection != -1 {
		endIdx = idx + len(searchKey) + nextSection
	}

	return source[idx+len(searchKey) : endIdx]
}
