//go:build functionaltest

package functional

import (
	"bufio"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/swizzleio/swiz/cmd/cmds"
)

type testResult struct {
	response   string
	failureMsg string
}

type ExpectResponseMatching int

const (
	Exact ExpectResponseMatching = iota
	Fuzzy
)

type ExpectResponse struct {
	Match    ExpectResponseMatching
	Output   string
	Response string
}

type RunFunctionalTest func(t *testing.T, appFs afero.Fs, resp string)

func RunTestWithMocks(t *testing.T, command []string, expect []ExpectResponse, timeoutSec time.Duration, handler RunFunctionalTest) {
	// Set up fixture
	appFs := cmds.SetupFixtures()

	// Set up pipes
	stdinR, stdinW, err := os.Pipe()
	assert.NoError(t, err)
	stdoutR, stdoutW, err := os.Pipe()
	assert.NoError(t, err)

	origStdin := os.Stdin
	origStdout := os.Stdout

	os.Stdin = stdinR
	os.Stdout = stdoutW

	// Start up output monitoring
	cmdDone := make(chan bool, 1)
	capturedOutputChan := make(chan testResult, 1)
	go handleStdout(stdoutR, stdinW, timeoutSec, expect, cmdDone, capturedOutputChan)

	// Run command
	os.Args = command

	go func() {
		ret := cmds.Execute()
		assert.Equal(t, 0, ret)
		assert.NoError(t, stdoutW.Close())
		cmdDone <- true
	}()

	// Capture output
	capturedOutput := <-capturedOutputChan

	// Clean up
	os.Stdin = origStdin
	os.Stdout = origStdout
	assert.NoError(t, stdinR.Close())
	assert.NoError(t, stdinW.Close())
	assert.NoError(t, stdoutR.Close())

	// Assert
	if capturedOutput.failureMsg != "" {
		assert.Failf(t, capturedOutput.failureMsg, capturedOutput.response)
	}
	handler(t, appFs, capturedOutput.response)
}

func handleStdout(stdout io.ReadCloser, stdin io.WriteCloser, timeoutSec time.Duration,
	expect []ExpectResponse, cmdDone chan bool, capturedOutputChan chan testResult) {
	if expect == nil {
		expect = []ExpectResponse{}
	}

	reader := bufio.NewReader(stdout)
	output := ""
	cmdRunning := true
	timeout := time.After(timeoutSec * time.Second)
	result := testResult{}
	for cmdRunning {
		line, err := reader.ReadString('\n') // TODO: Dumbsurvey may need to append a \n to every question
		if err != nil {
			if err != io.EOF {
				result.failureMsg = fmt.Sprintf("error reading stdout: %s", err)
			}
			break
		}

		output += line

		select {
		case <-cmdDone:
			cmdRunning = false
		case <-timeout:
			result.failureMsg = "timed out waiting for output"
			cmdRunning = false
		default:
			if len(expect) != 0 {
				resp, rErr := getExpectResponse(line, expect)
				if rErr != nil {
					result.failureMsg = rErr.Error()
					cmdRunning = false
				} else {
					_, wErr := stdin.Write([]byte(resp + "\n"))
					if wErr != nil {
						result.failureMsg = "failed to write to stdin"
						cmdRunning = false
					}
				}
			}
		}
	}

	result.response = output
	capturedOutputChan <- result
}

func getExpectResponse(output string, expect []ExpectResponse) (string, error) {
	for _, response := range expect {
		if response.Match == Exact {
			if response.Output == output {
				return response.Response, nil
			}
		} else if response.Match == Fuzzy {
			if fuzzy.Match(response.Output, output) {
				return response.Response, nil
			}
		}
	}

	return "", fmt.Errorf("response not found from output: %s", output)
}
