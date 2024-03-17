package model

import (
	"fmt"
	"github.com/swizzleio/swiz/internal/appconfig"
)

const (
	DefaultNamingScheme = "{{env_name:32}}-{{stack_name:32}}"
	DefaultEnclaveName  = "NameMe"
)

type StackConfigDef struct {
	Name       string `yaml:"name"`
	ConfigFile string `yaml:"config_file"`
	Order      int    `yaml:"order"`
}

type EnvironmentConfig struct {
	EnvDefName        string                  `yaml:"-"`
	Version           int                     `yaml:"version"`
	DefaultEnclave    string                  `yaml:"default_enclave"`
	NamingScheme      string                  `yaml:"naming_scheme"`
	EnclaveDefinition []Enclave               `yaml:"enclave_def"`
	StackCfgDef       []StackConfigDef        `yaml:"stack_cfg"`
	Stacks            map[string]*StackConfig `yaml:"-"`
}

type EnvironmentInfo struct {
	EnvironmentName string
	DeployStatus    DeployStatus
	StackInfo       []StackInfo
}

func GenerateFileName(stackName string) string {
	return fmt.Sprintf("%v/%v-cfg.yaml", appconfig.DefaultOutLocation, stackName)
}

func GenerateEnvironmentConfig(stacks []StackConfig, enclaves []Enclave, defaultEnclave string) EnvironmentConfig {
	stackConfigDef := []StackConfigDef{}
	for _, stack := range stacks {
		stackConfigDef = append(stackConfigDef, StackConfigDef{
			Name:       stack.Name,
			ConfigFile: GenerateFileName(stack.Name),
			Order:      stack.Order,
		})
	}

	return EnvironmentConfig{
		Version:           1,
		DefaultEnclave:    defaultEnclave,
		NamingScheme:      DefaultNamingScheme,
		EnclaveDefinition: enclaves,
		StackCfgDef:       stackConfigDef,
	}
}
