//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/swizzleio/swiz/test/functional/util"
	"strings"
	"testing"
)

func TestEnvHappyPathIntegrationSuite(t *testing.T) {
	suite.Run(t, new(EnvHappyPathIntegrationSuite))
}

// EnvHappyPathIntegrationSuite test suite
type EnvHappyPathIntegrationSuite struct {
	suite.Suite
}

func (*EnvHappyPathIntegrationSuite) SetupAllSuite() {
}

func (*EnvHappyPathIntegrationSuite) TearDownSuite() {
	// TODO: Check to see if there are any tagged stacks, if so, delete them
}

// TestEnvAADeploy tests deployment. A bit of an overview here. The AA/AB/AC/AD naming convention takes advantage of how
// testify suites order tests alphabetically. While there is a philosophy that tests should be able to be run independently
// when you are talking E2E suites, this begins to become problematic, especially when code structure is considered.
// It's very easy for an E2E scenario to completely turn into a big ball of mud. The answer here may be more of a framework
// that allows for E2E tests in a single test but this may also introduce long test times.
func (s *EnvHappyPathIntegrationSuite) TestEnvAADeploy() {
	cmd := []string{"swiz", "env", "deploy", "--name", "AwesomeEnv"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	stackStatus := []string{
		"Stack: AwesomeEnv-swizboot [Creating] - Create",
		"Stack: AwesomeEnv-swizstack1 [Creating] - Create",
		"Stack: AwesomeEnv-swizstack2 [Creating] - Create",
		"Stack: AwesomeEnv-swizstack3 [Creating] - Create",
		"",
	}

	assert.ElementsMatch(s.T(), stackStatus, lines)
}

func (s *EnvHappyPathIntegrationSuite) TestEnvABInfo() {
	cmd := []string{"swiz", "env", "info", "--name", "AwesomeEnv"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "Name: AwesomeEnv", lines[0])
	assert.Equal(s.T(), "Status: {AwesomeEnv Complete Complete }", lines[1])
	assert.Equal(s.T(), "Stacks [Status]:", lines[2])
	stackStatus := []string{
		"  AwesomeEnv-swizstack3 [Complete]",
		"  AwesomeEnv-swizstack2 [Complete]",
		"  AwesomeEnv-swizstack1 [Complete]",
		"  AwesomeEnv-swizboot [Complete]",
		"",
	}
	assert.ElementsMatch(s.T(), stackStatus, lines[len(lines)-5:])
}

func (s *EnvHappyPathIntegrationSuite) TestEnvACList() {
	cmd := []string{"swiz", "env", "list"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "AwesomeEnv", lines[0])
}

func (s *EnvHappyPathIntegrationSuite) TestEnvADDelete() {
	cmd := []string{"swiz", "env", "delete", "--name", "AwesomeEnv"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "Environment AwesomeEnv deleted", lines[0])
}
