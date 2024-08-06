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
	assert.Equal(t, "App config: dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vfi9zaW1wbGVjZmcvZW52LWRlZi55YW1sCmRpc2FibGVkX2NvbW1hbmRzOiBbXQo=", lines[len(lines)-4])
	assert.Equal(t, "Signature: 79c78e8152139418d5a3837241910b76466f3511cea945841531c69c5c71e9c6", lines[len(lines)-3])
	assert.Equal(t, "Word list: jawbone retraction orca inventive Dupont barbecue Pluto borderline sterling pandemic Mohawk holiness cranky miracle alone impetus cubic hemisphere chopper Babylon spyglass passenger crusade Jupiter backfield company southward October escape hideaway treadmill responsive", lines[len(lines)-2])
}
