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
	assert.Equal(t, "App config: dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vfi9zbGVlcGNmL2Vudi1kZWYueWFtbApkaXNhYmxlZF9jb21tYW5kczogW10K", lines[len(lines)-4])
	assert.Equal(t, "Signature: 8cd115d4637660341fa95066ea93945043221710f4ccc2ed1e75e88a80a4fe17", lines[len(lines)-3])
	assert.Equal(t, "Word list: offload scavenger backfield souvenir flatfoot impetus facial confidence billiard passenger drumbeat gossamer Trojan molasses Pluto embezzle crucial candidate banjo autopsy upshot revolver snapshot unify berserk impartial trauma maverick merit Pandora woodlark bookseller", lines[len(lines)-2])
}
