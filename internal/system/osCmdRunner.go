package system

import (
	"os/exec"
)

type OSCommandRunner struct {
	CommandRunner
}

func (r OSCommandRunner) Execute(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}
