package model

const (
	StackKeyEnvName    = "SwzEnv"
	StackKeyCreateDate = "SwzCreateDate"
	StackKeyCreateUser = "SwzCreateUser"
	StackKeyEnvDef     = "SwzEnvDef"
	StackKeyEnclave    = "SwzEnclave"
)

type StackConfig struct {
	Version      int `yaml:"version"`
	Name         string
	Order        int
	Parameters   map[string]string `yaml:"params"`
	TemplateFile string            `yaml:"template_file"`
}

type StackInfo struct {
	Name         string
	NextAction   NextAction
	DeployStatus DeployStatus
	Resources    []string
}
