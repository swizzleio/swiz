package model

type RemoteLaunchProfile struct {
	Appname      string
	Os           string
	InstanceName string
	Port         int
	TunnelPort   int
	Username     string
	Keyfile      string
	Width        int
	Height       int
}

type ClientLauncher interface {
	Launch(profile RemoteLaunchProfile) error
}
