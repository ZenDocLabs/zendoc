package system

type CommandRunner interface {
	Execute(dir string, name string, arg ...string) ([]byte, error)
}
