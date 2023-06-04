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
	RawName      string
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

func GenerateStackConfig(name string, templateFile string) StackConfig {
	return StackConfig{
		Version:      1,
		Name:         name,
		RawName:      name,
		Order:        1,
		Parameters:   map[string]string{},
		TemplateFile: templateFile,
	}
}
