package xprometheus

type Config struct {
	Namespace string
	Subsystem string
	Path      string
	Addr      string
}

func DefaultConfig() *Config {
	return &Config{
		Namespace: "default",
		Subsystem: "default",
		Path:      "/metrics",
		Addr:      ":8888",
	}
}
