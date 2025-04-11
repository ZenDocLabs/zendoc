package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOSCommandRunner_Execute_Success(t *testing.T) {
	runner := OSCommandRunner{}

	output, err := runner.Execute("", "echo", "hello world")
	assert.NoError(t, err)
	assert.Contains(t, string(output), "hello world")
}

func TestOSCommandRunner_Execute_Failure(t *testing.T) {
	runner := OSCommandRunner{}

	output, err := runner.Execute("", "fakecommand123")
	assert.Error(t, err)
	assert.Empty(t, output)
}

func TestOSCommandRunner_Execute_WithWorkingDir(t *testing.T) {
	runner := OSCommandRunner{}

	tempDir := os.TempDir()

	var output []byte
	var err error
	if runtime.GOOS == "windows" {
		output, err = runner.Execute(tempDir, "cmd", "/C", "cd")
	} else {
		output, err = runner.Execute(tempDir, "pwd")
	}

	assert.NoError(t, err)

	outputStr := strings.TrimSpace(string(output))

	normalizedTempDir := filepath.Clean(tempDir)
	normalizedOutput := filepath.Clean(outputStr)

	assert.Contains(t, normalizedOutput, normalizedTempDir)
}
