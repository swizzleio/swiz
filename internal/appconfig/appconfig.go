package appconfig

import (
	"encoding/base64"
	"fmt"
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

	return &cfg, nil
}

func Generate(enclave EnclaveDef, env EnvDef) error {
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
		EnclaveDefinition: []EnclaveDef{enclave},
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

func (a AppConfig) GetBase64() (b64 string, sig string, err error) {
	// Marshal YAML into StackConfig
	var out []byte
	out, err = yaml.Marshal(a)
	if err != nil {
		return "", "", err
	}

	b64 = base64.StdEncoding.EncodeToString(out)
	sig = security.GetWordList(b64)

	return b64, sig, nil
}
