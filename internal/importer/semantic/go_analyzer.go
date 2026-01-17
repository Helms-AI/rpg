package semantic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/kon1790/rpg/internal/importer/treesitter"
)

// GoAnalyzer provides semantic analysis for Go code
type GoAnalyzer struct {
	fset *token.FileSet
}

// NewGoAnalyzer creates a new Go semantic analyzer
func NewGoAnalyzer() *GoAnalyzer {
	return &GoAnalyzer{
		fset: token.NewFileSet(),
	}
}

// Language returns the language this analyzer handles
func (a *GoAnalyzer) Language() treesitter.Language {
	return treesitter.LanguageGo
}

// IsAvailable checks if Go tools are available (always true since we use stdlib)
func (a *GoAnalyzer) IsAvailable() bool {
	return true
}

// Analyze performs semantic analysis on a Go project directory
func (a *GoAnalyzer) Analyze(dir string) (*Analysis, error) {
	a.fset = token.NewFileSet()

	analysis := &Analysis{
		Language:  treesitter.LanguageGo,
		CallGraph: make(map[string][]string),
		TypeGraph: make(map[string][]string),
	}

	// Find project name from go.mod or directory name
	analysis.Name = a.findProjectName(dir)

	// Collect all directories containing Go files
	goDirs := make(map[string]bool)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip inaccessible paths
		}
		if info.IsDir() {
			// Skip vendor and hidden directories
			name := info.Name()
			if name == "vendor" || name == "testdata" || (len(name) > 0 && name[0] == '.') {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			goDirs[filepath.Dir(path)] = true
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	// Parse all Go files in each directory
	for goDir := range goDirs {
		pkgs, err := parser.ParseDir(a.fset, goDir, func(fi os.FileInfo) bool {
			return !strings.HasSuffix(fi.Name(), "_test.go")
		}, parser.ParseComments)
		if err != nil {
			analysis.Errors = append(analysis.Errors, AnalysisError{
				File:     goDir,
				Message:  fmt.Sprintf("parsing directory: %v", err),
				Severity: SeverityWarning,
			})
			continue
		}

		// Process each package
		for pkgName, pkg := range pkgs {
			for filename, file := range pkg.Files {
				fileAnalysis, err := a.analyzeASTFile(filename, file, pkgName)
				if err != nil {
					analysis.Errors = append(analysis.Errors, AnalysisError{
						File:     filename,
						Message:  err.Error(),
						Severity: SeverityWarning,
					})
					continue
				}

				analysis.Files = append(analysis.Files, fileAnalysis)

				// Aggregate types and functions
				analysis.Types = append(analysis.Types, fileAnalysis.Types...)
				analysis.Functions = append(analysis.Functions, fileAnalysis.Functions...)
			}
		}
	}

	// Build call graph
	a.buildCallGraph(analysis)

	// Build type graph
	a.buildTypeGraph(analysis)

	// Extract dependencies
	a.extractDependencies(analysis)

	return analysis, nil
}

// AnalyzeFile performs semantic analysis on a single Go file
func (a *GoAnalyzer) AnalyzeFile(path string, content []byte) (*FileAnalysis, error) {
	a.fset = token.NewFileSet()

	file, err := parser.ParseFile(a.fset, path, content, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	return a.analyzeASTFile(path, file, file.Name.Name)
}

// analyzeASTFile analyzes a parsed Go AST file
func (a *GoAnalyzer) analyzeASTFile(path string, file *ast.File, pkgName string) (*FileAnalysis, error) {
	analysis := &FileAnalysis{
		Path:    path,
		Package: pkgName,
	}

	// Extract imports
	for _, imp := range file.Imports {
		analysis.Imports = append(analysis.Imports, a.extractImport(imp))
	}

	// Extract types and functions
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					typ := a.extractType(s, path, d.Doc)
					analysis.Types = append(analysis.Types, typ)
				}
			}
		case *ast.FuncDecl:
			fn := a.extractFunction(d, path)
			analysis.Functions = append(analysis.Functions, fn)
		}
	}

	return analysis, nil
}

// extractImport extracts import information
func (a *GoAnalyzer) extractImport(imp *ast.ImportSpec) ResolvedImport {
	path := strings.Trim(imp.Path.Value, `"`)

	ri := ResolvedImport{
		Import: treesitter.Import{
			Path:    path,
			IsLocal: strings.HasPrefix(path, ".") || !strings.Contains(path, "."),
		},
	}

	if imp.Name != nil {
		ri.Alias = imp.Name.Name
		ri.PackageName = imp.Name.Name
	} else {
		// Package name is usually the last part of the path
		parts := strings.Split(path, "/")
		ri.PackageName = parts[len(parts)-1]
	}

	return ri
}

// extractType extracts type information from a TypeSpec
func (a *GoAnalyzer) extractType(spec *ast.TypeSpec, path string, doc *ast.CommentGroup) ResolvedType {
	pos := a.fset.Position(spec.Pos())

	rt := ResolvedType{
		TypeDef: treesitter.TypeDef{
			Name:     spec.Name.Name,
			IsPublic: ast.IsExported(spec.Name.Name),
			Location: treesitter.SourceLocation{
				File:      path,
				StartLine: pos.Line,
				EndLine:   pos.Line,
			},
		},
	}

	// Extract doc comment
	if doc != nil {
		rt.DocComment = doc.Text()
	} else if spec.Doc != nil {
		rt.DocComment = spec.Doc.Text()
	}

	// Determine type kind and extract details
	switch t := spec.Type.(type) {
	case *ast.StructType:
		rt.Kind = treesitter.TypeKindStruct
		rt.ResolvedFields = a.extractStructFields(t)
	case *ast.InterfaceType:
		rt.Kind = treesitter.TypeKindInterface
		rt.Methods, rt.MethodSignatures = a.extractInterfaceMethods(t)
	case *ast.Ident:
		rt.Kind = treesitter.TypeKindAlias
		rt.AliasOf = t.Name
	case *ast.SelectorExpr:
		rt.Kind = treesitter.TypeKindAlias
		rt.AliasOf = types.ExprString(t)
	default:
		rt.Kind = treesitter.TypeKindAlias
		rt.AliasOf = types.ExprString(spec.Type)
	}

	// Extract type parameters (generics)
	if spec.TypeParams != nil {
		for _, param := range spec.TypeParams.List {
			for _, name := range param.Names {
				rt.Generic = append(rt.Generic, name.Name)
			}
		}
	}

	return rt
}

// extractStructFields extracts fields from a struct type
func (a *GoAnalyzer) extractStructFields(st *ast.StructType) []ResolvedField {
	var fields []ResolvedField

	if st.Fields == nil {
		return fields
	}

	for _, field := range st.Fields.List {
		typeStr := types.ExprString(field.Type)

		rf := ResolvedField{
			ResolvedType: typeStr,
		}

		// Analyze type structure
		rf.IsPointer, rf.IsSlice, rf.IsMap, rf.ElementType, rf.KeyType = analyzeType(field.Type)

		// Get tag if present
		if field.Tag != nil {
			rf.Tags = strings.Trim(field.Tag.Value, "`")
		}

		// Handle named and embedded fields
		if len(field.Names) == 0 {
			// Embedded field
			rf.Name = typeStr
			fields = append(fields, rf)
		} else {
			for _, name := range field.Names {
				f := rf // Copy
				f.Name = name.Name
				fields = append(fields, f)
			}
		}
	}

	return fields
}

// analyzeType analyzes a type expression
func analyzeType(expr ast.Expr) (isPointer, isSlice, isMap bool, elementType, keyType string) {
	switch t := expr.(type) {
	case *ast.StarExpr:
		isPointer = true
		elementType = types.ExprString(t.X)
	case *ast.ArrayType:
		isSlice = true
		elementType = types.ExprString(t.Elt)
	case *ast.MapType:
		isMap = true
		keyType = types.ExprString(t.Key)
		elementType = types.ExprString(t.Value)
	}
	return
}

// extractInterfaceMethods extracts methods from an interface type
func (a *GoAnalyzer) extractInterfaceMethods(it *ast.InterfaceType) ([]string, []MethodSignature) {
	var methods []string
	var sigs []MethodSignature

	if it.Methods == nil {
		return methods, sigs
	}

	for _, method := range it.Methods.List {
		if len(method.Names) == 0 {
			// Embedded interface
			continue
		}

		for _, name := range method.Names {
			methods = append(methods, name.Name)

			sig := MethodSignature{
				Name:       name.Name,
				IsExported: ast.IsExported(name.Name),
			}

			// Extract function signature
			if ft, ok := method.Type.(*ast.FuncType); ok {
				sig.Parameters = a.extractParameters(ft.Params)
				sig.ReturnTypes = a.extractReturnTypes(ft.Results)
			}

			sigs = append(sigs, sig)
		}
	}

	return methods, sigs
}

// extractFunction extracts function information from a FuncDecl
func (a *GoAnalyzer) extractFunction(fn *ast.FuncDecl, path string) ResolvedFunction {
	pos := a.fset.Position(fn.Pos())

	rf := ResolvedFunction{
		FunctionDef: treesitter.FunctionDef{
			Name:     fn.Name.Name,
			IsPublic: ast.IsExported(fn.Name.Name),
			Location: treesitter.SourceLocation{
				File:      path,
				StartLine: pos.Line,
			},
		},
	}

	// Extract doc comment
	if fn.Doc != nil {
		rf.DocComment = fn.Doc.Text()
	}

	// Extract parameters
	if fn.Type.Params != nil {
		rf.ResolvedParameters = a.extractParameters(fn.Type.Params)
		for _, p := range rf.ResolvedParameters {
			rf.Parameters = append(rf.Parameters, p.Parameter)
		}
	}

	// Extract return types
	if fn.Type.Results != nil {
		rf.ResolvedReturnTypes = a.extractReturnTypes(fn.Type.Results)
		rf.ReturnType = strings.Join(rf.ResolvedReturnTypes, ", ")
	}

	// Build signature
	rf.Signature = a.buildFunctionSignature(fn)

	// Extract calls from body
	if fn.Body != nil {
		rf.Calls = a.extractCalls(fn.Body)
		rf.ResolvedCalls = a.resolveCalls(rf.Calls)
		rf.LocalVariables = a.extractLocalVariables(fn.Body)
		rf.Complexity = a.calculateComplexity(fn.Body)
	}

	return rf
}

// extractParameters extracts parameters from a FieldList
func (a *GoAnalyzer) extractParameters(fields *ast.FieldList) []ResolvedParameter {
	var params []ResolvedParameter

	if fields == nil {
		return params
	}

	for _, field := range fields.List {
		typeStr := types.ExprString(field.Type)

		// Check if variadic
		isVariadic := false
		if _, ok := field.Type.(*ast.Ellipsis); ok {
			isVariadic = true
		}

		if len(field.Names) == 0 {
			// Unnamed parameter
			params = append(params, ResolvedParameter{
				Parameter: treesitter.Parameter{
					Type:       typeStr,
					IsVariadic: isVariadic,
				},
				ResolvedType: typeStr,
			})
		} else {
			for _, name := range field.Names {
				params = append(params, ResolvedParameter{
					Parameter: treesitter.Parameter{
						Name:       name.Name,
						Type:       typeStr,
						IsVariadic: isVariadic,
					},
					ResolvedType: typeStr,
				})
			}
		}
	}

	return params
}

// extractReturnTypes extracts return types from a FieldList
func (a *GoAnalyzer) extractReturnTypes(fields *ast.FieldList) []string {
	var types []string

	if fields == nil {
		return types
	}

	for _, field := range fields.List {
		typeStr := a.typeString(field.Type)
		if len(field.Names) == 0 {
			types = append(types, typeStr)
		} else {
			// Named return values
			for range field.Names {
				types = append(types, typeStr)
			}
		}
	}

	return types
}

// typeString converts an ast.Expr to a string representation
func (a *GoAnalyzer) typeString(expr ast.Expr) string {
	return types.ExprString(expr)
}

// buildFunctionSignature builds a function signature string
func (a *GoAnalyzer) buildFunctionSignature(fn *ast.FuncDecl) string {
	var sb strings.Builder
	sb.WriteString("func ")

	// Add receiver if present
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		sb.WriteString("(")
		recv := fn.Recv.List[0]
		if len(recv.Names) > 0 {
			sb.WriteString(recv.Names[0].Name)
			sb.WriteString(" ")
		}
		sb.WriteString(types.ExprString(recv.Type))
		sb.WriteString(") ")
	}

	sb.WriteString(fn.Name.Name)
	sb.WriteString("(")

	// Parameters
	if fn.Type.Params != nil {
		params := make([]string, 0)
		for _, field := range fn.Type.Params.List {
			typeStr := types.ExprString(field.Type)
			if len(field.Names) == 0 {
				params = append(params, typeStr)
			} else {
				for _, name := range field.Names {
					params = append(params, name.Name+" "+typeStr)
				}
			}
		}
		sb.WriteString(strings.Join(params, ", "))
	}
	sb.WriteString(")")

	// Return types
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		sb.WriteString(" ")
		if len(fn.Type.Results.List) == 1 && len(fn.Type.Results.List[0].Names) == 0 {
			sb.WriteString(types.ExprString(fn.Type.Results.List[0].Type))
		} else {
			sb.WriteString("(")
			results := make([]string, 0)
			for _, field := range fn.Type.Results.List {
				typeStr := types.ExprString(field.Type)
				if len(field.Names) == 0 {
					results = append(results, typeStr)
				} else {
					for _, name := range field.Names {
						results = append(results, name.Name+" "+typeStr)
					}
				}
			}
			sb.WriteString(strings.Join(results, ", "))
			sb.WriteString(")")
		}
	}

	return sb.String()
}

// extractCalls extracts function calls from a block statement
func (a *GoAnalyzer) extractCalls(body *ast.BlockStmt) []string {
	calls := make(map[string]bool)

	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			switch fn := call.Fun.(type) {
			case *ast.Ident:
				calls[fn.Name] = true
			case *ast.SelectorExpr:
				calls[types.ExprString(fn)] = true
			}
		}
		return true
	})

	result := make([]string, 0, len(calls))
	for call := range calls {
		result = append(result, call)
	}
	return result
}

// resolveCalls converts call names to CallReferences
func (a *GoAnalyzer) resolveCalls(calls []string) []CallReference {
	var refs []CallReference

	for _, call := range calls {
		ref := CallReference{Name: call, ResolvedName: call}

		if strings.Contains(call, ".") {
			ref.IsMethod = true
			parts := strings.Split(call, ".")
			if len(parts) == 2 {
				ref.Package = parts[0]
				ref.ResolvedName = parts[1]
			}
		}

		refs = append(refs, ref)
	}

	return refs
}

// extractLocalVariables extracts local variable declarations
func (a *GoAnalyzer) extractLocalVariables(body *ast.BlockStmt) []Variable {
	var vars []Variable

	ast.Inspect(body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			if stmt.Tok == token.DEFINE {
				for i, lhs := range stmt.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						var typeStr string
						if i < len(stmt.Rhs) {
							typeStr = inferType(stmt.Rhs[i])
						}
						pos := a.fset.Position(ident.Pos())
						vars = append(vars, Variable{
							Name: ident.Name,
							Type: typeStr,
							Line: pos.Line,
						})
					}
				}
			}
		case *ast.DeclStmt:
			if decl, ok := stmt.Decl.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						typeStr := ""
						if vs.Type != nil {
							typeStr = types.ExprString(vs.Type)
						}
						for _, name := range vs.Names {
							pos := a.fset.Position(name.Pos())
							vars = append(vars, Variable{
								Name: name.Name,
								Type: typeStr,
								Line: pos.Line,
							})
						}
					}
				}
			}
		}
		return true
	})

	return vars
}

// inferType attempts to infer the type of an expression
func inferType(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		switch e.Kind {
		case token.INT:
			return "int"
		case token.FLOAT:
			return "float64"
		case token.STRING:
			return "string"
		case token.CHAR:
			return "rune"
		}
	case *ast.CompositeLit:
		return types.ExprString(e.Type)
	case *ast.CallExpr:
		return types.ExprString(e.Fun)
	}
	return ""
}

// calculateComplexity calculates cyclomatic complexity
func (a *GoAnalyzer) calculateComplexity(body *ast.BlockStmt) int {
	complexity := 1

	ast.Inspect(body, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt,
			*ast.CaseClause, *ast.CommClause:
			complexity++
		case *ast.BinaryExpr:
			if be, ok := n.(*ast.BinaryExpr); ok {
				if be.Op == token.LAND || be.Op == token.LOR {
					complexity++
				}
			}
		}
		return true
	})

	return complexity
}

// findProjectName finds the project name from go.mod or directory
func (a *GoAnalyzer) findProjectName(dir string) string {
	// Try to read go.mod
	goModPath := filepath.Join(dir, "go.mod")
	if data, err := os.ReadFile(goModPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				module := strings.TrimPrefix(line, "module ")
				module = strings.TrimSpace(module)
				parts := strings.Split(module, "/")
				return parts[len(parts)-1]
			}
		}
	}

	// Fall back to directory name
	return filepath.Base(dir)
}

// buildCallGraph builds the call graph from analysis
func (a *GoAnalyzer) buildCallGraph(analysis *Analysis) {
	for _, fn := range analysis.Functions {
		analysis.CallGraph[fn.Name] = fn.Calls
	}
}

// buildTypeGraph builds the type graph from analysis
func (a *GoAnalyzer) buildTypeGraph(analysis *Analysis) {
	for _, typ := range analysis.Types {
		if len(typ.ImplementsInterfaces) > 0 {
			analysis.TypeGraph[typ.Name] = typ.ImplementsInterfaces
		}
	}
}

// extractDependencies extracts external dependencies
func (a *GoAnalyzer) extractDependencies(analysis *Analysis) {
	seen := make(map[string]bool)

	for _, file := range analysis.Files {
		for _, imp := range file.Imports {
			if seen[imp.Path] {
				continue
			}
			seen[imp.Path] = true

			dep := Dependency{
				Path:     imp.Path,
				IsStdLib: !strings.Contains(imp.Path, "."),
				IsLocal:  imp.IsLocal,
			}
			analysis.Dependencies = append(analysis.Dependencies, dep)
		}
	}
}
