package client

import (
	"getswizzle.io/swiz/pkg/client/model"
	"getswizzle.io/swiz/pkg/client/osx"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/fshelper"
	"strings"
)

type Servicer interface {
	Launch(os string, profile model.RemoteLaunchProfile) error
}

type Service struct {
	clients map[string]model.ClientLauncher
}

func NewService() Servicer {
	svc := &Service{
		clients: map[string]model.ClientLauncher{},
	}

	fs := fshelper.NewFsHelper()

	svc.clients[common.OsOsx] = osx.NewOsxClient(fs)

	return svc
}

func (s Service) Launch(os string, profile model.RemoteLaunchProfile) error {
	client := s.clients[strings.ToLower(os)]
	if client == nil {
		return common.NotSupportedError{Subject: os}
	}

	return client.Launch(profile)
}
