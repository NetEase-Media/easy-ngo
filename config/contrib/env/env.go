package env

import (
	"strings"

	"github.com/NetEase-Media/easy-ngo/config"
)

const EnvConfigSourceName = "env"

type EnvConfigSource struct {
	envPrefix string
	bindEnv   []string
}

func New() *EnvConfigSource {
	return &EnvConfigSource{
		bindEnv: make([]string, 0),
	}
}

func (e *EnvConfigSource) Init(protocol string, config *config.Config) error {
	kvs := strings.Split(protocol, ";")
	for _, kv := range kvs {
		if strings.HasPrefix(kv, "prefix=") {
			e.envPrefix = strings.TrimPrefix(kv, "prefix=")
		} else if strings.HasPrefix(kv, "envName=") {
			e.bindEnv = strings.Split(strings.TrimPrefix(kv, "envName="), ",")
		}
	}
	config.Viper.AutomaticEnv()
	config.Viper.SetEnvPrefix(e.envPrefix)
	config.Viper.BindEnv(e.bindEnv...)
	return nil
}

func (e *EnvConfigSource) Read(config *config.Config) error {
	return nil
}
