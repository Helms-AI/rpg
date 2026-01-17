package treesitter

import (
	"testing"
)

func TestGoParser(t *testing.T) {
	code := []byte(`
package main

import (
	"fmt"
	"context"
)

// User represents a system user
type User struct {
	ID   int    // user identifier
	Name string // user name
}

// GetUser retrieves a user by ID
func GetUser(ctx context.Context, id int) (*User, error) {
	fmt.Println("Getting user", id)
	return &User{ID: id, Name: "Test"}, nil
}
`)

	parser := NewGoParser()
	result, err := parser.Parse(code, "main.go")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check functions
	if len(result.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(result.Functions))
	} else {
		fn := result.Functions[0]
		if fn.Name != "GetUser" {
			t.Errorf("Expected function name 'GetUser', got '%s'", fn.Name)
		}
		if len(fn.Parameters) != 2 {
			t.Errorf("Expected 2 parameters, got %d", len(fn.Parameters))
		}
		if fn.ReturnType == "" {
			t.Error("Expected return type to be set")
		}
	}

	// Check types
	if len(result.Types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(result.Types))
	} else {
		typ := result.Types[0]
		if typ.Name != "User" {
			t.Errorf("Expected type name 'User', got '%s'", typ.Name)
		}
		if typ.Kind != TypeKindStruct {
			t.Errorf("Expected struct kind, got '%s'", typ.Kind)
		}
		if len(typ.Fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(typ.Fields))
		}
	}

	// Check imports
	if len(result.Imports) != 2 {
		t.Errorf("Expected 2 imports, got %d", len(result.Imports))
	}
}

func TestTypeScriptParser(t *testing.T) {
	code := []byte(`
import { Context } from 'context';
import axios from 'axios';

interface User {
	id: number;
	name: string;
	email?: string;
}

async function getUser(ctx: Context, id: number): Promise<User> {
	const response = await axios.get('/users/' + id);
	return response.data;
}

class UserService {
	private cache: Map<number, User>;

	constructor() {
		this.cache = new Map();
	}

	async fetchUser(id: number): Promise<User | null> {
		return this.cache.get(id) ?? null;
	}
}
`)

	parser := NewTypeScriptParser()
	result, err := parser.Parse(code, "user.ts")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check functions
	if len(result.Functions) < 1 {
		t.Errorf("Expected at least 1 function, got %d", len(result.Functions))
	}

	// Check types (interface + class)
	if len(result.Types) < 2 {
		t.Errorf("Expected at least 2 types, got %d", len(result.Types))
	}

	// Check imports
	if len(result.Imports) < 2 {
		t.Errorf("Expected at least 2 imports, got %d", len(result.Imports))
	}
}

func TestPythonParser(t *testing.T) {
	code := []byte(`
from typing import Optional
import asyncio

class User:
    """Represents a system user"""

    def __init__(self, id: int, name: str):
        self.id = id
        self.name = name

    def get_display_name(self) -> str:
        return f"User: {self.name}"

async def get_user(user_id: int) -> Optional[User]:
    """Retrieves a user by ID"""
    await asyncio.sleep(0.1)
    return User(user_id, "Test")
`)

	parser := NewPythonParser()
	result, err := parser.Parse(code, "user.py")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check functions
	if len(result.Functions) < 1 {
		t.Errorf("Expected at least 1 function, got %d", len(result.Functions))
	}

	// Find async function
	foundAsync := false
	for _, fn := range result.Functions {
		if fn.Name == "get_user" && fn.IsAsync {
			foundAsync = true
			break
		}
	}
	if !foundAsync {
		t.Error("Expected to find async function 'get_user'")
	}

	// Check types (class)
	if len(result.Types) < 1 {
		t.Errorf("Expected at least 1 type, got %d", len(result.Types))
	}
}

func TestJavaParser(t *testing.T) {
	code := []byte(`
package com.example;

import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class User {
    private final int id;
    private String name;

    public User(int id, String name) {
        this.id = id;
        this.name = name;
    }

    public int getId() {
        return id;
    }

    public String getName() {
        return name;
    }
}

public interface UserService {
    Optional<User> getUser(int id);
    CompletableFuture<User> getUserAsync(int id);
}
`)

	parser := NewJavaParser()
	result, err := parser.Parse(code, "User.java")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check types (class + interface)
	if len(result.Types) < 2 {
		t.Errorf("Expected at least 2 types, got %d", len(result.Types))
	}

	// Check imports
	if len(result.Imports) < 2 {
		t.Errorf("Expected at least 2 imports, got %d", len(result.Imports))
	}
}

func TestRustParser(t *testing.T) {
	code := []byte(`
use std::collections::HashMap;
use tokio::sync::Mutex;

/// Represents a system user
pub struct User {
    pub id: i32,
    pub name: String,
}

/// User status enum
pub enum UserStatus {
    Active,
    Inactive,
    Pending,
}

impl User {
    pub fn new(id: i32, name: String) -> Self {
        User { id, name }
    }
}

/// Retrieves a user by ID
pub async fn get_user(id: i32) -> Option<User> {
    Some(User::new(id, "Test".to_string()))
}
`)

	parser := NewRustParser()
	result, err := parser.Parse(code, "user.rs")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check functions
	if len(result.Functions) < 1 {
		t.Errorf("Expected at least 1 function, got %d", len(result.Functions))
	}

	// Check types (struct + enum)
	if len(result.Types) < 2 {
		t.Errorf("Expected at least 2 types, got %d", len(result.Types))
	}

	// Verify enum variants
	for _, typ := range result.Types {
		if typ.Kind == TypeKindEnum {
			if len(typ.Variants) != 3 {
				t.Errorf("Expected 3 enum variants, got %d", len(typ.Variants))
			}
			break
		}
	}
}

func TestCSharpParser(t *testing.T) {
	code := []byte(`
using System;
using System.Threading.Tasks;

namespace MyApp
{
    public class User
    {
        public int Id { get; init; }
        public string Name { get; set; }

        public User(int id, string name)
        {
            Id = id;
            Name = name;
        }
    }

    public interface IUserService
    {
        Task<User?> GetUserAsync(int id);
        User? GetUser(int id);
    }

    public enum UserStatus
    {
        Active,
        Inactive,
        Pending
    }
}
`)

	parser := NewCSharpParser()
	result, err := parser.Parse(code, "User.cs")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check types (class + interface + enum)
	if len(result.Types) < 3 {
		t.Errorf("Expected at least 3 types, got %d", len(result.Types))
	}

	// Check imports
	if len(result.Imports) < 2 {
		t.Errorf("Expected at least 2 imports, got %d", len(result.Imports))
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		filename string
		expected Language
	}{
		{"main.go", LanguageGo},
		{"app.ts", LanguageTypeScript},
		{"app.tsx", LanguageTypeScript},
		{"script.py", LanguagePython},
		{"Main.java", LanguageJava},
		{"lib.rs", LanguageRust},
		{"Program.cs", LanguageCSharp},
		{"readme.md", ""},
		{"unknown.xyz", ""},
	}

	for _, tt := range tests {
		result := DetectLanguage(tt.filename)
		if result != tt.expected {
			t.Errorf("DetectLanguage(%s) = %s, expected %s", tt.filename, result, tt.expected)
		}
	}
}

func TestMultiLanguageParser(t *testing.T) {
	parser := NewParser()

	// Test that all parsers are registered
	languages := []Language{
		LanguageGo,
		LanguageTypeScript,
		LanguagePython,
		LanguageJava,
		LanguageRust,
		LanguageCSharp,
	}

	for _, lang := range languages {
		_, ok := parser.GetParser(lang)
		if !ok {
			t.Errorf("Parser for %s not registered", lang)
		}
	}
}
