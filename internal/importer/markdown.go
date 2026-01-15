package importer

import (
	"fmt"
	"strings"
)

// GenerateMarkdown generates a spec markdown file from extracted project data.
func GenerateMarkdown(project *ExtractedProject) string {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# %s\n\n", project.Name))

	// Description
	if project.Description != "" {
		sb.WriteString(project.Description)
		sb.WriteString("\n\n")
	}

	// Target Languages
	sb.WriteString("## Target Languages\n\n")
	sb.WriteString(fmt.Sprintf("- %s\n\n", project.DetectedLanguage))

	// Types section
	if len(project.Types) > 0 {
		sb.WriteString("## Types\n\n")
		for _, t := range project.Types {
			sb.WriteString(fmt.Sprintf("### %s\n\n", t.Name))

			if t.Description != "" {
				sb.WriteString(t.Description)
				sb.WriteString("\n\n")
			}

			switch t.Kind {
			case "enum":
				sb.WriteString("is one of:\n")
				for _, v := range t.Variants {
					sb.WriteString(fmt.Sprintf("- %s\n", v))
				}
			case "alias":
				sb.WriteString(fmt.Sprintf("alias of: %s\n", t.AliasOf))
			default: // struct, interface
				if len(t.Fields) > 0 {
					sb.WriteString("contains:\n")
					for _, f := range t.Fields {
						typeStr := f.Type
						if f.Optional {
							typeStr = "Optional " + typeStr
						}
						line := fmt.Sprintf("- %s: %s", f.Name, typeStr)
						if f.Description != "" {
							line += fmt.Sprintf(" (%s)", f.Description)
						}
						if f.Default != "" {
							line += fmt.Sprintf(" (defaults to %s)", f.Default)
						}
						sb.WriteString(line + "\n")
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	// Functions section
	if len(project.Functions) > 0 {
		sb.WriteString("## Functions\n\n")
		for _, fn := range project.Functions {
			// Function header with modifiers
			modifiers := ""
			if fn.IsAsync {
				modifiers += " [async]"
			}
			if fn.IsPure {
				modifiers += " [pure]"
			}
			sb.WriteString(fmt.Sprintf("### %s%s\n\n", fn.Name, modifiers))

			// Description
			if fn.Description != "" {
				sb.WriteString(fn.Description)
				sb.WriteString("\n\n")
			}

			// Parameters
			if len(fn.Parameters) > 0 {
				sb.WriteString("**accepts:**\n")
				for _, p := range fn.Parameters {
					line := fmt.Sprintf("- %s: %s", p.Name, p.Type)
					if p.Description != "" {
						line += fmt.Sprintf(" (%s)", p.Description)
					}
					if p.Default != "" {
						line += fmt.Sprintf(" (defaults to %s)", p.Default)
					}
					sb.WriteString(line + "\n")
				}
				sb.WriteString("\n")
			}

			// Return type
			if fn.Returns != "" && fn.Returns != "Nothing" {
				sb.WriteString(fmt.Sprintf("**returns:** %s\n\n", fn.Returns))
			}

			// Logic
			if fn.Logic != "" && fn.Logic != "// TODO: Add logic description" {
				sb.WriteString("**logic:**\n```\n")
				sb.WriteString(fn.Logic)
				sb.WriteString("\n```\n\n")
			}

			// Errors
			if len(fn.Errors) > 0 {
				sb.WriteString("**errors:**\n")
				for _, e := range fn.Errors {
					sb.WriteString(fmt.Sprintf("- %s\n", e))
				}
				sb.WriteString("\n")
			}
		}
	}

	// Tests section
	if len(project.Tests) > 0 {
		sb.WriteString("## Tests\n\n")

		// Group tests by function
		testsByFunc := make(map[string][]ExtractedTest)
		for _, test := range project.Tests {
			funcName := test.Function
			if funcName == "" {
				funcName = "general"
			}
			testsByFunc[funcName] = append(testsByFunc[funcName], test)
		}

		for funcName, tests := range testsByFunc {
			sb.WriteString(fmt.Sprintf("### %s\n\n", funcName))

			for _, test := range tests {
				sb.WriteString(fmt.Sprintf("#### test: %s\n", test.Name))

				if test.Given != nil {
					givenStr := formatTestValue(test.Given)
					sb.WriteString(fmt.Sprintf("given: %s\n", givenStr))
				}

				if test.When != "" {
					sb.WriteString(fmt.Sprintf("when: %s\n", test.When))
				}

				if test.Expect != nil {
					expectStr := formatTestValue(test.Expect)
					sb.WriteString(fmt.Sprintf("expect: %s\n", expectStr))
				}

				sb.WriteString("\n")
			}
		}
	}

	// Dependencies section (if any detected)
	if len(project.Dependencies) > 0 {
		sb.WriteString("## Dependencies\n\n")
		for _, dep := range project.Dependencies {
			sb.WriteString(fmt.Sprintf("- %s\n", dep))
		}
		sb.WriteString("\n")
	}

	// Warnings as comments (if any)
	if len(project.Warnings) > 0 {
		sb.WriteString("<!-- Import Warnings:\n")
		for _, w := range project.Warnings {
			sb.WriteString(fmt.Sprintf("  - %s\n", w))
		}
		sb.WriteString("-->\n")
	}

	return sb.String()
}

// formatTestValue formats a test value for markdown output.
func formatTestValue(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		// Quote strings
		return fmt.Sprintf(`"%s"`, v)
	case map[string]interface{}:
		// Format as YAML-like structure
		var parts []string
		for k, val := range v {
			parts = append(parts, fmt.Sprintf("  %s: %s", k, formatTestValue(val)))
		}
		return "\n" + strings.Join(parts, "\n")
	case []interface{}:
		var parts []string
		for _, val := range v {
			parts = append(parts, fmt.Sprintf("  - %s", formatTestValue(val)))
		}
		return "\n" + strings.Join(parts, "\n")
	default:
		return fmt.Sprintf("%v", v)
	}
}
