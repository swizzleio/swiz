package model

type StackConfig struct {
	Version      int `yaml:"version"`
	Name         string
	Parameters   map[string]interface{} `yaml:"params"`
	TemplateFile string                 `yaml:"template_file"`
}
