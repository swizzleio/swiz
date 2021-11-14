package exechelper

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ExecHelper interface {
	RunShellCmd(directory string, command string, args ...string) (err error)
	CancelShellCmd() error
}

type ExecHelp struct {
	cmd *exec.Cmd
}

// NewExecHelper creates a new execution helper. Currently only a single command is supported
func NewExecHelper() ExecHelper {
	return &ExecHelp{}
}

// RunShellCmd runs a shell command.
func (e *ExecHelp) RunShellCmd(directory string, command string, args ...string) (err error) {

	if e.cmd != nil {
		return fmt.Errorf("concurrent execution of multiple commands is not currently supported")
	}
	defer func() { e.cmd = nil }()

	//cmd := exec.Command("command", "open", "-n", "-F", "-W", "-a", appParams, param)
	e.cmd = exec.Command(command, args...)
	e.cmd.Dir = directory
	stderr, err := e.cmd.StderrPipe()
	log.SetOutput(os.Stderr)
	if err != nil {
		return err
	}

	// Launch command
	if err = e.cmd.Start(); err != nil {
		return err
	}

	// Read output and wait
	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("[DEBUG] %s\n", slurp)

	if err = e.cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// CancelShellCmd cancels a running shell command
func (e *ExecHelp) CancelShellCmd() error {
	defer func() { e.cmd = nil }()
	return e.cmd.Process.Kill()
}
