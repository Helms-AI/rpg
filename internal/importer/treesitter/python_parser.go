package treesitter

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

// PythonParser implements LanguageParser for Python
type PythonParser struct {
	baseParser
}

// NewPythonParser creates a new Python parser
func NewPythonParser() *PythonParser {
	return &PythonParser{
		baseParser: baseParser{
			lang:       LanguagePython,
			extensions: []string{".py"},
			tsLang:     python.GetLanguage(),
		},
	}
}

// Parse parses Python source code
func (p *PythonParser) Parse(code []byte, filename string) (*ParseResult, error) {
	tree, err := p.parseTree(code)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	root := tree.RootNode()
	result := &ParseResult{
		Language: LanguagePython,
		FileName: filename,
	}

	p.extractFunctions(code, root, filename, result)
	p.extractTypes(code, root, filename, result)
	p.extractImports(code, root, result)

	return result, nil
}

// extractFunctions extracts all function definitions
func (p *PythonParser) extractFunctions(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	funcNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "function_definition"
	})

	for _, node := range funcNodes {
		fn := p.parseFunctionNode(code, node, filename, root)
		if fn != nil {
			result.Functions = append(result.Functions, *fn)
		}
	}
}

// parseFunctionNode extracts function information
func (p *PythonParser) parseFunctionNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *FunctionDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
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
	if paramsNode != nil {
		fn.Parameters = p.parseParameters(code, paramsNode)
	}

	// Extract return type annotation
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

		// Extract docstring
		fn.DocComment = p.extractDocstring(code, bodyNode)
	}

	return fn
}

// parseParameters extracts parameters from a parameters node
func (p *PythonParser) parseParameters(code []byte, paramsNode *sitter.Node) []Parameter {
	var params []Parameter

	walkTree(paramsNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "identifier":
			// Simple parameter
			name := nodeText(code, n)
			if name != "self" && name != "cls" { // Skip self/cls
				params = append(params, Parameter{Name: name})
			}
			return false
		case "typed_parameter":
			param := p.parseTypedParameter(code, n)
			if param != nil && param.Name != "self" && param.Name != "cls" {
				params = append(params, *param)
			}
			return false
		case "default_parameter":
			param := p.parseDefaultParameter(code, n)
			if param != nil && param.Name != "self" && param.Name != "cls" {
				params = append(params, *param)
			}
			return false
		case "typed_default_parameter":
			param := p.parseTypedDefaultParameter(code, n)
			if param != nil && param.Name != "self" && param.Name != "cls" {
				params = append(params, *param)
			}
			return false
		case "list_splat_pattern", "dictionary_splat_pattern":
			param := p.parseSplatParameter(code, n)
			if param != nil {
				params = append(params, *param)
			}
			return false
		}
		return true
	})

	return params
}

// parseTypedParameter extracts a typed parameter
func (p *PythonParser) parseTypedParameter(code []byte, node *sitter.Node) *Parameter {
	nameNode := findChildByType(node, "identifier")
	if nameNode == nil {
		return nil
	}

	param := &Parameter{
		Name: nodeText(code, nameNode),
	}

	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		param.Type = p.extractTypeAnnotation(code, typeNode)
	}

	return param
}

// parseDefaultParameter extracts a parameter with default value
func (p *PythonParser) parseDefaultParameter(code []byte, node *sitter.Node) *Parameter {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	param := &Parameter{
		Name:       nodeText(code, nameNode),
		IsOptional: true,
	}

	valueNode := findChildByFieldName(node, "value")
	if valueNode != nil {
		param.DefaultValue = nodeText(code, valueNode)
	}

	return param
}

// parseTypedDefaultParameter extracts a typed parameter with default
func (p *PythonParser) parseTypedDefaultParameter(code []byte, node *sitter.Node) *Parameter {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	param := &Parameter{
		Name:       nodeText(code, nameNode),
		IsOptional: true,
	}

	typeNode := findChildByFieldName(node, "type")
	if typeNode != nil {
		param.Type = p.extractTypeAnnotation(code, typeNode)
	}

	valueNode := findChildByFieldName(node, "value")
	if valueNode != nil {
		param.DefaultValue = nodeText(code, valueNode)
	}

	return param
}

// parseSplatParameter extracts *args or **kwargs
func (p *PythonParser) parseSplatParameter(code []byte, node *sitter.Node) *Parameter {
	nameNode := findChildByType(node, "identifier")
	if nameNode == nil {
		return nil
	}

	param := &Parameter{
		Name:       nodeText(code, nameNode),
		IsVariadic: true,
	}

	if node.Type() == "dictionary_splat_pattern" {
		param.Name = "**" + param.Name
	} else {
		param.Name = "*" + param.Name
	}

	return param
}

// extractTypeAnnotation extracts a type from a type annotation node
func (p *PythonParser) extractTypeAnnotation(code []byte, node *sitter.Node) string {
	// Skip the "->" or ":" if present at the start
	text := nodeText(code, node)
	text = strings.TrimPrefix(text, "->")
	text = strings.TrimPrefix(text, ":")
	return strings.TrimSpace(text)
}

// buildSignature builds a function signature string
func (p *PythonParser) buildSignature(fn *FunctionDef) string {
	var sb strings.Builder

	if fn.IsAsync {
		sb.WriteString("async ")
	}
	sb.WriteString("def ")
	sb.WriteString(fn.Name)
	sb.WriteString("(")

	for i, param := range fn.Parameters {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(param.Name)
		if param.Type != "" {
			sb.WriteString(": ")
			sb.WriteString(param.Type)
		}
		if param.DefaultValue != "" {
			sb.WriteString(" = ")
			sb.WriteString(param.DefaultValue)
		}
	}
	sb.WriteString(")")

	if fn.ReturnType != "" {
		sb.WriteString(" -> ")
		sb.WriteString(fn.ReturnType)
	}

	return sb.String()
}

// extractDocstring extracts a docstring from a function body
func (p *PythonParser) extractDocstring(code []byte, bodyNode *sitter.Node) string {
	// First statement in body might be a string (docstring)
	if bodyNode.ChildCount() == 0 {
		return ""
	}

	firstStmt := bodyNode.Child(0)
	if firstStmt == nil {
		return ""
	}

	// Check if it's an expression statement containing a string
	if firstStmt.Type() == "expression_statement" {
		if firstStmt.ChildCount() > 0 {
			child := firstStmt.Child(0)
			if child.Type() == "string" {
				docstring := nodeText(code, child)
				// Clean up quotes
				docstring = strings.Trim(docstring, `"'`)
				docstring = strings.TrimPrefix(docstring, `""`)
				docstring = strings.TrimSuffix(docstring, `""`)
				return strings.TrimSpace(docstring)
			}
		}
	}

	return ""
}

// extractCallsFromBody extracts function calls from a function body
func (p *PythonParser) extractCallsFromBody(code []byte, bodyNode *sitter.Node) []string {
	calls := make(map[string]bool)

	callNodes := collectNodes(bodyNode, func(n *sitter.Node) bool {
		return n.Type() == "call"
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
func (p *PythonParser) calculateComplexity(bodyNode *sitter.Node) int {
	complexity := 1

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "if_statement", "elif_clause", "for_statement",
			"while_statement", "try_statement", "except_clause",
			"with_statement", "match_statement", "case_clause",
			"conditional_expression", "and", "or":
			complexity++
		}
		return true
	})

	return complexity
}

// extractTypes extracts class definitions
func (p *PythonParser) extractTypes(code []byte, root *sitter.Node, filename string, result *ParseResult) {
	classNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "class_definition"
	})

	for _, node := range classNodes {
		typeDef := p.parseClassNode(code, node, filename, root)
		if typeDef != nil {
			result.Types = append(result.Types, *typeDef)
		}
	}
}

// parseClassNode extracts class information
func (p *PythonParser) parseClassNode(code []byte, node *sitter.Node, filename string, root *sitter.Node) *TypeDef {
	nameNode := findChildByFieldName(node, "name")
	if nameNode == nil {
		return nil
	}

	name := nodeText(code, nameNode)
	typeDef := &TypeDef{
		Name:     name,
		Kind:     TypeKindClass,
		IsPublic: !strings.HasPrefix(name, "_"),
		Location: nodeLocation(filename, node),
		ASTHash:  hashNode(code, node),
	}

	// Extract base classes
	superclassNode := findChildByFieldName(node, "superclasses")
	if superclassNode != nil {
		bases := p.extractSuperclasses(code, superclassNode)
		if len(bases) > 0 {
			typeDef.Extends = bases[0]
			if len(bases) > 1 {
				typeDef.Implements = bases[1:]
			}
		}
	}

	// Extract body
	bodyNode := findChildByFieldName(node, "body")
	if bodyNode != nil {
		typeDef.DocComment = p.extractDocstring(code, bodyNode)
		typeDef.Fields, typeDef.Methods = p.extractClassMembers(code, bodyNode)
	}

	return typeDef
}

// extractSuperclasses extracts base classes
func (p *PythonParser) extractSuperclasses(code []byte, node *sitter.Node) []string {
	var bases []string

	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "identifier" || n.Type() == "attribute" {
			bases = append(bases, nodeText(code, n))
			return false
		}
		return true
	})

	return bases
}

// extractClassMembers extracts fields and methods from a class body
func (p *PythonParser) extractClassMembers(code []byte, bodyNode *sitter.Node) ([]Field, []string) {
	var fields []Field
	var methods []string
	fieldsSeen := make(map[string]bool)

	walkTree(bodyNode, func(n *sitter.Node) bool {
		switch n.Type() {
		case "function_definition":
			nameNode := findChildByFieldName(n, "name")
			if nameNode != nil {
				methods = append(methods, nodeText(code, nameNode))
			}

			// Extract fields from __init__
			name := nodeText(code, nameNode)
			if name == "__init__" {
				initFields := p.extractFieldsFromInit(code, n)
				for _, f := range initFields {
					if !fieldsSeen[f.Name] {
						fields = append(fields, f)
						fieldsSeen[f.Name] = true
					}
				}
			}
			return false

		case "expression_statement":
			// Class-level field assignment
			field := p.extractFieldFromAssignment(code, n)
			if field != nil && !fieldsSeen[field.Name] {
				fields = append(fields, *field)
				fieldsSeen[field.Name] = true
			}
			return false
		}
		return true
	})

	return fields, methods
}

// extractFieldsFromInit extracts self.x assignments from __init__
func (p *PythonParser) extractFieldsFromInit(code []byte, funcNode *sitter.Node) []Field {
	var fields []Field
	seen := make(map[string]bool)

	bodyNode := findChildByFieldName(funcNode, "body")
	if bodyNode == nil {
		return fields
	}

	walkTree(bodyNode, func(n *sitter.Node) bool {
		if n.Type() == "assignment" {
			leftNode := findChildByFieldName(n, "left")
			if leftNode != nil && leftNode.Type() == "attribute" {
				// Check if it's self.something
				objNode := findChildByFieldName(leftNode, "object")
				if objNode != nil && nodeText(code, objNode) == "self" {
					attrNode := findChildByFieldName(leftNode, "attribute")
					if attrNode != nil {
						name := nodeText(code, attrNode)
						if !seen[name] {
							fields = append(fields, Field{Name: name})
							seen[name] = true
						}
					}
				}
			}
		}
		return true
	})

	return fields
}

// extractFieldFromAssignment extracts a class-level field from an assignment
func (p *PythonParser) extractFieldFromAssignment(code []byte, node *sitter.Node) *Field {
	// Look for assignment or annotated assignment
	walkTree(node, func(n *sitter.Node) bool {
		if n.Type() == "assignment" || n.Type() == "annotated_assignment" {
			leftNode := findChildByFieldName(n, "left")
			if leftNode != nil && leftNode.Type() == "identifier" {
				field := &Field{
					Name: nodeText(code, leftNode),
				}

				// Get type annotation if present
				typeNode := findChildByFieldName(n, "type")
				if typeNode != nil {
					field.Type = p.extractTypeAnnotation(code, typeNode)
				}

				// Get value if present
				rightNode := findChildByFieldName(n, "right")
				if rightNode != nil {
					field.Default = nodeText(code, rightNode)
				}

				return false
			}
		}
		return true
	})

	return nil
}

// extractImports extracts import statements
func (p *PythonParser) extractImports(code []byte, root *sitter.Node, result *ParseResult) {
	importNodes := collectNodes(root, func(n *sitter.Node) bool {
		return n.Type() == "import_statement" || n.Type() == "import_from_statement"
	})

	for _, node := range importNodes {
		imps := p.parseImportStatement(code, node)
		result.Imports = append(result.Imports, imps...)
	}
}

// parseImportStatement extracts import information
func (p *PythonParser) parseImportStatement(code []byte, node *sitter.Node) []Import {
	var imports []Import

	if node.Type() == "import_statement" {
		// import foo, bar
		walkTree(node, func(n *sitter.Node) bool {
			if n.Type() == "dotted_name" {
				path := nodeText(code, n)
				imports = append(imports, Import{
					Path:    path,
					IsLocal: strings.HasPrefix(path, "."),
				})
				return false
			}
			if n.Type() == "aliased_import" {
				nameNode := findChildByFieldName(n, "name")
				aliasNode := findChildByFieldName(n, "alias")
				if nameNode != nil {
					imp := Import{
						Path:    nodeText(code, nameNode),
						IsLocal: false,
					}
					if aliasNode != nil {
						imp.Alias = nodeText(code, aliasNode)
					}
					imports = append(imports, imp)
				}
				return false
			}
			return true
		})
	} else {
		// from foo import bar, baz
		moduleNode := findChildByFieldName(node, "module_name")
		modulePath := ""
		if moduleNode != nil {
			modulePath = nodeText(code, moduleNode)
		}

		// Check for relative imports (leading dots)
		var prefix string
		for i := 0; i < int(node.ChildCount()); i++ {
			child := node.Child(i)
			if child.Type() == "relative_import" {
				prefix = nodeText(code, child)
				break
			}
		}

		fullPath := prefix + modulePath
		imp := Import{
			Path:    fullPath,
			IsLocal: strings.HasPrefix(fullPath, "."),
		}

		// Extract imported names
		walkTree(node, func(n *sitter.Node) bool {
			if n.Type() == "dotted_name" && n != moduleNode {
				imp.Items = append(imp.Items, nodeText(code, n))
				return false
			}
			if n.Type() == "aliased_import" {
				nameNode := findChildByFieldName(n, "name")
				if nameNode != nil {
					imp.Items = append(imp.Items, nodeText(code, nameNode))
				}
				return false
			}
			return true
		})

		imports = append(imports, imp)
	}

	return imports
}
