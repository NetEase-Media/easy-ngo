package config

var (
	configsourceBuilders = make(map[string]ConfigSource)
)

type ConfigSource interface {
	Read(config *Config) error
	Init(protocol string, config *Config) error
}

func Register(scheme string, creator ConfigSource) {
	configsourceBuilders[scheme] = creator
}
