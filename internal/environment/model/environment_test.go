package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment_GenerateFileName(t *testing.T) {
	name := GenerateFileName("neato")
	assert.Equal(t, "file://./out/neato-cfg.yaml", name)
}

func TestEnvironment_GenerateEnvironmentConfig(t *testing.T) {
	stacks := []StackConfig{
		{
			Name:  "neato",
			Order: 1,
		},
		{
			Name:  "meh",
			Order: 2,
		},
		{
			Name:  "something",
			Order: 2,
		},
	}

	enclaves := []Enclave{
		{
			Name: "dev",
		},
		{
			Name: "prod",
		},
	}

	cfg := GenerateEnvironmentConfig(stacks, enclaves, "dev")
	assert.Len(t, cfg.StackCfgDef, 3)
	assert.Equal(t, "meh", cfg.StackCfgDef[1].Name)
	assert.Equal(t, 2, cfg.StackCfgDef[1].Order)
	assert.Equal(t, "file://./out/meh-cfg.yaml", cfg.StackCfgDef[1].ConfigFile)
	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "dev", cfg.DefaultEnclave)
	assert.Equal(t, DefaultNamingScheme, cfg.NamingScheme)
	assert.Equal(t, enclaves, cfg.EnclaveDefinition)
}
