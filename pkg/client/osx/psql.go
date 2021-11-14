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

const datagripTemplate = `?xml version="1.0" encoding="UTF-8"?>
<project version="4">
  <component name="dataSourceStorageLocal" created-in="DB-212.5457.41">
    <data-source name="postgres@127.0.0.1:{{.TunnelPort}}" uuid="6cd6667b-104c-4e02-baa8-be3c6fc38a00">
      <database-info product="" version="" jdbc-version="" driver-name="" driver-version="" dbms="POSTGRES" exact-version="0" />
      <secret-storage>master_key</secret-storage>
      <user-name>postgres</user-name>
      <schema-mapping />
    </data-source>
  </component>
</project>
`

// genRdpFileString returns a generated rdp file string from the launch profile
func genDatagripFileString(profile model.RemoteLaunchProfile) (string, error) {
	tmpl, err := template.New("rdslaunch").Parse(rdpTemplate)
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
func (c OsxClient) launchDatagrip(profile model.RemoteLaunchProfile) error {
	return c.fs.RunInTempDir(false, func(tmpPath string) error {
		// Generate file name
		randuuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		filename := path.Join(tmpPath, fmt.Sprintf("%s.rdp", randuuid))
		log.Printf("[DEBUG] creating file %v", filename)
		rdpStr, err := genDatagripFileString(profile)
		if err != nil {
			return err
		}

		err = c.fs.WriteString(filename, rdpStr, 0600)

		return c.exec.RunShellCmd("", "command", "open", "-n", "-F", "-W", "-a", "/Applications/Microsoft Remote Desktop.app", filename)
	})
}
