package treesitter

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
)

// CSharpParser implements LanguageParser for C#
type CSharpParser struct {
	baseParser
}

// NewCSharpParser creates a new C# parser
func NewCSharpParser() *CSharpParser {
	return &CSharpParser{
		baseParser: baseParser{
			lang:       LanguageCSharp,
			extensions: []string{".cs"},
			tsLang:     csharp.GetLanguage(),
		},
	}
}

// Parse parses C# source code
func (p *CSharpParser) Parse(code []byte, filename string) (*ParseResult, error) {
	tree, err := p.parseTree(code)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()
	result := &ParseResult{
		Language: LanguageCSharp,
		FileName: filename,
	}

	p.extractTypes(code, root, filename, result)
	p.extractImports(code, root, result)

	return result, nil
}

// extractTypes extracts class, interface, struct, enum, and record definitions
func (p *CSharpParser) extractTypes(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	typeNodes := collectNodes(root, func(n *sitter.Node) bool {
		switch n.Type() {
		case "class_declaration", "interface_declaration", "struct_declaration",
			"enum_declaration", "record_declaration", "record_struct_declaration":
			return true
		}
		return false
	})

	for _, node := range typeNodes {
		typeDef := p.parseTypeNode(code, node, filename, root)
		if typeDef != nil {
			result.Types = append(result.Types, *typeDef)
			// Extract methods from the type
			methods := p.extractMethodsFromType(code, node, filename)
			result.Functions = append(result.Functions, methods...)
		}
	}
}

// parseTypeNode extracts type information
func (p *CSharpParser) parseTypeNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *TypeDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	typeDef := &TypeDef{
		Name:     name,
		Location: nodeLocation(filename, node),
		ASTHash:  hashNode(code, node),
	}

	// Determine visibility
	typeDef.IsPublic = p.hasModifier(code, node, "public")

	// Determine kind
	switch node.Type() {
	case "class_declaration":
		typeDef.Kind = TypeKindClass
	case "interface_declaration":
		typeDef.Kind = TypeKindInterface
	case "struct_declaration", "record_struct_declaration":
		typeDef.Kind = TypeKindStruct
	case "enum_declaration":
		typeDef.Kind = TypeKindEnum
		typeDef.Variants = p.extractEnumMembers(code, node)
	case "record_declaration":
		typeDef.Kind = TypeKindClass // Records are reference types
	}

	// Extract base types
	baseListNode := findChildByFieldName(node, "bases")
	if baseListNode != nil {
		bases := p.extractBaseTypes(code, baseListNode)
		if len(bases) > 0 {
			// First base could be a class (for class) or interface
			if typeDef.Kind == TypeKindClass {
				typeDef.Extends = bases[0]
				if len(bases) > 1 {
					typeDef.Implements = bases[1:]
				}
			} else {
				typeDef.Implements = bases
			}
		}
	}

	// Extract type parameters
	typeParamsNode := findChildByFieldName(node, "type_parameters")
	if typeParamsNode != nil {
		typeDef.Generic = p.extractTypeParameters(code, typeParamsNode)
	}

	// Extract fields
	typeDef.Fields = p.extractFields(code, node)

	// Extract method names
	typeDef.Methods = p.extractMethodNames(code, node)

	// Extract doc comment
	typeDef.DocComment = p.extractXmlDoc(code, node, root)

	return typeDef
}

// hasModifier checks if a declaration has a specific modifier
func (p *CSharpParser) hasModifier(code []byte, node *sitter.Node, modifier string) bool {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "modifier" {
			text := nodeText(code, child)
			if text == modifier {
				return true
			}
		}
	}
	return false
}

// extractBaseTypes extracts base class and interfaces
func (p *CSharpParser) extractBaseTypes(code []byte, node *sitter.Node) []string {
	var bases []string

	walkTree(node, func(n *sitter.Node) bool {
		switch n.Type() {
		case "identifier", "generic_name", "qualified_name":
			bases = append(bases, nodeText(code, n))
			return false
		}
		return true
	})

	return bases
}

// extractTypeParameters extracts generic type parameters
func (p *CSharpParser) extractTypeParameters(code []byte, node *sitter.Node) []string {
	var params []string

	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "type_parameter" {
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				params = append(params, nodeText(code, nameNode))
			} else {
				// Try identifier directly
				idNode := findChildByType(n, "identifier")
				if idNode != nil {
					params = append(params, nodeText(code, idNode))
				}
			}
			return false
		}
		return true
	})

	return params
}

// extractEnumMembers extracts enum member names
func (p *CSharpParser) extractEnumMembers(code []byte, node *sitter.Node) []string {
	var members []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return members
	}

	walkTree(bodyNode, func(n *sitter.Node) bool {
		if n.Type() == "enum_member_declaration" {
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				members = append(members, nodeText(code, nameNode))
			}
			return false
		}
		return true
	})

	return members
}

// extractFields extracts field and property declarations
func (p *CSharpParser) extractFields(code []byte, node *sitter.Node) []Field {
	var fields []Field

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return fields
	}

	// Extract field declarations
	fieldDecls := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "field_declaration"
	})

	for _, decl := range fieldDecls {
		declFields := p.parseFieldDeclaration(code, decl)
		fields = append(fields, declFields...)
	}

	// Extract property declarations
	propDecls := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "property_declaration"
	})

	for _, decl := range propDecls {
		prop := p.parsePropertyDeclaration(code, decl)
		if prop != nil {
			fields = append(fields, *prop)
		}
	}

	return fields
}

// parseFieldDeclaration extracts fields from a field declaration
func (p *CSharpParser) parseFieldDeclaration(code []byte, node *sitter.Node) []Field {
	var fields []Field

	// Get type
	typeNode := findChildByFieldName(node, "type")
	typeName := ""
	if typeNode != nil {
		typeName = nodeText(code, typeNode)
	}

	// Check modifiers
	isReadonly := p.hasModifier(code, node, "readonly")

	// Get declarators
	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "variable_declarator" {
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				field := Field{
					Name:       nodeText(code, nameNode),
					Type:       typeName,
					IsReadonly: isReadonly,
				}

				initNode := findChildByFieldName(n, "initializer")
				if initNode != nil {
					field.Default = nodeText(code, initNode)
				}

				fields = append(fields, field)
			}
			return false
		}
		return true
	})

	return fields
}

// parsePropertyDeclaration extracts a property as a field
func (p *CSharpParser) parsePropertyDeclaration(code []byte, node *sitter.Node) *Field {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	field := &Field{
		Name: nodeText(code, nameNode),
	}

	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		field.Type = nodeText(code, typeNode)
	}

	// Check if readonly (has get but no set, or has init)
	accessors := findChildByFieldName(node, "accessors")
	if accessors != nil {
		hasSet := false
		hasInit := false
		walkTree(accessors, func(n *sitter.Node) bool {
			if n.Type() == "accessor_declaration" {
				text := nodeText(code, n)
				if strings.Contains(text, "set") {
					hasSet = true
				}
				if strings.Contains(text, "init") {
					hasInit = true
				}
			}
			return true
		})
		field.IsReadonly = !hasSet || hasInit
	}

	return field
}

// extractMethodNames extracts method names from a type
func (p *CSharpParser) extractMethodNames(code []byte, node *sitter.Node) []string {
	var methods []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return methods
	}

	methodDecls := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "method_declaration" || n.Type() == "constructor_declaration"
	})

	for _, decl := range methodDecls {
		nameNode := findChildByFieldName(decl, "name")
		if nameNode != nil {
			methods = append(methods, nodeText(code, nameNode))
		}
	}

	return methods
}

// extractMethodsFromType extracts full method definitions
func (p *CSharpParser) extractMethodsFromType(code []byte, typeNode *sitter.Node, filename string) []FunctionDef {
	var functions []FunctionDef

	bodyNode := findChildByFieldName(typeNode, "body")
	if bodyNode == nil {
		return functions
	}

	methodDecls := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "method_declaration"
	})

	for _, decl := range methodDecls {
		fn := p.parseMethodDeclaration(code, decl, filename, typeNode)
		if fn != nil {
			functions = append(functions, *fn)
		}
	}

	return functions
}

// parseMethodDeclaration extracts a method definition
func (p *CSharpParser) parseMethodDeclaration(code []byte, node *sitter.Node, filename string, parent *sitter.Node) *FunctionDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	fn := &FunctionDef{
		Name:     name,
		IsPublic: p.hasModifier(code, node, "public"),
		IsStatic: p.hasModifier(code, node, "static"),
		IsAsync:  p.hasModifier(code, node, "async"),
		Location: nodeLocation(filename, node),
		ASTHash:  hashNode(code, node),
	}

	// Extract return type
	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		fn.ReturnType = nodeText(code, typeNode)
	}

	// Extract parameters
	paramsNode := findChildByFieldName(node, "parameters")
	if paramsNode != nil {
		fn.Parameters = p.parseParameters(code, paramsNode)
	}

	// Build signature
	fn.Signature = p.buildSignature(fn)

	// Extract body
	bodyNode := findChildByFieldName(node, "body")
	if bodyNode != nil {
		fn.Body = nodeText(code, bodyNode)
		fn.Calls = p.extractCallsFromBody(code, bodyNode)
		fn.Complexity = p.calculateComplexity(bodyNode)
	}

	// Extract doc comment
	fn.DocComment = p.extractXmlDoc(code, node, parent)

	return fn
}

// parseParameters extracts parameters
func (p *CSharpParser) parseParameters(code []byte, paramsNode *sitter.Node) []Parameter {
	var params []Parameter

	paramDecls := collectNodes(paramsNode, func(n *sitter.Node) bool {
		return n.Type() == "parameter"
	})

	for _, decl := range paramDecls {
		param := p.parseParameter(code, decl)
		if param != nil {
			params = append(params, *param)
		}
	}

	return params
}

// parseParameter extracts a single parameter
func (p *CSharpParser) parseParameter(code []byte, node *sitter.Node) *Parameter {
	param := &Parameter{}

	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		param.Type = nodeText(code, typeNode)
	}

	nameNode := findChildByFieldName(node, "name")
	if nameNode != nil {
		param.Name = nodeText(code, nameNode)
	}

	// Check for default value
	defaultNode := findChildByFieldName(node, "default")
	if defaultNode != nil {
		param.DefaultValue = nodeText(code, defaultNode)
		param.IsOptional = true
	}

	// Check for params keyword (variadic)
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "params" {
			param.IsVariadic = true
			break
		}
	}

	return param
}

// buildSignature builds a method signature string
func (p *CSharpParser) buildSignature(fn *FunctionDef) string {
	var sb strings.Builder

	if fn.IsPublic {
		sb.WriteString("public ")
	}
	if fn.IsStatic {
		sb.WriteString("static ")
	}
	if fn.IsAsync {
		sb.WriteString("async ")
	}

	sb.WriteString(fn.ReturnType)
	sb.WriteString(" ")
	sb.WriteString(fn.Name)
	sb.WriteString("(")

	for i, param := range fn.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Type)
		sb.WriteString(" ")
		sb.WriteString(param.Name)
		if param.DefaultValue != "" {
			sb.WriteString(" = ")
			sb.WriteString(param.DefaultValue)
		}
	}
	sb.WriteString(")")

	return sb.String()
}

// extractCallsFromBody extracts method calls from a method body
func (p *CSharpParser) extractCallsFromBody(code []byte, bodyNode *sitter.Node) []string {
	calls := make(map[string]bool)

	callNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "invocation_expression"
	})

	for _, callNode := range callNodes {
		// Get the method name
		walkTree(callNode, func(n *sitter.Node) bool {
			if n.Type() == "member_access_expression" {
				nameNode := findChildByFieldName(n, "name")
				if nameNode != nil {
					calls[nodeText(code, nameNode)] = true
				}
				return false
			}
			if n.Type() == "identifier" {
				calls[nodeText(code, n)] = true
				return false
			}
			return true
		})
	}

	var result []string
	for call := range calls {
		result = append(result, call)
	}
	return result
}

// calculateComplexity calculates cyclomatic complexity
func (p *CSharpParser) calculateComplexity(bodyNode *sitter.Node) int {
	complexity := 1

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "if_statement", "for_statement", "foreach_statement",
			"while_statement", "do_statement", "switch_statement",
			"catch_clause", "conditional_expression", "&&", "||",
			"switch_expression_arm":
			complexity++
		}
		return true
	})

	return complexity
}

// extractXmlDoc extracts XML documentation comments
func (p *CSharpParser) extractXmlDoc(code []byte, node *sitter.Node, root *sitter.Node) string {
	comment := getCommentAbove(code, node, root)

	// Parse XML doc format: /// <summary>...</summary>
	if strings.Contains(comment, "<summary>") {
		start := strings.Index(comment, "<summary>")
		end := strings.Index(comment, "</summary>")
		if start != -1 && end != -1 && end > start {
			summary := comment[start+9 : end]
			summary = strings.TrimSpace(summary)
			return summary
		}
	}

	return comment
}

// extractImports extracts using statements
func (p *CSharpParser) extractImports(code []byte, root *sitter.Node, result *ParseResult) {
	usingNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "using_directive"
	})

	for _, node := range usingNodes {
		imp := p.parseUsingDirective(code, node)
		if imp != nil {
			result.Imports = append(result.Imports, *imp)
		}
	}
}

// parseUsingDirective extracts import information
func (p *CSharpParser) parseUsingDirective(code []byte, node *sitter.Node) *Import {
	// Check for static using
	isStatic := false
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "static" {
			isStatic = true
			break
		}
	}

	// Get namespace
	var path string
	walkTree(node, func(n *sitter.Node) bool {
		switch n.Type() {
		case "qualified_name", "identifier":
			path = nodeText(code, n)
			return false
		}
		return true
	})

	if path == "" {
		return nil
	}

	imp := &Import{
		Path:    path,
		IsLocal: false,
	}

	if isStatic {
		imp.Alias = "static"
	}

	// Check for alias
	aliasNode := findChildByType(node, "name_equals")
	if aliasNode != nil {
		nameNode := findChildByType(aliasNode, "identifier")
		if nameNode != nil {
			imp.Alias = nodeText(code, nameNode)
		}
	}

	return imp
}
