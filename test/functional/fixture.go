//go:build functionaltest

package functional

import (
	"bufio"
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"io"
	"os"
	"strings"
	"testing"
	"time"

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
	ExactOnce
	Fuzzy
)

type ExpectResponseAction int

const (
	Response ExpectResponseAction = iota
	NoAction
)

const DefaultCmdTimeoutSec = 2

type ExpectResponse struct {
	Match    ExpectResponseMatching
	Output   string
	Response string
	Action   ExpectResponseAction
	used     bool
}

func RunCommandWithMocks(t *testing.T, command []string, expect []*ExpectResponse, timeoutSec time.Duration) (mocks cmds.FixtureMocks, resp string) {
	// Set up pipes
	stdinR, stdinW, err := os.Pipe()
	assert.NoError(t, err)
	stdoutR, stdoutW, err := os.Pipe()
	assert.NoError(t, err)

	// Set up fixture
	mocks = cmds.SetupFixtures(stdinR, stdoutW)

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
	assert.NoError(t, stdinR.Close())
	assert.NoError(t, stdinW.Close())
	assert.NoError(t, stdoutR.Close())

	// Assert
	if capturedOutput.failureMsg != "" {
		assert.Failf(t, capturedOutput.failureMsg, capturedOutput.response)
	}

	resp = capturedOutput.response
	return
}

func handleStdout(stdout io.ReadCloser, stdin io.WriteCloser, timeoutSec time.Duration,
	expect []*ExpectResponse, cmdDone chan bool, capturedOutputChan chan testResult) {
	if expect == nil {
		expect = []*ExpectResponse{}
	}

	reader := bufio.NewReader(stdout)
	output := ""
	cmdRunning := true
	timeout := time.After(timeoutSec * time.Second)
	result := testResult{}
	for cmdRunning {
		line, err := reader.ReadString('\n')
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
					if resp.Action == Response {
						outResp := strings.TrimSpace(resp.Response) + "\n"
						_, wErr := stdin.Write([]byte(outResp))
						output += outResp
						if wErr != nil {
							result.failureMsg = "failed to write to stdin"
							cmdRunning = false
						}
					}
				}
			}
		}
	}

	result.response = output
	capturedOutputChan <- result
}

func getExpectResponse(output string, expect []*ExpectResponse) (ExpectResponse, error) {
	var resp *ExpectResponse
	for _, response := range expect {
		if response.used {
			continue
		}

		if response.Match == Exact || response.Match == ExactOnce {
			if response.Output == output {
				resp = response
				break
			}
		} else if response.Match == Fuzzy {
			if fuzzy.Match(response.Output, output) {
				resp = response
				break
			}
		}
	}

	if resp == nil {
		return ExpectResponse{}, fmt.Errorf("response not found from output: %s", output)
	}

	if resp.Match == ExactOnce {
		resp.used = true
	}

	return *resp, nil
}
