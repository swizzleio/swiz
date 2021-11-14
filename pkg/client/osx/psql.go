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

const psqlTemplate = `#!/bin/sh
psql -h 127.0.0.1 -p {{.TunnelPort}} -U {{.Username}}
`

// genPsqlFileString returns a generated script string from a launch profile
func genPsqlFileString(profile model.RemoteLaunchProfile) (string, error) {
	tmpl, err := template.New("psqlaunch").Parse(psqlTemplate)
	if err != nil {
		return "", err
	}

	var templBuf bytes.Buffer
	err = tmpl.Execute(&templBuf, profile)

	// Create template string
	rdpStr := templBuf.String()
	return rdpStr, nil
}

// launchRdp launches an rdp app
func (c OsxClient) launchPgSql(profile model.RemoteLaunchProfile) error {
	return c.fs.RunInTempDir(false, func(tmpPath string) error {
		// Generate file name
		randuuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		filename := path.Join(tmpPath, fmt.Sprintf("%s.sh", randuuid))
		log.Printf("[DEBUG] creating file %v", filename)
		rdpStr, err := genPsqlFileString(profile)
		if err != nil {
			return err
		}

		err = c.fs.WriteString(filename, rdpStr, 0700)

		return c.exec.RunShellCmd("", "open", "-n", "-F", "-W", "-a", "Terminal", filename)
	})
}
