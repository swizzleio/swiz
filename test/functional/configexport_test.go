//go:build functionaltest

package functional

import (
	"fmt"
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
	fmt.Println(resp)
	lines := strings.Split(resp, "\n")
	assert.Equal(t, "App config: dmVyc2lvbjogMApkZWZhdWx0X2VudjogIiIKZW52X2RlZjogW10KZGlzYWJsZWRfY29tbWFuZHM6IFtdCg==", lines[len(lines)-4])
	assert.Equal(t, "Signature: b869bfa89aee048121b18f24e8f27fcacb358507afa2cfdfae5cdbb092e543be", lines[len(lines)-3])
	assert.Equal(t, "Word list: select guitarist slingshot paramount pupil universe adrift inventive blackjack photograph payday Capricorn trauma vagabond lockup revenue spheroid conformist music amusement rocker Pacific stagehand therapist robust fascinate suspense phonetic physique travesty crucial racketeer", lines[len(lines)-2])
}
