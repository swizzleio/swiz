package model

import "github.com/swizzleio/swiz/pkg/drivers/awswrap"

const (
	EncProvDummy = "DUMMY"
	EncProvAws   = "AWS"
	IacTypeDummy = "Dummy"
	IacTypeCf    = "Cloudformation"
)

type EnvBehavior struct {
	NoUpdateDeploy  *bool `yaml:"no_update_deploy,omitempty"`
	NoOrphanDelete  *bool `yaml:"no_orphan_delete,omitempty"`
	DeployAllStacks *bool `yaml:"deploy_all_stacks,omitempty"`
	FastDelete      *bool `yaml:"fast_delete,omitempty"`
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

func (e EncProvider) ToAwsConfig() awswrap.AwsConfiger {
	return awswrap.NewAwsConfig(e.Name, e.AccountId, e.Region)
}

func GenerateEnclave(config awswrap.AwsConfig, domainName string, params map[string]string) Enclave {
	deployAllStacks := true

	return Enclave{
		Name:            "",
		DefaultProvider: config.Profile,
		DefaultIac:      IacTypeCf,
		Providers: []EncProvider{
			{
				Name:       config.Profile,
				AccountId:  config.AccountId,
				Region:     config.Region,
				ProviderId: EncProvAws,
			},
		},
		EnvBehavior: EnvBehavior{
			DeployAllStacks: &deployAllStacks,
		},
		DomainName: domainName,
		Parameters: params,
	}
}
