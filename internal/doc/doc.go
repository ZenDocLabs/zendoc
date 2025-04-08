package doc

import (
	"fmt"
	"strings"
)

type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Return struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type FuncDoc struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Params      []Param `json:"params"`
	Return      *Return `json:"return"`
	Example     string  `json:"example"`
	Author      string  `json:"author"`
	Deprecated  string  `json:"deprecated"`
}

type FileDoc struct {
	FileName string    `json:"file_name"`
	Docs     []FuncDoc `json:"docs"`
}

type ProjectDoc struct {
	PackageDocs map[string][]FileDoc `json:"package_docs"`
}

func (p Param) String() string {
	return fmt.Sprintf("Param{Name: %s, Type: %s, Description: %s}", p.Name, p.Type, p.Description)
}

func (r Return) String() string {
	return fmt.Sprintf("Return{Type: %s, Description: %s}", r.Type, r.Description)
}

func (fd FuncDoc) String() string {
	var params []string
	for _, p := range fd.Params {
		params = append(params, p.String())
	}
	return fmt.Sprintf("FuncDoc{Name: %s, Description: %s, Params: [%s], Return: %s, Example: %s, Author: %s, Deprecated: %s}",
		fd.Name, fd.Description, strings.Join(params, ", "), fd.Return.String(), fd.Example, fd.Author, fd.Deprecated)
}

func (fd FileDoc) String() string {
	var docs []string
	for _, doc := range fd.Docs {
		docs = append(docs, doc.String())
	}
	return fmt.Sprintf("FileDoc{FileName: %s, Docs: [%s]}", fd.FileName, strings.Join(docs, ", "))
}

func (pd ProjectDoc) String() string {
	var packageDocs []string
	for pkg, files := range pd.PackageDocs {
		var filesStr []string
		for _, file := range files {
			filesStr = append(filesStr, file.String())
		}
		packageDocs = append(packageDocs, fmt.Sprintf("%s: [%s]", pkg, strings.Join(filesStr, ", ")))
	}
	return fmt.Sprintf("ProjectDoc{PackageDocs: {%s}}", strings.Join(packageDocs, ", "))
}
