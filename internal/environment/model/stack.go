package model

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
}
