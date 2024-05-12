//go:build functionaltest

package cmds

import (
	"github.com/spf13/afero"
	"github.com/swizzleio/swiz/internal/appconfig"
	appcli "github.com/swizzleio/swiz/pkg/cli"
)

// SetupFixtures configures the app for a functional test environment but mocks things
// like the OS and provides CLI output that is simplified for testability
func SetupFixtures() afero.Fs {
	cl = appcli.NewCli(appcli.NewDumbSurvey(nil, nil))
	appFs = afero.NewMemMapFs()
	appConfigMgr = appconfig.NewManage(appFs)

	return appFs
}
