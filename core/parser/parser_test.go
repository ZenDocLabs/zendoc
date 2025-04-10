package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocForFunction_WithStructReceiver(t *testing.T) {
	docParser := DocParser{}
	src := `
	// @description Adds two numbers
	// @param a int - first number
	// @param b int - second number
	// @return int - sum of two numbers
	// @author Alice
	// @deprecated Use AddV2 instead
	// @example Add(1, 10) -> 11
	func (c *Calculator) Add(a int, b int) int {
		return a + b
	}`

	node := parseFunc(t, src)
	fd := docParser.ParseDocForFunction(node)

	assert.NotNil(t, fd)
	assert.Equal(t, "Add", fd.Name)
	assert.Equal(t, "Calculator", fd.Struct)
	assert.Equal(t, "Adds two numbers", fd.Description)
	assert.Equal(t, 2, len(fd.Params))
	assert.Equal(t, "a", fd.Params[0].Name)
	assert.Equal(t, "int", fd.Params[0].Type)
	assert.Equal(t, "first number", fd.Params[0].Description)
	assert.Equal(t, "b", fd.Params[1].Name)
	assert.Equal(t, "int", fd.Params[1].Type)
	assert.Equal(t, "second number", fd.Params[1].Description)
	assert.NotNil(t, fd.Return)
	assert.Equal(t, "int", fd.Return.Type)
	assert.Equal(t, "sum of two numbers", fd.Return.Description)
	assert.Equal(t, "Alice", fd.Author)
	assert.Equal(t, "Add(1, 10) -> 11", fd.Example)
	assert.Equal(t, "Use AddV2 instead", fd.Deprecated)
}

func TestParseDocForFunction_WithNoReceiver(t *testing.T) {
	docParser := DocParser{}
	src := `
	// @description Multiplies two numbers
	// @param x int - first number
	// @param y int - second number
	// @return int - product of two numbers
	// @example Multiply(1, 2) -> 2
	func Multiply(x int, y int) int {
		return x * y
	}`

	node := parseFunc(t, src)
	fd := docParser.ParseDocForFunction(node)

	assert.NotNil(t, fd)
	assert.Equal(t, "Multiply", fd.Name)
	assert.Empty(t, fd.Struct)
	assert.Equal(t, "Multiplies two numbers", fd.Description)
	assert.Equal(t, 2, len(fd.Params))
	assert.Equal(t, "x", fd.Params[0].Name)
	assert.Equal(t, "int", fd.Params[0].Type)
	assert.Equal(t, "first number", fd.Params[0].Description)
	assert.Equal(t, "y", fd.Params[1].Name)
	assert.Equal(t, "int", fd.Params[1].Type)
	assert.Equal(t, "second number", fd.Params[1].Description)
	assert.NotNil(t, fd.Return)
	assert.Equal(t, "int", fd.Return.Type)
	assert.Equal(t, "product of two numbers", fd.Return.Description)
	assert.Equal(t, "Multiply(1, 2) -> 2", fd.Example)
	assert.Empty(t, fd.Author)
	assert.Empty(t, fd.Deprecated)
}

func TestParseDocForFunction_WithAllTags(t *testing.T) {
	docParser := DocParser{}
	src := `
	// @description Does something cool
	// @param input string - the input value
	// @return string - the output value
	// @example DoSomething("input")
	// @author Dorian
	// @deprecated Use DoSomethingElse instead
	func DoSomething(input string) string {
		return input
	}`

	node := parseFunc(t, src)
	fd := docParser.ParseDocForFunction(node)

	assert.NotNil(t, fd)
	assert.Equal(t, "DoSomething", fd.Name)
	assert.Equal(t, "Does something cool", fd.Description)
	assert.Equal(t, 1, len(fd.Params))
	assert.Equal(t, "input", fd.Params[0].Name)
	assert.Equal(t, "string", fd.Params[0].Type)
	assert.Equal(t, "the input value", fd.Params[0].Description)
	assert.NotNil(t, fd.Return)
	assert.Equal(t, "string", fd.Return.Type)
	assert.Equal(t, "the output value", fd.Return.Description)
	assert.Equal(t, "DoSomething(\"input\")", fd.Example)
	assert.Equal(t, "Dorian", fd.Author)
	assert.Equal(t, "Use DoSomethingElse instead", fd.Deprecated)
}

func TestSanitizeLines(t *testing.T) {
	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "// @description Hello"},
			{Text: "// @param name string - user name"},
			{Text: "/*\n* @return string - result\n*/"},
		},
	}

	lines := sanitizeLines(commentGroup)

	assert.Equal(t, []string{
		"@description Hello",
		"@param name string - user name",
		"@return string - result",
	}, lines)
}

func TestGetPackageName(t *testing.T) {
	content := `package dummy`
	tmpFile := writeTempFile(t, "dummy.go", content)
	defer os.Remove(tmpFile)

	name, err := getPackageName(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, "dummy", name)
}

// ---- Helpers functions ---- //
func parseFunc(t *testing.T, src string) *ast.FuncDecl {
	t.Helper()
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", "package main\n"+src, parser.ParseComments)
	assert.NoError(t, err)

	for _, decl := range node.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			return fn
		}
	}

	t.Fatal("no function found")
	return nil
}

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(t, err)
	return path
}

// Struct tests

func TestParseDocForStruct_WithAllTags(t *testing.T) {
	docParser := DocParser{}

	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "// @description This struct represents a user"},
			{Text: "// @author John"},
			{Text: "// @deprecated Use NewUser instead"},
			{Text: "// @field Name string - the user's name"},
			{Text: "// @field Age int - the user's age"},
		},
	}

	structDoc := docParser.ParseDocForStruct(commentGroup, "User")

	assert.NotNil(t, structDoc)
	assert.Equal(t, "User", structDoc.Name)
	assert.Equal(t, "This struct represents a user", structDoc.Description)
	assert.Equal(t, "John", structDoc.Author)
	assert.Equal(t, "Use NewUser instead", structDoc.Deprecated)
	assert.Equal(t, "struct", structDoc.Type)

	assert.Len(t, structDoc.Fields, 2)

	assert.Equal(t, "Name", structDoc.Fields[0].Name)
	assert.Equal(t, "string", structDoc.Fields[0].Type)
	assert.Equal(t, "the user's name", structDoc.Fields[0].Description)

	assert.Equal(t, "Age", structDoc.Fields[1].Name)
	assert.Equal(t, "int", structDoc.Fields[1].Type)
	assert.Equal(t, "the user's age", structDoc.Fields[1].Description)
}

func TestParseDocForStruct_NoTags(t *testing.T) {
	docParser := DocParser{}

	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "// no tags here"},
		},
	}

	structDoc := docParser.ParseDocForStruct(commentGroup, "Empty")
	assert.Nil(t, structDoc)
}

func TestParseDocForStruct_PartialTags(t *testing.T) {
	docParser := DocParser{}

	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "// @description Partial struct"},
			{Text: "// @field ID string - identifier"},
		},
	}

	structDoc := docParser.ParseDocForStruct(commentGroup, "Partial")

	assert.NotNil(t, structDoc)
	assert.Equal(t, "Partial", structDoc.Name)
	assert.Equal(t, "Partial struct", structDoc.Description)
	assert.Empty(t, structDoc.Author)
	assert.Empty(t, structDoc.Deprecated)
	assert.Len(t, structDoc.Fields, 1)
	assert.Equal(t, "ID", structDoc.Fields[0].Name)
	assert.Equal(t, "string", structDoc.Fields[0].Type)
	assert.Equal(t, "identifier", structDoc.Fields[0].Description)
}
