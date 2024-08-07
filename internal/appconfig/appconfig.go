package appconfig

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/swizzleio/swiz/pkg/configutil"
	"github.com/swizzleio/swiz/pkg/fileutil"
)

var DefaultFileName = "app-config.yaml"
var DefaultSwizDir = "file://~/.swiz"
var DefaultLocation = fmt.Sprintf("%v/%v", DefaultSwizDir, DefaultFileName) // Don't use filepath.Join here, it's not URI aware
var DefaultOutLocation = "file://./out"

type EnvDef struct {
	Name       string `yaml:"name"`
	EnvDefFile string `yaml:"env_def_file"`
}

type AppConfig struct {
	Version          int      `yaml:"version"`
	DefaultEnv       string   `yaml:"default_env"`
	EnvDefinition    []EnvDef `yaml:"env_def"`
	DisabledCommands []string `yaml:"disabled_commands"`
	BaseDir          string   `yaml:"base_dir"`
}

type Manage struct {
	appFs    afero.Fs
	ser      fileutil.SerializeHelper[AppConfig]
	fh       fileutil.FileHelper
	isLoaded bool
}

func NewManage(appFs afero.Fs) *Manage {
	return &Manage{
		appFs: appFs,
		ser:   fileutil.NewYamlHelper[AppConfig](appFs),
		fh:    fileutil.NewFileHelper(appFs),
	}
}

func (a *Manage) GenFromEnv(env EnvDef) *AppConfig {
	// Set defaults if not set
	env.Name = configutil.SetOrDefault[string](env.Name, "default")

	// Save app config to yaml
	cfg := &AppConfig{
		Version:          1,
		DefaultEnv:       env.Name,
		EnvDefinition:    []EnvDef{env},
		DisabledCommands: []string{},
	}

	a.ser.Set(*cfg)

	return cfg
}

// GenFromB64 generates an app config from a base64 string
func (a *Manage) GenFromB64(data string, save bool) error {

	err := a.ser.SetFromB64(data)
	if err != nil {
		return err
	}

	a.isLoaded = true

	if save {
		err = a.fh.CreateDirIfNotExist(DefaultLocation)
		if err != nil {
			return err
		}

		return a.ser.Save(DefaultLocation)
	}

	return nil
}

func (a *Manage) Load(location string) (*AppConfig, error) {
	isCustomLoc := true
	if location == "" {
		location = DefaultLocation
		isCustomLoc = false
	}

	// Open Yaml
	cfg, err := a.ser.Open(location)
	if err != nil {
		return nil, err
	}

	// Check to see if this is a custom location and there is no BaseDir override, if that's the case, assume the base
	// directory is the location of this file. This handles the case of a monorepo with the app-config.yaml being
	// custom defined
	if isCustomLoc &&
		cfg.BaseDir == "" {

		openUrl := fileutil.NewFileUrlHelper(a.appFs)
		cfg.BaseDir, err = openUrl.GetPathFromUrl(location, false)

		if err != nil {
			return nil, err
		}

		a.ser.Set(*cfg)
	}

	a.isLoaded = true

	return cfg, nil
}

// IsLoaded returns true if the app config is loaded
func (a *Manage) IsLoaded() bool {
	return a.isLoaded
}

// Get returns the app config
func (a *Manage) Get() AppConfig {
	return a.ser.Get()
}

// GetBase64 returns the base64 signature of the app config
func (a *Manage) GetBase64() (*fileutil.Base64Resp, error) {
	// Return base64
	return a.ser.GetBase64()
}
