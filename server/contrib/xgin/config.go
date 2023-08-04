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
	MetricsPath   string
	Mode          MODE
}

func DefaultConfig() *Config {
	return &Config{
		Host:          "0.0.0.0",
		Port:          8080,
		EnabledMetric: true,
		EnabledTrace:  false,
		MetricsPath:   "/metrics",
		Mode:          DEBUG,
	}
}
