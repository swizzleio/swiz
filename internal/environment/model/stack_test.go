package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack_GenerateStackConfig(t *testing.T) {
	params := map[string]string{
		"blah": "stuff",
		"boo":  "ahhh",
	}
	cfg := GenerateStackConfig("neato", "neato-cfg.yaml", params)

	assert.Equal(t, 1, cfg.Version)
	assert.Equal(t, "neato", cfg.Name)
	assert.Equal(t, "neato", cfg.RawName)
	assert.Equal(t, 1, cfg.Order)
	assert.Equal(t, "neato-cfg.yaml", cfg.TemplateFile)
	assert.Equal(t, "{{blah}}", cfg.Parameters["blah"])
	assert.Equal(t, "{{boo}}", cfg.Parameters["boo"])
}
