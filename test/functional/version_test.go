//go:build functionaltest
package functional

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/swizzleio/swiz/cmd/cmds"
	"io"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	// Create pipes
	stdinR, stdinW, _ := os.Pipe()
	stdoutR, stdoutW, _ := os.Pipe()

	// Backup original stdin and stdout
	origStdin := os.Stdin
	origStdout := os.Stdout

	// Redirect stdin and stdout
	os.Stdin = stdinR
	os.Stdout = stdoutW

	// Prepare to capture output
	outputC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, stdoutR)
		assert.NoError(t, err)
		outputC <- buf.String()
	}()

	// Provide input
	input := ""
	go func() {
		_, err := stdinW.Write([]byte(input))
		assert.NoError(t, err)
		assert.NoError(t, stdinW.Close())
	}()

	os.Args = []string{"swiz", "version"}

	// Run the version command
	ret := cmds.Execute()
	assert.Equal(t, 0, ret)

	// Clean up
	assert.NoError(t, stdoutW.Close())
	os.Stdin = origStdin
	os.Stdout = origStdout

	// Get the captured output. Hangs here...
	capturedOutput := <-outputC

	fmt.Println("Captured Output:", capturedOutput)
}
