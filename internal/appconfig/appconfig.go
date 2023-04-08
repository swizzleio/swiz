package appconfig

type AppConfig struct {
	Version       int `yaml:"version"`
	EnvDefinition []struct {
		Name       string `yaml:"name"`
		EnvDefFile string `yaml:"env_def_file"`
		Default    bool   `yaml:"default"`
	} `yaml:"env_def"`
	EnclaveDefinition []struct {
		Name       string `yaml:"name"`
		ProviderId string `yaml:"provider_id"`
		AccountId  string `yaml:"account_id"`
		Region     string `yaml:"region"`
		DomainName string `yaml:"domain_name"`
	} `yaml:"enclave_def"`
}

func Parse() {

}

func Generate() {

}

func Fetch() {
	// Decode base64 into file or URI
}
