package model

type RemoteLaunchProfile struct {
	Appname  string
	Port     int
	Username string
	Keyfile  string
	Width    int
	Height   int
}

type ClientLauncher interface {
	Launch(profile RemoteLaunchProfile) error
}
