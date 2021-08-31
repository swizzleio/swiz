package osx

import (
	"bytes"
	"getswizzle.io/swiz/pkg/client"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"
)

func Launch(profile client.RemoteLaunchProfile, templateStr string) error {
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
	cmd := exec.Command("command", "open", "-n", "-F", "-W", "-a", param)
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
