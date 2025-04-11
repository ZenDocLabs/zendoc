package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTestFile(t *testing.T) {
	assert.True(t, IsTestFile("example_test.go"), "Expected 'example_test.go' to be recognized as a test file")
	assert.False(t, IsTestFile("example.go"), "Expected 'example.go' to NOT be recognized as a test file")
	assert.False(t, IsTestFile("example_test.ts"), "Expected 'example_test.ts' to NOT be recognized as a test file")
}

func TestIsMainFile(t *testing.T) {
	assert.True(t, IsMainFile("main.go"), "Expected 'main.go' to be recognized as a main file")
	assert.False(t, IsMainFile("example.go"), "Expected 'example.go' to NOT be recognized as a main file")
	assert.False(t, IsMainFile("main_test.go"), "Expected 'main_test.go' to NOT be recognized as a main file")
}
