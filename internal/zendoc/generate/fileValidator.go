package generate

import (
	"strings"
)

/*
@description Check if a file is a Go test file by verifying if it ends with "_test.go"
@param filePath string - The path of the file to check
@return bool - true if the file is a test file, false otherwise
@example IsTestFile("example_test.go") => true
@author Dorian TERBAH
*/

func IsTestFile(filePath string) bool {
	return strings.HasSuffix(filePath, "_test.go")
}

/*
@description Check if a file is a Go main file by verifying if it ends with "main.go"
@param filepath string - The path of the file to check
@return bool - true if the file is a main file, false otherwise
@example IsMainFile("main.go") => true
@author Dorian TERBAH
*/

func IsMainFile(filepath string) bool {
	return strings.HasSuffix(filepath, "main.go")
}
