package osx

import (
	"bytes"
	"getswizzle.io/swiz/pkg/client/model"
	"getswizzle.io/swiz/pkg/exechelper"
	"log"
	"text/template"
)

// launchSsh launches an SSH app
func launchSsh(profile model.RemoteLaunchProfile, exec exechelper.ExecHelper) error {
	// Generate template
	// TODO: Sanitize template profile
	//tmpl, err := template.New("osxlaunch").Parse("/usr/bin/ssh -i {{.Keyfile}} {{.Username}}@127.0.0.1 -p {{.TunnelPort}}")
	tmpl, err := template.New("sshlaunch").Parse("ssh://{{.Username}}@127.0.0.1:{{.TunnelPort}}?keyfile={{.Keyfile}}")
	if err != nil {
		return err
	}

	var templBuf bytes.Buffer
	err = tmpl.Execute(&templBuf, profile)

	// Create command and redirect output
	param := templBuf.String()

	log.Printf("[DEBUG] starting ssh with connection %s\n", param)
	return exec.RunShellCmd("/usr/bin/", "open", "-n", "-F", "-W", "-a", "Terminal", param)
}
