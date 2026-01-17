package semantic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kon1790/rpg/internal/importer/treesitter"
)

func TestGoAnalyzerFile(t *testing.T) {
	code := []byte(`package main

import "fmt"

// User represents a user in the system
type User struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}

// Validator is an interface for validation
type Validator interface {
	Validate() error
}

// NewUser creates a new user
func NewUser(id int, name string) *User {
	return &User{ID: id, Name: name}
}

func main() {
	user := NewUser(1, "test")
	fmt.Println(user.Name)
}
`)

	analyzer := NewGoAnalyzer()
	if !analyzer.IsAvailable() {
		t.Skip("Go analyzer not available")
	}

	analysis, err := analyzer.AnalyzeFile("test.go", code)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Check functions
	if len(analysis.Functions) < 2 {
		t.Errorf("expected at least 2 functions, got %d", len(analysis.Functions))
	}

	foundNewUser := false
	for _, fn := range analysis.Functions {
		if fn.Name == "NewUser" {
			foundNewUser = true
			if len(fn.Parameters) != 2 {
				t.Errorf("NewUser should have 2 parameters, got %d", len(fn.Parameters))
			}
			if fn.ReturnType != "*User" {
				t.Errorf("NewUser return type should be *User, got %s", fn.ReturnType)
			}
		}
	}
	if !foundNewUser {
		t.Error("function NewUser not found")
	}

	// Check types
	if len(analysis.Types) < 2 {
		t.Errorf("expected at least 2 types, got %d", len(analysis.Types))
	}

	foundUser := false
	foundValidator := false
	for _, typ := range analysis.Types {
		if typ.Name == "User" {
			foundUser = true
			if typ.Kind != treesitter.TypeKindStruct {
				t.Errorf("User should be a struct, got %s", typ.Kind)
			}
			if len(typ.ResolvedFields) != 2 {
				t.Errorf("User should have 2 fields, got %d", len(typ.ResolvedFields))
			}
		}
		if typ.Name == "Validator" {
			foundValidator = true
			if typ.Kind != treesitter.TypeKindInterface {
				t.Errorf("Validator should be an interface, got %s", typ.Kind)
			}
		}
	}
	if !foundUser {
		t.Error("type User not found")
	}
	if !foundValidator {
		t.Error("type Validator not found")
	}

	// Check imports
	if len(analysis.Imports) != 1 {
		t.Errorf("expected 1 import, got %d", len(analysis.Imports))
	}
}

func TestDefaultRegistry(t *testing.T) {
	registry := DefaultRegistry()

	// Check that all languages are registered
	languages := []treesitter.Language{
		treesitter.LanguageGo,
		treesitter.LanguageTypeScript,
		treesitter.LanguagePython,
		treesitter.LanguageJava,
		treesitter.LanguageRust,
		treesitter.LanguageCSharp,
	}

	for _, lang := range languages {
		analyzer, ok := registry.Get(lang)
		if !ok {
			t.Errorf("analyzer for %s not registered", lang)
			continue
		}
		if analyzer.Language() != lang {
			t.Errorf("analyzer language mismatch: expected %s, got %s", lang, analyzer.Language())
		}
	}

	// Check that Go analyzer is always available
	goAnalyzer, _ := registry.Get(treesitter.LanguageGo)
	if !goAnalyzer.IsAvailable() {
		t.Error("Go analyzer should always be available")
	}
}

func TestTypeScriptAnalyzerFile(t *testing.T) {
	code := []byte(`
interface User {
  id: number;
  name: string;
}

function createUser(id: number, name: string): User {
  return { id, name };
}

export class UserService {
  private users: User[] = [];

  async getUser(id: number): Promise<User | undefined> {
    return this.users.find(u => u.id === id);
  }
}
`)

	analyzer := NewTypeScriptAnalyzer()
	analysis, err := analyzer.AnalyzeFile("test.ts", code)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Check functions
	foundCreateUser := false
	for _, fn := range analysis.Functions {
		if fn.Name == "createUser" {
			foundCreateUser = true
			if len(fn.Parameters) != 2 {
				t.Errorf("createUser should have 2 parameters, got %d", len(fn.Parameters))
			}
		}
	}
	if !foundCreateUser {
		t.Error("function createUser not found")
	}

	// Check types
	foundUser := false
	foundUserService := false
	for _, typ := range analysis.Types {
		if typ.Name == "User" {
			foundUser = true
		}
		if typ.Name == "UserService" {
			foundUserService = true
		}
	}
	if !foundUser {
		t.Error("interface User not found")
	}
	if !foundUserService {
		t.Error("class UserService not found")
	}
}

func TestPythonAnalyzerFile(t *testing.T) {
	code := []byte(`
from dataclasses import dataclass
from typing import Optional

@dataclass
class User:
    id: int
    name: str

class UserService:
    def __init__(self):
        self.users = []

    async def get_user(self, user_id: int) -> Optional[User]:
        for user in self.users:
            if user.id == user_id:
                return user
        return None

def create_user(id: int, name: str) -> User:
    return User(id=id, name=name)
`)

	analyzer := NewPythonAnalyzer()
	analysis, err := analyzer.AnalyzeFile("test.py", code)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Check functions
	foundCreateUser := false
	for _, fn := range analysis.Functions {
		if fn.Name == "create_user" {
			foundCreateUser = true
		}
	}
	if !foundCreateUser {
		t.Error("function create_user not found")
	}

	// Check types
	foundUser := false
	foundUserService := false
	for _, typ := range analysis.Types {
		if typ.Name == "User" {
			foundUser = true
		}
		if typ.Name == "UserService" {
			foundUserService = true
		}
	}
	if !foundUser {
		t.Error("class User not found")
	}
	if !foundUserService {
		t.Error("class UserService not found")
	}
}

func TestGoAnalyzerDirectory(t *testing.T) {
	// Create a temporary directory with Go files
	tempDir, err := os.MkdirTemp("", "go-analyzer-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write a go.mod file
	goMod := []byte("module testproject\n\ngo 1.21\n")
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), goMod, 0644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	// Write a main.go file
	mainGo := []byte(`package main

import "fmt"

type Config struct {
	Host string
	Port int
}

func main() {
	cfg := Config{Host: "localhost", Port: 8080}
	fmt.Printf("Starting server on %s:%d\n", cfg.Host, cfg.Port)
}
`)
	if err := os.WriteFile(filepath.Join(tempDir, "main.go"), mainGo, 0644); err != nil {
		t.Fatalf("failed to write main.go: %v", err)
	}

	// Write a service.go file
	serviceGo := []byte(`package main

type Service struct {
	name string
}

func NewService(name string) *Service {
	return &Service{name: name}
}

func (s *Service) Run() error {
	return nil
}
`)
	if err := os.WriteFile(filepath.Join(tempDir, "service.go"), serviceGo, 0644); err != nil {
		t.Fatalf("failed to write service.go: %v", err)
	}

	analyzer := NewGoAnalyzer()
	analysis, err := analyzer.Analyze(tempDir)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Check project name
	if analysis.Name != "testproject" {
		t.Errorf("expected project name 'testproject', got '%s'", analysis.Name)
	}

	// Check files
	if len(analysis.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(analysis.Files))
	}

	// Check types (Config and Service)
	if len(analysis.Types) < 2 {
		t.Errorf("expected at least 2 types, got %d", len(analysis.Types))
	}

	// Check functions (main, NewService, Run)
	if len(analysis.Functions) < 3 {
		t.Errorf("expected at least 3 functions, got %d", len(analysis.Functions))
	}

	// Check call graph was built
	if len(analysis.CallGraph) == 0 {
		t.Error("call graph should not be empty")
	}

	// Check dependencies were extracted
	if len(analysis.Dependencies) == 0 {
		t.Error("dependencies should not be empty")
	}

	// Check that fmt dependency exists
	foundFmt := false
	for _, dep := range analysis.Dependencies {
		if dep.Path == "fmt" {
			foundFmt = true
			if !dep.IsStdLib {
				t.Error("fmt should be marked as stdlib")
			}
		}
	}
	if !foundFmt {
		t.Error("fmt dependency not found")
	}
}
