//go:build functionaltest

package functional

import (
	"fmt"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"github.com/swizzleio/swiz/test/functional/util"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigImport(t *testing.T) {
	cmd := []string{"swiz", "config", "import", ""}
	expect := []*util.ExpectResponse{
		{
			Output:   "Paste the app config here \n",
			Response: "dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vfi9zbGVlcGNmL2Vudi1kZWYueWFtbApkaXNhYmxlZF9jb21tYW5kczogW10K",
		},
		{
			Output:   "Paste the signature here (or leave blank if you don't know) \n",
			Response: "8cd115d4637660341fa95066ea93945043221710f4ccc2ed1e75e88a80a4fe17",
		},
	}
	mocks, resp := util.RunCommandWithMocks(t, cmd, expect, util.DefaultCmdTimeoutSec, false, nil)
	fmt.Println(resp)

	// Get the home directory
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	fh := fileutil.NewFileUrlHelper(mocks.Fs)
	acBuf, err := fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
	assert.NoError(t, err)

	yAc, err := util.NewYamlTestUtil(string(acBuf))
	yAc.AssertEqual(t, 1, "version")
	yAc.AssertEqual(t, "SleepySleep", "default_env")
	yAc.AssertEqual(t, "SleepySleep", "env_def[0].name")
	yAc.AssertEqual(t, "file://~/sleepcf/env-def.yaml", "env_def[0].env_def_file")
	yAc.AssertEqual(t, []interface{}{}, "disabled_commands")
}

func TestConfigImportEmptySig(t *testing.T) {
	cmd := []string{"swiz", "config", "import", ""}
	expect := []*util.ExpectResponse{
		{
			Output:   "Paste the app config here \n",
			Response: "dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vfi9zbGVlcGNmL2Vudi1kZWYueWFtbApkaXNhYmxlZF9jb21tYW5kczogW10K",
		},
		{
			Output:   "Paste the signature here (or leave blank if you don't know) \n",
			Response: "",
		},
		{
			Output:   "Does one or the other match? \n",
			Response: "y",
		},
		{
			Output: "Verify one or the other matches the signature provided with the app config\n",
			Action: util.NoAction,
		},
		{
			Output: "Signature: 8cd115d4637660341fa95066ea93945043221710f4ccc2ed1e75e88a80a4fe17\n",
			Action: util.NoAction,
		},
		{
			Output: "Word list: offload scavenger backfield souvenir flatfoot impetus facial confidence billiard passenger drumbeat gossamer Trojan molasses Pluto embezzle crucial candidate banjo autopsy upshot revolver snapshot unify berserk impartial trauma maverick merit Pandora woodlark bookseller\n",
			Action: util.NoAction,
		},
		{
			Output: "",
			Action: util.NoAction,
		},
	}
	mocks, resp := util.RunCommandWithMocks(t, cmd, expect, util.DefaultCmdTimeoutSec, false, nil)

	lines := strings.Split(resp, "\n")
	assert.Equal(t, "Signature: 8cd115d4637660341fa95066ea93945043221710f4ccc2ed1e75e88a80a4fe17", lines[len(lines)-5])
	assert.Equal(t, "Word list: offload scavenger backfield souvenir flatfoot impetus facial confidence billiard passenger drumbeat gossamer Trojan molasses Pluto embezzle crucial candidate banjo autopsy upshot revolver snapshot unify berserk impartial trauma maverick merit Pandora woodlark bookseller", lines[len(lines)-4])

	// Get the home directory
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	fh := fileutil.NewFileUrlHelper(mocks.Fs)
	acBuf, err := fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
	assert.NoError(t, err)

	yAc, err := util.NewYamlTestUtil(string(acBuf))
	yAc.AssertEqual(t, 1, "version")
	yAc.AssertEqual(t, "SleepySleep", "default_env")
	yAc.AssertEqual(t, "SleepySleep", "env_def[0].name")
	yAc.AssertEqual(t, "file://~/sleepcf/env-def.yaml", "env_def[0].env_def_file")
	yAc.AssertEqual(t, []interface{}{}, "disabled_commands")
}

func TestConfigImportEmptySigReject(t *testing.T) {
	cmd := []string{"swiz", "config", "import", ""}
	expect := []*util.ExpectResponse{
		{
			Output:   "Paste the app config here \n",
			Response: "dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vfi9zbGVlcGNmL2Vudi1kZWYueWFtbApkaXNhYmxlZF9jb21tYW5kczogW10K",
		},
		{
			Output:   "Paste the signature here (or leave blank if you don't know) \n",
			Response: "",
		},
		{
			Output:   "Does one or the other match? \n",
			Response: "n",
		},
		{
			Output: "Verify one or the other matches the signature provided with the app config\n",
			Action: util.NoAction,
		},
		{
			Output: "Signature: 8cd115d4637660341fa95066ea93945043221710f4ccc2ed1e75e88a80a4fe17\n",
			Action: util.NoAction,
		},
		{
			Output: "Word list: offload scavenger backfield souvenir flatfoot impetus facial confidence billiard passenger drumbeat gossamer Trojan molasses Pluto embezzle crucial candidate banjo autopsy upshot revolver snapshot unify berserk impartial trauma maverick merit Pandora woodlark bookseller\n",
			Action: util.NoAction,
		},
		{
			Output: "",
			Action: util.NoAction,
		},
	}
	mocks, _ := util.RunCommandWithMocks(t, cmd, expect, util.DefaultCmdTimeoutSec, true, nil)

	// Get the home directory
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	fh := fileutil.NewFileUrlHelper(mocks.Fs)
	_, err = fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
	assert.Error(t, err)
}

func TestConfigEmptySigIncorrect(t *testing.T) {
	cmd := []string{"swiz", "config", "import", ""}
	expect := []*util.ExpectResponse{
		{
			Output:   "Paste the app config here \n",
			Response: "dmVyc2lvbjogMQpkZWZhdWx0X2VudjogU2xlZXB5U2xlZXAKZW52X2RlZjoKICAgIC0gbmFtZTogU2xlZXB5U2xlZXAKICAgICAgZW52X2RlZl9maWxlOiBmaWxlOi8vfi9zbGVlcGNmL2Vudi1kZWYueWFtbApkaXNhYmxlZF9jb21tYW5kczogW10K",
		},
		{
			Output:   "Paste the signature here (or leave blank if you don't know) \n",
			Response: "AllYourBaseAreBelongToUs",
		},
	}
	mocks, _ := util.RunCommandWithMocks(t, cmd, expect, util.DefaultCmdTimeoutSec, true, nil)

	// Get the home directory
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	fh := fileutil.NewFileUrlHelper(mocks.Fs)
	_, err = fh.OpenUrlWithBaseDir(homeDir, "file://.swiz/app-config.yaml")
	assert.Error(t, err)
}
