package osx

import (
	"getswizzle.io/swiz/pkg/client/model"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/exechelper"
	"getswizzle.io/swiz/pkg/fshelper"
	"strings"
)

type OsxClient struct {
	fs   fshelper.FsHelper
	exec exechelper.ExecHelper
}

// NewOsxClient returns a new OSX client
func NewOsxClient(fs fshelper.FsHelper, exec exechelper.ExecHelper) model.ClientLauncher {
	return &OsxClient{
		fs:   fs,
		exec: exec,
	}
}

// Launch launches a client based on the specified launch profile
func (c OsxClient) Launch(profile model.RemoteLaunchProfile) error {
	// Determine what to launch. Eventually this can possibly be more of a map similar to the aws system
	switch strings.ToLower(profile.Appname) {
	case common.RemoteAccessRdp:
		return c.launchRdp(profile)
	case common.RemoteAccessSsh:
		return c.launchSsh(profile)
	case common.RemoteAccessGuess:
		switch strings.ToLower(profile.Os) {
		case common.OsPgSql:
			return c.launchPgSql(profile)
		}
	}

	return common.NotSupportedError{Subject: profile.Appname}
}
