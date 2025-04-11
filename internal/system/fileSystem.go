package system

import "io/fs"

type FileSystem interface {
	FileExists(path string) bool
	WriteFile(path string, data []byte, perm uint32) error
	MkdirAll(path string, perm fs.FileMode) error
	Rename(oldPath, newPath string) error
}
