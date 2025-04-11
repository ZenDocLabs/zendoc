package system

import (
	"io/fs"
	"os"
)

type OSFileSystem struct {
	FileSystem
}

func (fs OSFileSystem) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (fs OSFileSystem) WriteFile(path string, data []byte, perm uint32) error {
	return os.WriteFile(path, data, os.FileMode(perm))
}

func (fs OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs OSFileSystem) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}
