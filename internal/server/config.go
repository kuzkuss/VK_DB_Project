package server

type Config struct {
	ServerAddr string `toml: "server_addr"`
	LogLevel string `toml: "log_level"`
}


func NewConfig() *Config {
	return &Config{
		ServerAddr: ":8000",
		LogLevel: "debug",
	}
}

