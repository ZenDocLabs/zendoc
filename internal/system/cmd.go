package system

/*
@description Interface for running system commands
@author Dorian TERBAH
*/
type CommandRunner interface {
	/*
	   @description Executes a system command in the specified directory.
	   @param dir string - The directory where the command will be executed.
	   @param name string - The name of the command to execute.
	   @param arg ...string - Additional arguments to pass to the command.
	   @author Dorian TERBAH
	   @return ([]byte, error) - The output from the executed command, and an error if the command fails.
	*/
	Execute(dir string, name string, arg ...string) ([]byte, error)
}
