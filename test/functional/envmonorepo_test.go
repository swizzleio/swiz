//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/swizzleio/swiz/cmd/cmds"
	"github.com/swizzleio/swiz/pkg/errtype"
	"github.com/swizzleio/swiz/test/functional/util"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEnvMonoRepoIntegrationSuite(t *testing.T) {
	suite.Run(t, new(EnvMonoRepoIntegrationSuite))
}

// EnvHappyPathIntegrationSuite test suite
type EnvMonoRepoIntegrationSuite struct {
	suite.Suite
}

func bootstrapEnvIntegration(mocks cmds.FixtureMocks) error {
	errList := errtype.ErrList{}
	homeDir, err := os.UserHomeDir()
	errList.Add(err)
	errList.Add(util.CopyDir(mocks.Fs, "./data/basefs/src", filepath.Join(homeDir, "src")))
	errList.Add(util.CopyDir(mocks.Fs, "./data/basefs/swiz", filepath.Join(homeDir, ".swiz")))

	return errList.ErrOrNil()
}

func (*EnvMonoRepoIntegrationSuite) SetupAllSuite() {
}

func (*EnvMonoRepoIntegrationSuite) TearDownSuite() {
	// TODO: Check to see if there are any tagged stacks, if so, delete them
}

// TestEnvAADeploy tests deployment. A bit of an overview here. The AA/AB/AC/AD naming convention takes advantage of how
// testify suites order tests alphabetically. While there is a philosophy that tests should be able to be run independently
// when you are talking E2E suites, this begins to become problematic, especially when code structure is considered.
// It's very easy for an E2E scenario to completely turn into a big ball of mud. The answer here may be more of a framework
// that allows for E2E tests in a single test but this may also introduce long test times.
func (s *EnvMonoRepoIntegrationSuite) TestEnvAADeploy() {
	cmd := []string{"swiz", "--appconfig", "file://~/src/monorepo/app-config.yaml", "env", "deploy", "--name", "AwesomeEnv"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	//fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "Stack: AwesomeEnv-swizboot [Creating] - Create", lines[0])
	assert.Equal(s.T(), "Stack: AwesomeEnv-swizsleep [Creating] - Create", lines[1])
}

func (s *EnvMonoRepoIntegrationSuite) TestEnvABInfo() {
	cmd := []string{"swiz", "--appconfig", "file://~/src/monorepo/app-config.yaml", "env", "info", "--name", "AwesomeEnv"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "Name: AwesomeEnv", lines[0])
	assert.Equal(s.T(), "Status: {AwesomeEnv Complete Complete }", lines[1])
	assert.Equal(s.T(), "Stacks [Status]:", lines[2])
	assert.Equal(s.T(), "  AwesomeEnv-swizsleep [Complete]", lines[3])
	assert.Equal(s.T(), "  AwesomeEnv-swizboot [Complete]", lines[4])
}

func (s *EnvMonoRepoIntegrationSuite) TestEnvACList() {
	cmd := []string{"swiz", "--appconfig", "file://~/src/monorepo/app-config.yaml", "env", "list"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "AwesomeEnv", lines[0])
}

func (s *EnvMonoRepoIntegrationSuite) TestEnvADDelete() {
	cmd := []string{"swiz", "--appconfig", "file://~/src/monorepo/app-config.yaml", "env", "delete", "--name", "AwesomeEnv"}

	_, resp := util.RunCommandWithMocks(s.T(), cmd, nil, 300, false, bootstrapEnvIntegration)

	fmt.Println(resp)

	lines := strings.Split(resp, "\n")
	assert.Equal(s.T(), "Environment AwesomeEnv deleted", lines[0])
}
