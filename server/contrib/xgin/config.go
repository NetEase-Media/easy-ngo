package xgin

import "github.com/NetEase-Media/easy-ngo/xmetrics"

type MODE string

const (
	DEBUG   MODE = "debug"
	RELEASE      = "release"
	TEST         = "test"
)

type Config struct {
	Host           string
	Port           int
	EnabledMetrics bool
	EnabledTracer  bool
	Mode           MODE
	Metrics        Metrics
}

type Metrics struct {
	Bucket           xmetrics.Bucket
	ExcludeByPrefix  []string
	ExcludeByRegular []string
	IncludeByPrefix  []string
	IncludeByRegular []string
}

func DefaultConfig() *Config {
	return &Config{
		Host:           "0.0.0.0",
		Port:           8080,
		EnabledMetrics: false,
		EnabledTracer:  false,
		Mode:           DEBUG,
	}
}
