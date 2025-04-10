package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking dependencies
type MockDocExporter struct {
	mock.Mock
}

func (m *MockDocExporter) Export(doc interface{}) error {
	args := m.Called(doc)
	return args.Error(0)
}

// Test wrapFileValidator
func TestWrapFileValidator(t *testing.T) {
	mockValidator := func(filePath string) bool {
		return filePath == "valid_file.go"
	}

	validator := wrapFileValidator(false, mockValidator)

	// Test condition met
	assert.False(t, validator("valid_file.go"))

	// Test condition not met
	assert.True(t, validator("invalid_file.go"))
}

// Test wrapRegexFileValidator
func TestWrapRegexFileValidator(t *testing.T) {
	regexValidator := wrapRegexFileValidator([]string{`_test\.go$`, `main\.go$`})

	// Test file that should match regex
	assert.False(t, regexValidator("file_test.go"))

	// Test file that shouldn't match regex
	assert.True(t, regexValidator("file.go"))
}
