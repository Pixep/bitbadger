package bitbadger

// Config holds server configuration
// to authenticate to the upstream repository
type Config struct {
	Username string
	Password string
}

var config Config

// SetConfig sets the global configuration
func SetConfig(conf Config) {
	config = conf
}

// GetConfig gets the global configuration
func GetConfig() Config {
	return config
}
