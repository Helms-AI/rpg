package treesitter

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"
)

// RustParser implements LanguageParser for Rust
type RustParser struct {
	baseParser
}

// NewRustParser creates a new Rust parser
func NewRustParser() *RustParser {
	return &RustParser{
		baseParser: baseParser{
			lang:       LanguageRust,
			extensions: []string{".rs"},
			tsLang:     rust.GetLanguage(),
		},
	}
}

// Parse parses Rust source code
func (p *RustParser) Parse(code []byte, filename string) (*ParseResult, error) {
	tree, err := p.parseTree(code)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()
	result := &ParseResult{
		Language: LanguageRust,
		FileName: filename,
	}

	p.extractFunctions(code, root, filename, result)
	p.extractTypes(code, root, filename, result)
	p.extractImports(code, root, result)

	return result, nil
}

// extractFunctions extracts function definitions
func (p *RustParser) extractFunctions(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	funcNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "function_item"
	})

	for _, node := range funcNodes {
		fn := p.parseFunctionNode(code, node, filename, root)
		if fn != nil {
			result.Functions = append(result.Functions, *fn)
		}
	}
}

// parseFunctionNode extracts function information
func (p *RustParser) parseFunctionNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *FunctionDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	fn := &FunctionDef{
		Name:     name,
		IsPublic: p.hasVisibility(code, node, "pub"),
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
	if paramsNode != nil {
		fn.Parameters = p.parseParameters(code, paramsNode)
	}

	// Extract return type
	returnTypeNode := findChildByFieldName(node, "return_type")
	if returnTypeNode != nil {
		fn.ReturnType = p.extractReturnType(code, returnTypeNode)
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
	fn.DocComment = p.extractDocComment(code, node, root)

	return fn
}

// hasVisibility checks if a node has a specific visibility modifier
func (p *RustParser) hasVisibility(code []byte, node *sitter.Node, visibility string) bool {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "visibility_modifier" {
			text := nodeText(code, child)
			return strings.HasPrefix(text, visibility)
		}
	}
	return false
}

// parseParameters extracts parameters from a parameters node
func (p *RustParser) parseParameters(code []byte, paramsNode *sitter.Node) []Parameter {
	var params []Parameter

	walkTree(paramsNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "parameter":
			param := p.parseParameter(code, n)
			if param != nil {
				params = append(params, *param)
			}
			return false
		case "self_parameter":
			// Skip self parameter but note it exists
			return false
		}
		return true
	})

	return params
}

// parseParameter extracts a single parameter
func (p *RustParser) parseParameter(code []byte, node *sitter.Node) *Parameter {
	param := &Parameter{}

	patternNode := findChildByFieldName(node, "pattern")
	if patternNode != nil {
		param.Name = nodeText(code, patternNode)
	}

	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		param.Type = nodeText(code, typeNode)
	}

	return param
}

// extractReturnType extracts the return type
func (p *RustParser) extractReturnType(code []byte, node *sitter.Node) string {
	// Skip the "->" token
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() != "->" {
			return nodeText(code, child)
		}
	}
	return nodeText(code, node)
}

// buildSignature builds a function signature string
func (p *RustParser) buildSignature(fn *FunctionDef) string {
	var sb strings.Builder

	if fn.IsPublic {
		sb.WriteString("pub ")
	}
	if fn.IsAsync {
		sb.WriteString("async ")
	}
	sb.WriteString("fn ")
	sb.WriteString(fn.Name)
	sb.WriteString("(")

	for i, param := range fn.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
		sb.WriteString(": ")
		sb.WriteString(param.Type)
	}
	sb.WriteString(")")

	if fn.ReturnType != "" {
		sb.WriteString(" -> ")
		sb.WriteString(fn.ReturnType)
	}

	return sb.String()
}

// extractCallsFromBody extracts function calls from a function body
func (p *RustParser) extractCallsFromBody(code []byte, bodyNode *sitter.Node) []string {
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
func (p *RustParser) calculateComplexity(bodyNode *sitter.Node) int {
	complexity := 1

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "if_expression", "if_let_expression", "for_expression",
			"while_expression", "while_let_expression", "loop_expression",
			"match_expression", "match_arm", "&&", "||":
			complexity++
		}
		return true
	})

	return complexity
}

// extractDocComment extracts /// or //! doc comments
func (p *RustParser) extractDocComment(code []byte, node *sitter.Node, root *sitter.Node) string {
	comment := getCommentAbove(code, node, root)
	// Rust doc comments start with /// or //!
	lines := strings.Split(comment, "\n")
	var docLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "/") {
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimPrefix(line, "/")
			line = strings.TrimPrefix(line, "!")
			docLines = append(docLines, strings.TrimSpace(line))
		}
	}
	return strings.Join(docLines, "\n")
}

// extractTypes extracts type definitions
func (p *RustParser) extractTypes(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	typeNodes := collectNodes(root, func(n *sitter.Node) bool {
		switch n.Type() {
		case "struct_item", "enum_item", "trait_item", "type_item", "impl_item":
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

// parseTypeNode extracts type information
func (p *RustParser) parseTypeNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *TypeDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil && node.Type() != "impl_item" {
		return nil
	}

	var name string
	if nameNode != nil {
		name = nodeText(code, nameNode)
	}

	typeDef := &TypeDef{
		Name:     name,
		IsPublic: p.hasVisibility(code, node, "pub"),
		Location: nodeLocation(filename, node),
		ASTHash:  hashNode(code, node),
	}

	switch node.Type() {
	case "struct_item":
		typeDef.Kind = TypeKindStruct
		typeDef.Fields = p.extractStructFields(code, node)
		typeDef.Generic = p.extractTypeParameters(code, node)

	case "enum_item":
		typeDef.Kind = TypeKindEnum
		typeDef.Variants = p.extractEnumVariants(code, node)
		typeDef.Generic = p.extractTypeParameters(code, node)

	case "trait_item":
		typeDef.Kind = TypeKindInterface
		typeDef.Methods = p.extractTraitMethods(code, node)
		typeDef.Generic = p.extractTypeParameters(code, node)

	case "type_item":
		typeDef.Kind = TypeKindAlias
		typeNode := findChildByFieldName(node, "type")
		if typeNode != nil {
			typeDef.AliasOf = nodeText(code, typeNode)
		}

	case "impl_item":
		// Skip impl blocks for now - they don't define new types
		return nil
	}

	typeDef.DocComment = p.extractDocComment(code, node, root)

	return typeDef
}

// extractStructFields extracts fields from a struct
func (p *RustParser) extractStructFields(code []byte, node *sitter.Node) []Field {
	var fields []Field

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return fields
	}

	fieldNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "field_declaration"
	})

	for _, fieldNode := range fieldNodes {
		field := p.parseFieldDeclaration(code, fieldNode)
		if field != nil {
			fields = append(fields, *field)
		}
	}

	return fields
}

// parseFieldDeclaration extracts a single field
func (p *RustParser) parseFieldDeclaration(code []byte, node *sitter.Node) *Field {
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

	// Check for pub visibility
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "visibility_modifier" {
			// Field is public
			break
		}
	}

	return field
}

// extractEnumVariants extracts enum variant names
func (p *RustParser) extractEnumVariants(code []byte, node *sitter.Node) []string {
	var variants []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return variants
	}

	variantNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "enum_variant"
	})

	for _, variantNode := range variantNodes {
		nameNode := findChildByFieldName(variantNode, "name")
		if nameNode != nil {
			variants = append(variants, nodeText(code, nameNode))
		}
	}

	return variants
}

// extractTraitMethods extracts method names from a trait
func (p *RustParser) extractTraitMethods(code []byte, node *sitter.Node) []string {
	var methods []string

	bodyNode := findChildByFieldName(node, "body")
	if bodyNode == nil {
		return methods
	}

	funcNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "function_signature_item" || n.Type() == "function_item"
	})

	for _, funcNode := range funcNodes {
		nameNode := findChildByFieldName(funcNode, "name")
		if nameNode != nil {
			methods = append(methods, nodeText(code, nameNode))
		}
	}

	return methods
}

// extractTypeParameters extracts generic type parameters
func (p *RustParser) extractTypeParameters(code []byte, node *sitter.Node) []string {
	var params []string

	typeParamsNode := findChildByFieldName(node, "type_parameters")
	if typeParamsNode == nil {
		return params
	}

	walkTree(typeParamsNode, func(n *sitter.Node) bool {
		if n.Type() == "type_identifier" || n.Type() == "lifetime" {
			params = append(params, nodeText(code, n))
			return false
		}
		return true
	})

	return params
}

// extractImports extracts use statements
func (p *RustParser) extractImports(code []byte, root *sitter.Node, result *ParseResult) {
	useNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "use_declaration"
	})

	for _, node := range useNodes {
		imps := p.parseUseDeclaration(code, node)
		result.Imports = append(result.Imports, imps...)
	}
}

// parseUseDeclaration extracts import information from a use statement
func (p *RustParser) parseUseDeclaration(code []byte, node *sitter.Node) []Import {
	var imports []Import

	// Get the full use path
	walkTree(node, func(n *sitter.Node) bool {
		switch n.Type() {
		case "scoped_identifier", "identifier":
			path := nodeText(code, n)
			imports = append(imports, Import{
				Path:    path,
				IsLocal: strings.HasPrefix(path, "crate::") || strings.HasPrefix(path, "self::") || strings.HasPrefix(path, "super::"),
			})
			return false
		case "use_wildcard":
			// Handle glob imports
			parent := n.Parent()
			if parent != nil {
				path := nodeText(code, parent)
				imports = append(imports, Import{
					Path:  path,
					Items: []string{"*"},
				})
			}
			return false
		case "use_list":
			// Handle grouped imports: use foo::{bar, baz}
			// The path should come from the parent
			return true
		}
		return true
	})

	return imports
}
