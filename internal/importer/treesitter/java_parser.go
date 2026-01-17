package treesitter

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
)

// JavaParser implements LanguageParser for Java
type JavaParser struct {
	baseParser
}

// NewJavaParser creates a new Java parser
func NewJavaParser() *JavaParser {
	return &JavaParser{
		baseParser: baseParser{
			lang:       LanguageJava,
			extensions: []string{".java"},
			tsLang:     java.GetLanguage(),
		},
	}
}

// Parse parses Java source code
func (p *JavaParser) Parse(code []byte, filename string) (*ParseResult, error) {
	tree, err := p.parseTree(code)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()
	result := &ParseResult{
		Language: LanguageJava,
		FileName: filename,
	}

	p.extractTypes(code, root, filename, result)
	p.extractImports(code, root, result)

	return result, nil
}

// extractTypes extracts class, interface, enum, and record definitions
func (p *JavaParser) extractTypes(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	typeNodes := collectNodes(root, func(n *sitter.Node) bool {
		switch n.Type() {
		case "class_declaration", "interface_declaration",
			"enum_declaration", "record_declaration", "annotation_type_declaration":
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
func (p *JavaParser) parseTypeNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *TypeDef {
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
	typeDef.IsPublic = p.hasModifier(node, "public")

	// Determine kind
	switch node.Type() {
	case "class_declaration":
		typeDef.Kind = TypeKindClass
	case "interface_declaration":
		typeDef.Kind = TypeKindInterface
	case "enum_declaration":
		typeDef.Kind = TypeKindEnum
		typeDef.Variants = p.extractEnumConstants(code, node)
	case "record_declaration":
		typeDef.Kind = TypeKindClass // Records are special classes
	case "annotation_type_declaration":
		typeDef.Kind = TypeKindInterface
	}

	// Extract superclass
	superNode := findChildByFieldName(node, "superclass")
	if superNode != nil {
		typeDef.Extends = p.extractTypeName(code, superNode)
	}

	// Extract interfaces
	interfacesNode := findChildByFieldName(node, "interfaces")
	if interfacesNode != nil {
		typeDef.Implements = p.extractTypeList(code, interfacesNode)
	}

	// Extract type parameters (generics)
	typeParamsNode := findChildByFieldName(node, "type_parameters")
	if typeParamsNode != nil {
		typeDef.Generic = p.extractTypeParameters(code, typeParamsNode)
	}

	// Extract fields
	typeDef.Fields = p.extractFields(code, node)

	// Extract method names
	typeDef.Methods = p.extractMethodNames(code, node)

	// Extract doc comment
	typeDef.DocComment = p.extractJavadoc(code, node, root)

	return typeDef
}

// hasModifier checks if a declaration has a specific modifier
func (p *JavaParser) hasModifier(node *sitter.Node, modifier string) bool {
	modifiersNode := findChildByFieldName(node, "modifiers")
	if modifiersNode == nil {
		return false
	}

	for i := 0; i < int(modifiersNode.ChildCount()); i++ {
		child := modifiersNode.Child(i)
		if child.Type() == modifier {
			return true
		}
	}
	return false
}

// extractTypeName extracts a type name from a node
func (p *JavaParser) extractTypeName(code []byte, node *sitter.Node) string {
	// Skip keywords like "extends", "implements"
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "type_identifier" || child.Type() == "generic_type" {
			return nodeText(code, child)
		}
	}
	return nodeText(code, node)
}

// extractTypeList extracts a list of type names
func (p *JavaParser) extractTypeList(code []byte, node *sitter.Node) []string {
	var types []string

	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "type_identifier" || n.Type() == "generic_type" {
			types = append(types, nodeText(code, n))
			return false
		}
		return true
	})

	return types
}

// extractTypeParameters extracts generic type parameters
func (p *JavaParser) extractTypeParameters(code []byte, node *sitter.Node) []string {
	var params []string

	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "type_parameter" {
			nameNode := findChildByType(n, "type_identifier")
			if nameNode != nil {
				params = append(params, nodeText(code, nameNode))
			}
			return false
		}
		return true
	})

	return params
}

// extractEnumConstants extracts enum variant names
func (p *JavaParser) extractEnumConstants(code []byte, node *sitter.Node) []string {
	var variants []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return variants
	}

	walkTree(bodyNode, func(n *sitter.Node) bool {
		if n.Type() == "enum_constant" {
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				variants = append(variants, nodeText(code, nameNode))
			}
			return false
		}
		return true
	})

	return variants
}

// extractFields extracts field declarations
func (p *JavaParser) extractFields(code []byte, node *sitter.Node) []Field {
	var fields []Field

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return fields
	}

	fieldDecls := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "field_declaration"
	})

	for _, decl := range fieldDecls {
		declFields := p.parseFieldDeclaration(code, decl)
		fields = append(fields, declFields...)
	}

	return fields
}

// parseFieldDeclaration extracts fields from a field declaration
func (p *JavaParser) parseFieldDeclaration(code []byte, node *sitter.Node) []Field {
	var fields []Field

	// Get type
	typeNode := findChildByFieldName(node, "type")
	typeName := ""
	if typeNode != nil {
		typeName = nodeText(code, typeNode)
	}

	// Check modifiers
	isReadonly := p.hasModifier(node, "final")

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

				valueNode := findChildByFieldName(n, "value")
				if valueNode != nil {
					field.Default = nodeText(code, valueNode)
				}

				fields = append(fields, field)
			}
			return false
		}
		return true
	})

	return fields
}

// extractMethodNames extracts method names from a type
func (p *JavaParser) extractMethodNames(code []byte, node *sitter.Node) []string {
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
func (p *JavaParser) extractMethodsFromType(code []byte, typeNode *sitter.Node, filename string) []FunctionDef {
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
func (p *JavaParser) parseMethodDeclaration(code []byte, node *sitter.Node, filename string, parent *sitter.Node) *FunctionDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	fn := &FunctionDef{
		Name:     name,
		IsPublic: p.hasModifier(node, "public"),
		IsStatic: p.hasModifier(node, "static"),
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

	// Extract javadoc
	fn.DocComment = p.extractJavadoc(code, node, parent)

	return fn
}

// parseParameters extracts parameters from a formal_parameters node
func (p *JavaParser) parseParameters(code []byte, paramsNode *sitter.Node) []Parameter {
	var params []Parameter

	paramDecls := collectNodes(paramsNode, func(n *sitter.Node) bool {
		return n.Type() == "formal_parameter" || n.Type() == "spread_parameter"
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
func (p *JavaParser) parseParameter(code []byte, node *sitter.Node) *Parameter {
	param := &Parameter{}

	if node.Type() == "spread_parameter" {
		param.IsVariadic = true
	}

	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		param.Type = nodeText(code, typeNode)
	}

	nameNode := findChildByFieldName(node, "name")
	if nameNode != nil {
		param.Name = nodeText(code, nameNode)
	}

	return param
}

// buildSignature builds a method signature string
func (p *JavaParser) buildSignature(fn *FunctionDef) string {
	var sb strings.Builder

	if fn.IsPublic {
		sb.WriteString("public ")
	}
	if fn.IsStatic {
		sb.WriteString("static ")
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
	}
	sb.WriteString(")")

	return sb.String()
}

// extractCallsFromBody extracts method calls from a method body
func (p *JavaParser) extractCallsFromBody(code []byte, bodyNode *sitter.Node) []string {
	calls := make(map[string]bool)

	callNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "method_invocation"
	})

	for _, callNode := range callNodes {
		nameNode := findChildByFieldName(callNode, "name")
		if nameNode != nil {
			callName := nodeText(code, nameNode)
			calls[callName] = true
		}
	}

	var result []string
	for call := range calls {
		result = append(result, call)
	}
	return result
}

// calculateComplexity calculates cyclomatic complexity
func (p *JavaParser) calculateComplexity(bodyNode *sitter.Node) int {
	complexity := 1

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "if_statement", "for_statement", "enhanced_for_statement",
			"while_statement", "do_statement", "switch_expression",
			"catch_clause", "ternary_expression", "&&", "||":
			complexity++
		case "switch_block_statement_group":
			// Each case adds complexity
			complexity++
		}
		return true
	})

	return complexity
}

// extractJavadoc extracts javadoc comment above a node
func (p *JavaParser) extractJavadoc(code []byte, node *sitter.Node, root *sitter.Node) string {
	// Look for block_comment or line_comment above
	comment := getCommentAbove(code, node, root)
	if strings.HasPrefix(comment, "*") {
		// Clean up javadoc format
		lines := strings.Split(comment, "\n")
		var cleaned []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "*")
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "@") {
				cleaned = append(cleaned, line)
			}
		}
		return strings.Join(cleaned, " ")
	}
	return comment
}

// extractImports extracts import declarations
func (p *JavaParser) extractImports(code []byte, root *sitter.Node, result *ParseResult) {
	importNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "import_declaration"
	})

	for _, node := range importNodes {
		imp := p.parseImportDeclaration(code, node)
		if imp != nil {
			result.Imports = append(result.Imports, *imp)
		}
	}
}

// parseImportDeclaration extracts import information
func (p *JavaParser) parseImportDeclaration(code []byte, node *sitter.Node) *Import {
	// Check for static import
	isStatic := false
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "static" {
			isStatic = true
			break
		}
	}

	// Get the import path
	var path string
	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "scoped_identifier" || n.Type() == "identifier" {
			if path == "" {
				path = nodeText(code, n)
			}
			return false
		}
		return true
	})

	if path == "" {
		return nil
	}

	imp := &Import{
		Path:    path,
		IsLocal: false, // Java imports are always external in terms of module system
	}

	if isStatic {
		imp.Alias = "static"
	}

	// Check for wildcard
	if strings.HasSuffix(path, ".*") {
		imp.Items = []string{"*"}
	}

	return imp
}
