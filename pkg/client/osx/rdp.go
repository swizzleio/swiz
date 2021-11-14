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

// File format: https://docs.microsoft.com/en-us/windows-server/remote/remote-desktop-services/clients/rdp-files
const rdpTemplate = `targetisaadjoined:i:0
hubdiscoverygeourl:s:
redirected video capture encoding quality:i:0
camerastoredirect:s:
gatewaybrokeringtype:i:0
use redirection server name:i:0
alternate shell:s:
disable themes:i:0
disable cursor setting:i:1
remoteapplicationname:s:
resourceprovider:s:
disable menu anims:i:1
remoteapplicationcmdline:s:
promptcredentialonce:i:0
gatewaycertificatelogonauthority:s:
audiocapturemode:i:0
prompt for credentials on client:i:0
gatewayhostname:s:
remoteapplicationprogram:s:
gatewayusagemethod:i:2
screen mode id:i:1
use multimon:i:0
authentication level:i:2
desktopwidth:i:{{.Width}}
desktopheight:i:{{.Height}}
redirectsmartcards:i:0
redirectclipboard:i:1
forcehidpioptimizations:i:0
full address:s:127.0.0.1:{{.TunnelPort}}
drivestoredirect:s:
loadbalanceinfo:s:
networkautodetect:i:1
enablecredsspsupport:i:1
redirectprinters:i:0
autoreconnection enabled:i:1
session bpp:i:32
administrative session:i:0
audiomode:i:0
bandwidthautodetect:i:1
authoring tool:s:
connection type:i:7
remoteapplicationmode:i:0
disable full window drag:i:0
gatewayusername:s:
shell working directory:s:
wvd endpoint pool:s:
remoteapplicationappid:s:
username:s:{{.Username}}
allow font smoothing:i:1
connect to console:i:0
disable wallpaper:i:0
gatewayaccesstoken:s:
`

// genRdpFileString returns a generated rdp file string from the launch profile
func genRdpFileString(profile model.RemoteLaunchProfile) (string, error) {
	tmpl, err := template.New("rdplaunch").Parse(rdpTemplate)
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
func (c OsxClient) launchRdp(profile model.RemoteLaunchProfile) error {
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

		err = c.fs.WriteString(filename, rdpStr, 0600)

		return c.exec.RunShellCmd("", "command", "open", "-n", "-F", "-W", "-a", "/Applications/Microsoft Remote Desktop.app", filename)
	})
}
