package appconfig

import (
	"fmt"
	"github.com/swizzleio/swiz/pkg/configutil"
	"github.com/swizzleio/swiz/pkg/fileutil"
)

var DefaultFileName = "app-config.yaml"
var DefaultLocation = fmt.Sprintf("file://~/.swiz/%v", DefaultFileName)
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
	BaseDir          string   `yaml:"-"`
}

type Manage struct {
	ser      fileutil.SerializeHelper[AppConfig]
	fh       fileutil.FileHelper
	isLoaded bool
}

func NewManage() *Manage {
	return &Manage{
		ser: fileutil.NewYamlHelper[AppConfig](),
		fh:  fileutil.NewFileHelper(),
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
	if location == "" {
		location = DefaultLocation
	}

	// Open Yaml
	cfg, err := a.ser.Open(location)
	if err != nil {
		return nil, err
	}

	openUrl := fileutil.NewFileUrlHelper()

	cfg.BaseDir, err = openUrl.GetPathFromUrl(location, false)
	if err != nil {
		return nil, err
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
