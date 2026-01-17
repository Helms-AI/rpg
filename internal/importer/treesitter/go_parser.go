package treesitter

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// GoParser implements LanguageParser for Go
type GoParser struct {
	baseParser
}

// NewGoParser creates a new Go parser
func NewGoParser() *GoParser {
	return &GoParser{
		baseParser: baseParser{
			lang:       LanguageGo,
			extensions: []string{".go"},
			tsLang:     golang.GetLanguage(),
		},
	}
}

// Parse parses Go source code
func (p *GoParser) Parse(code []byte, filename string) (*ParseResult, error) {
	tree, err := p.parseTree(code)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()
	result := &ParseResult{
		Language: LanguageGo,
		FileName: filename,
	}

	// Extract all declarations
	p.extractFunctions(code, root, filename, result)
	p.extractTypes(code, root, filename, result)
	p.extractImports(code, root, result)
	p.extractConstants(code, root, filename, result)

	return result, nil
}

// extractFunctions extracts all function declarations
func (p *GoParser) extractFunctions(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	// Find function_declaration and method_declaration nodes
	funcNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "function_declaration" || n.Type() == "method_declaration"
	})

	for _, node := range funcNodes {
		fn := p.parseFunctionNode(code, node, filename, root)
		if fn != nil {
			result.Functions = append(result.Functions, *fn)
		}
	}
}

// parseFunctionNode extracts function information from a function_declaration node
func (p *GoParser) parseFunctionNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *FunctionDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)

	// Skip test functions and unexported functions based on use case
	fn := &FunctionDef{
		Name:     name,
		IsPublic: isExported(name, LanguageGo),
		Location: nodeLocation(filename, node),
		ASTHash:  hashNode(code, node),
	}

	// Extract receiver for methods
	receiverNode := findChildByFieldName(node, "receiver")
	if receiverNode != nil {
		// This is a method - could extract receiver type if needed
	}

	// Extract parameters
	paramsNode := findChildByFieldName(node, "parameters")
	if paramsNode != nil {
		fn.Parameters = p.parseParameters(code, paramsNode)
		fn.Signature = p.buildSignature(name, fn.Parameters, "")
	}

	// Extract return type
	resultNode := findChildByFieldName(node, "result")
	if resultNode != nil {
		fn.ReturnType = nodeText(code, resultNode)
		fn.Signature = p.buildSignature(name, fn.Parameters, fn.ReturnType)
	}

	// Extract body for call analysis
	bodyNode := findChildByFieldName(node, "body")
	if bodyNode != nil {
		fn.Body = nodeText(code, bodyNode)
		fn.Calls = p.extractCallsFromBody(code, bodyNode)
		fn.Complexity = p.calculateComplexity(bodyNode)
	}

	// Extract documentation comment
	fn.DocComment = getCommentAbove(code, node, root)

	return fn
}

// parseParameters extracts parameters from a parameter_list node
func (p *GoParser) parseParameters(code []byte, paramsNode *sitter.Node) []Parameter {
	var params []Parameter

	paramDecls := findChildrenByType(paramsNode, "parameter_declaration")
	for _, decl := range paramDecls {
		// In Go, multiple names can share a type: (a, b int)
		names := findChildrenByType(decl, "identifier")
		typeNode := findChildByType(decl, "type_identifier")
		if typeNode == nil {
			// Try other type nodes
			for i := 0; i < int(decl.ChildCount()); i++ {
				child := decl.Child(i)
				if isTypeNode(child) {
					typeNode = child
					break
				}
			}
		}

		typeName := ""
		if typeNode != nil {
			typeName = nodeText(code, typeNode)
		}

		// Check for variadic
		isVariadic := false
		if findChildByType(decl, "variadic_parameter_declaration") != nil {
			isVariadic = true
		}

		for _, nameNode := range names {
			params = append(params, Parameter{
				Name:       nodeText(code, nameNode),
				Type:       typeName,
				IsVariadic: isVariadic,
			})
		}
	}

	return params
}

// isTypeNode checks if a node represents a type
func isTypeNode(node *sitter.Node) bool {
	switch node.Type() {
	case "type_identifier", "pointer_type", "slice_type", "array_type",
		"map_type", "channel_type", "function_type", "struct_type",
		"interface_type", "qualified_type":
		return true
	}
	return false
}

// buildSignature builds a function signature string
func (p *GoParser) buildSignature(name string, params []Parameter, returnType string) string {
	var sb strings.Builder
	sb.WriteString("func ")
	sb.WriteString(name)
	sb.WriteString("(")

	for i, param := range params {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
		sb.WriteString(" ")
		sb.WriteString(param.Type)
	}
	sb.WriteString(")")

	if returnType != "" {
		sb.WriteString(" ")
		sb.WriteString(returnType)
	}

	return sb.String()
}

// extractCallsFromBody extracts function calls from a function body
func (p *GoParser) extractCallsFromBody(code []byte, bodyNode *sitter.Node) []string {
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
func (p *GoParser) calculateComplexity(bodyNode *sitter.Node) int {
	complexity := 1 // Base complexity

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "if_statement", "for_statement", "range_statement",
			"switch_statement", "case_clause", "select_statement",
			"binary_expression": // && and || add to complexity
			if n.Type() == "binary_expression" {
				// Only count && and ||
				for i := 0; i < int(n.ChildCount()); i++ {
					child := n.Child(i)
					if child.Type() == "&&" || child.Type() == "||" {
						complexity++
						break
					}
				}
			} else {
				complexity++
			}
		}
		return true
	})

	return complexity
}

// extractTypes extracts type declarations
func (p *GoParser) extractTypes(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	typeDecls := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "type_declaration"
	})

	for _, decl := range typeDecls {
		specs := findChildrenByType(decl, "type_spec")
		for _, spec := range specs {
			typeDef := p.parseTypeSpec(code, spec, filename, root)
			if typeDef != nil {
				result.Types = append(result.Types, *typeDef)
			}
		}
	}
}

// parseTypeSpec extracts type information from a type_spec node
func (p *GoParser) parseTypeSpec(code []byte, spec *sitter.Node, filename string, root *sitter.Node) *TypeDef {
	nameNode := findChildByFieldName(spec, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	typeDef := &TypeDef{
		Name:       name,
		IsPublic:   isExported(name, LanguageGo),
		Location:   nodeLocation(filename, spec),
		ASTHash:    hashNode(code, spec),
		DocComment: getCommentAbove(code, spec, root),
	}

	// Determine type kind and extract details
	typeNode := findChildByFieldName(spec, "type")
	if typeNode == nil {
		return typeDef
	}

	switch typeNode.Type() {
	case "struct_type":
		typeDef.Kind = TypeKindStruct
		typeDef.Fields = p.extractStructFields(code, typeNode)
	case "interface_type":
		typeDef.Kind = TypeKindInterface
		typeDef.Methods = p.extractInterfaceMethods(code, typeNode)
	case "type_identifier", "pointer_type", "slice_type", "array_type", "map_type":
		typeDef.Kind = TypeKindAlias
		typeDef.AliasOf = nodeText(code, typeNode)
	default:
		typeDef.Kind = TypeKindAlias
		typeDef.AliasOf = nodeText(code, typeNode)
	}

	return typeDef
}

// extractStructFields extracts fields from a struct_type node
func (p *GoParser) extractStructFields(code []byte, structNode *sitter.Node) []Field {
	var fields []Field

	fieldDeclList := findChildByType(structNode, "field_declaration_list")
	if fieldDeclList == nil {
		return fields
	}

	fieldDecls := findChildrenByType(fieldDeclList, "field_declaration")
	for _, decl := range fieldDecls {
		// Extract field names and type
		var names []string
		var fieldType string
		var tags string

		for i := 0; i < int(decl.ChildCount()); i++ {
			child := decl.Child(i)
			switch child.Type() {
			case "field_identifier":
				names = append(names, nodeText(code, child))
			case "raw_string_literal", "interpreted_string_literal":
				tags = nodeText(code, child)
			default:
				if isTypeNode(child) {
					fieldType = nodeText(code, child)
				}
			}
		}

		// Handle embedded fields (no name, just type)
		if len(names) == 0 && fieldType != "" {
			names = []string{fieldType}
		}

		for _, name := range names {
			fields = append(fields, Field{
				Name: name,
				Type: fieldType,
				Tags: tags,
			})
		}
	}

	return fields
}

// extractInterfaceMethods extracts method signatures from an interface_type node
func (p *GoParser) extractInterfaceMethods(code []byte, ifaceNode *sitter.Node) []string {
	var methods []string

	// Find method_spec nodes inside interface
	walkTree(ifaceNode, func(n *sitter.Node) bool {
		if n.Type() == "method_spec" {
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				methods = append(methods, nodeText(code, nameNode))
			}
		}
		return true
	})

	return methods
}

// extractImports extracts import declarations
func (p *GoParser) extractImports(code []byte, root *sitter.Node, result *ParseResult) {
	importDecls := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "import_declaration"
	})

	for _, decl := range importDecls {
		// Handle import specs
		specs := collectNodes(decl, func(n *sitter.Node) bool {
			return n.Type() == "import_spec"
		})

		for _, spec := range specs {
			imp := p.parseImportSpec(code, spec)
			if imp != nil {
				result.Imports = append(result.Imports, *imp)
			}
		}
	}
}

// parseImportSpec extracts import information from an import_spec node
func (p *GoParser) parseImportSpec(code []byte, spec *sitter.Node) *Import {
	pathNode := findChildByFieldName(spec, "path")
	if pathNode == nil {
		return nil
	}

	path := strings.Trim(nodeText(code, pathNode), `"`)

	imp := &Import{
		Path:    path,
		IsLocal: !strings.Contains(path, ".") || strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../"),
	}

	// Check for alias
	nameNode := findChildByFieldName(spec, "name")
	if nameNode != nil {
		imp.Alias = nodeText(code, nameNode)
	}

	return imp
}

// extractConstants extracts constant declarations
func (p *GoParser) extractConstants(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	constDecls := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "const_declaration"
	})

	for _, decl := range constDecls {
		specs := findChildrenByType(decl, "const_spec")
		for _, spec := range specs {
			consts := p.parseConstSpec(code, spec, filename, root)
			result.Constants = append(result.Constants, consts...)
		}
	}
}

// parseConstSpec extracts constant information from a const_spec node
func (p *GoParser) parseConstSpec(code []byte, spec *sitter.Node, filename string, root *sitter.Node) []Constant {
	var consts []Constant

	// Collect names
	var names []*sitter.Node
	var typeName string
	var value string

	for i := 0; i < int(spec.ChildCount()); i++ {
		child := spec.Child(i)
		switch child.Type() {
		case "identifier":
			names = append(names, child)
		case "type_identifier":
			typeName = nodeText(code, child)
		case "expression_list":
			value = nodeText(code, child)
		default:
			// Could be a value expression
			if len(names) > 0 && value == "" {
				value = nodeText(code, child)
			}
		}
	}

	for _, nameNode := range names {
		name := nodeText(code, nameNode)
		consts = append(consts, Constant{
			Name:       name,
			Type:       typeName,
			Value:      value,
			IsPublic:   isExported(name, LanguageGo),
			Location:   nodeLocation(filename, spec),
			DocComment: getCommentAbove(code, spec, root),
		})
	}

	return consts
}
