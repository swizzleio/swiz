package model

type StackConfigDef struct {
	Name       string `yaml:"name"`
	ConfigFile string `yaml:"config_file"`
	Order      int    `yaml:"order"`
}

type EnvironmentConfig struct {
	EnvDefName        string
	Version           int              `yaml:"version"`
	DefaultEnclave    string           `yaml:"default_enclave"`
	EnclaveDefinition []Enclave        `yaml:"enclave_def"`
	StackCfgDef       []StackConfigDef `yaml:"stack_cfg"`
	Stacks            map[string]*StackConfig
}

type EnvironmentInfo struct {
	EnvironmentName   string
	DeployStatus      DeployStatus
	StackDeployStatus []DeployStatus
}
