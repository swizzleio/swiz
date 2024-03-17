//go:build functionaltest

package functional

import (
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/swizzleio/swiz/cmd/cmds"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"swiz": cmds.Execute,
	}))
}

func Test_Version(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "scripts/version",
	})
}

func Test_Config(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "scripts/config",
	})
}
