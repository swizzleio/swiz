package model

type EnvBehavior struct {
	NoUpdateDeploy bool `yaml:"no_update_deploy"`
	NoOrphanDelete bool `yaml:"no_orphan_delete"`
	FastDelete     bool `yaml:"fast_delete"`
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
	Providers       []EncProvider     `yaml:"providers"`
	EnvBehavior     EnvBehavior       `yaml:"env_behavior"`
	DomainName      string            `yaml:"domain_name"`
	Parameters      map[string]string `yaml:"params"`
}
