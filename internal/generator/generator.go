package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/languages"
	"github.com/kon1790/rpg/internal/specparser"
)

// Generator handles code generation from spec analysis.
type Generator struct {
	registry *languages.Registry
}

// NewGenerator creates a new code generator.
func NewGenerator(registry *languages.Registry) *Generator {
	return &Generator{
		registry: registry,
	}
}

// Generate generates code files from a spec analysis for the target language.
func (g *Generator) Generate(spec *specparser.SpecAnalysis, language, outputDir string) ([]GeneratedFile, error) {
	adapter, err := g.registry.Get(language)
	if err != nil {
		return nil, fmt.Errorf("unsupported language %s: %w", language, err)
	}

	var files []GeneratedFile

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get project structure
	projectFiles := adapter.GetProjectStructure(spec.Name, len(spec.Tests) > 0)

	// Generate types
	typeFiles := g.generateTypes(spec, adapter, outputDir)
	files = append(files, typeFiles...)

	// Generate functions (grouped by receiver/module)
	funcFiles := g.generateFunctions(spec, adapter, outputDir)
	files = append(files, funcFiles...)

	// Generate tests
	if len(spec.Tests) > 0 {
		testFiles := g.generateTests(spec, adapter, outputDir)
		files = append(files, testFiles...)
	}

	// Generate project files (go.mod, package.json, etc.)
	projFiles := g.generateProjectFiles(spec, adapter, projectFiles, outputDir)
	files = append(files, projFiles...)

	// Write all files
	for i, f := range files {
		fullPath := filepath.Join(outputDir, f.Path)

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			continue
		}

		if err := os.WriteFile(fullPath, []byte(f.Content), 0644); err != nil {
			continue
		}
		files[i].Size = len(f.Content)
	}

	return files, nil
}

// generateTypes generates type definition files.
func (g *Generator) generateTypes(spec *specparser.SpecAnalysis, adapter languages.LanguageAdapter, outputDir string) []GeneratedFile {
	if len(spec.Types) == 0 {
		return nil
	}

	lang := adapter.GetLanguage()
	var files []GeneratedFile

	// Group types or generate one file depending on language
	var content strings.Builder

	// Add package/module declaration based on language
	switch lang.ID {
	case "go":
		content.WriteString(fmt.Sprintf("package %s\n\n", toPackageName(spec.Name)))
	case "typescript":
		content.WriteString("// Type definitions\n\n")
	case "python":
		content.WriteString("from dataclasses import dataclass\nfrom typing import Optional, List, Dict, Any\n\n")
	case "java":
		content.WriteString(fmt.Sprintf("package %s;\n\n", toPackageName(spec.Name)))
	case "rust":
		content.WriteString("use serde::{Deserialize, Serialize};\n\n")
	case "csharp":
		content.WriteString(fmt.Sprintf("namespace %s\n{\n", toPascalCase(spec.Name)))
	}

	// Generate each type
	for _, t := range spec.Types {
		typeCode := g.generateType(t, lang)
		content.WriteString(typeCode)
		content.WriteString("\n")
	}

	// Close namespace for C#
	if lang.ID == "csharp" {
		content.WriteString("}\n")
	}

	// Determine file path based on language
	var filePath string
	switch lang.ID {
	case "go":
		filePath = "types.go"
	case "typescript":
		filePath = "src/types.ts"
	case "python":
		filePath = "src/types.py"
	case "java":
		filePath = fmt.Sprintf("src/main/java/%s/Types.java", toPackageName(spec.Name))
	case "rust":
		filePath = "src/types.rs"
	case "csharp":
		filePath = "src/Types.cs"
	default:
		filePath = "types.txt"
	}

	var elements []string
	for _, t := range spec.Types {
		elements = append(elements, t.Name)
	}

	files = append(files, GeneratedFile{
		Path:     filePath,
		Content:  content.String(),
		Category: "type",
		Elements: elements,
	})

	return files
}

// generateType generates code for a single type.
func (g *Generator) generateType(t specparser.SpecType, lang languages.Language) string {
	var sb strings.Builder

	// Add doc comment
	if t.Description != "" {
		switch lang.ID {
		case "go", "java", "csharp", "rust":
			sb.WriteString(fmt.Sprintf("// %s\n", t.Description))
		case "python":
			// Python docstrings are inside the class
		case "typescript":
			sb.WriteString(fmt.Sprintf("/** %s */\n", t.Description))
		}
	}

	switch t.Kind {
	case "struct", "class":
		sb.WriteString(g.generateStruct(t, lang))
	case "interface":
		sb.WriteString(g.generateInterface(t, lang))
	case "enum":
		sb.WriteString(g.generateEnum(t, lang))
	default:
		sb.WriteString(g.generateStruct(t, lang))
	}

	return sb.String()
}

// generateStruct generates a struct/class definition.
func (g *Generator) generateStruct(t specparser.SpecType, lang languages.Language) string {
	var sb strings.Builder

	switch lang.ID {
	case "go":
		sb.WriteString(fmt.Sprintf("type %s struct {\n", t.Name))
		for _, f := range t.Fields {
			fieldType := mapType(f.Type, lang.ID)
			jsonTag := fmt.Sprintf("`json:\"%s\"`", toSnakeCase(f.Name))
			sb.WriteString(fmt.Sprintf("\t%s %s %s\n", toPascalCase(f.Name), fieldType, jsonTag))
		}
		sb.WriteString("}\n")

	case "typescript":
		sb.WriteString(fmt.Sprintf("export interface %s {\n", t.Name))
		for _, f := range t.Fields {
			fieldType := mapType(f.Type, lang.ID)
			optional := ""
			if !f.Required {
				optional = "?"
			}
			sb.WriteString(fmt.Sprintf("  %s%s: %s;\n", toCamelCase(f.Name), optional, fieldType))
		}
		sb.WriteString("}\n")

	case "python":
		sb.WriteString("@dataclass\n")
		sb.WriteString(fmt.Sprintf("class %s:\n", t.Name))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("    \"\"\"%s\"\"\"\n", t.Description))
		}
		if len(t.Fields) == 0 {
			sb.WriteString("    pass\n")
		} else {
			for _, f := range t.Fields {
				fieldType := mapType(f.Type, lang.ID)
				if !f.Required {
					fieldType = fmt.Sprintf("Optional[%s]", fieldType)
					sb.WriteString(fmt.Sprintf("    %s: %s = None\n", toSnakeCase(f.Name), fieldType))
				} else {
					sb.WriteString(fmt.Sprintf("    %s: %s\n", toSnakeCase(f.Name), fieldType))
				}
			}
		}

	case "java":
		sb.WriteString(fmt.Sprintf("public class %s {\n", t.Name))
		for _, f := range t.Fields {
			fieldType := mapType(f.Type, lang.ID)
			sb.WriteString(fmt.Sprintf("    private %s %s;\n", fieldType, toCamelCase(f.Name)))
		}
		sb.WriteString("\n")
		// Generate getters/setters
		for _, f := range t.Fields {
			fieldType := mapType(f.Type, lang.ID)
			fieldName := toCamelCase(f.Name)
			pascalName := toPascalCase(f.Name)
			sb.WriteString(fmt.Sprintf("    public %s get%s() { return %s; }\n", fieldType, pascalName, fieldName))
			sb.WriteString(fmt.Sprintf("    public void set%s(%s %s) { this.%s = %s; }\n\n", pascalName, fieldType, fieldName, fieldName, fieldName))
		}
		sb.WriteString("}\n")

	case "rust":
		sb.WriteString("#[derive(Debug, Clone, Serialize, Deserialize)]\n")
		sb.WriteString(fmt.Sprintf("pub struct %s {\n", t.Name))
		for _, f := range t.Fields {
			fieldType := mapType(f.Type, lang.ID)
			if !f.Required {
				fieldType = fmt.Sprintf("Option<%s>", fieldType)
			}
			sb.WriteString(fmt.Sprintf("    pub %s: %s,\n", toSnakeCase(f.Name), fieldType))
		}
		sb.WriteString("}\n")

	case "csharp":
		sb.WriteString(fmt.Sprintf("    public class %s\n    {\n", t.Name))
		for _, f := range t.Fields {
			fieldType := mapType(f.Type, lang.ID)
			if !f.Required {
				fieldType += "?"
			}
			sb.WriteString(fmt.Sprintf("        public %s %s { get; set; }\n", fieldType, toPascalCase(f.Name)))
		}
		sb.WriteString("    }\n")
	}

	return sb.String()
}

// generateInterface generates an interface definition.
func (g *Generator) generateInterface(t specparser.SpecType, lang languages.Language) string {
	var sb strings.Builder

	switch lang.ID {
	case "go":
		sb.WriteString(fmt.Sprintf("type %s interface {\n", t.Name))
		for _, m := range t.Methods {
			sb.WriteString(fmt.Sprintf("\t%s\n", m))
		}
		sb.WriteString("}\n")

	case "typescript":
		sb.WriteString(fmt.Sprintf("export interface %s {\n", t.Name))
		for _, m := range t.Methods {
			sb.WriteString(fmt.Sprintf("  %s;\n", m))
		}
		sb.WriteString("}\n")

	case "python":
		sb.WriteString("from abc import ABC, abstractmethod\n\n")
		sb.WriteString(fmt.Sprintf("class %s(ABC):\n", t.Name))
		for _, m := range t.Methods {
			sb.WriteString("    @abstractmethod\n")
			sb.WriteString(fmt.Sprintf("    def %s(self): pass\n", toSnakeCase(m)))
		}

	case "java":
		sb.WriteString(fmt.Sprintf("public interface %s {\n", t.Name))
		for _, m := range t.Methods {
			sb.WriteString(fmt.Sprintf("    %s;\n", m))
		}
		sb.WriteString("}\n")

	case "rust":
		sb.WriteString(fmt.Sprintf("pub trait %s {\n", t.Name))
		for _, m := range t.Methods {
			sb.WriteString(fmt.Sprintf("    fn %s(&self);\n", toSnakeCase(m)))
		}
		sb.WriteString("}\n")

	case "csharp":
		sb.WriteString(fmt.Sprintf("    public interface %s\n    {\n", t.Name))
		for _, m := range t.Methods {
			sb.WriteString(fmt.Sprintf("        %s;\n", m))
		}
		sb.WriteString("    }\n")
	}

	return sb.String()
}

// generateEnum generates an enum definition.
func (g *Generator) generateEnum(t specparser.SpecType, lang languages.Language) string {
	var sb strings.Builder

	switch lang.ID {
	case "go":
		sb.WriteString(fmt.Sprintf("type %s string\n\nconst (\n", t.Name))
		for _, v := range t.Values {
			sb.WriteString(fmt.Sprintf("\t%s%s %s = \"%s\"\n", t.Name, toPascalCase(v.Name), t.Name, v.Name))
		}
		sb.WriteString(")\n")

	case "typescript":
		sb.WriteString(fmt.Sprintf("export enum %s {\n", t.Name))
		for _, v := range t.Values {
			if v.Value != "" {
				sb.WriteString(fmt.Sprintf("  %s = %s,\n", v.Name, v.Value))
			} else {
				sb.WriteString(fmt.Sprintf("  %s,\n", v.Name))
			}
		}
		sb.WriteString("}\n")

	case "python":
		sb.WriteString("from enum import Enum\n\n")
		sb.WriteString(fmt.Sprintf("class %s(Enum):\n", t.Name))
		for _, v := range t.Values {
			if v.Value != "" {
				sb.WriteString(fmt.Sprintf("    %s = %s\n", v.Name, v.Value))
			} else {
				sb.WriteString(fmt.Sprintf("    %s = \"%s\"\n", v.Name, v.Name))
			}
		}

	case "java":
		sb.WriteString(fmt.Sprintf("public enum %s {\n", t.Name))
		for i, v := range t.Values {
			if i < len(t.Values)-1 {
				sb.WriteString(fmt.Sprintf("    %s,\n", v.Name))
			} else {
				sb.WriteString(fmt.Sprintf("    %s\n", v.Name))
			}
		}
		sb.WriteString("}\n")

	case "rust":
		sb.WriteString("#[derive(Debug, Clone, Serialize, Deserialize)]\n")
		sb.WriteString(fmt.Sprintf("pub enum %s {\n", t.Name))
		for _, v := range t.Values {
			sb.WriteString(fmt.Sprintf("    %s,\n", toPascalCase(v.Name)))
		}
		sb.WriteString("}\n")

	case "csharp":
		sb.WriteString(fmt.Sprintf("    public enum %s\n    {\n", t.Name))
		for _, v := range t.Values {
			sb.WriteString(fmt.Sprintf("        %s,\n", toPascalCase(v.Name)))
		}
		sb.WriteString("    }\n")
	}

	return sb.String()
}

// generateFunctions generates function/method files.
func (g *Generator) generateFunctions(spec *specparser.SpecAnalysis, adapter languages.LanguageAdapter, outputDir string) []GeneratedFile {
	if len(spec.Functions) == 0 {
		return nil
	}

	lang := adapter.GetLanguage()
	var files []GeneratedFile

	// Group functions by receiver (for methods) or generate one file
	var content strings.Builder

	// Add package/module declaration
	switch lang.ID {
	case "go":
		content.WriteString(fmt.Sprintf("package %s\n\n", toPackageName(spec.Name)))
	case "typescript":
		content.WriteString("// Functions\n\n")
	case "python":
		content.WriteString("from typing import Optional, Any\n\n")
	case "java":
		content.WriteString(fmt.Sprintf("package %s;\n\npublic class Service {\n", toPackageName(spec.Name)))
	case "rust":
		content.WriteString("use crate::types::*;\n\n")
	case "csharp":
		content.WriteString(fmt.Sprintf("namespace %s\n{\n    public static class Service\n    {\n", toPascalCase(spec.Name)))
	}

	// Generate each function
	for _, f := range spec.Functions {
		funcCode := g.generateFunction(f, lang)
		content.WriteString(funcCode)
		content.WriteString("\n")
	}

	// Close class/namespace for Java/C#
	if lang.ID == "java" {
		content.WriteString("}\n")
	} else if lang.ID == "csharp" {
		content.WriteString("    }\n}\n")
	}

	// Determine file path
	var filePath string
	switch lang.ID {
	case "go":
		filePath = "service.go"
	case "typescript":
		filePath = "src/service.ts"
	case "python":
		filePath = "src/service.py"
	case "java":
		filePath = fmt.Sprintf("src/main/java/%s/Service.java", toPackageName(spec.Name))
	case "rust":
		filePath = "src/service.rs"
	case "csharp":
		filePath = "src/Service.cs"
	default:
		filePath = "service.txt"
	}

	var elements []string
	for _, f := range spec.Functions {
		elements = append(elements, f.Name)
	}

	files = append(files, GeneratedFile{
		Path:     filePath,
		Content:  content.String(),
		Category: "function",
		Elements: elements,
	})

	return files
}

// generateFunction generates code for a single function.
func (g *Generator) generateFunction(f specparser.SpecFunction, lang languages.Language) string {
	var sb strings.Builder

	// Add doc comment
	if f.Description != "" {
		switch lang.ID {
		case "go":
			sb.WriteString(fmt.Sprintf("// %s %s\n", f.Name, f.Description))
		case "typescript":
			sb.WriteString(fmt.Sprintf("/**\n * %s\n */\n", f.Description))
		case "python":
			// Python docstrings are inside the function
		case "java", "csharp":
			sb.WriteString(fmt.Sprintf("    /**\n     * %s\n     */\n", f.Description))
		case "rust":
			sb.WriteString(fmt.Sprintf("/// %s\n", f.Description))
		}
	}

	// Generate function signature and body
	switch lang.ID {
	case "go":
		sb.WriteString(g.generateGoFunction(f))
	case "typescript":
		sb.WriteString(g.generateTypeScriptFunction(f))
	case "python":
		sb.WriteString(g.generatePythonFunction(f))
	case "java":
		sb.WriteString(g.generateJavaFunction(f))
	case "rust":
		sb.WriteString(g.generateRustFunction(f))
	case "csharp":
		sb.WriteString(g.generateCSharpFunction(f))
	}

	return sb.String()
}

// generateGoFunction generates a Go function.
func (g *Generator) generateGoFunction(f specparser.SpecFunction) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range f.Parameters {
		paramType := mapType(p.Type, "go")
		params = append(params, fmt.Sprintf("%s %s", toCamelCase(p.Name), paramType))
	}

	// Build return types
	var returns []string
	for _, r := range f.Returns {
		returns = append(returns, mapType(r.Type, "go"))
	}

	// Add error return if function has errors
	if len(f.Errors) > 0 && !containsError(returns) {
		returns = append(returns, "error")
	}

	returnStr := ""
	if len(returns) == 1 {
		returnStr = returns[0]
	} else if len(returns) > 1 {
		returnStr = fmt.Sprintf("(%s)", strings.Join(returns, ", "))
	}

	sb.WriteString(fmt.Sprintf("func %s(%s) %s {\n", f.Name, strings.Join(params, ", "), returnStr))

	// Generate body with TODO
	if f.Logic != "" {
		sb.WriteString(fmt.Sprintf("\t// %s\n", f.Logic))
	}
	sb.WriteString("\t// TODO: Implement\n")
	if len(returns) > 0 {
		if containsError(returns) {
			sb.WriteString("\treturn nil, nil\n")
		} else {
			sb.WriteString(fmt.Sprintf("\treturn %s\n", defaultValue(returns[0], "go")))
		}
	}
	sb.WriteString("}\n")

	return sb.String()
}

// generateTypeScriptFunction generates a TypeScript function.
func (g *Generator) generateTypeScriptFunction(f specparser.SpecFunction) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range f.Parameters {
		paramType := mapType(p.Type, "typescript")
		params = append(params, fmt.Sprintf("%s: %s", toCamelCase(p.Name), paramType))
	}

	// Build return type
	returnType := "void"
	if len(f.Returns) > 0 {
		returnType = mapType(f.Returns[0].Type, "typescript")
	}

	asyncPrefix := ""
	if f.IsAsync {
		asyncPrefix = "async "
		if returnType != "void" {
			returnType = fmt.Sprintf("Promise<%s>", returnType)
		} else {
			returnType = "Promise<void>"
		}
	}

	sb.WriteString(fmt.Sprintf("export %sfunction %s(%s): %s {\n", asyncPrefix, toCamelCase(f.Name), strings.Join(params, ", "), returnType))

	if f.Logic != "" {
		sb.WriteString(fmt.Sprintf("  // %s\n", f.Logic))
	}
	sb.WriteString("  // TODO: Implement\n")
	if returnType != "void" && returnType != "Promise<void>" {
		sb.WriteString(fmt.Sprintf("  return %s;\n", defaultValue(returnType, "typescript")))
	}
	sb.WriteString("}\n")

	return sb.String()
}

// generatePythonFunction generates a Python function.
func (g *Generator) generatePythonFunction(f specparser.SpecFunction) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range f.Parameters {
		paramType := mapType(p.Type, "python")
		params = append(params, fmt.Sprintf("%s: %s", toSnakeCase(p.Name), paramType))
	}

	// Build return type
	returnType := ""
	if len(f.Returns) > 0 {
		returnType = fmt.Sprintf(" -> %s", mapType(f.Returns[0].Type, "python"))
	}

	asyncPrefix := ""
	if f.IsAsync {
		asyncPrefix = "async "
	}

	sb.WriteString(fmt.Sprintf("%sdef %s(%s)%s:\n", asyncPrefix, toSnakeCase(f.Name), strings.Join(params, ", "), returnType))

	if f.Description != "" {
		sb.WriteString(fmt.Sprintf("    \"\"\"%s\"\"\"\n", f.Description))
	}

	if f.Logic != "" {
		sb.WriteString(fmt.Sprintf("    # %s\n", f.Logic))
	}
	sb.WriteString("    # TODO: Implement\n")
	sb.WriteString("    pass\n")

	return sb.String()
}

// generateJavaFunction generates a Java method.
func (g *Generator) generateJavaFunction(f specparser.SpecFunction) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range f.Parameters {
		paramType := mapType(p.Type, "java")
		params = append(params, fmt.Sprintf("%s %s", paramType, toCamelCase(p.Name)))
	}

	// Build return type
	returnType := "void"
	if len(f.Returns) > 0 {
		returnType = mapType(f.Returns[0].Type, "java")
	}

	visibility := "public"
	if !f.IsPublic {
		visibility = "private"
	}

	sb.WriteString(fmt.Sprintf("    %s static %s %s(%s) {\n", visibility, returnType, toCamelCase(f.Name), strings.Join(params, ", ")))

	if f.Logic != "" {
		sb.WriteString(fmt.Sprintf("        // %s\n", f.Logic))
	}
	sb.WriteString("        // TODO: Implement\n")
	if returnType != "void" {
		sb.WriteString(fmt.Sprintf("        return %s;\n", defaultValue(returnType, "java")))
	}
	sb.WriteString("    }\n")

	return sb.String()
}

// generateRustFunction generates a Rust function.
func (g *Generator) generateRustFunction(f specparser.SpecFunction) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range f.Parameters {
		paramType := mapType(p.Type, "rust")
		params = append(params, fmt.Sprintf("%s: %s", toSnakeCase(p.Name), paramType))
	}

	// Build return type
	returnType := ""
	if len(f.Returns) > 0 {
		returnType = fmt.Sprintf(" -> %s", mapType(f.Returns[0].Type, "rust"))
	}

	visibility := "pub "
	if !f.IsPublic {
		visibility = ""
	}

	asyncPrefix := ""
	if f.IsAsync {
		asyncPrefix = "async "
	}

	sb.WriteString(fmt.Sprintf("%s%sfn %s(%s)%s {\n", visibility, asyncPrefix, toSnakeCase(f.Name), strings.Join(params, ", "), returnType))

	if f.Logic != "" {
		sb.WriteString(fmt.Sprintf("    // %s\n", f.Logic))
	}
	sb.WriteString("    // TODO: Implement\n")
	sb.WriteString("    todo!()\n")
	sb.WriteString("}\n")

	return sb.String()
}

// generateCSharpFunction generates a C# method.
func (g *Generator) generateCSharpFunction(f specparser.SpecFunction) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range f.Parameters {
		paramType := mapType(p.Type, "csharp")
		params = append(params, fmt.Sprintf("%s %s", paramType, toCamelCase(p.Name)))
	}

	// Build return type
	returnType := "void"
	if len(f.Returns) > 0 {
		returnType = mapType(f.Returns[0].Type, "csharp")
	}

	asyncPrefix := ""
	if f.IsAsync {
		asyncPrefix = "async "
		if returnType != "void" {
			returnType = fmt.Sprintf("Task<%s>", returnType)
		} else {
			returnType = "Task"
		}
	}

	visibility := "public"
	if !f.IsPublic {
		visibility = "private"
	}

	sb.WriteString(fmt.Sprintf("        %s static %s%s %s(%s)\n        {\n", visibility, asyncPrefix, returnType, toPascalCase(f.Name), strings.Join(params, ", ")))

	if f.Logic != "" {
		sb.WriteString(fmt.Sprintf("            // %s\n", f.Logic))
	}
	sb.WriteString("            // TODO: Implement\n")
	sb.WriteString("            throw new NotImplementedException();\n")
	sb.WriteString("        }\n")

	return sb.String()
}

// generateTests generates test files.
func (g *Generator) generateTests(spec *specparser.SpecAnalysis, adapter languages.LanguageAdapter, outputDir string) []GeneratedFile {
	if len(spec.Tests) == 0 {
		return nil
	}

	lang := adapter.GetLanguage()
	var files []GeneratedFile

	var content strings.Builder

	// Add test framework imports based on language
	switch lang.ID {
	case "go":
		content.WriteString(fmt.Sprintf("package %s\n\nimport \"testing\"\n\n", toPackageName(spec.Name)))
	case "typescript":
		content.WriteString("import { describe, it, expect } from 'vitest';\n\n")
	case "python":
		content.WriteString("import pytest\n\n")
	case "java":
		content.WriteString(fmt.Sprintf("package %s;\n\nimport org.junit.jupiter.api.Test;\nimport static org.junit.jupiter.api.Assertions.*;\n\npublic class Tests {\n", toPackageName(spec.Name)))
	case "rust":
		content.WriteString("#[cfg(test)]\nmod tests {\n    use super::*;\n\n")
	case "csharp":
		content.WriteString(fmt.Sprintf("namespace %s.Tests\n{\n    using Xunit;\n\n    public class Tests\n    {\n", toPascalCase(spec.Name)))
	}

	// Generate each test
	for _, t := range spec.Tests {
		testCode := g.generateTest(t, lang)
		content.WriteString(testCode)
		content.WriteString("\n")
	}

	// Close class/module
	switch lang.ID {
	case "java":
		content.WriteString("}\n")
	case "rust":
		content.WriteString("}\n")
	case "csharp":
		content.WriteString("    }\n}\n")
	}

	// Determine file path
	var filePath string
	switch lang.ID {
	case "go":
		filePath = "service_test.go"
	case "typescript":
		filePath = "src/service.test.ts"
	case "python":
		filePath = "tests/test_service.py"
	case "java":
		filePath = fmt.Sprintf("src/test/java/%s/Tests.java", toPackageName(spec.Name))
	case "rust":
		filePath = "src/tests.rs"
	case "csharp":
		filePath = "tests/Tests.cs"
	default:
		filePath = "tests.txt"
	}

	var elements []string
	for _, t := range spec.Tests {
		elements = append(elements, t.Name)
	}

	files = append(files, GeneratedFile{
		Path:     filePath,
		Content:  content.String(),
		Category: "test",
		Elements: elements,
	})

	return files
}

// generateTest generates code for a single test.
func (g *Generator) generateTest(t specparser.SpecTest, lang languages.Language) string {
	var sb strings.Builder

	testName := toPascalCase(strings.ReplaceAll(t.Name, " ", "_"))

	switch lang.ID {
	case "go":
		sb.WriteString(fmt.Sprintf("func Test%s(t *testing.T) {\n", testName))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("\t// %s\n", t.Description))
		}
		for _, g := range t.Given {
			sb.WriteString(fmt.Sprintf("\t// Given: %s\n", g.Description))
		}
		if t.When != "" {
			sb.WriteString(fmt.Sprintf("\t// When: %s\n", t.When))
		}
		for _, a := range t.Then {
			sb.WriteString(fmt.Sprintf("\t// Then: %s\n", a.Description))
		}
		sb.WriteString("\t// TODO: Implement test\n")
		sb.WriteString("}\n")

	case "typescript":
		sb.WriteString(fmt.Sprintf("describe('%s', () => {\n", t.Name))
		sb.WriteString(fmt.Sprintf("  it('%s', () => {\n", t.Description))
		for _, g := range t.Given {
			sb.WriteString(fmt.Sprintf("    // Given: %s\n", g.Description))
		}
		if t.When != "" {
			sb.WriteString(fmt.Sprintf("    // When: %s\n", t.When))
		}
		for _, a := range t.Then {
			sb.WriteString(fmt.Sprintf("    // Then: %s\n", a.Description))
		}
		sb.WriteString("    // TODO: Implement test\n")
		sb.WriteString("  });\n")
		sb.WriteString("});\n")

	case "python":
		sb.WriteString(fmt.Sprintf("def test_%s():\n", toSnakeCase(testName)))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("    \"\"\"%s\"\"\"\n", t.Description))
		}
		for _, g := range t.Given {
			sb.WriteString(fmt.Sprintf("    # Given: %s\n", g.Description))
		}
		if t.When != "" {
			sb.WriteString(fmt.Sprintf("    # When: %s\n", t.When))
		}
		for _, a := range t.Then {
			sb.WriteString(fmt.Sprintf("    # Then: %s\n", a.Description))
		}
		sb.WriteString("    # TODO: Implement test\n")
		sb.WriteString("    pass\n")

	case "java":
		sb.WriteString(fmt.Sprintf("    @Test\n    public void test%s() {\n", testName))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("        // %s\n", t.Description))
		}
		for _, g := range t.Given {
			sb.WriteString(fmt.Sprintf("        // Given: %s\n", g.Description))
		}
		if t.When != "" {
			sb.WriteString(fmt.Sprintf("        // When: %s\n", t.When))
		}
		for _, a := range t.Then {
			sb.WriteString(fmt.Sprintf("        // Then: %s\n", a.Description))
		}
		sb.WriteString("        // TODO: Implement test\n")
		sb.WriteString("    }\n")

	case "rust":
		sb.WriteString(fmt.Sprintf("    #[test]\n    fn test_%s() {\n", toSnakeCase(testName)))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("        // %s\n", t.Description))
		}
		for _, g := range t.Given {
			sb.WriteString(fmt.Sprintf("        // Given: %s\n", g.Description))
		}
		if t.When != "" {
			sb.WriteString(fmt.Sprintf("        // When: %s\n", t.When))
		}
		for _, a := range t.Then {
			sb.WriteString(fmt.Sprintf("        // Then: %s\n", a.Description))
		}
		sb.WriteString("        // TODO: Implement test\n")
		sb.WriteString("        todo!()\n")
		sb.WriteString("    }\n")

	case "csharp":
		sb.WriteString(fmt.Sprintf("        [Fact]\n        public void Test%s()\n        {\n", testName))
		if t.Description != "" {
			sb.WriteString(fmt.Sprintf("            // %s\n", t.Description))
		}
		for _, g := range t.Given {
			sb.WriteString(fmt.Sprintf("            // Given: %s\n", g.Description))
		}
		if t.When != "" {
			sb.WriteString(fmt.Sprintf("            // When: %s\n", t.When))
		}
		for _, a := range t.Then {
			sb.WriteString(fmt.Sprintf("            // Then: %s\n", a.Description))
		}
		sb.WriteString("            // TODO: Implement test\n")
		sb.WriteString("            throw new NotImplementedException();\n")
		sb.WriteString("        }\n")
	}

	return sb.String()
}

// generateProjectFiles generates project configuration files.
func (g *Generator) generateProjectFiles(spec *specparser.SpecAnalysis, adapter languages.LanguageAdapter, projectFiles []languages.ProjectFile, outputDir string) []GeneratedFile {
	lang := adapter.GetLanguage()
	var files []GeneratedFile

	switch lang.ID {
	case "go":
		files = append(files, GeneratedFile{
			Path:     "go.mod",
			Content:  fmt.Sprintf("module %s\n\ngo 1.21\n", toPackageName(spec.Name)),
			Category: "config",
		})

	case "typescript":
		files = append(files, GeneratedFile{
			Path: "package.json",
			Content: fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "type": "module",
  "main": "dist/index.js",
  "scripts": {
    "build": "tsc",
    "test": "vitest"
  },
  "devDependencies": {
    "typescript": "^5.0.0",
    "vitest": "^1.0.0"
  }
}
`, toPackageName(spec.Name)),
			Category: "config",
		})

		files = append(files, GeneratedFile{
			Path: "tsconfig.json",
			Content: `{
  "compilerOptions": {
    "target": "ES2020",
    "module": "ESNext",
    "moduleResolution": "node",
    "strict": true,
    "outDir": "./dist",
    "rootDir": "./src"
  },
  "include": ["src/**/*"]
}
`,
			Category: "config",
		})

	case "python":
		files = append(files, GeneratedFile{
			Path: "pyproject.toml",
			Content: fmt.Sprintf(`[project]
name = "%s"
version = "1.0.0"
requires-python = ">=3.10"

[build-system]
requires = ["setuptools>=61.0"]
build-backend = "setuptools.build_meta"

[tool.pytest.ini_options]
testpaths = ["tests"]
`, toPackageName(spec.Name)),
			Category: "config",
		})

	case "rust":
		files = append(files, GeneratedFile{
			Path: "Cargo.toml",
			Content: fmt.Sprintf(`[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[dependencies]
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
`, toPackageName(spec.Name)),
			Category: "config",
		})

		files = append(files, GeneratedFile{
			Path:     "src/lib.rs",
			Content:  "pub mod types;\npub mod service;\n",
			Category: "config",
		})
	}

	return files
}

// Helper functions

func mapType(pseudoType, lang string) string {
	// Normalize the type
	t := strings.ToLower(strings.TrimSpace(pseudoType))

	// Handle arrays/lists
	if strings.HasPrefix(t, "[]") || strings.HasPrefix(t, "list[") || strings.HasPrefix(t, "array[") {
		innerType := strings.TrimPrefix(t, "[]")
		innerType = strings.TrimPrefix(innerType, "list[")
		innerType = strings.TrimPrefix(innerType, "array[")
		innerType = strings.TrimSuffix(innerType, "]")
		mapped := mapType(innerType, lang)
		switch lang {
		case "go":
			return "[]" + mapped
		case "typescript":
			return mapped + "[]"
		case "python":
			return fmt.Sprintf("List[%s]", mapped)
		case "java":
			return fmt.Sprintf("List<%s>", mapped)
		case "rust":
			return fmt.Sprintf("Vec<%s>", mapped)
		case "csharp":
			return fmt.Sprintf("List<%s>", mapped)
		}
	}

	// Handle maps
	if strings.HasPrefix(t, "map[") || strings.HasPrefix(t, "dict[") {
		switch lang {
		case "go":
			return "map[string]interface{}"
		case "typescript":
			return "Record<string, any>"
		case "python":
			return "Dict[str, Any]"
		case "java":
			return "Map<String, Object>"
		case "rust":
			return "HashMap<String, serde_json::Value>"
		case "csharp":
			return "Dictionary<string, object>"
		}
	}

	// Handle optional types
	if strings.HasPrefix(t, "optional[") || strings.HasSuffix(t, "?") {
		innerType := strings.TrimPrefix(t, "optional[")
		innerType = strings.TrimSuffix(innerType, "]")
		innerType = strings.TrimSuffix(innerType, "?")
		mapped := mapType(innerType, lang)
		switch lang {
		case "go":
			return "*" + mapped
		case "typescript":
			return mapped + " | null"
		case "python":
			return fmt.Sprintf("Optional[%s]", mapped)
		case "java":
			return mapped // Java uses null
		case "rust":
			return fmt.Sprintf("Option<%s>", mapped)
		case "csharp":
			return mapped + "?"
		}
	}

	// Map basic types
	typeMap := map[string]map[string]string{
		"string": {
			"go": "string", "typescript": "string", "python": "str",
			"java": "String", "rust": "String", "csharp": "string",
		},
		"str": {
			"go": "string", "typescript": "string", "python": "str",
			"java": "String", "rust": "String", "csharp": "string",
		},
		"int": {
			"go": "int", "typescript": "number", "python": "int",
			"java": "int", "rust": "i32", "csharp": "int",
		},
		"integer": {
			"go": "int", "typescript": "number", "python": "int",
			"java": "int", "rust": "i32", "csharp": "int",
		},
		"int64": {
			"go": "int64", "typescript": "number", "python": "int",
			"java": "long", "rust": "i64", "csharp": "long",
		},
		"float": {
			"go": "float64", "typescript": "number", "python": "float",
			"java": "double", "rust": "f64", "csharp": "double",
		},
		"float64": {
			"go": "float64", "typescript": "number", "python": "float",
			"java": "double", "rust": "f64", "csharp": "double",
		},
		"bool": {
			"go": "bool", "typescript": "boolean", "python": "bool",
			"java": "boolean", "rust": "bool", "csharp": "bool",
		},
		"boolean": {
			"go": "bool", "typescript": "boolean", "python": "bool",
			"java": "boolean", "rust": "bool", "csharp": "bool",
		},
		"any": {
			"go": "interface{}", "typescript": "any", "python": "Any",
			"java": "Object", "rust": "serde_json::Value", "csharp": "object",
		},
		"bytes": {
			"go": "[]byte", "typescript": "Uint8Array", "python": "bytes",
			"java": "byte[]", "rust": "Vec<u8>", "csharp": "byte[]",
		},
		"error": {
			"go": "error", "typescript": "Error", "python": "Exception",
			"java": "Exception", "rust": "anyhow::Error", "csharp": "Exception",
		},
		"void": {
			"go": "", "typescript": "void", "python": "None",
			"java": "void", "rust": "()", "csharp": "void",
		},
		"date": {
			"go": "time.Time", "typescript": "Date", "python": "datetime",
			"java": "LocalDate", "rust": "chrono::NaiveDate", "csharp": "DateTime",
		},
		"datetime": {
			"go": "time.Time", "typescript": "Date", "python": "datetime",
			"java": "LocalDateTime", "rust": "chrono::DateTime<chrono::Utc>", "csharp": "DateTime",
		},
		"uuid": {
			"go": "string", "typescript": "string", "python": "str",
			"java": "UUID", "rust": "uuid::Uuid", "csharp": "Guid",
		},
	}

	if langMap, ok := typeMap[t]; ok {
		if mapped, ok := langMap[lang]; ok {
			return mapped
		}
	}

	// Return original type (might be a custom type)
	return toPascalCase(pseudoType)
}

func defaultValue(typeName, lang string) string {
	t := strings.ToLower(typeName)

	switch {
	case strings.Contains(t, "string") || t == "str":
		return `""`
	case strings.Contains(t, "int") || strings.Contains(t, "number") || strings.Contains(t, "float"):
		return "0"
	case strings.Contains(t, "bool"):
		if lang == "python" {
			return "False"
		}
		return "false"
	case strings.Contains(t, "[]") || strings.Contains(t, "list") || strings.Contains(t, "vec"):
		switch lang {
		case "go":
			return "nil"
		case "typescript":
			return "[]"
		case "python":
			return "[]"
		case "java":
			return "new ArrayList<>()"
		case "rust":
			return "Vec::new()"
		case "csharp":
			return "new List<>()"
		}
	}

	switch lang {
	case "go":
		return "nil"
	case "typescript", "python":
		return "null"
	case "java", "csharp":
		return "null"
	case "rust":
		return "Default::default()"
	}

	return "nil"
}

func containsError(types []string) bool {
	for _, t := range types {
		if strings.ToLower(t) == "error" {
			return true
		}
	}
	return false
}

func toPackageName(name string) string {
	// Convert to lowercase, replace non-alphanumeric with empty
	result := strings.ToLower(name)
	result = strings.ReplaceAll(result, "-", "")
	result = strings.ReplaceAll(result, "_", "")
	result = strings.ReplaceAll(result, " ", "")
	return result
}

func toPascalCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(pascal[:1]) + pascal[1:]
	}
	return pascal
}

func toSnakeCase(s string) string {
	words := splitWords(s)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

func splitWords(s string) []string {
	// Split on underscores, hyphens, and camelCase boundaries
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")

	// Insert underscore before uppercase letters in camelCase
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(s[i-1])
			if prev >= 'a' && prev <= 'z' {
				result = append(result, '_')
			}
		}
		result = append(result, r)
	}

	return strings.Split(string(result), "_")
}
