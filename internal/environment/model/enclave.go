package model

type EnvBehavior struct {
	NoUpdateDeploy bool `yaml:"no_update_deploy"`
	NoOrphanDelete bool `yaml:"no_orphan_delete"`
	FastDelete     bool `yaml:"fast_delete"`
}

type Enclave struct {
	Name        string      `yaml:"name"`
	ProviderId  string      `yaml:"provider_id"`
	AccountId   string      `yaml:"account_id"`
	Region      string      `yaml:"region"`
	EnvBehavior EnvBehavior `yaml:"env_behavior"`
}
