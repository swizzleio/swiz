package model

import "github.com/swizzleio/swiz/pkg/drivers/awswrap"

type EnvBehavior struct {
	NoUpdateDeploy *bool `yaml:"no_update_deploy"`
	NoOrphanDelete *bool `yaml:"no_orphan_delete"`
	FastDelete     *bool `yaml:"fast_delete"`
}

type EncProvider struct {
	Name       string `yaml:"name"`
	ProviderId string `yaml:"provider_id"`
	AccountId  string `yaml:"account_id"`
	Region     string `yaml:"region"`
}

type Enclave struct {
	Name            string            `yaml:"name"`
	DefaultProvider string            `yaml:"default_provider"`
	DefaultIac      string            `yaml:"default_iac"`
	Providers       []EncProvider     `yaml:"providers"`
	EnvBehavior     EnvBehavior       `yaml:"env_behavior"`
	DomainName      string            `yaml:"domain_name"`
	Parameters      map[string]string `yaml:"params"`
}

func (e Enclave) GetProvider(providerName string) *EncProvider {
	if providerName == "" {
		providerName = e.DefaultProvider
	}

	for _, p := range e.Providers {
		if p.Name == providerName {
			return &p
		}
	}
	return nil
}

func (e EncProvider) ToAwsConfig() awswrap.AwsConfig {
	return awswrap.AwsConfig{
		Profile:   e.Name,
		AccountId: e.AccountId,
		Region:    e.Region,
	}
}
