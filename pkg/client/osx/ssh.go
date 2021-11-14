package osx

import (
	"bytes"
	"fmt"
	"getswizzle.io/swiz/pkg/client/model"
	"github.com/google/uuid"
	"log"
	"path"
	"text/template"
)

const sshTemplate = `#!/bin/sh
ssh -i {{.Keyfile}} {{.Username}}@127.0.0.1 -p {{.TunnelPort}}
`

func genSshFileString(profile model.RemoteLaunchProfile) (string, error) {
	// TODO: These params need to be sanitized!
	tmpl, err := template.New("sshlaunch").Parse(sshTemplate)
	if err != nil {
		return "", err
	}

	var templBuf bytes.Buffer
	err = tmpl.Execute(&templBuf, profile)

	// Create template string
	rdpStr := templBuf.String()
	return rdpStr, nil
}

// launchSsh launches an SSH app
func (c OsxClient) launchSsh(profile model.RemoteLaunchProfile) error {
	return c.fs.RunInTempDir(false, func(tmpPath string) error {
		// Generate file name
		randuuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		filename := path.Join(tmpPath, fmt.Sprintf("%s.sh", randuuid))
		log.Printf("[DEBUG] creating file %v", filename)
		rdpStr, err := genSshFileString(profile)
		if err != nil {
			return err
		}

		err = c.fs.WriteString(filename, rdpStr, 0700)

		return c.exec.RunShellCmd("", "open", "-n", "-F", "-W", "-a", "Terminal", filename)
	})
}
