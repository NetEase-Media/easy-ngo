package xecho

type Config struct {
	Host string
	Port int
}

func DefaultConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 8080,
	}
}

func EnvConfig() *Config {
	config := FileConfig()
	return config
}

func FileConfig() *Config {
	config := DefaultConfig()
	return config
}
