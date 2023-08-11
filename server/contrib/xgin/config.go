package xgin

type MODE string

const (
	DEBUG   MODE = "debug"
	RELEASE      = "release"
	TEST         = "test"
)

type Config struct {
	Host          string
	Port          int
	EnabledMetric bool
	EnabledTrace  bool
	Mode          MODE
}

func DefaultConfig() *Config {
	return &Config{
		Host:          "0.0.0.0",
		Port:          8080,
		EnabledMetric: false,
		EnabledTrace:  false,
		Mode:          DEBUG,
	}
}
