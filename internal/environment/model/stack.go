package model

import (
	"github.com/swizzleio/swiz/pkg/fileutil"
	"gopkg.in/yaml.v3"
)

type StackConfig struct {
	Version      int                    `yaml:"version"`
	Parameters   map[string]interface{} `yaml:"params"`
	TemplateFile string                 `yaml:"template_file"`
}

func NewFromConfig(cfg StackConfig) *Stack {
	return &Stack{}
}

func NewFromLocation(location string) (*Stack, error) {

	// Open URL
	data, err := fileutil.OpenUrl(location)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into StackConfig
	cfg := StackConfig{}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return NewFromConfig(cfg), nil
}

type Stack struct {
}

func (s *Stack) Create() {
}

func (s *Stack) Update() {
}

func (s *Stack) Delete() {
}

func (s *Stack) GetOutput() {
}
