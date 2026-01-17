package treesitter

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// LanguageParser defines the interface for language-specific parsing
type LanguageParser interface {
	// Parse parses source code and returns the parse result
	Parse(code []byte, filename string) (*ParseResult, error)
	// Language returns the language this parser handles
	Language() Language
	// FileExtensions returns file extensions handled by this parser
	FileExtensions() []string
}

// Parser is the main tree-sitter parser that delegates to language-specific parsers
type Parser struct {
	parsers map[Language]LanguageParser
}

// NewParser creates a new multi-language parser
func NewParser() *Parser {
	p := &Parser{
		parsers: make(map[Language]LanguageParser),
	}

	// Register all language parsers
	p.RegisterParser(NewGoParser())
	p.RegisterParser(NewTypeScriptParser())
	p.RegisterParser(NewPythonParser())
	p.RegisterParser(NewJavaParser())
	p.RegisterParser(NewRustParser())
	p.RegisterParser(NewCSharpParser())

	return p
}

// RegisterParser registers a language-specific parser
func (p *Parser) RegisterParser(lp LanguageParser) {
	p.parsers[lp.Language()] = lp
}

// Parse parses source code in the given language
func (p *Parser) Parse(code []byte, filename string, lang Language) (*ParseResult, error) {
	parser, ok := p.parsers[lang]
	if !ok {
		return nil, fmt.Errorf("no parser registered for language: %s", lang)
	}
	return parser.Parse(code, filename)
}

// ParseAuto attempts to auto-detect language from filename and parse
func (p *Parser) ParseAuto(code []byte, filename string) (*ParseResult, error) {
	lang := DetectLanguage(filename)
	if lang == "" {
		return nil, fmt.Errorf("could not detect language for file: %s", filename)
	}
	return p.Parse(code, filename, lang)
}

// GetParser returns the parser for a specific language
func (p *Parser) GetParser(lang Language) (LanguageParser, bool) {
	parser, ok := p.parsers[lang]
	return parser, ok
}

// DetectLanguage detects the language from a filename
func DetectLanguage(filename string) Language {
	lower := strings.ToLower(filename)

	switch {
	case strings.HasSuffix(lower, ".go"):
		return LanguageGo
	case strings.HasSuffix(lower, ".ts"), strings.HasSuffix(lower, ".tsx"):
		return LanguageTypeScript
	case strings.HasSuffix(lower, ".py"):
		return LanguagePython
	case strings.HasSuffix(lower, ".java"):
		return LanguageJava
	case strings.HasSuffix(lower, ".rs"):
		return LanguageRust
	case strings.HasSuffix(lower, ".cs"):
		return LanguageCSharp
	default:
		return ""
	}
}

// getTreeSitterLanguage returns the tree-sitter language for a Language
func getTreeSitterLanguage(lang Language) *sitter.Language {
	switch lang {
	case LanguageGo:
		return golang.GetLanguage()
	case LanguageTypeScript:
		return typescript.GetLanguage()
	case LanguagePython:
		return python.GetLanguage()
	case LanguageJava:
		return java.GetLanguage()
	case LanguageRust:
		return rust.GetLanguage()
	case LanguageCSharp:
		return csharp.GetLanguage()
	default:
		return nil
	}
}

// baseParser provides common parsing functionality
type baseParser struct {
	lang       Language
	extensions []string
	tsLang     *sitter.Language
}

func (b *baseParser) Language() Language {
	return b.lang
}

func (b *baseParser) FileExtensions() []string {
	return b.extensions
}

// parseTree creates a tree-sitter parse tree from source code
func (b *baseParser) parseTree(code []byte) (*sitter.Tree, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(b.tsLang)

	tree := parser.Parse(nil, code)
	if tree == nil {
		return nil, fmt.Errorf("tree-sitter parse failed: returned nil tree")
	}

	return tree, nil
}

// nodeText extracts the text content of a node
func nodeText(code []byte, node *sitter.Node) string {
	return string(code[node.StartByte():node.EndByte()])
}

// nodeLocation extracts the source location of a node
func nodeLocation(filename string, node *sitter.Node) SourceLocation {
	start := node.StartPoint()
	end := node.EndPoint()
	return SourceLocation{
		File:      filename,
		StartLine: int(start.Row) + 1, // tree-sitter uses 0-based line numbers
		EndLine:   int(end.Row) + 1,
		StartCol:  int(start.Column),
		EndCol:    int(end.Column),
	}
}

// hashNode creates a hash of a node's content for comparison
func hashNode(code []byte, node *sitter.Node) string {
	content := code[node.StartByte():node.EndByte()]
	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:8]) // Use first 8 bytes for shorter hash
}

// findChildByType finds the first child of a node with the given type
func findChildByType(node *sitter.Node, typeName string) *sitter.Node {
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == typeName {
			return child
		}
	}
	return nil
}

// findChildrenByType finds all children of a node with the given type
func findChildrenByType(node *sitter.Node, typeName string) []*sitter.Node {
	var children []*sitter.Node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == typeName {
			children = append(children, child)
		}
	}
	return children
}

// findChildByFieldName finds a child node by its field name
func findChildByFieldName(node *sitter.Node, fieldName string) *sitter.Node {
	return node.ChildByFieldName(fieldName)
}

// walkTree walks the tree depth-first and calls fn for each node
func walkTree(node *sitter.Node, fn func(*sitter.Node) bool) {
	if !fn(node) {
		return
	}
	for i := 0; i < int(node.ChildCount()); i++ {
		walkTree(node.Child(i), fn)
	}
}

// collectNodes collects all nodes matching a predicate
func collectNodes(root *sitter.Node, predicate func(*sitter.Node) bool) []*sitter.Node {
	var nodes []*sitter.Node
	walkTree(root, func(node *sitter.Node) bool {
		if predicate(node) {
			nodes = append(nodes, node)
		}
		return true
	})
	return nodes
}

// getCommentAbove extracts documentation comments above a node
func getCommentAbove(code []byte, node *sitter.Node, parent *sitter.Node) string {
	if parent == nil {
		return ""
	}

	// Find this node's index in parent
	nodeIndex := -1
	for i := 0; i < int(parent.ChildCount()); i++ {
		if parent.Child(i) == node {
			nodeIndex = i
			break
		}
	}

	if nodeIndex <= 0 {
		return ""
	}

	// Look for comment nodes before this one
	var comments []string
	for i := nodeIndex - 1; i >= 0; i-- {
		sibling := parent.Child(i)
		siblingType := sibling.Type()

		if siblingType == "comment" || siblingType == "line_comment" ||
			siblingType == "block_comment" || siblingType == "doc_comment" {
			comment := nodeText(code, sibling)
			comments = append([]string{cleanComment(comment)}, comments...)
		} else if siblingType != "newline" && siblingType != "\n" {
			break
		}
	}

	return strings.TrimSpace(strings.Join(comments, "\n"))
}

// cleanComment removes comment delimiters
func cleanComment(comment string) string {
	// Remove common comment prefixes
	comment = strings.TrimPrefix(comment, "//")
	comment = strings.TrimPrefix(comment, "#")
	comment = strings.TrimPrefix(comment, "///")
	comment = strings.TrimPrefix(comment, "/**")
	comment = strings.TrimSuffix(comment, "*/")
	comment = strings.TrimPrefix(comment, "/*")

	// Clean up each line for block comments
	lines := strings.Split(comment, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, " ")
		lines[i] = line
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// isExported checks if a name is exported/public based on language conventions
func isExported(name string, lang Language) bool {
	if name == "" {
		return false
	}

	switch lang {
	case LanguageGo:
		// Go: exported if starts with uppercase
		return name[0] >= 'A' && name[0] <= 'Z'
	case LanguagePython:
		// Python: private if starts with underscore
		return !strings.HasPrefix(name, "_")
	default:
		// Default: assume public
		return true
	}
}
