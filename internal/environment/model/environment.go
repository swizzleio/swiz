package model

type EnvironmentConfig struct {
	Version int `yaml:"version"`
	Config  []struct {
		Enclave    string                 `yaml:"enclave"`
		Parameters map[string]interface{} `yaml:"params"`
	} `yaml:"config"`
	Stacks []struct {
		Name       string `yaml:"name"`
		ConfigFile string `yaml:"config_file"`
	}
}
