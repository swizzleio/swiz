//go:build functionaltest

package functional

import (
	"github.com/swizzleio/swiz/cmd/cmds"
	"github.com/swizzleio/swiz/test/functional/util"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDisableCommand(t *testing.T) {
	cmd := []string{"swiz", "config"}

	_, resp := util.RunCommandWithMocks(t, cmd, nil, util.DefaultCmdTimeoutSec, false, func(mocks cmds.FixtureMocks) error {
		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)
		return util.CopyFile(mocks.Fs, "./data/basefs/src/monorepo/app-config.yaml", filepath.Join(homeDir, ".swiz", "app-config.yaml"))
	})

	assert.NotContains(t, resp, "generate")
}
