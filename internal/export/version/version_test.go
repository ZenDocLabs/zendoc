package version

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dterbah/zendoc/internal/export/helper"
	"github.com/stretchr/testify/assert"
)

func TestUpdateVersions_FileDoesNotExist(t *testing.T) {
	versionPath := "/tmp/test_versions.json"
	currentVersion := "1.0.0"

	defer os.Remove(versionPath)

	if helper.IsFileExist(versionPath) {
		t.Fatal("Test file should not exist before running test")
	}

	err := UpdateVersions(versionPath, currentVersion)

	assert.NoError(t, err)
	assert.FileExists(t, versionPath)

	data, err := os.ReadFile(versionPath)
	assert.NoError(t, err)

	var versions Versions
	err = json.Unmarshal(data, &versions)
	assert.NoError(t, err)
	assert.Contains(t, versions.Versions, currentVersion)
}

func TestUpdateVersions_VersionAlreadyExists(t *testing.T) {
	// Arrange
	versionPath := "/tmp/test_versions.json"
	currentVersion := "1.0.0"
	versions := Versions{Versions: []string{currentVersion}}
	err := saveVersions(versionPath, versions)
	assert.NoError(t, err)
	defer os.Remove(versionPath)

	// Act
	err = UpdateVersions(versionPath, currentVersion)

	// Assert
	assert.NoError(t, err)

	data, err := os.ReadFile(versionPath)
	assert.NoError(t, err)

	var updatedVersions Versions
	err = json.Unmarshal(data, &updatedVersions)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(updatedVersions.Versions))
}

func TestUpdateVersions_AddNewVersion(t *testing.T) {
	versionPath := "/tmp/test_versions.json"
	currentVersion := "1.1.0"

	defer os.Remove(versionPath)

	err := UpdateVersions(versionPath, currentVersion)

	assert.NoError(t, err)
	assert.FileExists(t, versionPath)

	data, err := os.ReadFile(versionPath)
	assert.NoError(t, err)

	var versions Versions
	err = json.Unmarshal(data, &versions)
	assert.NoError(t, err)
	assert.Contains(t, versions.Versions, currentVersion)
}

func TestSaveVersions_ErrorOnFileWrite(t *testing.T) {
	invalidPath := "/tmp/non/existent/directory/test_invalid_write_versions.json"
	versions := Versions{Versions: []string{"1.0.0"}}

	err := saveVersions(invalidPath, versions)

	assert.Error(t, err)
}
