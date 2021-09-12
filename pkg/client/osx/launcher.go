package osx

import (
	"bytes"
	"getswizzle.io/swiz/pkg/client/model"
	"getswizzle.io/swiz/pkg/common"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type OsxClient struct {
}

// Launch launches a client based on the specified launch profile
func (c OsxClient) Launch(profile model.RemoteLaunchProfile) error {
	switch strings.ToLower(profile.Appname) {
	case "rdp":
		return c.launchOsx(profile, "/Applications/Microsoft Remote Desktop Beta.app", "rdp://full%20address=s:127.0.0.1:{{.Port}}&audiomode=i:0&disable%20themes=i:1&desktopwidth:i:{{.Width}}&desktopheight:i:{{.Height}}&screen%20mode%20id=i:1&username=s:{{.Username}}&prompt%20for%20credentials%20on%20client:i:0")
	case "ssh":
		return c.launchOsx(profile, "ssh", "-i {{.Keyfile}} {{.Username}}@127.0.0.1:{{.Port}}")
	}

	return common.NotSupportedError{Subject: profile.Appname}
}

// launchOsx launches a command in OSX
func (c OsxClient) launchOsx(profile model.RemoteLaunchProfile, appName string, templateStr string) error {
	// Generate template
	// TODO: Sanitize template profile (specifically username)
	tmpl, err := template.New("osxlaunch").Parse(templateStr)
	if err != nil {
		return err
	}

	var templBuf bytes.Buffer
	err = tmpl.Execute(&templBuf, profile)

	// Create command and redirect output
	param := templBuf.String()
	cmd := exec.Command("command", "open", "-n", "-F", "-W", "-a", appName, param)
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
