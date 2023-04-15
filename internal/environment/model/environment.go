package model

type EnvironmentConfig struct {
	EnvDefName string
	Version    int `yaml:"version"`
	Config     []struct {
		Enclave    string                 `yaml:"enclave"`
		Parameters map[string]interface{} `yaml:"params"`
	} `yaml:"config"`
	StackCfgDef []struct {
		Name       string `yaml:"name"`
		ConfigFile string `yaml:"config_file"`
	} `yaml:"stack_cfg"`
	Stacks map[string]*StackConfig
}
