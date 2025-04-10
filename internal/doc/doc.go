package doc

/*
Struct param
*/
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Return struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

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

type StructDoc struct {
	BaseDoc
	Fields []StructField `json:"fields"`
}

type FileDoc struct {
	FileName string `json:"filename"`
	Path     string `json:"path"`
	// should be BaseDoc objects
	Docs []any `json:"docs"`
}

type ProjectDoc struct {
	PackageDocs map[string][]FileDoc `json:"packageDocs"`
}
