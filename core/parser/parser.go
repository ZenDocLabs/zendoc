package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dterbah/zendoc/internal/doc"
	"github.com/fatih/color"
)

const GO_EXTENSION = ".go"

type DocParserFileValidator = func(string) bool
type DocParserFunctionValidator = func(string) bool

/*
@description Struct responsible for orchestrating validation logic when parsing documentation from Go source files. It holds a list of validators for files and functions to modularize and organize parsing rules and behaviors.
@field FileValidators []DocParserFileValidator - A list of validators applied at the file level (e.g. checking file-level tags, imports, etc.)
@field FunctionValidators []DocParserFunctionValidator - A list of validators specifically designed to validate function-level documentation (e.g. param/return tag parsing, required fields, etc.)
*/
type DocParser struct {
	FileValidators     []DocParserFileValidator
	FunctionValidators []DocParserFunctionValidator
}

/*
@description Check if a file is valid for documentation generation using the provided file validators
@param filepath string - The file path to validate
@return bool - true if the file is valid, false otherwise
*/
func (docParser DocParser) isValidateFileForDoc(filepath string) bool {
	for _, validator := range docParser.FileValidators {
		if !validator(filepath) {
			return false
		}
	}

	return true
}

/*
@description Check if a function name is valid for documentation generation using the provided function validators
@param name string - The function name to validate
@return bool - true if the function is valid, false otherwise
*/
func (docParser DocParser) isValidateFunction(name string) bool {
	for _, validator := range docParser.FunctionValidators {
		if !validator(name) {
			return false
		}
	}

	return true
}

/*
@description Recursively parse documentation in a directory and its subdirectories
@param dirPath string - The root path to scan
@param currentPath string - The relative path used for output (maintains relative structure)
@return *doc.ProjectDoc, error - The parsed project documentation and an error if something went wrong
@example ParseDocForDir("./myproject", "")
*/
func (docParser DocParser) ParseDocForDir(dirPath string, currentPath string) (*doc.ProjectDoc, error) {
	entries, err := os.ReadDir(dirPath)
	projectDoc := &doc.ProjectDoc{
		PackageDocs: make(map[string][]doc.FileDoc),
	}
	if err != nil {
		return nil, fmt.Errorf("error when listing the files of the dir %s", dirPath)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.Type().IsRegular() {
			fileName := entry.Name()
			if filepath.Ext(fullPath) == GO_EXTENSION {
				if !docParser.isValidateFileForDoc(fileName) {
					color.HiYellow("File \"%s\" skipped", fileName)
					continue
				}

				color.Green("File \"%s\" being processed...", path.Base(fileName))

				pckName, fileDoc := docParser.ParseDocForFile(fullPath)
				_, ok := projectDoc.PackageDocs[pckName]
				if !ok {
					projectDoc.PackageDocs[pckName] = []doc.FileDoc{}
				}
				fileDoc.Path = filepath.Join(currentPath, entry.Name())
				if len(fileDoc.Docs) > 0 {
					projectDoc.PackageDocs[pckName] = append(projectDoc.PackageDocs[pckName], *fileDoc)
				}
			}
		} else if entry.Type().IsDir() {
			dirDoc, err := docParser.ParseDocForDir(fullPath, filepath.Join(currentPath, entry.Name()))
			if err != nil {
				return nil, fmt.Errorf("error when retrieving doc of the directory %s", fullPath)
			}
			// iterate through the doc
			for pckName, docs := range dirDoc.PackageDocs {
				_, ok := projectDoc.PackageDocs[pckName]
				if !ok {
					projectDoc.PackageDocs[pckName] = []doc.FileDoc{}
				}

				projectDoc.PackageDocs[pckName] = append(projectDoc.PackageDocs[pckName], docs...)
			}
		}
	}

	removeEmptyPackages(projectDoc)

	return projectDoc, nil
}

func removeEmptyPackages(projectDoc *doc.ProjectDoc) {
	for pckName, doc := range projectDoc.PackageDocs {
		if len(doc) == 0 {
			delete(projectDoc.PackageDocs, pckName)
		}
	}
}

// @description Parse the documentation for a single file
//
// @param filePath string - The file path
// @deprecated Just a small test
// @author Dorian TERBAH
// @return (string, []doc.FuncDoc) - The associated doc for the file. If no package is mentioned, it return an empty string and nil
// @example ParseDocForFile("myfile.go")
func (docParser DocParser) ParseDocForFile(filePath string) (string, *doc.FileDoc) {
	// retrieve package name
	packageName, err := getPackageName(filePath)
	if err != nil {
		return "", nil
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)

	docs := []any{}

	if err != nil {
		panic(err)
	}

	for _, decl := range node.Decls {
		funcDecl, isFunction := decl.(*ast.FuncDecl)
		if isFunction {
			if !docParser.isValidateFunction(funcDecl.Name.Name) {
				color.HiYellow("	Skip function \"%s\"", funcDecl.Name.Name)
				continue
			}

			fd := docParser.ParseDocForFunction(funcDecl)
			if fd != nil {
				docs = append(docs, *fd)
			}

			continue
		}

		genDecl, ok := decl.(*ast.GenDecl)
		if ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				_, ok = typeSpec.Type.(*ast.StructType)
				if ok {
					if genDecl.Doc != nil {
						sd := docParser.ParseDocForStruct(genDecl.Doc, typeSpec.Name.Name)
						if sd != nil {
							docs = append(docs, *sd)
						}
					}
				}
			}
		}

	}

	return packageName, &doc.FileDoc{
		Docs:     docs,
		FileName: filepath.Base(filePath),
	}
}

/*
@description Parse documentation for a struct
@param function *ast.CommentGroup - The comments line associated to the struct
@param name string - Name of the struct
@author Dorian TERBAH
@return *doc.StructDoc - Associated function documentation object, or nil if there is no comments with tags
*/
func (docParser DocParser) ParseDocForStruct(structComments *ast.CommentGroup, name string) *doc.StructDoc {
	descriptionRegex := regexp.MustCompile(`@description\s*(.*)`)
	authorRegex := regexp.MustCompile(`@author\s*(.*)`)
	deprecatedRegex := regexp.MustCompile(`@deprecated\s*(.*)`)
	fieldRegex := regexp.MustCompile(`@field\s+(\w+)\s+(.+?)\s*-\s*(.*)`)

	lines := sanitizeLines(structComments)
	if len(lines) == 0 {
		return nil
	}

	sd := &doc.StructDoc{
		Fields: []doc.StructField{},
	}

	sd.Name = name
	sd.Type = "struct"

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@description"):
			if matches := descriptionRegex.FindStringSubmatch(line); len(matches) == 2 {
				sd.Description = matches[1]
			}

		case strings.HasPrefix(line, "@author"):
			if matches := authorRegex.FindStringSubmatch(line); len(matches) == 2 {
				sd.Author = matches[1]
			}

		case strings.HasPrefix(line, "@deprecated"):
			if matches := deprecatedRegex.FindStringSubmatch(line); len(matches) == 2 {
				sd.Deprecated = matches[1]
			}

		case strings.HasPrefix(line, "@field"):
			if matches := fieldRegex.FindStringSubmatch(line); len(matches) == 4 {
				sd.Fields = append(sd.Fields, doc.StructField{
					Name:        matches[1],
					Type:        matches[2],
					Description: matches[3],
				})
			}
		}
	}

	// if no documentation is available
	if sd.Description == "" && sd.Author == "" && sd.Deprecated == "" && len(sd.Fields) == 0 {
		return nil
	}

	return sd
}

/*
@description Parse documentation for a single function
@param function *ast.FuncDecl - The function to parse
@author Dorian TERBAH
@return *doc.FuncDoc - Associated function documentation object, or nil if there is not tagged comments
*/
func (docParser DocParser) ParseDocForFunction(function *ast.FuncDecl) *doc.FuncDoc {
	paramRegex := regexp.MustCompile(`@param\s+(\w+)\s+(.+?)\s*-\s*(.*)`)
	returnRegex := regexp.MustCompile(`@return\s+(.+?)\s*-\s*(.*)`)
	exampleRegex := regexp.MustCompile(`@example\s*(.*)`)
	descriptionRegex := regexp.MustCompile(`@description\s*(.*)`)
	authorRegex := regexp.MustCompile(`@author\s*(.*)`)
	deprecatedRegex := regexp.MustCompile(`@deprecated\s*(.*)`)

	fd := &doc.FuncDoc{
		Params: []doc.Param{},
	}

	fd.Name = function.Name.Name
	fd.Type = "function"

	// Check if it's a method associated with a struct
	if function.Recv != nil && len(function.Recv.List) > 0 {
		// The type can be *ast.Ident (e.g., T) or *ast.StarExpr (e.g., *T)
		switch expr := function.Recv.List[0].Type.(type) {
		case *ast.Ident:
			fd.Struct = expr.Name
		case *ast.StarExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				fd.Struct = ident.Name
			}
		}
	}

	lines := sanitizeLines(function.Doc)
	if len(lines) == 0 {
		return nil
	}

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@param"):
			if matches := paramRegex.FindStringSubmatch(line); len(matches) == 4 {
				fd.Params = append(fd.Params, doc.Param{
					Name:        matches[1],
					Type:        matches[2],
					Description: matches[3],
				})
			}
		case strings.HasPrefix(line, "@return"):
			if matches := returnRegex.FindStringSubmatch(line); len(matches) == 3 {
				fd.Return = &doc.Return{
					Type:        matches[1],
					Description: matches[2],
				}
			}
		case strings.HasPrefix(line, "@example"):
			if matches := exampleRegex.FindStringSubmatch(line); len(matches) == 2 {
				fd.Example = matches[1]
			}
		case strings.HasPrefix(line, "@description"):
			if matches := descriptionRegex.FindStringSubmatch(line); len(matches) == 2 {
				fd.Description = matches[1]
			}
		case strings.HasPrefix(line, "@author"):
			if matches := authorRegex.FindStringSubmatch(line); len(matches) == 2 {
				fd.Author = matches[1]
			}
		case strings.HasPrefix(line, "@deprecated"):
			if matches := deprecatedRegex.FindStringSubmatch(line); len(matches) == 2 {
				fd.Deprecated = matches[1]
			}
		}
	}

	return fd
}

/*
@description Sanitize and flatten comment lines (block or single-line) to a slice of clean strings
@param doc *ast.CommentGroup - The group of AST comments to sanitize
@return []string - The cleaned lines
*/
func sanitizeLines(doc *ast.CommentGroup) []string {
	lines := []string{}

	if doc != nil {
		for _, comment := range doc.List {
			text := comment.Text

			if strings.HasPrefix(text, "//") {
				// line comment
				line := strings.TrimPrefix(text, "//")
				line = strings.TrimSpace(line)
				lines = append(lines, line)
			} else if strings.HasPrefix(text, "/*") {
				// block comment
				block := strings.TrimPrefix(text, "/*")
				block = strings.TrimSuffix(block, "*/")
				blockLines := strings.Split(block, "\n")

				for _, line := range blockLines {
					line = strings.TrimSpace(strings.TrimPrefix(line, "*"))
					if line != "" {
						lines = append(lines, line)
					}
				}
			}
		}
	}

	return lines
}

/*
@description Retrieve the package name of a file
@param filePath string - The path of the file
@author Dorian TERBAH
@return (string, error) - The associated package name and an error if the parsing failed
@example getPackageName("./parser.go") => parser
*/
func getPackageName(filePath string) (string, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, filePath, nil, parser.PackageClauseOnly)

	if err != nil {
		return "", nil
	}

	return node.Name.Name, nil
}
