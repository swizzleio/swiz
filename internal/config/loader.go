package config

import "log"

// LoadConfig loads the application config file and sets any embedded config items
func LoadConfig() (Config, error) {
	cfg.Version = "0.0"

	return GetConfig(), nil
}

// LoadConfigMust loads the config and panics if it could not be loaded
func LoadConfigMust() Config {
	cfg, err := LoadConfig()

	if err != nil {
		log.Fatalf("loading config. %v", err)
	}

	return cfg
}
