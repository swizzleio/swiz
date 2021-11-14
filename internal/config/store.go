package config

import (
	clientmodel "getswizzle.io/swiz/pkg/client/model"
	"getswizzle.io/swiz/pkg/common"
	"getswizzle.io/swiz/pkg/fshelper"
	inframodel "getswizzle.io/swiz/pkg/infra/model"
	"getswizzle.io/swiz/pkg/network"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	defaultHeight   = 900
	defaultWidth    = 1200
	defaultFileName = "/.swiz/cfg.json"
)

type Storer interface {
	Load(filename string) (Config, error)
	MustLoad(filename string) Config
	Save(filename string) error
	GenerateDefaults(linuxFlavor string, linuxHostKey string, bastionHosts []Bastion)
	GetHostLaunchInfo(host inframodel.TargetInstance) (*HostLaunchInfo, error)
	Get() Config
}

type Store struct {
	fs  fshelper.FsHelper
	cfg Config
}

func NewConfig() Storer {
	return &Store{
		fs: fshelper.NewFsHelper(),
		cfg: Config{
			Version:       1,
			Os:            runtime.GOOS,
			BastionHosts:  []Bastion{},
			LaunchProfile: map[string]*LaunchProfile{},
		},
	}
}

func NewConfigMustLoad(filename string) Storer {
	store := &Store{
		fs: fshelper.NewFsHelper(),
		cfg: Config{
			Version:       1,
			Os:            runtime.GOOS,
			BastionHosts:  []Bastion{},
			LaunchProfile: map[string]*LaunchProfile{},
		},
	}

	_ = store.MustLoad(filename)
	return store
}

// LoadConfig loads the application config file and sets any embedded config items
func (c *Store) Load(filename string) (Config, error) {
	loadFilename, err := c.filenameOrDefault(filename)
	if nil != err {
		return c.Get(), err
	}

	err = c.fs.ReadJson(loadFilename, &c.cfg)

	return c.Get(), err
}

// LoadConfigMust loads the config and panics if it could not be loaded
func (c *Store) MustLoad(filename string) Config {
	cfg, err := c.Load(filename)

	if err != nil {
		log.Fatalf("loading config. %v", err)
	}

	return cfg
}

// Save saves the app config
func (c Store) Save(filename string) error {
	saveFilename, err := c.filenameOrDefault(filename)
	if nil != err {
		return err
	}

	return c.fs.WriteJson(saveFilename, c.cfg)
}

// GenerateDefaults resets the config file
func (c *Store) GenerateDefaults(linuxFlavor string, linuxHostKey string, bastionHosts []Bastion) {
	// Reset values
	c.cfg.LaunchProfile = map[string]*LaunchProfile{}
	c.cfg.BastionHosts = bastionHosts

	// Update bastion ids if they are empty
	for i, bastion := range c.cfg.BastionHosts {
		_, err := uuid.Parse(bastion.Id)
		if err != nil {
			c.cfg.BastionHosts[i].Id = uuid.New().String()
		}
	}

	// Determine linux username
	linuxUser := "root"
	switch linuxFlavor {
	case common.LinuxUbuntu:
		linuxUser = "ubuntu"
	case common.LinuxAmazon:
		linuxUser = "ec2-user"
	}

	// Init os values
	c.cfg.LaunchProfile[common.OsLinux] = &LaunchProfile{
		DefaultApp:  common.RemoteAccessSsh,
		DefaultUser: linuxUser,
		DefaultPort: 22,
		Auth: AuthInfo{
			KeyFilename: linuxHostKey,
		},
		Width:  0,
		Height: 0,
	}
	c.cfg.LaunchProfile[common.OsWindows] = &LaunchProfile{
		DefaultApp:  common.RemoteAccessRdp,
		DefaultUser: "Administrator",
		DefaultPort: 3389,
		Auth:        AuthInfo{},
		Width:       defaultWidth,
		Height:      defaultHeight,
	}
}

type HostLaunchInfo struct {
	ClientConfig     clientmodel.RemoteLaunchProfile
	HostString       string
	Os               string
	BastionAddr      string
	BastionSignature string
	BastionAuth      AuthInfo
}

// GetHostLaunchInfo returns all of the information needed to launch a tunnel
// TODO: This needs to be moved into it's own service
func (c Store) GetHostLaunchInfo(host inframodel.TargetInstance) (*HostLaunchInfo, error) {

	// Get launch profile
	launchProfile := c.cfg.LaunchProfile[strings.ToLower(host.Os)]
	if launchProfile == nil {
		return nil, common.NotFoundError{
			Subject: host.Os,
		}
	}

	// Fetch private endpoint
	var hostEndpoint network.Endpoint
	for _, e := range host.Endpoints {
		if e.IsPrivate() {
			hostEndpoint = e
			break
		}
	}
	if hostEndpoint.User == "" {
		hostEndpoint.User = launchProfile.DefaultUser
	}
	if hostEndpoint.Port == 0 {
		hostEndpoint.Port = launchProfile.DefaultPort
	}

	// Fetch random bastion host
	randomIndex := rand.Intn(len(c.cfg.BastionHosts))
	bastionHost := c.cfg.BastionHosts[randomIndex]

	return &HostLaunchInfo{
		HostString:       hostEndpoint.String(),
		BastionAddr:      bastionHost.Addr,
		BastionSignature: bastionHost.Signature,
		BastionAuth:      bastionHost.BastionAuth,
		Os:               c.cfg.Os,
		ClientConfig: clientmodel.RemoteLaunchProfile{
			Appname:  launchProfile.DefaultApp,
			Port:     launchProfile.DefaultPort,
			Username: hostEndpoint.User,
			Keyfile:  launchProfile.Auth.KeyFilename,
			Width:    launchProfile.Width,
			Height:   launchProfile.Height,
		},
	}, nil
}

// Get returns the app config
func (c Store) Get() Config {
	return c.cfg
}

// filenameOrDefault returns the filename or, if empty, a default path
func (Store) filenameOrDefault(filename string) (string, error) {
	if filename == "" {
		// Default is the user's home directory
		dirname, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return path.Join(dirname, defaultFileName), nil
	}

	return filename, nil
}
