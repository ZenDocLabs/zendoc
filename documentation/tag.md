# Zendoc Documentation Tags

Zendoc uses a special documentation format with tags that look similar to Javadoc. These tags help structure your documentation and provide specific information about your Go code. Below is a comprehensive guide to all available documentation tags with examples.

## Tag Format

Documentation blocks start with `/*` and end with `*/`. Each tag begins with `@` and is followed by its content.

## Available Tags

### @description

Used to provide a general description of a function, method, or struct.

**Example (Function):**

```go
/*
@description Retrieve the package name of a file
*/
func getPackageName(filePath string) (string, error) {
    // Function implementation
}
```

**Example (Struct):**

```go
/*
@description Struct responsible for orchestrating validation logic when parsing documentation from Go source files
*/
type DocParser struct {
    // Fields
}
```

### @param

Documents a function or method parameter. Format: `@param paramName type - Description`

**Example:**

```go
/*
@param filePath string - The path of the file
*/
func getPackageName(filePath string) (string, error) {
    // Function implementation
}
```

### @return

Documents the return value(s) of a function or method. Format: `@return (type) - Description`

**Example:**

```go
/*
@return (string, error) - The associated package name and an error if the parsing failed
*/
func getPackageName(filePath string) (string, error) {
    // Function implementation
}
```

### @author

Specifies the author of the code.

**Example:**

```go
/*
@author Dorian TERBAH
*/
func getPackageName(filePath string) (string, error) {
    // Function implementation
}
```

### @example

Provides an example of how to use the function or method.

**Example:**

```go
/*
@example getPackageName("./parser.go") => parser
*/
func getPackageName(filePath string) (string, error) {
    // Function implementation
}
```

### @field

Documents a field in a struct. Format: `@field fieldName type - Description`

**Example:**

```go
/*
@field FileValidators []DocParserFileValidator - A list of validators applied at the file level
*/
type DocParser struct {
    FileValidators []DocParserFileValidator
}
```

### @deprecated

Indicates that a function, method, or struct is deprecated and should not be used.

**Example:**

```go
/*
@deprecated Use newGetPackageName() instead
*/
func getPackageName(filePath string) (string, error) {
    // Function implementation
}
```

## Complete Documentation Examples

### Function Documentation

```go
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
```

### Struct Documentation

```go
/*
@description Struct responsible for orchestrating validation logic when parsing documentation from Go source files. It holds a list of validators for files and functions to modularize and organize parsing rules and behaviors.
@field FileValidators []DocParserFileValidator - A list of validators applied at the file level (e.g. checking file-level tags, imports, etc.)
@field FunctionValidators []DocParserFunctionValidator - A list of validators specifically designed to validate function-level documentation (e.g. param/return tag parsing, required fields, etc.)
@author Dorian TERBAH
*/
type DocParser struct {
    FileValidators     []DocParserFileValidator
    FunctionValidators []DocParserFunctionValidator
}
```

### Interface Documentation

```go
/*
@description Interface for running system commands
@author Dorian TERBAH
*/
type CommandRunner interface {
	/*
	   @description Executes a system command in the specified directory.
	   @param dir string - The directory where the command will be executed.
	   @param name string - The name of the command to execute.
	   @param arg ...string - Additional arguments to pass to the command.
	   @author Dorian TERBAH
	   @return ([]byte, error) - The output from the executed command, and an error if the command fails.
	*/
	Execute(dir string, name string, arg ...string) ([]byte, error)
}
```

## Best Practices

1. Always include a `@description` tag to provide context
2. Document all parameters with `@param` tags
3. Document return values with `@return` tags
4. Use `@example` to show usage patterns
5. Include `@author` information for maintainability
6. For structs, document each field with `@field` tags

Following these guidelines will ensure that Zendoc generates comprehensive and useful documentation for your Go project.
