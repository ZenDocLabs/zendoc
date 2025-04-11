package app_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dterbah/zendoc/internal/export/app"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAppConfig_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "VERSION.json")

	err := app.UpdateAppConfig(filePath, "v1.0.0", "Initial release")
	assert.NoError(t, err)

	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	var config app.AppConfig
	err = json.Unmarshal(data, &config)
	assert.NoError(t, err)

	assert.Equal(t, []string{"v1.0.0"}, config.Versions)
	assert.Equal(t, "Initial release", config.Description)
}

func TestUpdateAppConfig_AppendNewVersion(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "VERSION.json")

	initialConfig := app.AppConfig{
		Versions:    []string{"v1.0.0"},
		Description: "Initial release",
	}
	data, _ := json.MarshalIndent(initialConfig, "", "  ")
	_ = os.WriteFile(filePath, data, 0644)

	err := app.UpdateAppConfig(filePath, "v1.1.0", "Initial release")
	assert.NoError(t, err)

	newData, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	var config app.AppConfig
	_ = json.Unmarshal(newData, &config)

	assert.ElementsMatch(t, []string{"v1.0.0", "v1.1.0"}, config.Versions)
	assert.Equal(t, "Initial release", config.Description)
}

func TestUpdateAppConfig_DuplicateVersion(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "VERSION.json")

	initialConfig := app.AppConfig{
		Versions:    []string{"v1.0.0"},
		Description: "Initial release",
	}
	data, _ := json.MarshalIndent(initialConfig, "", "  ")
	_ = os.WriteFile(filePath, data, 0644)

	err := app.UpdateAppConfig(filePath, "v1.0.0", "Initial release")
	assert.NoError(t, err)

	newData, _ := os.ReadFile(filePath)
	var config app.AppConfig
	_ = json.Unmarshal(newData, &config)

	assert.Equal(t, 1, len(config.Versions))
	assert.Contains(t, config.Versions, "v1.0.0")
	assert.Equal(t, "Initial release", config.Description)
}
