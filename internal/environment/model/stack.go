package model

import "fmt"

const (
	StackKeyEnvName    = "SwzEnv"
	StackKeyCreateDate = "SwzCreateDate"
	StackKeyCreateUser = "SwzCreateUser"
	StackKeyEnvDef     = "SwzEnvDef"
	StackKeyEnclave    = "SwzEnclave"
)

type StackConfig struct {
	Version      int               `yaml:"version"`
	Name         string            `yaml:"-"`
	RawName      string            `yaml:"-"`
	Order        int               `yaml:"-"`
	Parameters   map[string]string `yaml:"params"`
	TemplateFile string            `yaml:"template_file"`
}

type StackInfo struct {
	Name         string
	NextAction   NextAction
	DeployStatus DeployStatus
	Resources    []string
}

func GenerateStackConfig(name string, templateFile string, params map[string]string) StackConfig {
	defaultParams := map[string]string{}
	for k, _ := range params {
		defaultParams[k] = fmt.Sprintf("{{%v}}", k)
	}

	return StackConfig{
		Version:      1,
		Name:         name,
		RawName:      name,
		Order:        1,
		Parameters:   defaultParams,
		TemplateFile: templateFile,
	}
}
