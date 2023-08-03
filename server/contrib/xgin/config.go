package xgin

type Config struct {
	Host          string
	Port          int
	EnabledMetric bool
	EnabledTracer bool
	MetricsPath   string
}

func DefaultConfig() *Config {
	return &Config{
		Host:          "0.0.0.0",
		Port:          8080,
		EnabledMetric: true,
		EnabledTracer: false,
		MetricsPath:   "/metrics",
	}
}
