package osx

import (
	"bytes"
	"fmt"
	"getswizzle.io/swiz/pkg/client/model"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/fshelper"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

type OsxClient struct {
	fs fshelper.FsHelper
}

// NewOsxClient returns a new OSX client
func NewOsxClient(fs fshelper.FsHelper) model.ClientLauncher {
	return &OsxClient{
		fs: fs,
	}
}

// Launch launches a client based on the specified launch profile
func (c OsxClient) Launch(profile model.RemoteLaunchProfile) error {
	switch strings.ToLower(profile.Appname) {
	case common.RemoteAccessRdp:
		return c.fs.RunInTempDir(false, func(tmpPath string) error {
			// Generate file name
			randuuid, err := uuid.NewUUID()
			if err != nil {
				return err
			}

			filename := path.Join(tmpPath, fmt.Sprintf("%s.rdp", randuuid))
			log.Printf("[DEBUG] creating file %v", filename)
			rdpStr, err := genRdpFileString(profile)
			if err != nil {
				return err
			}

			err = c.fs.WriteString(filename, rdpStr)

			return c.launchOsx("/Applications/Microsoft Remote Desktop.app", filename)
		})

	case common.RemoteAccessSsh:
		// Generate template
		// TODO: Sanitize template profile
		tmpl, err := template.New("osxlaunch").Parse("-i {{.Keyfile}} {{.Username}}@127.0.0.1:{{.Port}}")
		if err != nil {
			return err
		}

		var templBuf bytes.Buffer
		err = tmpl.Execute(&templBuf, profile)

		// Create command and redirect output
		param := templBuf.String()

		return c.launchOsx("ssh", param)
	}

	return common.NotSupportedError{Subject: profile.Appname}
}

// launchOsx launches a command in OSX
func (c OsxClient) launchOsx(cmdApp string, param string) error {

	//cmd := exec.Command("command", "open", "-n", "-F", "-W", "-a", appParams, param)
	cmd := exec.Command("command", "open", "-n", "-F", "-W", cmdApp, param)
	stderr, err := cmd.StderrPipe()
	log.SetOutput(os.Stderr)
	if err != nil {
		return err
	}

	// Launch command
	if err = cmd.Start(); err != nil {
		return err
	}

	// Read output and wait
	// TODO: Add cancel
	slurp, _ := ioutil.ReadAll(stderr)
	log.Printf("%v\n", slurp)

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
