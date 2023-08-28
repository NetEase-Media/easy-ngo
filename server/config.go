package server

type MODE string

const (
	DEBUG   MODE = "debug"
	RELEASE      = "release"
	TEST         = "test"
)

type Config struct {
	Host          string
	Port          int
	EnabledTracer bool
	Mode          MODE
	Metrics       Metrics
}

type Metrics struct {
	Enabled bool
	Bucket  Bucket
	Path    Path
}

type Path struct {
	Include Include
	Exclude Exclude
}

type Include struct {
	Prefix []string
	Regex  []string
}

type Exclude struct {
	Prefix []string
	Regex  []string
}

type Bucket struct {
	Start, Factor float64
	Count         int
}

func DefaultConfig() *Config {
	return &Config{
		Host: "0.0.0.0",
		Port: 8080,
		Metrics: Metrics{
			Enabled: false,
		},
		EnabledTracer: false,
		Mode:          DEBUG,
	}
}
