//go:build functionaltest

package functional

import (
	"github.com/swizzleio/swiz/cmd/cmds"
	"github.com/swizzleio/swiz/test/functional/util"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigExport(t *testing.T) {
	cmd := []string{"swiz", "config", "export"}

	_, resp := util.RunCommandWithMocks(t, cmd, nil, util.DefaultCmdTimeoutSec, false, func(mocks cmds.FixtureMocks) error {
		homeDir, err := os.UserHomeDir()
		assert.NoError(t, err)
		return util.CopyDir(mocks.Fs, "./data/basefs/swiz", filepath.Join(homeDir, ".swiz"))
	})

	lines := strings.Split(resp, "\n")
	assert.Equal(t, "App config: dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vLi9jZnNpbXBsZWNmZy9lbnYtZGVmLnlhbWwKZGlzYWJsZWRfY29tbWFuZHM6IFtdCmJhc2VfZGlyOiB+L3NyYwo=", lines[len(lines)-4])
	assert.Equal(t, "Signature: 9fed093a15de3c00c7fc41991d35afe9415f950b20166fd7636f67b9eb7fdaea", lines[len(lines)-3])
	assert.Equal(t, "Word list: quota unify Algol corrosion backfield telephone cobra adroitness soybean Wilmington cranky nebula Belfast conformist rocker ultimate cranky forever preclude armistice bison bodyguard gremlin stethoscope flatfoot hemisphere freedom proximate trouble integrate surmount undaunted", lines[len(lines)-2])
}
