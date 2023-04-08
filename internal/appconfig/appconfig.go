package appconfig

type AppConfig struct {
	Version       int    `yaml:"version"`
	EnvDefinition string `yaml:"envDef"`
}

func Parse() {

}

func Generate() {

}

func Fetch() {
	// Decode base64 into file or URI
}
