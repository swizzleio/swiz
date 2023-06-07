package appconfig

import (
	"encoding/base64"
	"fmt"
	"github.com/swizzleio/swiz/pkg/configutil"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"github.com/swizzleio/swiz/pkg/security"
	"gopkg.in/yaml.v3"
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

type Base64Resp struct {
	Encoded   string
	WordList  string
	Signature string
}

func Parse(location string) (*AppConfig, error) {

	if location == "" {
		location = DefaultLocation
	}

	// Open URL
	data, err := fileutil.OpenUrl(location)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into AppConfig
	cfg := AppConfig{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	cfg.BaseDir, err = fileutil.GetPathFromUrl(location, false)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Generate(env EnvDef) AppConfig {
	// Set defaults if not set
	env.Name = configutil.SetOrDefault[string](env.Name, "default")

	// Save app config to yaml
	return AppConfig{
		Version:          1,
		DefaultEnv:       env.Name,
		EnvDefinition:    []EnvDef{env},
		DisabledCommands: []string{},
	}
}

func Fetch(data string) error {
	err := fileutil.CreateDirIfNotExist(DefaultLocation)
	if err != nil {
		return err
	}

	// Decode base64
	b64 := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(b64, []byte(data))
	if err != nil {
		return err
	}

	return fileutil.WriteUrl(DefaultLocation, b64[:n])
}

func (a AppConfig) GetBase64() (*Base64Resp, error) {
	// Marshal YAML
	out, err := yaml.Marshal(a)
	if err != nil {
		return nil, err
	}

	// Get encoding and signature
	retVal := &Base64Resp{}
	retVal.Encoded = base64.StdEncoding.EncodeToString(out)
	retVal.Signature, retVal.WordList = security.GetSha256AndWordList(retVal.Encoded)

	return retVal, nil
}
