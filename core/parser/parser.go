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
)

const GO_EXTENSION = ".go"

type DocParserFileValidator = func(string) bool

type DocParser struct {
	FileValidators []DocParserFileValidator
}

func (docParser DocParser) ParseDocForDir(dirPath string) (*doc.ProjectDoc, error) {
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
			if filepath.Ext(fullPath) == GO_EXTENSION {
				pckName, fileDoc := docParser.ParseDocForFile(fullPath)
				_, ok := projectDoc.PackageDocs[pckName]
				if !ok {
					projectDoc.PackageDocs[pckName] = []doc.FileDoc{}
				}

				projectDoc.PackageDocs[pckName] = append(projectDoc.PackageDocs[pckName], *fileDoc)
			}
		} else if entry.Type().IsDir() {
			dirDoc, err := docParser.ParseDocForDir(fullPath)
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

	return projectDoc, nil
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

	docs := []doc.FuncDoc{}

	if err != nil {
		panic(err)
	}

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fd := docParser.ParseDocForFunction(funcDecl)

		if fd != nil {
			docs = append(docs, *fd)
		}
	}

	return packageName, &doc.FileDoc{
		Docs:     docs,
		FileName: path.Base(filePath),
	}
}

/*
@description Parse documentation for a single function
@param function *ast.FuncDecl - The function to parse
@author Dorian TERBAH
@return *doc.FuncDoc - Associated function documentation object
*/
func (docParser DocParser) ParseDocForFunction(function *ast.FuncDecl) *doc.FuncDoc {
	// Regex for tags in documentation
	paramRegex := regexp.MustCompile(`@param\s+(\w+)\s+(.+?)\s*-\s*(.*)`)
	returnRegex := regexp.MustCompile(`@return\s+(.+?)\s*-\s*(.*)`)
	exampleRegex := regexp.MustCompile(`@example\s*(.*)`)
	descriptionRegex := regexp.MustCompile(`@description\s*(.*)`)
	authorRegex := regexp.MustCompile(`@author\s*(.*)`)
	deprecatedRegex := regexp.MustCompile(`@deprecated\s*(.*)`)

	fd := &doc.FuncDoc{
		Name:   function.Name.Name,
		Params: []doc.Param{},
	}

	lines := sanitizeLines(function.Doc)
	if len(lines) == 0 {
		// no documentation available for this function, skip it
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

		default:
		}
	}

	return fd
}

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
