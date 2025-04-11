package config

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) FileExists(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

func (m *MockFileSystem) WriteFile(filename string, content []byte, perm uint32) error {
	args := m.Called(filename, content, perm)
	return args.Error(0)
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockFileSystem) Rename(oldPath, newPath string) error {
	args := m.Called(oldPath, newPath)
	return args.Error(0)
}

func TestGetConfiguration_Success(t *testing.T) {
	expectedConfig := &Config{
		ProjectConfig: ProjectConfig{
			Name:        "Test Project",
			Description: "A test project",
			Version:     "1.0.0",
			GitLink:     "https://github.com/test/test.git",
			MainBranch:  "main",
			DocPath:     "/docs",
		},
		DocConfig: DocConfig{
			IncludePrivate: true,
			IncludeTests:   true,
			IncludeMain:    true,
			ExcludeFiles:   []string{"file1.go", "file2.go"},
		},
	}

	// Simule la lecture du fichier
	fileContent, _ := json.MarshalIndent(expectedConfig, "", "  ")
	err := os.WriteFile(ZENDOC_CONFIG_FILE, fileContent, 0644)
	if err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}
	defer os.Remove(ZENDOC_CONFIG_FILE)

	config, err := GetConfiguration()
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}

func TestGetConfiguration_FileNotFound(t *testing.T) {
	os.Remove(ZENDOC_CONFIG_FILE)
	config, err := GetConfiguration()
	assert.Error(t, err)
	assert.Nil(t, config)
}
