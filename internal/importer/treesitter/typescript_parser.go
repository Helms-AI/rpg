package treesitter

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// TypeScriptParser implements LanguageParser for TypeScript
type TypeScriptParser struct {
	baseParser
}

// NewTypeScriptParser creates a new TypeScript parser
func NewTypeScriptParser() *TypeScriptParser {
	return &TypeScriptParser{
		baseParser: baseParser{
			lang:       LanguageTypeScript,
			extensions: []string{".ts", ".tsx"},
			tsLang:     typescript.GetLanguage(),
		},
	}
}

// Parse parses TypeScript source code
func (p *TypeScriptParser) Parse(code []byte, filename string) (*ParseResult, error) {
	tree, err := p.parseTree(code)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()
	result := &ParseResult{
		Language: LanguageTypeScript,
		FileName: filename,
	}

	p.extractFunctions(code, root, filename, result)
	p.extractTypes(code, root, filename, result)
	p.extractImports(code, root, result)

	return result, nil
}

// extractFunctions extracts all function declarations
func (p *TypeScriptParser) extractFunctions(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	// Find function declarations, arrow functions, and method definitions
	funcNodes := collectNodes(root, func(n *sitter.Node) bool {
		switch n.Type() {
		case "function_declaration", "method_definition", "arrow_function",
			"function_expression", "generator_function_declaration":
			return true
		}
		return false
	})

	for _, node := range funcNodes {
		fn := p.parseFunctionNode(code, node, filename, root)
		if fn != nil {
			result.Functions = append(result.Functions, *fn)
		}
	}
}

// parseFunctionNode extracts function information
func (p *TypeScriptParser) parseFunctionNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *FunctionDef {
	var name string

	// Extract name based on node type
	nameNode := findChildByFieldName(node, "name")
	if nameNode != nil {
		name = nodeText(code, nameNode)
	}

	// For arrow functions in variable declarations, get parent variable name
	if name == "" && node.Type() == "arrow_function" {
		parent := node.Parent()
		if parent != nil && parent.Type() == "variable_declarator" {
			parentNameNode := findChildByFieldName(parent, "name")
			if parentNameNode != nil {
				name = nodeText(code, parentNameNode)
			}
		}
	}

	if name == "" {
		return nil // Anonymous functions without assignment
	}

	fn := &FunctionDef{
		Name:     name,
		IsPublic: !strings.HasPrefix(name, "_"),
		Location: nodeLocation(filename, node),
		ASTHash:  hashNode(code, node),
	}

	// Check for async
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "async" {
			fn.IsAsync = true
			break
		}
	}

	// Extract parameters
	paramsNode := findChildByFieldName(node, "parameters")
	if paramsNode == nil {
		// Try formal_parameters for some function types
		paramsNode = findChildByType(node, "formal_parameters")
	}
	if paramsNode != nil {
		fn.Parameters = p.parseParameters(code, paramsNode)
	}

	// Extract return type
	returnTypeNode := findChildByFieldName(node, "return_type")
	if returnTypeNode != nil {
		fn.ReturnType = p.extractTypeAnnotation(code, returnTypeNode)
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

	// Extract documentation
	fn.DocComment = getCommentAbove(code, node, root)

	return fn
}

// parseParameters extracts parameters from a formal_parameters node
func (p *TypeScriptParser) parseParameters(code []byte, paramsNode *sitter.Node) []Parameter {
	var params []Parameter

	walkTree(paramsNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "required_parameter", "optional_parameter", "rest_parameter":
			param := p.parseParameter(code, n)
			if param != nil {
				params = append(params, *param)
			}
			return false // Don't descend into parameters
		}
		return true
	})

	return params
}

// parseParameter extracts a single parameter
func (p *TypeScriptParser) parseParameter(code []byte, node *sitter.Node) *Parameter {
	param := &Parameter{}

	// Handle rest parameter
	if node.Type() == "rest_parameter" {
		param.IsVariadic = true
	}

	// Handle optional parameter
	if node.Type() == "optional_parameter" {
		param.IsOptional = true
	}

	// Extract pattern (name)
	patternNode := findChildByFieldName(node, "pattern")
	if patternNode != nil {
		param.Name = nodeText(code, patternNode)
	} else {
		// Try direct identifier
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "identifier" {
				param.Name = nodeText(code, child)
				break
			}
		}
	}

	// Extract type annotation
	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		param.Type = p.extractTypeAnnotation(code, typeNode)
	}

	// Extract default value
	valueNode := findChildByFieldName(node, "value")
	if valueNode != nil {
		param.DefaultValue = nodeText(code, valueNode)
		param.IsOptional = true
	}

	return param
}

// extractTypeAnnotation extracts a type from a type_annotation node
func (p *TypeScriptParser) extractTypeAnnotation(code []byte, node *sitter.Node) string {
	// Skip the ":" if present
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() != ":" {
			return nodeText(code, child)
		}
	}
	return nodeText(code, node)
}

// buildSignature builds a function signature string
func (p *TypeScriptParser) buildSignature(fn *FunctionDef) string {
	var sb strings.Builder

	if fn.IsAsync {
		sb.WriteString("async ")
	}
	sb.WriteString("function ")
	sb.WriteString(fn.Name)
	sb.WriteString("(")

	for i, param := range fn.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
		if param.IsOptional {
			sb.WriteString("?")
		}
		if param.Type != "" {
			sb.WriteString(": ")
			sb.WriteString(param.Type)
		}
	}
	sb.WriteString(")")

	if fn.ReturnType != "" {
		sb.WriteString(": ")
		sb.WriteString(fn.ReturnType)
	}

	return sb.String()
}

// extractCallsFromBody extracts function calls from a function body
func (p *TypeScriptParser) extractCallsFromBody(code []byte, bodyNode *sitter.Node) []string {
	calls := make(map[string]bool)

	callNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "call_expression"
	})

	for _, callNode := range callNodes {
		funcNode := findChildByFieldName(callNode, "function")
		if funcNode != nil {
			callName := nodeText(code, funcNode)
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
func (p *TypeScriptParser) calculateComplexity(bodyNode *sitter.Node) int {
	complexity := 1

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "if_statement", "for_statement", "for_in_statement",
			"while_statement", "do_statement", "switch_statement",
			"case", "catch_clause", "ternary_expression",
			"&&", "||", "??":
			complexity++
		}
		return true
	})

	return complexity
}

// extractTypes extracts type definitions
func (p *TypeScriptParser) extractTypes(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	// Find interfaces, type aliases, classes, and enums
	typeNodes := collectNodes(root, func(n *sitter.Node) bool {
		switch n.Type() {
		case "interface_declaration", "type_alias_declaration",
			"class_declaration", "enum_declaration":
			return true
		}
		return false
	})

	for _, node := range typeNodes {
		typeDef := p.parseTypeNode(code, node, filename, root)
		if typeDef != nil {
			result.Types = append(result.Types, *typeDef)
		}
	}
}

// parseTypeNode extracts type information from a type declaration node
func (p *TypeScriptParser) parseTypeNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *TypeDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	typeDef := &TypeDef{
		Name:       name,
		IsPublic:   !strings.HasPrefix(name, "_"),
		Location:   nodeLocation(filename, node),
		ASTHash:    hashNode(code, node),
		DocComment: getCommentAbove(code, node, root),
	}

	switch node.Type() {
	case "interface_declaration":
		typeDef.Kind = TypeKindInterface
		typeDef.Fields, typeDef.Methods = p.extractInterfaceMembers(code, node)
		typeDef.Extends = p.extractExtends(code, node)

	case "type_alias_declaration":
		typeDef.Kind = TypeKindAlias
		valueNode := findChildByFieldName(node, "value")
		if valueNode != nil {
			typeDef.AliasOf = nodeText(code, valueNode)
		}

	case "class_declaration":
		typeDef.Kind = TypeKindClass
		typeDef.Fields, typeDef.Methods = p.extractClassMembers(code, node)
		typeDef.Extends = p.extractExtends(code, node)
		typeDef.Implements = p.extractImplements(code, node)

	case "enum_declaration":
		typeDef.Kind = TypeKindEnum
		typeDef.Variants = p.extractEnumMembers(code, node)
	}

	return typeDef
}

// extractInterfaceMembers extracts fields and methods from an interface
func (p *TypeScriptParser) extractInterfaceMembers(code []byte, node *sitter.Node) ([]Field, []string) {
	var fields []Field
	var methods []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return fields, methods
	}

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "property_signature":
			field := p.parsePropertySignature(code, n)
			if field != nil {
				fields = append(fields, *field)
			}
			return false
		case "method_signature":
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				methods = append(methods, nodeText(code, nameNode))
			}
			return false
		}
		return true
	})

	return fields, methods
}

// parsePropertySignature extracts a property from an interface
func (p *TypeScriptParser) parsePropertySignature(code []byte, node *sitter.Node) *Field {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	field := &Field{
		Name: nodeText(code, nameNode),
	}

	// Check for optional
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "?" {
			field.IsOptional = true
			break
		}
	}

	// Extract type
	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		field.Type = p.extractTypeAnnotation(code, typeNode)
	}

	return field
}

// extractClassMembers extracts fields and methods from a class
func (p *TypeScriptParser) extractClassMembers(code []byte, node *sitter.Node) ([]Field, []string) {
	var fields []Field
	var methods []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return fields, methods
	}

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "public_field_definition", "field_definition":
			field := p.parseFieldDefinition(code, n)
			if field != nil {
				fields = append(fields, *field)
			}
			return false
		case "method_definition":
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				methods = append(methods, nodeText(code, nameNode))
			}
			return false
		}
		return true
	})

	return fields, methods
}

// parseFieldDefinition extracts a field from a class
func (p *TypeScriptParser) parseFieldDefinition(code []byte, node *sitter.Node) *Field {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	field := &Field{
		Name: nodeText(code, nameNode),
	}

	// Check for readonly
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "readonly" {
			field.IsReadonly = true
			break
		}
	}

	// Extract type
	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		field.Type = p.extractTypeAnnotation(code, typeNode)
	}

	return field
}

// extractEnumMembers extracts enum variants
func (p *TypeScriptParser) extractEnumMembers(code []byte, node *sitter.Node) []string {
	var variants []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return variants
	}

	walkTree(bodyNode, func(n *sitter.Node) bool {
		if n.Type() == "enum_assignment" || n.Type() == "property_identifier" {
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				variants = append(variants, nodeText(code, nameNode))
			} else if n.Type() == "property_identifier" {
				variants = append(variants, nodeText(code, n))
			}
		}
		return true
	})

	return variants
}

// extractExtends extracts the extended type/class
func (p *TypeScriptParser) extractExtends(code []byte, node *sitter.Node) string {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "extends_clause" || child.Type() == "extends_type_clause" {
			// Get the type after extends
			for j := 0; j < int(child.ChildCount()); j++ {
				grandchild := child.Child(j)
				if grandchild.Type() != "extends" {
					return nodeText(code, grandchild)
				}
			}
		}
	}
	return ""
}

// extractImplements extracts implemented interfaces
func (p *TypeScriptParser) extractImplements(code []byte, node *sitter.Node) []string {
	var implements []string

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "implements_clause" {
			for j := 0; j < int(child.ChildCount()); j++ {
				grandchild := child.Child(j)
				if grandchild.Type() != "implements" && grandchild.Type() != "," {
					implements = append(implements, nodeText(code, grandchild))
				}
			}
		}
	}

	return implements
}

// extractImports extracts import declarations
func (p *TypeScriptParser) extractImports(code []byte, root *sitter.Node, result *ParseResult) {
	importNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "import_statement"
	})

	for _, node := range importNodes {
		imp := p.parseImportStatement(code, node)
		if imp != nil {
			result.Imports = append(result.Imports, *imp)
		}
	}
}

// parseImportStatement extracts import information
func (p *TypeScriptParser) parseImportStatement(code []byte, node *sitter.Node) *Import {
	sourceNode := findChildByFieldName(node, "source")
	if sourceNode == nil {
		return nil
	}

	path := strings.Trim(nodeText(code, sourceNode), `"'`)
	imp := &Import{
		Path:    path,
		IsLocal: strings.HasPrefix(path, ".") || strings.HasPrefix(path, "/"),
	}

	// Extract named imports
	clauseNode := findChildByType(node, "import_clause")
	if clauseNode != nil {
		namedImports := findChildByType(clauseNode, "named_imports")
		if namedImports != nil {
			walkTree(namedImports, func(n *sitter.Node) bool {
				if n.Type() == "import_specifier" {
					nameNode := findChildByFieldName(n, "name")
					if nameNode != nil {
						imp.Items = append(imp.Items, nodeText(code, nameNode))
					}
				}
				return true
			})
		}

		// Default import
		for i := 0; i < int(clauseNode.ChildCount()); i++ {
			child := clauseNode.Child(i)
			if child.Type() == "identifier" {
				imp.Alias = nodeText(code, child)
				break
			}
		}
	}

	return imp
}
