package model

type EnvironmentConfigDef struct {
	Parameters map[string]string `yaml:"params"`
}

type StackConfigDef struct {
	Name       string `yaml:"name"`
	ConfigFile string `yaml:"config_file"`
	Order      int    `yaml:"order"`
}

type EnvironmentConfig struct {
	EnvDefName  string
	Version     int                             `yaml:"version"`
	Config      map[string]EnvironmentConfigDef `yaml:"config"`
	StackCfgDef []StackConfigDef                `yaml:"stack_cfg"`
	Stacks      map[string]*StackConfig
}

type EnvironmentInfo struct {
	EnvironmentName   string
	DeployStatus      DeployStatus
	StackDeployStatus []DeployStatus
}
