package generate

import (
	"strings"

	"github.com/dterbah/zendoc/core/parser"
)

func IsTestFile(filePath string) bool {
	return strings.HasSuffix(filePath, "_test.go")
}

func IsTestFileWrapper(includeTest bool) parser.DocParserFileValidator {
	if includeTest {
		return IsTestFile
	}

	return func(filepath string) bool {
		return true
	}
}
