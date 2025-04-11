package system

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSFileSystem_FileExists(t *testing.T) {
	fs := OSFileSystem{}

	tempFile := "/tmp/testfile_exists"
	err := os.WriteFile(tempFile, []byte("test data"), 0644)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	assert.True(t, fs.FileExists(tempFile))

	assert.False(t, fs.FileExists("/tmp/non_existent_file"))
}

func TestOSFileSystem_WriteFile(t *testing.T) {
	fs := OSFileSystem{}

	// Define test file path and data
	tempFile := "/tmp/testfile_write"
	data := []byte("This is a test")

	err := fs.WriteFile(tempFile, data, 0644)
	assert.NoError(t, err)
	defer os.Remove(tempFile)

	content, err := os.ReadFile(tempFile)
	assert.NoError(t, err)
	assert.Equal(t, data, content)
}

func TestOSFileSystem_MkdirAll(t *testing.T) {
	fs := OSFileSystem{}

	dirPath := "/tmp/testdir/subdir"

	err := fs.MkdirAll(dirPath, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll("/tmp/testdir")

	_, err = os.Stat(dirPath)
	assert.NoError(t, err)
}

func TestOSFileSystem_Rename(t *testing.T) {
	fs := OSFileSystem{}

	oldFilePath := "/tmp/testfile_old"
	newFilePath := "/tmp/testfile_new"

	err := os.WriteFile(oldFilePath, []byte("test data"), 0644)
	assert.NoError(t, err)
	defer os.Remove(oldFilePath)

	err = fs.Rename(oldFilePath, newFilePath)
	assert.NoError(t, err)
	defer os.Remove(newFilePath)

	_, err = os.Stat(newFilePath)
	assert.NoError(t, err)
	_, err = os.Stat(oldFilePath)
	assert.True(t, os.IsNotExist(err))
}
