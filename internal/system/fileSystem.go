package system

import "io/fs"

/*
@description Interface for file system operations
@author Dorian TERBAH
*/
type FileSystem interface {
	/*
	   @description Checks if a file exists at the specified path.
	   @param path string - The file path to check.
	   @author Dorian TERBAH
	   @return bool - Returns true if the file exists, false otherwise.
	*/
	FileExists(path string) bool

	/*
	   @description Writes data to a file at the specified path with the given permissions.
	   @param path string - The path where the file should be written.
	   @param data []byte - The data to write to the file.
	   @param perm uint32 - The file permissions to set.
	   @author Dorian TERBAH
	   @return error - An error if the file could not be written.
	*/
	WriteFile(path string, data []byte, perm uint32) error

	/*
	   @description Creates a directory and all necessary parent directories with the specified permissions.
	   @param path string - The path of the directory to create.
	   @param perm fs.FileMode - The permissions to set for the directory.
	   @author Dorian TERBAH
	   @return error - An error if the directory could not be created.
	*/
	MkdirAll(path string, perm fs.FileMode) error

	/*
	   @description Renames a file or directory from oldPath to newPath.
	   @param oldPath string - The current path of the file or directory.
	   @param newPath string - The new path to rename the file or directory to.
	   @author Dorian TERBAH
	   @return error - An error if the file or directory could not be renamed.
	*/
	Rename(oldPath, newPath string) error
}
