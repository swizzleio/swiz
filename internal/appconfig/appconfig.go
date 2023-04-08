package appconfig

import (
	"fmt"
	"github.com/swizzleio/swiz/pkg/configutil"
	"github.com/swizzleio/swiz/pkg/drivers/awswrap"
	"github.com/swizzleio/swiz/pkg/fileutil"
	"gopkg.in/yaml.v3"
	"os"
)

var DefaultLocaton = "file://~/.swiz/appconfig.yaml"

var ProviderIds = []string{"aws"}

type EnvDef struct {
	Name       string `yaml:"name"`
	EnvDefFile string `yaml:"env_def_file"`
	Default    bool   `yaml:"default"`
}

type EnclaveDef struct {
	Name       string `yaml:"name"`
	ProviderId string `yaml:"provider_id"`
	AccountId  string `yaml:"account_id"`
	Region     string `yaml:"region"`
	DomainName string `yaml:"domain_name"`
}

type AppConfig struct {
	Version           int          `yaml:"version"`
	EnvDefinition     []EnvDef     `yaml:"env_def"`
	EnclaveDefinition []EnclaveDef `yaml:"enclave_def"`
}

func Parse(location string) (*AppConfig, error) {

	if location == "" {
		location = DefaultLocaton
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

	return &cfg, nil
}

func Generate(enclave EnclaveDef, env EnvDef) error {
	if err := os.MkdirAll("~/.swiz", 0755); err != nil {
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
		EnclaveDefinition: []EnclaveDef{enclave},
	}

	return fileutil.YamlToLocation(DefaultLocaton, cfg)
}

func Fetch() {
	// Decode base64 into file or URI
}
