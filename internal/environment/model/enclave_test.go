package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swizzleio/swiz/pkg/drivers/awswrap"
)

func TestEnclave_GetProvider(t *testing.T) {
	enc := Enclave{
		DefaultProvider: "something",
		Providers: []EncProvider{
			{
				Name: "foobar",
			},
			{
				Name: "something",
			},
		},
	}
	tests := []struct {
		name     string
		provider string
		want     *EncProvider
	}{
		{
			name:     "default provider",
			provider: "",
			want: &EncProvider{
				Name: "something",
			},
		},
		{
			name:     "foobar provider",
			provider: "foobar",
			want: &EncProvider{
				Name: "foobar",
			},
		},
		{
			name:     "empty provider",
			provider: "nothing",
			want:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enc.GetProvider(tt.provider)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEnclave_ToAwsConfig(t *testing.T) {
	encProvider := EncProvider{
		Name:      "foobar",
		AccountId: "1234567890",
		Region:    "us-west-2",
	}

	cfger := encProvider.ToAwsConfig()
	cfg, ok := cfger.(*awswrap.AwsConfig)
	assert.True(t, ok)
	assert.Equal(t, "foobar", cfg.Profile)
	assert.Equal(t, "1234567890", cfg.AccountId)
	assert.Equal(t, "us-west-2", cfg.Region)
}

func TestEnclave_GenerateEnclave(t *testing.T) {
	cfg := awswrap.AwsConfig{
		Profile:   "foobar",
		AccountId: "1234567890",
		Region:    "us-west-2",
	}
	paramMap := map[string]string{
		"blah": "stuff",
		"boo":  "ahhh",
	}

	enc := GenerateEnclave(cfg, "example.com", paramMap)
	assert.Equal(t, "foobar", enc.DefaultProvider)
	assert.Equal(t, IacTypeCf, enc.DefaultIac)
	assert.Equal(t, "example.com", enc.DomainName)
	assert.Equal(t, paramMap, enc.Parameters)
	assert.True(t, *enc.EnvBehavior.DeployAllStacks)
	assert.Len(t, enc.Providers, 1)
	prov := enc.Providers[0]
	assert.Equal(t, "foobar", prov.Name)
	assert.Equal(t, "1234567890", prov.AccountId)
	assert.Equal(t, "us-west-2", prov.Region)
	assert.Equal(t, EncProvAws, prov.ProviderId)
}
