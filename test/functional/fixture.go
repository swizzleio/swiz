//go:build functionaltest

package functional

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/swizzleio/swiz/cmd/cmds"
)

type RunFunctionalTest func(t *testing.T, appFs afero.Fs, resp string)

func RunTestWithMocks(t *testing.T, command []string, expect map[string]string, handler RunFunctionalTest) {
	appFs := cmds.SetupFixtures()

	stdinR, stdinW, err := os.Pipe()
	assert.NoError(t, err)
	stdoutR, stdoutW, err := os.Pipe()
	assert.NoError(t, err)

	origStdin := os.Stdin
	origStdout := os.Stdout

	os.Stdin = stdinR
	os.Stdout = stdoutW

	cmdDone := make(chan bool)
	capturedOutputChan := make(chan string)

	go handleStdout(t, stdoutR, stdinW, expect, cmdDone, capturedOutputChan)

	os.Args = command

	cmds.Execute()
	assert.NoError(t, stdoutW.Close()) // Close writer to signal EOF to stdoutR
	cmdDone <- true                    // Signal command completion

	capturedOutput := <-capturedOutputChan

	os.Stdin = origStdin
	os.Stdout = origStdout
	assert.NoError(t, stdinR.Close())
	assert.NoError(t, stdinW.Close())
	assert.NoError(t, stdoutR.Close())

	ret := 0 // Mock return from cmds.Execute()

	handler(t, appFs, capturedOutput)

	assert.Equal(t, 0, ret)
	close(capturedOutputChan)
	close(cmdDone)
}

func handleStdout(t *testing.T, stdout io.ReadCloser, stdin io.WriteCloser, expect map[string]string, cmdDone chan bool, capturedOutputChan chan string) {
	scanner := bufio.NewScanner(stdout)
	output := ""
	cmdRunning := true
	for scanner.Scan() && cmdRunning {
		line := scanner.Text()
		output += line

		select {
		case <-cmdDone:
			cmdRunning = false
		default:
			if response, ok := expect[line]; ok {
				_, err := stdin.Write([]byte(response + "\n"))
				assert.NoError(t, err)
			}
		}
	}

	capturedOutputChan <- output
}
