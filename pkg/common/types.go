package common

const (
	OsWindows = "windows"
	OsOsx     = "darwin"
	OsLinux   = "linux"
	OsPgSql   = "postgresql"
	OsMySql   = "mysql"
	OsMariaDb = "mariadb"
)

// Currently this is used for a future capability where the actual app can be overridden.
// Until then, this is a bit of ugly tech debt...
const (
	RemoteAccessGuess = "guess"
	RemoteAccessSsh   = "ssh"
	RemoteAccessRdp   = "rdp"
)

const (
	LinuxAmazon = "Amazon"
	LinuxCentos = "Centos"
	LinuxUbuntu = "Ubuntu"
	LinuxOther  = "Other"
)
