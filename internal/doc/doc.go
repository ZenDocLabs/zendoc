package doc

/*
@description Struct to represent a function or struct field parameter in the documentation system.
@author Dorian TERBAH
@field Name string - The name of the parameter
@field Type string - The Go type of the parameter
@field Description string - A description of what the parameter represents
*/
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

/*
@description Struct to represent the return value of a documented function in the documentation system. Supports a single return type only for simplicity.@author
@author Dorian TERBAH
@field Type string - The Go type of the return value
@field Description string - A description of the return value
*/
type Return struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

/*
@description Base struct shared by all documentation types, providing common metadata fields such as name, author, and description.@author
@author Dorian TERBAH
@field Name string - The name of the documented item (function, struct, etc.)
@field Description string - A description of what this item does
@field Author string - The author of the item or its documentation
@field Deprecated string - A deprecation message, if the item is deprecated
@field Type string - The type of the documented item (e.g. 'function', 'struct')
*/
type BaseDoc struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Deprecated  string `json:"deprecated"`
	Type        string `json:"type"`
}

/*
@description Struct to represent the documentation associated to a function
@author Dorian TERBAH
@field Params []Param - The params of the function
@field Return *Return - The return type of the function, if it exists
@field Example string - An example of the usage of this function
*/
type FuncDoc struct {
	BaseDoc
	Params  []Param `json:"params"`
	Return  *Return `json:"return"`
	Example string  `json:"example"`
	Struct  string  `json:"struct,omitempty"`
}

type StructField = Param

/*
@description Struct to represent the documentation associated to a Go struct
@author Dorian TERBAH
@field Fields []StructField - The fields that belong to the struct, with their type and description
*/
type StructDoc struct {
	BaseDoc
	Fields []StructField `json:"fields"`
}

type InterfaceDoc struct {
	BaseDoc
	Methods []FuncDoc `json:"methods,omitempty"`
}

/*
@description Struct to represent the documentation of a Go file in a package
@author Dorian TERBAH
@field FileName string - The name of the file
@field Path string - The full path to the file
@field Docs []any - The documentation items contained in this file (functions, structs, etc.)
*/
type FileDoc struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
	// should be BaseDoc objects
	Docs []any `json:"docs"`
}

/*
@description Struct to represent the entire documentation of a project, organized by package name and files within each package
@author Dorian TERBAH
@field PackageDocs map[string][]FileDoc - A mapping from package names to their documented files
*/
type ProjectDoc struct {
	PackageDocs map[string][]FileDoc `json:"packageDocs"`
}
