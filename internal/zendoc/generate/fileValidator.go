package generate

import (
	"strings"
)

func IsTestFile(filePath string) bool {
	return strings.HasSuffix(filePath, "_test.go")
}

func IsMainFile(filepath string) bool {
	return strings.HasSuffix(filepath, "main.go")
}
