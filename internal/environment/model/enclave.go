package model

type Enclave struct {
	Name       string `yaml:"name"`
	ProviderId string `yaml:"provider_id"`
	AccountId  string `yaml:"account_id"`
	Region     string `yaml:"region"`
}
