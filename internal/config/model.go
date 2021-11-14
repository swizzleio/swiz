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

// Bastion defines each bastion host
type Bastion struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Addr        string   `json:"addr"`
	BastionAuth AuthInfo `json:"bastion_auth"`
	Signature   string   `json:"signature"`
}

// Config contains the main application configuration
type Config struct {
	Version       int                       `json:"version"` // Allows for future upgradability
	Os            string                    `json:"-"`
	BastionHosts  []Bastion                 `json:"bastion_hosts"`
	LaunchProfile map[string]*LaunchProfile `json:"launch_profile"`
}

// GetBastionById returns a bastion host by id
func (c Config) GetBastionById(id string) *Bastion {
	for _, bastion := range c.BastionHosts {
		if bastion.Id == id {
			return &bastion
		}
	}

	return nil
}

// GetBastionByAddr returns a bastion host by address
func (c Config) GetBastionByAddr(addr string) *Bastion {
	for _, bastion := range c.BastionHosts {
		if bastion.Addr == addr {
			return &bastion
		}
	}

	return nil
}

// GetBastionAddrList returns a list of the bastion host addresses
func (c Config) GetBastionAddrList() []string {
	addrList := []string{}
	for _, bastion := range c.BastionHosts {
		addrList = append(addrList, bastion.Addr)
	}

	return addrList
}
