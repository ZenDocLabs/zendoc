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

func TestParseDocForFile(t *testing.T) {
	docParser := DocParser{}
	content := `package testpkg

	// @description Does something
	func TestFunc() {}

	// Not a function
	var something = 42
	`

	tmpFile := writeTempFile(t, "file.go", content)
	defer os.Remove(tmpFile)

	pkg, fileDoc := docParser.ParseDocForFile(tmpFile)

	assert.Equal(t, "testpkg", pkg)
	assert.NotNil(t, fileDoc)
	assert.Equal(t, 1, len(fileDoc.Docs))
	assert.Equal(t, "TestFunc", fileDoc.Docs[0].Name)
	assert.Equal(t, "Does something", fileDoc.Docs[0].Description)
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
