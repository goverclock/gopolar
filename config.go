package gopolar

type Config struct {
	tunnels []Tunnel
}

// read config from $HOME/.config/gopolar/gopolar.toml
// create config file if not exist
func NewConfig() *Config {
	return &Config{}
}
