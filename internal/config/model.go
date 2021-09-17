package config

// AuthInfo provides info about the remote auth mechanism
type AuthInfo struct {
	KeyFilename string `json:"key_filename"`
}

// LaunchProfile defines how a remote application will be launched
type LaunchProfile struct {
	DefaultApp  string   `json:"default_app"`
	DefaultUser string   `json:"default_user"`
	DefaultPort int      `json:"default_port"`
	Auth        AuthInfo `json:"auth"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
}

// Config contains the main application configuration
type Config struct {
	Version       int                       `json:"version"` // Allows for future upgradability
	Os            string                    `json:"-"`
	BastionAddrs  []string                  `json:"bastion_addrs"`
	BastionAuth   AuthInfo                  `json:"bastion_auth"`
	LaunchProfile map[string]*LaunchProfile `json:"launch_profile"`
}
