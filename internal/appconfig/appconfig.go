package appconfig

import (
	"encoding/base64"
	"fmt"
	"github.com/swizzleio/swiz/internal/environment/model"
	"github.com/swizzleio/swiz/pkg/configutil"
	"github.com/swizzleio/swiz/pkg/drivers/awswrap"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"github.com/swizzleio/swiz/pkg/security"
	"gopkg.in/yaml.v3"
)

var DefaultLocation = "file://~/.swiz/appconfig.yaml"

var ProviderIds = []string{"aws"}

type EnvDef struct {
	Name       string `yaml:"name"`
	EnvDefFile string `yaml:"env_def_file"`
	Default    bool   `yaml:"default"`
}

type AppConfig struct {
	Version           int             `yaml:"version"`
	EnvDefinition     []EnvDef        `yaml:"env_def"`
	EnclaveDefinition []model.Enclave `yaml:"enclave_def"`
	DefaultEnclave    string          `yaml:"default_enclave"`
	DisabledCommands  []string        `yaml:"disabled_commands"`
	BaseDir           string
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

	cfg.BaseDir, err = fileutil.GetPathFromUrl(location)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Generate(enclave model.Enclave, env EnvDef) error {
	err := fileutil.CreateDirIfNotExist(DefaultLocation)
	if err != nil {
		return err
	}

	// Set defaults if not set
	enclave.Name = configutil.SetOrDefault[string](enclave.Name, awswrap.DefaultAccountName)
	enclave.ProviderId = configutil.SetOrDefault[string](enclave.ProviderId, ProviderIds[0])
	enclave.Region = configutil.SetOrDefault[string](enclave.Name, awswrap.DefaultRegion)

	env.Name = configutil.SetOrDefault[string](env.Name, "default")
	env.Default = configutil.SetOrDefault[bool](env.Default, true)

	// Check for mandatory fields
	if enclave.AccountId == "" || env.EnvDefFile == "" {
		return fmt.Errorf("missing mandatory fields of account id or env def file")
	}

	// Save app config to yaml
	cfg := AppConfig{
		Version:           1,
		EnvDefinition:     []EnvDef{env},
		EnclaveDefinition: []model.Enclave{enclave},
	}

	return fileutil.YamlToLocation(DefaultLocation, cfg)
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

	// TODO, in the command line, verify signature
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
