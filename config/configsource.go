package config

var (
	configsourceBuilders = make(map[string]ConfigSource)
)

type ConfigSource interface {
	Read() error
	Init(protocol string) error
}

func Register(scheme string, creator ConfigSource) {
	configsourceBuilders[scheme] = creator
}
